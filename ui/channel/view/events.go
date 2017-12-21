package view

import (
	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/cceckman/discoirc/ui/widgets"
	"github.com/marcusolsson/tui-go"
)

type EventsView struct {
	*widgets.TailBox

	Renderer channel.EventRenderer
}

func (v *EventsView) SetEvents(evs []data.Event) {
	w := make([]tui.Widget, len(evs))
	for i, e := range evs {
		w[i] = v.Renderer(e)
	}
	v.SetContents(w...)
}