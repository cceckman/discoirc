package mocks

import (
	"github.com/cceckman/discoirc/backend"
)

var _ backend.FilteredStateReceiver = &Channel{}

// NewChannel returns a new mock Channel.
func NewChannel(network, name string) *Channel {
	return &Channel{
		Client: NewClient(),

		Network: network,
		Channel: name,
	}
}

// Channel is a mock Channel view. It supports the FilteredSubscriber interface.
type Channel struct {
	*Client

	Network, Channel string
}

// Filter specifies the network and target this Channel is watching for.
func (c *Channel) Filter() (string, string) {
	return c.Network, c.Channel
}
