/*
Package rules is responsible for evaluating a Go board according to the game rules.
The rules followed for this implementation come from the British Go Association,
and can be found at https://www.britgo.org/intro/intro2.html
*/
package rules
/*
import (
	_ "fmt"
	"github.com/tengen-io/server/models"
	"reflect"
	"sort"
)

// A String is a chain of Stones on a Go Board. A string is defined as any
// set of Stones for which each Stone is adjacent to at least one other Stone.
type String []models.Stone

// Run determines whether any Stones should be removed from the board, based on
// game rules.
func Run(boardSize int, stones []models.Stone, stone models.Stone) ([]String, error) {
	strings := getStrings(boardSize, stones)

	toRemove := make([]String, 0)

	for _, str := range strings {
		_, numLiberties := findLiberties(boardSize, stones, str)

		if numLiberties == 0 {
			toRemove = append(toRemove, str)
		}
	}

	if len(toRemove) == 1 && contains(toRemove[0], stone) {
		return []String{}, selfCaptureError{}
	}
	if len(toRemove) > 1 {
		var updatedRemove []String
		for _, str := range toRemove {
			if !contains(str, stone) {
				updatedRemove = append(updatedRemove, str)
			}
		}
		toRemove = updatedRemove
	}

	return toRemove, nil
}

func findLiberties(boardSize int, stones []models.Stone, str String) ([]models.Stone, int) {
	liberties := make([]models.Stone, 0)

	for _, s := range str {
		nearby := getNearby(boardSize, stones, s)

		for i, ns := range nearby {
			if ns.Color == "" && !contains(liberties, ns) {
				liberties = append(liberties, nearby[i])
			}
		}
	}

	return liberties, len(liberties)
}

func getStrings(boardSize int, stones []models.Stone) []String {
	strings := make([]String, 0)

	for _, stone := range stones {
		newString := getString(boardSize, stones, stone)

		sort.Slice(newString, func(i, j int) bool {
			if newString[i].Y < newString[j].Y {
				return true
			}
			if newString[i].Y > newString[j].Y {
				return false
			}
			return newString[i].X < newString[j].X
		})

		if !containsString(strings, newString) {
			strings = append(strings, newString)
		}
	}

	return strings
}

func getString(boardSize int, stones []models.Stone, stone models.Stone) String {
	acc := []models.Stone{stone}
	str := make(String, 0)

	for len(acc) > 0 {
		var s models.Stone
		s, acc = acc[len(acc)-1], acc[:len(acc)-1]

		if !contains(str, s) && s.Color == stone.Color {
			str = append(str, s)
			nearby := getNearby(boardSize, stones, s)
			acc = append(acc, nearby...)
		}
	}

	return str
}

func getNearby(boardSize int, stones []models.Stone, stone models.Stone) []models.Stone {
	up := models.Stone{X: stone.X, Y: stone.Y + 1}
	down := models.Stone{X: stone.X, Y: stone.Y - 1}
	left := models.Stone{X: stone.X - 1, Y: stone.Y}
	right := models.Stone{X: stone.X + 1, Y: stone.Y}

	nearbyStones := []models.Stone{up, down, left, right}
	validStones := make([]models.Stone, 0)

	for _, s := range nearbyStones {
		if isInbounds(boardSize, s) {
			existingStone := find(stones, s)

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
		if reflect.DeepEqual(v, item) {
			return true
		}
	}
	return false
}

func containsString(strings []String, str String) bool {
	for _, s := range strings {
		if reflect.DeepEqual(s, str) {
			return true
		}
	}
	return false
}

type selfCaptureError struct{}

func (e selfCaptureError) Error() string {
	return "Move is suicidal"
} */
