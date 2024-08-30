package main

import (
	"fmt"
	"jackob101/run/common"
	"jackob101/run/widgets"
	"log/slog"
	"os"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	Width  int = 0
	Height int = 0
)

type Root struct {
	content         MainView
	keybinds        []common.Keybind
	keybindsDisplay widgets.KeybindDisplay
}

func (m Root) handleKeybind(msg tea.KeyMsg) tea.Msg {
	for _, keybindE := range m.keybinds {
		for _, keyE := range keybindE.Keys {
			if keyE == msg.String() {
				if keybindE.Msg == nil {
					return msg
				}
				return keybindE.Msg
			}
		}
	}
	return msg
}

func (m Root) Init() tea.Cmd {
	return m.content.Init()
}

func (m Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("Message",
		"Type", reflect.TypeOf(msg),
		"Value", fmt.Sprintf("%+v", msg))
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case common.SetKeybinds:
		m.keybinds = msg.Keybinds
		m.keybindsDisplay, _ = m.keybindsDisplay.Update(msg)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		mappedMessage := m.handleKeybind(msg)
		switch mappedMessage.(type) {
		case tea.KeyMsg:
		default:
			return m, wrapMsg(mappedMessage)
		}
	case tea.WindowSizeMsg:
		Width = msg.Width
		Height = msg.Height - lipgloss.Height(m.keybindsDisplay.View())

		return m, wrapMsg(common.ContentSectionResize{
			Width:  Width,
			Height: Height,
		})
	}

	var cmd tea.Cmd
	m.content, cmd = m.content.Update(msg)
	cmds = append(cmds, cmd)

	m.keybindsDisplay, cmd = m.keybindsDisplay.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Root) View() string {
	keybindsDisplayView := m.keybindsDisplay.View()
	contentView := m.content.View()
	return lipgloss.JoinVertical(0, contentView, keybindsDisplayView)
}

func wrapMsg(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func main() {
	fo, _ := os.OpenFile("out.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer fo.Close()

	logginLevel := new(slog.LevelVar)
	slog.SetDefault(slog.New(
		slog.NewTextHandler(fo, &slog.HandlerOptions{Level: logginLevel})))
	logginLevel.Set(slog.LevelDebug)
	slog.Info("Logger configured")

	root := Root{
		content:         initialModel(),
		keybinds:        []common.Keybind{},
		keybindsDisplay: widgets.KeybindDisplay{},
	}

	p := tea.NewProgram(root,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Failed to run program %v", err)
		os.Exit(1)
	}
}
