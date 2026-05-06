package gui

import (
	"fmt"
	"iceSlidingPuzzle/src/filehandler"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// Map library
type LibraryPage struct {
	main        *MainUI
	searchQuery string
}

func NewLibraryPage(m *MainUI) *LibraryPage {
	return &LibraryPage{
		main:        m,
		searchQuery: m.librarySearchQuery,
	}
}

func (p *LibraryPage) Build() fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorSurfaceContainerLow)

	header := buildAppHeader(PageLibrary, p.main.NavigateTo)
	sidebar := p.buildSidebar()
	mainContent := p.buildMainContent()
	statusBar := buildStatusBar(
		[]string{fmt.Sprintf("Total: %d items", p.filteredLibraryCount())},
		[]string{},
	)

	// Main layout: sidebar + content
	body := container.NewHSplit(sidebar, mainContent)
	body.SetOffset(0.18)

	fullLayout := container.NewBorder(
		header,
		statusBar,
		nil, nil,
		body,
	)

	return container.NewStack(bg, fullLayout)
}

// Sidebar
func (p *LibraryPage) buildSidebar() fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorWhite)

	// Right border line
	rightBorder := canvas.NewRectangle(ColorOutlineVariant)
	rightBorder.SetMinSize(fyne.NewSize(1, 1))

	navTitle := sectionTitle("NAVIGATION")

	// Active item
	libBg := canvas.NewRectangle(ColorSidebarActive)
	libBg.CornerRadius = 4
	libIcon := canvas.NewText("📂", ColorPrimary)
	libIcon.TextSize = 16
	libLabel := canvas.NewText("Library", ColorOnSurface)
	libLabel.TextSize = 13
	libLabel.TextStyle = fyne.TextStyle{Bold: true}
	libItem := container.NewStack(
		libBg,
		container.NewPadded(container.NewHBox(libIcon, libLabel)),
	)

	// Solver
	solverIcon := canvas.NewText("🎯", ColorSlate600)
	solverIcon.TextSize = 16
	solverLabel := canvas.NewText("Solver", ColorSlate600)
	solverLabel.TextSize = 13
	solverBtn := widget.NewButton("", func() { p.main.NavigateTo(PageSolver) })
	solverBtn.Importance = widget.LowImportance
	solverItem := container.NewStack(
		solverBtn,
		container.NewPadded(container.NewHBox(solverIcon, solverLabel)),
	)

	navItems := container.NewVBox(
		navTitle,
		vSpacer(4),
		libItem,
		solverItem,
		vSpacer(2),
	)

	sidebarContent := container.NewBorder(
		container.NewPadded(navItems),
		nil, nil, nil, nil,
	)

	return container.NewStack(bg, rightBorder, sidebarContent)
}

// Main content
func (p *LibraryPage) buildMainContent() fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorWhite)

	toolbar := p.buildToolbar()
	grid := p.buildMapGrid()

	content := container.NewBorder(
		toolbar,
		nil, nil, nil,
		container.NewScroll(container.NewPadded(grid)),
	)

	return container.NewStack(bg, content)
}

func (p *LibraryPage) openImportDialog() {
	fd := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, p.main.window)
			return
		}
		if uc == nil {
			return
		}
		defer uc.Close()
		srcPath := uc.URI().Path()
		base := filepath.Base(srcPath)

		board, perr := filehandler.ParseBoard(srcPath)
		if perr != nil {
			dialog.ShowError(fmt.Errorf("Invalid map file: %w", perr), p.main.window)
			return
		}

		maxDim := board.N
		if board.M > maxDim {
			maxDim = board.M
		}
		difficulty := Easy
		if maxDim > 12 {
			difficulty = Hard
		} else if maxDim > 7 {
			difficulty = Intermediate
		}

		entry := &MapEntry{
			Filename:   base,
			FullPath:   srcPath,
			Width:      board.M,
			Height:     board.N,
			Difficulty: difficulty,
		}

		p.main.library = append(p.main.library, entry)
		p.main.SelectMap(entry)
	}, p.main.window)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
	fd.Show()
}

