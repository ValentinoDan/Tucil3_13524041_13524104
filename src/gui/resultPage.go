package gui

import (
	"fmt"
	"time"

	"iceSlidingPuzzle/src/filehandler"
	"iceSlidingPuzzle/src/puzzle"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type ResultPage struct {
	main  *MainUI
	pback *Pback
}

func NewResultPage(m *MainUI) *ResultPage {
	p := &ResultPage{main: m}
	if m.currentBoard != nil && len(m.solverPath) > 1 {
		p.pback = NewPback(m.currentBoard, m.solverPath)
	}
	return p
}

func (p *ResultPage) Build() fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorBackground)

	header := buildAppHeader(PageResult, p.main.NavigateTo)
	leftPane := p.buildVisualizerPane()
	rightPane := p.buildExecutionLogPane()

	split := container.NewHSplit(leftPane, rightPane)
	split.SetOffset(0.60)

	fullLayout := container.NewBorder(
		header,
		nil,
		nil, nil,
		split,
	)

	return container.NewStack(bg, fullLayout)
}

func (p *ResultPage) buildVisualizerPane() fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorWhite)
	rightBorder := canvas.NewRectangle(ColorOutlineVariant)
	rightBorder.SetMinSize(fyne.NewSize(1, 1))

	// Pane header
	paneBg := canvas.NewRectangle(ColorWhite)
	bottomLine := canvas.NewRectangle(ColorSlate100)
	bottomLine.SetMinSize(fyne.NewSize(1, 1))

	titleText := canvas.NewText("Visualizer", ColorSlate800)
	titleText.TextSize = 13
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	// Solution found
	badgeBg := canvas.NewRectangle(ColorTertiaryBg)
	badgeBg.CornerRadius = 20
	badgeBg.StrokeColor = ColorTertiaryBorder
	badgeBg.StrokeWidth = 1
	checkIcon := canvas.NewText("✓", ColorTertiary)
	checkIcon.TextSize = 12
	checkIcon.TextStyle = fyne.TextStyle{Bold: true}
	badgeText := canvas.NewText("SOLUTION FOUND", ColorTertiary)
	badgeText.TextSize = 10
	badgeText.TextStyle = fyne.TextStyle{Bold: true}
	badgeContent := container.NewHBox(checkIcon, badgeText)
	badge := container.NewStack(
		container.NewGridWrap(fyne.NewSize(160, 26), badgeBg),
		container.NewCenter(badgeContent),
	)
	if p.main.solverResult == nil || !p.main.solverResult.Found {
		badge = container.NewStack(
			container.NewGridWrap(fyne.NewSize(160, 26), canvas.NewRectangle(ColorRedBg)),
			container.NewCenter(canvas.NewText("NO SOLUTION", ColorRedText)),
		)
	}

	paneHeader := container.NewStack(
		paneBg,
		container.NewVBox(
			container.NewPadded(container.NewBorder(nil, nil, titleText, badge)),
			bottomLine,
		),
	)

	// Grid visualization area
	vizArea := p.buildVisualizationArea()

	// Footer stats
	footer := p.buildVisualizerFooter()

	content := container.NewBorder(
		paneHeader,
		footer,
		nil, nil,
		vizArea,
	)

	return container.NewStack(bg, rightBorder, content)
}

