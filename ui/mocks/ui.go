package mocks

import (
	"github.com/marcusolsson/tui-go"
)

func NewUI() *UI {
	ui := &UI{}

	return ui

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
		if ui.Root == nil {
			return
		}
		ui.Root.OnKeyEvent(ev)
	}
}

// UI implements a subset of the tui.UI functionality for use in tests.
type UI struct {
	Root tui.Widget

	Painter *tui.Painter

	HasQuit bool
}

// Repaint re-renders if the painter and root are not nil.
func (ui *UI) Repaint() {
	if ui.Painter != nil && ui.Root != nil {
		ui.Painter.Repaint(ui.Root)
	}

}

func (ui *UI) Update(f func()) {
	f()
}

func (ui *UI) SetWidget(w tui.Widget) {
	ui.Root = w
	ui.Repaint()
}

func (ui *UI) Quit() {
	ui.HasQuit = true
}
