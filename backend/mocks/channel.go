package mocks

import (
	"github.com/cceckman/discoirc/backend"
)

var _ backend.FilteredStateReceiver = &Channel{}

func NewChannel(network, name string) *Channel {
	return &Channel{
		Client: NewClient(),

		Network: network,
		Channel: name,
	}
}

type Channel struct {
	*Client

	Network, Channel string
}

func (c *Channel) Filter() (string, string) {
	return c.Network, c.Channel
}
