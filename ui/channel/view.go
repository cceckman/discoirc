package channel

import (
	"image"

	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/widgets"
	"github.com/marcusolsson/tui-go"
)

var _ View = &view{}

// DefaultRenderer is the default way to render Widgets.
func DefaultRenderer(e data.Event) tui.Widget {
	r := tui.NewLabel(e.Contents)
	r.SetWordWrap(true)
	r.SetSizePolicy(tui.Expanding, tui.Minimum)
	return r
}

// view implements View as a tui Widget.
type view struct {
	// root element
	*tui.Box

	// Second-level elements
	topic  *tui.Label
	events *EventsView
	// status bar
	connection *tui.Label
	name       *tui.Label
	mode       *tui.Label
	// input bar
	nick  *tui.Label
	input *tui.Entry

	controller Controller
}

func (v *view) OnKeyEvent(ev tui.KeyEvent) {
	if ev.Key == tui.KeyCtrlC && v.controller != nil {
		v.controller.Quit()
	}
	v.Box.OnKeyEvent(ev)
}

func (v *view) handleInput(entry *tui.Entry) {
	if v.controller != nil {
		v.controller.Input(entry.Text())
		entry.SetText("")
	}
}

func (v *view) SetRenderer(e EventRenderer) {
	v.events.Renderer = e
}
func (v *view) SetTopic(t string) {
	v.topic.SetText(t)
}
func (v *view) SetConnection(s string) {
	v.connection.SetText(s)
}
func (v *view) SetName(s string) {
	v.name.SetText(s)
}
func (v *view) SetMode(s string) {
	v.mode.SetText(s)
}
func (v *view) SetNick(s string) {
	v.nick.SetText(s)
}
func (v *view) SetEvents(events []data.Event) {
	v.events.SetEvents(events)
}

func (v *view) Attach(c Controller) {
	v.controller = c
	// Set initial size
	v.controller.Resize(v.events.Size().Y)
}

func (v *view) Resize(size image.Point) {
	oldSize := v.events.Size()
	v.Box.Resize(size)
	if v.events.Size().Y > oldSize.Y && v.controller != nil {
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

// New returns a new View.
func NewView() View {
	// construct V
	v := &view{
		topic: tui.NewLabel(""),
		events: &EventsView{
			TailBox:  widgets.NewTailBox(),
			Renderer: DefaultRenderer,
		},
		connection: tui.NewLabel(""),
		name:       tui.NewLabel(""),
		mode:       tui.NewLabel(""),
		nick:       tui.NewLabel(""),
		input:      tui.NewEntry(),
	}
	v.topic.SetSizePolicy(tui.Expanding, tui.Minimum)
	v.events.SetSizePolicy(tui.Expanding, tui.Expanding)
	v.connection.SetSizePolicy(tui.Minimum, tui.Minimum)
	v.name.SetSizePolicy(tui.Minimum, tui.Minimum)
	v.mode.SetSizePolicy(tui.Minimum, tui.Minimum)
	v.nick.SetSizePolicy(tui.Minimum, tui.Minimum)
	v.input.SetSizePolicy(tui.Expanding, tui.Minimum)

	v.input.OnSubmit(v.handleInput)
	v.input.SetFocused(true)

	rspacer := tui.NewSpacer()

	inputBar := tui.NewHBox(
		tui.NewLabel("<"),
		v.nick,
		tui.NewLabel("> "),
		v.input,
	)

	v.Box = tui.NewVBox(
		&reversedBox{
			Box: tui.NewHBox(v.topic),
		},
		v.events,
		&reversedBox{
			Box: tui.NewHBox(
				v.connection,
				tui.NewLabel(" "),
				v.name,
				tui.NewLabel(" "),
				v.mode,
				rspacer,
			),
		},
		inputBar,
	)
	return v
}

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
