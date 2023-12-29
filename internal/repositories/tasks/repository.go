package tasks

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"restAPI/internal/entities"
	"restAPI/internal/model"
	"time"
)

var (
	ErrNoTask        = errors.New("task not found")
	ErrEmptyTable    = errors.New("table is empty")
	ErrNoTasksByDate = errors.New("no tasks found by this date")
	ErrNoTaskByTag   = errors.New("no tasks found by this tag")
)

type Repository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewRepository(db *sqlx.DB, log *slog.Logger) *Repository {
	return &Repository{
		db:  db,
		log: log,
	}
}

func (r *Repository) getOrCreateTagID(tag string) (int64, error) {
	op := "getTagID"
	tagID := make([]int64, 0, 1)
	query := "SELECT tag_id FROM tags WHERE tag = $1"
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

func (r *Repository) insertTag(tag string) (int64, error) {
	op := "insertTag"
	var tagID int64
	query := "INSERT INTO tags (tag) VALUES ($1) RETURNING  tag_id"
	err := r.db.Get(&tagID, query, tag)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return tagID, nil
}

func (r *Repository) insertInTagInTask(taskID, tagID int64) error {
	op := "insertInTagInTask"
	query := "INSERT INTO tags_in_task (tag_id, task_id) VALUES ($1, $2)"
	_, err := r.db.Exec(query, tagID, taskID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *Repository) CreateTask(text string, tags []string, date time.Time) (int64, error) {
	op := "CreateTask"
	var taskID int64
	query := "INSERT INTO tasks (task, date) VALUES ($1, $2) RETURNING task_id"
	err := r.db.Get(&taskID, query, text, date)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	for _, tag := range tags {
		tagID, err := r.getOrCreateTagID(tag)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		err = r.insertInTagInTask(taskID, tagID)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}

	}
	return taskID, nil
}

func (r *Repository) GetTask(taskID int64) (model.Task, error) {
	op := "GetTask"
	query := "SELECT task_id, task, date FROM tasks WHERE task_id = $1"
	task := make([]model.Task, 0, 1)
	err := r.db.Select(&task, query, taskID)
	if len(task) == 0 {
		return model.Task{}, ErrNoTask
	}
	if err != nil {
		return model.Task{}, fmt.Errorf("%s: %w", op, err)
	}
	tags := make([]string, 0)
	query = "SELECT tag FROM tags " +
		"JOIN tags_in_task ON tags_in_task.tag_id = tags.tag_id WHERE tags_in_task.task_id = $1"
	err = r.db.Select(&tags, query, taskID)
	task[0].Tags = tags
	return task[0], nil
}

func (r *Repository) checkAndDelete(tagID int64) error {
	op := "checkAndDelete"
	query := "SELECT tag_id FROM tags_in_task WHERE tag_id = $1"
	tags := make([]int64, 0)
	err := r.db.Select(&tags, query, tagID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if len(tags) != 0 {
		return nil
	}
	query = "DELETE FROM tags WHERE tag_id = $1"
	_, err = r.db.Exec(query, tagID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil

}

func (r *Repository) DeleteTask(taskID int64) error {
	op := "DeleteTask"
	query := "DELETE FROM tasks WHERE task_id = $1"
	res, err := r.db.Exec(query, taskID)
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
	query = "SELECT tag_id FROM tags_in_task WHERE task_id = $1"
	tagsID := make([]int64, 0)
	err = r.db.Select(&tagsID, query, taskID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if len(tagsID) == 0 {
		return nil
	}
	query = "DELETE FROM tags_in_task WHERE task_id = $1"
	res, err = r.db.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	for _, tagID := range tagsID {
		err = r.checkAndDelete(tagID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}

func (r *Repository) DeleteAllTasks() error {
	op := "DeleteAllTasks"
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
	return nil
}

func (r *Repository) uniteTasks(rawTasks []entities.TaskWithTag) []model.Task {
	taskMap := make(map[int64]*model.Task)
	for _, rawTask := range rawTasks {
		_, ok := taskMap[rawTask.ID]
		if !ok {
			taskMap[rawTask.ID] = &model.Task{
				ID:   rawTask.ID,
				Text: rawTask.Task,
				Date: rawTask.Date,
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

func (r *Repository) GetAllTasks() ([]model.Task, error) {
	op := "GetAllTasks"
	rawTasks := make([]entities.TaskWithTag, 0)
	query := `SELECT tasks.task_id, task, date, tags.tag AS tag FROM tasks 
              LEFT OUTER JOIN tags_in_task 
                  ON tasks.task_id = tags_in_task.task_id 
    		  LEFT OUTER JOIN tags  
    		      ON tags.tag_id = tags_in_task.tag_id`
	err := r.db.Select(&rawTasks, query)
	r.log.Debug("rawTasks", slog.Any("string", rawTasks))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(rawTasks) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrEmptyTable)
	}
	tasks := r.uniteTasks(rawTasks)
	return tasks, nil
}

func (r *Repository) GetTasksByDate(day, month, year int) ([]model.Task, error) {
	op := "GetTasksByDate"
	rawTasks := make([]entities.TaskWithTag, 0)
	query := `SELECT tasks.task_id, task, date, tags.tag AS tag FROM tasks 
              LEFT OUTER JOIN tags_in_task 
                  ON tasks.task_id = tags_in_task.task_id 
    		  LEFT OUTER JOIN tags  
    		      ON tags.tag_id = tags_in_task.tag_id
    		  WHERE EXTRACT(YEAR FROM date) = $1 AND EXTRACT(MONTH FROM date) = $2 AND EXTRACT(DAY FROM date) = $3`
	err := r.db.Select(&rawTasks, query, year, month, day)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(rawTasks) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrNoTasksByDate)
	}
	tasks := r.uniteTasks(rawTasks)
	return tasks, nil
}

func (r *Repository) GetTasksByTag(tag string) ([]model.Task, error) {
	op := "GetTasksByTag"
	rawTasks := make([]entities.TaskWithTag, 0)
	query := `SELECT tasks.task_id, task, date, tags.tag AS tag FROM tasks 
              LEFT OUTER JOIN tags_in_task 
                  ON tasks.task_id = tags_in_task.task_id 
    		  LEFT OUTER JOIN tags  
    		      ON tags.tag_id = tags_in_task.tag_id
			  WHERE tag = $1`
	err := r.db.Select(&rawTasks, query, tag)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(rawTasks) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrNoTaskByTag)
	}
	tasks := r.uniteTasks(rawTasks)
	return tasks, nil

}