func (p *LibraryPage) buildToolbar() fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorWhite)
	bottomLine := canvas.NewRectangle(ColorDivider)
	bottomLine.SetMinSize(fyne.NewSize(1, 1))

	importBtn := widget.NewButton("⬆ Import", func() {
		p.openImportDialog()
	})
	importBtn.Importance = widget.HighImportance

	leftGroup := container.NewHBox(importBtn)

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("🔍 Search file...")
	searchEntry.SetText(p.searchQuery)
	searchEntry.MultiLine = false
	searchEntry.Wrapping = fyne.TextTruncate
	searchEntry.OnChanged = func(value string) {
		p.searchQuery = value
		p.main.librarySearchQuery = value
		p.main.showPage(PageLibrary)
	}

	sep := canvas.NewRectangle(ColorOutlineVariant)
	sep.SetMinSize(fyne.NewSize(1, 24))

	rightGroup := container.NewHBox(
		container.NewGridWrap(fyne.NewSize(300, 30), searchEntry),
		sep,
	)

	toolbarRow := container.NewBorder(nil, nil, leftGroup, rightGroup)

	content := container.NewVBox(
		container.NewPadded(toolbarRow),
		bottomLine,
	)

	return container.NewStack(bg, content)
}

func (p *LibraryPage) buildMapGrid() fyne.CanvasObject {
	cards := make([]fyne.CanvasObject, 0, len(p.main.library)+1)

	for _, entry := range p.main.library {
		if p.searchQuery != "" && !strings.Contains(strings.ToLower(entry.Filename), strings.ToLower(p.searchQuery)) {
			continue
		}
		e := entry
		cards = append(cards, p.buildMapCard(e))
	}

	// Add new map card
	cards = append(cards, p.buildAddCard())

	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(220, 150)), cards...)
	return grid
}

func (p *LibraryPage) filteredLibraryCount() int {
	if p.searchQuery == "" {
		return len(p.main.library)
	}

	count := 0
	query := strings.ToLower(p.searchQuery)
	for _, entry := range p.main.library {
		if strings.Contains(strings.ToLower(entry.Filename), query) {
			count++
		}
	}
	return count
}

// Card
func (p *LibraryPage) buildMapCard(entry *MapEntry) fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorWhite)
	bg.CornerRadius = 8
	bg.StrokeColor = ColorOutlineVariant
	bg.StrokeWidth = 1.5

	// File icon area
	iconBg := canvas.NewRectangle(ColorBlueBg)
	iconBg.CornerRadius = 6
	iconText := canvas.NewText("📄", ColorPrimary)
	iconText.TextSize = 20
	iconWidget := container.NewStack(
		container.NewGridWrap(fyne.NewSize(40, 40), iconBg),
		container.NewCenter(iconText),
	)

	// Difficulty badge
	var badge fyne.CanvasObject
	switch entry.Difficulty {
	case Easy:
		badge = easyBadge()
	case Intermediate:
		badge = intermediateBadge()
	case Hard:
		badge = hardBadge()
	}

	topRow := container.NewBorder(nil, nil, iconWidget, badge)

	// Filename
	filename := canvas.NewText(entry.Filename, ColorSlate800)
	filename.TextSize = 13
	filename.TextStyle = fyne.TextStyle{Bold: true}

	// Dimensions
	dims := canvas.NewText(fmt.Sprintf("%d × %d Dimensions", entry.Width, entry.Height), ColorSlate500)
	dims.TextSize = 11

	info := container.NewVBox(filename, dims)

	// Buttons
	selectBtn := primaryButton("Select", func() {
		p.main.SelectMap(entry)
	})

	btnRow := container.NewBorder(nil, nil, nil, selectBtn)

	cardContent := container.NewVBox(
		topRow,
		info,
		btnRow,
	)

	return container.NewStack(bg, container.NewPadded(cardContent))
}

// Add new map card
func (p *LibraryPage) buildAddCard() fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorWhite)
	bg.CornerRadius = 8
	bg.StrokeColor = ColorOutlineVariant
	bg.StrokeWidth = 2
	// dashed effect via slightly transparent fill
	bg.FillColor = ColorTransparent

	addBg := canvas.NewRectangle(ColorSlate100)
	addBg.CornerRadius = 6
	plusText := canvas.NewText("+", ColorSlate500)
	plusText.TextSize = 22
	plusText.TextStyle = fyne.TextStyle{Bold: true}
	addIcon := container.NewCenter(
		container.NewStack(
			container.NewGridWrap(fyne.NewSize(40, 40), addBg),
			container.NewCenter(plusText),
		),
	)

	label := canvas.NewText("New Map", ColorSlate600)
	label.TextSize = 12
	label.TextStyle = fyne.TextStyle{Bold: true}

	addBtn := widget.NewButton("", func() {
		p.openImportDialog()
	})
	addBtn.Importance = widget.LowImportance

	content := container.NewCenter(
		container.NewVBox(addIcon, label),
	)

	return container.NewStack(
		bg,
		addBtn,
		content,
	)
}
