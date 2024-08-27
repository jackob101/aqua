package main

import (
	"bufio"
	"fmt"
	"jackob101/run/common"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return titleStyle.Padding(0, 1).BorderStyle(b)
	}()
)

type liveoutput struct {
	sub                chan string
	commandDisplayName string
	command            string
	lines              string
	finished           bool
	viewport           viewport.Model
	helpMenu           common.Keybinds
}

type newline struct {
	content string
}

type loFinished struct{}

type loClosed struct{}

func (m liveoutput) listenForNewline() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("bash", "-c", m.command)
		stdout, _ := cmd.StdoutPipe()
		time.Sleep(1 * time.Second)
		err := cmd.Start()
		if err != nil {
			m.sub <- err.Error()
			return loFinished{}
		}
		buf := bufio.NewReader(stdout)
		for {
			line, _, err := buf.ReadLine()
			if err != nil {
				return loFinished{}
			}
			m.sub <- string(line)
		}
	}
}

func (m liveoutput) waitForNewline() tea.Cmd {
	return func() tea.Msg {
		content := <-m.sub
		return newline{content}
	}
}

func NewLiveoutput(cmd Command) liveoutput {
	return liveoutput{
		sub:                make(chan string),
		commandDisplayName: cmd.displayName,
		command:            cmd.cmd,
		lines:              "",
		finished:           false,
		viewport:           viewport.Model{},
		helpMenu: common.NewKeybindHandler([]common.Keybind{
			{
				Message:            Liveoutput_Quit{},
				Description:        "Close liveoutput",
				DisplayInShortMenu: true,
				Keys:               []string{"esc"},
			},
		}, Width),
	}
}

func (m liveoutput) GetKeybinds() common.Keybinds {
	return common.Keybinds{}
}

func (m liveoutput) Init() tea.Cmd {
	return tea.Batch(m.listenForNewline(),
		m.waitForNewline(),
		wrapMsg(LoadViewport{}),
	)
}

func (m liveoutput) Update(msg tea.Msg) (liveoutput, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		var cmd tea.Cmd
		m.helpMenu, cmd = m.helpMenu.Update(msg)
		cmds = append(cmds, cmd)
	case newline:
		m.lines += msg.content
		if !strings.HasSuffix("\n", msg.content) {
			m.lines += "\n"
		}
		m.viewport.SetContent(m.lines)
		m.viewport.GotoBottom()
		cmds = append(cmds, m.waitForNewline())
	case loFinished:
		m.finished = true
	case LoadViewport:
		m.initViewport()
	case tea.WindowSizeMsg:
		m.initViewport()
		cmds = append(cmds, viewport.Sync(m.viewport))
	case Liveoutput_Quit:
		m.finished = true
		return m, func() tea.Msg { return loClosed{} }
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m liveoutput) View() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s",
		m.headerView(),
		m.viewport.View(),
		m.footerView(),
		m.helpMenu.View(),
	)
}

func (lo *liveoutput) initViewport() {
	helpMenuHeight := lipgloss.Height(lo.helpMenu.View())
	headerHeight := lipgloss.Height(lo.headerView())
	footerHeight := lipgloss.Height(lo.footerView())
	verticalMarginHeight := headerHeight + footerHeight + helpMenuHeight
	lo.viewport = viewport.New(Width, Height-verticalMarginHeight)
	lo.viewport.YPosition = headerHeight
	lo.viewport.HighPerformanceRendering = false
	lo.viewport.SetContent(lo.lines)

	lo.viewport.YPosition = headerHeight + 1
}

func (m liveoutput) headerView() string {
	title := titleStyle.Render(m.commandDisplayName)
	line := strings.Repeat("─", max(0, Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m liveoutput) footerView() string {
	var message string
	if m.finished {
		message = "Finished"
	} else {
		message = "Running"
	}
	info := infoStyle.Render(message)
	line := strings.Repeat("─", max(0, Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, info, line)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
