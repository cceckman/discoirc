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

// Match chekcs the scope against
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
// SortEvents produces an EventList from the Events.
func SortEvents(es []Event) EventList {
	r := make([]Event, len(es))
	copy(r, es)
	sort.Sort(EventList(r))
	return r
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

// Len implements sort.Interface for EventList.
func (e EventList) Len() int { return len(e) }

// Less implements sort.Interface for EventList
func (e EventList) Less(i, j int) bool {
	return e[i].ID().Seq < e[j].ID().Seq
}

// Swap implements sort.Interface for EventList
func (e EventList) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
