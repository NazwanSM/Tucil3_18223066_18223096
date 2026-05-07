package search

type Algorithm uint8

const (
	UCS   Algorithm = iota 
	GBFS                   
	AStar                  
)

func (a Algorithm) String() string {
	switch a {
	case UCS:
		return "UCS"
	case GBFS:
		return "GBFS"
	case AStar:
		return "A*"
	default:
		return "unknown"
	}
}

func priorityFor(a Algorithm) func(*Node) int {
	switch a {
	case UCS:
		return func(n *Node) int { return n.G }
	case GBFS:
		return func(n *Node) int { return n.H }
	case AStar:
		return func(n *Node) int { return n.G + n.H }
	default:
		return func(n *Node) int { return n.G + n.H }
	}
}
