// Package model provides models for chat and connection updates.
package model

import (
	"context"
	"log"
)

type Channel interface {
	Events

	Name() string
	Network() string

	// SendMessage sends a message to the channel.
	SendMessage(string)
	// Events reports the logged/timestamped Events that come in after the
	// listener begins.
	Events(context.Context) <-chan Event

	// State represents update to the current state of the channel: changes
	// to the user's nick and mode, the channel topic, etc.
	// These may also come in as Events.
	// The returned channel will return an initial state.
	State(context.Context) <-chan ChannelState
}

type ChannelState struct {
	Connected         bool
	Nick, Mode, Topic string
	next              chan ChannelState
}

type Notification struct {
	Latest Event
	next   chan Notification
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

	eventSubscribe chan chan Notification
	send           chan string

	updateState chan ChannelState
	connected   chan bool

	stateSubscribe chan chan ChannelState
	channelState   chan ChannelState
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

func (c *MockChannel) Events(ctx context.Context) <-chan Event {
	c.log.Printf("added listener for events in channel %s / %s", c.Network(), c.Name())
	result := make(chan Event)
	go func() {
		// Block on subscription request.
		var notices chan Notification
		select {
		case <-ctx.Done():
			return
		case notices = <-c.eventSubscribe:
			// Have a channel to listen on.
		}

		for {
			select {
			case <-ctx.Done():
				return
			case notice := <-notices:
				// Put it back in the broadcast channel immediately.
				notices <- notice
				// Await the next notice
				notices = notice.next
				// And forward the semantic content
				result <- notice.Latest
			}
		}
	}()
	return result
}

func (c *MockChannel) SendMessage(msg string) {
	c.log.Printf("awaiting send for message \"%s\"", msg)
	c.send <- msg
}

// eventLoop is the handler for messages / events.
func (c *MockChannel) eventLoop() {
	events := []Event{}
	var epoch int
	var seq uint
	next := make(chan Notification, 1)

	for {
		select {
		case c.request <- EventList(events):
			// handled state request
		case c.eventSubscribe <- next:
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
				next:   make(chan Notification, 1),
			}
			next <- notice
			next = notice.next
		}
	}

}

// State listens for updates to the channel's state.
func (c *MockChannel) State(ctx context.Context) <-chan ChannelState {
	r := make(chan ChannelState, 1)
	go func() {
		defer close(r)

		handle := <-c.stateSubscribe

		for {
			// Do initial read
			select {
			case <-ctx.Done():
				return
			case newState := <-handle:
				// Made a local copy of the state.
				// Put it back for the next listener, update our local handle.
				handle <- newState
				handle = newState.next
				// Notify client with lossy write to the channel.
				select {
				case r <- newState:
					// have enqueued.
				case _ = <-r:
					// dequeued old state; enqueue the new one.
					r <- newState
				}
			}
		}
	}()

	return r
}

// stateLoop is the handler for updates to the channel state.
func (c *MockChannel) stateLoop() {
	var state = ChannelState{
		next: make(chan ChannelState, 1),
	}

	last := make(chan ChannelState, 1)
	last <- state

	for {
		select {
		case c.stateSubscribe <- last:
			// Sent handler to last state to the new listener.
			continue
			// It'll pick up the next field and continue the stream from there.
		case con := <-c.connected:
			state.Connected = con
		case update := <-c.updateState:
			if update.Nick != "" {
				state.Nick = update.Nick
			}
			if update.Mode != "" {
				state.Mode = update.Mode
			}
			if update.Topic != "" {
				state.Topic = update.Topic
			}
		}
		// Broadcast an update.
		last = state.next
		state.next = make(chan ChannelState, 1)
		last <- state
	}
}

func NewMockChannel(log *log.Logger, network, name string) *MockChannel {
	c := &MockChannel{
		log:            log,
		name:           name,
		network:        network,
		request:        make(chan Events),
		eventSubscribe: make(chan chan Notification, 1),
		send:           make(chan string),

		updateState: make(chan ChannelState),
		connected:   make(chan bool),

		stateSubscribe: make(chan chan ChannelState, 1),
	}

	// Event loop.
	go c.eventLoop()
	go c.stateLoop()

	return c
}
