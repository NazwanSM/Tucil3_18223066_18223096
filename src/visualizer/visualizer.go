package visualizer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"iceslide/src/board"
	"iceslide/src/heuristic"
	"iceslide/src/search"
	"iceslide/src/state"
	"iceslide/src/successor"
)

type Frame struct {
	Index     int
	Pos       state.State
	Action    string 
	StepCost  int   
	TotalCost int  
	Collected []bool 
}

func BuildFrames(b *board.Board, path []*search.Node) []Frame {
	if len(path) == 0 {
		return nil
	}
	frames := make([]Frame, len(path))
	for i, node := range path {
		action := "Initial"
		stepCost := 0
		if i > 0 {
			action = directionName(node.Action)
			stepCost = node.G - path[i-1].G
		}
		frames[i] = Frame{
			Index:     i,
			Pos:       node.State,
			Action:    action,
			StepCost:  stepCost,
			TotalCost: node.G,
			Collected: collectedSet(node.State.NextDigit, b.MaxDigit()),
		}
	}
	return frames
}

func PrintSummary(w io.Writer, alg search.Algorithm, hName string, b *board.Board, start state.State, res *search.Result) {
	fmt.Fprintln(w, "== Initial Board ==")
	fmt.Fprint(w, RenderInitial(b, start))

	fmt.Fprintln(w)
	fmt.Fprintln(w, "== Solution Summary ==")
	fmt.Fprintf(w, "Algorithm       : %s\n", alg)
	fmt.Fprintf(w, "Heuristic       : %s\n", hName)

	if res.Goal == nil {
		fmt.Fprintln(w, "Result          : NO PATH FOUND")
	} else {
		fmt.Fprintf(w, "Path            : %s\n", PathString(res.Path))
		fmt.Fprintf(w, "Cost            : %d\n", res.Cost)
		fmt.Fprintf(w, "Steps           : %d\n", len(res.Path)-1)
	}
	fmt.Fprintf(w, "Iterations      : %d\n", res.NodesExpanded)
	fmt.Fprintf(w, "Nodes generated : %d\n", res.NodesGenerated)
	fmt.Fprintf(w, "Duration        : %s\n", res.Duration)
}

func PathString(path []*search.Node) string {
	if len(path) <= 1 {
		return "(no moves)"
	}
	parts := make([]string, 0, len(path)-1)
	for i := 1; i < len(path); i++ {
		parts = append(parts, path[i].Action.String())
	}
	return strings.Join(parts, " -> ")
}

func RenderInitial(b *board.Board, s state.State) string {
	return renderBoardAt(b, s, collectedSet(0, b.MaxDigit()))
}

func RenderFrame(b *board.Board, f Frame) string {
	return renderBoardAt(b, f.Pos, f.Collected)
}

func Playback(scanner *bufio.Scanner, w io.Writer, b *board.Board, path []*search.Node) error {
	frames := BuildFrames(b, path)
	if len(frames) == 0 {
		fmt.Fprintln(w, "(tidak ada frame untuk diputar)")
		return nil
	}

	idx := 0
	for {
		printFrame(w, b, frames[idx], len(frames)-1)
		fmt.Fprint(w, "[n/Enter] next  [p] prev  [j] jump  [q] keluar > ")

		if !scanner.Scan() {
			return scanner.Err()
		}
		cmd := strings.ToLower(strings.TrimSpace(scanner.Text()))

		switch cmd {
		case "n", "":
			if idx < len(frames)-1 {
				idx++
			} else {
				fmt.Fprintln(w, "(sudah di langkah terakhir)")
			}
		case "p":
			if idx > 0 {
				idx--
			} else {
				fmt.Fprintln(w, "(sudah di langkah pertama)")
			}
		case "j":
			fmt.Fprint(w, "Pada step berapa anda ingin melakukan playback : ")
			if !scanner.Scan() {
				return scanner.Err()
			}
			target, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil || target < 0 || target > len(frames)-1 {
				fmt.Fprintf(w, "(input tidak valid; range 0..%d)\n", len(frames)-1)
				continue
			}
			idx = target
		case "q":
			return nil
		default:
			fmt.Fprintln(w, "(perintah tidak dikenal)")
		}
	}
}

func PromptYesNo(scanner *bufio.Scanner, w io.Writer, prompt string) (bool, error) {
	for {
		fmt.Fprintf(w, "%s (Ya/Tidak): ", prompt)
		if !scanner.Scan() {
			return false, scanner.Err()
		}
		switch strings.ToLower(strings.TrimSpace(scanner.Text())) {
		case "ya", "y", "yes":
			return true, nil
		case "tidak", "t", "n", "no", "":
			return false, nil
		default:
			fmt.Fprintln(w, "(jawab dengan Ya atau Tidak)")
		}
	}
}

func PromptLine(scanner *bufio.Scanner, w io.Writer, prompt string) (string, error) {
	fmt.Fprint(w, prompt)
	if !scanner.Scan() {
		return "", scanner.Err()
	}
	return strings.TrimSpace(scanner.Text()), nil
}

