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
	Scope

	Presence    Presence
	ChannelMode string

	Topic   string
	Members int

	Unread      int
	LastMessage Seq
}


// ChannelStateEvent is an Event indicating a change in a channel's state.
type ChannelStateEvent struct {
	ChannelState

	// Line is the IRC line indicating this change.
	Line string

	seq Seq
}

var _ Event = &ChannelStateEvent{}

// Scope returns the channel's Scope.
func (c *ChannelStateEvent) Scope() Scope { return c.ChannelState.Scope }
// String implements fmt.Stringer
func (c *ChannelStateEvent) String() string { return c.Line }
// Seq returns the event's sequence ID within the scope
func (c *ChannelStateEvent) Seq() Seq { return c.seq }
