package common

import (
	"jackob101/run/common/dto"

	tea "github.com/charmbracelet/bubbletea"
)

type SelectedCommandEntry struct {
	Cmd         string
	Description string
	DisplayName string
}

type LoadViewport struct{}

type ContentSectionResize struct {
	Width  int
	Height int
}

type (
	LoadedCommands struct {
		Cmds []dto.Command
	}
	ShowErrorScreen struct {
		Err error
	}
)

// Command list

type (
	CommandListDown         struct{}
	CommandListUp           struct{}
	CommandListQuit         struct{}
	CommandListFilter       struct{}
	CommandListSelect       struct{}
	CommandListFilterToggle struct{}
	CommandListSelected     struct {
		Cmd dto.Command
	}
)

type (
	LiveoutputClose           struct{}
	LiveoutputClosed          struct{}
	LiveoutputCommandFinished struct{}
	LiveoutputCommandStop     struct{}
	LiveoutputCommandRestart  struct{}
	LiveoutputToggleDetails   struct{}
	LiveoutputUp              struct{}
	LiveoutputDown            struct{}
	LiveoutputOpenEditor      struct{}
	LiveoutputEditorOpened    struct{}
	LiveoutputEditorClosed    struct{}
)

type (
	ConfirmationDialogLeft     struct{}
	ConfirmationDialogRight    struct{}
	ConfirmationDialogSelect   struct{}
	ConfirmationDialogSelected struct {
		Value bool
	}
)

type SetKeybinds struct {
	Keybinds []Keybind
}

func SetKeybindsCmd(keybinds []Keybind) tea.Cmd {
	return MakeCmd(SetKeybinds{
		Keybinds: keybinds,
	})
}

type Keybind struct {
	Keys        []string
	Description string
	Msg         tea.Msg
}

func NewKeybind(msg tea.Msg, description string, keys ...string) Keybind {
	return Keybind{
		Keys:        keys,
		Description: description,
		Msg:         msg,
	}
}
