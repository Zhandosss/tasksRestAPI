package repositories

import (
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"restAPI/internal/model"
)

var (
	ErrNoTask           = errors.New("task not found")
	ErrEmptyTable       = errors.New("table is empty")
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrNoSuchUser       = errors.New("there is no such user")
	ErrTwoSameLoginInDb = errors.New("there is two same user logins in db")
	ErrWrongPassword    = errors.New("wrong password")
)

type Task interface {
	DeleteAllTasks() error
	GetAllTasks() ([]model.Task, error)
	CreateTask(task model.Task) (int64, error)
	DeleteTask(taskID, userID int64) error
	DeleteAllByUser(userID int64) error
	GetTask(taskID, userID int64) (model.Task, error)
	GetAllByUser(userID int64) ([]model.Task, error)
	GetTasksByDate(day, month, year int, userID int64) ([]model.Task, error)
	GetTasksByTag(tag string, userID int64) ([]model.Task, error)
	UpdateTask(taskID, userID int64, Text string, Tags []string) error
}

type Authorization interface {
	CreateUser(user model.User) (int64, error)
	GetUser(login, password string) (model.User, error)
}

type Repository struct {
	Task
	Authorization
}

func New(db *sqlx.DB, log *slog.Logger) *Repository {
	return &Repository{
		Task:          NewTaskPostgres(db, log),
		Authorization: NewAuthPostgres(db, log),
	}
}
