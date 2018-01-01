// Package client contains the MVC for the client's overall state.
package client

import (
	"github.com/marcusolsson/tui-go"
)

// View is a top-level view of the client state.
type View interface {
	tui.Widget
	tui.FocusChain
}

// UIController handles UI events from a client View.
// All its methods should be called within an Update closure.
type UIController interface {
	Update(func())
	ActivateChannel(network, channel string)
	SetWidget(tui.Widget)
	Quit()
}
