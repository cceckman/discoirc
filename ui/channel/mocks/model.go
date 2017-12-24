package mocks

import (
	"sync"

	"github.com/cceckman/discoirc/data"
	"github.com/cceckman/discoirc/ui/channel"
)

var _ channel.Model = &Model{}

// Model implements channel.Model for testing.
// It has some restrictions - in particular, it can only support one reader for Follow or Channel.
type Model struct {
	Received []string

	sync.RWMutex
	Channel    data.Channel
	Events     data.EventList
	Controller channel.ModelController
}

// EventsEndingAt satisfies the channel.Model interface.
func (m *Model) EventsEndingAt(end data.EventID, n int) []data.Event {
	m.RLock()
	defer m.RUnlock()
	return m.Events.SelectSizeMax(uint(n), end)
}

// AddEvent adds an event to the history of this Model.
// It triggers any Follow()ers.
func (m *Model) AddEvent(s string) {
	m.Lock()
	defer m.Unlock()

	var eid data.EventID
	if len(m.Events) > 0 {
		lastEid := m.Events[len(m.Events)-1]
		eid.Epoch = lastEid.Epoch + 1
	}

	// keep sorted.
	m.Events = data.NewEvents(append(m.Events, data.Event{
		EventID:  eid,
		Contents: s,
	}))

	if m.Controller != nil {
		m.Controller.UpdateContents(m.Events[len(m.Events)-1])
	}
}

func (m *Model) UpdateChannel(d data.Channel) {
	m.Lock()
	defer m.Unlock()
	m.Channel = d

	if m.Controller != nil {
		m.Controller.UpdateMeta(m.Channel)
	}
}

func (m *Model) Attach(c channel.ModelController) {
	m.Lock()
	defer m.Unlock()
	m.Controller = c

	if len(m.Events) > 0 {
		m.Controller.UpdateContents(m.Events[len(m.Events)-1])
	}

	m.Controller.UpdateMeta(m.Channel)
}

func (m *Model) Send(s string) error {
	m.Received = append(m.Received, s)
	// TODO support testing with an error returned
	// TODO support reflecting back into events
	return nil
}
