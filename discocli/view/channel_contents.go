package view

import (
	"context"
	"fmt"

	"github.com/jroimartin/gocui"
)

// ChannelContents is the ViewModel for the Channel window's status line.
type ChannelContents struct {
	*Channel
	laidOut chan struct{}
}

var _ gocui.Manager = &ChannelContents{}

// Layout sets up the ChannelContents view.
func (c *ChannelContents) Layout(g *gocui.Gui) error {
	// Full width; no frame. Position relative to ChannelStatusView.
	status, err := g.View(ChannelStatusView)
	if err != nil {
		return fmt.Errorf("error in laying out %s: could not find status bar: %v", ChannelContentsView, err)
	}

	_, oy := status.Origin()
	_, dy := status.Size()
	maxX, _ := g.Size()

	ax, ay, bx, by := -1, oy+dy-2, maxX, -1

	c.Log.Printf("%s: laying out at (%d, %d) (%d, %d)", ChannelContentsView, ax, ay, bx, by)
	// No border at the bottom of the terminal, full width.
	v, err := g.SetView(ChannelContentsView, ax, ay, bx, by)
	switch err {
	case nil:
		// pass
	case gocui.ErrUnknownView:
		c.Log.Printf("%s [start] initial setup", ChannelContentsView)
		defer c.Log.Printf("%s [done] initial setup", ChannelContentsView)
		defer close(c.laidOut)
		v.Frame = false
		v.Editable = false
		// TODO better handling of scroll / etc. behavior
		v.Autoscroll = true
	default:
		return fmt.Errorf("error in laying out %s: %v", ChannelContentsView, err)
	}

	return nil
}

// NewContents creates a new ViewModel for the status bar and starts threads to update it.
func (vm *Channel) NewContents() gocui.Manager {
	result := &ChannelContents{
		Channel: vm,
		laidOut: make(chan struct{}),
	}
	go result.Listen()
	return result
}

// Listen waits for a backend connection, then listens for updates to the messages until the UI is gone.
func (c *ChannelContents) Listen() {
	<-c.connected
	<-c.laidOut
	ctx, cancel := context.WithCancel(context.Background())

	count := 0

	for n := range c.channel.Await(ctx) {
		if n.Messages == count {
			continue
		}
		count = n.Messages

		// Only enqueue one update at a time, to make sure we don't go backwards in time.
		done := make(chan struct{})
		c.Gui.Update(func(g *gocui.Gui) error {
			defer close(done)
			v, err := g.View(ChannelContentsView)
			switch {
			case err == gocui.ErrUnknownView:
				cancel()
				return nil
			case err != nil:
				return err
			}
			// TODO: Refactor this controller, s.t. this doesn't take place in the UI thread.
			// I think the way to do that is to have a select on (resize events, new message events).
			_, lines := v.Size()
			// TODO: Support scrollback; keep relative position of the view.
			messages := c.channel.GetMessages(0, uint(lines))
			v.Clear()
			for _, m := range messages {
				fmt.Fprintln(v, m)
			}
			return nil
		})
		<-done
	}
}
