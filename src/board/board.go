package board

type Board struct {
	width    int
	height   int
	grid     [][]Cell
	maxDigit int 
}

func New(grid [][]Cell) *Board {
	height := len(grid)
	width := 0
	if height > 0 {
		width = len(grid[0])
	}
	maxDigit := -1
	for _, row := range grid {
		for _, c := range row {
			if c.Type == CellNumber && c.Digit > maxDigit {
				maxDigit = c.Digit
			}
		}
	}
	return &Board{
		width:    width,
		height:   height,
		grid:     grid,
		maxDigit: maxDigit,
	}
}

func (b *Board) Width() int { return b.width }

func (b *Board) Height() int { return b.height }

func (b *Board) MaxDigit() int { return b.maxDigit }

func (b *Board) At(x, y int) Cell { return b.grid[y][x] }

func (b *Board) InBounds(x, y int) bool {
	return x >= 0 && x < b.width && y >= 0 && y < b.height
}
