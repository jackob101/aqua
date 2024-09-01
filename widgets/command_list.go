package widgets

import (
	"fmt"
	"jackob101/run/common"
	"jackob101/run/common/dto"
	"log/slog"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	commandListEntryHeight      = 3
	commandListVerticalPadding  = 2
	commandListVerticalMargin   = 1
	commandListHorizontalMargin = 1
)

var (
	commandListStyle = lipgloss.NewStyle().
				Margin(commandListVerticalMargin, commandListHorizontalMargin)
	filterKeyStyle = lipgloss.NewStyle().
			Bold(true)
	listTitleStyle   = lipgloss.NewStyle().Inherit(filterKeyStyle)
	filterValueStyle = lipgloss.NewStyle().
				Faint(true)

	filterInputActiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#a9b665"))
	listStyle = func(height int, width int) lipgloss.Style {
		return lipgloss.NewStyle().
			Margin(1, commandListVerticalPadding).
			Height(height).
			Width(width)
	}
	selectedArrow = func() lipgloss.Border {
		b := lipgloss.NormalBorder()
		b.Right = "┃"
		return b
	}()
	entryDescriptionStyle = lipgloss.NewStyle().
				MarginLeft(4).
				Faint(true)
	entryContainerStyle = lipgloss.NewStyle().
				PaddingBottom(1)
	selectedEntryStyle = lipgloss.NewStyle().
				Border(selectedArrow, false, false, false, true).
				Foreground(lipgloss.Color("#a9b665"))
	entryStyle = lipgloss.NewStyle().
			Margin(0, 0, 0, 1)
)

type CommandList struct {
	cmds       []dto.Command
	filtered   []dto.Command
	selected   int
	top        int
	width      int
	height     int
	filtering  bool
	filteredBy string
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
			Title: "This is test command 9 filtering test",
		},
		{
			Cmd:   "This is test command 10 This is some very very long title that should't fit whole into the list so it should be cut",
			Title: "This is test command 10 This is some very very long title that should't fit whole into the list so it should be cut ",
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
		filtered: items,
		cmds:     items,
		width:    max(width-(commandListHorizontalMargin*2), 0),
		height:   max(height-(commandListVerticalMargin*2), 0),
	}
}

func (m CommandList) determinePageSize() int {
	return ((m.height + 1 - (commandListVerticalPadding * 2)) - lipgloss.Height(m.headerView())) / commandListEntryHeight
}

func (m *CommandList) next() {
	if len(m.filtered) > m.selected+1 {
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

func (m *CommandList) goToTop() {
	m.selected = 0
	m.top = 0
}

func (m CommandList) getKeybinds() []common.Keybind {
	if m.filtering {
		return []common.Keybind{
			// common.NewKeybind(common.CommandListSelect{}, "select", "enter"),
			common.NewKeybind(common.CommandListFilterToggle{}, "accept filter", "enter", "esc"),
		}
	} else {
		return []common.Keybind{
			common.NewKeybind(common.CommandListDown{}, "down", "down", "j"),
			common.NewKeybind(common.CommandListUp{}, "up", "up", "k"),
			common.NewKeybind(common.CommandListSelect{}, "select", "enter"),
			common.NewKeybind(common.CommandListFilterToggle{}, "filter", "/"),
			common.NewKeybind(common.CommandListQuit{}, "quit", "q"),
		}
	}
}

func (m CommandList) Init() tea.Cmd {
	return common.MakeCmd(common.SetKeybinds{
		Keybinds: m.getKeybinds(),
	})
}

func (m CommandList) Update(msg tea.Msg) (CommandList, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.filtering {
			switch msg.Type {
			case tea.KeyRunes, tea.KeySpace:
				m.filteredBy += string(msg.Runes)
			case tea.KeyBackspace:
				if len(m.filteredBy) > 0 {
					m.filteredBy = m.filteredBy[:len(m.filteredBy)-1]
				}
			}

			filtered := []dto.Command{}
			for _, e := range m.cmds {
				a := strings.ToLower(e.Title)
				b := strings.ToLower(m.filteredBy)
				if strings.Contains(a, b) {
					filtered = append(filtered, e)
				}
			}
			m.filtered = filtered
			m.goToTop()
		}
	case common.CommandListUp:
		m.prev()
	case common.CommandListDown:
		m.next()
	case common.CommandListQuit:
		cmds = append(cmds, tea.Quit)
	case common.CommandListSelect:
		return m, m.getSelected
	case common.CommandListFilterToggle:
		m.filtering = !m.filtering
		cmds = append(cmds, common.SetKeybindsCmd(m.getKeybinds()))
	case common.ContentSectionResize:
		slog.Info("content resize ", "Width", msg.Width, "Height", msg.Height)
		m.width = msg.Width - (commandListHorizontalMargin * 2)
		m.height = msg.Height - (commandListVerticalMargin * 2)
	}

	return m, tea.Batch(cmds...)
}

func (m CommandList) View() string {
	headerView := m.headerView()
	currentTotalView := m.currentTotalView()
	listView := m.commandListView(m.height - lipgloss.Height(headerView) - lipgloss.Height(currentTotalView))

	slog.Info("Allowed height", "value", m.height)
	slog.Info("command_list height", "value", lipgloss.Height(headerView)+lipgloss.Height(listView)+lipgloss.Height(currentTotalView))
	slog.Info("Heights", "headerView", lipgloss.Height(headerView),
		"currentTotalView", lipgloss.Height(currentTotalView),
		"listView", lipgloss.Height(listView))

	return commandListStyle.Render(lipgloss.JoinVertical(0, headerView, listView, currentTotalView))
}

func (m CommandList) commandListView(height int) string {
	entries := []string{}
	pageSize := min(len(m.filtered)-m.top, m.determinePageSize())

	// The '- 1' is because each list entry have padding on the left
	maxEntryWidth := m.width - (commandListHorizontalMargin * 2) - 1 - 3

	for i := 0; i < pageSize; i++ {
		e := m.filtered[m.top+i]
		title := common.ElipsisizeText(e.Title, maxEntryWidth)
		command := entryDescriptionStyle.Render(common.ElipsisizeText(e.Cmd, maxEntryWidth-entryDescriptionStyle.GetMarginLeft()))

		entryContent := lipgloss.JoinVertical(0, title, command)

		var entryView string
		if m.top+i == m.selected && !m.filtering {
			entryView = selectedEntryStyle.Render(entryContent)
		} else {
			entryView = entryStyle.Render(entryContent)
		}

		entries = append(entries, entryView)
	}

	listView := ""

	for _, e := range entries {
		listView += entryContainerStyle.Render(e) + "\n"
	}

	return listStyle(height, m.width).MaxHeight(height).Render(listView)
}

func (m CommandList) currentTotalView() string {
	return fmt.Sprintf("%d/%d", m.selected+1, len(m.cmds))
}

func (m CommandList) headerView() string {
	var headerView string
	if m.filtering || m.filteredBy != "" {
		key := filterKeyStyle.Render("Filter: ")
		value := filterValueStyle.Render(m.filteredBy)
		var cursor string
		if m.filtering {
			cursor = filterValueStyle.Render("_")
		} else {
			cursor = ""
		}
		headerView = lipgloss.JoinHorizontal(0, key, value, cursor)
	} else {
		headerView = listTitleStyle.Render("Select command...")
	}

	if m.filtering {
		headerView = filterInputActiveStyle.Render(headerView)
	}

	return headerView
}
