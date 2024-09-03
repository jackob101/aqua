package widgets

import (
	"jackob101/run/common"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m ErrorScreen) getErrorScreenStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ea6962")).
		Align(lipgloss.Center, lipgloss.Center).
		Width(m.Width).
		Height(m.Height)
}

var errorScreenKeybinds = []common.Keybind{
	common.NewKeybind(tea.Quit(), "quit", "esc"),
}

type ErrorScreen struct {
	Err    error
	Width  int
	Height int
}

func (m ErrorScreen) Init() tea.Cmd {
	return common.SetKeybindsCmd(errorScreenKeybinds)
}

func (m ErrorScreen) Update(msg tea.Msg) (ErrorScreen, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case common.ContentSectionResize:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	return m, tea.Batch(cmds...)
}

func (m ErrorScreen) View() string {
	slog.Info("pane content", "value", m)
	return m.getErrorScreenStyle().Render(m.Err.Error())
}
