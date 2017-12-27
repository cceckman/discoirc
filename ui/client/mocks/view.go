package mocks

import (
	"github.com/marcusolsson/tui-go"
	"github.com/cceckman/discoirc/ui/client"
)

type View struct {
	tui.Widget
	tui.FocusChain
	// above intentionally nil

	Networks   map[string]*Network
	Controller client.UIController
}

func (v *View) GetNetwork(name string) client.NetworkView {
	if _, ok := v.Networks[name]; !ok {
		v.Networks[name] = &Network{
			name:     name,
			Channels: make(map[string]*Channel),
		}
	}
	return v.Networks[name]
}
func (v *View) RemoveNetwork(name string) {
	delete(v.Networks, name)
}

func (v *View) Attach(ctl client.UIController) {
	v.Controller = ctl
}

var _ client.View = &View{}

// Network is a mock client.NetworkView for tests.
type Network struct {
	tui.Widget // intentionally nil

	name       string
	Nick       string
	Connection string
	Channels   map[string]*Channel
}

func (n *Network) Name() string           { return n.name }
func (n *Network) SetNick(s string)       { n.Nick = s }
func (n *Network) SetConnection(s string) { n.Connection = s }
func (n *Network) GetChannel(name string) client.ChannelView {
	if _, ok := n.Channels[name]; !ok {
		n.Channels[name] = &Channel{
			name: name,
		}
	}

	return n.Channels[name]
}
func (n *Network) RemoveChannel(name string) {
	delete(n.Channels, name)
}

var _ client.NetworkView = &Network{}

// Channel is a mock client.ChannelView for tests.
type Channel struct {
	tui.Widget // intentionally nil

	name    string
	Mode    string
	Unread  int
	Members int
}

func (c *Channel) Name() string     { return c.name }
func (c *Channel) SetMode(m string) { c.Mode = m }
func (c *Channel) SetUnread(u int)  { c.Unread = u }
func (c *Channel) SetMembers(n int) { c.Members = n }

var _ client.ChannelView = &Channel{}
