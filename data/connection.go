package data

// ConnectionState represents the status of a user's connection to an IRC network.
type ConnectionState int

const (
	// Disconnected indicates the client is not connected to the network.
	Disconnected ConnectionState = iota
	// Connecting indicates the client is attempting to connect to the network.
	Connecting
	// Connected indicates the client has connected to the network.
	Connected
)

// Connection represents a user's relation with an IRC network.
type Connection struct {
	Network string
	Nick    string
	State   ConnectionState
}
