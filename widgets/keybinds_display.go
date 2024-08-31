package widgets

import (
	"jackob101/run/common"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	keybindsDisplayBoxStyle = lipgloss.NewStyle().
				Padding(0, 1, 1, 1)

	keybindsDisplayDescriptionStyle = lipgloss.NewStyle().
					Faint(true)
)

type KeybindDisplay struct {
	keybinds []common.Keybind
	width    int
}

func (m KeybindDisplay) Init() tea.Cmd {
	return nil
}

func (m KeybindDisplay) Update(msg tea.Msg) (KeybindDisplay, tea.Cmd) {
	switch msg := msg.(type) {
	case common.ContentSectionResize:
		m.width = msg.Width
	case common.SetKeybinds:
		m.keybinds = msg.Keybinds
	}
	return m, nil
}

func (m KeybindDisplay) View() string {
	keybindsDisplayView := ""
	for _, e := range m.keybinds {
		keys := ""
		for _, keyEntry := range e.Keys {
			if len(keys) == 0 {
				keys += keyEntry
			} else {
				keys += "/" + keyEntry
			}
		}
		if len(keybindsDisplayView) == 0 {
			keybindsDisplayView += keys + " " + keybindsDisplayDescriptionStyle.Render(e.Description)
		} else {
			keybindsDisplayView += " · " + keys + " " + keybindsDisplayDescriptionStyle.Render(e.Description)
		}
	}

	return lipgloss.JoinVertical(0, strings.Repeat("─", m.width),
		keybindsDisplayBoxStyle.Width(m.width).Render(keybindsDisplayView),
	)
}
