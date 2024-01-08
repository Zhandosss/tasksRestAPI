package service

import (
	"errors"
	"restAPI/internal/model"
	"restAPI/internal/repositories"
	"time"
)

var (
	ErrWrongPassword = errors.New("wrong password")
	ErrTokenClaims   = errors.New("token claims are not of type *tokenClaims")
)

type Task interface {
	DeleteAllTasks() error
	GetAllTasks() ([]model.Task, error)
	CreateTask(text string, tags []string, date time.Time, ownerID int64) (int64, error)
	DeleteTask(taskID, userID int64) error
	GetTask(taskID, userID int64) (model.Task, error)
	GetTasksByDate(day, month, year int, userID int64) ([]model.Task, error)
	GetTasksByTag(tag string, userID int64) ([]model.Task, error)
}

type Authorization interface {
	CreateUser(user model.User) (int64, error)
	GenerateToken(login, password string) (string, error)
	ParseToken(inputToken string) (int64, error)
}

type Service struct {
	Task
	Authorization
}

func New(rep *repositories.Repository) *Service {
	return &Service{
		Task:          NewTaskService(rep.Task),
		Authorization: NewAuthService(rep.Authorization),
	}
}
