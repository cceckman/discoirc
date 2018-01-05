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
		Nets:  make(map[string]data.NetworkState),
		Chans: make(map[ChannelIdent]data.ChannelState),
		await: make(chan func()),
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
	Nets  map[string]data.NetworkState
	Chans map[ChannelIdent]data.ChannelState

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

func (c *Client) Close() {
	c.Join(func() {
		close(c.await)
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
		c.Chans[ChannelIdent{
			Network: d.Network,
			Channel: d.Channel,
		}] = d

		if c.Archive != nil {
			_ = c.Archive.EventsBefore(
				d.Network, d.Channel,
				1, d.LastMessage.EventID)
		}
	})
}
