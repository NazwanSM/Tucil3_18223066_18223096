package board

type CellType uint8

const (
	CellIce    CellType = iota 
	CellRock                  
	CellLava                  
	CellExit                  
	CellNumber                
)

type Cell struct {
	Type  CellType
	Cost  int
	Digit int
}
