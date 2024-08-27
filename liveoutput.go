package main

import (
	"bufio"
	"fmt"
	"jackob101/run/styles"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
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

type LiveoutputKeybinds struct {
	Close key.Binding
	Stop  key.Binding
}

var keybinds = LiveoutputKeybinds{
	Close: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Close liveoutput"),
	),
	Stop: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Stops the command"),
	),
}

func (k LiveoutputKeybinds) ShortHelp() []key.Binding {
	return []key.Binding{k.Close, k.Stop}
}

func (k LiveoutputKeybinds) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Close, k.Stop},
		{},
	}
}

type liveoutput struct {
	sub                chan string
	subStop            chan bool
	commandDisplayName string
	command            string
	lines              string
	finished           bool
	viewport           viewport.Model
	helpMenu           help.Model
	runtime            stopwatch.Model
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
			select {
			case value := <-m.subStop:
				if value {
					return loFinished{}
				}
			default:
			}
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
	helpMenu := help.New()
	helpMenu.Styles = styles.MenuStyles
	return liveoutput{
		sub:                make(chan string),
		subStop:            make(chan bool, 1),
		commandDisplayName: cmd.displayName,
		command:            cmd.cmd,
		lines:              "",
		finished:           false,
		viewport:           viewport.Model{},
		helpMenu:           helpMenu,
		runtime:            stopwatch.New(),
	}
}

func (m liveoutput) Init() tea.Cmd {
	return tea.Batch(m.listenForNewline(),
		m.waitForNewline(),
		wrapMsg(LoadViewport{}),
		m.runtime.Init(),
	)
}

func (m liveoutput) Update(msg tea.Msg) (liveoutput, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybinds.Stop):
			m.finished = true
			select {
			case m.subStop <- true:
			default:
			}
			cmds = append(cmds, wrapMsg(loFinished{}))
			cmds = append(cmds, m.runtime.Stop())
		case key.Matches(msg, keybinds.Close):
			m.finished = true
			return m, func() tea.Msg { return loClosed{} }
		}
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
		cmds = append(cmds, m.runtime.Stop())
	case LoadViewport:
		m.initViewport()
	case tea.WindowSizeMsg:
		m.initViewport()
		cmds = append(cmds, viewport.Sync(m.viewport))
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.runtime, cmd = m.runtime.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m liveoutput) View() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s",
		m.headerView(),
		m.viewport.View(),
		m.footerView(),
		styles.MenuStyle().Render(m.helpMenu.View(keybinds)),
	)
}

func (lo *liveoutput) initViewport() {
	helpMenuHeight := lipgloss.Height(styles.MenuStyle().Render(lo.helpMenu.View(keybinds)))
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
	message += " - " + m.runtime.View()
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
