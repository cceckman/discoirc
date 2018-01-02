package data

type NetworkState struct {
	Network, Nick string
	State         ConnectionState
	// TODO: make UserMode and ChannelMode their own types
	UserMode string
}
