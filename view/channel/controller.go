package channel

import (
	"context"

	"github.com/cceckman/discoirc/model"
	"github.com/marcusolsson/tui-go"
)

type Controller struct {
	View View
	UI   tui.UI

	msgSend chan string
	resize  chan int
}

// Send queues a message for sending.
func (ctl *Controller) Send(msg string) {
	ctl.msgSend <- msg
}

func (ctl *Controller) Start(ctx context.Context, getCh chan model.Channel) {
	// Wait for channel to connect.
	ch := <-getCh

	ctl.sendMessages(ctx, ch)
	newRange := ctl.rerange(ctx, ch)
	ctl.updateContents(ctx, ch, newRange)

	// Update UI to indicate connection.
	go ctl.UI.Update(func() {
		ctl.View.Connect(ctl)
	})
}

// sendMessages is the send loop, passing messages from the UI to the network.
func (ctl *Controller) sendMessages(ctx context.Context, ch model.Channel) {
	go func() {
		queuedMessages := []string{}
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
}

// reerange loop, getting new contents fo the UI.
func (ctl *Controller) rerange(ctx context.Context, ch model.Channel) chan uint {
	newRange := make(chan uint, 1)
	go func() {
		defer close(newRange)

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
			case newSize, ok := <-ctl.View.ContentSize():
				if !ok {
					return
				}
				if uint(newSize.X) == size {
					// Don't need to resize; ignore.
					break
				}
				size = uint(newSize.X)
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
				// New messages received. Redraw. We signal that as resize.
				select {
				case newRange <- size:
					// pass
				case _ = <-newRange:
					newRange <- size
				}
			}
		}
	}()
	return newRange
}

// updateContents updates the contents of the View
func (ctl *Controller) updateContents(ctx context.Context, ch model.Channel, update chan uint) {
	go func() {
		updateDone := make(chan *struct{})
		for size := range update {
			messages := ch.GetMessages(0, size)
			// Schedule the GUI update and block on its completion before
			// continuing to pick up the new size.
			go ctl.UI.Update(func() {
				ctl.View.SetContents(messages)
				updateDone <- nil
			})
			_ = <-updateDone
		}
	}()
}

func New(ctx context.Context, view View, ui tui.UI, client model.Client, network, channel string) {
	ctl := &Controller{
		View:    view,
		UI:      ui,
		msgSend: make(chan string, 1),
	}

	// Queue UI update with channel location.
	go ctl.UI.Update(func() {
		ctl.View.SetLocation(network, channel)
	})

	gotChannel := make(chan model.Channel)
	// Start / establish connection in background.
	go func() {
		gotChannel <- client.Connection(network).Channel(channel)
	}()

	go ctl.Start(ctx, gotChannel)
}
