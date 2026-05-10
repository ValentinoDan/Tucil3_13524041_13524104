package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// SolverPage is the main Solver screen (Screen 3)
type SolverPage struct {
	main         *MainUI
	selectedAlgo SolverAlgorithm
	algoRadios   []*widget.Check
}

func NewSolverPage(m *MainUI) *SolverPage {
	return &SolverPage{
		main:         m,
		selectedAlgo: m.selectedAlgo,
	}
}

func (p *SolverPage) Build() fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorBackground)

	header := buildAppHeader(PageSolver, p.main.NavigateTo)
	leftPane := p.buildGridPane()
	rightPane := p.buildControlPane()

	// Divider between panes
	divider := canvas.NewRectangle(ColorSurfaceContainerHigh)
	divider.SetMinSize(fyne.NewSize(8, 1))

	body := container.NewHSplit(leftPane, rightPane)
	body.SetOffset(0.67)

	fullLayout := container.NewBorder(
		header,
		nil,
		nil, nil,
		body,
	)

	return container.NewStack(bg, fullLayout)
}

// Puzzle grid
func (p *SolverPage) buildGridPane() fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorWhite)

	grid := p.buildPuzzleGrid()

	scrolled := container.NewScroll(
		container.NewCenter(
			container.NewPadded(grid),
		),
	)

	return container.NewStack(bg, scrolled)
}

func (p *SolverPage) buildPuzzleGrid() fyne.CanvasObject {
	b := p.main.currentBoard
	if b == nil {
		// placeholder (no map loaded yet)
		placeholder := widget.NewLabelWithStyle("📋 Load a map to see the puzzle grid", fyne.TextAlignCenter, fyne.TextStyle{})
		return container.NewCenter(placeholder)
	}

	cardBg := canvas.NewRectangle(ColorWhite)
	cardBg.CornerRadius = 8
	cardBg.StrokeColor = ColorOutlineVariant
	cardBg.StrokeWidth = 1

	maxCells := b.M
	if b.N > maxCells {
		maxCells = b.N
	}
	var cs float32
	switch {
	case maxCells <= 5:
		cs = 60
	case maxCells <= 8:
		cs = 48
	case maxCells <= 12:
		cs = 36
	case maxCells <= 16:
		cs = 28
	default:
		cs = 22
	}
	cellSize := fyne.NewSize(cs, cs)
	gap := float32(2)

	cells := make([]fyne.CanvasObject, 0, b.M*b.N)
	for row := 0; row < b.N; row++ {
		for col := 0; col < b.M; col++ {
			cell := p.buildCell(b.Grid[row][col], cellSize)
			cells = append(cells, cell)
		}
	}

	gridLayout := container.New(newGridLayoutWithSize(b.M, cellSize, gap), cells...)

	gridWrapper := container.NewStack(
		container.NewGridWrap(
			fyne.NewSize(float32(b.M)*(cs+gap)+16, float32(b.N)*(cs+gap)+16),
			cardBg,
		),
		container.NewPadded(gridLayout),
	)

	return gridWrapper
}

func (p *SolverPage) buildCell(cellRune rune, size fyne.Size) fyne.CanvasObject {
	switch cellRune {
	case 'X': // Obstacle
		bg := canvas.NewRectangle(ColorObstacle)
		bg.CornerRadius = 3
		bg.StrokeColor = ColorObstacleBorder
		bg.StrokeWidth = 1
		return container.NewGridWrap(size, container.NewStack(bg))

	case 'L': // Lava
		bg := canvas.NewRectangle(ColorRedBg)
		bg.CornerRadius = 3
		bg.StrokeColor = ColorRedBorder
		bg.StrokeWidth = 1

		flame := canvas.NewText("⚠", ColorRedText)
		flame.TextSize = 12
		flame.TextStyle = fyne.TextStyle{Bold: true}

		return container.NewGridWrap(size, container.NewStack(bg, container.NewCenter(flame)))

	case 'Z': // Player
		bg := canvas.NewRectangle(ColorIceTile)
		bg.CornerRadius = 3
		bg.StrokeColor = ColorIceBorder
		bg.StrokeWidth = 1

		orb := canvas.NewCircle(ColorPlayerOrb)
		orb.StrokeColor = ColorPlayerBorder
		orb.StrokeWidth = 2

		return container.NewGridWrap(size, container.NewStack(bg, container.NewCenter(
			container.NewGridWrap(fyne.NewSize(size.Width*0.75, size.Height*0.75), orb),
		)))

	case 'O': // Goal
		bg := canvas.NewRectangle(ColorIceTile)
		bg.CornerRadius = 3
		bg.StrokeColor = ColorIceBorder
		bg.StrokeWidth = 1

		orb := canvas.NewCircle(ColorGoalOrb)
		orb.StrokeColor = ColorGoalBorder
		orb.StrokeWidth = 2

		flagText := canvas.NewText("F", ColorOnSurface)
		flagText.TextSize = 11
		flagText.TextStyle = fyne.TextStyle{Bold: true}

		return container.NewGridWrap(size, container.NewStack(bg, container.NewCenter(
			container.NewStack(
				container.NewGridWrap(fyne.NewSize(size.Width*0.75, size.Height*0.75), orb),
				container.NewCenter(flagText),
			),
		)))
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		bg := canvas.NewRectangle(ColorIceTile)
		bg.CornerRadius = 3
		bg.StrokeColor = ColorIceBorder
		bg.StrokeWidth = 1

		numberText := canvas.NewText(string(cellRune), ColorPrimary)
		numberText.TextSize = 14
		numberText.TextStyle = fyne.TextStyle{Bold: true}

		return container.NewGridWrap(size, container.NewStack(bg, container.NewCenter(numberText)))

	default: // Empty
		bg := canvas.NewRectangle(ColorIceTile)
		bg.CornerRadius = 3
		bg.StrokeColor = ColorIceBorder
		bg.StrokeWidth = 1
		return container.NewGridWrap(size, container.NewStack(bg))
	}
}

