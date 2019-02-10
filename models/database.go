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

type PostgresDB struct {
	config *PostgresDBConfig
	*sql.DB
}

type DB interface {
	CheckPw(username, password string) (*User, error)
	GetUser(id interface{}) (*User, error)
	CreateUser(username, email, password, passwordConfirm string) (*User, error)
	GetGame(gameId interface{}) (*Game, error)
	GetGames(userId interface{}) ([]*Game, error)
	CreateGame(userId int, opponent *User) (*Game, error)
	UpdateGame(userId int, game *Game, toAdd Stone, toRemove []Stone) error
	Pass(userId int, game *Game) (*Game, error)
}

// NewPostgresDB creates a connection to the postgres database
func NewPostgresDB(config *PostgresDBConfig) (*PostgresDB, error) {
	conn, err := sql.Open("postgres", config.Url())
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, err
	}
	db := &PostgresDB{config, conn}

	return db, nil
}

// CheckPw compares the given password against the encrypted password for the given User.
func (db *PostgresDB) CheckPw(username, password string) (*User, error) {
	user, err := db.GetUser(username)

	if err != nil {
		return user, err
	}

	inputPassword := []byte(password)
	userPassword := []byte(user.encryptedPassword)
	err = bcrypt.CompareHashAndPassword(userPassword, inputPassword)
	return user, HandleError(err)
}
