// Package channel includes top-level types and interfaces for the channel cnotents UI.
package channel

import (
	"github.com/cceckman/discoirc/data"
	"github.com/marcusolsson/tui-go"
)

// EventRenderer is a function that converts a DiscoIRC event
// (e.g. message) into an tui.Widget suitable for display.
type EventRenderer func(data.Event) tui.Widget

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

	Attach(UIController)
}

// UIController is a type which can receive updates from a view.
type UIController interface {
	// Accepts input from the user. Must be non-blocking.
	Input(string)

	// Resize indicates a change in the number of lines available for display.
	// Must be non-blocking.
	Resize(n int)

	// TODO: Deferred: Scrolling
	// TODO: Deferred: Localization of connection / presence
}

// ModelController is a type which can receive updates from a Model.
type ModelController interface {
	// UpdateMeta indicates a change in the channel state.
	UpdateMeta(data.Channel)

	// UpdateContents indicates a new Event has arrived.
	UpdateContents(data.Event)
}

type Controller interface {
	UIController
	ModelController
}

// Model implements the Model of a channel.
type Model interface {
	// Returns up to N events ending at this ID.
	EventsEndingAt(end data.EventID, n int) []data.Event
	// TODO: maybe use EventsList instead

	// Send sends the message to the channel.
	Send(string) error

	// Attach uses the ModelController for future updates.
	Attach(ModelController)
}

