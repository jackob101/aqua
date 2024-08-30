package common

import (
	"fmt"
	"log/slog"

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

var confirmationKeybinds = []Keybind{
	NewKeybind(ConfirmationDialogLeft{}, "Left", "left", "h"),
	NewKeybind(ConfirmationDialogRight{}, "Right", "right", "l"),
	NewKeybind(ConfirmationDialogSelect{}, "Select", "enter"),
}

// TODO: These keybinds are kinda annoying. Should decide if keybinds should be always shown at
// the bottom of the screen like in zellij.
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
	slog.Info("Init")
	return SetKeybindsCmd(confirmationKeybinds)
}

func (m Confirmation) Update(msg tea.Msg) (Confirmation, tea.Cmd) {
	switch msg := msg.(type) {
	case ConfirmationDialogLeft:
		m.selected = true
	case ConfirmationDialogRight:
		m.selected = false
	case ConfirmationDialogSelect:
		return m, MakeCmd(ConfirmationDialogSelected{Value: m.selected})
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
