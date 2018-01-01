// Package backend defines the types that UI components can use to get updates
// on chat state.
//
// There are a few different planned implementations: "demo", which generates
// exemplary events / state internally; "local", which starts IRC clients within
// the process; and "daemon", which connects to another process that terminates
// the IRC connections, performs logging, etc.
package backend

import (
	"context"

	"github.com/cceckman/discoirc/data"
)

// A DataPublisher allows UI components to subscribe to updates.
type DataPublisher interface {
	Subscribe(context.Context, StateReceiver)
	SubscribeFiltered(FilteredStateReceiver)
}

// StateReceiver receives updates about one or more networks and channels.
type StateReceiver interface {
	UpdateNetwork(data.NetworkState)
	UpdateChannel(data.ChannelState)
}

// FilteredStateReceiver only receives updates for a particular channel and its network.
type FilteredStateReceiver interface {
	StateReceiver

	Filter() (network, channel string)
}

// EventsArchive allows lookup of previous event entries.
type EventsArchive interface {
	EventsBefore(n int, last data.EventID) []data.Event
}

// Backend supports notification on new events, and lookup of prevents events.
type Backend interface {
	DataPublisher
	EventsArchive
}
