package main

import (
	"database/sql"
	_ "fmt"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func ConnectDB() (*sql.DB, error) {
	connStr := "user=postgres password=postgres dbname=go_stop_go"
	return sql.Open("postgres", connStr)
}

func CheckPw(db *sql.DB, username string, password string) (*User, error) {
	user, err := GetUser(db, username)

	if err != nil {
		return user, err
	}

	inputPassword := []byte(password)
	userPassword := []byte(user.encryptedPassword)
	err = bcrypt.CompareHashAndPassword(userPassword, inputPassword)
	return user, HandleError(err)
}

func GetUser(db *sql.DB, identifier interface{}) (*User, error) {
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

func CreateUser(db *sql.DB, username string, email string, password string, passwordConfirm string) (*User, error) {
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

func GetGame(db *sql.DB, gameId int) (*Game, error) {
	rows, _ := db.Query("SELECT * from games where id = $1", gameId)
	games, _ := parseGameRows(rows)
	if len(games) == 0 {
		return nil, gameNotFoundError{}
	}
	game := &games[0]

	if err := buildGame(db, game); err != nil {
		return nil, err
	}

	return game, nil
}

func GetGames(db *sql.DB, userId int) ([]*Game, error) {
	rows, _ := db.Query("SELECT games.* FROM games JOIN players P ON P.user_id = $1", userId)
	games, _ := parseGameRows(rows)

	gameRefs := make([]*Game, 0)
	for _, game := range games {
		buildGame(db, &game)
		gameRefs = append(gameRefs, &game)
	}

	return gameRefs, nil
}

func CreateGame(db *sql.DB, userId, opponentId int) (*Game, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	time := pq.FormatTimestamp(time.Now())

	rows, err := db.Query("INSERT INTO games VALUES (nextval('games_id_seq'), $1, $2, $3, $4) RETURNING *", "not-started", nil, time, time)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	games, _ := parseGameRows(rows)
	game := &games[0]

	player1, err := createPlayer(db, userId, game.Id, "active", "black", time)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	user, _ := GetUser(db, userId)
	player1.User = user
	player1.Game = game

	player2, err := createPlayer(db, opponentId, game.Id, "user-pending", "white", time)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	user, _ = GetUser(db, opponentId)
	player2.User = user
	player2.Game = game

	game.Players = []Player{*player1, *player2}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return game, nil
}

func createPlayer(db *sql.DB, userId, gameId int, status, color string, time []byte) (*Player, error) {
	rows, err := db.Query("INSERT INTO players VALUES (nextval('players_id_seq'), $1, $2, $3, $4, $5, $6, $7, $8) RETURNING *", userId, gameId, status, color, "{}", false, time, time)

	if err != nil {
		return nil, HandleError(err)
	}

	players, _ := parsePlayerRows(rows)

	return &players[0], nil
}

func buildGame(db *sql.DB, game *Game) error {
	rows, _ := db.Query("SELECT * from players where game_id = $1", game.Id)
	players, _ := parsePlayerRows(rows)

	// We save an extra query by getting both users
	rows, _ = db.Query("SELECT DISTINCT users.* FROM users JOIN players P ON (P.game_id = $1)", game.Id)
	users, _ := parseUserRows(rows)

	for i, player := range players {
		for j, user := range users {
			if player.userId == user.Id {
				players[i].User = &users[j]
			}
		}
	}

	game.Players = players

	return nil
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

func parseGameRows(rows *sql.Rows) ([]Game, error) {
	defer rows.Close()

	games := make([]Game, 0)
	for rows.Next() {
		var game Game
		rows.Scan(
			&game.Id,
			&game.Status,
			&game.PlayerTurnId,
			&game.InsertedAt,
			&game.UpdatedAt,
		)
		games = append(games, game)
	}

	return games, rows.Err()
}

func parsePlayerRows(rows *sql.Rows) ([]Player, error) {
	defer rows.Close()

	players := make([]Player, 0)
	for rows.Next() {
		var player Player
		rows.Scan(
			&player.Id,
			&player.userId,
			&player.gameId,
			&player.Status,
			&player.Color,
			&player.Stats,
			&player.HasPassed,
			&player.InsertedAt,
			&player.UpdatedAt,
		)
		players = append(players, player)
	}

	return players, rows.Err()
}
