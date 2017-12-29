package channel

import (
	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/widgets"
	"github.com/marcusolsson/tui-go"
)

// EventsView displays the last data.Event objects it contains.
type EventsView struct {
	*widgets.TailBox

	Renderer EventRenderer
}

func (v *EventsView) SetEvents(evs []data.Event) {
	w := make([]tui.Widget, len(evs))
	for i, e := range evs {
		w[i] = v.Renderer(e)
	}
	v.SetContents(w...)
}
