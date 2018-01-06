package channel

import (
	"image"

	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/widgets"
)

// EventRenderer is a function that converts a DiscoIRC event
// (e.g. message) into an tui.Widget suitable for display.
type EventRenderer func(data.Event) tui.Widget

// EventsProvider returns up to n events ending with 'last'.
type EventsProvider interface {
	EventsBefore(net, target string, n int, last data.EventID) []data.Event
}

// NewEventsWidget returns a new EventsWidget.
func NewEventsWidget(network, target string, in EventsProvider) *EventsWidget {
	return &EventsWidget{
		TailBox:  widgets.NewTailBox(),
		Renderer: DefaultRenderer,

		network: network,
		target:  target,
		source:  in,
	}
}

// EventsWidget displays the last data.Event objects it contains.
type EventsWidget struct {
	*widgets.TailBox

	source EventsProvider
	last   data.EventID

	network, target string

	Renderer EventRenderer
}

// SetLast sets the last Event displayed. If it's newer than the previous value,
// it may cause request a backfill of its contents.
func (v *EventsWidget) SetLast(new data.EventID) {
	if v.last != new && v.source != nil {
		v.last = new
		v.refreshContents()
	}
}

// refreshContents redraws the contents of the EventsWidget,
func (v *EventsWidget) refreshContents() {
	// TODO:
	// 1. Assume EventsSince may take a long time; handle it in a non-blocking way.
	// 2. Handle single-new-message more gracefully, i.e. without redrawing
	//    all of the widgets.
	events := v.source.EventsBefore(
		v.network, v.target,
		v.TailBox.Size().Y, v.last)

	w := make([]tui.Widget, len(events))
	for i, e := range events {
		w[i] = v.Renderer(e)
	}
	v.SetContents(w...)
}

// Resize handles resizing of the Widget. It may trigger a refresh of the
// Widget's contents.
func (v *EventsWidget) Resize(size image.Point) {
	oldSize := v.Size()
	v.TailBox.Resize(size)
	if v.TailBox.Size().Y > oldSize.Y && v.source != nil {
		v.refreshContents()
	}
}

// reversedBox is a Box that applies the "reversed" style to its contents.
type reversedBox struct {
	*tui.Box
}

func (rb *reversedBox) Draw(p *tui.Painter) {
	p.WithStyle("reversed", func(p *tui.Painter) {
		rb.Box.Draw(p)
	})
}
