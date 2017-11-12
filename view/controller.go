// Package view contains a controller that manages other views. Its subpackages
// provide views and controllers for specific views.

package view

import (
	"context"
	"fmt"
	"log"

	"github.com/cceckman/discoirc/model"
	"github.com/cceckman/discoirc/view/channel"
	"github.com/marcusolsson/tui-go"
)

// ConsoleSession is the controller for an instance of the discoirc TUI.
// It provides methods for reconfiguring the UI to different views,
// as well as interfaces common to those views.
type ConsoleSession struct {
	tui.UI
	*log.Logger
	model.Client

	viewRequest chan ViewRequest
}

func NewConsoleSession(logger *log.Logger, client model.Client) *ConsoleSession {
	r := &ConsoleSession{
		UI:     tui.New(splashRequest(0).New(nil, nil)),
		Logger: logger,
		Client: client,
		// Use a nonblocking buffer s.t. updates to the view are nonblocking.
		viewRequest: make(chan ViewRequest, 1),
	}
	go r.handleViewChange()

	return r
}

func (cs *ConsoleSession) handleViewChange() {
	ctx, cancel := context.WithCancel(context.Background())

	for req := range cs.viewRequest {
		cs.Printf("laying out new view: %s", req)
		// stop old threads
		cancel()
		ctx, cancel = context.WithCancel(context.Background())
		// Set the new view.
		newRoot := req.New(ctx, cs)
		await := make(chan struct{})
		cs.UI.Update(func() {
			close(await)
			cs.UI.SetWidget(newRoot)
			// Ensure that exit is accounted for.
			cs.UI.SetKeybinding("Esc", func() { cs.UI.Quit() })
		})
		<-await
	}
}

// OpenChannel requests that the view change to that of the given channel.
// It's asynchronous, safe call from any thread.
func (cs *ConsoleSession) OpenChannel(network, channel string) {
	cs.viewRequest <- &channelViewRequest{
		Network: network,
		Channel: channel,
	}
}

// Splash sets the root view to the slash screen.
func (cs *ConsoleSession) Splash() {
	cs.viewRequest <- splashRequest(0)
}

// A ViewRequest is a request to jump to a specific view.
type ViewRequest interface {
	fmt.Stringer

	// New constructs a new widget for the TUI to use as the root.
	New(context.Context, *ConsoleSession) tui.Widget
}

// A channelViewRequest is a request to open a channel view.
type channelViewRequest struct {
	Network, Channel string
}

func (c *channelViewRequest) String() string {
	return fmt.Sprintf("channel view: %s / %s", c.Network, c.Channel)
}

func (c *channelViewRequest) New(ctx context.Context, cs *ConsoleSession) tui.Widget {
	return channel.New(ctx, cs.Logger, cs.UI, cs.Client, c.Network, c.Channel)
}
