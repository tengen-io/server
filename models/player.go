package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Player struct {
	Id        int
	UserId    int
	GameId    int
	Status    string
	Color     string
	Stats     *Stats
	HasPassed bool
	User      *User
	Timestamps
}

type Stats struct{}

// Scan implements a Scanner for Player Stats, so that they can be decoded from a json string
// from the db.
func (stats *Stats) Scan(src interface{}) error {
	err := json.Unmarshal(src.([]byte), stats)
	if err != nil {
		return err
	}

	return nil
}

type Timestamps struct {
	InsertedAt time.Time
	UpdatedAt  time.Time
}

func createPlayer(tx *sql.Tx, userId, gameId interface{}, status, color string, time []byte) (*Player, error) {
	stats, err := json.Marshal(Stats{})
	rows, err := tx.Query("INSERT INTO players (user_id, game_id, status, color, stats, inserted_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *", userId, gameId, status, color, stats, time, time)

	if err != nil {
		return nil, HandleError(err)
	}

	players, _ := parsePlayerRows(rows)

	return &players[0], nil
}

func parsePlayerRows(rows *sql.Rows) ([]Player, error) {
	defer rows.Close()

	players := make([]Player, 0)
	for rows.Next() {
		var player Player
		rows.Scan(
			&player.Id,
			&player.UserId,
			&player.GameId,
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
