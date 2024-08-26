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
	return c.description
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
		0,
		0,
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			return m, wrapMsg(SelectedCommandEntry{command: m.list.SelectedItem().(Command)})
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m CommandList) View() string {
	return m.list.View()
}
