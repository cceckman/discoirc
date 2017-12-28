// Package client contains the MVC for the client's overall state.
package client

import (
	"github.com/marcusolsson/tui-go"
)

// View is a top-level view of the client state.
type View interface {
	tui.Widget
	tui.FocusChain

	// GetNetwork gets the NetworkView of the network with the given name.
	// It creates the view if one does not already exist.
	GetNetwork(string) NetworkView
	// RemoveNetwork removes the NetworkView with the given name.
	RemoveNetwork(string)

	Attach(UIController)
}

// NetworkView is a view of a particular network's state.
type NetworkView interface {
	tui.Widget
	Name() string

	SetNick(string)
	SetConnection(string)

	// GetChannel gets the ChannelView of the channel with the given name.
	// It creates a view if one does not already exist.
	GetChannel(string) ChannelView
	RemoveChannel(string)
}

// ChannelView is the view of a particular channel-in-network's state.
type ChannelView interface {
	tui.Widget
	Name() string

	SetMode(string)
	SetUnread(int)
	SetMembers(int)
}

// UIController handles UI events from a client View.
type UIController interface {
	ActivateChannel(network, channel string)
	Quit()
}
