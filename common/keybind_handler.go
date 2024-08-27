package common

import (
	"strings"

	help "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var menuStyle = func() lipgloss.Style {
	b := lipgloss.NewStyle().
		UnsetForeground().
		Faint(true).
		Padding(1)
	return b
}

var menuStyles = help.Styles{
	Ellipsis:       lipgloss.Style{},
	ShortKey:       menuStyle(),
	ShortDesc:      menuStyle(),
	ShortSeparator: lipgloss.Style{},
	FullKey:        menuStyle(),
	FullDesc:       menuStyle(),
	FullSeparator:  lipgloss.Style{},
}

type Keybinds struct {
	Keybinds []Keybind
	helpMenu help.Model
}

type Keybind struct {
	Message            interface{}
	Description        string
	DisplayInShortMenu bool
	Keys               []string
}

func NewKeybindHandler(keybinds []Keybind, width int) Keybinds {
	helpMenu := help.New()
	helpMenu.Styles = menuStyles
	helpMenu.Width = width
	return Keybinds{
		Keybinds: keybinds,
		helpMenu: helpMenu,
	}
}

func (m Keybinds) Update(msg tea.KeyMsg) (Keybinds, tea.Cmd) {
	for _, e := range m.Keybinds {
		for _, keyE := range e.Keys {
			if keyE == msg.String() {
				return m, func() tea.Msg {
					return e.Message
				}
			}
		}
	}
	return m, nil
}

func (m Keybinds) Init() tea.Cmd {
	return nil
}

func (k Keybinds) ShortHelp() []key.Binding {
	fullHelp := k.FullHelp()
	if len(fullHelp) != 0 {
		return fullHelp[0]
	}
	return []key.Binding{}
}

func (k Keybinds) FullHelp() [][]key.Binding {
	keys := [][]key.Binding{}
	oneRow := []key.Binding{}
	// TODO: Handle DisplayInShortMenu
	for _, e := range k.Keybinds {
		oneRow = append(oneRow, key.NewBinding(
			key.WithKeys(e.Keys...),
			key.WithHelp(strings.Join(e.Keys, "/"), e.Description),
		))
	}
	keys = append(keys, oneRow)
	return keys
}

func (m Keybinds) View() string {
	helpView := m.helpMenu.View(m)
	return menuStyle().Render(helpView)
}
