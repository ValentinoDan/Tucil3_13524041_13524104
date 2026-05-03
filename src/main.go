package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"iceSlidingPuzzle/src/algorithm"
	"iceSlidingPuzzle/src/filehandler"
	"iceSlidingPuzzle/src/gui"
	"iceSlidingPuzzle/src/puzzle"
)

func main() {
	myApp := app.NewWithID("com.arcticsolver.app")
	myApp.Settings().SetTheme(theme.LightTheme()) // Paksa Light Theme agar text berwarna hitam (kontras dengan background biru)

	myWindow := myApp.NewWindow("Arctic Solver")
	myWindow.Resize(fyne.NewSize(450, 750)) // Ukuran ala mobile app
	// myWindow.SetFixedSize(true) // Kunci ukuran agar tidak dimaximize (mencegah layout hancur)

	var currentBoard *puzzle.Board
	var boardCanvas *gui.BoardCanvas
	var tabs *container.AppTabs

	// Background color (Soft Ice Blue)
	bgColor := color.NRGBA{R: 215, G: 235, B: 245, A: 255}
	bgRect := canvas.NewRectangle(bgColor)

	// ==========================================
	// SHARED UI ELEMENTS
	// ==========================================
	statusLabel := widget.NewLabelWithStyle("No Map Loaded", fyne.TextAlignCenter, fyne.TextStyle{})
	statsLabel := widget.NewLabelWithStyle("--", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	canvasContainer := container.NewStack()

	logList := container.NewVBox()
	logScroll := container.NewVScroll(logList)

	// ==========================================
	// SOLVE BUTTON LOGIC
	// ==========================================
	var solveBtn *widget.Button
	solveBtn = widget.NewButtonWithIcon("SOLVE", theme.MediaPlayIcon(), func() {
		if currentBoard == nil {
			dialog.ShowError(fmt.Errorf("Please load a map first from Library"), myWindow)
			return
		}

		statusLabel.SetText("Solving...")
		statsLabel.SetText("--")
		logList.Objects = nil // clear execution log
		solveBtn.Disable()    // Mencegah double-click

		// Jalankan UCS di Goroutine agar UI tidak Freeze!
		go func() {
			start := time.Now()
			// Panggil algoritma
			goalNode, iter := algorithm.UCS(currentBoard)
			duration := time.Since(start)

			if goalNode == nil {
				statusLabel.SetText("No Solution Found!")
				statsLabel.SetText(fmt.Sprintf("Iterations: %d | Time: %d ms", iter, duration.Milliseconds()))
				solveBtn.Enable()
				return
			}

			// Reconstruct path
			var pathNodes []*puzzle.Node
			curr := goalNode
			for curr != nil {
				pathNodes = append([]*puzzle.Node{curr}, pathNodes...)
				curr = curr.Parent
			}

			var points []puzzle.Point
			for _, n := range pathNodes {
				points = append(points, n.State.Pos)
			}

			// Populate Execution Log in the Stats tab
			logList.Add(widget.NewLabelWithStyle("Execution Log:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
			for i, n := range pathNodes {
				if i == 0 {
					continue
				}
				stepCost := n.Cost - pathNodes[i-1].Cost
				logItem := fmt.Sprintf("[%d] Move %s\n      Cost: +%d | Sliding to [%d, %d]", i, n.Dir, stepCost, n.State.Pos.Row, n.State.Pos.Col)

				card := widget.NewCard(fmt.Sprintf("Step %d", i), "", widget.NewLabel(logItem))
				logList.Add(card)
			}

			statusLabel.SetText("✅ GOAL REACHED!")
			statsLabel.SetText(fmt.Sprintf("Time: %d ms | Iterations: %d | Total Cost: %d", duration.Milliseconds(), iter, goalNode.Cost))

			// Animate
			if boardCanvas != nil {
				boardCanvas.AnimatePath(points)
			}

			solveBtn.Enable()
		}()
	})
	solveBtn.Disable()

	// ==========================================
	// TAB 1: LIBRARY PAGE
	// ==========================================
	importBtn := widget.NewButtonWithIcon("Import New .txt Map", theme.FolderOpenIcon(), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			filename := reader.URI().Path()
			board, err := filehandler.ParseBoard(filename)
			if err != nil {
				dialog.ShowError(fmt.Errorf("failed parsing board: %v", err), myWindow)
				return
			}

			currentBoard = board
			statusLabel.SetText(fmt.Sprintf("Loaded: %s", reader.URI().Name()))
			statsLabel.SetText("--")
			logList.Objects = nil

			// Draw Map Grid dengan ukuran yang dinamis
			maxWidth := float32(380)
			maxHeight := float32(380)
			
			gridSizeX := maxWidth / float32(currentBoard.M)
			gridSizeY := maxHeight / float32(currentBoard.N)
			
			gridSize := gridSizeX
			if gridSizeY < gridSize {
				gridSize = gridSizeY
			}
			// Batasi ukuran maksimal cell agar map kecil (contoh: 5x5) tidak terlihat raksasa
			if gridSize > 60 {
				gridSize = 60
			}

			boardCanvas = gui.NewBoardCanvas(currentBoard, gridSize)
			
			// Bungkus grid dengan Scroll + Center (Sekarang otomatis berada di tengah!)
			scrollGrid := container.NewScroll(container.NewCenter(boardCanvas.Container))
			canvasContainer.Objects = []fyne.CanvasObject{scrollGrid}
			canvasContainer.Refresh()

			solveBtn.Enable()

			// Auto-switch to Solver Tab
			tabs.SelectIndex(1)

		}, myWindow)
	})

	libraryContent := container.NewVBox(
		widget.NewLabelWithStyle("Saved Puzzles", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Manage and select your test case files."),
		widget.NewSeparator(),
		container.NewPadded(importBtn),
	)

	// ==========================================
	// TAB 2: SOLVER PAGE
	// ==========================================
	solverContent := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("SOLVER DASHBOARD", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			statusLabel,
			statsLabel,
			widget.NewSeparator(),
		),
		container.NewPadded(solveBtn),
		nil,
		nil,
		canvasContainer, // Grid fills the remaining space
	)

	// ==========================================
	// TAB 3: STATS PAGE
	// ==========================================
	statsContent := container.NewBorder(
		widget.NewLabelWithStyle("Algorithm Analysis & Logs", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		container.NewPadded(logScroll),
	)

	// ==========================================
	// ASSEMBLE APP TABS
	// ==========================================
	tabs = container.NewAppTabs(
		container.NewTabItemWithIcon("Library", theme.StorageIcon(), container.NewPadded(libraryContent)),
		container.NewTabItemWithIcon("Solver", theme.ComputerIcon(), container.NewPadded(solverContent)),
		container.NewTabItemWithIcon("Stats", theme.InfoIcon(), container.NewPadded(statsContent)),
	)
	tabs.SetTabLocation(container.TabLocationBottom) // Bottom Navigation Bar

	// Wrap in Stack to maintain background color
	content := container.NewStack(bgRect, tabs)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
