// Package channel includes views and controllers for the Channel view
package channel

import (
	"fmt"
	"image"

	"github.com/cceckman/tui-go"
)

const (
	disconnected = "disconnected"
	connected    = "connected"

	joining = "joining"
	kicked  = "kicked"
	left    = "left"
)

// View is the root of a Channel. All its methods should be called from the UI thread.
type View interface {
	tui.Widget

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

func NewView(network, channel string) View {
	w := &view{
		Contents: &Contents{
			List:       tui.NewList(),
			SizeUpdate: make(chan image.Point, 1),
		},

		Nick:  tui.NewLabel(""),
		Input: tui.NewEntry(),

		Status: newStatusBar(network, channel),
	}

	// Layout
	w.Contents.SetSizePolicy(tui.Expanding, tui.Expanding)
	w.Status.SetSizePolicy(tui.Expanding, tui.Preferred)
	w.Nick.SetSizePolicy(tui.Preferred, tui.Preferred)
	w.Input.SetSizePolicy(tui.Expanding, tui.Preferred)

	w.Widget = tui.NewVBox(
		w.Contents,
		w.Status,
		tui.NewHBox(w.Nick, w.Input),
	)

	// Initialization
	w.Input.SetFocused(true)

	return w
}

// view is the root of the Channel view.
type view struct {
	Contents *Contents

	Input *tui.Entry
	Nick  *tui.Label

	Status *statusBar

	tui.Widget // root widget
}

type statusBar struct {
	*tui.Box

	Con  *tui.Label
	Mode *tui.Label
}

func newStatusBar(network, channel string) *statusBar {
	r := &statusBar{}

	// elements, in order:
	// network [ <> connected <> ] / #channel [ <> joined <> ]

	networkLabel := tui.NewLabel(fmt.Sprintf("%s [", network))
	networkLabel.SetSizePolicy(tui.Preferred, tui.Preferred)

	r.Con = tui.NewLabel(connected)
	r.Con.SetSizePolicy(tui.Preferred, tui.Preferred)

	channelLabel := tui.NewLabel(fmt.Sprintf("] / %s [", channel))
	channelLabel.SetSizePolicy(tui.Preferred, tui.Preferred)

	r.Mode = tui.NewLabel(joining)
	r.Mode.SetSizePolicy(tui.Preferred, tui.Preferred)

	endcap := tui.NewLabel("]")
	endcap.SetSizePolicy(tui.Expanding, tui.Preferred)

	r.Box = tui.NewHBox(networkLabel, r.Con, channelLabel, r.Mode, endcap)
	return r
}

func (m *statusBar) Draw(p *tui.Painter) {
	p.WithStyle("reverse", m.Box.Draw)
}

// Connect updates the UI to show the connection is active.
func (v *view) Connect(ctl *Controller) {
	v.Input.OnSubmit(func(entry *tui.Entry) {
		ctl.Send(entry.Text())
		entry.SetText("")
	})
	v.Status.Con.SetText(connected)
}

func (v *view) Disconnect() {
	v.Input.OnSubmit(func(_ *tui.Entry) {})
	v.Status.Con.SetText(disconnected)
}

func (v *view) SetNick(nick string) {
	v.Nick.SetText(fmt.Sprintf("<%s> ", nick))
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
