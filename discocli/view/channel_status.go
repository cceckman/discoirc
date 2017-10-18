package view

import (
	"context"
	"fmt"

	"github.com/jroimartin/gocui"
)

// ChannelStatus is the ViewModel for the Channel window's status line.
type ChannelStatus struct {
	*Channel
	laidOut chan struct{}
}

var _ gocui.Manager = &ChannelStatus{}

// Layout sets up the ChannelStatus view.
func (c *ChannelStatus) Layout(g *gocui.Gui) error {
	// Full width; no frame. Position relative to ChannelInputView.
	input, err := g.View(ChannelInputView)
	if err != nil {
		return err
	}
	_, y := input.Origin()
	maxX, _ := g.Size()
	ax, ay, bx, by := -1, y-2, maxX, y-1
	// No border at the bottom of the terminal, full width.
	v, err := g.SetView(ChannelStatusView, ax, ay, bx, by)
	switch err {
	case nil:
		return nil
	case gocui.ErrUnknownView:
		c.Log.Printf("%s [start] initial setup", ChannelStatusView)
		defer c.Log.Printf("%s [done] initial setup", ChannelStatusView)
		defer close(c.laidOut)
		v.Frame = false
		v.Editable = false
		// TODO more color customization
		v.BgColor, v.FgColor = gocui.ColorWhite, gocui.ColorBlack

	default:
		return err
	}
	return nil
}

// NewStatus creates a new ViewModel for the status bar and starts threads to update it.
func (vm *Channel) NewStatus() gocui.Manager {
	result := &ChannelStatus{
		Channel: vm,
		laidOut: make(chan struct{}),
	}
	go result.Listen()
	return result
}

// Listen waits for a backend connection, then listens for updates to the topic until the UI is gone.
func (c *ChannelStatus) Listen() {
	<-c.connected
	<-c.laidOut
	ctx, cancel := context.WithCancel(context.Background())

	for n := range c.channel.Await(ctx) {
		topic := n.Topic
		// TODO: There's a subtle race here; if topic updates are arriving quickly, there may be more than
		// one Update queued- and they may execute in any order, so the topic may go backards.
		c.Gui.Update(func(g *gocui.Gui) error {
			v, err := g.View(ChannelStatusView)
			switch err {
			case gocui.ErrUnknownView:
				cancel()
			case nil:
				//pass
			default:
				return err
			}
			v.Clear()
			if _, err := fmt.Fprintf(v, "[%s] %s | %s", c.Connection, c.Channel, topic); err != nil {
				return err
			}
			return nil
		})
	}
}
