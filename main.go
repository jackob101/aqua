package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type CommandList struct {
	cmds   []Command
	cursor int
	lo     *liveoutput
}

type Command struct {
	cmd         string
	description string
	displayName string
}

func initialModel() CommandList {
	return CommandList{
		lo: nil,
		cmds: []Command{{
			cmd:         "echo Test from run",
			description: "This is test command",
			displayName: "Test Command",
		}, {
			cmd:         "echo \"Test from run 2\"",
			description: "This is test command 2",
			displayName: "Test Command 2",
		}},
	}
}

func (m CommandList) Init() tea.Cmd {
	return nil
}

func (m CommandList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case loClosed:
		m.lo = nil
		return m, nil
	case tea.KeyMsg:
		// TODO: This should work like this. Make one component responsible for passing udpates
		// around. For example if LO is nil then display list. Add maybe some global keybinds on it
		// and the rest will be handled by child components
		if m.lo == nil {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.cmds)-1 {
					m.cursor++
				}

			case "enter", " ":
				{
					command := m.cmds[m.cursor]
					m.lo = &liveoutput{
						sub:                make(chan string),
						commandDisplayName: command.displayName,
						command:            command.cmd,
					}
					return m, m.lo.Init()
				}
			case "ctrl+l":
				return m, wrapMsg(tea.ClearScreen())
			}
		}
	}

	if m.lo != nil {
		mlo, cmd := m.lo.Update(msg)
		var test liveoutput = mlo.(liveoutput)
		m.lo = &test
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func wrapMsg(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func (m CommandList) View() string {
	if m.lo != nil {
		return m.lo.View()
	}
	s := "Please select command \n"

	for i, choice := range m.cmds {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice.displayName)
	}

	s += "\nPress q to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Failed to run program %v", err)
		os.Exit(1)
	}
}
