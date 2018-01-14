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

// NetworkState represents the state of a user's relation to a network.
type NetworkState struct {
	Scope

	Nick     string
	State    ConnectionState
	UserMode string
}

// NetworkStateEvent can be part of an Event, indicating a change in the network's state.
type NetworkStateEvent struct {
	NetworkState

	// Line is the IRC line indicating this change.
	Line string
}

var _ EventContents = &NetworkStateEvent{}
func (l *NetworkStateEvent) String() string { return l.Line }

