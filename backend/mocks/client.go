// Package mocks provides mocks to use for testing backends.
package mocks

import (
	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/data"
)

var _ backend.StateReceiver = &Client{}

type ChannelIdent struct {
	Network, Channel string
}

func NewClient() *Client {
	c := &Client{
		Nets:     make(map[string]data.NetworkState),
		Chans:    make(map[ChannelIdent]data.ChannelState),
		Contents: make(map[ChannelIdent][]data.Event),
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
	Nets     map[string]data.NetworkState
	Chans    map[ChannelIdent]data.ChannelState
	Contents map[ChannelIdent][]data.Event

	await chan func()

	// If Archives is nonempty, UpdateChannel will cause a query to
	// EventsBefore to be made on each update.
	Archive backend.EventsArchive
}

func (c *Client) UpdateNetwork(d data.NetworkState) {
	c.Join(func() {
		c.Nets[d.Network] = d
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

func (c *Client) UpdateChannel(d data.ChannelState) {
	c.Join(func() {
		cid := ChannelIdent{
			Network: d.Network,
			Channel: d.Channel,
		}
		c.Chans[cid] = d

		if c.Archive != nil {
			c.Contents[cid] = c.Archive.EventsBefore(
				d.Network, d.Channel,
				100, d.LastMessage.EventID)
		}
	})
}
