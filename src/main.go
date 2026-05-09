package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"iceslide/src/board"
	"iceslide/src/gui"
	"iceslide/src/heuristic"
	"iceslide/src/search"
	"iceslide/src/state"
	"iceslide/src/successor"
	"iceslide/src/visualizer"
)

const defaultSavePath = "test/solusi.txt"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	useGUI := flag.Bool("gui", false, "jalankan mode GUI")
	algName := flag.String("alg", "", "algoritma: ucs | gbfs | astar | idastar")
	hName := flag.String("h", "", "heuristik: zero | manhattan | remaining | euclidean | maxcombined")
	outPath := flag.String("out", defaultSavePath, "default file output untuk fitur simpan solusi")
	flag.Parse()

	if *useGUI || flag.NArg() == 0 {
		gui.Run()
		return nil
	}

	args := flag.Args()
	pr, err := board.LoadFile(args[0])
	if err != nil {
		return fmt.Errorf("load %s: %w", args[0], err)
	}
	start := state.State{X: pr.StartX, Y: pr.StartY, NextDigit: 0}

	scanner := bufio.NewScanner(os.Stdin)

	alg, err := chooseAlg(scanner, os.Stdout, *algName)
	if err != nil {
		return err
	}
	h, hLabel, err := chooseHeuristic(scanner, os.Stdout, alg, *hName)
	if err != nil {
		return err
	}

	res, err := search.Search(search.Config{
		Board:     pr.Board,
		Start:     start,
		Algorithm: alg,
		Heuristic: h,
		Generator: successor.SlideGenerator{},
	})
	if err != nil {
		return err
	}

	visualizer.PrintSummary(os.Stdout, alg, hLabel, pr.Board, start, res)

	return visualizer.RunInteractive(
		scanner, os.Stdout,
		alg, hLabel,
		pr.Board, start, res,
		*outPath,
	)
}

func chooseAlg(scanner *bufio.Scanner, w *os.File, flagVal string) (search.Algorithm, error) {
	if flagVal != "" {
		return parseAlg(flagVal)
	}
	return visualizer.SelectAlgorithm(scanner, w)
}

func chooseHeuristic(scanner *bufio.Scanner, w *os.File, alg search.Algorithm, flagVal string) (heuristic.Func, string, error) {
	if alg == search.UCS {
		return heuristic.Zero, "Zero (UCS)", nil
	}
	if flagVal != "" {
		return parseHeuristic(flagVal)
	}
	return visualizer.SelectHeuristic(scanner, w, alg)
}

func parseAlg(name string) (search.Algorithm, error) {
	switch strings.ToLower(name) {
	case "ucs":
		return search.UCS, nil
	case "gbfs":
		return search.GBFS, nil
	case "astar", "a*":
		return search.AStar, nil
	case "idastar", "ida*":
		return search.IDAStar, nil
	default:
		return 0, fmt.Errorf("algoritma tidak dikenal: %q", name)
	}
}

func parseHeuristic(name string) (heuristic.Func, string, error) {
	switch strings.ToLower(name) {
	case "zero":
		return heuristic.Zero, "Zero", nil
	case "manhattan":
		return heuristic.Manhattan, "Manhattan", nil
	case "remaining":
		return heuristic.RemainingDigits, "RemainingDigits", nil
	case "euclidean":
		return heuristic.Euclidean, "Euclidean", nil
	case "maxcombined":
		return heuristic.MaxCombined, "Max Combined", nil
	default:
		return nil, "", fmt.Errorf("heuristik tidak dikenal: %q", name)
	}
}
