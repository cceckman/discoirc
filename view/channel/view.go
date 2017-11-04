// Package channel includes views and controllers for the Channel view
package channel

import (
	"fmt"
	"image"

	"github.com/marcusolsson/tui-go"
)

const (
	nocon       = "disconnected"
	okcon       = "connected"
	defaultMode = "default"
)

// View is the root of a Channel. All its methods should be called from the UI thread.
type View interface {
	tui.Widget

	// SetLocation updates the view with the current network and channel.
	SetLocation(network, channel string)

	// Connect shows teh UI that the connection is active.
	Connect(ctl *Controller)

	// Disconnect shows the UI that the connection is inactive.
	Disconnect()

	// SetContents updates the Contents view.
	SetContents([]string)

	// ContentSize provides a channel from which the most content size can be read.
	ContentSize() <-chan image.Point

	// SetNick sets the user's name in this channel.
	SetNick(string)
}

func NewView() View {
	w := &view{
		Contents: &Contents{
			List: tui.NewList(),
			SizeUpdate: make(chan image.Point, 1),
		},

		NetBar: tui.NewStatusBar(""),
		modeBar: &modeBar{
			StatusBar: tui.NewStatusBar(""),
			con:       nocon,
			input:     defaultMode,
		},

		Nick: tui.NewLabel(""),
		Input:  tui.NewEntry(),
	}

	// Layout
	w.Contents.SetSizePolicy(tui.Expanding, tui.Expanding)
	w.NetBar.SetSizePolicy(tui.Expanding, tui.Preferred)
	w.modeBar.SetSizePolicy(tui.Expanding, tui.Preferred)
	w.Nick.SetSizePolicy(tui.Preferred, tui.Preferred)
	w.Input.SetSizePolicy(tui.Expanding, tui.Preferred)

	w.Widget = tui.NewVBox(
		w.Contents,
		tui.NewHBox(w.NetBar, w.modeBar),
		tui.NewHBox(w.Nick, w.Input),
	)

	// Initialization
	w.Input.SetFocused(true)
	w.modeBar.render()

	return w
}

// view is the root of the Channel view.
type view struct {
	Contents *Contents

	NetBar   *tui.StatusBar
	modeBar  *modeBar

	Input    *tui.Entry
	Nick *tui.Label

	tui.Widget // root widget
}

// Connect updates the UI to show the connection is active.
func (v *view) Connect(ctl *Controller) {
	v.Input.OnSubmit(func(entry *tui.Entry) {
		ctl.Send(entry.Text())
		entry.SetText("")
	})
	v.modeBar.SetConnected(true)
}

func (v *view) Disconnect() {
	v.Input.OnSubmit(func(_ *tui.Entry) {})
	v.modeBar.SetConnected(false)
}

func (v *view) SetLocation(network, channel string) {
	v.NetBar.SetText(fmt.Sprintf("%s / %s", network, channel))
}

func (v *view) SetNick(nick string) {
	v.Nick.SetText(fmt.Sprintf("<%s>", nick))
}

type modeBar struct {
	*tui.StatusBar

	con, input string
}

func (m *modeBar) Draw(p *tui.Painter) {
	p.WithStyle("reverse", m.StatusBar.Draw)
}

func (m *modeBar) render() {
	m.SetPermanentText(fmt.Sprintf("[%s] [%s]", m.con, m.input))
}

func (m *modeBar) SetConnected(connected bool) {
	if connected {
		m.con = okcon
	} else {
		m.input = nocon
	}
	m.render()
}

type Contents struct {
	*tui.List

	SizeUpdate chan image.Point
}

func (v *view) SetContents(s []string) {
	v.Contents.RemoveItems()
	v.Contents.AddItems(s...)
}

func (v *view) ContentSize() <-chan image.Point {
	return v.Contents.SizeUpdate
}

func (c *Contents) Resize(size image.Point) {
	// Non-blocking, lossy send.
	select {
	case c.SizeUpdate <- size:
		// Sent the value. Mission accomplished.
	case _ = <-c.SizeUpdate:
		// There was a cached value. Send a new one.
		c.SizeUpdate <- size
	}
	c.List.Resize(size)
}
