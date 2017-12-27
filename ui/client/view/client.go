package view

import (
	"sort"

	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/ui/client"
)

var _ client.View = &Client{}

func New() *Client {
	c := &Client{
		networksBox: tui.NewVBox(tui.NewSpacer()),
	}

	c.Widget = c.networksBox

	return c
}

type Client struct {
	tui.Widget
	networksBox *tui.Box

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
			c.networksBox.Insert(i, v)
			return v
		}
	}
	return nil
}

func (c *Client) RemoveNetwork(name string) {
	for i, v := range c.networks {
		if v.name == name {
			c.networks = append(c.networks[0:i], c.networks[i+1:]...)
			c.networksBox.Remove(i)
			return
		}
	}
	return
}

type netByName []*Network

func (n netByName) Len() int           { return len(n) }
func (n netByName) Less(i, j int) bool { return n[i].name < n[j].name }
func (n netByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
