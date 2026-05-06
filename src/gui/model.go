package gui

import (
	"iceSlidingPuzzle/src/filehandler"
	"os"
	"strings"
)

// Difficulty levels
type Difficulty int

const (
	Easy         Difficulty = iota
	Intermediate Difficulty = iota
	Hard         Difficulty = iota
)

func (d Difficulty) String() string {
	switch d {
	case Easy:
		return "Easy"
	case Intermediate:
		return "Intermediate"
	case Hard:
		return "Hard"
	}
	return "Unknown"
}

// MapEntry represents a map file in the library
type MapEntry struct {
	Filename   string
	Width      int
	Height     int
	Difficulty Difficulty
}

// SolverAlgorithm types
type SolverAlgorithm int

const (
	AlgorithmUCS   SolverAlgorithm = iota
	AlgorithmGBFS  SolverAlgorithm = iota
	AlgorithmAStar SolverAlgorithm = iota
)

func (a SolverAlgorithm) String() string {
	switch a {
	case AlgorithmUCS:
		return "Uniform Cost Search (UCS)"
	case AlgorithmGBFS:
		return "Greedy Best-First (GBFS)"
	case AlgorithmAStar:
		return "A* Search Optimizer"
	}
	return "Unknown"
}

// Direction for moves (GUI-facing)
type Direction int

const (
	Up    Direction = iota
	Down  Direction = iota
	Left  Direction = iota
	Right Direction = iota
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "Move Up"
	case Down:
		return "Move Down"
	case Left:
		return "Move Left"
	case Right:
		return "Move Right"
	}
	return "Unknown"
}

// SolverStep represents a single step in the solution
type SolverStep struct {
	StepNum   int
	Direction Direction
	Units     int
	Cost      int
}

// SolverResult holds the result of a solver run
type SolverResult struct {
	Found      bool
	Steps      []SolverStep
	TotalCost  int
	TotalMoves int
	DurationMs int
	Weight     int
	Level      int
	Seed       string
}

// DefaultLibrary enumerates local `.txt` map files and attempts to
// parse them to populate width/height. This replaces the old
// hardcoded list and lets the UI reflect actual files.
func DefaultLibrary() []*MapEntry {
	entries := make([]*MapEntry, 0)

	files, err := os.ReadDir(".")
	if err != nil {
		return entries
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".txt") {
			continue
		}
		board, err := filehandler.ParseBoard(name)
		if err != nil {
			continue
		}
		diff := Easy
		maxDim := board.N
		if board.M > maxDim {
			maxDim = board.M
		}
		if maxDim <= 7 {
			diff = Easy
		} else if maxDim <= 12 {
			diff = Intermediate
		} else {
			diff = Hard
		}
		entries = append(entries, &MapEntry{Filename: name, Width: board.M, Height: board.N, Difficulty: diff})
	}

	return entries
}

// DefaultSolverResult returns nil. The UI should construct real
// results from the solver run; fallbacks should be handled by callers.
func DefaultSolverResult() *SolverResult {
	return nil
}
