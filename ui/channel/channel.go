// Package channel includes top-level types and interfaces for the channel cnotents UI.
package channel

import (
	"context"
	"github.com/cceckman/discoirc/data"
	"github.com/marcusolsson/tui-go"
)

type EventRenderer func(data.Event) tui.Widget

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
}

type Controller interface {
	// Parse input string, do appropriate stuff to backend
	Input(string)

	// Resize indicates the number of lines now available for messages.
	Resize(n int)

	// (Asynchronous) scroll
	Scroll(up bool)

	// TODO: Deferred: Localization of connection / presence
}
