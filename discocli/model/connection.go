package model

// Connection represents a relationship with an IRC network (or server).
// It may be connected to one or more channels.
type Connection interface {
	Channel(name string) Channel
	Channels() []string

	// TODO add server info
}
