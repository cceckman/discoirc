package mocks

import (
	"github.com/marcusolsson/tui-go"
)

func NewUI() *UI {
	return &UI{
		UpdateCounter: NewUpdateCounter(),
	}

}

// Type produces an approprate keypress against its root for each character in the input.
// Each is a separate event in the UI thread.
func (ui *UI) Type(s string) {
	for _, rn := range s {
		var ev tui.KeyEvent
		if rn != '\n' {
			ev = tui.KeyEvent{
				Key:  tui.KeyRune,
				Rune: rn,
			}
		} else {
			ev = tui.KeyEvent{
				Key: tui.KeyEnter,
			}
		}
		ui.Update(func(){
			if ui.Root == nil {
				return
			}
			ui.Root.OnKeyEvent(ev)
		})
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
