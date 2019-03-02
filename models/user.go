package models

/*
import (
	"database/sql"
	"github.com/badoux/checkmail"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const MinPasswordLength int = 8

type AuthUser struct {
	Jwt  string
	User *User
}

// GetUser gets a User by id or username.
func (db *PostgresDB) GetUser(identifier interface{}) (*User, error) {
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
func (db *PostgresDB) CreateUser(username, email, password, passwordConfirm string) (*User, error) {
	if len(password) < MinPasswordLength {
		return nil, passwordTooShortError{}
	}
	if password != passwordConfirm {
		return nil, passwordMismatchError{}
	}
	if err := checkmail.ValidateFormat(email); err != nil {
		return nil, invalidEmailError{}
	}

	pw, err := bcrypt.GenerateFromPassword([]byte(password), db.config.BcryptCost)

	if err != nil {
		return nil, err
	}
	time := pq.FormatTimestamp(time.Now())

	rows, err := db.Query("INSERT INTO users (username, email, encrypted_password, inserted_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING *", username, email, pw, time, time)
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
}*/
