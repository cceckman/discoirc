package data

// Presence is a user's presence, or lack thereof, in a channel.
type Presence int

const (
	// NotPresent indicates the user is not in the channel.
	NotPresent Presence = iota
	// Joined indicates the user is in the channel.
	Joined
)

// ChannelState summarizes the state of a channel.
type ChannelState struct {
	Presence Presence
	Mode     string

	Topic   string
	Members int

	Unread      int
	LastMessage Seq
}

// ChannelStateEvent is an Event indicating a change in a channel's state.
type ChannelStateEvent struct {
	ChannelState
	EventID

	// Line is the IRC line indicating this change.
	Line string
}

var _ Event = &ChannelStateEvent{}

// ID returns the scope & sequence of this Event.
func (e *ChannelStateEvent) ID() *EventID {
	return &e.EventID
}

// String implments fmt.Stringer.
func (e *ChannelStateEvent) String() string { return e.Line }
