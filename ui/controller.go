package ui

import (
	"context"

	"github.com/marcusolsson/tui-go"

	clientView "github.com/cceckman/discoirc/ui/client/view"
	"github.com/cceckman/discoirc/ui/widgets"
)

// UI is the subset of the tui.UI interface that the Controller uses directly or passes through
type UI interface {
	Update(func())
	SetWidget(tui.Widget)
	Quit()
}

func New(ctx context.Context, ui UI) *Controller {
	c := &Controller{
		UI:       ui,
		toClient: make(chan struct{}),
	}

	go c.mainLoop(ctx)

	return c
}

// Controller is the master controller for a discoirc UI.
// It manages the lifecycle of other controllers and views.
type Controller struct {
	UI

	toClient chan struct{}
}

// ActivateChannel closes the current view, and replaces it with a view of the
// given channel in the given network.
func (c *Controller) ActivateChannel(network, channel string) {
	// TODO
}

// ActivateClient closes the current view, and replaces it with a view of all
// active sessions of this client.
func (c *Controller) ActivateClient() {
	c.toClient <- struct{}{}
}

// mainLoop runs outside the UI thread so that it can manage the
// lifecycle of background threads (via Context).
func (c *Controller) mainLoop(ctx context.Context) {
	join(c.UI, func() {
		c.UI.SetWidget(widgets.NewSplash(c))
	})
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.toClient:
			join(c.UI, c.activateClient)
		}
		// TODO: Clean up previous threads (context)
	}
}

// join is a utility function that calls Update(f), but waits
// until its completion. It's the synchronous version of ui.Update.
func join(ui UI, f func()) {
	blk := make(chan struct{})
	ui.Update(func(){
		f()
		close(blk)
	})
	<-blk
}

func (c *Controller) activateClient() {
	// TODO: Reset any global keybindings
	client := clientView.New()
	client.Attach(c)
	c.SetWidget(client)
}
