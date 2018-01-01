package channel

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/glog"
	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/data"
)

// UIController is the interface that a higher-level controller must provide.
// All of its methods should be accessed via Update.
type UIController interface {
	// Update runs the provided closure in the UI event loop.
	Update(func())

	// SetWidget sets the provided widget as the root of the UI.
	SetWidget(tui.Widget)

	// ActivateClient switches the global view to the Client view.
	ActivateClient()

	Quit()
}

// View is a user-facing display of an IRC channel.
type View interface {
	tui.Widget

	SetTopic(string)
	SetNick(string)
	SetConnection(string)
	SetName(string)
	SetMode(string)
	SetEvents([]data.Event)

	// SetRenderer passes in the function used to render Events in
	// the channel contents display.
	SetRenderer(EventRenderer)

	Attach(Controller)
}

// Model holds and updates the state of a channel.
type Model interface {
	// Returns up to N events ending at this ID.
	EventsEndingAt(end data.EventID, n int) []data.Event
	// TODO: maybe use EventsList instead

	// Send sends the message to the channel.
	Send(string) error

	// Attach uses the ModelController for future updates.
	Attach(Controller)
}

var _ Controller = &controller{}

// controller implements a channel Controller.
type controller struct {
	ui    UIController
	view  View
	model Model

	// Async communication
	sizeUpdate chan int
	input      chan string
	metaUpdate chan data.Channel
	newEvent   chan data.EventID
}

// New returns a new Controller.
func New(ctx context.Context, ui UIController, v View, m Model) Controller {
	c := &controller{
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

	c.ui.Update(func() {
		c.ui.SetWidget(c.view)
	})

	return c
}

func (c *controller) Quit() {
	c.ui.Update(func() {
		c.ui.Quit()
	})
}

// TODO: Support localization
// updateMeta updates a View with channel metadata.
func updateMeta(d data.Channel, v View) {
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

func (c *controller) awaitMetaUpdate(ctx context.Context) {
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

func (c *controller) awaitEvents(ctx context.Context) {
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
func (c *controller) awaitInput(ctx context.Context) {
	nextMessage := make(chan string)
	defer close(nextMessage)
	go func() {
		for m := range nextMessage {
			c.model.Send(m)
		}
	}()

	queue := []string{}

	// TODO: Don't do this with an inline function.
	// Has to be at the moment because it closes over queue,
	// but there's surely a better way to handle.
	handleMessage := func(m string) {
		lower := strings.ToLower(m)

		if strings.HasPrefix(lower, "/client") {
			c.ui.Update(func() {
				c.ui.ActivateClient()
			})
			return
		}
		if strings.HasPrefix(lower, "/quit") {
			c.Quit()
		}

		queue = append(queue, m)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case m := <-c.input:
			handleMessage(m)
		}
		if len(queue) > 0 {
			select {
			case <-ctx.Done():
				return
			case m := <-c.input:
				handleMessage(m)
			case nextMessage <- queue[0]:
				queue = queue[1:]
			}
		}
	}
}

// UpdateContents indicates a new Event has arrived.
func (c *controller) UpdateContents(d data.Event) {
	glog.V(1).Infof("UpdateContents received message %v", d)
	c.newEvent <- d.EventID
}

// Input accepts input from the user.
func (c *controller) Input(s string) {
	c.input <- s
}

// Resize notes a change in the number of events displayed.
func (c *controller) Resize(n int) {
	glog.V(1).Infof("Controller got resize: %d", n)
	select {
	case c.sizeUpdate <- n:
		// Sent update.
	case <-c.sizeUpdate:
		c.sizeUpdate <- n
	}
}

// UpdateMeta receives an update about the channel's state.
func (c *controller) UpdateMeta(d data.Channel) {
	select {
	case c.metaUpdate <- d:
		// Sent update.
	case <-c.metaUpdate:
		c.metaUpdate <- d
	}
}