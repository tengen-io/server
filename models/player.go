package models

import (
	"time"
    "database/sql"
)

type PlayerParsed struct {
	Id        int    `json:id`
	UserId    int    `json:userId`
	GameId    int    `json:gameId`
	Status    string `json:status`
	Color     string `json:color`
	Stats     string  `json:stats`
	HasPassed bool   `json:hasPassed`
	Timestamps
}

type Player struct {
	User      *User
	PlayerParsed
}

type Timestamps struct {
	InsertedAt time.Time `json:insertedAt`
	UpdatedAt  time.Time `json:updatedAt`
}

func createPlayer(tx *sql.Tx, userId, gameId int, status, color string, time []byte) (*Player, error) {
	rows, err := tx.Query("INSERT INTO players VALUES (nextval('players_id_seq'), $1, $2, $3, $4, $5, $6, $7, $8) RETURNING *", userId, gameId, status, color, "{}", false, time, time)

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
		var player PlayerParsed
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
		players = append(players, Player{&User{}, player})
	}

	return players, rows.Err()
}
