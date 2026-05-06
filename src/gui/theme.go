package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ArcticTheme struct{}

var _ fyne.Theme = (*ArcticTheme)(nil)

func (t *ArcticTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 246, G: 250, B: 255, A: 255} 
	case theme.ColorNameForeground:
		return color.NRGBA{R: 21, G: 29, B: 34, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0, G: 94, B: 151, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 0, G: 94, B: 151, A: 255}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 192, G: 199, B: 210, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 237, G: 244, B: 253, A: 255}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 112, G: 120, B: 130, A: 255}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 192, G: 199, B: 210, A: 255}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 20}
	case theme.ColorNameHeaderBackground:
		return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	case theme.ColorNameHover:
		return color.NRGBA{R: 210, G: 221, B: 232, A: 255}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 207, G: 229, B: 255, A: 255}
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 225, G: 233, B: 241, A: 255}
	case theme.ColorNameError:
		return color.NRGBA{R: 186, G: 26, B: 26, A: 255}
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 15, G: 105, B: 0, A: 255}
	case theme.ColorNameWarning:
		return color.NRGBA{R: 0, G: 119, B: 190, A: 255}
	default:
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	}
}

func (t *ArcticTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *ArcticTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *ArcticTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameScrollBar:
		return 10
	case theme.SizeNameScrollBarSmall:
		return 4
	case theme.SizeNameText:
		return 13
	case theme.SizeNameInputBorder:
		return 1.5
	case theme.SizeNameInnerPadding:
		return 8
	case theme.SizeNameLineSpacing:
		return 4
	}
	return theme.DefaultTheme().Size(name)
}

type tappableRow struct {
    widget.BaseWidget
    content  fyne.CanvasObject
    onTapped func()
}

// Hover
func newTappableRow(content fyne.CanvasObject, onTapped func()) *tappableRow {
    t := &tappableRow{content: content, onTapped: onTapped}
    t.ExtendBaseWidget(t)
    return t
}

func (t *tappableRow) CreateRenderer() fyne.WidgetRenderer {
    return widget.NewSimpleRenderer(t.content)
}

func (t *tappableRow) Tapped(_ *fyne.PointEvent) {
    if t.onTapped != nil {
        t.onTapped()
    }
}

func (t *tappableRow) TappedSecondary(_ *fyne.PointEvent) {}

// Color constants
var (
	ColorPrimary              = color.NRGBA{R: 0, G: 94, B: 151, A: 255}
	ColorPrimaryDark          = color.NRGBA{R: 0, G: 74, B: 121, A: 255}
	ColorPrimaryLight         = color.NRGBA{R: 207, G: 229, B: 255, A: 255}
	ColorPrimaryContainer     = color.NRGBA{R: 0, G: 119, B: 190, A: 255}
	ColorBackground           = color.NRGBA{R: 246, G: 250, B: 255, A: 255}
	ColorSurface              = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	ColorSurfaceContainer     = color.NRGBA{R: 231, G: 239, B: 247, A: 255}
	ColorSurfaceContainerLow  = color.NRGBA{R: 237, G: 244, B: 253, A: 255}
	ColorSurfaceContainerHigh = color.NRGBA{R: 225, G: 233, B: 241, A: 255}
	ColorOutline              = color.NRGBA{R: 112, G: 120, B: 130, A: 255}
	ColorOutlineVariant       = color.NRGBA{R: 192, G: 199, B: 210, A: 255}
	ColorOnSurface            = color.NRGBA{R: 21, G: 29, B: 34, A: 255}
	ColorOnSurfaceVariant     = color.NRGBA{R: 64, G: 71, B: 81, A: 255}
	ColorSlate400             = color.NRGBA{R: 148, G: 163, B: 184, A: 255}
	ColorSlate500             = color.NRGBA{R: 100, G: 116, B: 139, A: 255}
	ColorSlate600             = color.NRGBA{R: 71, G: 85, B: 105, A: 255}
	ColorSlate700             = color.NRGBA{R: 51, G: 65, B: 85, A: 255}
	ColorSlate800             = color.NRGBA{R: 30, G: 41, B: 59, A: 255}
	ColorGreen500             = color.NRGBA{R: 34, G: 197, B: 94, A: 255}
	ColorGreenBg              = color.NRGBA{R: 240, G: 253, B: 244, A: 255}
	ColorGreenText            = color.NRGBA{R: 21, G: 128, B: 61, A: 255}
	ColorGreenBorder          = color.NRGBA{R: 187, G: 247, B: 208, A: 255}
	ColorYellowBg             = color.NRGBA{R: 254, G: 252, B: 232, A: 255}
	ColorYellowText           = color.NRGBA{R: 161, G: 98,  B: 7,   A: 255}
	ColorYellowBorder         = color.NRGBA{R: 253, G: 230, B: 138, A: 255}
	ColorBlueBg               = color.NRGBA{R: 239, G: 246, B: 255, A: 255}
	ColorBlueText             = color.NRGBA{R: 29, G: 78, B: 216, A: 255}
	ColorBlueBorder           = color.NRGBA{R: 191, G: 219, B: 254, A: 255}
	ColorRedBg                = color.NRGBA{R: 254, G: 242, B: 242, A: 255}
	ColorRedText              = color.NRGBA{R: 185, G: 28, B: 28, A: 255}
	ColorRedBorder            = color.NRGBA{R: 254, G: 202, B: 202, A: 255}
	ColorTertiary             = color.NRGBA{R: 15, G: 105, B: 0, A: 255}
	ColorTertiaryBg           = color.NRGBA{R: 236, G: 255, B: 224, A: 255}
	ColorTertiaryBorder       = color.NRGBA{R: 187, G: 247, B: 160, A: 255}
	ColorIceTile              = color.NRGBA{R: 240, G: 244, B: 248, A: 255}
	ColorIceBorder            = color.NRGBA{R: 209, G: 217, B: 224, A: 255}
	ColorObstacle             = color.NRGBA{R: 219, G: 227, B: 235, A: 255}
	ColorObstacleBorder       = color.NRGBA{R: 176, G: 186, B: 197, A: 255}
	ColorPlayerOrb            = color.NRGBA{R: 0, G: 241, B: 253, A: 255}
	ColorPlayerBorder         = color.NRGBA{R: 0, G: 220, B: 230, A: 255}
	ColorGoalOrb              = color.NRGBA{R: 121, G: 255, B: 91, A: 255}
	ColorGoalBorder           = color.NRGBA{R: 42, G: 229, B: 0, A: 255}
	ColorTransparent          = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	ColorHeaderBg             = color.NRGBA{R: 248, G: 249, B: 250, A: 255}
	ColorDivider              = color.NRGBA{R: 225, G: 233, B: 241, A: 255}
	ColorCardHoverBorder      = color.NRGBA{R: 0, G: 119, B: 190, A: 255}
	ColorSidebarActive        = color.NRGBA{R: 207, G: 229, B: 255, A: 255}
	ColorWhite                = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	ColorSlate50              = color.NRGBA{R: 248, G: 250, B: 252, A: 255}
	ColorSlate100             = color.NRGBA{R: 241, G: 245, B: 249, A: 255}
	ColorSlate200             = color.NRGBA{R: 226, G: 232, B: 240, A: 255}
)
