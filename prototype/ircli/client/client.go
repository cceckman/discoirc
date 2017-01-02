// Package client provides a multi-server IRC client.
package client

import (
	"sync"
	irc "github.com/fluffle/goirc/client"
)

// C is a client connected to (potentially several) IRC servers and rooms.
type C interface {
	LoadConfigs() map[string]*irc.Config
	Connect() <-chan error
	ConnectedNetworks() []string
	// Disconnect() error

	// Set the current target to channel or nick 'target' on the given network.
	// SetTarget(network, target error) error
	// Send(msg string) error
}

type client struct{
	connections map[string]*irc.Conn
	sync.Mutex
}

func NewClient() C {
	return &client{}
}
