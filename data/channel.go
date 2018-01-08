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
