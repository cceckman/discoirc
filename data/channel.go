package data

// Presence is a user's presence, or lack thereof, in a channel.
type Presence int

const (
	NotPresent Presence = iota
	Joined
)

// Channel represents the state of a user's presence in a channel.
type Channel struct {
	Name       string
	Connection Connection

	Topic    string
	Presence Presence
	Mode     string
}

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
