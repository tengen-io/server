package game

type Color byte

type point struct {
	x int
	y int
}

const (
	White Color = iota
	Black Color = iota
)

type Game struct {
	board        *Board
	captures     []int
	currentColor Color
	move         int
	ko           *point
}

func NewGame(size int) *Game {
	return &Game{
		board:        NewBoard(size),
		captures:     make([]int, 2),
		currentColor: Black,
		move:         0,
	}
}

func (g *Game) PlayMove(x int, y int) error {
	// Ensure the position is empty
	old := g.board.GetNode(x, y)
	if old != empty {
		return NonEmptyError{}
	}

	if g.ko != nil && x == g.ko.x && y == g.ko.y {
		return KoViolationError{}
	}

	// count new liberties around the stone
	north, south, east, west := g.board.GetNeighbors(x, y)
	newLiberties := 0
	isolated := true
	potentialSuicide := false
	for _, neighbor := range []node{north, south, east, west} {
		if neighbor == empty {
			newLiberties += 1
		} else if neighbor == toNode(g.currentColor) {
			isolated = false
		}
	}

	if isolated && newLiberties == 0 {
		potentialSuicide = true
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
			potentialSuicide = true
		}
	}

	// check for suicide
	if len(toRemove) == 0 && potentialSuicide {
		return SuicideError{}
	}

	// The move is valid at this point. Clear Ko.
	g.ko = nil

	// Check for new ko.
	if isolated && len(toRemove) == 1 && len(toRemove[0]) == 1 {
		koX, koY := g.board.coord(toRemove[0][0])
		g.ko = &point{koX, koY}
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

type NonEmptyError struct{}

func (e NonEmptyError) Error() string {
	return "position is not empty"
}

type KoViolationError struct{}

func (e KoViolationError) Error() string {
	return "ko violation"
}

type SuicideError struct{}

func (e SuicideError) Error() string {
	return "move is suicidal"
}
