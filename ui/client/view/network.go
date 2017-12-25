package view

import (
	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/ui/client"
)

var _ client.NetworkView = &Network{}

// NewNetwork gives a new Network view.
func NewNetwork(name string) *Network {
	r := &Network{
		name:       name,
		nameWidget: tui.NewLabel(name),
		nickWidget: tui.NewLabel(""),
		connWidget: tui.NewLabel(""),
	}

	r.Widget = tui.NewHBox(
		r.nameWidget,
		tui.NewLabel(": "),
		r.connWidget,
		tui.NewLabel(" "),
		tui.NewSpacer(),
		r.nickWidget,
	)

	return r
}

// Network is the implementation of a NetworkView.
type Network struct {
	name string
	tui.Widget

	nameWidget *tui.Label

	nickWidget *tui.Label
	connWidget *tui.Label
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
