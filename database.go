package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func ConnectDB() (*sql.DB, error) {
	connStr := "user=postgres password=postgres dbname=go_stop_go"
	return sql.Open("postgres", connStr)
}

func GetUser(db *sql.DB, username string) (*User, error) {
	rows, err := db.Query("SELECT * FROM users WHERE username = $1", username)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.encryptedPassword,
			&user.InsertedAt,
			&user.UpdatedAt); err != nil {
			log.Println(err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
	}

	if len(users) > 0 {
		return &users[0], nil
	} else {
		return nil, nil
	}
}
