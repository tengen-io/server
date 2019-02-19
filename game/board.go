package game

import (
	"fmt"
	"sort"
)

type node byte
type nodestring []int

type Board struct {
	size int
	board []node
}

const (
	empty node = iota
	white node = iota
	black node = iota
	edge node = iota
	// TODO(eac): does ko need to be here?
	//	Ko = iota
)

func NewBoard(size int) *Board {
	board := make([]node, size * size)
	return &Board {
		size: size,
		board: board,
	}
}

func (b *Board) GetNode(x int, y int) node {
	if x < 0 || x >= b.size {
		return edge
	}

	if y < 0 || y >= b.size {
		return edge
	}

	return b.board[b.idx(x,y)]
}

func (b *Board) SetNode(x int, y int, value node) {
	b.board[b.idx(x,y)] = value
}

func (b *Board) GetString(x int, y int) (nodestring, error) {
	node := b.GetNode(x, y)
	if node == white || node == black {
		return b.findString(b.idx(x, y)), nil
	}

	return nil, fmt.Errorf("cannot get string of a %+v node", node)
}

func (b *Board) GetNeighbors(x int, y int) (node, node, node, node) {
	north := b.GetNode(x, y+1)
	south := b.GetNode(x, y-1)
	east := b.GetNode(x+1, y)
	west := b.GetNode(x-1, y)

	return north, south, east, west
}

func (b *Board) IsInString(haystack nodestring, x int, y int) bool {
	i := sort.SearchInts(haystack, b.idx(x, y))
	return i < len(haystack)
}

func (b *Board) GetStringAndNeighbors(x int, y int) (nodestring, []nodestring) {
	origin := b.GetNode(x, y)
	var string nodestring = nil
	if origin == white || origin == black {
		string = b.findString(b.idx(x, y))
	}
	neighbors := make([]nodestring, 0)

	// Technically we can dedupe strings more efficiently by unrolling this loop but
	// we can optimize this later if necessary
	isDupe := func(idx int) bool {
		if string != nil && sort.SearchInts(string, idx) < len(string) {
			return true
		}

		for _, str := range neighbors {
			i := sort.SearchInts(str, idx)
			if i < len(str) && str[i] == idx {
				return true
			}
		}
		return false
	}

	north := b.GetNode(x, y+1)
	if north == white || north == black {
		neighbors = append(neighbors, b.findString(b.idx(x, y+1)))
	}

	south := b.GetNode(x, y-1)
	if south == white || south == black && !isDupe(b.idx(x, y-1)) {
		neighbors = append(neighbors, b.findString(b.idx(x, y-1)))
	}

	east := b.GetNode(x+1, y)
	if east == white || east == black && !isDupe(b.idx(x+1, y)) {
		neighbors = append(neighbors, b.findString(b.idx(x+1, y)))
	}

	west := b.GetNode(x-1, y)
	if west == white || west == black && !isDupe(b.idx(x-1, y)) {
		neighbors = append(neighbors, b.findString(b.idx(x-1, y)))
	}

	return string, neighbors
}

func (b *Board) CountLiberties(string nodestring) int {
	rv := 0
	seen := make(map[int]bool)
	for _, idx := range string {
		x, y := b.coord(idx)
		north := b.GetNode(x, y+1)
		south := b.GetNode(x, y-1)
		east := b.GetNode(x+1, y)
		west := b.GetNode(x-1, y)

		if north == empty && !seen[b.idx(x, y+1)] {
			rv += 1
			seen[b.idx(x, y+1)] = true
		}

		if south == empty && !seen[b.idx(x, y-1)] {
			rv += 1
			seen[b.idx(x, y-1)] = true
		}

		if east == empty && !seen[b.idx(x+1, y)] {
			rv += 1
			seen[b.idx(x+1, y)] = true
		}

		if west == empty && !seen[b.idx(x-1, y)] {
			rv += 1
			seen[b.idx(x-1, y)] = true
		}
	}

	return rv
}

func (b *Board) StringColor(string nodestring) node {
	return b.board[string[0]]
}

func (b *Board) findString(idx int) nodestring {
	rv := make([]int, 0)
	stack := make([]int, 0)
	seen := make(map[int]bool)
	stack = append(stack, idx)
	origColor := b.board[idx]

	for len(stack) > 0 {
		i := stack[0]
		stack = stack[1:]
		seen[i] = true
		color := b.board[idx]
		if color == origColor {
			rv = append(rv, i)
		}

		northIdx := i + b.size
		southIdx := i - b.size
		eastIdx := i + 1
		westIdx := i - 1

		if northIdx < b.size * b.size && !seen[northIdx] && b.board[northIdx] == origColor {
			stack = append(stack, northIdx)
		}

		if southIdx > 0 && !seen[southIdx] && b.board[southIdx] == origColor {
			stack = append(stack, southIdx)
		}

		if eastIdx % b.size > idx % b.size && !seen[eastIdx] && b.board[eastIdx] == origColor {
			stack = append(stack, eastIdx)
		}

		if westIdx >= 0 && westIdx % b.size < idx % b.size && !seen[westIdx] && b.board[westIdx] == origColor {
			stack = append(stack, westIdx)
		}
	}

	sort.Ints(rv)
	return rv
}

func (b *Board) RemoveString(string nodestring) int {
	for _, idx := range string {
		b.board[idx] = empty
	}

	return len(string)
}

func (b* Board) idx(x int, y int) int {
	return x + y*b.size
}

func (b *Board) coord(idx int) (int, int) {
	return idx % b.size, idx / b.size
}
