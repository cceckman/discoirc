package view

import (
	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
	"github.com/cceckman/discoirc/ui/widgets"
	"github.com/marcusolsson/tui-go"
	"image"
)

var _ channel.View = &V{}

func DefaultRenderer(e data.Event) tui.Widget {
	r := tui.NewLabel(e.Contents)
	r.SetWordWrap(true)
	r.SetSizePolicy(tui.Expanding, tui.Minimum)
	return r
}

// V implements channel.View using tui-go.
type V struct {
	// root element
	*tui.Box

	// Top-level display elements:
	ui tui.UI

	// Second-level elements
	topic  *tui.Label
	events *EventsView
	// status bar
	connection *tui.Label
	presence   *tui.Label
	mode       *tui.Label
	// input bar
	nick  *tui.Label
	input *tui.Entry

	controller channel.Controller
}

func (v *V) handleInput(entry *tui.Entry) {
	if v.controller != nil {
		v.controller.Input(entry.Text())
		entry.SetText("")
	}
}

func (v *V) SetRenderer(e channel.EventRenderer) {
	v.events.Renderer = e
}
func (v *V) SetTopic(t string) {
	v.topic.SetText(t)
}
func (v *V) SetConnection(s string) {
	v.connection.SetText(s)
}
func (v *V) SetPresence(s string) {
	v.presence.SetText(s)
}
func (v *V) SetMode(s string) {
	v.mode.SetText(s)
}
func (v *V) SetNick(s string) {
	v.nick.SetText(s)
}
func (v *V) SetEvents(events []data.Event) {
	v.events.SetEvents(events)
}

func (v *V) Attach(c channel.Controller) {
	v.controller = c
}

func (v *V) Resize(size image.Point) {
	eventsSize := v.events.Size()
	v.Box.Resize(size)
	if eventsSize.Y > v.events.Size().Y && v.controller != nil {
		// events box got bigger. Request an update.
		v.controller.Resize(v.events.Size().Y)
	}
}

type reversedBox struct {
	*tui.Box
}

func (rb *reversedBox) Draw(p *tui.Painter) {
	p.WithStyle("reversed", func(p *tui.Painter) {
		rb.Box.Draw(p)
	})
}

// New returns a new channel.View and assigns it to the current UI.
func New(ui tui.UI) channel.View {
	// construct V
	v := &V{
		ui: ui,

		topic: tui.NewLabel(""),
		events: &EventsView{
			TailBox:  widgets.NewTailBox(),
			Renderer: DefaultRenderer,
		},
		connection: tui.NewLabel(""),
		presence:   tui.NewLabel(""),
		mode:       tui.NewLabel(""),
		nick:       tui.NewLabel(""),
		input:      tui.NewEntry(),
	}
	v.topic.SetSizePolicy(tui.Expanding, tui.Minimum)
	v.events.SetSizePolicy(tui.Expanding, tui.Expanding)
	v.connection.SetSizePolicy(tui.Minimum, tui.Minimum)
	v.presence.SetSizePolicy(tui.Minimum, tui.Minimum)
	v.mode.SetSizePolicy(tui.Minimum, tui.Minimum)
	v.nick.SetSizePolicy(tui.Minimum, tui.Minimum)
	v.input.SetSizePolicy(tui.Expanding, tui.Minimum)

	v.input.OnSubmit(v.handleInput)

	rspacer := tui.NewLabel(" ")
	rspacer.SetSizePolicy(tui.Expanding, tui.Preferred)

	inputBar := tui.NewHBox(v.nick, v.input)

	v.Box = tui.NewVBox(
		&reversedBox{
			Box: tui.NewHBox(v.topic),
		},
		v.events,
		&reversedBox{
			Box: tui.NewHBox(
				v.connection,
				tui.NewLabel(" "),
				v.presence,
				tui.NewLabel(" "),
				v.mode,
				rspacer,
			),
		},
		inputBar,
	)
	return v
}