func SelectAlgorithm(scanner *bufio.Scanner, w io.Writer) (search.Algorithm, error) {
	for {
		fmt.Fprintln(w, "== Pilih Algoritma ==")
		fmt.Fprintln(w, "  1. Uniform Cost Search (UCS)")
		fmt.Fprintln(w, "  2. Greedy Best-First Search (GBFS)")
		fmt.Fprintln(w, "  3. A*")
		fmt.Fprint(w, "Pilihan [1-3]: ")
		if !scanner.Scan() {
			return 0, scanner.Err()
		}
		switch strings.TrimSpace(scanner.Text()) {
		case "1":
			return search.UCS, nil
		case "2":
			return search.GBFS, nil
		case "3":
			return search.AStar, nil
		default:
			fmt.Fprintln(w, "(input tidak valid; pilih 1, 2, atau 3)")
		}
	}
}

func SelectHeuristic(scanner *bufio.Scanner, w io.Writer, alg search.Algorithm) (heuristic.Func, string, error) {
	if alg == search.UCS {
		return heuristic.Zero, "Zero (UCS)", nil
	}
	for {
		fmt.Fprintln(w, "== Pilih Heuristik ==")
		fmt.Fprintln(w, "  1. Manhattan")
		fmt.Fprintln(w, "  2. Remaining Digits")
		fmt.Fprint(w, "Pilihan [1-2]: ")
		if !scanner.Scan() {
			return nil, "", scanner.Err()
		}
		switch strings.TrimSpace(scanner.Text()) {
		case "1":
			return heuristic.Manhattan, "Manhattan", nil
		case "2":
			return heuristic.RemainingDigits, "RemainingDigits", nil
		default:
			fmt.Fprintln(w, "(input tidak valid; pilih 1 atau 2)")
		}
	}
}

func RunInteractive(
	scanner *bufio.Scanner,
	w io.Writer,
	alg search.Algorithm,
	hName string,
	b *board.Board,
	start state.State,
	res *search.Result,
	defaultSavePath string,
) error {
	if res == nil || res.Goal == nil {
		return nil
	}

	play, err := PromptYesNo(scanner, w, "Apakah Anda ingin melakukan playback?")
	if err != nil {
		return err
	}
	if play {
		if err := Playback(scanner, w, b, res.Path); err != nil {
			return err
		}
	}

	save, err := PromptYesNo(scanner, w, "Apakah Anda ingin menyimpan solusi?")
	if err != nil {
		return err
	}
	if !save {
		return nil
	}

	outPath, err := PromptLine(scanner, w,
		fmt.Sprintf("Nama file output [%s]: ", defaultSavePath))
	if err != nil {
		return err
	}
	if outPath == "" {
		outPath = defaultSavePath
	}

	if err := SaveSolution(outPath, alg, hName, b, start, res); err != nil {
		return fmt.Errorf("save %s: %w", outPath, err)
	}
	fmt.Fprintf(w, "Solusi disimpan pada %s\n", outPath)
	return nil
}

func SaveSolution(path string, alg search.Algorithm, hName string, b *board.Board, start state.State, res *search.Result) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	defer bw.Flush()

	PrintSummary(bw, alg, hName, b, start, res)

	if res.Goal == nil {
		return nil
	}

	fmt.Fprintln(bw)
	fmt.Fprintln(bw, "== Step-by-step ==")
	for _, fr := range BuildFrames(b, res.Path) {
		fmt.Fprintln(bw)
		printFrame(bw, b, fr, len(res.Path)-1)
	}
	return nil
}

func printFrame(w io.Writer, b *board.Board, f Frame, last int) {
	fmt.Fprintf(w, "\n--- Step %d / %d  (%s) ---\n", f.Index, last, f.Action)
	fmt.Fprint(w, RenderFrame(b, f))
	if f.Index > 0 {
		fmt.Fprintf(w, "Move: %s   step cost: %d   total cost: %d\n",
			f.Action, f.StepCost, f.TotalCost)
	}
}

func renderBoardAt(b *board.Board, pos state.State, collected []bool) string {
	var sb strings.Builder
	for y := 0; y < b.Height(); y++ {
		for x := 0; x < b.Width(); x++ {
			if x == pos.X && y == pos.Y {
				sb.WriteRune('Z')
				continue
			}
			c := b.At(x, y)
			if c.Type == board.CellNumber && c.Digit < len(collected) && collected[c.Digit] {
				sb.WriteRune('*')
				continue
			}
			sb.WriteRune(glyphFor(c))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func collectedSet(nextDigit, maxDigit int) []bool {
	size := maxDigit + 1
	if size < 0 {
		size = 0
	}
	out := make([]bool, size)
	for d := 0; d < nextDigit && d < size; d++ {
		out[d] = true
	}
	return out
}

func glyphFor(c board.Cell) rune {
	switch c.Type {
	case board.CellRock:
		return 'X'
	case board.CellLava:
		return 'L'
	case board.CellExit:
		return 'O'
	case board.CellNumber:
		return rune('0' + c.Digit)
	default:
		return '*'
	}
}

func directionName(d successor.Direction) string {
	switch d {
	case successor.Up:
		return "Up"
	case successor.Down:
		return "Down"
	case successor.Left:
		return "Left"
	case successor.Right:
		return "Right"
	default:
		return "?"
	}
}
