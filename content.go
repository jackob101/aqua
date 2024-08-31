package main

import (
	"jackob101/run/common"
	"jackob101/run/widgets"

	tea "github.com/charmbracelet/bubbletea"
)

type MainView struct {
	commandList *widgets.CommandList
	lo          *liveoutput
}

func initialModel() MainView {
	list := widgets.NewCommandList(Width, Height)
	return MainView{
		commandList: &list,
	}
}

func (m MainView) Init() tea.Cmd {
	return m.commandList.Init()
}

func (m MainView) Update(msg tea.Msg) (MainView, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case common.SelectedCommandEntry:
		{
			m.commandList = nil
			lo := NewLiveoutput(msg.Cmd, msg.DisplayName)
			m.lo = &lo
			return m, m.lo.Init()
		}
	case common.CommandListSelected:
		m.commandList = nil
		lo := NewLiveoutput(msg.Cmd.Cmd, msg.Cmd.Title)
		m.lo = &lo
		return m, m.lo.Init()
	case common.LiveoutputClosed:
		m.lo = nil
		cmdList := widgets.NewCommandList(Width, Height)
		m.commandList = &cmdList
		cmds = append(cmds, cmdList.Init())
	}

	if m.commandList != nil {
		mResp, cmd := m.commandList.Update(msg)
		m.commandList = &mResp
		cmds = append(cmds, cmd)
	}

	if m.lo != nil {
		lo, cmd := m.lo.Update(msg)
		m.lo = &lo
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m MainView) View() string {
	if m.commandList != nil {
		return m.commandList.View()
	} else if m.lo != nil {
		return m.lo.View()
	}
	return "Missing output"
}
