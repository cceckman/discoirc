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
	Network, Channel string
	Presence         Presence
	// TODO: This representation is incorrect.
	// A user's mode is for the network, not to the channel.
	ChannelMode, UserMode string

	Topic   string
	Members int

	Unread      int
	LastMessage Event
}
