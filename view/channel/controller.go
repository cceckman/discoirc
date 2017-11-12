package channel

import (
	"context"

	"github.com/golang/glog"
	"github.com/marcusolsson/tui-go"

	"github.com/cceckman/discoirc/model"
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
	// Wait for Channel object to be availabe.
	ch := <-getCh

	ctl.sendMessages(ctx, ch)
	ctl.updateContents(ctx, ch, ctl.rerange(ctx, ch))
	ctl.updateState(ctx, ch)
}

// updateState updates the static portions of the view: connection state, nick, etc.
func (ctl *Controller) updateState(ctx context.Context, ch model.Channel) {
	go func() {
		stateChan := ch.State(ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case newState := <-stateChan:
				glog.V(1).Infof("controller got state notification: %+v", newState)
				await := make(chan struct{})
				ctl.UI.Update(func() {
					defer close(await)
					if newState.Connected {
						ctl.View.Connect(ctl)
					} else {
						ctl.View.Disconnect()
					}
					ctl.View.SetMode(newState.Mode)
					ctl.View.SetNick(newState.Nick)
					ctl.View.SetTopic(newState.Topic)
				})
				<-await
			}
		}
	}()
}

// sendMessages is the send loop, passing messages from the UI to the network.
func (ctl *Controller) sendMessages(ctx context.Context, ch model.Channel) {
	messages := make(chan string)
	go func() {
		for msg := range messages {
			ch.SendMessage(msg)
		}
	}()
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
				case messages <- hd:
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

// rerange is the layout loop, handling resize and message-received events.
func (ctl *Controller) rerange(ctx context.Context, ch model.Channel) chan int {
	newRange := make(chan int, 1)
	go func() {
		defer close(newRange)

		// Listen for resize events
		notices := ch.Events(ctx)
		var size int
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
				if newSize.Y == size {
					// Don't need to resize; ignore.
					break
				}
				glog.V(1).Infof("rerange: size changed to %v", newSize)
				size = newSize.Y
				// Resize with non-blocking / lossy write.
				select {
				case newRange <- size:
					// pass
				case _ = <-newRange:
					newRange <- size
				}
			case _, ok := <-notices:
				glog.V(1).Infof("rerange: received notice of new content")
				if !ok {
					return
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
func (ctl *Controller) updateContents(ctx context.Context, ch model.Channel, update chan int) {
	go func() {
		updateDone := make(chan *struct{})
		for size := range update {
			glog.V(1).Infof("updateContents got size update of %v", size)
			events := ch.SelectSize(uint(size))
			messages := make([]string, len(events))
			for i, event := range events {
				messages[i] = event.Contents
			}
			// Schedule the GUI update and block on its completion before
			// continuing to pick up the new size.
			go ctl.UI.Update(func() {
				if len(events) == 0 {
					glog.V(1).Infof("showing zero events")
				} else {
					last := len(events) - 1
					glog.V(1).Infof("showing %d messages, from (%d, %d) to (%d, %d)",
						len(messages), events[0].Epoch, events[0].Seq, events[last].Epoch, events[last].Seq)
				}
				ctl.View.SetContents(messages)
				updateDone <- nil
			})
			_ = <-updateDone
		}
	}()
}

func New(ctx context.Context, ui tui.UI, client model.Client, network, channel string) tui.Widget {
	ctl := &Controller{
		View:    NewView(network, channel),
		UI:      ui,
		msgSend: make(chan string, 1),
	}

	gotChannel := make(chan model.Channel)
	// Start / establish connection in background.
	go func() {
		gotChannel <- client.Connection(network).Channel(channel)
	}()

	go ctl.Start(ctx, gotChannel)
	return ctl.View
}
