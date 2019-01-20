package models

import (
	"database/sql"
	"encoding/json"
	_ "fmt"
	"github.com/lib/pq"
	"time"
)

type Game struct {
	Id           int
	Status       string
	PlayerTurnId int
	Board        *Board
	Players      []Player
	Timestamps
}

func (game Game) CurrentPlayer(userId int) (Player, Player) {
	var otherPlayer Player
	var currentPlayer Player
	for i, player := range game.Players {
		if player.UserId != userId {
			otherPlayer = game.Players[i]
		} else {
			currentPlayer = game.Players[i]
		}
	}
	return currentPlayer, otherPlayer
}

// GetGame gets a Game by id.
func (db *DB) GetGame(gameId interface{}) (*Game, error) {
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

// GetGames gets a list of Games for a given User id.
func (db *DB) GetGames(userId interface{}) ([]*Game, error) {
	var games []Game
	if userId == nil {
		rows, _ := db.Query("SELECT * FROM games ORDER BY updated_at limit 20")
		games, _ = parseGameRows(rows)
	} else {
		rows, _ := db.Query("SELECT DISTINCT games.* FROM games JOIN players P ON P.user_id = $1", userId)
		games, _ = parseGameRows(rows)
	}

	gameRefs := make([]*Game, 0)
	for i, _ := range games {
		buildGame(db, &games[i])
		gameRefs = append(gameRefs, &games[i])
	}

	return gameRefs, nil
}

// CreateGame builds all the necessary information to start a game, including
// associated Player entries.
func (db *DB) CreateGame(userId, opponentId interface{}) (*Game, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	time := pq.FormatTimestamp(time.Now())

	board := Board{Size: SmallBoardSize, Stones: []Stone{}}
	encodedBoard, err := json.Marshal(board)

	// Create Game
	rows, err := tx.Query("INSERT INTO games VALUES (nextval('games_id_seq'), $1, $2, $3, $4, $5) RETURNING *", "active", nil, encodedBoard, time, time)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	games, _ := parseGameRows(rows)
	game := &games[0]

	// Create Player 1 (inviter)
	player1, err := createPlayer(tx, userId, game.Id, "active", "black", time)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	user, _ := db.GetUser(userId)
	player1.User = user

	// Create Player 2 (invitee)
	player2, err := createPlayer(tx, opponentId, game.Id, "user-pending", "white", time)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	user, _ = db.GetUser(opponentId)
	player2.User = user

	game.Players = []Player{*player1, *player2}

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	_, err = db.Exec("UPDATE games SET player_turn_id = $1 WHERE id = $2", player1.Id, game.Id)

	game.PlayerTurnId = player1.Id

	return game, nil
}

// Pass is a game action where a player decides that they cannot make a
// move on their turn. If both players pass, the game ends.
func (db *DB) Pass(userId int, game *Game) (*Game, error) {
	currentPlayer, otherPlayer := game.CurrentPlayer(userId)

	time := pq.FormatTimestamp(time.Now())

	if otherPlayer.HasPassed {
		tx, _ := db.Begin()
		_, err := tx.Exec("UPDATE players SET (has_passed, updated_at) = ($1, $2) WHERE id = $3 RETURNING *", true, time, currentPlayer.Id)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		rows, err := tx.Query("UPDATE games SET (status, updated_at) = ($1, $2) WHERE id = $3 RETURNING *", "complete", time, game.Id)

		if err != nil {
			return nil, err
		}

		games, _ := parseGameRows(rows)
		game = &games[0]

		buildGame(db, game)

		if err := tx.Commit(); err != nil {
			return nil, err
		}

		return game, nil
	} else {
		tx, _ := db.Begin()
		_, err := tx.Exec("UPDATE players SET (has_passed, updated_at) = ($1, $2) WHERE id = $3 RETURNING *", true, time, currentPlayer.Id)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		rows, err := tx.Query("UPDATE games SET (player_turn_id, updated_at) = ($1, $2) WHERE id = $3 RETURNING *", otherPlayer.Id, time, game.Id)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}

		games, _ := parseGameRows(rows)
		game = &games[0]

		buildGame(db, game)

		if err := tx.Commit(); err != nil {
			return nil, err
		}

		return game, nil
	}
}

func (db *DB) UpdateBoard(userId int, game *Game) (*Game, error) {
	_, otherPlayer := game.CurrentPlayer(userId)
	time := pq.FormatTimestamp(time.Now())
	board, _ := json.Marshal(game.Board)
	rows, err := db.Query("UPDATE games SET (board, player_turn_id, updated_at) = ($1, $2, $3) where id = $4 RETURNING *", board, otherPlayer.Id, time, game.Id)

	if err != nil {
		return nil, err
	}

	games, _ := parseGameRows(rows)

	game = &games[0]
	buildGame(db, game)

	return game, nil
}

func buildGame(db *DB, game *Game) error {
	rows, _ := db.Query("SELECT * from players where game_id = $1", game.Id)
	players, _ := parsePlayerRows(rows)

	// We save an extra query by getting both users
	rows, _ = db.Query("SELECT DISTINCT users.* FROM users JOIN players P ON (P.game_id = $1)", game.Id)
	users, _ := parseUserRows(rows)

	for i, player := range players {
		for j, user := range users {
			if player.UserId == user.Id {
				players[i].User = &users[j]
			}
		}
	}

	game.Players = players

	return nil
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
			&game.Board,
			&game.InsertedAt,
			&game.UpdatedAt,
		)
		games = append(games, game)
	}

	return games, rows.Err()
}