func (p *ResultPage) buildVisualizationArea() fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorSlate50)
	b := p.main.currentBoard
	if b == nil {
		return container.NewStack(bg, container.NewCenter(widget.NewLabelWithStyle("Load a map first", fyne.TextAlignCenter, fyne.TextStyle{})))
	}

	maxCells := b.M
	if b.N > maxCells {
		maxCells = b.N
	}

	var cellSize float32
	switch {
	case maxCells <= 5:
		cellSize = 60
	case maxCells <= 8:
		cellSize = 48
	case maxCells <= 12:
		cellSize = 36
	case maxCells <= 16:
		cellSize = 28
	default:
		cellSize = 22
	}
	cellGap := float32(2)
	pad := float32(12)
	boardWidth := pad*2 + float32(b.M)*cellSize + float32(b.M-1)*cellGap
	boardHeight := pad*2 + float32(b.N)*cellSize + float32(b.N-1)*cellGap

	board := container.NewWithoutLayout()
	for row := 0; row < b.N; row++ {
		for col := 0; col < b.M; col++ {
			cell := p.makeBoardCell(b.Grid[row][col], fyne.NewSize(cellSize, cellSize))
			x := pad + float32(col)*(cellSize+cellGap)
			y := pad + float32(row)*(cellSize+cellGap)
			cell.Move(fyne.NewPos(x, y))
			board.Add(cell)
		}
	}

	board.Resize(fyne.NewSize(boardWidth, boardHeight))

	pathOverlay := container.NewWithoutLayout()
	boardStack := container.NewStack(canvas.NewRectangle(ColorSlate50), board, pathOverlay)
	boardStack.Resize(fyne.NewSize(boardWidth, boardHeight))
	fixedBoard := container.NewGridWrap(fyne.NewSize(boardWidth, boardHeight), boardStack)

	// Playback controls
	controls := p.buildPlaybackControls(pathOverlay, cellSize, pad, cellGap)

	// Label
	result := p.main.solverResult
	if result == nil {
		result = &SolverResult{}
	}
	pathTitle := canvas.NewText("PATH TRAVERSAL", ColorSlate800)
	pathTitle.TextSize = 12
	pathTitle.TextStyle = fyne.TextStyle{Bold: true}
	pathTitle.Alignment = fyne.TextAlignCenter
	pathSub := canvas.NewText(fmt.Sprintf("%d steps detected", result.TotalMoves), ColorSlate500)
	pathSub.TextSize = 11
	pathSub.Alignment = fyne.TextAlignCenter

	labelArea := container.NewVBox(
		container.NewCenter(pathTitle),
		container.NewCenter(pathSub),
	)

	vizContent := container.NewCenter(
		container.NewVBox(
			fixedBoard,
			vSpacer(12),
			controls,
			vSpacer(8),
			labelArea,
			vSpacer(12),
		),
	)

	return container.NewStack(bg, container.NewScroll(vizContent))
}

