// Package testhelper provides mocks to use for testing backends.
package testhelper

import (
	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/data"
)

var _ backend.Receiver = &Client{}

// NewClient returns a new Client.
func NewClient() *Client {
	c := &Client{
		Nets:     make(map[data.Scope]data.NetworkState),
		Chans:    make(map[data.Scope]data.ChannelState),
		Contents: make(map[data.Scope][]data.Event),
		await:    make(chan func()),
	}
	go func() {
		for f := range c.await {
			f()
		}
	}()
	return c
}

// Client is a mocked-up client.View.
// Its Update* methods run functions in a separate thread.
type Client struct {
	Nets     map[data.Scope]data.NetworkState
	Chans    map[data.Scope]data.ChannelState
	Contents map[data.Scope][]data.Event

	await chan func()

	// If Archives is nonempty, UpdateChannel will cause a query to
	// EventsBefore to be made on each update.
	Archive backend.EventsArchive
}

// Filter implements the Filter interface by accepting everything
func (c *Client) Filter() data.Filter {
	return data.Filter{}
}

// Receive receives the new state of a network or channel.
func (c *Client) Receive(e data.Event) {
	switch e := e.(type) {
	case *data.NetworkStateEvent:
		c.updateNetwork(e.ID().Scope, e.NetworkState)
	case *data.ChannelStateEvent:
		c.updateChannel(e.ID().Scope, e.ChannelState)
	}
}

func (c *Client) updateNetwork(scope data.Scope, state data.NetworkState) {
	c.Join(func() {
		c.Nets[scope] = state
	})
}

func (c *Client) updateChannel(scope data.Scope, state data.ChannelState) {
	c.Join(func() {
		c.Chans[scope] = state

		if c.Archive != nil {
			c.Contents[scope] = c.Archive.EventsBefore(
				scope, 100, state.LastMessage)
		}
	})
}

// Join runs the closure in the same thread as updates, and returns once it completes.
func (c *Client) Join(f func()) {
	blk := make(chan struct{})
	c.await <- func() {
		f()
		close(blk)
	}
	<-blk
}
