package testhelper

import (
	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/data"
)

var _ backend.Receiver = &Channel{}

// Channel is a mock Channel view. It supports the FilteredSubscriber interface.
type Channel struct {
	*Client
	Scope data.Scope
}

// NewChannel returns a new mock Channel.
func NewChannel(network, name string) *Channel {
	return &Channel{
		Client: NewClient(),
		Scope: data.Scope{
			Net:  network,
			Name: name,
		},
	}
}

// Filter specifies the network and target this Channel is watching for.
func (c *Channel) Filter() data.Filter {
	return data.Filter{
		Scope:     c.Scope,
		MatchNet:  true,
		MatchName: true,
	}
}
