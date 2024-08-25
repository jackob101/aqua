package main

import (
	"bufio"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type liveoutput struct {
	sub                chan string
	commandDisplayName string
	command            string
	lines              []string
	finished           bool
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

func (m liveoutput) Init() tea.Cmd {
	return tea.Batch(wrapMsg(tea.ClearScreen()),
		m.listenForNewline(),
		m.waitForNewline())
}

func (m liveoutput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch message := msg.(type) {
	case tea.KeyMsg:
		return m, wrapMsg(loClosed{})
	case newline:
		m.lines = append(m.lines, message.content)
		return m, m.waitForNewline()
	case loFinished:
		m.finished = true
		return m, nil
	}
	return m, nil
}

func (m liveoutput) View() string {
	s := m.command
	s = s + "\n----Output----\n\n"
	for _, v := range m.lines {
		s = s + v + "\n"
	}
	if m.finished {
		s = s + "\n------------\nPress any key to continue\n"
	}
	return s
}
