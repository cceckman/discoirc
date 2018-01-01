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
	return &Client{
		Nets:  make(map[string]data.NetworkState),
		Chans: make(map[ChannelIdent]data.ChannelState),
	}
}

type Client struct {
	Nets  map[string]data.NetworkState
	Chans map[ChannelIdent]data.ChannelState
}

func (c *Client) UpdateNetwork(d data.NetworkState) {
	c.Nets[d.Network] = d
}

func (c *Client) UpdateChannel(d data.ChannelState) {
	c.Chans[ChannelIdent{
		Network: d.Network,
		Channel: d.Channel,
	}] = d
}
