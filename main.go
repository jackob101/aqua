package main

import (
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type MainView struct {
	commandList *CommandList
	lo          *liveoutput
}

func initialModel() MainView {
	list := NewCommandList()
	return MainView{
		commandList: &list,
	}
}

func (m MainView) Init() tea.Cmd {
	return nil
}

func (m MainView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	if m.commandList != nil {
		mResp, cmd := m.commandList.Update(msg)
		typedMResp := mResp.(CommandList)
		m.commandList = &typedMResp
		cmds = append(cmds, cmd)
	}

	if m.lo != nil {
		mResp, cmd := m.lo.Update(msg)
		typedMResp := mResp.(liveoutput)
		m.lo = &typedMResp
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case SelectedCommandEntry:
		{
			m.commandList = nil
			m.lo = &liveoutput{
				sub:                make(chan string),
				commandDisplayName: msg.command.displayName,
				command:            msg.command.cmd,
			}
			return m, m.lo.Init()
		}
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

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Failed to run program %v", err)
		os.Exit(1)
	}
}
