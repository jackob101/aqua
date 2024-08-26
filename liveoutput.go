package main

import (
	"bufio"
	"fmt"
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
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type liveoutput struct {
	sub                chan string
	commandDisplayName string
	command            string
	lines              string
	finished           bool
	viewport           viewport.Model
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
	}
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
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m liveoutput) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (lo *liveoutput) initViewport() {
	headerHeight := lipgloss.Height(lo.headerView())
	footerHeight := lipgloss.Height(lo.footerView())
	verticalMarginHeight := headerHeight + footerHeight
	lo.viewport = viewport.New(Width, Height-verticalMarginHeight)
	lo.viewport.YPosition = headerHeight
	lo.viewport.HighPerformanceRendering = false
	lo.viewport.SetContent(lo.lines)

	lo.viewport.YPosition = headerHeight + 1
}

func (m liveoutput) headerView() string {
	title := titleStyle.Render("Mr. Pager")
	line := strings.Repeat("─", max(0, Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m liveoutput) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", 0.0*100))
	line := strings.Repeat("─", max(0, Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
