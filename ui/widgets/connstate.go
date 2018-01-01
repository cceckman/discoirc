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

func NewConnState() *ConnState {
	return &ConnState{
		Label: tui.NewLabel("?"),
	}
}

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
