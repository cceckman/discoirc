package girc

import (
	"fmt"

	"github.com/lrstanley/girc"

	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/data"
)

// Network encapsulates event logging and a connection against an IRC network.
type Network struct {
	State data.NetworkState

	// Config is the configuration that should be used to connect to this
	// network.
	Config girc.Config

	// Client is an active or available IRC connection wrapper.
	*girc.Client

	log map[data.Scope]*Log
}

func NewNetwork(name string, c *girc.Config) *Network {
	netScope := data.Scope{Net: name}
	n := &Network{
		State: data.NetworkState{
			Scope: netScope,
		},
		log: map[data.Scope]*Log{
			netScope: &Log{Scope: netScope},
		},
	}

	n.Reconfigure(c)

	return n
}

// Reconfigure closes the current client and creates a new one with the updated
// configuration.
func (n *Network) Reconfigure(c *girc.Config) {
	if c != nil {
		n.Config = *c
	}

	if n.Client != nil {
		go n.Client.Close()
		// TODO: Capture the error from Connect().
	}
	n.Client = girc.New(n.Config)

	// Set up handlers for the new client
	n.Client.Handlers.Add(girc.CONNECTED, n.OnConnect)
	n.Client.Handlers.Add(girc.DISCONNECTED, n.OnDisconnect)
}

// Connect attempts to connect, and blocks until the connection is closed.
func (n *Network) Connect() error {
	n.Reconfigure(nil)
	n.State.State = data.Connecting
	return n.Client.Connect()
}

func (n *Network) OnConnect(c *girc.Client, ev girc.Event) {
	n.State.State = data.Connected
	n.State.Nick = c.GetNick()

	e := data.Event{
		EventContents: &data.NetworkStateEvent{
			NetworkState: n.State,
			Line:         ev.String(),
		},
	}
	for _, log := range n.log {
		log.Append(e)
	}
}
func (n *Network) OnDisconnect(c *girc.Client, ev girc.Event) {
	n.State.State = data.Disconnected
	n.State.Nick = c.GetNick()

	e := data.Event{
		EventContents: &data.NetworkStateEvent{
			NetworkState: n.State,
			Line:         ev.String(),
		},
	}
	for _, log := range n.log {
		log.Append(e)
	}
}

var _ backend.EventsArchive = &Network{}

func (n *Network) EventsBefore(_ data.Scope, _ int, _ data.Seq) data.EventList {
	panic(fmt.Errorf("not implemented"))
}
