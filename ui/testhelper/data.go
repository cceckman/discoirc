package testhelper

import (
	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/data"
)

type event struct {
	seq      data.Seq
	Contents string
}

func (e *event) String() string    { return e.Contents }
func (e *event) Scope() data.Scope { return data.Scope{Net: "Shaxnet", Name: "#hamlet"} }
func (e *event) Seq() data.Seq     { return e.seq }

// Events is a set of data.Events used by tests as filler data - a Lorem.
// Specifically, it's a few lines and stage directions from the first two scenes
// of Shakespeare's Hamlet.
var Events data.EventList

func init() {
	d := []event{
		{data.Seq(1), "TOPIC Act I, Scene 1"},
		{data.Seq(2), "JOIN barnardo"},
		{data.Seq(3), "JOIN francisco"},
		{data.Seq(4), "<barnardo> Who's there?"},
		{data.Seq(5), "<francisco> Nay answer me: Stand & vnfold your selfe"},
		{data.Seq(6), "<barnardo> Long liue the King"},
		{data.Seq(7), "<claudius> Welcome, dear Rosencrantz and Guildenstern!"},
		{data.Seq(8), "<gertrude> Good gentlemen, he hath much talk'd of you;"},
		{data.Seq(9), "<rosencrantz> Both your majesties"},
	}
	es := make([]data.Event, len(d))
	for i, v := range d {
		v := v
		es[i] = data.Event(&v)
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
