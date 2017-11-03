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

type Events interface {
	SelectSize(uint) []Event
	SelectSizeMax(uint, EventID) []Event
	SelectMinSize(EventID, uint) []Event
	SelectMinMax(EventID, EventID) []Event
}

// NewEvents returns a new EventList, ensuring that it is sorted.
func NewEvents(es []Event) EventList {
	r := make([]Event, len(es))
	copy(r, es)
	sort.Sort(EventList(r))
	return r
}

type EventList []Event

// SelectSize selects the most recent n events.
func (e EventList) SelectSize(n uint) []Event {
	start := len(e) - int(n)
	if start < 0 {
		start = 0
	}
	return e[start:]
}

// SelectMaxSize selects at most n Events ending at max.
func (e EventList) SelectSizeMax(n uint, max EventID) []Event {
	// Find the first element > Max
	end := sort.Search(len(e), func(i int) bool {
		if e[i].Epoch == max.Epoch {
			return e[i].Seq > max.Seq
		}
		return e[i].Epoch > max.Epoch
	})

	start := end - int(n)
	if start < 0 {
		start  = 0
	}
	return e[start:end]
}

// SelectMinSize selects at most n Events starting from min.
func (e EventList) SelectMinSize(min EventID, n uint) []Event {
	// Find the first element >= Min
	start := sort.Search(len(e), func(i int) bool {
		if e[i].Epoch == min.Epoch {
			return e[i].Seq >= min.Seq
		}
		return e[i].Epoch > min.Epoch
	})

	end := start + int(n)
	if end > len(e) {
		end = len(e)
	}
	return e[start:end]
}


// Select returns a slice from its receiver with those within the EventRange.
func (e EventList) SelectMinMax(min, max EventID) []Event {
	// Events must already be sorted.

	// Find the first element >= Min
	start := sort.Search(len(e), func(i int) bool {
		if e[i].Epoch == min.Epoch {
			return e[i].Seq >= min.Seq
		}
		return e[i].Epoch > min.Epoch
	})

	// Find the first element >= Max
	end := sort.Search(len(e), func(i int) bool {
		if e[i].Epoch == max.Epoch {
			return e[i].Seq > max.Seq
		}
		return e[i].Epoch > max.Epoch
	})

	return e[start:end]
}

func Less(a EventID, b EventID) bool {
	if a.Epoch == b.Epoch {
		return a.Seq < b.Seq
	}
	if a.Epoch < b.Epoch {
		return true
	}
	return false
}

func (e EventList) Len() int { return len(e) }
func (e EventList) Less(i, j int) bool {
	return Less(e[i].EventID, e[j].EventID)
}
func (e EventList) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

// An Event represents an event in IRC, e.g. a message.
type Event struct {
	EventID

	Contents string
}

func (e *Event) String() string {
	return e.Contents
}
