package common

type SelectedCommandEntry struct {
	Cmd         string
	Description string
	DisplayName string
}

type LoadViewport struct{}

type Liveoutput_Quit struct{}

type Confirmation_Selected struct {
	Selected bool
}
