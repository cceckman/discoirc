// Package channel provides the channel contents UI.
package channel

import (
	"github.com/cceckman/discoirc/data"
	"github.com/marcusolsson/tui-go"
)

// EventRenderer is a function that converts a DiscoIRC event
// (e.g. message) into an tui.Widget suitable for display.
type EventRenderer func(data.Event) tui.Widget

// UIController is a type which can receive updates from a view.
type UIController interface {
	// Accepts input from the user. Must be non-blocking.
	Input(string)

	// Resize indicates a change in the number of lines available for display.
	// Must be non-blocking.
	Resize(n int)

	// TODO: Deferred: Scrolling
	// TODO: Deferred: Localization of connection / presence

	Quit()
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
