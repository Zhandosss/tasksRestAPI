package service

import (
	"fmt"
	"restAPI/internal/model"
	"restAPI/internal/repositories"
	"time"
)

type TaskService struct {
	rep repositories.Task
}

func NewTaskService(rep repositories.Task) *TaskService {
	return &TaskService{
		rep: rep,
	}
}

func (s *TaskService) CreateTask(text string, tags []string, date time.Time, ownerID int64) (int64, error) {
	id, err := s.rep.CreateTask(text, tags, date, ownerID)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	return id, nil
}

func (s *TaskService) GetTask(taskID, userID int64) (model.Task, error) {
	task, err := s.rep.GetTask(taskID, userID)
	if err != nil {
		return model.Task{}, fmt.Errorf("%w", err)
	}
	return task, nil
}

func (s *TaskService) DeleteTask(taskID, userID int64) error {
	err := s.rep.DeleteTask(taskID, userID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (s *TaskService) DeleteAllTasks() error {
	err := s.rep.DeleteAllTasks()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (s *TaskService) GetAllTasks() ([]model.Task, error) {
	tasks, err := s.rep.GetAllTasks()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return tasks, nil
}

func (s *TaskService) GetTasksByDate(day, month, year int, userID int64) ([]model.Task, error) {
	tasks, err := s.rep.GetTasksByDate(day, month, year, userID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return tasks, nil
}

func (s *TaskService) GetTasksByTag(tag string, userID int64) ([]model.Task, error) {
	tasks, err := s.rep.GetTasksByTag(tag, userID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return tasks, nil
}
