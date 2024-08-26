package main

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	Width  int = 0
	Height int = 0
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
	slog.Debug("Message",
		"Type", reflect.TypeOf(msg),
		"Value", fmt.Sprintf("%+v", msg))
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case SelectedCommandEntry:
		{
			m.commandList = nil
			lo := NewLiveoutput(msg.command)
			m.lo = &lo
			return m, m.lo.Init()
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			{
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		Width = msg.Width
		Height = msg.Height
	case loClosed:
		m.lo = nil
		newCommandList := NewCommandList()
		// cmds = append(cmds, newCommandList.Init())
		m.commandList = &newCommandList
	}

	if m.commandList != nil {
		mResp, cmd := m.commandList.Update(msg)
		typedMResp := mResp.(CommandList)
		m.commandList = &typedMResp
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

	p := tea.NewProgram(initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Failed to run program %v", err)
		os.Exit(1)
	}
}
