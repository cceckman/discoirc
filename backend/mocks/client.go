package mocks

import (
	"sync"

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
	go func(){
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

	wg sync.WaitGroup
	await chan func()
}

func (c *Client) UpdateNetwork(d data.NetworkState) {
	c.wg.Add(1)
	c.await <- func() {
		c.Nets[d.Network] = d
		c.wg.Done()
	}
}

// Close waits for all outstanding updates to complete, then collects this Client.
func (c *Client) Close() {
	c.wg.Wait()
	c.Join(func(){
		close(c.await)
	})
}

// Join runs the closure in the same thread as updates, and returns once it completes.
func (c *Client) Join(f func ()) {
	blk := make(chan struct{})
	c.wg.Add(1)
	c.await <- func() {
		f()
		close(blk)
		c.wg.Done()
	}
	<-blk
}

func (c *Client) UpdateChannel(d data.ChannelState) {
	c.wg.Add(1)
	c.await <- func() {
		c.Chans[ChannelIdent{
			Network: d.Network,
			Channel: d.Channel,
		}] = d
		c.wg.Done()
	}
}
