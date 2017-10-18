package view

import (
	"context"
	"fmt"
	"log"

	"github.com/cceckman/discoirc/discocli/model"
	"github.com/jroimartin/gocui"
)


const (
	// View names.
	ChannelInputView = "channel input"
	ChannelStatusView = "channel status"
	ChannelContentsView = "channel contents"
)


// So: what's the flow here?
// Either the process startup, or some action in a different top-level view, decides that a given
// *gocui.Gui should be using a particular model.Client as a backend, with logging to a log.Logger.

// Context provides data necessary for all Windows.
type Context struct {
	Gui *gocui.Gui
	Log *log.Logger
	Backend model.Client
}

// Window is a top-level view, e.g. Channel or Session.
type Window interface {
	// Start replaces the Gui with this Window, or returns an error.
	Start() error
}


// Channel is the ViewModel for the Channel view.
type Channel struct {
	*Context

	Connection, Channel string

	// channel is not necessarily populated until 'connected' is closed.p
	// connected blocks some operations until the channel is properly connected from the client side.
	connected chan struct{}
	channel   model.Channel
}

func (vm *Channel) validate() error {
	if vm.Gui == nil {
		return errors.New("no Gui provided")
	}
	if vm.Log == nil {
		return errors.New("no Logger provided")
		}
	if vm.Connection == "" {
		return errors.New("no Connection provided")
	}
	if vm.Channel == "" {
		return errors.New("no Channel provided")
	}
	return nil
}

func (vm *Channel) Start() error {
	if err := vm.validate(); err != nil {
		return err
	}
	// Start client connection.
	vm.connected = make(chan struct)
	go func() {
		defer close(vm.connected)
		vm.channel = vm.Client.Connection(vm.Connection).Channel(vm.Channel)
	}()

	// Attach ViewModels.
	vm.Gui.SetManagers(
		vm.NewInput(),
		vm.NewStatus(),
		vm.NewContents(),
		QuitManager,
	)

	return nil
}

// QuitManager is a Manager that provides a Ctrl+C quit handler.
var QuitManager gocui.Manager = gocui.ManagerFunc(func(g *gocui.Gui) error {
		g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
			return gocui.ErrQuit
		})
		return nil
	})

func (m *ChannelViewModel) addMessagesHandlers(g *gocui.Gui) {
	<-m.connected

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		notifications := m.channel.Await(ctx)
		count := 0
		for n := range notifications {
			// Suppress topic-only updates.
			if n.Messages == count {
				continue
			}
			count = n.Messages

			done := make(chan struct{})
			g.Update(func(g *gocui.Gui) error {
				defer close(done)
				v, err := g.View("messages")
				switch {
				case err == gocui.ErrUnknownView:
					cancel()
					return nil
				case err != nil:
					return err
				}
				// TODO: Refactor this controller, s.t. this doesn't take place in the UI thread.
				_, lines := v.Size()
				// TODO: allow scrollback. Part of refactoring.
				messages := m.channel.GetMessages(0, uint(lines))
				v.Clear()
				for _, m := range messages {
					fmt.Fprintln(v, m)
				}
				return nil
			})
			<-done
		}
	}()
}

func (m *ChannelViewModel) layoutMessages(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	ax, ay, bx, by := 0, 0, maxX-1, maxY-3
	if v, err := g.SetView("messages", ax, ay, bx, by); err != nil {
		if err != gocui.ErrUnknownView {
			// unknown error, percolate it up
			return err
		}
		// Initialize the view.
		m.Log.Print("Channel/messages: [start] initial layout")
		defer m.Log.Print("Channel/messages: [done] initial layout")
		v.Editable = false
		v.Frame = true

		// TODO attach controller/model here, instead of fake init.
		v.Autoscroll = true
		go m.addMessagesHandlers(g)
	}

	return nil
}
