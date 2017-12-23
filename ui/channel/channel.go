// Package channel includes top-level types and interfaces for the channel cnotents UI.
package channel

import (
	"context"
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
	SetPresence(string)
	SetMode(string)
	SetEvents([]data.Event)

	// SetRenderer passes in the function used to render Events in
	// the channel contents display.
	SetRenderer(EventRenderer)

	// Attach indicates the Controller should be used for responses to UI events.
	Attach(Controller)
}

// A Controller handles receiving inputs from a View and updating the View with new contents.
type Controller interface {
	// Accepts input from the user. Non-blocking; safe to run from UI thread.
	Input(string)

	// Resize indicates the number of lines now available for messages.
	Resize(n int)

	// TODO: Deferred: Scrolling
	// TODO: Deferred: Localization of connection / presence
}

type Model interface {
	// Channel reports metadata about the channel.
	// It MUST return an initial value.
	Channel(ctx context.Context) <-chan data.Channel
	// Returns up to N events ending at this ID
	EventsEndingAt(end data.EventID, n int) []data.Event
	// TODO: use EventsList instead

	// Receives new events as they come in.
	// MUST return the most recent event, if any, when initialized.
	Follow(ctx context.Context) <-chan data.Event
	Send(string) error
}

