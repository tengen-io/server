package main

import (
	"database/sql"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func ConnectDB() (*sql.DB, error) {
	connStr := "user=postgres password=postgres dbname=go_stop_go"
	return sql.Open("postgres", connStr)
}

func GetUser(db *sql.DB, username string) (*User, error) {
	rows, _ := db.Query("SELECT * FROM users WHERE username = $1", username)

	users, err := parseUserRows(rows)

	if err != nil {
		return nil, err
	} else if len(users) > 0 {
		return &users[0], nil
	} else {
		return nil, userNotFoundError{}
	}
}

func CreateUser(db *sql.DB, username string, email string, password string, passwordConfirm string) (*User, error) {
	if password != passwordConfirm {
		return nil, passwordMismatchError{}
	}
	pw, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return nil, err
	}
	time := pq.FormatTimestamp(time.Now())

	_, err = db.Query("INSERT INTO users VALUES (nextval('users_id_seq'), $1, $2, $3, $4, $5)", username, email, pw, time, time)
	if err != nil {
		return nil, handleError(err)
	}

	user, _ := GetUser(db, username)

	return user, nil
}

func CheckPw(db *sql.DB, username string, password string) (*User, error) {
	user, err := GetUser(db, username)

	if err != nil {
		return user, err
	}

	inputPassword := []byte(password)
	userPassword := []byte(user.encryptedPassword)
	err = bcrypt.CompareHashAndPassword(userPassword, inputPassword)
	return user, handleError(err)
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

func handleError(e error) error {
	if e == nil {
		return e
	}

	switch e.Error() {
	case "pq: duplicate key value violates unique constraint \"users_username_key\"":
		return usernameTakenError{e}
	case "pq: duplicate key value violates unique constraint \"users_email_key\"":
		return emailTakenError{e}
	case "crypto/bcrypt: hashedPassword is not the hash of the given password":
		return invalidLoginError{e}
	default:
		return e
	}
}

type passwordMismatchError struct {
	Err error
}
type usernameTakenError struct {
	Err error
}
type emailTakenError struct {
	Err error
}
type invalidEmailError struct {
	Err error
}
type invalidLoginError struct {
	Err error
}
type userNotFoundError struct{}

func (e passwordMismatchError) Error() string {
	return "Passwords do not match"
}

func (e usernameTakenError) Error() string {
	return "Username is already taken"
}

func (e emailTakenError) Error() string {
	return "Email is already taken"
}

func (e invalidEmailError) Error() string {
	return "Email format is invalid"
}

func (e invalidLoginError) Error() string {
	return "Invalid login"
}

func (e userNotFoundError) Error() string {
	return "User not found"
}
