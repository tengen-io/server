package game

import "fmt"

type Color byte

const (
	White Color = iota
	Black Color = iota
)

type Game struct {
	board *Board
	captures []int
	currentColor Color
	move int
}

func NewGame(size int) *Game {
	return &Game {
		board: NewBoard(size),
		captures: make([]int, 2),
		currentColor: Black,
		move: 0,
	}
}

func (g *Game) PlayMove(x int, y int) error {
	// Ensure the position is empty
	old := g.board.GetNode(x, y)
	if old != empty {
		return fmt.Errorf("position has a stone already")
	}

	// count new liberties around the stone
	north, south, east, west := g.board.GetNeighbors(x, y)
	newLiberties := 0
	isolated := true
	for _, neighbor := range []node{north, south, east, west} {
		if neighbor == empty {
			newLiberties += 1
		} else if neighbor == toNode(g.currentColor) {
			isolated = false
		}
	}

	if isolated && newLiberties == 0 {
		return fmt.Errorf("move is suicidal")
	}

	// we have to check the surrounding groups to ensure move legality and compute captures
	_, neighborStrings := g.board.GetStringAndNeighbors(x, y)
	toRemove := make([]nodestring, 0)

	for _, string := range neighborStrings {
		liberties := g.board.CountLiberties(string)
		stringColor := g.board.StringColor(string)

		if liberties == 1 && toNode(opp(g.currentColor)) == stringColor {
			// capture
			toRemove = append(toRemove, string)
		} else if liberties == 1 && toNode(g.currentColor) == stringColor && newLiberties == 0 {
			return fmt.Errorf("move is suicidal")
		}
	}

	// check for suicide
	if len(toRemove) == 0 {
		north, south, east, west := g.board.GetNeighbors(x, y)
		if north != empty && south != empty && east != empty && west != empty {
			return fmt.Errorf("move is suicidal")
		}
	}

	// if we haven't exploded yet, remove all captured strings
	for _, string := range toRemove {
		removed := g.board.RemoveString(string)
		g.captures[g.currentColor] += removed
	}

	// set the node
	g.board.SetNode(x, y, toNode(g.currentColor))

	// swap color, increment move count
	g.currentColor = opp(g.currentColor)
	g.move += 1
	return nil
}

func (g *Game) Pass() {
	g.currentColor = opp(g.currentColor)
	g.move += 1
}

func opp(c Color) Color {
	if c == White {
		return Black
	}
	return White
}

func toNode(c Color) node {
	if c == White {
		return white
	}

	return black
}
