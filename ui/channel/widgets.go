package channel

import (
	"image"

	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/widgets"
)

// EventsProvider returns up to n events ending with 'last'.
type EventsProvider interface {
	EventsBefore(n int, last data.EventID) []data.Event
}


func NewEventsWidget(in EventsProvider) *EventsWidget {
	return &EventsWidget{
		TailBox: widgets.NewTailBox(),
		Renderer: DefaultRenderer,
		source: in,
	}
}

// EventsWidget displays the last data.Event objects it contains.
type EventsWidget struct {
	*widgets.TailBox

	source   EventsProvider
	last data.EventID

	Renderer EventRenderer
}

func (v *EventsWidget) SetLast(new data.EventID) {
	if v.last != new && v.source != nil {
		v.last = new
		v.refreshContents()
	}
}

// refreshContents redraws the contents of the EventsWidget,
func (v *EventsWidget) refreshContents(){
	// TODO:
	// 1. Assume EventsSince may take a long time; handle it in a non-blocking way.
	// 2. Handle single-new-message more gracefully, i.e. without redrawing
	//    all of the widgets.
	events := v.source.EventsBefore(v.TailBox.Size().Y, v.last)

	w := make([]tui.Widget, len(events))
	for i, e := range events {
		w[i] = v.Renderer(e)
	}
	v.SetContents(w...)
}

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
