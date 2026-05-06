package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func colorRect(c color.Color) *canvas.Rectangle {
	r := canvas.NewRectangle(c)
	return r
}

// App Header
func buildAppHeader(currentPage Page, onNavigate func(Page)) fyne.CanvasObject {
	bg := colorRect(ColorWhite)

	// Logo + Title
	snowflakeIcon := canvas.NewText("❄", ColorPrimary)
	snowflakeIcon.TextSize = 20
	snowflakeIcon.TextStyle = fyne.TextStyle{Bold: true}
	title := canvas.NewText("Arctic Solver", ColorOnSurface)
	title.TextSize = 15
	title.TextStyle = fyne.TextStyle{Bold: true}
	logoGroup := container.NewHBox(snowflakeIcon, widget.NewSeparator(), title)

	// Nav buttons
	makeNavBtn := func(label string, page Page) *widget.Button {
		btn := widget.NewButton(label, func() { onNavigate(page) })
		if page == currentPage {
			btn.Importance = widget.HighImportance
		} else {
			btn.Importance = widget.LowImportance
		}
		return btn
	}

	solverBtn := makeNavBtn("Solver", PageSolver)
	libraryBtn := makeNavBtn("Library", PageLibrary)

	divider := canvas.NewLine(ColorOutlineVariant)
	divider.StrokeWidth = 1

	navGroup := container.NewHBox(
		solverBtn,
		libraryBtn,
		widget.NewSeparator(),
	)

	header := container.NewBorder(nil, nil, logoGroup, navGroup)

	// Bottom border line
	bottomLine := canvas.NewLine(ColorDivider)
	bottomLine.StrokeWidth = 2

	content := container.NewVBox(
		container.NewPadded(header),
		bottomLine,
	)

	return container.NewStack(bg, content)
}

// Status bar
func buildStatusBar(leftItems []string, rightItems []string) fyne.CanvasObject {
	bg := colorRect(ColorWhite)
	bottomBorder := canvas.NewLine(ColorDivider)

	makeLabel := func(s string) *canvas.Text {
		t := canvas.NewText(s, ColorSlate500)
		t.TextSize = 10
		return t
	}

	leftObjs := []fyne.CanvasObject{}
	for i, s := range leftItems {
		leftObjs = append(leftObjs, makeLabel(s))
		if i < len(leftItems)-1 {
			sep := canvas.NewLine(ColorOutlineVariant)
			sep.StrokeWidth = 1
			leftObjs = append(leftObjs, sep)
		}
	}

	rightObjs := []fyne.CanvasObject{}
	for _, s := range rightItems {
		rightObjs = append(rightObjs, makeLabel(s))
	}

	if len(rightObjs) > 0 {
		dot := canvas.NewCircle(ColorGreen500)
		dot.Resize(fyne.NewSize(6, 6))
		connLabel := makeLabel("Connected")
		rightObjs = append(rightObjs, dot, connLabel)
	}

	leftRow := container.NewHBox(leftObjs...)
	rightRow := container.NewHBox(rightObjs...)
	row := container.NewBorder(nil, nil, leftRow, rightRow)

	statusContent := container.NewVBox(
		bottomBorder,
		container.NewPadded(row),
	)

	return container.NewStack(
		container.NewGridWrap(fyne.NewSize(1280, 48), bg),
		statusContent,
	)
}

type badgeStyle struct {
	bg     color.Color
	fg     color.Color
	border color.Color
}

func difficultyBadge(label string, style badgeStyle) fyne.CanvasObject {
	bg := canvas.NewRectangle(style.bg)
	bg.CornerRadius = 3
	bg.StrokeColor = style.border
	bg.StrokeWidth = 1

	text := canvas.NewText(label, style.fg)
	text.TextSize = 9
	text.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewStack(bg, container.NewPadded(text))
}

func easyBadge() fyne.CanvasObject {
	return difficultyBadge("EASY", badgeStyle{
		bg:     ColorGreenBg,
		fg:     ColorGreenText,
		border: ColorGreenBorder,
	})
}

func intermediateBadge() fyne.CanvasObject {
	return difficultyBadge("INTERMEDIATE", badgeStyle{
		bg:     ColorYellowBg,
		fg:     ColorYellowText,
		border: ColorYellowBorder,
	})
}

func hardBadge() fyne.CanvasObject {
	return difficultyBadge("HARD", badgeStyle{
		bg:     ColorRedBg,
		fg:     ColorRedText,
		border: ColorRedBorder,
	})
}

func primaryButton(label string, onTap func()) *widget.Button {
	btn := widget.NewButton(label, onTap)
	btn.Importance = widget.HighImportance
	return btn
}

func sectionCard(content fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorWhite)
	bg.CornerRadius = 8
	bg.StrokeColor = ColorOutlineVariant
	bg.StrokeWidth = 1
	return container.NewStack(bg, container.NewPadded(content))
}

func statBlock(label, value string) fyne.CanvasObject {
	lbl := canvas.NewText(label, ColorOutline)
	lbl.TextSize = 9
	lbl.TextStyle = fyne.TextStyle{Bold: true}

	val := canvas.NewText(value, ColorPrimary)
	val.TextSize = 18
	val.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewVBox(lbl, val)
}

func metricBlock(label, value string) fyne.CanvasObject {
	lbl := canvas.NewText(label, ColorSlate500)
	lbl.TextSize = 10
	lbl.TextStyle = fyne.TextStyle{Bold: true}

	val := canvas.NewText(value, ColorPrimary)
	val.TextSize = 13
	val.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewVBox(lbl, val)
}

func hLine(c color.Color) fyne.CanvasObject {
	line := canvas.NewRectangle(c)
	line.SetMinSize(fyne.NewSize(1, 1))
	return line
}

func vSpacer(h float32) fyne.CanvasObject {
	spacer := layout.NewSpacer()
	_ = spacer
	r := canvas.NewRectangle(ColorTransparent)
	r.SetMinSize(fyne.NewSize(1, h))
	return r
}

func hSpacer(w float32) fyne.CanvasObject {
	r := canvas.NewRectangle(ColorTransparent)
	r.SetMinSize(fyne.NewSize(w, 1))
	return r
}

func sectionTitle(text string) *canvas.Text {
	t := canvas.NewText(text, ColorSlate400)
	t.TextSize = 10
	t.TextStyle = fyne.TextStyle{Bold: true}
	return t
}