// Control panel
func (p *SolverPage) buildControlPane() fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorSurfaceContainerLow)
	borderLine := canvas.NewRectangle(ColorOutlineVariant)
	borderLine.SetMinSize(fyne.NewSize(1, 1))

	mapCard := p.buildMapConfigCard()
	engineCard := p.buildExecutionEngineCard()
	metricsCard := p.buildMetricsCard()

	scrollContent := container.NewVBox(
		vSpacer(8),
		mapCard,
		vSpacer(8),
		engineCard,
		vSpacer(8),
		metricsCard,
		vSpacer(16),
	)

	scrolled := container.NewScroll(container.NewPadded(scrollContent))

	return container.NewStack(bg, container.NewBorder(nil, nil, borderLine, nil, scrolled))
}

// Map config card
func (p *SolverPage) buildMapConfigCard() fyne.CanvasObject {
	cardBg := canvas.NewRectangle(ColorWhite)
	cardBg.CornerRadius = 8
	cardBg.StrokeColor = ColorOutlineVariant
	cardBg.StrokeWidth = 1

	// Title
	titleIcon := canvas.NewText("🗺", ColorOnSurfaceVariant)
	titleIcon.TextSize = 14
	titleText := canvas.NewText("MAP", ColorOnSurfaceVariant)
	titleText.TextSize = 11
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleRow := container.NewHBox(titleIcon, titleText)

	if p.main.selectedMap == nil {
		emptyText := canvas.NewText("No map selected", ColorSlate400)
		emptyText.TextSize = 12
		content := container.NewVBox(titleRow, vSpacer(8), container.NewCenter(emptyText))
		return container.NewStack(cardBg, container.NewPadded(content))
	}

	// Badge difficulty
	var badge fyne.CanvasObject
	switch p.main.selectedMap.Difficulty {
	case Easy:
		badge = easyBadge()
	case Intermediate:
		badge = intermediateBadge()
	case Hard:
		badge = hardBadge()
	}

	mapName := canvas.NewText(p.main.selectedMap.Filename, ColorSlate800)
	mapName.TextSize = 14
	mapName.TextStyle = fyne.TextStyle{Bold: true}

	dimText := canvas.NewText(
		fmt.Sprintf("%d × %d grid", p.main.selectedMap.Height, p.main.selectedMap.Width),
		ColorSlate500,
	)
	dimText.TextSize = 11

	nameRow := container.NewHBox(mapName, badge)
	sep := canvas.NewRectangle(ColorOutlineVariant)
	sep.SetMinSize(fyne.NewSize(0, 1))

	content := container.NewVBox(
		titleRow,
		sep,
		nameRow,
		dimText,
	)

	return container.NewStack(cardBg, container.NewPadded(content))
}

