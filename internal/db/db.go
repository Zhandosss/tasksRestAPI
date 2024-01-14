package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"restAPI/internal/config"
)

func New(config *config.DB) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("auth=%s password=%s host=%s dbname=%s sslmode=require",
		config.User, config.Password, config.Host, config.Name)
	conn, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		return nil, fmt.Errorf("sqlx connect: %w", err)
	}
	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}
	return conn, nil

}
