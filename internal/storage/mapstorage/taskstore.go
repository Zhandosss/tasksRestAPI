package mapstorage

import (
	"errors"
	"fmt"
	"restAPI/internal/model"
	"sync"
	"time"
)

var (
	ErrTaskStoreIsEmpty = errors.New("task store is empty")
)

type TaskStore struct {
	sync.Mutex
	tasks  map[int]*model.Task
	nextId int
}

func NewTaskStore() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[int]*model.Task)
	ts.nextId = 1
	return ts

}

func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) (int, error) { // +
	ts.Lock()
	defer ts.Unlock()

	t := &model.Task{
		Id:   ts.nextId,
		Text: text,
		Due:  due,
	}
	t.Tags = make([]string, len(tags))
	copy(t.Tags, tags)
	ts.tasks[ts.nextId] = t
	if _, ok := ts.tasks[ts.nextId]; !ok {
		err := errors.New("failed to create task")
		return 0, fmt.Errorf("%w", err)
	}
	ts.nextId++
	return t.Id, nil
}

func (ts *TaskStore) GetTask(id int) (model.Task, error) {
	ts.Lock()
	defer ts.Unlock()
	if t, ok := ts.tasks[id]; !ok {
		err := fmt.Errorf("task with id=%d not found", id)
		return model.Task{}, fmt.Errorf("%w", err)
	} else {
		return *t, nil
	}
}

func (ts *TaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()
	if _, ok := ts.tasks[id]; !ok {
		err := fmt.Errorf("task with id=%d not found", id)
		return fmt.Errorf("%w", err)
	}
	delete(ts.tasks, id)
	return nil
}

func (ts *TaskStore) DeleteAllTasks() error { //+
	ts.Lock()
	defer ts.Unlock()
	if len(ts.tasks) == 0 {
		return fmt.Errorf("%w", ErrTaskStoreIsEmpty)
	}
	ts.tasks = make(map[int]*model.Task)
	ts.nextId = 1
	return nil
}

func (ts *TaskStore) GetAllTasks() ([]model.Task, error) { //+
	ts.Lock()
	defer ts.Unlock()
	ans := make([]model.Task, 0, len(ts.tasks))
	for _, t := range ts.tasks {
		ans = append(ans, *t)
	}
	if len(ans) == 0 {
		return []model.Task{}, fmt.Errorf("%w", ErrTaskStoreIsEmpty)
	}
	return ans, nil
}

func (ts *TaskStore) GetTasksByTag(tag string) ([]model.Task, error) {
	ts.Lock()
	defer ts.Unlock()
	ans := make([]model.Task, 0)
	for _, t := range ts.tasks {
		for _, tagOfT := range t.Tags {
			if tag == tagOfT {
				ans = append(ans, *t)
				break
			}
		}
	}
	if len(ans) == 0 {
		err := errors.New("dont found tasks by tag")
		return []model.Task{}, fmt.Errorf("%w", err)
	}
	return ans, nil
}

func (ts *TaskStore) GetTasksByDueDate(year int, month time.Month, day int) ([]model.Task, error) {
	ts.Lock()
	defer ts.Unlock()
	ans := make([]model.Task, 0)
	for _, t := range ts.tasks {
		y, m, d := t.Due.Date()
		if y == year && month == m && d == day {
			ans = append(ans, *t)
		}
	}
	if len(ans) == 0 {
		err := errors.New("dont found tasks by date")
		return []model.Task{}, fmt.Errorf("%w", err)
	}
	return ans, nil
}
