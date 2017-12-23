// Package controller provides the channel Controller.
package controller

import (
	"context"
	"fmt"

	"github.com/marcusolsson/tui-go"
	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
)

// UIUpdater is a subset of the tui.UI interface- just the bit that update a UI.
type UIUpdater interface {
	Update(func())
}

var _ UIUpdater = tui.UI(nil)


var _ channel.Controller = &C{}

// C implements a channel Controller.
type C struct {
	ui    UIUpdater
	view  channel.View
	model channel.Model

	// Async communication
	sizeUpdate chan int
	input      chan string
}

// New returns a new Controller.
func New(ctx context.Context, ui UIUpdater, v channel.View, m channel.Model) channel.Controller {
	c := &C{
		ui:         ui,
		view:       v,
		model:      m,
		sizeUpdate: make(chan int, 1),
		input:      make(chan string, 1),
	}

	if v != nil {
		v.Attach(c)
	}

	go c.updateEvents(ctx)
	go c.updateChannelMeta(ctx)
	go c.handleInput(ctx)
	return c
}

// TODO: Support localization
func updateMeta(d data.Channel, v channel.View) {
	v.SetTopic(d.Topic)

	connStrings := map[data.ConnectionState]string{
		data.Disconnected: "∅",
		data.Connecting:   "…",
		data.Connected:    "✓",
	}
	conn := d.Connection
	v.SetConnection(fmt.Sprintf("%s: %s", conn.Network, connStrings[conn.State]))
	v.SetNick(conn.Nick)

	v.SetPresence(d.Name)
	if d.Presence == data.NotPresent {
		v.SetMode("[parted]")
	} else {
		v.SetMode(d.Mode)
	}
}

func (c *C) updateChannelMeta(ctx context.Context) {
	metadata := c.model.Channel(ctx)

	// newData allows any pending updates to always get the most up-to-date
	// data.Channel.
	newData := make(chan data.Channel, 1)
	defer close(newData)
	for {
		select {
		case <-ctx.Done():
			return
		case m := <-metadata:
			select {
			case <-newData:
				// Updater thread was waiting on an update. Give it better data.
				newData <- m
			case newData <- m:
				// can write to it; need to launch a new thread.
				c.ui.Update(func() {
					updateMeta(<-newData, c.view)
				})
			}
		}
	}
}

// updateEvents is a thread that handles updating the events display.
func (c *C) updateEvents(ctx context.Context) {
	follow := c.model.Follow(ctx)

	size := 1             // desired N of events to display
	var last data.EventID // last event in the display

	for {
		fetch := false // need to refresh events?
		select {
		case <-ctx.Done():
			return
		case size = <-c.sizeUpdate:
			fetch = true
		case ev := <-follow:
			// TODO: be more efficient; don't do a full-fetch if only an
			// update is needed.
			last = ev.EventID
			fetch = true
		}

		// TODO perform fetch asynchronously; assume it may take
		// a relatively long time.
		if fetch {
			events := c.model.EventsEndingAt(last, size)
			// TODO perform update asynchronously
			c.ui.Update(func() {
				c.view.SetEvents(events)
			})
		}
	}
}

// handleInput is a thread that handles queueing inputted messages for processing.
func (c *C) handleInput(ctx context.Context) {
	nextMessage := make(chan string)
	defer close(nextMessage)
	go func() {
		for m := range nextMessage {
			c.model.Send(m)
		}
	}()

	queue := []string{}

	for {
		select {
		case <-ctx.Done():
			return
		case m := <-c.input:
			// TODO: Parse for non-sending operations before sending to model.
			queue = append(queue, m)
		}
		if len(queue) > 0 {
			select {
			case <-ctx.Done():
				return
			case m := <-c.input:
				queue = append(queue, m)
			case nextMessage <- queue[0]:
				queue = queue[1:]
			}
		}
	}
}

// Input accepts input from the user.
func (c *C) Input(s string) {
	c.input <- s
}

// Resize notes a change in the number of events displayed.
func (c *C) Resize(n int) {
	select {
	case c.sizeUpdate <- n:
		// pass
	case <- c.sizeUpdate:
		c.sizeUpdate <- n
	}
}