func (p *ResultPage) buildPlaybackControls(overlay *fyne.Container, cellSize, pad, gap float32) fyne.CanvasObject {
	if p.pback == nil {
		return container.NewCenter(canvas.NewText("No path to play", ColorSlate400))
	}

	pb := p.pback

	// Progress bar
	progressBar := canvas.NewRectangle(ColorSlate200)
	progressBar.SetMinSize(fyne.NewSize(300, 4))
	progressFill := canvas.NewRectangle(ColorPrimary)
	progressFill.Resize(fyne.NewSize(0, 4))

	progressTrack := container.NewWithoutLayout()
	progressTrack.Add(progressBar)
	progressTrack.Add(progressFill)
	progressBar.Resize(fyne.NewSize(300, 4))
	progressTrack.Resize(fyne.NewSize(300, 4))

	// Step label
	stepLabel := canvas.NewText("0 / 0", ColorSlate500)
	stepLabel.TextSize = 10

	updateOverlay := func(step int) {
		if step < 0 || step >= len(pb.pathPoints) {
			return
		}
		overlay.Objects = p.buildPathObjects(pb.pathPoints[:step+1], cellSize, pad, gap)
		overlay.Refresh()

		prog := float32(pb.GetProgress())
		progressFill.Resize(fyne.NewSize(300*prog, 4))
		progressFill.Refresh()

		stepLabel.Text = fmt.Sprintf("%d / %d", pb.GetcurrStep(), pb.GetTotalSteps())
		stepLabel.Refresh()
	}

	// Balik ke awal
	prevBtn := widget.NewButton("⏮", func() {
		pb.PrevStep()
		updateOverlay(pb.GetcurrStep())
	})
	prevBtn.Importance = widget.LowImportance

	playBtn := widget.NewButton("▶", func() {
		if pb.GetState() == PbackPlaying {
			pb.Pause()
		} else {
			pb.Play()
		}
	})
	playBtn.Importance = widget.HighImportance

	// Maju ke akhir
	nextBtn := widget.NewButton("⏭", func() {
		pb.NextStep()
		updateOverlay(pb.GetcurrStep())
	})
	nextBtn.Importance = widget.LowImportance

	stopBtn := widget.NewButton("⏹", func() {
		pb.Stop()
	})
	stopBtn.Importance = widget.LowImportance

	skipBackBtn := widget.NewButton("«", func() {
		pb.GoToStep(0)
		updateOverlay(0)
	})
	skipBackBtn.Importance = widget.LowImportance

	skipFwdBtn := widget.NewButton("»", func() {
		target := pb.GetTotalSteps()
		pb.GoToStep(target)
		updateOverlay(target)
	})
	skipFwdBtn.Importance = widget.LowImportance

	// Set callbacks sekali saja
	pb.SetCallbacks(func(step int) {
		fyne.Do(func() { updateOverlay(step) })
	}, func(state PbackState) {
		fyne.Do(func() {
			if state == PbackPlaying {
				playBtn.SetText("⏸")
			} else {
				playBtn.SetText("▶")
			}
			playBtn.Refresh()
		})
	})

	// Ticker untuk drive Pback.Update()
	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			fyne.Do(func() { pb.Update() })
		}
	}()

	// Play 1x dulu setelah build
	updateOverlay(0)
	pb.Play()

	progressRow := container.NewGridWrap(fyne.NewSize(300, 4), progressTrack)

	speedLabel := canvas.NewText("Speed: 250ms", ColorSlate500)
	speedLabel.TextSize = 10

	speedSlider := widget.NewSlider(50, 1000)
	speedSlider.Value = 750
	speedSlider.Step = 50
	speedSlider.OnChanged = func(v float64) {
		delay := int(1050 - v)
		pb.SetSpeed(delay)
		speedLabel.Text = fmt.Sprintf("Speed: %dms", delay)
		speedLabel.Refresh()
	}

	speedRow := container.NewVBox(
		container.NewCenter(speedLabel),
		container.NewGridWrap(fyne.NewSize(300, 20), speedSlider),
	)

	return container.NewCenter(
		container.NewVBox(
			container.NewCenter(progressRow),
			vSpacer(6),
			container.NewCenter(
				container.NewHBox(stopBtn, skipBackBtn, prevBtn, playBtn, nextBtn, skipFwdBtn),
			),
			vSpacer(4),
			container.NewCenter(stepLabel),
			speedRow,
		),
	)
}

func (p *ResultPage) makeBoardCell(cellRune rune, size fyne.Size) fyne.CanvasObject {
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
		icon := canvas.NewText("🔥", ColorRedText)
		icon.TextSize = 12
		icon.TextStyle = fyne.TextStyle{Bold: true}
		return container.NewGridWrap(size, container.NewStack(bg, container.NewCenter(icon)))
	case 'Z': // Player
		bg := canvas.NewRectangle(ColorPrimaryLight)
		bg.CornerRadius = 3
		bg.StrokeColor = ColorPrimary
		bg.StrokeWidth = 1
		orb := canvas.NewCircle(ColorPrimary)
		return container.NewGridWrap(size, container.NewStack(bg, container.NewCenter(container.NewGridWrap(fyne.NewSize(size.Width*0.5, size.Height*0.5), orb))))
	case 'O': // Goal
		bg := canvas.NewRectangle(ColorTertiaryBg)
		bg.CornerRadius = 3
		bg.StrokeColor = ColorTertiary
		bg.StrokeWidth = 1.5
		checkText := canvas.NewText("✓", ColorTertiary)
		checkText.TextSize = 22
		checkText.TextStyle = fyne.TextStyle{Bold: true}
		return container.NewGridWrap(size, container.NewStack(bg, container.NewCenter(checkText)))
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		bg := canvas.NewRectangle(ColorWhite)
		bg.StrokeColor = ColorSlate200
		bg.StrokeWidth = 1
		numberText := canvas.NewText(string(cellRune), ColorPrimary)
		numberText.TextSize = 14
		numberText.TextStyle = fyne.TextStyle{Bold: true}
		return container.NewGridWrap(size, container.NewStack(bg, container.NewCenter(numberText)))
	default: // Empty
		bg := canvas.NewRectangle(ColorWhite)
		bg.StrokeColor = ColorSlate200
		bg.StrokeWidth = 1
		return container.NewGridWrap(size, container.NewStack(bg))
	}
}

