package models

import (
	"database/sql"
	"encoding/json"
)

type Stone struct {
	Id     int    `json:"id"`
	gameId int    `json:"gameId"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Color  string `json:"color"`
	Game   *Game
	Timestamps
}

// Scan implements a Scanner for a Stone, so that it can be decoded from a json string
// from the db.
func (stone *Stone) Scan(src interface{}) error {
	err := json.Unmarshal(src.([]byte), stone)
	if err != nil {
		return err
	}

	return nil
}

func parseStoneRows(rows *sql.Rows) ([]Stone, error) {
	defer rows.Close()

	stones := make([]Stone, 0)
	for rows.Next() {
		var stone Stone
		rows.Scan(
			&stone.Id,
			&stone.gameId,
			&stone.X,
			&stone.Y,
			&stone.Color,
			&stone.InsertedAt,
			&stone.UpdatedAt,
		)
		stones = append(stones, stone)
	}

	return stones, rows.Err()
}
