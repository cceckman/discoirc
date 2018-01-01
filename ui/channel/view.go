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
	r := tui.NewLabel(e.Contents)
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
	ui              UIController
	sender          backend.Sender
	network, target string

	// root element
	*tui.Box

	// Second-level elements
	topic  *tui.Label
	events *EventsWidget
	// status bar
	// network     *tui.Label
	connState *widgets.ConnState
	// target     *tui.Label
	channelMode *tui.Label
	// input bar
	nick  *tui.Label
	input *tui.Entry
}

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
		v.sender.Send(v.network, v.target, m)
	}
}

func (v *View) SetRenderer(e EventRenderer) {
	v.events.Renderer = e
}

func (v *View) UpdateNetwork(n data.NetworkState) {
	if n.Network != v.network {
		return
	}

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

func (v *View) UpdateChannel(d data.ChannelState) {
	if d.Network != v.network || d.Channel != v.target {
		return
	}

	update := func() {
		v.topic.SetText(d.Topic)
		v.channelMode.SetText(d.ChannelMode)
		v.events.SetLast(d.LastMessage.EventID)
	}
	if v.ui != nil {
		v.ui.Update(update)
	} else {
		update()
	}

}

// Filter indicates the network and target this widget should receive updates for.
func (v *View) Filter() (string, string) {
	return v.network, v.target
}

// New returns a new View. It must be run from the main (UI) thread.
func NewView(network, target string, ui UIController, backend backend.Backend) *View {
	// construct V
	v := &View{
		ui:     ui,
		sender: backend,

		network: network,
		target:  target,

		topic:       tui.NewLabel(""),
		events:      NewEventsWidget(network, target, backend),
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
				tui.NewLabel(network),
				tui.NewLabel(": "),
				v.connState,
				tui.NewLabel(" "),
				tui.NewLabel(target),
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
		backend.SubscribeFiltered(v)
	}

	return v
}
