package search

import (
	"container/heap"
	"errors"
	"time"
	"math"

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
	if cfg.Algorithm == IDAStar {
		return searchIDAStar(cfg)
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

func searchIDAStar(cfg Config) (*Result, error) {
	startTime := time.Now()
	
	root := &Node{
		State: cfg.Start,
		G: 0,
		H: cfg.Heuristic(cfg.Start, cfg.Board),
	}

	bound := root.H 

	expanded := 0
	generated := 1
	pathVisited := make(map[uint64]bool)

	var dfs func(node *Node, currentBound int) (bool, int, *Node)
	dfs = func(node *Node, currentBound int) (bool, int, *Node) {
		f := node.G + node.H
		
		if f > currentBound {
			return false, f, nil
		}
		if isGoal(cfg.Board, node) {
			return true, f, node
		}

		expanded++
		minOutlier := math.MaxInt32
		pathVisited[node.State.Key()] = true

		for _, succ := range cfg.Generator.Generate(cfg.Board, node.State) {
			if pathVisited[succ.State.Key()] {
				continue 
			}

			child := &Node{
				State: succ.State,
				Parent: node,
				Action: succ.Direction,
				G: node.G + succ.StepCost,
				H: cfg.Heuristic(succ.State, cfg.Board),
			}
			generated++

			found, newBound, goalNode := dfs(child, currentBound)
			if found {
				return true, newBound, goalNode
			}
			if newBound < minOutlier {
				minOutlier = newBound
			}
		}

		delete(pathVisited, node.State.Key())
		return false, minOutlier, nil
	}

	for {
		found, newBound, goalNode := dfs(root, bound)
		if found {
			return &Result{
				Goal: goalNode,
				Path: goalNode.Path(),
				Cost: goalNode.G,
				NodesExpanded: expanded,
				NodesGenerated: generated,
				Duration: time.Since(startTime),
			}, nil
		}

		if newBound == math.MaxInt32 {
			return &Result{
				Goal: nil,
				Path: nil,
				Cost: 0,
				NodesExpanded:  expanded,
				NodesGenerated: generated,
				Duration:       time.Since(startTime),
			}, nil
		}
		
		bound = newBound 
	}
}

func isGoal(b *board.Board, n *Node) bool {
	if !b.InBounds(n.State.X, n.State.Y) {
		return false
	}
	cell := b.At(n.State.X, n.State.Y)
	return cell.Type == board.CellExit && n.State.NextDigit > b.MaxDigit()
}
