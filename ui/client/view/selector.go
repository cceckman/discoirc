package view

import (
	"github.com/cceckman/discoirc/ui/widgets"
	"github.com/marcusolsson/tui-go"
	"image"
)

/*
Sketch: How should selection work? Some options:

1 External list, a la focus.go. Reflect the flattened network+channel structure
	in an external structure; keep it updated as elements are added / removed.
	I don't like this because it involves keeping the two lists synced,
	i.e. storing in two places.
2 Implementing FocusChain within Client, Network, and/or Channel.
	This is nice because, well, that's what it's meant for; but it's not a very
	natural fit for Widget vs. other stuff
	But maybe it works, implementing in Client;
	case on *Network / *Channel / *Client, select the appropriate next.
3 Implementing a custom API. May be a more natural API, but less integrated with
	tui.

Let's go with 2 to begin with.
*/


func newIndicator() *indicator {
	r := &indicator{}
	r.SetFill(' ')
	return r
}

type indicator struct {
	widgets.Filler
}

func (f *indicator) SizeHint() image.Point {
	return image.Pt(1, 0)
}

func (f *indicator) SizePolicy() (tui.SizePolicy, tui.SizePolicy) {
	return tui.Maximum, tui.Preferred
}
