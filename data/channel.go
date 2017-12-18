package data


// Presence is a user's presence, or lack thereof, in a channel.
type Presence int
const (
	NotPresent Presence = iota
	Joined
)

// Channel represents the state of a user's presence in a channel.
type Channel struct {
	Connection Connection

	Topic string
	Presence Presence
	Mode string
}
