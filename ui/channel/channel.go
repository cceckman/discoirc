// Package channel provides the channel contents UI.
package channel

import (
	"github.com/cceckman/discoirc/data"
	"github.com/marcusolsson/tui-go"
)

// EventRenderer is a function that converts a DiscoIRC event
// (e.g. message) into an tui.Widget suitable for display.
type EventRenderer func(data.Event) tui.Widget

