package repositories

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"restAPI/internal/entities"
	"restAPI/internal/model"
)

type TaskPostgres struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewTaskPostgres(db *sqlx.DB, log *slog.Logger) *TaskPostgres {
	return &TaskPostgres{
		db:  db,
		log: log,
	}
}

func (r *TaskPostgres) getOrCreateTagID(tag string) (int64, error) {
	op := "getTagID"
	tagID := make([]int64, 0, 1)
	query := "SELECT id FROM tags WHERE tag = $1"
	err := r.db.Select(&tagID, query, tag)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if len(tagID) == 1 {
		return tagID[0], nil
	}
	id, err := r.insertTag(tag)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (r *TaskPostgres) insertTag(tag string) (int64, error) {
	op := "insertTag"
	var tagID int64
	query := "INSERT INTO tags (tag) VALUES ($1) RETURNING id"
	err := r.db.Get(&tagID, query, tag)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return tagID, nil
}

func (r *TaskPostgres) insertInTagInTask(taskID, tagID int64) error {
	op := "insertInTagInTask"
	query := "INSERT INTO tags_in_task (tag_id, task_id) VALUES ($1, $2)"
	_, err := r.db.Exec(query, tagID, taskID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *TaskPostgres) CreateTask(task model.Task) (int64, error) {
	op := "CreateTask"
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	var taskID int64
	query := "INSERT INTO tasks (task, date, owner_id) VALUES ($1, $2, $3) RETURNING id"
	err = r.db.Get(&taskID, query, task.Text, task.Date, task.OwnerID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	for _, tag := range task.Tags {
		tagID, err := r.getOrCreateTagID(tag)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		err = r.insertInTagInTask(taskID, tagID)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}

	}
	if tx.Commit() != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return taskID, nil
}

func (r *TaskPostgres) GetTask(taskID, userID int64) (model.Task, error) {
	op := "GetTask"
	tx, err := r.db.Begin()
	if err != nil {
		return model.Task{}, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	query := "SELECT id, task, date, owner_id FROM tasks WHERE id = $1 AND owner_id = $2"
	task := make([]model.Task, 0, 1)
	err = r.db.Select(&task, query, taskID, userID)
	if err != nil {
		return model.Task{}, fmt.Errorf("%s: %w", op, err)
	}
	if len(task) == 0 {
		return model.Task{}, ErrNoTask
	}
	tags := make([]string, 0)
	query = `SELECT tag FROM tags
		     JOIN tags_in_task 
		         ON tags_in_task.tag_id = tags.id 
             WHERE tags_in_task.task_id = $1`
	err = r.db.Select(&tags, query, taskID)
	task[0].Tags = tags
	if tx.Commit() != nil {
		return model.Task{}, fmt.Errorf("%s: %w", op, err)
	}
	return task[0], nil
}

func (r *TaskPostgres) DeleteTask(taskID, userID int64) error {
	op := "DeleteTask"
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	query := "DELETE FROM tasks WHERE id = $1 AND owner_id = $2"
	res, err := r.db.Exec(query, taskID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, ErrNoTask)
	}
	query = "DELETE FROM tags_in_task WHERE task_id = $1"
	res, err = r.db.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tx.Commit() != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *TaskPostgres) DeleteAllByUser(userID int64) error {
	op := "DeleteAllByUser"
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	query := "SELECT id FROM tasks WHERE owner_id = $1"
	tasks := make([]int64, 0)
	err = r.db.Select(&tasks, query, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if len(tasks) == 0 {
		if tx.Commit() != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	}
	for _, taskID := range tasks {
		err = r.DeleteTask(taskID, userID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	if tx.Commit() != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *TaskPostgres) DeleteAllTasks() error {
	op := "DeleteAllTasks"
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	query := "DELETE FROM tasks"
	res, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, ErrEmptyTable)
	}
	query = "TRUNCATE TABLE tags_in_task"
	res, err = r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	query = "TRUNCATE TABLE tags"
	res, err = r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tx.Commit() != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *TaskPostgres) uniteTasks(rawTasks []entities.TaskWithTag) []model.Task {
	taskMap := make(map[int64]*model.Task)
	for _, rawTask := range rawTasks {
		_, ok := taskMap[rawTask.ID]
		if !ok {
			taskMap[rawTask.ID] = &model.Task{
				ID:      rawTask.ID,
				Text:    rawTask.Task,
				Date:    rawTask.Date,
				OwnerID: rawTask.OwnerID,
			}
		}
		if rawTask.Tag == nil {
			continue
		}
		taskMap[rawTask.ID].Tags = append(taskMap[rawTask.ID].Tags, *rawTask.Tag)
	}
	tasks := make([]model.Task, 0, len(taskMap))
	for _, task := range taskMap {
		tasks = append(tasks, *task)
	}
	return tasks

}

func (r *TaskPostgres) GetAllTasks() ([]model.Task, error) {
	op := "GetAllTasks"
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	rawTasks := make([]entities.TaskWithTag, 0)
	query := `SELECT tasks.id, task, date, tags.tag AS tag, owner_id FROM tasks 
              LEFT OUTER JOIN tags_in_task 
                  ON tasks.id = tags_in_task.task_id 
    		  LEFT OUTER JOIN tags  
    		      ON tags.id = tags_in_task.tag_id`
	err = r.db.Select(&rawTasks, query)
	r.log.Debug("rawTasks", slog.Any("string", rawTasks))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(rawTasks) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrEmptyTable)
	}
	tasks := r.uniteTasks(rawTasks)
	if tx.Commit() != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tasks, nil
}

func (r *TaskPostgres) GetAllByUser(userID int64) ([]model.Task, error) {
	op := "GetAllByUser"
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	rawTasks := make([]entities.TaskWithTag, 0)
	query := `SELECT tasks.id, task, date, tags.tag AS tag, owner_id FROM tasks
			  LEFT OUTER JOIN tags_in_task
 				ON tasks.id = tags_in_task.task_id
			  LEFT OUTER JOIN tags
			  	ON tags.id = tags_in_task.tag_id
			  WHERE owner_id = $1`
	err = r.db.Select(&rawTasks, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(rawTasks) == 0 {
		if tx.Commit() != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return []model.Task{}, nil
	}
	tasks := r.uniteTasks(rawTasks)
	if tx.Commit() != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tasks, nil
}

func (r *TaskPostgres) GetTasksByDate(day, month, year int, userID int64) ([]model.Task, error) {
	op := "GetTasksByDate"
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	rawTasks := make([]entities.TaskWithTag, 0)
	query := `SELECT tasks.id, task, date, tags.tag, owner_id AS tag FROM tasks 
              LEFT OUTER JOIN tags_in_task 
                  ON tasks.id = tags_in_task.task_id 
    		  LEFT OUTER JOIN tags  
    		      ON tags.id = tags_in_task.tag_id
    		  WHERE EXTRACT(YEAR FROM date) = $1 AND EXTRACT(MONTH FROM date) = $2 AND EXTRACT(DAY FROM date) = $3 AND owner_id = $4`
	err = r.db.Select(&rawTasks, query, year, month, day, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(rawTasks) == 0 {
		if tx.Commit() != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return []model.Task{}, nil
	}
	tasks := r.uniteTasks(rawTasks)
	if tx.Commit() != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tasks, nil
}

func (r *TaskPostgres) GetTasksByTag(tag string, userID int64) ([]model.Task, error) {
	op := "GetTasksByTag"
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	rawTasks := make([]entities.TaskWithTag, 0)
	query := `SELECT tasks.id, task, date, tags.tag, owner_id AS tag FROM tasks 
              LEFT OUTER JOIN tags_in_task 
                  ON tasks.id = tags_in_task.task_id 
    		  LEFT OUTER JOIN tags  
    		      ON tags.id = tags_in_task.tag_id
			  WHERE tag = $1 AND owner_id = $2`
	err = r.db.Select(&rawTasks, query, tag, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(rawTasks) == 0 {
		if tx.Commit() != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return []model.Task{}, nil
	}
	tasks := r.uniteTasks(rawTasks)
	if tx.Commit() != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tasks, nil

}

func (r *TaskPostgres) deleteFromTagsInTask(taskID int64, tagsToDelete []string) error {
	op := "tagUpdate"
	query := fmt.Sprintf("DELETE FROM tags_in_task WHERE task_id = %d AND tag_id IN (?)", taskID)
	query, args, err := sqlx.In(query, tagsToDelete)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	query = r.db.Rebind(query)
	res, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, errors.New("can't delete tags"))
	}
	return nil
}

func (r *TaskPostgres) tagUpdate(taskID int64, tags []string) error {
	op := "tagUpdate"
	newTagsMap := make(map[string]struct{})
	for _, tag := range tags {
		newTagsMap[tag] = struct{}{}
	}
	oldTags := make([]string, 0)
	query := `SELECT tags.tag, tags.id FROM tags_in_task
			 LEFT JOIN tags on tags_in_task.tag_id = tags.id
			 WHERE tags_in_task.task_id = $1`
	err := r.db.Select(&oldTags, query, taskID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	tagsToDelete := make([]string, len(oldTags))
	for _, tag := range oldTags {
		if _, ok := newTagsMap[tag]; ok {
			delete(newTagsMap, tag)
		} else {
			tagsToDelete = append(tagsToDelete, tag)
		}
	}
	if err = r.deleteFromTagsInTask(taskID, tagsToDelete); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	for tag := range newTagsMap {
		tagID, err := r.getOrCreateTagID(tag)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		err = r.insertInTagInTask(taskID, tagID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}

func (r *TaskPostgres) UpdateTask(taskID, userID int64, text string, tags []string) error {
	op := "Update"
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	query := `UPDATE tasks
			  SET task = $1
			  WHERE id = $2 AND owner_id = $3`
	res, err := r.db.Exec(query, text, taskID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, ErrNoTask)
	}
	err = r.tagUpdate(taskID, tags)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tx.Commit() != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
