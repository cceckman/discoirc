package mocks

import (
	"context"
	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
	"sync"
)

var _ channel.Model = &Model{}

// Model implements channel.Model for testing.
// It has some restrictions - in particular, it can only support one reader for Follow or Channel.
type Model struct {
	Received []string

	Meta chan data.Channel

	sync.RWMutex
	events   data.EventList
	newEvent chan struct{}
}

func NewModel() *Model {
	r := &Model{
		Meta:     make(chan data.Channel, 1),
		newEvent: make(chan struct{}, 1),
	}

	return r
}

// Follow picks up the most recent data.Event when one is added.
// mock.Model only supports one Follow()ing thread.
func (m *Model) Follow(ctx context.Context) <-chan data.Event {
	result := make(chan data.Event)
	go func() {
		defer close(result)
		for {
			select {
			case <-ctx.Done():
				return
			case <-m.newEvent:
				m.RLock()
				next := m.events[len(m.events)-1]
				m.RUnlock()
				result <- next
			}
		}
	}()
	return result
}

// EventsEndingAt satisfies the channel.Model interface.
func (m *Model) EventsEndingAt(end data.EventID, n int) []data.Event {
	m.RLock()
	defer m.RUnlock()

	return m.events.SelectSizeMax(uint(n), end)
}

// AddEvent adds an event to the history of this Model.
// It triggers any Follow()ers.
func (m *Model) AddEvent(e data.Event) {
	m.Lock()
	defer m.Unlock()
	// keep sorted.
	m.events = data.NewEvents(append(m.events, e))
	go func() {
		select {
		case <-m.newEvent:
			// Clear any waiting trigger
			m.newEvent <- struct{}{}
		case m.newEvent <- struct{}{}:
			// trigger, if nothing else is waiting
		}
	}()
}

// Channel reflects updates to the channel's state.
// It does NOT support multiple listeners / multiple calls.
func (m *Model) Channel(ctx context.Context) <-chan data.Channel {
	result := make(chan data.Channel)
	go func() {
		defer close(result)
		for {
			select {
			case <-ctx.Done():
				return
			case update := <-m.Meta:
				result <- update
			}
		}
	}()
	return result
}

func (m *Model) Send(s string) error {
	m.Received = append(m.Received, s)
	// TODO support testing with an error returned
	return nil
}
