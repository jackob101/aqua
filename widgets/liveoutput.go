package widgets

import (
	"bufio"
	"jackob101/run/common"
	"jackob101/run/styles"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	gutterStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Right).
			MarginRight(1).
			Faint(true)
	detailsValueStyle = lipgloss.NewStyle().
				Faint(true)
	detailsKeyNameStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Bold(true)
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
		common.NewKeybind(common.LiveoutputToggleDetails{}, "toggle details", "d"),
		common.NewKeybind(common.LiveoutputUp{}, "scroll up", "up", "k"),
		common.NewKeybind(common.LiveoutputDown{}, "scroll down", "down", "j"),
		common.NewKeybind(common.LiveoutputOpenEditor{}, "open editor", "e"),
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
	lines                []string
	lineCount            int
	finished             bool
	viewport             viewport.Model
	viewportOffset       int
	viewportHeight       int
	helpMenu             help.Model
	runtime              stopwatch.Model
	closeConfirmation    *common.Confirmation
	showDetails          bool
	width                int
	height               int
}

type newline struct {
	content string
}

func (m liveoutput) listenForNewline() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("bash", "-c", m.command)
		stdout, _ := cmd.StdoutPipe()
		time.Sleep(500 * time.Millisecond)
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
					if cmd.Cancel != nil {
						if err := cmd.Cancel(); err != nil {
							slog.Error("I'am too stupid to handle this error", "value", err.Error())
						}
					}
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

	return m.runtime.Stop()
}

func (m *liveoutput) refreshKeybinds() tea.Cmd {
	return common.SetKeybindsCmd(m.getLiveoutputKeybinds())
}

func (m *liveoutput) restartCommand() tea.Cmd {
	if m.finished {
		m.lines = []string{}
		m.commandOutputChannel = make(chan string)
		m.finished = false
		m.viewport.SetContent(strings.Join(m.lines, "\n"))
		return tea.Batch(m.listenForNewline(),
			m.waitForNewline(),
			m.runtime.Reset(),
			m.runtime.Start(),
			common.SetKeybindsCmd(m.getLiveoutputKeybinds()),
		)
	}
	return nil
}

func NewLiveoutput(cmd string, displayName string, width int, height int) liveoutput {
	helpMenu := help.New()
	helpMenu.Styles = styles.MenuStyles
	lo := liveoutput{
		commandOutputChannel: make(chan string),
		commandStopChannel:   make(chan bool, 1),
		commandDisplayName:   displayName,
		command:              cmd,
		lines:                []string{},
		lineCount:            0,
		finished:             false,
		viewport:             viewport.Model{},
		helpMenu:             helpMenu,
		runtime:              stopwatch.New(),
		showDetails:          true,
		width:                width,
		height:               height,
	}
	lo.viewportHeight = lo.height - lo.getDetailsViewHeight()
	return lo
}

func (m liveoutput) openEditor() tea.Cmd {
	lines := strings.Join(m.lines, "\n")
	os.WriteFile("/tmp/aqua_tmp.txt", []byte(lines), 0666)
	cmd := exec.Command("nvim", "/tmp/aqua_tmp.txt")
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return common.LiveoutputEditorClosed{}
	})
}

func (m liveoutput) Init() tea.Cmd {
	return tea.Batch(m.listenForNewline(),
		m.waitForNewline(),
		m.runtime.Init(),
		common.MakeCmd(common.SetKeybinds{Keybinds: m.getLiveoutputKeybinds()}),
	)
}

