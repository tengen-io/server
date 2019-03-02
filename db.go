package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	dbUrlFmt string = "postgres://%s:%s@%s:%d/%s?sslmode=disable"
)

type PostgresDBConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Database     string
}

func (c *PostgresDBConfig) Url() string {
	return fmt.Sprintf(dbUrlFmt, c.User, c.Password, c.Host, c.Port, c.Database)
}

func NewPostgresDb(config *PostgresDBConfig) (*sqlx.DB, error) {
	conn, err := sqlx.Open("postgres", config.Url())
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
