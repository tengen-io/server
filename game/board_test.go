package game

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoard_SetNode(t *testing.T) {
	board := NewBoard(5)
	board.SetNode(2, 2, white)
	board.SetNode(3, 4, black)

	assert.Equal(t, board.board[12], white, "must be equal")
	assert.Equal(t, board.board[23], black, "must be equal")
}

func TestBoard_GetNode(t *testing.T) {
	board := NewBoard(5)
	board.board[10] = white
	board.board[3] = black

	assert.Equal(t, board.GetNode(0, 2), white, "must be equal")
	assert.Equal(t, board.GetNode(3, 0), black, "must be equal")
	assert.Equal(t, board.GetNode(4, 2), empty, "must be equal")
	assert.Equal(t, board.GetNode(2, 5), edge, "must be equal")
	assert.Equal(t, board.GetNode(-1, 3), edge, "must be equal")
}

func TestBoard_GetNeighboringStrings(t *testing.T) {
	board := NewBoard(5)
	board.SetNode(1, 3, black)
	board.SetNode(0, 3, black)
	board.SetNode(0, 2, black)
	board.SetNode(2, 3, black)

	board.SetNode(3, 2, white)
	board.SetNode(3, 1, white)
	board.SetNode(3, 3, white)

	board.SetNode(0, 0, white)
	board.SetNode(1, 0, white)
	board.SetNode(1, 1, white)
	board.SetNode(1, 2, white)

	_, strings := board.GetStringAndNeighbors(2, 2)
	assertStringContains(t, strings, []int{0, 1, 6, 11})
	assertStringContains(t, strings, []int{10, 15, 16, 17})
	assertStringContains(t, strings, []int{8, 13, 18})
}

func TestBoard_CountLiberties(t *testing.T) {
	board := NewBoard(5)
	board.SetNode(0, 4, black)
	board.SetNode(1, 4, black)
	board.SetNode(0, 3, black)
	board.SetNode(0, 2, white)
	board.SetNode(1, 3, white)
	board.SetNode(2, 4, white)

	board.SetNode(2, 0, black)
	board.SetNode(2, 1, black)
	board.SetNode(2, 2, black)
	board.SetNode(3, 2, black)
	board.SetNode(4, 2, black)
	board.SetNode(4, 1, black)
	board.SetNode(4, 0, black)
	board.SetNode(3, 0, black)

	string1, err := board.GetString(0, 4)
	assert.Nil(t, err)
	assert.Equal(t, board.CountLiberties(string1), 0)
}

func TestBoard_RemoveString(t *testing.T) {
	board := NewBoard(5)
	board.SetNode(2, 0, white)
	board.SetNode(2, 1, white)
	board.SetNode(2, 2, white)
	board.SetNode(3, 0, black)
	board.SetNode(3, 1, black)

	string, err := board.GetString(2, 0)
	assert.Nil(t, err)
	board.RemoveString(string)

	assert.Equal(t, empty, board.GetNode(2, 0))
	assert.Equal(t, empty, board.GetNode(2, 1))
	assert.Equal(t, empty, board.GetNode(2, 2))
	assert.Equal(t, black, board.GetNode(3, 0))
	assert.Equal(t, black, board.GetNode(3, 1))
}

// for debugging
func printBoard(board *Board) {
	for y := board.size - 1; y >= 0; y-- {
		for x := 0; x < board.size; x++ {
			node := board.board[board.idx(x, y)]
			var nodeStr = "."
			if node == white {
				nodeStr = "W"
			}
			if node == black {
				nodeStr = "B"
			}
			fmt.Printf("%s ", nodeStr)
		}
		fmt.Print("\n")
	}
}

func assertStringContains(t *testing.T, strings []nodestring, string nodestring) {
	for _, toCompare := range strings {
		if reflect.DeepEqual(toCompare, string) {
			return
		}
	}

	t.Errorf("%+v did not match any of %+v", string, strings)
}
