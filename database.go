package main

import (
	"database/sql"
	// "fmt"
	_ "github.com/lib/pq"
	// "log"
)

func Connect() (*sql.DB, error) {
	connStr := "user=postgres password=postgres dbname=go_stop_go sslmode=verify-full"
	return sql.Open("postgres", connStr)
}
