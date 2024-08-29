package common

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	boxStyle = func(message string) lipgloss.Style {
		border := lipgloss.RoundedBorder()
		return lipgloss.NewStyle().
			BorderStyle(border).
			Width(len(message)+10).
			Align(lipgloss.Center, lipgloss.Center).
			Padding(1)
	}

	buttonStyle = func() lipgloss.Style {
		border := lipgloss.RoundedBorder()
		return lipgloss.NewStyle().
			BorderStyle(border).
			Align(lipgloss.Center).
			Width(10).
			Padding(0)
	}()

	selectedButtonStyle = func() lipgloss.Style {
		color := lipgloss.Color("#a9b665")
		border := lipgloss.RoundedBorder()
		return lipgloss.NewStyle().
			BorderStyle(border).
			BorderForeground(color).
			Align(lipgloss.Center).
			Width(10).
			Foreground(color).
			Padding(0)
	}()
)

type ConfirmationKeybinds struct {
	Left   key.Binding
	Right  key.Binding
	Accept key.Binding
}

var keybinds = ConfirmationKeybinds{
	Left:   key.NewBinding(key.WithKeys("left", "h")),
	Right:  key.NewBinding(key.WithKeys("right", "l")),
	Accept: key.NewBinding(key.WithKeys("enter")),
}

type Confirmation struct {
	message  string
	width    int
	height   int
	selected bool
}

func NewConfirmation(message string, width int, height int) Confirmation {
	return Confirmation{
		message:  message,
		width:    width,
		height:   height,
		selected: true,
	}
}

func (m Confirmation) Init() tea.Cmd {
	return nil
}

func (m Confirmation) Update(msg tea.Msg) (Confirmation, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybinds.Right):
			m.selected = false
		case key.Matches(msg, keybinds.Left):
			m.selected = true
		case key.Matches(msg, keybinds.Accept):
			return m, func() tea.Msg {
				return Confirmation_Selected{
					Selected: m.selected,
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m Confirmation) View() string {
	var yesBtn string
	var noBtn string
	if m.selected {
		yesBtn = selectedButtonStyle.Render("Yes")
		noBtn = buttonStyle.Render("No")
	} else {
		yesBtn = buttonStyle.Render("Yes")
		noBtn = selectedButtonStyle.Render("No")
	}
	confirmationBox := boxStyle(m.message).Render(fmt.Sprintf("%s\n\n%s", m.message, lipgloss.JoinHorizontal(lipgloss.Center, yesBtn, noBtn)))
	return lipgloss.Place(m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		confirmationBox)
}
