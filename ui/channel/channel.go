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

	// (Asynchronous) scroll
	Scroll(up bool)

	// TODO: Deferred: Localization of connection / presence
}

type Model interface {
	// Includes nick, channelMode, connection state, topic
	Channel(ctx context.Context) <-chan data.Channel
	// Returns up to N events starting at EpochId
	EventsStartingAt(start data.EventID, n int) []data.Event
	// Returns up to N events ending at EpochId
	EventsEndingAt(end data.EventID, n int) []data.Event
	// Streams events starting at EpochId
	Follow(ctx context.Context, start data.EventID) <-chan data.Event
	Send(e data.Event) error
}