func (p *ResultPage) buildPathObjects(path []puzzle.Point, cellSize float32, pad, gap float32) []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0, len(path)*2)
	for i := 1; i < len(path); i++ {
		prev := path[i-1]
		cur := path[i]
		prevX := pad + float32(prev.Col)*(cellSize+gap) + cellSize/2
		prevY := pad + float32(prev.Row)*(cellSize+gap) + cellSize/2
		curX := pad + float32(cur.Col)*(cellSize+gap) + cellSize/2
		curY := pad + float32(cur.Row)*(cellSize+gap) + cellSize/2

		seg := canvas.NewRectangle(ColorPrimaryContainer)
		if prev.Row == cur.Row {
			x := prevX
			if curX < prevX {
				x = curX
			}
			seg.Move(fyne.NewPos(x, prevY-2))
			seg.Resize(fyne.NewSize(absFloat(curX-prevX), 4))
		} else {
			y := prevY
			if curY < prevY {
				y = curY
			}
			seg.Move(fyne.NewPos(prevX-2, y))
			seg.Resize(fyne.NewSize(4, absFloat(curY-prevY)))
		}
		objects = append(objects, seg)
	}

	if len(path) > 0 {
		last := path[len(path)-1]
		dot := canvas.NewCircle(ColorPrimary)
		dot.Resize(fyne.NewSize(16, 16))
		dot.Move(fyne.NewPos(
			pad+float32(last.Col)*(cellSize+gap)+cellSize/2-8,
			pad+float32(last.Row)*(cellSize+gap)+cellSize/2-8,
		))
		objects = append(objects, dot)
	}

	return objects
}

func absFloat(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func (p *ResultPage) buildVisualizerFooter() fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorHeaderBg)
	topLine := canvas.NewRectangle(ColorSlate200)
	topLine.SetMinSize(fyne.NewSize(1, 1))

	result := p.main.solverResult
	if result == nil {
		result = &SolverResult{}
	}

	movesLabel := canvas.NewText("TOTAL MOVES", ColorSlate500)
	movesLabel.TextSize = 9
	movesLabel.TextStyle = fyne.TextStyle{Bold: true}
	movesVal := canvas.NewText(fmt.Sprintf("%d Steps", result.TotalMoves), ColorPrimary)
	movesVal.TextSize = 13
	movesVal.TextStyle = fyne.TextStyle{Bold: true}
	movesBlock := container.NewVBox(movesLabel, movesVal)

	costLabel := canvas.NewText("TOTAL COST", ColorSlate500)
	costLabel.TextSize = 9
	costLabel.TextStyle = fyne.TextStyle{Bold: true}
	costVal := canvas.NewText(fmt.Sprintf("%d Credits", result.TotalCost), ColorPrimary)
	costVal.TextSize = 13
	costVal.TextStyle = fyne.TextStyle{Bold: true}
	costBlock := container.NewVBox(costLabel, costVal)

	timeLabel := canvas.NewText("TOTAL TIME", ColorSlate500)
	timeLabel.TextSize = 9
	timeLabel.TextStyle = fyne.TextStyle{Bold: true}
	timeVal := canvas.NewText(fmt.Sprintf("%.3f ms", result.DurationMs), ColorPrimary)
	timeVal.TextSize = 13
	timeVal.TextStyle = fyne.TextStyle{Bold: true}
	timeBlock := container.NewVBox(timeLabel, timeVal)

	footerContent := container.NewPadded(
		container.NewHBox(movesBlock, hSpacer(24), costBlock, hSpacer(24), timeBlock),
	)

	return container.NewStack(
		bg,
		container.NewVBox(topLine, footerContent),
	)
}

