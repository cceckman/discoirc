package mocks

import (
	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/data"
)

type message string

func (m message) String() string { return string(m) }

// Events is a set of data.Events used by tests as filler data - a Lorem.
// Specifically, it's a few lines and stage directions from the first two scenes
// of Shakespeare's Hamlet.
var Events data.EventList

func init() {
	d := []string{
		"TOPIC Act I, Scene 1",
		"JOIN barnardo",
		"JOIN francisco",
		"<barnardo> Who's there?",
		"<francisco> Nay answer me: Stand & vnfold your selfe",
		"<barnardo> Long liue the King",
		"<claudius> Welcome, dear Rosencrantz and Guildenstern!",
		"<gertrude> Good gentlemen, he hath much talk'd of you;",
		"<rosencrantz> Both your majesties",
	}
	es := make([]data.Event, len(d))
	for i, v := range d {
		es[i] = data.Event{
			Scope:    data.Scope{Net: "Shaxnet", Name: "#hamlet"},
			Seq:      data.Seq(i),
			EventContents: message(v),
		}
	}
	Events = data.SortEvents(es)
}

// Backend is a mock implementor of the backend.Backend interface.
type Backend struct {
	Receiver backend.StateReceiver

	events data.EventList

	Sent []string
}

// Subscribe implements backend.Backend
func (b *Backend) Subscribe(r backend.StateReceiver) {
	b.Receiver = r
}

// EventsBefore implements backend.Backend
func (b *Backend) EventsBefore(s data.Scope, n int, last data.Seq) data.EventList {
	return b.events.SelectSizeMax(n, last)
}

// Send implements backend.Backend
func (b *Backend) Send(_ data.Scope, message string) {
	b.Sent = append(b.Sent, message)
}

// NewBackend returns a new, mock, Backend
func NewBackend() *Backend {
	return &Backend{
		events: Events,
	}
}
