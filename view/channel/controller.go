package channel

import (
	"context"

	"github.com/cceckman/discoirc/model"
	"github.com/marcusolsson/tui-go"
)

type Controller struct {
	View *View
	UI tui.UI

	Client           model.Client

	msgSend chan string
	resize chan int
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


	// Resize loop: get the current size.
	size := make(chan int)
	go func() {
		for {
			var x int
			select {
			case <-ctx.Done():
				return
			case x = <-ctl.View.Contents.Resized:
				// pass
			case size <- x:
				// pass
			}
		}
	}()

	// Receive/resize loop; pass messages to UI.
	go func() {
		// Wait for channel to be ready.
		ch := <- gotChannel
		gotChannel <- ch

		// Listen for resize events
		notices := ch.Await(ctx)
		msgCount := 0
		var messages []string
		var sz int
		for {
			select {
			case <-ctx.Done():
				return
			case sz = <-size:
				// Resize, so refetch.
				// TODO: allow for starting not at the end.
				messages = ch.GetMessages(0, uint(sz))
			case notice := <-notices:
				if notice.Messages != msgCount {
					msgCount = notice.Messages
					messages = ch.GetMessages(0, uint(sz))
				}
			}
			ctl.View.Contents.Set(messages)
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
