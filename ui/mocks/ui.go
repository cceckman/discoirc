package mocks

import (
	"github.com/marcusolsson/tui-go"
)

func NewUI() *UI {
	return &UI{
		UpdateCounter: NewUpdateCounter(),
	}

}

// UI implements a subset of the tui.UI functionality for use in tests.
type UI struct {
	*UpdateCounter

	Root tui.Widget

	HasQuit bool
}

func (ui *UI) SetWidget(w tui.Widget) {
	ui.Root = w
}

func (ui *UI) Quit() {
	ui.HasQuit = true
}
