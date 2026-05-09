package board

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type ParseResult struct {
	Board  *Board
	StartX int
	StartY int
}

func LoadFile(path string) (*ParseResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	return Parse(f)
}

func Parse(r io.Reader) (*ParseResult, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1<<20)

	N, M, err := parseDims(scanner)
	if err != nil {
		return nil, err
	}

	grid, startX, startY, err := parseGrid(scanner, N, M)
	if err != nil {
		return nil, err
	}

	if err := parseCosts(scanner, grid, N, M); err != nil {
		return nil, err
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	return &ParseResult{
		Board:  New(grid),
		StartX: startX,
		StartY: startY,
	}, nil
}

func parseDims(s *bufio.Scanner) (int, int, error) {
	if !s.Scan() {
		if err := s.Err(); err != nil {
			return 0, 0, fmt.Errorf("read dims: %w", err)
		}
		return 0, 0, errors.New("empty input")
	}
	fields := strings.Fields(s.Text())
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("dims line: want \"N M\", got %q", s.Text())
	}
	N, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("dims N: %w", err)
	}
	M, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("dims M: %w", err)
	}
	if N <= 0 || M <= 0 {
		return 0, 0, fmt.Errorf("dims must be positive: N=%d M=%d", N, M)
	}
	return N, M, nil
}

func parseGrid(s *bufio.Scanner, N, M int) ([][]Cell, int, int, error) {
	grid := make([][]Cell, N)
	startX, startY := -1, -1
	for y := 0; y < N; y++ {
		if !s.Scan() {
			return nil, 0, 0, fmt.Errorf("expected %d board rows, got %d", N, y)
		}
		runes := []rune(s.Text())
		if len(runes) < M {
			return nil, 0, 0, fmt.Errorf("board row %d: want %d cols, got %d", y, M, len(runes))
		}
		row := make([]Cell, M)
		for x := 0; x < M; x++ {
			cell, isStart, err := glyphToCell(runes[x])
			if err != nil {
				return nil, 0, 0, fmt.Errorf("board row %d col %d: %w", y, x, err)
			}
			if isStart {
				if startX != -1 {
					return nil, 0, 0, fmt.Errorf("multiple start tiles: (%d,%d) and (%d,%d)", startX, startY, x, y)
				}
				startX, startY = x, y
			}
			row[x] = cell
		}
		grid[y] = row
	}
	if startX == -1 {
		return nil, 0, 0, errors.New("no start tile (Z) found")
	}
	return grid, startX, startY, nil
}

func parseCosts(s *bufio.Scanner, grid [][]Cell, N, M int) error {
	for y := 0; y < N; y++ {
		if !s.Scan() {
			return fmt.Errorf("expected %d cost rows, got %d", N, y)
		}
		fields := strings.Fields(s.Text())
		if len(fields) < M {
			return fmt.Errorf("cost row %d: want %d values, got %d", y, M, len(fields))
		}
		for x := 0; x < M; x++ {
			cost, err := strconv.Atoi(fields[x])
			if err != nil {
				return fmt.Errorf("cost row %d col %d: %w", y, x, err)
			}
			grid[y][x].Cost = cost
		}
	}
	return nil
}

func glyphToCell(ch rune) (Cell, bool, error) {
	switch ch {
	case 'X':
		return Cell{Type: CellRock}, false, nil
	case '*':
		return Cell{Type: CellIce}, false, nil
	case 'L':
		return Cell{Type: CellLava}, false, nil
	case 'O':
		return Cell{Type: CellExit}, false, nil
	case 'Z':
		return Cell{Type: CellIce}, true, nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return Cell{Type: CellNumber, Digit: int(ch - '0')}, false, nil
	default:
		return Cell{}, false, fmt.Errorf("unknown glyph %q", ch)
	}
}
