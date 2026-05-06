package main

import (
	"iceSlidingPuzzle/src/gui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.NewWithID("com.Arctic")
	a.Settings().SetTheme(&gui.ArcticTheme{})

	w := a.NewWindow("Arctic Solver")
	w.Resize(fyne.NewSize(1280, 720))
	w.SetFixedSize(false)

	mainUI := gui.NewMainUI(a, w)
	w.SetContent(mainUI.Build())

	w.ShowAndRun()
}
