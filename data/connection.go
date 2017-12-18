package data

// ConnectionState represents the status of a user's connection to an IRC network.
type ConnectionState int

const (
	Disconnected ConnectionState = iota
	Connecting
	Connected
)


// Connection represents a user's relation with an IRC network.
type Connection struct {
	Nick string
	State ConnectionState
}
