package models

import (
	"encoding/json"
)

const SmallBoardSize int = 13
const RegBoardSize int = 19

type Board struct {
	Size      int     `json:"size"`
	LastTaker *Stone  `json:"lastTaker"`
	Stones    []Stone `json:"stones"`
}

type Stone struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Color string `json:"color"`
}

// Scan implements a Scanner for a Board, so that it can be decoded from a json string
// from the db.
func (board *Board) Scan(src interface{}) error {
	err := json.Unmarshal(src.([]byte), board)
	if err != nil {
		return err
	}

	return nil
}
