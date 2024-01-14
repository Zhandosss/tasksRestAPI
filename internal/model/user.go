package model

type User struct {
	ID         int64  `json:"-" db:"user_id"`
	FirstName  string `json:"first_name" db:"firstname"`
	SecondName string `json:"second_name" db:"secondname"`
	Login      string `json:"login" db:"login"`
	Password   string `json:"password" db:"password_hash"`
}
