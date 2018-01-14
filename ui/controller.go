package ui

import (
	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/cceckman/discoirc/ui/client"
)

// UI is the subset of the tui.UI interface that the Controller uses directly or passes through
type UI interface {
	Update(func())
	SetWidget(tui.Widget)
	Quit()
}

// New returns a new Controller.
func New(ui UI, be backend.Backend) *Controller {
	c := &Controller{
		UI:      ui,
		backend: be,
	}

	return c
}

// Controller is the master controller for a discoirc UI.
// It manages the lifecycle of other controllers and views.
type Controller struct {
	UI

	backend backend.Backend
}

// ActivateChannel closes the current view, and replaces it with a view of the
// given channel in the given network.
func (c *Controller) ActivateChannel(network, target string) {
	channel.New(
		data.Scope{Net: network, Name: target},
		c, c.backend,
	)
}

// ActivateClient closes the current view, and replaces it with a view of all
// active sessions of this client.
// Must be run from the UI thread.
func (c *Controller) ActivateClient() {
	client.New(c, c.backend)
}
