package common

import (
	"strings"

	help "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type MenuType string

const (
	Short MenuType = "string"
	Long  MenuType = "long"
	Both  MenuType = "both"
)

type Keybinds[T any] struct {
	Keybinds []Keybind[T]
	helpMenu help.Model
}

type Keybind[T any] struct {
	Callback    func(m Keybinds[T]) (Keybinds[T], tea.Cmd)
	Description string
	Menu        MenuType
	Keys        []string
}

func NewHelpMenu[T any](keybinds []Keybind[T], width int) Keybinds[T] {
	helpMenu := help.New()
	helpMenu.Width = width
	return Keybinds[T]{
		Keybinds: keybinds,
		helpMenu: helpMenu,
	}
}

func (m Keybinds[T]) Update(msg tea.KeyMsg) (Keybinds[T], tea.Cmd) {
	for _, e := range m.Keybinds {
		for _, keyE := range e.Keys {
			if keyE == msg.String() {
				return e.Callback(m)
			}
		}
	}
	return m, nil
}

func (m Keybinds[T]) Init() tea.Cmd {
	return nil
}

func (k Keybinds[T]) ShortHelp() []key.Binding {
	fullHelp := k.FullHelp()
	if len(fullHelp) != 0 {
		return fullHelp[0]
	}
	return []key.Binding{}
}

func (k Keybinds[T]) FullHelp() [][]key.Binding {
	keys := [][]key.Binding{}
	oneRow := []key.Binding{}
	for _, e := range k.Keybinds {
		oneRow = append(oneRow, key.NewBinding(
			key.WithKeys(e.Keys...),
			key.WithHelp(strings.Join(e.Keys, "/"), e.Description),
		))
	}
	keys = append(keys, oneRow)
	return keys
}

func (m Keybinds[T]) View() string {
	helpView := m.helpMenu.View(m)
	return helpView
}
