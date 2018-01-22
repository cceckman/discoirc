package testhelper

import (
	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/data"
)

var _ backend.StateReceiver = &Channel{}

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

// Channel is a mock Channel view. It supports the FilteredSubscriber interface.
type Channel struct {
	*Client
	Scope data.Scope
}

// Filter specifies the network and target this Channel is watching for.
func (c *Channel) Filter() data.Filter {
	return data.Filter{
		Scope:     c.Scope,
		MatchNet:  true,
		MatchName: true,
	}
}
