package channel

import (
	"strings"

	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/widgets"
	"github.com/marcusolsson/tui-go"
)

// DefaultRenderer is the default way to render Widgets.
func DefaultRenderer(e data.Event) tui.Widget {
	r := tui.NewLabel(e.String())
	r.SetWordWrap(true)
	r.SetSizePolicy(tui.Expanding, tui.Minimum)
	return r
}

// UIController provides an interface to a global control layer.
type UIController interface {
	Update(func())
	SetWidget(tui.Widget)
	Quit()

	ActivateClient()
}

// View implements the channel view.
type View struct {
	ui     UIController
	sender backend.Sender
	scope  data.Scope

	// root element
	*tui.Box

	// Second-level elements
	topic  *tui.Label
	events *EventsWidget
	// status bar
	connState   *widgets.ConnState
	channelMode *tui.Label
	// input bar
	nick  *tui.Label
	input *tui.Entry
}

// OnKeyEvent handles key presses.
func (v *View) OnKeyEvent(ev tui.KeyEvent) {
	if ev.Key == tui.KeyCtrlC && v.ui != nil {
		v.ui.Quit()
	}
	v.Box.OnKeyEvent(ev)
}

// handleInput handles input from the user.
func (v *View) handleInput(entry *tui.Entry) {
	m := entry.Text()
	defer entry.SetText("")
	lower := strings.ToLower(m)

	if strings.HasPrefix(lower, "/client") && v.ui != nil {
		v.ui.ActivateClient()
		return
	}
	if strings.HasPrefix(lower, "/quit") && v.ui != nil {
		v.ui.Quit()
		return
	}
	if v.sender != nil {
		v.sender.Send(v.scope, m)
	}
}

// SetRenderer sets the function that turns Events into Widgets.
func (v *View) SetRenderer(e EventRenderer) {
	v.events.Renderer = e
}

// UpdateNetwork receives the new state of the network.
func (v *View) UpdateNetwork(n data.NetworkState) {
	update := func() {
		v.nick.SetText(n.Nick)
		v.connState.Set(n.State)
	}

	if v.ui != nil {
		v.ui.Update(update)
	} else {
		update()
	}
}

// UpdateChannel receives the new state of the channel.
func (v *View) UpdateChannel(d data.ChannelState) {
	update := func() {
		v.topic.SetText(d.Topic)
		v.channelMode.SetText(d.Mode)
		v.events.SetLast(d.LastMessage)
	}
	if v.ui != nil {
		v.ui.Update(update)
	} else {
		update()
	}

}

// Filter returns the match rule for this view.
func (v *View) Filter() data.Filter {
	return data.Filter{
		Scope:     v.scope,
		MatchNet:  true,
		MatchName: true,
	}
}

// New returns a new View. It must be run from the main (UI) thread.
func New(s data.Scope, ui UIController, backend backend.Backend) *View {
	// construct V
	v := &View{
		ui:     ui,
		sender: backend,
		scope:  s,

		topic:       tui.NewLabel(""),
		events:      NewEventsWidget(s, backend),
		connState:   widgets.NewConnState(),
		channelMode: tui.NewLabel(""),
		nick:        tui.NewLabel(""),
		input:       tui.NewEntry(),
	}
	v.topic.SetSizePolicy(tui.Expanding, tui.Minimum)
	v.events.SetSizePolicy(tui.Expanding, tui.Expanding)
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
				tui.NewLabel(s.Net),
				tui.NewLabel(": "),
				v.connState,
				tui.NewLabel(" "),
				tui.NewLabel(s.Name),
				tui.NewLabel(": "),
				v.channelMode,
				rspacer,
			),
		},
		inputBar,
	)

	if ui != nil {
		ui.SetWidget(v)
	}

	if backend != nil {
		go backend.Subscribe(v)
	}

	return v
}
