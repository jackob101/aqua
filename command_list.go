package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type CommandList struct {
	cmds []Command
	list list.Model
}

type Command struct {
	cmd         string
	description string
	displayName string
}

func (c Command) Title() string {
	return c.displayName
}

func (c Command) Description() string {
	return c.cmd
}

func (c Command) FilterValue() string {
	return c.displayName
}

func NewCommandList() CommandList {
	items := []list.Item{
		Command{
			cmd:         "echo Test from run",
			description: "This is test command",
			displayName: "Test Command",
		}, Command{
			cmd:         "./test.sh",
			description: "This is test command 2",
			displayName: "Test Command 2",
		},
	}

	commandList := list.New(items,
		list.NewDefaultDelegate(),
		Width,
		Height,
	)
	commandList.Title = "Commands"

	return CommandList{
		cmds: []Command{},
		list: commandList,
	}
}

func (m CommandList) Init() tea.Cmd {
	return nil
}

func (m CommandList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		var cmd tea.Cmd
		m, cmd = m.handleKeybind(msg)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m CommandList) View() string {
	return m.list.View()
}

func (m CommandList) handleKeybind(msg tea.KeyMsg) (CommandList, tea.Cmd) {
	switch msg.String() {
	case "enter":
		selected := m.list.SelectedItem()
		if selected != nil {
			return m, wrapMsg(SelectedCommandEntry{command: selected.(Command)})
		}
	}
	return m, nil
}
