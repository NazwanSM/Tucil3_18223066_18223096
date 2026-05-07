package search

import (
	"container/heap"
	"errors"
	"time"

	"iceslide/src/board"
	"iceslide/src/heuristic"
	"iceslide/src/state"
	"iceslide/src/successor"
)

type Result struct {
	Goal *Node
	Path []*Node
	Cost int
	NodesExpanded int
	NodesGenerated int
	Duration time.Duration
}

type Config struct {
	Board *board.Board
	Start state.State
	Algorithm Algorithm
	Heuristic heuristic.Func
	Generator successor.Generator
}

func Search(cfg Config) (*Result, error) {
	if cfg.Board == nil {
		return nil, errors.New("search: nil board")
	}
	if cfg.Generator == nil {
		return nil, errors.New("search: nil generator")
	}
	if cfg.Heuristic == nil {
		return nil, errors.New("search: nil heuristic")
	}

	startTime := time.Now()

	root := &Node{
		State: cfg.Start,
		G: 0,
		H: cfg.Heuristic(cfg.Start, cfg.Board),
	}

	pq := newPriorityQueue(priorityFor(cfg.Algorithm))
	heap.Push(pq, root)

	bestG := map[uint64]int{cfg.Start.Key(): 0}
	expanded := 0
	generated := 1

	for pq.Len() > 0 {
		node := heap.Pop(pq).(*Node)

		if g, ok := bestG[node.State.Key()]; ok && g < node.G {
			continue
		}

		if isGoal(cfg.Board, node) {
			return &Result{
				Goal: node,
				Path: node.Path(),
				Cost: node.G,
				NodesExpanded: expanded,
				NodesGenerated: generated,
				Duration: time.Since(startTime),
			}, nil
		}

		expanded++

		for _, succ := range cfg.Generator.Generate(cfg.Board, node.State) {
			childG := node.G + succ.StepCost
			childKey := succ.State.Key()

			if g, ok := bestG[childKey]; ok && g <= childG {
				continue
			}
			bestG[childKey] = childG

			child := &Node{
				State: succ.State,
				Parent: node,
				Action: succ.Direction,
				G: childG,
				H: cfg.Heuristic(succ.State, cfg.Board),
			}
			heap.Push(pq, child)
			generated++
		}
	}

	return &Result{
		Goal: nil,
		Path: nil,
		Cost: 0,
		NodesExpanded: expanded,
		NodesGenerated: generated,
		Duration: time.Since(startTime),
	}, nil
}

func isGoal(b *board.Board, n *Node) bool {
	if !b.InBounds(n.State.X, n.State.Y) {
		return false
	}
	cell := b.At(n.State.X, n.State.Y)
	return cell.Type == board.CellExit && n.State.NextDigit > b.MaxDigit()
}
