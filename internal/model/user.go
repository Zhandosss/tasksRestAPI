package model

type User struct {
	ID        int64  `json:"-" db:"id"`
	FirstName string `json:"first_name" db:"firstname"`
	LastName  string `json:"last_name" db:"lastname"`
	Login     string `json:"login" db:"login"`
	Password  string `json:"password" db:"password_hash"`
}
