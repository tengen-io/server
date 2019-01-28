package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"time"
)

type Game struct {
	Id           int
	Status       string
	PlayerTurnId int
	BoardSize    int
	LastTaker    *Stone
	Players      []Player
	Stones       []Stone
	Timestamps
}

const SmallBoardSize int = 13
const RegBoardSize int = 19

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
		rows, _ := db.Query("SELECT DISTINCT ON (games.id) games.* FROM games JOIN players P ON P.user_id = $1", userId)
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
func (db *DB) CreateGame(userId int, opponent *User) (*Game, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	time := pq.FormatTimestamp(time.Now())

	// Create Game
	rows, err := tx.Query("INSERT INTO games (status, board_size, inserted_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING *", "active", RegBoardSize, time, time)
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

	// Create Player 2 (invitee)
	_, err = createPlayer(tx, opponent.Id, game.Id, "user-pending", "white", time)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	_, err = db.Exec("UPDATE games SET player_turn_id = $1 WHERE id = $2", player1.Id, game.Id)

	game, _ = db.GetGame(game.Id)

	return game, nil
}

// Pass is a game action where a player decides that they cannot make a
// move on their turn. If both players pass, the game ends.
func (db *DB) Pass(userId int, game *Game) (*Game, error) {
	currentPlayer, otherPlayer := game.CurrentPlayer(userId)

	time := pq.FormatTimestamp(time.Now())

	if otherPlayer.HasPassed {
		tx, _ := db.Begin()

		_, err := tx.Exec("UPDATE players SET (has_passed, updated_at) = ($1, $2) WHERE id = $3", true, time, currentPlayer.Id)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}

		_, err = tx.Exec("UPDATE games SET (status, updated_at) = ($1, $2) WHERE id = $3", "complete", time, game.Id)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}

		if err := tx.Commit(); err != nil {
			return nil, err
		}
	} else {
		tx, _ := db.Begin()
		_, err := tx.Exec("UPDATE players SET (has_passed, updated_at) = ($1, $2) WHERE id = $3", true, time, currentPlayer.Id)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		_, err = tx.Exec("UPDATE games SET (player_turn_id, updated_at) = ($1, $2) WHERE id = $3", otherPlayer.Id, time, game.Id)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}

		if err := tx.Commit(); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	game, _ = db.GetGame(game.Id)

	return game, nil
}

func (db *DB) UpdateBoard(userId int, game *Game, toAdd Stone, toRemove []Stone) (*Game, error) {
	_, otherPlayer := game.CurrentPlayer(userId)
	time := pq.FormatTimestamp(time.Now())
	tx, _ := db.Begin()
	_, err := tx.Exec("UPDATE games SET (player_turn_id, updated_at) = ($1, $2) where id = $3", otherPlayer.Id, time, game.Id)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	ids := make([]string, 0)
	for _, stone := range toRemove {
		ids = append(ids, fmt.Sprintf("%d", stone.Id))
	}
	allIds := strings.Join(ids, ", ")

	_, err = tx.Exec("DELETE FROM stones WHERE id IN ($1)", allIds)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO stones (game_id, x, y, color, inserted_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", game.Id, toAdd.X, toAdd.Y, toAdd.Color, time, time)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if game.LastTaker != nil {
		encoded, _ := json.Marshal(game.LastTaker)
		_, err := tx.Exec("UPDATE games SET last_taker = $1", encoded)

		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	game, err = db.GetGame(game.Id)

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

	rows, _ = db.Query("SELECT * FROM stones where game_id = $1", game.Id)
	stones, _ := parseStoneRows(rows)

	game.Stones = stones

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
			&game.BoardSize,
			&game.LastTaker,
			&game.InsertedAt,
			&game.UpdatedAt,
		)
		games = append(games, game)
	}

	return games, rows.Err()
}
