// Package channel includes the Widgets for the channel view.
package channel

import (
	"context"
	"github.com/cceckman/discoirc/model"
	"github.com/marcusolsson/tui-go"
)

func New(ctx context.Context, client model.Client, network, channel string) tui.Widget {
	w := &Root{
		Input : tui.NewEntry(),
		Status: tui.NewStatusBar(""),
		Contents: tui.NewVBox(),
	}

	// Layout
	w.Contents.SetSizePolicy(tui.Expanding, tui.Expanding)
	w.Status.SetSizePolicy(tui.Expanding, tui.Preferred)
	w.Input.SetSizePolicy(tui.Expanding, tui.Preferred)
	w.Input.SetFocused(true)

	w.Status.SetPermanentText("status bar")

	w.Widget = tui.NewVBox(w.Contents, w.Status, w.Input)
	return w
}

type Root struct {
	tui.Widget // root widget

	Input *tui.Entry
	Status *tui.StatusBar
	Contents *tui.Box

	Events model.Events
}
