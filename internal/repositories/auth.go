package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"restAPI/internal/model"
)

type AuthPostgres struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewAuthPostgres(db *sqlx.DB, log *slog.Logger) *AuthPostgres {
	return &AuthPostgres{
		db:  db,
		log: log,
	}
}

func (r *AuthPostgres) validateNewUser(login string) (bool, error) {
	op := "validateLogin"
	query := `SELECT login FROM users WHERE login = $1`
	users := make([]string, 0)
	err := r.db.Select(&users, query, login)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if len(users) != 0 {
		return false, nil
	}
	return true, nil
}

func (r *AuthPostgres) CreateUser(user model.User) (int64, error) {
	op := "CreateUser"
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	var userID int64
	ok, err := r.validateNewUser(user.Login)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if !ok {
		return 0, fmt.Errorf("%s: %w", op, ErrUserAlreadyExist)
	}
	query := `INSERT INTO users (firstname, secondname, login, password_hash) 
		      VALUES ($1, $2, $3, $4) RETURNING id`
	err = r.db.Get(&userID, query, user.FirstName, user.SecondName, user.Login, user.Password)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if tx.Commit() != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return userID, nil
}

func (r *AuthPostgres) GetUser(login, password string) (model.User, error) {
	op := "GetUser"
	tx, err := r.db.Begin()
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	user := make([]model.User, 0)
	query := `SELECT id, firstname, secondname, login, password_hash
			 FROM users WHERE login = $1`
	err = r.db.Select(&user, query, login)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}
	if len(user) == 0 {
		return model.User{}, fmt.Errorf("%s: %w", op, ErrNoSuchUser)
	}
	if len(user) != 1 {
		return model.User{}, fmt.Errorf("%s: %w", op, ErrTwoSameLoginInDb)
	}
	if tx.Commit() != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user[0], nil
}
