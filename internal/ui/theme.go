package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type darkTheme struct {
	isActive func() bool
}

var _ fyne.Theme = (*darkTheme)(nil)

func (m darkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.Black
	case theme.ColorNameButton:
		return color.Black
	case theme.ColorNameForeground:
		if m.isActive() {
			return color.RGBA{R: 0x40, G: 0x14, B: 0x54, A: 0xFF}
		}
		return color.RGBA{R: 0x40, G: 0x54, B: 0x14, A: 0xFF}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m darkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m darkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m darkTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 200.0
	}
	return theme.DefaultTheme().Size(name)
}