// Algo options card
func (p *SolverPage) buildExecutionEngineCard() fyne.CanvasObject {
	cardBg := canvas.NewRectangle(ColorWhite)
	cardBg.CornerRadius = 8
	cardBg.StrokeColor = ColorOutlineVariant
	cardBg.StrokeWidth = 1

	// Title
	titleIcon := canvas.NewText("⌨", ColorOnSurfaceVariant)
	titleIcon.TextSize = 14
	titleText := canvas.NewText("ALGORITHMS", ColorOnSurfaceVariant)
	titleText.TextSize = 11
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleRow := container.NewHBox(titleIcon, titleText)

	// Algorithm radio options
	algorithms := []SolverAlgorithm{
		AlgorithmUCS,
		AlgorithmGBFS,
		AlgorithmAStar,
		AlgorithmIdaStar,
	}

	algoItems := make([]fyne.CanvasObject, 0, len(algorithms))
	for _, algo := range algorithms {
		a := algo
		isSelected := a == p.selectedAlgo

		rowBg := canvas.NewRectangle(ColorWhite)
		rowBg.CornerRadius = 4
		rowBg.StrokeColor = ColorOutlineVariant
		rowBg.StrokeWidth = 1
		if isSelected {
			rowBg.StrokeColor = ColorPrimary
			rowBg.StrokeWidth = 2
		}

		// Radio indicator
		radioOuter := canvas.NewCircle(ColorOutlineVariant)
		if isSelected {
			radioOuter = canvas.NewCircle(ColorPrimary)
		}
		radioInner := canvas.NewCircle(ColorTransparent)
		if isSelected {
			radioInner = canvas.NewCircle(ColorWhite)
		}

		radioWidget := container.NewStack(
			container.NewGridWrap(fyne.NewSize(16, 16), radioOuter),
			container.NewCenter(container.NewGridWrap(fyne.NewSize(8, 8), radioInner)),
		)

		algoLabel := canvas.NewText(a.String(), ColorOnSurface)
		algoLabel.TextSize = 13

		btn := widget.NewButton("", func() {
			p.selectedAlgo = a
			p.main.selectedAlgo = a
			p.main.showPage(PageSolver)
		})
		btn.Importance = widget.LowImportance

		rowContent := container.NewPadded(container.NewHBox(radioWidget, hSpacer(8), algoLabel))
		row := newTappableRow(container.NewStack(rowBg, rowContent), func() {
			p.selectedAlgo = a
			p.main.selectedAlgo = a
			p.main.showPage(PageSolver)
		})

		algoItems = append(algoItems, row)
	}

	// Start solver button
	startBtn := widget.NewButton("▶  START SOLVER", func() {
		p.main.RunSolver()
	})
	startBtn.Importance = widget.HighImportance

	algoList := container.NewVBox(algoItems...)

	sep := canvas.NewRectangle(ColorOutlineVariant)
	sep.SetMinSize(fyne.NewSize(0, 1))

	content := container.NewVBox(
		titleRow,
		sep,
		algoList,
		vSpacer(6),
		startBtn,
	)

	return container.NewStack(cardBg, container.NewPadded(content))
}

// Metric card
func (p *SolverPage) buildMetricsCard() fyne.CanvasObject {
	cardBg := canvas.NewRectangle(ColorWhite)
	cardBg.CornerRadius = 8
	cardBg.StrokeColor = ColorOutlineVariant
	cardBg.StrokeWidth = 1

	result := p.main.solverResult
	if result == nil {
		result = &SolverResult{}
	}
	durationBlock := statBlock("DURATION", fmt.Sprintf("%.3fms", result.DurationMs))
	stepsBlock := statBlock("STEPS", fmt.Sprintf("%d", result.TotalMoves))
	costBlock := statBlock("COST", fmt.Sprintf("%d", result.TotalCost))

	sep1 := canvas.NewRectangle(ColorOutlineVariant)
	sep1.SetMinSize(fyne.NewSize(1, 40))
	sep2 := canvas.NewRectangle(ColorOutlineVariant)
	sep2.SetMinSize(fyne.NewSize(1, 40))

	metricsRow := container.NewHBox(
		container.NewCenter(durationBlock),
		sep1,
		container.NewCenter(stepsBlock),
		sep2,
		container.NewCenter(costBlock),
	)

	return container.NewStack(cardBg, container.NewPadded(metricsRow))
}

type gridLayoutWithSize struct {
	cols     int
	cellSize fyne.Size
	gap      float32
}

func newGridLayoutWithSize(cols int, cellSize fyne.Size, gap float32) fyne.Layout {
	return &gridLayoutWithSize{cols: cols, cellSize: cellSize, gap: gap}
}

func (g *gridLayoutWithSize) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	col := 0
	row := 0
	for _, obj := range objects {
		x := float32(col) * (g.cellSize.Width + g.gap)
		y := float32(row) * (g.cellSize.Height + g.gap)
		obj.Move(fyne.NewPos(x, y))
		obj.Resize(g.cellSize)
		col++
		if col >= g.cols {
			col = 0
			row++
		}
	}
}

func (g *gridLayoutWithSize) MinSize(objects []fyne.CanvasObject) fyne.Size {
	rows := (len(objects) + g.cols - 1) / g.cols
	w := float32(g.cols)*(g.cellSize.Width+g.gap) - g.gap
	h := float32(rows)*(g.cellSize.Height+g.gap) - g.gap
	return fyne.NewSize(w, h)
}
