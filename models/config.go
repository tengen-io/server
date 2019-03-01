package models

import (
	"fmt"
)

const (
	dbUrlFmt string = "postgres://%s:%s@%s:%d/%s?sslmode=disable"
)

type PostgresDBConfig struct {
	Host       string
	Port       int
	User       string
	Password   string
	Database   string
	BcryptCost int
}

func (c *PostgresDBConfig) Url() string {
	return fmt.Sprintf(dbUrlFmt, c.User, c.Password, c.Host, c.Port, c.Database)
}
