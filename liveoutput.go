package main

import (
	"bufio"
	"fmt"
	"jackob101/run/common"
	"jackob101/run/styles"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.NormalBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.NormalBorder()
		b.Left = "┤"
		return titleStyle.Padding(0, 1).BorderStyle(b)
	}()
)

func (m liveoutput) getLiveoutputKeybinds() []common.Keybind {
	liveoutputKeybinds := []common.Keybind{
		common.NewKeybind(common.LiveoutputClose{}, "close", "esc"),
	}

	if m.finished {
		liveoutputKeybinds = append(liveoutputKeybinds,
			common.NewKeybind(common.LiveoutputCommandRestart{}, "restart command", "r"),
		)
	} else {
		liveoutputKeybinds = append(
			liveoutputKeybinds, common.NewKeybind(common.LiveoutputCommandStop{}, "stop command", "s"),
		)
	}
	return liveoutputKeybinds
}

type liveoutput struct {
	commandOutputChannel chan string
	commandStopChannel   chan bool
	commandDisplayName   string
	command              string
	lines                string
	finished             bool
	viewport             viewport.Model
	helpMenu             help.Model
	runtime              stopwatch.Model
	closeConfirmation    *common.Confirmation
}

type newline struct {
	content string
}

// TODO: This is not killing the command... But is should!
func (m liveoutput) listenForNewline() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("bash", "-c", m.command)
		stdout, _ := cmd.StdoutPipe()
		time.Sleep(1 * time.Second)
		err := cmd.Start()
		if err != nil {
			m.commandOutputChannel <- err.Error()
			return common.LiveoutputCommandFinished{}
		}
		buf := bufio.NewReader(stdout)
		for {
			select {
			case value := <-m.commandStopChannel:
				if value {
					return common.LiveoutputCommandFinished{}
				}
			default:
			}
			line, _, err := buf.ReadLine()
			if err != nil {
				return common.LiveoutputCommandFinished{}
			}
			m.commandOutputChannel <- string(line)
		}
	}
}

func (m liveoutput) waitForNewline() tea.Cmd {
	return func() tea.Msg {
		content := <-m.commandOutputChannel
		return newline{content}
	}
}

func (m *liveoutput) stopCommand() tea.Cmd {
	m.finished = true
	select {
	case m.commandStopChannel <- true:
	default:
	}

	return tea.Batch(m.runtime.Stop(),
		common.SetKeybindsCmd(m.getLiveoutputKeybinds()),
	)
}

func (m *liveoutput) restartCommand() tea.Cmd {
	if m.finished {
		m.lines = ""
		m.commandOutputChannel = make(chan string)
		m.finished = false
		m.viewport.SetContent(m.lines)
		return tea.Batch(m.listenForNewline(),
			m.waitForNewline(),
			m.runtime.Reset(),
			m.runtime.Start(),
			common.SetKeybindsCmd(m.getLiveoutputKeybinds()),
		)
	}
	return nil
}

func NewLiveoutput(cmd string, displayName string) liveoutput {
	helpMenu := help.New()
	helpMenu.Styles = styles.MenuStyles
	return liveoutput{
		commandOutputChannel: make(chan string),
		commandStopChannel:   make(chan bool, 1),
		commandDisplayName:   displayName,
		command:              cmd,
		lines:                "",
		finished:             false,
		viewport:             viewport.Model{},
		helpMenu:             helpMenu,
		runtime:              stopwatch.New(),
	}
}

func (m liveoutput) Init() tea.Cmd {
	return tea.Batch(m.listenForNewline(),
		m.waitForNewline(),
		wrapMsg(common.LoadViewport{}),
		m.runtime.Init(),
		common.MakeCmd(common.SetKeybinds{Keybinds: m.getLiveoutputKeybinds()}),
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
	case common.LiveoutputCommandFinished:
		m.finished = true
		cmds = append(cmds, m.runtime.Stop())
		cmds = append(cmds, common.SetKeybindsCmd(m.getLiveoutputKeybinds()))
	case common.LoadViewport:
		m.initViewport()
	case common.ContentSectionResize:
		m.initViewport()
		cmds = append(cmds, viewport.Sync(m.viewport))
	case common.LiveoutputClose:
		confirmationDialogHeight := lipgloss.Height(m.viewport.View())
		confirmationDialog := common.NewConfirmation("Do You want to close liveouput?",
			Width,
			confirmationDialogHeight)
		m.closeConfirmation = &confirmationDialog
		cmds = append(cmds, confirmationDialog.Init())
	case common.LiveoutputCommandRestart:
		cmds = append(cmds, m.restartCommand())
	case common.LiveoutputCommandStop:
		cmds = append(cmds, m.stopCommand())
	case common.ConfirmationDialogSelected:
		if m.closeConfirmation != nil {
			m.closeConfirmation = nil
			if msg.Value {
				cmds = append(cmds, m.stopCommand())
				cmds = append(cmds, common.MakeCmd(common.LiveoutputClosed{}))
			} else {
				m.closeConfirmation = nil
				cmds = append(cmds, common.SetKeybindsCmd(m.getLiveoutputKeybinds()))
			}
		}

	}

	if m.closeConfirmation != nil {
		newConfirmation, cmd := m.closeConfirmation.Update(msg)
		m.closeConfirmation = &newConfirmation
		cmds = append(cmds, cmd)
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.runtime, cmd = m.runtime.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m liveoutput) View() string {
	var mainBody string
	if m.closeConfirmation != nil {
		mainBody = m.closeConfirmation.View()
	} else {
		mainBody = m.viewport.View()
	}
	return fmt.Sprintf("%s\n%s",
		m.headerView(),
		mainBody,
		// styles.MenuStyle().Render(m.helpMenu.View(keybinds)),
	)
}

func (lo *liveoutput) initViewport() {
	// helpMenuHeight := lipgloss.Height(styles.MenuStyle().Render(lo.helpMenu.View(keybinds)))
	headerHeight := lipgloss.Height(lo.headerView())
	verticalMarginHeight := headerHeight
	lo.viewport = viewport.New(Width, Height-verticalMarginHeight)
	lo.viewport.YPosition = headerHeight
	lo.viewport.HighPerformanceRendering = false
	lo.viewport.SetContent(lo.lines)

	lo.viewport.YPosition = headerHeight + 1
}

func (m liveoutput) headerView() string {
	cmdDisplay := titleStyle.Render(m.commandDisplayName)

	var message string
	if m.finished {
		message = "Finished"
	} else {
		message = "Running"
	}
	message += " - " + m.runtime.View()
	info := infoStyle.Render(message)
	line := strings.Repeat("─", max(0, Width-lipgloss.Width(cmdDisplay)-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, cmdDisplay, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
