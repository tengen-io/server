package rules

import (
	"github.com/camirmas/go_stop/models"
	// "fmt"
)

// A String is a chain of Stones on a Go Board. A string is defined as any
// set of Stones for which each Stone is adjacent to at least one other Stone.
type String []models.Stone

func Run(board *models.Board, stone models.Stone) ([]String, error) {
	strings := getStrings(board)
	toRemove := make([]String, 0)

	for _, str := range strings {
		if !contains(str, stone) {
			_, numLiberties := findLiberties(board, str)

			if numLiberties == 0 {
				toRemove = append(toRemove, str)
			}
		}
	}

	return toRemove, nil
}

func findLiberties(board *models.Board, str String) ([]models.Stone, int) {
	liberties := make([]models.Stone, 0)

	for _, s := range str {
		nearby := getNearby(board, s)

		for i, ns := range nearby {
			if ns.Color == "" && !contains(liberties, ns) {
				liberties = append(liberties, nearby[i])
			}
		}
	}

	return liberties, len(liberties)
}

func getStrings(board *models.Board) []String {
	strings := make([]String, 0)

	for _, stone := range board.Stones {
		str := getString(board, stone)
		strings = append(strings, str)
	}

	return strings
}

func getString(board *models.Board, stone models.Stone) String {
	acc := []models.Stone{stone}
	str := make(String, 0)

	for len(acc) > 0 {
		var s models.Stone
		s, acc = acc[len(acc)-1], acc[:len(acc)-1]

		if !contains(str, s) && s.Color == stone.Color {
			str = append(str, s)
			nearby := getNearby(board, s)
			acc = append(acc, nearby...)
		}
	}

	return str
}

func getNearby(board *models.Board, stone models.Stone) []models.Stone {
	up := models.Stone{X: stone.X, Y: stone.Y + 1}
	down := models.Stone{X: stone.X, Y: stone.Y - 1}
	left := models.Stone{X: stone.X - 1, Y: stone.Y}
	right := models.Stone{X: stone.X + 1, Y: stone.Y}

	nearbyStones := []models.Stone{up, down, left, right}
	validStones := make([]models.Stone, 0)

	for _, s := range nearbyStones {
		if isInbounds(board.Size, s) {
			existingStone := find(board.Stones, s)

			if existingStone != nil {
				validStones = append(validStones, *existingStone)
			} else {
				validStones = append(validStones, s)
			}
		}
	}

	return validStones
}

func isInbounds(size int, stone models.Stone) bool {
	if stone.X < 0 || stone.X >= size {
		return false
	}
	if stone.Y < 0 || stone.Y >= size {
		return false
	}
	return true
}

func find(stones []models.Stone, stone models.Stone) *models.Stone {
	for _, s := range stones {
		if s.X == stone.X && s.Y == stone.Y {
			return &s
		}
	}
	return nil
}

func contains(list []models.Stone, item models.Stone) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
