package gui

import (
	"iceSlidingPuzzle/src/puzzle"
	"time"
)

// Cell type
const (
	PBCellEmpty = iota
	PBCellWall
	PBCellLava
	PBCellPoint
	PBCellGoal
	PBCellPlayer
)

// tracks the Curr Pback state
type PbackState int

const (
	PbackStopped PbackState = iota
	PbackPlaying
	PbackPaused
)

// handles solution animation and Pback control
type Pback struct {
	// Board and path
	board       *puzzle.Board
	pathPoints  []puzzle.Point // start to goal
	state       PbackState
	currStep    int
	totalSteps  int
	speed       time.Duration // ms
	startTime   time.Time
	pausedTime  time.Time
	totalPaused time.Duration

	onStepChange  func(int)        // called when step changes
	onStateChange func(PbackState) // called when play/pause/stop state changes
}

// new Pback controller
func NewPback(board *puzzle.Board, pathPoints []puzzle.Point) *Pback {
	return &Pback{
		board:         board,
		pathPoints:    pathPoints,
		state:         PbackStopped,
		currStep:      0,
		totalSteps:    len(pathPoints) - 1,
		speed:         250 * time.Millisecond, // default 250ms per step
		onStepChange:  func(int) {},
		onStateChange: func(PbackState) {},
	}
}

func (p *Pback) SetSpeed(ms int) {
	p.speed = time.Duration(ms) * time.Millisecond
}

func (p *Pback) SetCallbacks(onStep func(int), onStateChange func(PbackState)) {
	if onStep != nil {
		p.onStepChange = onStep
	}

	if onStateChange != nil {
		p.onStateChange = onStateChange
	}
}

func (p *Pback) Play() {
	if p.state == PbackPlaying {
		return
	}

	if p.currStep >= p.totalSteps {
		p.currStep = 0 // restart
	}

	p.state = PbackPlaying
	p.startTime = time.Now()
	if !p.pausedTime.IsZero() {
		p.totalPaused += time.Since(p.pausedTime)
		p.pausedTime = time.Time{}
	}

	p.onStateChange(PbackPlaying)
}

func (p *Pback) Pause() {
	if p.state != PbackPlaying {
		return
	}

	p.state = PbackPaused
	p.pausedTime = time.Now()
	p.onStateChange(PbackPaused)
}

func (p *Pback) Stop() {
	p.state = PbackStopped
	p.currStep = 0
	p.startTime = time.Time{}
	p.pausedTime = time.Time{}
	p.totalPaused = 0
	p.onStepChange(0)
	p.onStateChange(PbackStopped)
}

func (p *Pback) NextStep() {
	if p.currStep < p.totalSteps {
		p.currStep++
		p.onStepChange(p.currStep)
	}
}

func (p *Pback) PrevStep() {
	if p.currStep > 0 {
		p.currStep--
		p.onStepChange(p.currStep)
	}
}

func (p *Pback) GoToStep(step int) {
	if step >= 0 && step <= p.totalSteps {
		p.currStep = step
		p.onStepChange(step)
	}
}

func (p *Pback) Update() bool {
	if p.state != PbackPlaying {
		return false
	}

	if p.currStep >= p.totalSteps {
		p.Stop()
		return false
	}

	// calc time
	elapsed := time.Since(p.startTime) - p.totalPaused
	expectedStep := int(elapsed / p.speed)
	targetStep := expectedStep

	if targetStep > p.totalSteps {
		targetStep = p.totalSteps
	}

	if targetStep != p.currStep {
		p.currStep = targetStep
		p.onStepChange(p.currStep)

		if p.currStep >= p.totalSteps {
			p.Stop()
			return false
		}
	}

	return true
}

// returns the player position at curr step
func (p *Pback) GetCurrPosition() puzzle.Point {
	if p.currStep < len(p.pathPoints) {
		return p.pathPoints[p.currStep]
	}
	return puzzle.Point{Row: -1, Col: -1}
}

// returns the grid with player pos at curr step
func (p *Pback) GetCurrGrid() [][]int {
	grid := make([][]int, p.board.N)
	for r := 0; r < p.board.N; r++ {
		grid[r] = make([]int, p.board.M)
		for c := 0; c < p.board.M; c++ {
			ch := p.board.Grid[r][c]
			switch ch {
			case 'X':
				grid[r][c] = PBCellWall
			case 'L':
				grid[r][c] = PBCellLava
			case 'O':
				grid[r][c] = PBCellGoal
			case 'Z':
				grid[r][c] = PBCellPlayer
			default:
				if ch >= '0' && ch <= '9' {
					grid[r][c] = PBCellPoint
				} else {
					grid[r][c] = PBCellEmpty
				}
			}
		}
	}

	// Place player at curr position
	pos := p.GetCurrPosition()
	if pos.Row >= 0 && pos.Row < p.board.N && pos.Col >= 0 && pos.Col < p.board.M {
		grid[pos.Row][pos.Col] = PBCellPlayer
	}

	return grid
}

// Getter Funcs
func (p *Pback) GetState() PbackState {
	return p.state
}

func (p *Pback) GetcurrStep() int {
	return p.currStep
}

func (p *Pback) GetTotalSteps() int {
	return p.totalSteps
}

func (p *Pback) GetProgress() float64 {
	if p.totalSteps == 0 {
		return 0
	}
	return float64(p.currStep) / float64(p.totalSteps)
}

func (p *Pback) IsFinished() bool {
	return p.currStep >= p.totalSteps
}