func (m liveoutput) Update(msg tea.Msg) (liveoutput, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case newline:
		msg.content = strings.ReplaceAll(msg.content, "\\r\\n", "")
		if len(m.lines) >= m.height {
			m.viewportOffset++
		}
		m.lines = append(m.lines, msg.content)
		// m.viewport.SetContent(strings.Join(m.lines, "\n"))
		// m.viewport.GotoBottom()
		cmds = append(cmds, m.waitForNewline())
	case common.LiveoutputCommandFinished:
		m.finished = true
		cmds = append(cmds, m.runtime.Stop())
		if m.closeConfirmation == nil {
			cmds = append(cmds, common.SetKeybindsCmd(m.getLiveoutputKeybinds()))
		}
	case common.ContentSectionResize:
		m.width = msg.Width
		m.height = msg.Height
		cmds = append(cmds, viewport.Sync(m.viewport))
	case common.LiveoutputClose:
		confirmationDialogHeight := lipgloss.Height(m.viewport.View())
		confirmationDialog := common.NewConfirmation("Do You want to close liveouput?",
			m.width,
			confirmationDialogHeight)
		m.closeConfirmation = &confirmationDialog
		cmds = append(cmds, confirmationDialog.Init())
	case common.LiveoutputCommandRestart:
		cmds = append(cmds, m.restartCommand())
	case common.LiveoutputCommandStop:
		cmds = append(cmds, m.stopCommand())
		cmds = append(cmds, m.refreshKeybinds())
	case common.LiveoutputToggleDetails:
		m.showDetails = !m.showDetails
		detailsHeight := m.getDetailsViewHeight()
		m.viewport.Height = m.height - detailsHeight
		m.viewportHeight = m.height - detailsHeight
	case common.LiveoutputUp:
		m.viewportOffset = max(m.viewportOffset-1, 0)
	case common.LiveoutputDown:
		m.viewportOffset = min(m.viewportOffset+1, len(m.lines)-m.viewportHeight)
		if m.viewportOffset < 0 {
			m.viewportOffset = 0
		}
	case common.LiveoutputOpenEditor:
		cmds = append(cmds, m.openEditor())
		cmds = append(cmds, common.MakeCmd(common.LiveoutputEditorOpened{}))
	case common.ConfirmationDialogSelected:
		if m.closeConfirmation != nil {
			m.closeConfirmation = nil
			if msg.Value {
				cmds = append(cmds, m.stopCommand())
				cmds = append(cmds, common.MakeCmd(common.LiveoutputClosed{}))
			} else {
				m.closeConfirmation = nil
				cmds = append(cmds, m.refreshKeybinds())
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
		mainBody = lipgloss.NewStyle().
			Width(m.width).
			Height(m.viewportHeight).
			AlignVertical(lipgloss.Center).
			Render(m.closeConfirmation.View())
	} else {
		mainBody = lipgloss.NewStyle().
			Width(m.width).
			Height(m.viewportHeight).
			Render(m.viewportView())
	}

	if m.showDetails {
		return lipgloss.JoinVertical(
			0,
			m.detailsView(),
			mainBody,
		)
	} else {
		return mainBody
	}
}

func (m liveoutput) getDetailsViewHeight() int {
	if m.showDetails {
		return lipgloss.Height(m.detailsView())
	} else {
		return 0
	}
}

func (m liveoutput) getStatus() string {
	if m.finished {
		return "Finished"
	} else {
		return "Running"
	}
}

func (m liveoutput) viewportView() string {
	gutterNumberCount := common.IntDigits(len(m.lines))
	gutterSize := gutterNumberCount + gutterStyle.GetMarginRight()
	lines := []string{}
	appendedLines := 0
	for _, e := range m.lines[m.viewportOffset:] {
		if appendedLines >= m.viewportHeight {
			break
		}
		lineWidth := lipgloss.Width(e)
		chunkSize := m.width - gutterSize
		if lineWidth > chunkSize {
			for i := 0; i < lipgloss.Width(e); i += chunkSize {
				if appendedLines >= m.viewportHeight {
					break
				}
				part := e[i:min(lineWidth, i+chunkSize)]
				lines = append(lines, part)
				appendedLines++
			}
		} else {
			lines = append(lines, e)
			appendedLines++
		}
	}

	result := ""
	for i, e := range lines {
		gutter := gutterStyle.Width(gutterNumberCount).Render(strconv.Itoa(i + m.viewportOffset))
		result += gutter + e

		if i+1 != m.viewportHeight {
			result += "\n"
		}
	}

	return result
}

func (m liveoutput) detailsView() string {
	cmdTitleKey := "Command name: "
	cmdTitleValue := m.commandDisplayName
	cmdKey := "Command: "
	cmdValue := m.command
	statusKey := "Status: "
	statusValue := m.getStatus()
	timeKey := "Time: "
	timeValue := m.runtime.View()

	details := []string{
		lipgloss.JoinHorizontal(0, detailsKeyNameStyle.Render(cmdTitleKey), detailsValueStyle.Render(cmdTitleValue)),
		lipgloss.JoinHorizontal(0, detailsKeyNameStyle.Render(cmdKey), detailsValueStyle.Render(cmdValue)),
		lipgloss.JoinHorizontal(0, detailsKeyNameStyle.Render(statusKey), detailsValueStyle.Render(statusValue)),
		lipgloss.JoinHorizontal(0, detailsKeyNameStyle.Render(timeKey), detailsValueStyle.Render(timeValue)),
	}

	detailView := lipgloss.JoinVertical(0, details...)

	separator := strings.Repeat("─", m.width)
	return lipgloss.JoinVertical(0, detailView, separator)
}
