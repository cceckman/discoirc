package view

import (
	"sort"

	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/ui/client"
)

var _ client.NetworkView = &Network{}

// NewNetwork gives a new Network view.
func NewNetwork(name string) *Network {
	r := &Network{
		name:            name,
		indicatorWidget: newIndicator(),
		nameWidget:      tui.NewLabel(name),
		nickWidget:      tui.NewLabel(""),
		connWidget:      tui.NewLabel(""),
		chanWidget:      tui.NewVBox(),
	}

	r.Box = tui.NewVBox(
		tui.NewHBox(
			r.indicatorWidget,
			r.nameWidget,
			tui.NewLabel(": "),
			r.connWidget,
			tui.NewLabel(" "),
			tui.NewSpacer(),
			r.nickWidget,
		),
		r.chanWidget,
	)

	return r
}

// Network is the implementation of a NetworkView.
type Network struct {
	name string
	*tui.Box

	channels []*Channel

	indicatorWidget *indicator
	nameWidget      *tui.Label
	nickWidget      *tui.Label
	connWidget      *tui.Label
	chanWidget      *tui.Box
}

func (n *Network) SetFocused(focus bool) {
	n.Box.SetFocused(true)
	if focus {
		n.indicatorWidget.SetFill('>')
	} else {
		n.indicatorWidget.SetFill(' ')
	}
}

func (n *Network) Name() string {
	return n.name
}

func (n *Network) SetNick(s string) {
	n.nickWidget.SetText(s)
}

func (n *Network) SetConnection(s string) {
	n.connWidget.SetText(s)
}

func (n *Network) GetChannel(name string) client.ChannelView {
	for _, v := range n.channels {
		if v.name == name {
			return v
		}
	}
	// Add new network; insert into widget
	c := NewChannel(name)
	n.channels = append(n.channels, c)
	sort.Sort(chanByName(n.channels))
	for i, v := range n.channels {
		if v.name == name {
			n.chanWidget.Insert(i, v)
			return v
		}
	}
	return nil
}

func (n *Network) RemoveChannel(name string) {
	for i, v := range n.channels {
		if v.name == name {
			n.channels = append(n.channels[0:i], n.channels[i+1:]...)
			n.chanWidget.Remove(i)
			return
		}
	}
	return
}

type chanByName []*Channel

func (n chanByName) Len() int           { return len(n) }
func (n chanByName) Less(i, j int) bool { return n[i].name < n[j].name }
func (n chanByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
