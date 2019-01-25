/*
Package models is responsible for application data models as well as database
connections and queries. The current implementation is written to use postgres,
but uses an interface that could support a different database.
*/
package models

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	Config *DbConfig
	*sql.DB
}

type Database interface {
	CheckPw(username, password string) (*User, error)
	GetUser(id interface{}) (*User, error)
	CreateUser(username, email, password, passwordConfirm string) (*User, error)
	GetGame(gameId interface{}) (*Game, error)
	GetGames(userId interface{}) ([]*Game, error)
	CreateGame(userId int, opponent *User) (*Game, error)
	UpdateBoard(userId int, game *Game) (*Game, error)
	Pass(userId int, game *Game) (*Game, error)
}

// ConnectDB creates a connection with the postgres database, using credentials
// pulled from environment variables:
//
//	POSTGRES_DB - reads value
//	POSTGRES_USER - reads a file
//	POSTGRES_PASSWORD - reads a file
// 	POSTGRES_HOST - reads value
func ConnectDB() (*DB, error) {
	config := setupConfig()
	conn, err := sql.Open("postgres", config.DbUrl)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, err
	}
	db := &DB{config, conn}

	return db, nil
}

// CheckPw compares the given password against the encrypted password for the given User.
func (db *DB) CheckPw(username, password string) (*User, error) {
	user, err := db.GetUser(username)

	if err != nil {
		return user, err
	}

	inputPassword := []byte(password)
	userPassword := []byte(user.encryptedPassword)
	err = bcrypt.CompareHashAndPassword(userPassword, inputPassword)
	return user, HandleError(err)
}
