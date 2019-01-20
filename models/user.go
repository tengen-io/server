package models

import (
	"database/sql"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const MinPasswordLength int = 8

type User struct {
	Id                int
	Username          string
	Email             string
	encryptedPassword string
	Games             []Game
	Timestamps
}

type AuthUser struct {
	Jwt  string
	User *User
}

// GetUser gets a User by id or username.
func (db *DB) GetUser(identifier interface{}) (*User, error) {
	var rows *sql.Rows
	switch identifier.(type) {
	case int:
		rows, _ = db.Query("SELECT * FROM users WHERE id = $1", identifier)
	case string:
		rows, _ = db.Query("SELECT * FROM users WHERE username = $1", identifier)
	}

	users, _ := parseUserRows(rows)
	if len(users) == 0 {
		return nil, userNotFoundError{}
	}

	return &users[0], nil
}

// CreateUser starts a new User account.
func (db *DB) CreateUser(username, email, password, passwordConfirm string) (*User, error) {
	if len(password) < MinPasswordLength {
		return nil, passwordTooShortError{}
	}
	if password != passwordConfirm {
		return nil, passwordMismatchError{}
	}
	pw, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return nil, err
	}
	time := pq.FormatTimestamp(time.Now())

	rows, err := db.Query("INSERT INTO users VALUES (nextval('users_id_seq'), $1, $2, $3, $4, $5) RETURNING *", username, email, pw, time, time)
	if err != nil {
		return nil, HandleError(err)
	}
	users, _ := parseUserRows(rows)
	user := &users[0]

	return user, nil
}

func parseUserRows(rows *sql.Rows) ([]User, error) {
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.encryptedPassword,
			&user.InsertedAt,
			&user.UpdatedAt,
		)
		users = append(users, user)
	}

	return users, rows.Err()
}
