package ui

import (
	"context"

	"github.com/marcusolsson/tui-go"
)

// UI is the subset of the tui.UI interface that a Controller uses.
type UI interface {
	Update(func())
	SetWidget(tui.Widget)
}

func New(_ context.Context, ui UI) *Controller {
	c := &Controller{
		UI: ui,
	}

	return c
}

// Controller is the master controller for a discoirc UI.
// It manages the lifecycle of other controllers and views.
type Controller struct {
	UI
}

// ActivateChannel closes the current view, and replaces it with a view of the
// given channel in the given network.
func (c *Controller) ActivateChannel(network, channel string) {
	// TODO
}

// ActivateClient closes the current view, and replaces it with a view of all
// active sessions of this client.
func (c *Controller) ActivateClient() {

}
