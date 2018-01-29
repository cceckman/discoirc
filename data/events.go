package data

import (
	"fmt"
	"sort"
)

// Scope is the scope in which the event occurred.
type Scope struct {
	// Net is the network this event occurred in, or an empty string
	// if it's a discoirc-internal event.
	Net string
	// Name is the name of user / channel where this event occurred.
	Name string
}

// Filter is a pattern for mmatching events.
type Filter struct {
	Scope
	MatchNet, MatchName bool
}

// Match checks if the given scope is within the filter.
func (f *Filter) Match(s Scope) bool {
	net := !f.MatchNet || (s.Net == f.Net)
	name := !f.MatchName || (s.Name == f.Name)
	return net && name
}

// Seq identifies the order of an Event within a Scope.
type Seq int64

// EventID is an identifier for an event.
type EventID struct {
	Scope
	Seq
}

// Event is something that occurred in the IRC client.
type Event interface {
	fmt.Stringer

	ID() *EventID
}

// EventList implements the Events interface for an slice of Events.
type EventList []Event

// SelectSizeMax selects at most n Events, ending at max.
func (e EventList) SelectSizeMax(n int, max Seq) EventList {
	// Find the first element > Max
	end := sort.Search(len(e), func(i int) bool {
		return e[i].ID().Seq > max
	})

	start := end - int(n)
	if start < 0 {
		start = 0
	}
	return e[start:end]
}
