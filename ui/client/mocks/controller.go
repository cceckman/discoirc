package mocks

import (
	"github.com/cceckman/discoirc/ui/client"
)

var _ client.UIController = &Controller{}

// Activation is a record of an ActivateChannel call.
type Activation struct {
	Network, Channel string
}

// Controller implements the UIController interface, for testing.
type Controller struct {
	Activations []Activation
}

func (c *Controller) ActivateChannel(network, channel string) {
	c.Activations = append(c.Activations, Activation{
		Network: network,
		Channel: channel,
	})
}
