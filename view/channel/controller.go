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

	// Send loop: pass sent messages to the client.
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

		// Send any queued messages to the client.
		for {
			for len(queuedMessages) > 0 {
				// Pop a message, attempt to forward it.
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

	newRange := make(chan uint, 1)
	// Receive/resize loop; get new contents for UI.
	// TODO: This should be its own method and/or class.
	go func() {
		defer close(newRange)
		// Wait for channel to be ready.
		ch := <-gotChannel
		gotChannel <- ch

		// Listen for resize events
		notices := ch.Await(ctx)
		var size uint
		var msgCount int
		// Await resize or more messages received.
		// TODO support non-zero start.
		for {
			select {
			case <-ctx.Done():
				return
			case newSize, ok := <-ctl.View.Contents.Resized:
				if !ok {
					return
				}
				if uint(newSize) == size {
					// Don't need to resize; ignore.
					break
				}
				size = uint(newSize)
				// Resize with non-blocking / lossy write.
				select {
				case newRange <- size:
					// pass
				case _ = <-newRange:
					newRange <- size
				}
			case notice, ok := <-notices:
				if !ok {
					return
				}
				if notice.Messages == msgCount {
					// Don't need to resize; ignore
					break
				}
				select {
				case newRange <- size:
					// pass
				case _ = <-newRange:
					newRange <- size
				}
			}
		}
	}()

	// Refresh thread: read the new range, get its contents, and update the
	// display.
	newContents := make(chan []string, 1)
	go func() {
		defer close(newContents)
		// Wait for channel to be ready.
		ch := <-gotChannel
		gotChannel <- ch

		updateDone := make(chan *struct{})
		for size := range newRange {
			messages := ch.GetMessages(0, size)
			// Schedule the GUI update and block on its completion before
			// continuing to pick up the new size.
			ctl.UI.Update(func() {
				ctl.View.Contents.Set(messages)
				updateDone <- nil
			})
			_ = <-updateDone
		}
	}()

}

func New(ctx context.Context, view *View, ui tui.UI, client model.Client, network, channel string) {
	ctl := &Controller{
		Client:  client,
		View:    view,
		UI:      ui,
		msgSend: make(chan string, 1),
	}

	ctl.Start(ctx, network, channel)
}
