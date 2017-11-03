// Package model provides models for chat and connection updates.
package model

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Channel interface {
	Events

	Name() string
	Network() string

	Send(string)
	Updates(context.Context) <-chan Event
}

type Notification struct {
	Latest Event
	Next   chan Notification
}

func (c *MockChannel) Name() string {
	return c.name
}

func (c *MockChannel) Network() string {
	return c.network
}

// MockChannel implements the Channel interface.
type MockChannel struct {
	name    string
	network string
	log     *log.Logger

	request chan Events

	subscribe chan chan Notification
	send      chan string
}

func (c *MockChannel) SelectSize(n uint) []Event {
	return (<-c.request).SelectSize(n)
}

func (c *MockChannel) SelectSizeMax(n uint, e EventID) []Event {
	return (<-c.request).SelectSizeMax(n, e)
}

func (c *MockChannel) SelectMinSize(e EventID, n uint) []Event {
	return (<-c.request).SelectMinSize(e, n)
}

func (c *MockChannel) SelectMinMax(min, max EventID) []Event {
	return (<-c.request).SelectMinMax(min, max)
}

func (c *MockChannel) Updates(ctx context.Context) <-chan Event {
	c.log.Printf("added listener to channel %s / %s", c.Network(), c.Name())
	result := make(chan Event)
	go func() {
		// Block on subscription request.
		var notices chan Notification
		select {
		case <-ctx.Done():
			return
		case notices = <-c.subscribe:
			// Have a channel to listen on.
		}

		for {
			select {
			case <-ctx.Done():
				return
			case notice := <-notices:
				// Put it back in the broadcast channel immediately.
				notices <- notice
				result <- notice.Latest
			}
		}
	}()
	return result
}

func (c *MockChannel) Send(msg string) {
	c.log.Printf("awaiting send for message \"%s\"", msg)
	c.send <- msg
}

func NewMockChannel(log *log.Logger, network, name string) Channel {
	c := &MockChannel{
		log:       log,
		name:      name,
		network:   network,
		request:   make(chan Events),
		subscribe: make(chan chan Notification, 1),
		send:      make(chan string),
	}

	go func() {
		events := []Event{}
		var epoch int
		var seq uint
		next := make(chan Notification, 1)

		for {
			select {
			case c.request <- EventList(events):
				// handled state request
			case c.subscribe <- next:
				// handled subscribe request.
			case msg := <-c.send:
				c.log.Printf("sending new message: \"%s\"", msg)
				// Add to buffer
				event := Event{
					EventID: EventID{
						Epoch: epoch,
						Seq:   seq,
					},
					Contents: msg,
				}
				events = append(events, event)
				seq++
				// And notify
				notice := Notification{
					Latest: event,
					Next:   make(chan Notification, 1),
				}
				next <- notice
				next = notice.Next
			}
		}
	}()

	return c
}

// MessageGenerator sends message to a Channel.
func MessageGenerator(logger *log.Logger, max uint, c Channel) {
	go func() {
		logger.Print("Chat/messages: [start] counting bottles")
		defer logger.Print("Chat/messages: [done] counting bottles")
		for i := max; i >= 0; i-- {
			time.Sleep(time.Millisecond * 500)

			msg := fmt.Sprintf("%d bottles of beer on the wall, %d bottles of beer...", i, i)
			logger.Print("Chat/messages: [sending] : ", msg)
			c.Send(msg)
		}
	}()
}
