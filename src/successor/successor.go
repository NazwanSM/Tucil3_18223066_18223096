package successor

import (
	"iceslide/src/board"
	"iceslide/src/state"
)

type Direction uint8

const (
	Up Direction = iota
	Down
	Left
	Right
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "U"
	case Down:
		return "D"
	case Left:
		return "L"
	case Right:
		return "R"
	default:
		return "?"
	}
}

type Successor struct {
	State state.State
	StepCost int
	Direction Direction
}

type Generator interface {
	Generate(b *board.Board, s state.State) []Successor
}

type SlideGenerator struct{}

var allDirections = [...]Direction{Up, Down, Left, Right}

func (g SlideGenerator) Generate(b *board.Board, s state.State) []Successor {
	out := make([]Successor, 0, len(allDirections))
	for _, d := range allDirections {
		if succ, ok := g.slide(b, s, d); ok {
			out = append(out, succ)
		}
	}
	return out
}

func (g SlideGenerator) slide(b *board.Board, s state.State, d Direction) (Successor, bool) {
	dx, dy := delta(d)
	x, y := s.X, s.Y
	nextDigit := s.NextDigit
	cost := 0

	for {
		nx, ny := x+dx, y+dy

		if !b.InBounds(nx, ny) {
			return Successor{}, false
		}

		cell := b.At(nx, ny)

		switch cell.Type {
		case board.CellRock:
			if x == s.X && y == s.Y {
				return Successor{}, false
			}
			return Successor{
				State: state.State{X: x, Y: y, NextDigit: nextDigit},
				StepCost: cost,
				Direction: d,
			}, true

		case board.CellLava:
			return Successor{}, false

		case board.CellNumber:
			if cell.Digit > nextDigit {
				return Successor{}, false
			}
			if cell.Digit == nextDigit {
				nextDigit++
			}
			cost += cell.Cost
			x, y = nx, ny

		case board.CellExit, board.CellIce:
			cost += cell.Cost
			x, y = nx, ny
		}
	}
}

func delta(d Direction) (int, int) {
	switch d {
	case Up:
		return 0, -1
	case Down:
		return 0, 1
	case Left:
		return -1, 0
	case Right:
		return 1, 0
	default:
		return 0, 0
	}
}
