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
	State ConnectionState

	Nick     string
	UserMode string
}

// NetworkStateEvent is an Event indicating a change in the network's state.
type NetworkStateEvent struct {
	NetworkState
	EventID

	// Line is the IRC line indicating this change.
	Line string
}

var _ Event = &NetworkStateEvent{}

// ID returns the scope & sequence of this Event.
func (e *NetworkStateEvent) ID() *EventID {
	return &e.EventID
}

// String implments fmt.Stringer.
func (e *NetworkStateEvent) String() string { return e.Line }
