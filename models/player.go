package models

import (
	"time"
    "database/sql"
)

type Player struct {
	Id        int    `json:id`
	userId    int    `json:userId`
	gameId    int    `json:gameId`
	Status    string `json:status`
	Color     string `json:color`
	Stats     Stats  `json:stats`
	HasPassed bool   `json:hasPassed`
	User      *User
	Timestamps
}

type Timestamps struct {
	InsertedAt time.Time `json:insertedAt`
	UpdatedAt  time.Time `json:updatedAt`
}

type Stats struct{}

func createPlayer(tx *sql.Tx, userId, gameId int, status, color string, time []byte) (*Player, error) {
	rows, err := tx.Query("INSERT INTO players VALUES (nextval('players_id_seq'), $1, $2, $3, $4, $5, $6, $7, $8) RETURNING *", userId, gameId, status, color, "{}", false, time, time)

	if err != nil {
		return nil, HandleError(err)
	}

	players, _ := parsePlayerRows(rows)

	return &players[0], nil
}
