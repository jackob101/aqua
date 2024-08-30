package widgets

import (
	"fmt"
	"jackob101/run/common"
	"jackob101/run/common/dto"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	commandListEntryHeight     = 3
	commandListVerticalPadding = 2
)

var commandListKeybinds = []common.Keybind{
	common.NewKeybind(common.CommandListDown{}, "down", "down", "j"),
	common.NewKeybind(common.CommandListUp{}, "up", "up", "k"),
	common.NewKeybind(common.CommandListSelect{}, "select", "enter"),
	common.NewKeybind(nil, "filter", "/"),
	common.NewKeybind(common.CommandListQuit{}, "quit", "q"),
}

var (
	listStyle = func(height int, width int) lipgloss.Style {
		return lipgloss.NewStyle().
			Padding(1, commandListVerticalPadding).
			Height(height).
			Width(width)
	}
	selectedArrow = func() lipgloss.Border {
		b := lipgloss.NormalBorder()
		b.Right = "┃"
		return b
	}()
	entryDescriptionStyle = lipgloss.NewStyle().
				Faint(true)
	entryContainerStyle = lipgloss.NewStyle().
				PaddingBottom(1)
	selectedEntryStyle = lipgloss.NewStyle().
				Border(selectedArrow, false, false, false, true).
				Foreground(lipgloss.Color("#a9b665"))
	entryStyle = lipgloss.NewStyle().
			Padding(0, 0, 0, 1)
)

type CommandList struct {
	cmds     []dto.Command
	selected int
	top      int
	width    int
	height   int
}

func NewCommandList(width int, height int) CommandList {
	items := []dto.Command{
		{
			Cmd:   "echo Test from run",
			Title: "This is test command",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 2",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 3",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 4",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 5",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 6",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 7",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 8",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 9",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 10",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 11",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 12",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 13",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 14",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 15",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 16",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 17",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 18",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 19",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 20",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 21",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 22",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 23",
		},
		{
			Cmd:   "./test.sh",
			Title: "This is test command 24",
		},
	}

	return CommandList{
		cmds:   items,
		width:  width,
		height: height,
	}
}

func (m CommandList) determinePageSize() int {
	return (m.height + 1 - (commandListVerticalPadding * 2)) / commandListEntryHeight
}

func (m *CommandList) next() {
	if len(m.cmds) > m.selected+1 {
		m.selected++
	}
	if m.top+(m.determinePageSize()-1) < m.selected {
		m.top++
	}
}

func (m *CommandList) prev() {
	if m.selected > 0 {
		m.selected--
	}
	if m.selected < m.top {
		m.top = m.selected
	}
}

func (m CommandList) getSelected() tea.Msg {
	return common.CommandListSelected{
		Cmd: m.cmds[m.selected],
	}
}

func (m CommandList) Init() tea.Cmd {
	return common.MakeCmd(common.SetKeybinds{
		Keybinds: commandListKeybinds,
	})
}

func (m CommandList) Update(msg tea.Msg) (CommandList, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case common.CommandListUp:
		m.prev()
	case common.CommandListDown:
		m.next()
	case common.CommandListQuit:
		cmds = append(cmds, tea.Quit)
	case common.CommandListSelect:
		return m, m.getSelected
	case common.ContentSectionResize:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, tea.Batch(cmds...)
}

func (m CommandList) View() string {
	entries := []string{}

	pageSize := min(len(m.cmds)-m.top, m.determinePageSize())

	for i := 0; i < pageSize; i++ {
		e := m.cmds[m.top+i]
		var cmdView string
		if m.top+i == m.selected {
			cmdView = selectedEntryStyle.Render(fmt.Sprintf("%s\n\t%s", e.Title, entryDescriptionStyle.Render(e.Cmd)))
		} else {
			cmdView = entryStyle.Render(fmt.Sprintf("%s\n\t%s", e.Title, entryDescriptionStyle.Render(e.Cmd)))
		}

		entries = append(entries, cmdView)
	}

	listView := ""

	for _, e := range entries {
		listView += entryContainerStyle.Render(e) + "\n"
	}

	listView += fmt.Sprintf("%d/%d", m.selected+1, len(m.cmds))

	return listStyle(m.height, m.width).Render(listView)
}
