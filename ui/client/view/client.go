package view

import (
	"sort"

	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/ui/client"
)

var _ client.ClientView = &Client{}

func New() *Client {
	return &Client{
		Box: tui.NewVBox(tui.NewSpacer()),
	}
}

type Client struct {
	*tui.Box
	networks []*Network
}

func (c *Client) GetNetwork(name string) client.NetworkView {
	for _, v := range c.networks {
		if v.name == name {
			return v
		}
	}
	// Add new network; insert into widget
	n := NewNetwork(name)
	c.networks = append(c.networks, n)
	sort.Sort(netByName(c.networks))
	for i, v := range c.networks {
		if v.name == name {
			c.Box.Insert(i, v)
			return v
		}
	}
	return nil
}

func (c *Client) RemoveNetwork(name string) {
	for i, v := range c.networks {
		if v.name == name {
			c.networks = append(c.networks[0:i], c.networks[i+1:]...)
			c.Box.Remove(i)
			return
		}
	}
	return
}

type netByName []*Network

func (n netByName) Len() int           { return len(n) }
func (n netByName) Less(i, j int) bool { return n[i].name < n[j].name }
func (n netByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
