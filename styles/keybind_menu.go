package styles

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

var MenuStyle = func() lipgloss.Style {
	b := lipgloss.NewStyle().
		UnsetForeground().
		Faint(true).
		Padding(1)
	return b
}

var MenuStyles = help.Styles{
	Ellipsis:       lipgloss.Style{},
	ShortKey:       defaultFontStyle,
	ShortDesc:      defaultFontStyle,
	ShortSeparator: lipgloss.Style{},
	FullKey:        defaultFontStyle,
	FullDesc:       defaultFontStyle,
	FullSeparator:  lipgloss.Style{},
}

var defaultFontStyle = lipgloss.NewStyle().
	UnsetForeground().
	Faint(true)
