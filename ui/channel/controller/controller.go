// Package controller provides the channel Controller.
package controller

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/marcusolsson/tui-go"
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
	metaUpdate chan data.Channel
	newEvent   chan data.EventID
}

// New returns a new Controller.
func New(ctx context.Context, ui UIUpdater, v channel.View, m channel.Model) channel.Controller {
	c := &C{
		ui:         ui,
		view:       v,
		model:      m,
		sizeUpdate: make(chan int, 1),
		metaUpdate: make(chan data.Channel, 1),
		input:      make(chan string),
		newEvent:   make(chan data.EventID),
	}

	go c.awaitMetaUpdate(ctx)
	go c.awaitInput(ctx)
	go c.awaitEvents(ctx)

	if v != nil {
		v.Attach(c)
	}
	if m != nil {
		m.Attach(c)
	}

	return c
}

// TODO: Support localization
// updateMeta updates a channel.View with channel metadata.
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

	v.SetName(d.Name)
	if d.Presence == data.NotPresent {
		v.SetMode("∅")
	} else {
		v.SetMode(d.Mode)
	}
}

func (c *C) awaitMetaUpdate(ctx context.Context) {
	// newData allows any pending updates to always get the most up-to-date
	// data.Channel.
	newData := make(chan data.Channel, 1)
	defer close(newData)
	for {
		select {
		case <-ctx.Done():
			return
		case m := <-c.metaUpdate:
			select {
			case <-newData:
				// Flushed existing data; give it better data.
				newData <- m
				// Don't launch a new updater thread, there's one waiting.
			case newData <- m:
				// No existing data; launch a new updater thread.
				go c.ui.Update(func() {
					updateMeta(<-newData, c.view)
				})
			}
		}
	}
}

func (c *C) awaitEvents(ctx context.Context) {
	size := 0             // desired N of events to display
	var last data.EventID // last event in the display

	for {
		fetch := false // need to refresh events?
		select {
		case <-ctx.Done():
			return
		case size = <-c.sizeUpdate:
			fetch = true
		case last = <-c.newEvent:
			// TODO: be more efficient; don't do a full-fetch if only an
			// update is needed.
			// TODO: support scrolling
			fetch = true
		}
		// TODO perform fetch asynchronously; assume it may take
		// a relatively long time.
		if fetch {
			glog.V(1).Infof("Controller fetching new contents: %d after %v", size, last)
			events := c.model.EventsEndingAt(last, size)
			await := make(chan struct{})
			// TODO perform update asynchronously
			c.ui.Update(func() {
				c.view.SetEvents(events)
				await <- struct{}{}
			})
			<-await
		}
	}
}

// awaitInput is a thread that handles queueing inputted messages for processing.
func (c *C) awaitInput(ctx context.Context) {
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

// UpdateContents indicates a new Event has arrived.
func (c *C) UpdateContents(d data.Event) {
	glog.V(1).Infof("UpdateContents received message %v", d)
	c.newEvent <- d.EventID
}

// Input accepts input from the user.
func (c *C) Input(s string) {
	c.input <- s
}

// Resize notes a change in the number of events displayed.
func (c *C) Resize(n int) {
	glog.V(1).Infof("Controller got resize: %d", n)
	select {
	case c.sizeUpdate <- n:
		// Sent update.
	case <-c.sizeUpdate:
		c.sizeUpdate <- n
	}
}

// UpdateMeta receives an update about the channel's state.
func (c *C) UpdateMeta(d data.Channel) {
	select {
	case c.metaUpdate <- d:
		// Sent update.
	case <-c.metaUpdate:
		c.metaUpdate <- d
	}
}
