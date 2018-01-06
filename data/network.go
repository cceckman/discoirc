package data

// NetworkState represents the state of a user's relation to a network.
type NetworkState struct {
	Network, Nick string
	State         ConnectionState
	// TODO: make UserMode and ChannelMode their own types
	UserMode string
}
