//
// editgui.go
// provides GUI objects, etc. for this prototype

package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/cceckman/discoirc/prototype/bufchan"
	"github.com/jroimartin/gocui"
)

const (
	messagesView = "messages"
	inputView    = "input"
	noticeView   = "notices"
)

// SetupModelView sets up managers for a gocui.Gui, but does not start the main loop.
func SetupModelView(g *gocui.Gui) error {
	// Create a context that closing the UI terminates.
	ctx, cancel := context.WithCancel(context.Background())
	_ = ctx

	// Start window layout-er.
	// Note: manager must be provided before setting keybindings (e.g. below.)
	mv := New(g)
	go mv.Start(ctx)
	g.SetManager(mv)

	// Global handler for ctrl+c.
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(_ *gocui.Gui, _ *gocui.View) error {
			cancel()
			return gocui.ErrQuit
		},
	); err != nil {
		log.Println(err)
		return err
	}

	// Bind 'enter' to close, on the notice view.
	err := g.SetKeybinding(noticeView, gocui.KeyEnter, gocui.ModNone, closeView)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// New returns a new ModelView
func New(g *gocui.Gui) *ModelView {
	return &ModelView{
		ui: g,
	}
}

// ModelView is a view manager.
// It should be Start-ed before being attached as a Manager.
type ModelView struct {
	ui *gocui.Gui

	notices  *bufchan.Bufchan
	input    *bufchan.Bufchan
	messages *bufchan.Bufchan
}

// Type enforcement.
var _ gocui.Manager = &ModelView{}
var _ gocui.Editor = &ModelView{}

// Start begins operations that run outside the main thread. It should be run in a background thread (i.e. go m.Start())
func (m *ModelView) Start(ctx context.Context) {
	m.notices = bufchan.New(ctx)
	m.input = bufchan.New(ctx)
	m.messages = bufchan.New(ctx)

	// Generate some output for testing.
	// go m.genOuts(ctx)

	go m.WatchInput(ctx)
	go m.WriteMessages(ctx)
	go m.WriteNotices(ctx)
}

// genOuts writes numbers to the input channel.
func (m *ModelView) genOuts(ctx context.Context) {
	tick := time.NewTicker(time.Second * 2)
	defer tick.Stop()
	for i := 0; true; i++ {
		select {
		case <-ctx.Done():
			return
		case m.input.In() <- fmt.Sprintf(" %d\n", i):
			// Delay until the tick.
			<-tick.C
		}
	}
}

// WatchInput watches the input channel, and demuxes into 'messages' and 'notices'.
func (m *ModelView) WatchInput(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case input := <-m.input.Out():
			// TODO: I'm being unfriendly; RTL should absolutely be supported by this app.
			tr := strings.TrimRightFunc(input, unicode.IsSpace)
			if len(tr) == 0 {
				continue
			}

			if tr[0] == '!' {
				m.notices.In() <- tr[1:]
			} else {
				m.messages.In() <- tr
			}
		}
	}
}

// WriteMessages listens on the relevant channel, and writes messages to the UI.
func (m *ModelView) WriteMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-m.messages.Out():
			if len(message) == 0 {
				continue
			}
			m.ui.Execute(func(g *gocui.Gui) error {
				if v, err := g.View(messagesView); err != nil {
					log.Println(err)
					return err
				} else {
					fmt.Fprintf(v, "\"%s\"\n", message)
				}
				return nil
			})
		}
	}
}

// closeView is a callback that closes the View it is called on.
// It is idempotent; it can be run on already-removed Views.
func closeView(g *gocui.Gui, v *gocui.View) error {
	// Clean up the view, OK if it already doesn't exist.
	if err := g.DeleteView(v.Name()); err != nil  {
		if err == gocui.ErrUnknownView {
			return nil
		}
		log.Println(err)
		return err
	}
	return nil
}

// displayNotice displays a message in a 'notice' box.
func displayNotice(notice string) func(*gocui.Gui) error {
	return func(g *gocui.Gui) error {
		maxX, maxY := g.Size()
		l := len(notice) / 2
		if v, err := g.SetView(
			noticeView,
			maxX/2-l-1, maxY/2,
			maxX/2+l+1, maxY/2+2,
		); err != nil {
			if err != gocui.ErrUnknownView {
				log.Println(err)
				return err
			}
			// TODO: This isn't quite the right handling of "a new notice"...
			// This overwrites whatever is there, which means repeated notices can get squashed by
			// each other.
			v.Clear()
			v.SetCursor(0, 0)
			g.SetViewOnTop(noticeView)
			fmt.Fprintln(v, notice)

			if _, err := g.SetCurrentView(noticeView); err != nil {
				log.Println(err)
				return err
			}
		}
		return nil
	}
}

// WriteNotices listens on the relevant channel, and writes pop-up notifications to the UI.
func (m *ModelView) WriteNotices(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case notice := <-m.notices.Out():
			if len(notice) == 0 {
				continue
			}
			m.ui.Execute(displayNotice(notice))

		}
	}
}

// Edit implements gocui.Editor for ModelView.
// When a line is entered from the input, the buffer is cleared, and the input is sent to m.input.
func (m *ModelView) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyEnter:
		// Commit this line to the input channel.
		s := v.Buffer()
		m.input.In() <- s
		v.Clear()
		v.SetCursor(0, 0)
		/* // Scrolling disabled, at the moment...
		case key == gocui.KeyArrowDown:
			v.MoveCursor(0, 1, false)
		case key == gocui.KeyArrowUp:
			v.MoveCursor(0, -1, false)
		*/
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
}

// Layout implements gocui.Manager for ModelView.
func (m *ModelView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	inputHeight := 2

	// Input view. Sink to the bottom of the screen.
	if v, err := g.SetView(inputView, 0, maxY-inputHeight, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			log.Println(err)
			return err
		}
		v.Editable = true
		v.Frame = false
		v.Wrap = true
		v.Editor = m
	}

	// If the notice view exists, it gets focus...
	focus := noticeView
	if _, err := g.View(noticeView); err == gocui.ErrUnknownView {
		// ...otherwise, put focus on the input view, to get keyboard input.
		// TODO: allow focus to swap between input and messages.
		focus = inputView
	}
	if _, err := g.SetCurrentView(focus); err != nil {
		log.Println(err)
		return err
	}

	// Messages view: auto-scrolling, from m.messages.
	// Set its bottom edge to just above the input view.
	if v, err := g.SetView(messagesView, 0, 0, maxX-1, maxY-inputHeight); err != nil {
		if err != gocui.ErrUnknownView {
			log.Println(err)
			return err
		}
		v.Autoscroll = true
		v.Title = "Messages"
	}
	return nil
}
