package model

// An EventId is a unique, sequenceable identifier for an event.
// Within a scope, e.g. a channel, events are strictly ordered: by epoch, then by sequence.
// An Epoch may be ticked over when a log file rolls over, when a server disconnects, etc.
// Epochs may be negative because it may extend indefinitley into the past, e.g. via logfiles
// or an external log server.
type EventId struct {
	Epoch int
	Seq uint
}

// An EventRange is a range of events. It represents all events which may have occurred
// within the given range.
type EventRange struct {
	Min, Max EventId
}


// An Event represents an event in IRC, e.g. a message.
type Event struct {
	id EventId

	contents string
}
