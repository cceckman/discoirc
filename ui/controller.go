package ui

import (
	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/cceckman/discoirc/ui/client"
)

// UI is the subset of the tui.UI interface that the Controller uses directly or passes through
type UI interface {
	Update(func())
	SetWidget(tui.Widget)
	Quit()
}

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
	channel.New(network, target, c, c.backend)
}

// ActivateClient closes the current view, and replaces it with a view of all
// active sessions of this client.
// Must be run from the UI thread.
func (c *Controller) ActivateClient() {
	client.New(c, c.backend)
}

// Update runns the update in the UI thread.
// Unlike the underlying tui library, this Update call is synchronous; it only
// completes after the callback is run.
func (c *Controller) Update(f func()) {
	blk := make(chan struct{})
	c.UI.Update(func() {
		f()
		close(blk)
	})
	<-blk
}
