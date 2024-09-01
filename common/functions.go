package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func MakeCmd(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func Max(e ...int) int {
	maxValue := 0
	for _, e := range e {
		if e > maxValue {
			maxValue = e
		}
	}
	return maxValue
}

func Map[T any, B any](entries []T, mapper func(T) B) []B {
	result := []B{}
	for _, e := range entries {
		result = append(result, mapper(e))
	}
	return result
}

func Reduce[T any, B any](entries []T, init B, reducer func(B, T) B) B {
	var result B
	for _, e := range entries {
		result = reducer(result, e)
	}
	return result
}

func ElipsisizeText(text string, maxWidth int) string {
	if lipgloss.Width(text) >= maxWidth {
		return text[:maxWidth] + "..."
	} else {
		return text
	}
}
