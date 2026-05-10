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
	board         *puzzle.Board
	pathPoints    []puzzle.Point // start to goal
	state         PbackState
	currStep      int
	totalSteps    int
	speed         time.Duration // ms
	lastStepTime  time.Time
	pausedElapsed time.Duration
	finishTime    time.Time

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

	p.finishTime = time.Time{}
	p.state = PbackPlaying
	p.lastStepTime = time.Now()

	p.onStateChange(PbackPlaying)
}

func (p *Pback) Pause() {
	if p.state != PbackPlaying {
		return
	}

	p.pausedElapsed = time.Since(p.lastStepTime)
	p.state = PbackPaused
	p.onStateChange(PbackPaused)
}

func (p *Pback) Stop() {
	p.state = PbackStopped
	p.currStep = 0
	p.lastStepTime = time.Time{}
	p.pausedElapsed = 0
	p.finishTime = time.Time{}
	p.onStepChange(0)
	p.onStateChange(PbackStopped)
}

func (p *Pback) NextStep() {
	if p.currStep < p.totalSteps {
		p.currStep++
		p.lastStepTime = time.Now()
		p.onStepChange(p.currStep)
	}
}

func (p *Pback) PrevStep() {
	if p.currStep > 0 {
		p.currStep--
		p.lastStepTime = time.Now()
		p.onStepChange(p.currStep)
	}
}

func (p *Pback) GoToStep(step int) {
	if step >= 0 && step <= p.totalSteps {
		p.currStep = step
		p.lastStepTime = time.Now() // reset timer
		p.finishTime = time.Time{}
		p.onStepChange(step)
	}
}

func (p *Pback) Update() bool {
	if p.state != PbackPlaying {
		return false
	}

	if p.currStep >= p.totalSteps {
		// delay 1s
		if p.finishTime.IsZero() {
			p.finishTime = time.Now()
		}
		if time.Since(p.finishTime) >= 1*time.Second {
			p.Stop()
		}
		return false
	}

	// calc time
	elapsed := time.Since(p.lastStepTime)
	if elapsed >= p.speed {
		p.lastStepTime = time.Now()
		p.currStep++
		p.onStepChange(p.currStep)
	}

	return true
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
