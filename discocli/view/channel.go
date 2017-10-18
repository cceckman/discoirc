// Package chat provides the Channel view/model/viewmodel for the IRC channel view.
package view

import (
	"context"
	"fmt"
	"log"

	"github.com/cceckman/discoirc/discocli/model"
	"github.com/jroimartin/gocui"
)

// ChannelViewModel is the ViewModel for the Channel view.
type ChannelViewModel struct {
	Connection, Channel string

	Log  *log.Logger

	// channel is not necessarily populated until 'connected' is closed.
	channel   model.Channel
	connected chan struct{}
}

func NewChannelViewModel(connection, channel string, client model.Client, log *log.Logger) gocui.Manager {
	result := &ChannelViewModel{
		Connection: connection,
		Channel:    channel,
		Log:        log,

		connected: make(chan struct{}),
	}
	go result.Connect(client)
	return result
}

func (m *ChannelViewModel) Connect(client model.Client) {
	// Connect in the background.
	// Once done, signal to waiting layout routines that it's OK to add handlers.
	defer close(m.connected)

	// TODO allow for an error in connection.
	m.channel = client.Connection(m.Connection).Channel(m.Channel)
}

var _ gocui.Manager = &ChannelViewModel{}

// Layout sets up the Channel view. It creates new views as necessary, including starting threads.
func (m *ChannelViewModel) Layout(g *gocui.Gui) error {
	m.Log.Print("Channel: [start] layout")
	defer m.Log.Print("Channel: [done] layout")
	// Create three views: input, status, messages.
	// Create them in that order, since 'input' is fixed-len, but 'status' and 'messages' may need to flex.
	if err := m.layoutInput(g); err != nil {
		return err
	}
	if err := m.layoutStatus(g); err != nil {
		return err
	}
	if err := m.layoutMessages(g); err != nil {
		return err
	}
	g.SetCurrentView("input")
	return nil
}

func Quit(*gocui.Gui, *gocui.View) error {
	return gocui.ErrQuit
}

func (m *ChannelViewModel) addInputHandlers(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit)
		return nil
	})
	<-m.connected
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("input")
		switch {
		case err == gocui.ErrUnknownView:
			return nil
		case err != nil:
			return err
		}
		NewMessageEditor(m.channel, v)
		return nil
	})
}

func (m *ChannelViewModel) layoutInput(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	// These are good!
	ax, ay, bx, by := -1, maxY-2, maxX, maxY
	v, err := g.SetView("input", ax, ay, bx, by)
	switch err {
	case nil:
		return nil
	case gocui.ErrUnknownView:
		m.Log.Print("Channel/input: [start] initial layout")
		defer m.Log.Print("Channel/input: [done] initial layout")
		v.Frame = false
		go m.addInputHandlers(g)
	default:
		return err
	}
	return nil
}

func (m *ChannelViewModel) addStatusHandlers(g *gocui.Gui) {
	<-m.connected

	// Create a stream of topics to update.
	topics := make(chan string)

	ctx, cancel := context.WithCancel(context.Background())
	// Worker: Consume notifications, issue updates to the channel as available.
	go func() {
		defer close(topics)
		notifications := m.channel.Await(ctx)

		// Initialize by getting the current topic.
		topicReceived := m.channel.GetTopic()
		topics <- topicReceived
		topicSent := topicReceived

		for {
			if topicReceived == topicSent {
				// Just await a new notification.
				select {
				case <-ctx.Done():
					return
				case n := <-notifications:
					topicReceived = n.Topic
				}
			} else {
				// We have a topic to send, but may get a notification in the mean time.
				select {
				case <-ctx.Done():
					return
				case n := <-notifications:
					topicReceived = n.Topic
				case topics <- topicReceived:
					topicSent = topicReceived
				}
			}
		}
	}()

	// Worker: Take topics from the queue, issue updates to the UI thread.
	// Shut down the context (and therefore the producer) if the window is gone.
	go func() {
		for topic := range topics {
			// Serialize topic updates, since Update isn't serialized itself.
			done := make(chan struct{})
			g.Update(func(g *gocui.Gui) error {
				defer close(done)
				v, err := g.View("status")
				switch {
				case err == gocui.ErrUnknownView:
					// If the view is gone, stop working.
					cancel()
					return nil
				case err != nil:
					return err
				}
				v.Clear()
				_, err = fmt.Fprintf(v, "%s %s %s", m.Connection, m.Channel, topic)
				return err
			})
			<-done
		}
	}()
}

func (m *ChannelViewModel) layoutStatus(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	ax, ay, bx, by := -1, maxY-3, maxX, maxY-1
	if v, err := g.SetView("status", ax, ay, bx, by); err != nil {
		if err != gocui.ErrUnknownView {
			// unknown error, percolate it up
			return err
		}
		// Initialize the view.
		m.Log.Print("Channel/status: [start] initial layout")
		defer m.Log.Print("Channel/status: [done] initial layout")

		v.Editable = false
		v.Frame = false
		// TODO more color customization
		v.BgColor, v.FgColor = gocui.ColorWhite, gocui.ColorBlack
		// TODO attach controller/model here, instead of static init.
		if _, err := fmt.Fprintf(v, "%s / %s", m.Connection, m.Channel); err != nil {
			return err
		}

		go m.addStatusHandlers(g)
	}
	return nil
}

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
