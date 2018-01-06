package widgets

import (
	"github.com/cceckman/discoirc/data"
	"github.com/marcusolsson/tui-go"
)

// ConnState is a Widget that displays a single character indicating the
// current connection state.
type ConnState struct {
	*tui.Label
}

// NewConnState returns a ConnState widget.
func NewConnState() *ConnState {
	return &ConnState{
		Label: tui.NewLabel("?"),
	}
}

// Set sets the current connection state.
func (c *ConnState) Set(state data.ConnectionState) {
	connStrings := map[data.ConnectionState]string{
		data.Disconnected: "∅",
		data.Connecting:   "…",
		data.Connected:    "✓",
	}
	if text, ok := connStrings[state]; ok {
		c.Label.SetText(text)
	} else {
		c.Label.SetText("?")
	}
}
