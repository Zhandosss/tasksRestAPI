package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New(config *Config) (*sqlx.DB, error) {
	var dataSource string
	switch config.ConnType {
	case "string":
		dataSource = config.ConnString
	case "parameters":
		dataSource = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			config.Username, config.Password, config.Host, config.Port, config.DB)
	}
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
