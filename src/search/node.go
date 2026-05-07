package search

import (
	"iceslide/src/state"
	"iceslide/src/successor"
)

type Node struct {
	State  state.State
	Parent *Node
	Action successor.Direction
	G int
	H int
}

func (n *Node) F() int { return n.G + n.H }

func (n *Node) Path() []*Node {
	depth := 0
	for cur := n; cur != nil; cur = cur.Parent {
		depth++
	}
	path := make([]*Node, depth)
	for cur := n; cur != nil; cur = cur.Parent {
		depth--
		path[depth] = cur
	}
	return path
}
