package widgets

import (
	"jackob101/run/common"

	tea "github.com/charmbracelet/bubbletea"
)

type MainView struct {
	commandList *CommandList
	width       int
	height      int
	lo          *liveoutput
}

func NewContentPane(width int, height int) MainView {
	list := NewCommandList(width, height)
	return MainView{
		commandList: &list,
	}
}

func (m MainView) Init() tea.Cmd {
	return m.commandList.Init()
}

func (m MainView) Update(msg tea.Msg) (MainView, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msgi := msg.(type) {
	case common.SelectedCommandEntry:
		{
			m.commandList = nil
			lo := NewLiveoutput(msgi.Cmd, msgi.DisplayName, m.width, m.height)
			m.lo = &lo
			return m, m.lo.Init()
		}
	case common.CommandListSelected:
		m.commandList = nil
		lo := NewLiveoutput(msgi.Cmd.Cmd, msgi.Cmd.Title, m.width, m.height)
		m.lo = &lo
		return m, m.lo.Init()
	case common.LiveoutputClosed:
		m.lo = nil
		cmdList := NewCommandList(m.width, m.height)
		m.commandList = &cmdList
		cmds = append(cmds, cmdList.Init())
	case common.ContentSectionResize:
		m.width = msgi.Width
		m.height = msgi.Height
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