// Execution log
func (p *ResultPage) buildExecutionLogPane() fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorWhite)

	// Header
	headerBg := canvas.NewRectangle(ColorWhite)
	bottomLine := canvas.NewRectangle(ColorSlate100)
	bottomLine.SetMinSize(fyne.NewSize(1, 1))

	titleText := canvas.NewText("Execution Log", ColorSlate800)
	titleText.TextSize = 13
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	histIcon := canvas.NewText("🕐", ColorSlate400)
	histIcon.TextSize = 14

	paneHeader := container.NewStack(
		headerBg,
		container.NewVBox(
			container.NewPadded(container.NewBorder(nil, nil, titleText, histIcon)),
			bottomLine,
		),
	)

	// Table header
	tableHeaderBg := canvas.NewRectangle(ColorSlate50)
	stepHdr := canvas.NewText("STEP", ColorSlate500)
	stepHdr.TextSize = 10
	stepHdr.TextStyle = fyne.TextStyle{Bold: true}
	actionHdr := canvas.NewText("ACTION", ColorSlate500)
	actionHdr.TextSize = 10
	actionHdr.TextStyle = fyne.TextStyle{Bold: true}

	tableHeader := container.NewStack(
		tableHeaderBg,
		container.NewPadded(container.NewGridWithColumns(3, stepHdr, actionHdr)),
	)

	// Table rows
	result := p.main.solverResult
	if result == nil {
		result = &SolverResult{}
	}
	rows := make([]fyne.CanvasObject, 0, len(result.Steps))

	for _, step := range result.Steps {
		s := step
		rowBg := canvas.NewRectangle(ColorWhite)
		bottomBorder := canvas.NewRectangle(ColorSlate100)
		bottomBorder.SetMinSize(fyne.NewSize(1, 1))

		stepNum := canvas.NewText(fmt.Sprintf("%02d", s.StepNum), ColorSlate400)
		stepNum.TextSize = 12
		stepNum.TextStyle = fyne.TextStyle{Monospace: true}

		actionName := canvas.NewText(s.Direction.String(), ColorSlate800)
		actionName.TextSize = 12
		actionName.TextStyle = fyne.TextStyle{Bold: true}
		actionUnit := canvas.NewText(
			fmt.Sprintf("%d tile%s", s.Tiles, pluralS(s.Tiles)),
			ColorSlate500,
		)
		actionUnit.TextSize = 10
		actionCol := container.NewVBox(actionName, actionUnit)

		rowContent := container.NewPadded(
			container.NewGridWithColumns(3, stepNum, actionCol),
		)

		row := container.NewVBox(
			container.NewStack(rowBg, rowContent),
			bottomBorder,
		)
		rows = append(rows, row)
	}

	tableBody := container.NewVBox(rows...)
	scrollableTable := container.NewScroll(tableBody)

	saveBtn := widget.NewButton("💾  SAVE OUTPUT", func() {
		if p.main.solverResult == nil || !p.main.solverResult.Found {
			dialog.ShowError(fmt.Errorf("No solution to save"), p.main.window)
			return
		}
		fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, p.main.window)
				return
			}
			if writer == nil {
				return
			}
			filePath := writer.URI().Path()
			writer.Close()
			err = filehandler.SaveOutputTxt(p.main.currentBoard, p.main.solverPath, p.main.solverResult.TotalCost, p.main.solverResult.Iterations, p.main.solverResult.DurationMs, filePath)
			if err != nil {
				dialog.ShowError(err, p.main.window)
			} else {
				dialog.ShowInformation("Success", "Output saved successfully", p.main.window)
			}
		}, p.main.window)
		fd.SetFileName("output.txt")
		fd.Show()
	})
	saveBtn.Importance = widget.HighImportance
	backBtn := widget.NewButton("📊  BACK TO DASHBOARD", func() {
		p.main.NavigateTo(PageSolver)
	})
	backBtn.Importance = widget.HighImportance
	btnGrid := container.NewGridWithColumns(2, saveBtn, backBtn)
	actionArea := container.NewStack(
		canvas.NewRectangle(ColorHeaderBg),
		container.NewPadded(btnGrid),
	)

	content := container.NewBorder(
		container.NewVBox(paneHeader, tableHeader),
		actionArea,
		nil, nil,
		scrollableTable,
	)

	return container.NewStack(bg, content)
}

func pluralS(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
