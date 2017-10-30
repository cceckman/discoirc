package channel

import (
	"context"

	"github.com/cceckman/discoirc/model"
	"github.com/marcusolsson/tui-go"
)

type Controller struct {
	View *View
	UI   tui.UI

	Client model.Client

	msgSend chan string
	resize  chan int
}

// Send queues a message for sending.
func (ctl *Controller) Send(msg string) {
	ctl.msgSend <- msg
}

func (ctl *Controller) Start(ctx context.Context, network, channel string) {
	gotChannel := make(chan model.Channel, 1)

	// Start / establish connection in background.
	go func() {
		gotChannel <- ctl.Client.Connection(network).Channel(channel)
	}()

	// Update UI with channel.
	ctl.UI.Update(func() {
		ctl.View.SetLocation(network, channel)
	})

	// Connection state update.
	// TODO: handle this more generally- as connected / disconnected state.
	go func() {
		select {
		case <-ctx.Done():
			return
		case ch := <-gotChannel:
			go ctl.UI.Update(func() {
				ctl.View.Connect(ctl)
			})
			gotChannel <- ch
		}
	}()

	// Send loop: pass messages through to client.
	go func() {
		queuedMessages := []string{}
		var ch model.Channel
		// Wait for connection to establish.
	connectLoop:
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-ctl.msgSend:
				queuedMessages = append(queuedMessages, msg)
			case c := <-gotChannel:
				ch = c
				gotChannel <- c
				break connectLoop
			}
		}
		// Start sending messages.
		for {
			for len(queuedMessages) > 0 {
				// Pop a message in order to send to network.
				var hd string
				hd, queuedMessages = queuedMessages[0], queuedMessages[1:]
				select {
				case <-ctx.Done():
					return
				case msg := <-ctl.msgSend:
					queuedMessages = append(queuedMessages, msg)
				case ch.MessageInput() <- hd:
					// pass
				}
			}
			// Wait for another message to come in.
			select {
			case <-ctx.Done():
				return
			case msg := <-ctl.msgSend:
				queuedMessages = append(queuedMessages, msg)
			}
		}
	}()

	// Resize loop: get the most recent size, pass it on.
	// Isolates the UI thread (writing to ctl.View.Contents.Resized)
	// from the getting-messages thread (reading from size)
	sizeUpdate := make(chan uint)
	go func() {
		defer close(sizeUpdate)
		var newSize int
		var lastSize int
		for {
			// Upper loop: read a new value
			select {
			case <-ctx.Done():
				return
			case newSize = <-ctl.View.Contents.Resized:
				// pass
			}
			if newSize == lastSize {
				continue
			}
			// Lower loop: pass on old value, or read a new value
		lowerLoop:
			for {
				select {
				case <-ctx.Done():
					return
				case newSize = <-ctl.View.Contents.Resized:
					// pass
				case sizeUpdate <- uint(newSize):
					lastSize = newSize
					break lowerLoop
				}
			}
		}
	}()

	newMessages := make(chan []string)
	// Refresh thread: Schedule message updates against UI thread.
	// Use an independent channel for this s.t. writes are ordered regardless of
	// whether Update is.
	go func() {
		for {
			var messages []string
			var ok bool
			// Upper loop: await new messages, have none to send.
			select {
			case <-ctx.Done():
				return
			case messages, ok = <-newMessages:
				if !ok {
					return
				}
			}
			// Schedule a UI thread that joins on contents.
			contents := make(chan []string)
			ctl.UI.Update(func() {
				ctl.View.Contents.Set(<-contents)
			})
			defer close(contents)
			// Lower loop: await either more messages, or the UI thread to consume.
		lowerLoop:
			for {
				select {
				case <-ctx.Done():
					return
				case messages, ok = <-newMessages:
					if !ok {
						return
					}
				case contents <- messages:
					break lowerLoop
				}
			}
		}
	}()

	// Receive/resize loop; get new contents for UI.
	// TODO: This should be its own method and/or class.
	go func() {
		defer close(newMessages)
		// Wait for channel to be ready.
		ch := <-gotChannel
		gotChannel <- ch

		// Listen for resize events
		notices := ch.Await(ctx)
		var size uint
		var ok bool
		var msgCount int
		for {
			select {
			case <-ctx.Done():
				return
			case size, ok = <-sizeUpdate:
				if !ok {
					return
				}
				// Resize event. Re-fetch messages.
				// TODO: allow for non-zero-index, i.e. using EventRange.
				newMessages <- ch.GetMessages(0, size)
			case notice, ok := <-notices:
				if !ok {
					return
				}
				if notice.Messages != msgCount {
					msgCount = notice.Messages
					newMessages <- ch.GetMessages(0, size)
				}
			}
		}
	}()
}

func New(ctx context.Context, view *View, ui tui.UI, client model.Client, network, channel string) {
	ctl := &Controller{
		Client:  client,
		View:    view,
		UI:      ui,
		msgSend: make(chan string),
	}

	ctl.Start(ctx, network, channel)
}
