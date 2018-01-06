package widgets

import (
	"github.com/marcusolsson/tui-go"
)

const splashImg string = `
             ||
             ||
           <><><>
         <><><><><>
        <><><><><><>
        <><><><><><>
        <><><><><><>
         <><><><><>
           <><><>

          discoirc

github.com/cceckman/discoirc

`

// Quitter is a type that supports the ui.Quit() operation.
type Quitter interface {
	Quit()
}

type splash struct {
	tui.Widget

	ui Quitter
}

// OnKeyEvent handles keypress events.
func (s *splash) OnKeyEvent(ev tui.KeyEvent) {
	// TODO: put all of these into a single "isQuitEvent" utility funcion.
	if ev.Key == tui.KeyCtrlC {
		s.ui.Quit()
		return
	}
	s.Widget.OnKeyEvent(ev)
}

// NewSplash returns a new Splash widget.
func NewSplash(ui Quitter) tui.Widget {
	return &splash{
		ui: ui,
		Widget: tui.NewHBox(
			tui.NewSpacer(),
			tui.NewVBox(
				tui.NewSpacer(),
				tui.NewLabel(splashImg),
				tui.NewSpacer(),
			),
			tui.NewSpacer(),
		),
	}
}
