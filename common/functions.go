package common

import tea "github.com/charmbracelet/bubbletea"

func MakeCmd(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
