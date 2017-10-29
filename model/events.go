package model

import (
	"sort"
)

// An EventID is a unique, sequenceable identifier for an event.
// Within a scope, e.g. a channel, events are lexicographically ordered-
// by epoch, then by sequence.
// An Epoch may be ticked over when a log file rolls over, when a server
// disconnects, or in other circumstances. Usually, it represents when there was
// a potential discontinuity in events.
// Epochs may be negative because they may extend indefinitely into the past,
// e.g. via logfiles or an external log server.
type EventID struct {
	Epoch int
	Seq   uint
}

// An EventRange selects a range of events.
type EventRange struct {
	Min, Max EventID
}

// NewEvents returns a new Event object, ensuring that it is sorted.
func NewEvents(es []Event) Events {
	r := make([]Event, len(es))
	copy(r, es)
	sort.Sort(Events(r))
	return Events(r)
}

// Events is an ordered set of events.
type Events []Event

// Select returns a slice from its receiver with those within the EventRange.
func (e Events) Select(r EventRange) Events {
	// Events must already be sorted.

	// Find the first element >= Min
	start := sort.Search(len(e), func(i int) bool {
		if e[i].ID.Epoch == r.Min.Epoch {
			return e[i].ID.Seq >= r.Min.Seq
		}
		return e[i].ID.Epoch > r.Min.Epoch
	})

	// Find the first element >= Max
	end := sort.Search(len(e), func(i int) bool {
		if e[i].ID.Epoch == r.Max.Epoch {
			return e[i].ID.Seq > r.Max.Seq
		}
		return e[i].ID.Epoch > r.Max.Epoch
	})

	return Events(e[start:end])
}

func (e Events) Len() int { return len(e) }
func (e Events) Less(i, j int) bool {
	a, b := e[i].ID, e[j].ID
	if a.Epoch == b.Epoch {
		return a.Seq < b.Seq
	}
	if a.Epoch < b.Epoch {
		return true
	}
	return false
}
func (e Events) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

// An Event represents an event in IRC, e.g. a message.
type Event struct {
	ID EventID

	Contents string
}

func (e *Event) String() string {
	return e.Contents
}
