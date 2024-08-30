package widgets

import (
	"jackob101/run/common"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	keybindsDisplaySeparatorStyle = func() lipgloss.Border {
		b := lipgloss.NormalBorder()
		b.Top = "̅─"
		return b
	}
	keybindsDisplayBoxStyle = lipgloss.NewStyle().
				Padding(0, 1, 1, 1).
				Border(keybindsDisplaySeparatorStyle(), true, false, false, false)

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
	return keybindsDisplayBoxStyle.Width(m.width).Render(keybindsDisplayView)
}
