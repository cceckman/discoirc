//
// modelview.go
// provides the ModelView class.

package mvvm

import (
	"context"
	"fmt"
	"log"

	"github.com/cceckman/discoirc/prototype/bufchan"
	"github.com/jroimartin/gocui"
)

const (
	messagesView = "messages"
	inputView    = "input"
	noticeView   = "notices"
)

// ModelView provides Go interfaces to the UI behavior.
type ModelView struct {
	// These channels handle the interface between the UI and Model.
	input *bufchan.Bufchan

	// ui is the UI object this ModelView is attached to.
	ui *gocui.Gui
	// models is a channel for receiving Attach requests.
	models chan Model
}

// attachToGui sets up a ModelView for the provided Gui.
func (mv *ModelView) attachToGui() error {
	// Create a context that closing the UI terminates.
	ctx, cancel := context.WithCancel(context.Background())

	// Start window layout-er.
	// Note: manager must be provided before setting keybindings (e.g. below.)
	go mv.start(ctx)
	mv.ui.SetManager(mv)

	// Global handler for ctrl+c.
	if err := mv.ui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(_ *gocui.Gui, _ *gocui.View) error {
			cancel()
			return gocui.ErrQuit
		},
	); err != nil {
		log.Println(err)
		return err
	}

	// Bind 'enter' to close, on the notice view.
	err := mv.ui.SetKeybinding(noticeView, gocui.KeyEnter, gocui.ModNone, closeView)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// New returns a new ModelView, attached to the given GUI.
func New(g *gocui.Gui) (*ModelView, error) {
	r := &ModelView{
		ui: g,
		models: make(chan Model),
	}
	return r, r.attachToGui()
}

// Attach starts running this Model against this ModelView.
func (mv *ModelView) Attach(m Model) {
	mv.models <- m
}

// Notice sends a notice to the user.
// It is non-blocking; it will return before the notice is actually displayed.
func (mv *ModelView) Notice(notice string) {
	mv.ui.Execute(displayNotice(notice))
}

// UserInput returns a channel on which the user's input can be read.
// It is not mirrored- first consumer gets it!
func (mv *ModelView) UserInput() <-chan string {
	return mv.input.Out()
}

// Message writes a message to the user's log.
// It is non-blocking; it will return before the message is actually displayed.
func (mv *ModelView) Message(message string) {
	if len(message) == 0 {
		return
	}
	mv.ui.Execute(func(g *gocui.Gui) error {
		if v, err := g.View(messagesView); err != nil {
			log.Println(err)
			return err
		} else {
			fmt.Fprintf(v, "%s\n", message)
		}
		return nil
	})
}

// Type enforcement.
var _ gocui.Manager = &ModelView{}
var _ gocui.Editor = &ModelView{}

// start begins background operations: taking inputs, and attaching Models.
func (mv *ModelView) start(ctx context.Context) {
	mv.input = bufchan.New(ctx)

	// Hang around, listening for Models to be attached.
	for {
		select {
		case <-ctx.Done():
			return
		case m := <-mv.models:
			go m.Run(ctx, mv)
		}
	}
}


// Edit implements gocui.Editor for ModelView.
// When a line is entered from the input, the buffer is cleared, and the input is sent to m.Input.
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

	// Messages view: auto-scrolling, from m.Messages.
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

// closeView is a callback that closes the View it is called on.
// It is idempotent; it can be run on already-removed Views.
func closeView(g *gocui.Gui, v *gocui.View) error {
	// Clean up the view, OK if it already doesn't exist.
	if err := g.DeleteView(v.Name()); err != nil {
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
