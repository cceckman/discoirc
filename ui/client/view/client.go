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

func (c *Client) FocusDefault() tui.Widget {
	return c
}

func (c *Client) FocusNext(w tui.Widget) tui.Widget {
	switch w := w.(type) {
	case *Client:
		// If we're the current selection, and we have networks to select,
		// select the first network.
		if w == c && len(c.networks) > 0 {
			return c.networks[0]
		}
	case *Network:
		n := w
		// If the network knows what's next, use it.
		if next := n.focusNext(w); next != nil {
			return next
		}
		// Otherwise, select the next network in the chain.
		if next := c.nextNetwork(w); next != nil {
			return next
		}
	case *Channel:
		ch := w
		// If the network knows what the next thing is, use it.
		if next := ch.network.focusNext(w); next != nil {
			return next
		}
		// Otherwise, move on to the network after the channel's.
		if next := c.nextNetwork(ch.network); next != nil {
			return next
		}
	}
	// Final default: the Client itself.
	return c.FocusDefault()
}

func (c *Client) FocusPrev(w tui.Widget) tui.Widget {
	switch w := w.(type) {
	case *Client:
		// Select the last network.
		if len(c.networks) > 0 {
			next := c.networks[len(c.networks)-1].focusPrev(w)
			if next != nil {
				return next
			}
		}
	case *Network:
		// The network won't know what's previous to it;
		// it's either another channel, or another network.
		for i, n := range c.networks {
			if n == w  && i-1 >= 0{
				next := c.networks[i-1].focusPrev(w)
				if next != nil {
					return next
				}
			}
		}
		// No match? Wrap around to last network.
		if len(c.networks) > 0 {
			next := c.networks[len(c.networks)-1].focusPrev(w)
			if next != nil {
				return next
			}
		}
	case *Channel:
		ch := w
		// The network should know what to do here - either another
		// channel, or the network itself.
		if next := ch.network.focusPrev(w); next != nil {
			return next
		}
	}
	// Final default: the Client itself.
	return c.FocusDefault()

}

// nextNetwork picks the next network in the list, wrapping around to the top.
func (c *Client) nextNetwork(w *Network) *Network {
	for i, n := range c.networks {
		if n == w {
			if i+1 < len(c.networks) {
				return c.networks[i+1]
			}
		}
	}
	// Roll around to top network.
	if len(c.networks) > 0 {
		return c.networks[0]
	}
	return nil
}

type netByName []*Network

func (n netByName) Len() int           { return len(n) }
func (n netByName) Less(i, j int) bool { return n[i].name < n[j].name }
func (n netByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
