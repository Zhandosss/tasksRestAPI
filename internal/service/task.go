package service

import (
	"fmt"
	"restAPI/internal/model"
	"restAPI/internal/repositories"
)

type TaskService struct {
	rep repositories.Task
}

func NewTaskService(rep repositories.Task) *TaskService {
	return &TaskService{
		rep: rep,
	}
}

func (s *TaskService) CreateTask(task model.Task) (int64, error) {
	id, err := s.rep.CreateTask(task)
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

func (s *TaskService) GetAllByUser(userID int64) ([]model.Task, error) {
	tasks, err := s.rep.GetAllByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return tasks, nil
}

func (s *TaskService) DeleteTask(taskID, userID int64) error {
	err := s.rep.DeleteTask(taskID, userID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (s *TaskService) DeleteAllByUser(userID int64) error {
	err := s.rep.DeleteAllByUser(userID)
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

func (s *TaskService) UpdateTask(taskID, userID int64, Text string, Tags []string) error {
	err := s.rep.UpdateTask(taskID, userID, Text, Tags)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil

}
