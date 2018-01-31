// Package backend defines the types that UI components can use to get updates
// on chat state.
//
// There are a few different planned implementations: "demo", which generates
// exemplary events / state internally; "local", which starts IRC clients within
// the process; and "daemon", which connects to another process that terminates
// the IRC connections, performs logging, etc.
package backend

import (
	"github.com/cceckman/discoirc/data"
)

// A DataPublisher allows components to subscribe to updates.
// Only one subscriber of any sort may be active at a time, from a single
// DataPublisher.
type DataPublisher interface {
	Subscribe(Receiver)
}

// Receiver receives updates about one or more networks and channels.
type Receiver interface {
	Receive(data.Event)

	Filter() data.Filter
}

// EventsArchive allows lookup of previous event entries.
// TODO: Use EventID as selector instead.
type EventsArchive interface {
	EventsBefore(s data.Scope, n int, last data.Seq) data.EventList
}

// Sender sends a message on the given network to the given target (channel or user).
type Sender interface {
	Send(s data.Scope, message string)
}

// Backend supports the full set of backend functionality.
type Backend interface {
	DataPublisher
	EventsArchive
	Sender
}
