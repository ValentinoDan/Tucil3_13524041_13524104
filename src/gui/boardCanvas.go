package gui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"iceSlidingPuzzle/src/puzzle"
)

type boardLayout struct {
	width, height float32
}

func (b *boardLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	// Manual positioning, do nothing here
}

func (b *boardLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(b.width, b.height)
}

type BoardCanvas struct {
	Board      *puzzle.Board
	GridSize   float32
	Container  *fyne.Container
	PlayerRect *canvas.Rectangle
}

func NewBoardCanvas(b *puzzle.Board, gridSize float32) *BoardCanvas {
	bc := &BoardCanvas{
		Board:    b,
		GridSize: gridSize,
	}

	boardWidth := float32(b.M) * gridSize
	boardHeight := float32(b.N) * gridSize
	bc.Container = container.New(&boardLayout{width: boardWidth, height: boardHeight})
	
	bc.drawGrid()
	return bc
}

func (bc *BoardCanvas) drawGrid() {
	// Colors inspired by the mockup
	wallColor := color.NRGBA{R: 74, G: 85, B: 104, A: 255}     // Dark slate
	lavaColor := color.NRGBA{R: 229, G: 62, B: 62, A: 255}      // Red
	floorColor := color.NRGBA{R: 190, G: 227, B: 248, A: 120}   // Soft ice blue
	goalColor := color.NRGBA{R: 72, G: 187, B: 120, A: 255}     // Green
	playerColor := color.NRGBA{R: 49, G: 130, B: 206, A: 255}   // Solid Blue

	padding := float32(4) // padding between cells

	for row := 0; row < bc.Board.N; row++ {
		for col := 0; col < bc.Board.M; col++ {
			char := bc.Board.Grid[row][col]

			cellColor := floorColor
			if puzzle.IsWall(puzzle.Point{Row: row, Col: col}, bc.Board) {
				cellColor = wallColor
			} else if puzzle.IsLava(puzzle.Point{Row: row, Col: col}, bc.Board) {
				cellColor = lavaColor
			} else if char == 'O' {
				cellColor = goalColor
			}

			// Draw Background Rect with rounded corners
			rect := canvas.NewRectangle(cellColor)
			rect.CornerRadius = 8
			rect.Resize(fyne.NewSize(bc.GridSize-padding, bc.GridSize-padding))
			rect.Move(fyne.NewPos(float32(col)*bc.GridSize+padding/2, float32(row)*bc.GridSize+padding/2))
			bc.Container.Add(rect)

			// Texts
			if puzzle.IsWall(puzzle.Point{Row: row, Col: col}, bc.Board) {
				txt := canvas.NewText("X", color.White)
				txt.TextSize = bc.GridSize * 0.4
				txt.Alignment = fyne.TextAlignCenter
				txt.Resize(fyne.NewSize(bc.GridSize, bc.GridSize))
				txt.Move(fyne.NewPos(float32(col)*bc.GridSize, float32(row)*bc.GridSize+bc.GridSize*0.1))
				bc.Container.Add(txt)
			} else if puzzle.IsLava(puzzle.Point{Row: row, Col: col}, bc.Board) {
				txt := canvas.NewText("L", color.White)
				txt.TextSize = bc.GridSize * 0.4
				txt.Alignment = fyne.TextAlignCenter
				txt.Resize(fyne.NewSize(bc.GridSize, bc.GridSize))
				txt.Move(fyne.NewPos(float32(col)*bc.GridSize, float32(row)*bc.GridSize+bc.GridSize*0.1))
				bc.Container.Add(txt)
			} else if char == 'O' {
				txt := canvas.NewText("O", color.Black)
				txt.TextSize = bc.GridSize * 0.4
				txt.Alignment = fyne.TextAlignCenter
				txt.Resize(fyne.NewSize(bc.GridSize, bc.GridSize))
				txt.Move(fyne.NewPos(float32(col)*bc.GridSize, float32(row)*bc.GridSize+bc.GridSize*0.1))
				bc.Container.Add(txt)
			} else if char >= '0' && char <= '9' {
				// Checkpoint background
				bgRect := canvas.NewRectangle(color.White)
				bgRect.CornerRadius = 8
				bgRect.Resize(fyne.NewSize(bc.GridSize-padding, bc.GridSize-padding))
				bgRect.Move(fyne.NewPos(float32(col)*bc.GridSize+padding/2, float32(row)*bc.GridSize+padding/2))
				bc.Container.Add(bgRect)

				// Checkpoint text
				txt := canvas.NewText(string(char), playerColor)
				txt.TextSize = bc.GridSize * 0.4
				txt.TextStyle = fyne.TextStyle{Bold: true}
				txt.Alignment = fyne.TextAlignCenter
				txt.Resize(fyne.NewSize(bc.GridSize, bc.GridSize))
				txt.Move(fyne.NewPos(float32(col)*bc.GridSize, float32(row)*bc.GridSize+bc.GridSize*0.1))
				bc.Container.Add(txt)
			}
		}
	}

	// Draw Player Node
	bc.PlayerRect = canvas.NewRectangle(playerColor)
	bc.PlayerRect.CornerRadius = float32(bc.GridSize * 0.3)
	playerSize := bc.GridSize * 0.6
	bc.PlayerRect.Resize(fyne.NewSize(playerSize, playerSize))
	
	startCol, startRow := float32(bc.Board.Start.Col), float32(bc.Board.Start.Row)
	offset := (bc.GridSize - playerSize) / 2
	bc.PlayerRect.Move(fyne.NewPos(startCol*bc.GridSize+offset, startRow*bc.GridSize+offset))
	bc.Container.Add(bc.PlayerRect)
}

func (bc *BoardCanvas) AnimatePath(points []puzzle.Point) {
	go func() {
		for _, p := range points {
			playerSize := bc.GridSize * 0.6
			offset := (bc.GridSize - playerSize) / 2
			targetX := float32(p.Col)*bc.GridSize + offset
			targetY := float32(p.Row)*bc.GridSize + offset
			
			startPos := bc.PlayerRect.Position()
			steps := 15
			dx := (targetX - startPos.X) / float32(steps)
			dy := (targetY - startPos.Y) / float32(steps)

			for i := 0; i < steps; i++ {
				bc.PlayerRect.Move(fyne.NewPos(startPos.X+dx*float32(i+1), startPos.Y+dy*float32(i+1)))
				bc.PlayerRect.Refresh()
				time.Sleep(15 * time.Millisecond)
			}
			time.Sleep(150 * time.Millisecond) // slight pause after arriving
		}
	}()
}
