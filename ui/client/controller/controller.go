package controller

import (
	"context"

	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/ui/client"
)

// UI is the subset of tui.UI that a controller requries.
type UI interface {
	SetWidget(tui.Widget)
	Update(func())
}

var _ UI = tui.UI(nil)

func New(_ context.Context, _ UI) client.UIController {
	return &C{}

}

// C implements a Controller for the Client view.
type C struct { }

func (ctl *C) ActivateChannel(network, name string) {
	// TODO
}

var _ client.UIController = &C{}

