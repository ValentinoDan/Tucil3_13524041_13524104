package gui

import (
	"iceSlidingPuzzle/src/algorithm"
	"iceSlidingPuzzle/src/filehandler"
	"iceSlidingPuzzle/src/puzzle"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// Page identifiers
type Page int

const (
	PageLibrary Page = iota
	PageSolver  Page = iota
	PageResult  Page = iota
)

type MainUI struct {
	app         fyne.App
	window      fyne.Window
	currentPage Page
	mainContent *fyne.Container

	// shared state
	selectedMap        *MapEntry
	selectedAlgo       SolverAlgorithm
	solverResult       *SolverResult
	currentBoard       *puzzle.Board
	solverPath         []puzzle.Point
	librarySearchQuery string
	library            []*MapEntry

	// pages
	libraryPage *LibraryPage
	solverPage  *SolverPage
	resultPage  *ResultPage
}

func NewMainUI(a fyne.App, w fyne.Window) *MainUI {
	m := &MainUI{
		app:          a,
		window:       w,
		currentPage:  PageLibrary,
		selectedAlgo: AlgorithmUCS,
		currentBoard: nil, // nil until a map is loaded
		solverResult: nil,
		library:      DefaultLibrary(),
	}
	return m
}

func (m *MainUI) Build() fyne.CanvasObject {
	m.mainContent = container.NewStack()
	m.showPage(PageLibrary)
	return m.mainContent
}

func (m *MainUI) showPage(p Page) {
	m.currentPage = p
	var content fyne.CanvasObject

	switch p {
	case PageLibrary:
		m.libraryPage = NewLibraryPage(m)
		content = m.libraryPage.Build()
	case PageSolver:
		m.solverPage = NewSolverPage(m)
		content = m.solverPage.Build()
	case PageResult:
		m.resultPage = NewResultPage(m)
		content = m.resultPage.Build()
	}

	m.mainContent.Objects = []fyne.CanvasObject{content}
	m.mainContent.Refresh()
}

func (m *MainUI) NavigateTo(p Page) {
	m.showPage(p)
}

func (m *MainUI) SelectMap(entry *MapEntry) {
	m.selectedMap = entry

	board, err := filehandler.ParseBoard(entry.FullPath)
	if err == nil {
		m.currentBoard = board
		m.solverResult = nil
		m.solverPath = nil
	}

	m.showPage(PageSolver)
}

func (m *MainUI) RunSolver() {
	startedAt := time.Now()

	// fallback to default result (no map condition)
	if m.selectedMap == nil {
		m.solverResult = &SolverResult{Found: false}
		m.solverPath = nil
		m.showPage(PageResult)
		return
	}

	// Parse selected map file into a board
	board, err := filehandler.ParseBoard(m.selectedMap.FullPath)
	if err != nil {
		// failed to parse — show empty result
		m.solverResult = &SolverResult{Found: false}
		m.solverPath = nil
		m.showPage(PageResult)
		return
	}

	var path []puzzle.Point
	var totalCost int
	var found bool
	var iter int

	switch m.selectedAlgo {
	case AlgorithmUCS:
		node, iterTemp := algorithm.UCS(board)
		iter = iterTemp
		if node != nil {
			// reconstruct path from node
			var rev []puzzle.Point
			for n := node; n != nil; n = n.Parent {
				rev = append(rev, n.State.Pos)
			}
			// reverse
			for i := len(rev) - 1; i >= 0; i-- {
				path = append(path, rev[i])
			}
			totalCost = node.Cost
			found = true
		}
	case AlgorithmGBFS:
		path, totalCost, iter = algorithm.GreedyBestFirstSearch(board)
		found = len(path) > 0
	case AlgorithmAStar:
		path, totalCost, iter = algorithm.AStarSearch(board)
		found = len(path) > 0
	case AlgorithmIdaStar:
		path, totalCost, iter = algorithm.IDAStarSearch(board)
		found = len(path) > 0
	default:
		m.solverResult = &SolverResult{Found: false}
		m.solverPath = nil
		m.showPage(PageResult)
		return
	}

	// Build SolverResult steps from path
	steps := make([]SolverStep, 0)
	stepNum := 1
	i := 1
	for i < len(path) {
		prev := path[i-1]
		cur := path[i]

		var dir Direction
		if cur.Row < prev.Row {
			dir = Up
		} else if cur.Row > prev.Row {
			dir = Down
		} else if cur.Col < prev.Col {
			dir = Left
		} else {
			dir = Right
		}

		// hitung jumlah sel yang dilewati selama slide
		units := 0
		stepCost := 0
		for i < len(path) {
			p := path[i]
			pp := path[i-1]
			var d Direction
			if p.Row < pp.Row {
				d = Up
			} else if p.Row > pp.Row {
				d = Down
			} else if p.Col < pp.Col {
				d = Left
			} else {
				d = Right
			}
			if d != dir {
				break
			}

			// hitung jarak antar dua point
			dist := algorithm.Abs(p.Row-pp.Row) + algorithm.Abs(p.Col-pp.Col)
			units += dist

			for s := 1; s <= dist; s++ {
				var r, c int
				switch dir {
				case Up:
					r, c = pp.Row-s, pp.Col
				case Down:
					r, c = pp.Row+s, pp.Col
				case Left:
					r, c = pp.Row, pp.Col-s
				case Right:
					r, c = pp.Row, pp.Col+s
				}
				stepCost += board.Cost[r][c]
			}
			i++
		}

		steps = append(steps, SolverStep{
			StepNum:   stepNum,
			Direction: dir,
			Units:     units,
			Cost:      stepCost,
		})
		stepNum++
	}

	res := &SolverResult{
		Found:      found,
		Steps:      steps,
		TotalCost:  totalCost,
		TotalMoves: stepNum - 1,
		DurationMs: float64(time.Since(startedAt).Microseconds()) / 1000.0,
		Weight:     0,
		Level:      0,
		Seed:       "",
		Iterations: iter,
	}
	m.solverResult = res
	m.solverPath = path
	m.currentBoard = board
	m.showPage(PageResult)
}
