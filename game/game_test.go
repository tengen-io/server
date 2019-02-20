package game

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type pointcolor struct {
	x     int
	y     int
	color node
}

func TestGame_PlayMove(t *testing.T) {
	testCases := []struct {
		name             string
		moves            []*point
		expectedBoard    []pointcolor
		expectedCaptures []int
	}{
		{
			"add moves",
			[]*point{{1, 1}, {2, 2}, {1, 0}, {4, 2}},
			[]pointcolor{{1, 1, black}, {2, 2, white}, {1, 0, black}, {4, 2, white}},
			[]int{0, 0},
		},
		{
			"capture corner",
			[]*point{{0, 0}, {1, 0}, nil, {0, 1}},
			[]pointcolor{{0, 0, empty}, {1, 0, white}, {0, 1, white}},
			[]int{0, 1},
		},
		{
			"capture side",
			[]*point{{1, 0}, {0, 0}, {2, 0}, {1, 1}, nil, {2, 1}, nil, {3, 0}},
			[]pointcolor{{1, 0, empty}, {0, 0, white}, {2, 0, empty}, {1, 1, white}, {2, 1, white}, {3, 0, white}},
			[]int{0, 2},
		},
		{
			"capture center",
			[]*point{{2, 2}, {2, 1}, {3, 2}, {3, 1}, {2, 3}, {1, 2}, {3, 3}, {1, 3}, nil, {2, 4}, nil, {3, 4}, nil, {4, 3}, nil, {4, 2}},
			[]pointcolor{{2, 2, empty}, {3, 2, empty}, {2, 3, empty}, {3, 3, empty}, {1, 2, white}, {1, 3, white}, {2, 4, white}, {3, 4, white}, {4, 3, white}, {4, 2, white}, {2, 1, white}, {3, 1, white}},
			[]int{0, 4},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			game := NewGame(5)
			for _, move := range testCase.moves {
				if move == nil {
					game.Pass()
				} else {
					assert.Nil(t, game.PlayMove(move.x, move.y))
				}
			}

			for _, expected := range testCase.expectedBoard {
				assert.Equal(t, expected.color, game.board.GetNode(expected.x, expected.y))
			}

			assert.Equal(t, len(testCase.moves), game.move)
		})
	}
}

func TestGame_PlayMove_ExistingStone(t *testing.T) {
	game := NewGame(5)
	assert.Nil(t, game.PlayMove(2, 2))
	err := game.PlayMove(2, 2)
	assert.EqualError(t, err, NonEmptyError{}.Error())
}

func TestGame_Pass(t *testing.T) {
	game := NewGame(5)
	game.Pass()
	assert.Equal(t, game.currentColor, White)
	assert.Equal(t, game.move, 1)
}

func TestGame_PlayMove_Suicide(t *testing.T) {
	testCases := []struct {
		name  string
		moves []*point
	}{
		{
			"suicide in corner single stone",
			[]*point{{0, 1}, nil, {1, 0}, {0, 0}},
		},
		{
			"suicide in corner string",
			[]*point{{0, 2}, {0, 0}, {2, 0}, {1, 0}, {1, 1}, {0, 1}},
		},
		{
			"suicide on edge single stone",
			[]*point{nil, {1, 0}, nil, {2, 1}, nil, {3, 0}, {2, 0}},
		},
		{
			"suicide on edge string",
			[]*point{{1, 0}, {0, 0}, {2, 0}, {1, 1}, nil, {2, 1}, nil, {3, 1}, nil, {4, 0}, {3, 0}},
		},
		{
			"suicide in center single stone",
			[]*point{nil, {0, 2}, nil, {1, 3}, nil, {2, 2}, nil, {1, 1}, {1, 2}},
		},
		{
			"suicide in center string",
			[]*point{nil, {0, 2}, {1, 2}, {1, 3}, {2, 3}, {2, 4}, {3, 2}, {3, 3}, {2, 1}, {4, 2}, nil, {3, 1}, nil, {2, 0}, nil, {1, 1}, {2, 2}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			game := NewGame(5)
			var err error
			for _, move := range testCase.moves {
				if move == nil {
					game.Pass()
				} else {
					err = game.PlayMove(move.x, move.y)
				}
			}

			assert.EqualError(t, err, SuicideError{}.Error())
		})
	}
}

func TestGame_PlayMove_Ko(t *testing.T) {
	game := NewGame(5)
	assert.Nil(t, game.PlayMove(2, 2))
	assert.Nil(t, game.PlayMove(1, 2))
	assert.Nil(t, game.PlayMove(3, 3))
	assert.Nil(t, game.PlayMove(2, 3))
	assert.Nil(t, game.PlayMove(4, 2))
	assert.Nil(t, game.PlayMove(2, 1))
	assert.Nil(t, game.PlayMove(3, 1))
	assert.Nil(t, game.PlayMove(3, 2))

	assert.NotNil(t, game.ko)
	assert.Equal(t, game.ko.x, 2)
	assert.Equal(t, game.ko.y, 2)

	err := game.PlayMove(2, 2)
	assert.EqualError(t, err, KoViolationError{}.Error())

	assert.Nil(t, game.PlayMove(4, 4))
	assert.Nil(t, game.ko)
	game.Pass()
	assert.Nil(t, game.PlayMove(2, 2))

	assert.NotNil(t, game.ko)
	assert.Equal(t, game.ko.x, 3)
	assert.Equal(t, game.ko.y, 2)
}

func printGame(game *Game) {
	for y := game.board.size - 1; y >= 0; y-- {
		for x := 0; x < game.board.size; x++ {
			node := game.board.board[game.board.idx(x, y)]
			var nodeStr = "."
			if game.ko != nil && game.ko.x == x && game.ko.y == y {
				nodeStr = "X"
			} else if node == white {
				nodeStr = "W"
			} else if node == black {
				nodeStr = "B"
			}
			fmt.Printf("%s ", nodeStr)
		}
		fmt.Print("\n")
	}
	fmt.Printf("Move: %d\n", game.move)
	fmt.Printf("B Captures: %d\n", game.captures[Black])
	fmt.Printf("W Captures: %d\n", game.captures[White])
}
