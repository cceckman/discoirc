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

// NetworkStateEvent is an Event indicating a change in the network's state.
type NetworkStateEvent struct {
	NetworkState

	// Line is the IRC line indicating this change.
	Line string
	seq  Seq
}

var _ Event = &NetworkStateEvent{}

// Scope returns a Scope for this event.
func (l *NetworkStateEvent) Scope() Scope { return l.NetworkState.Scope }

// String implments fmt.Stringer.
func (l *NetworkStateEvent) String() string { return l.Line }

// Seq returns the sequence ID of this event within its scope
func (l *NetworkStateEvent) Seq() Seq { return l.seq }
