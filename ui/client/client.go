// Package client contains the MVC for the client's overall state.
package client

import (

	"github.com/marcusolsson/tui-go"
)

// ClientView is a top-level view of the client state.
type ClientView interface {
	tui.Widget

	// GetNetwork gets the NetworkView of the network with the given name.
	// It creates the view if one does not already exist.
	GetNetwork(string) NetworkView
	// RemoveNetwork removes the NetworkView with the given name.
	RemoveNetwork(string)
}

// NetworkView is a view of a particular network's state.
type NetworkView interface {
	Name() string

	SetNick(string)
	SetConnection(string)

	// GetChannel gets the ChannelView of the channel with the given name.
	// It creates a view if one does not already exist.
	/*
	   	GetChannel(string) ChannelView
	   	 RemoveChannel(string)
	   }

	   // ChannelView is the view of a particular channel-in-network's state.
	   type ChannelView interface {
	   	Name() string

	   	SetMode(string)
	   	SetUnreadCount(string)
	   	SetMembers(string)
	*/
}
