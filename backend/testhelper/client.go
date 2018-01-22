// Package mocks provides mocks to use for testing backends.
package testhelper

import (
	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/data"
)

var _ backend.StateReceiver = &Client{}

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

// UpdateNetwork receives the new state of the network.
func (c *Client) UpdateNetwork(d data.NetworkState) {
	c.Join(func() {
		c.Nets[d.Scope] = d
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

// UpdateChannel receives the new state of the channel.
func (c *Client) UpdateChannel(d data.ChannelState) {
	c.Join(func() {
		c.Chans[d.Scope] = d

		if c.Archive != nil {
			c.Contents[d.Scope] = c.Archive.EventsBefore(
				d.Scope, 100, d.LastMessage)
		}
	})
}
