package heuristic

import (
	"iceslide/src/board"
	"iceslide/src/state"
)

type Func func(s state.State, b *board.Board) int

func Zero(s state.State, b *board.Board) int { return 0 }

func Manhattan(s state.State, b *board.Board) int {
	tx, ty, ok := findTarget(s, b)
	if !ok {
		return 0
	}
	return abs(tx-s.X) + abs(ty-s.Y)
}

func RemainingDigits(s state.State, b *board.Board) int {
	remaining := b.MaxDigit() + 1 - s.NextDigit
	if remaining < 0 {
		remaining = 0
	}
	return remaining + 1 
}

func findTarget(s state.State, b *board.Board) (int, int, bool) {
	if s.NextDigit <= b.MaxDigit() {
		for y := 0; y < b.Height(); y++ {
			for x := 0; x < b.Width(); x++ {
				c := b.At(x, y)
				if c.Type == board.CellNumber && c.Digit == s.NextDigit {
					return x, y, true
				}
			}
		}
	}
	for y := 0; y < b.Height(); y++ {
		for x := 0; x < b.Width(); x++ {
			if b.At(x, y).Type == board.CellExit {
				return x, y, true
			}
		}
	}
	return 0, 0, false
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
