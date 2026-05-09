package gui

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"iceslide/src/board"
	"iceslide/src/heuristic"
	"iceslide/src/search"
	"iceslide/src/state"
	"iceslide/src/successor"
	"iceslide/src/visualizer"
)

func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveIndex)
	mux.HandleFunc("/api/solve", handleSolve)

	addr := "127.0.0.1:8765"
	fmt.Println("GUI berjalan di http://" + addr)
	fmt.Println("Tekan Ctrl+C untuk keluar.")

	go func() {
		time.Sleep(300 * time.Millisecond)
		openBrowser("http://" + addr)
	}()

	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Println("server error:", err)
	}
}

func openBrowser(url string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
}

type cellJSON struct {
	Type  string `json:"t"`
	Digit int    `json:"d"`
	Cost  int    `json:"c"`
}

type frameJSON struct {
	Index     int    `json:"index"`
	Action    string `json:"action"`
	StepCost  int    `json:"stepCost"`
	TotalCost int    `json:"totalCost"`
	PX        int    `json:"px"`
	PY        int    `json:"py"`
	Collected []bool `json:"collected"`
}

type solveResponse struct {
	Found          bool        `json:"found"`
	Path           string      `json:"path"`
	Cost           int         `json:"cost"`
	Steps          int         `json:"steps"`
	Iterations     int         `json:"iterations"`
	NodesGenerated int         `json:"nodesGenerated"`
	DurationMs     float64     `json:"durationMs"`
	Algorithm      string      `json:"algorithm"`
	Heuristic      string      `json:"heuristic"`
	Width          int         `json:"width"`
	Height         int         `json:"height"`
	Cells          [][]cellJSON `json:"cells"`
	Frames         []frameJSON `json:"frames"`
	Error          string      `json:"error,omitempty"`
}

func handleSolve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		jsonError(w, "parse form: "+err.Error())
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		jsonError(w, "baca file: "+err.Error())
		return
	}
	defer file.Close()

	pr, err := board.Parse(file)
	if err != nil {
		jsonError(w, "parse board: "+err.Error())
		return
	}

	algStr := r.FormValue("alg")
	heurStr := r.FormValue("heur")

	alg := parseAlg(algStr)
	h, hLabel := parseHeuristic(alg, heurStr)

	start := state.State{X: pr.StartX, Y: pr.StartY, NextDigit: 0}

	res, err := search.Search(search.Config{
		Board:     pr.Board,
		Start:     start,
		Algorithm: alg,
		Heuristic: h,
		Generator: successor.SlideGenerator{},
	})
	if err != nil {
		jsonError(w, "search: "+err.Error())
		return
	}

	b := pr.Board
	cells := make([][]cellJSON, b.Height())
	for y := 0; y < b.Height(); y++ {
		cells[y] = make([]cellJSON, b.Width())
		for x := 0; x < b.Width(); x++ {
			c := b.At(x, y)
			cells[y][x] = cellJSON{
				Type:  cellTypeName(c.Type),
				Digit: c.Digit,
				Cost:  c.Cost,
			}
		}
	}

	resp := solveResponse{
		Algorithm:      algStr,
		Heuristic:      hLabel,
		Width:          b.Width(),
		Height:         b.Height(),
		Cells:          cells,
		Iterations:     res.NodesExpanded,
		NodesGenerated: res.NodesGenerated,
		DurationMs:     float64(res.Duration.Microseconds()) / 1000.0,
	}

	if res.Goal != nil {
		resp.Found = true
		resp.Path = visualizer.PathString(res.Path)
		resp.Cost = res.Cost
		resp.Steps = len(res.Path) - 1

		frames := visualizer.BuildFrames(b, res.Path)
		resp.Frames = make([]frameJSON, len(frames))
		for i, f := range frames {
			resp.Frames[i] = frameJSON{
				Index:     f.Index,
				Action:    f.Action,
				StepCost:  f.StepCost,
				TotalCost: f.TotalCost,
				PX:        f.Pos.X,
				PY:        f.Pos.Y,
				Collected: f.Collected,
			}
		}
	}

	json.NewEncoder(w).Encode(resp)
}

func jsonError(w http.ResponseWriter, msg string) {
	resp := solveResponse{Error: msg}
	json.NewEncoder(w).Encode(resp)
}

func cellTypeName(ct board.CellType) string {
	switch ct {
	case board.CellRock:
		return "rock"
	case board.CellLava:
		return "lava"
	case board.CellExit:
		return "exit"
	case board.CellNumber:
		return "number"
	default:
		return "ice"
	}
}

func parseAlg(s string) search.Algorithm {
	switch strings.ToLower(s) {
	case "ucs":
		return search.UCS
	case "gbfs":
		return search.GBFS
	case "idastar", "ida*":
		return search.IDAStar
	default:
		return search.AStar
	}
}

func parseHeuristic(alg search.Algorithm, s string) (heuristic.Func, string) {
	if alg == search.UCS {
		return heuristic.Zero, "Zero (UCS)"
	}
	switch strings.ToLower(s) {
	case "remaining":
		return heuristic.RemainingDigits, "Remaining Digits"
	case "euclidean":
		return heuristic.Euclidean, "Euclidean"
	case "maxcombined":
		return heuristic.MaxCombined, "Max Combined"
	default:
		return heuristic.Manhattan, "Manhattan"
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, indexHTML)
}

