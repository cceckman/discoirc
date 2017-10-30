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

func NewView() *View {
	w := &View{
		Input:  tui.NewEntry(),
		NetBar: tui.NewStatusBar(""),
		ModeBar: ModeBar{
			StatusBar: tui.NewStatusBar(""),
			con:       nocon,
			input:     defaultMode,
		},
		Contents: &Contents{
			List:    tui.NewList(),
			Resized: make(chan int),
		},
	}

	// Layout
	w.Contents.SetSizePolicy(tui.Expanding, tui.Expanding)
	w.NetBar.SetSizePolicy(tui.Expanding, tui.Preferred)
	w.ModeBar.SetSizePolicy(tui.Expanding, tui.Preferred)
	w.Input.SetSizePolicy(tui.Expanding, tui.Preferred)

	// Initialization
	w.Input.SetFocused(true)
	w.ModeBar.render()

	w.Widget = tui.NewVBox(
		w.Contents,
		tui.NewHBox(w.NetBar, w.ModeBar),
		w.Input,
	)
	return w
}

// View is the root of the Channel view.
type View struct {
	tui.Widget // root widget

	Input    *tui.Entry
	NetBar   *tui.StatusBar
	ModeBar  ModeBar
	Contents *Contents
}

// Connect updates the UI to show the connection is active.
func (v *View) Connect(ctl *Controller) {
	v.Input.OnSubmit(func(entry *tui.Entry) {
		ctl.Send(entry.Text())
		entry.SetText("")
	})
	v.ModeBar.SetConnected(true)
}

func (v *View) SetLocation(network, channel string) {
	v.NetBar.SetText(fmt.Sprintf("%s / %s", network, channel))
}

type ModeBar struct {
	*tui.StatusBar

	con, input string
}

func (m *ModeBar) render() {
	m.SetPermanentText(fmt.Sprintf("[%s] [%s]", m.con, m.input))
}

func (m *ModeBar) SetConnected(connected bool) {
	if connected {
		m.con = okcon
	} else {
		m.input = nocon
	}
	m.render()
}

type Contents struct {
	*tui.List

	Resized chan int
}

func (c *Contents) Set(s []string) {
	c.RemoveItems()
	c.AddItems(s...)
}

func (c *Contents) Resize(size image.Point) {
	c.Resized <- size.X
	c.List.Resize(size)
}
