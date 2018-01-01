package mocks

import (
	"context"
	"fmt"

	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/data"
)

var Events = data.NewEvents([]data.Event{
	data.Event{EventID: data.EventID{Epoch: 1, Seq: 1}, Contents: "TOPIC Act I, Scene 1"},
	data.Event{EventID: data.EventID{Epoch: 1, Seq: 2}, Contents: "JOIN barnardo"},
	data.Event{EventID: data.EventID{Epoch: 1, Seq: 3}, Contents: "JOIN francisco"},
	data.Event{EventID: data.EventID{Epoch: 1, Seq: 4}, Contents: "<barnardo> Who's there?"},
	data.Event{EventID: data.EventID{Epoch: 1, Seq: 5}, Contents: "<francisco> Nay answer me: Stand & vnfold your selfe"},
	data.Event{EventID: data.EventID{Epoch: 1, Seq: 6}, Contents: "<barnardo> Long liue the King"},
	data.Event{EventID: data.EventID{Epoch: 2, Seq: 1}, Contents: "<claudius> Welcome, dear Rosencrantz and Guildenstern!"},
	data.Event{EventID: data.EventID{Epoch: 2, Seq: 2}, Contents: "<gertrude> Good gentlemen, he hath much talk'd of you;"},
	data.Event{EventID: data.EventID{Epoch: 2, Seq: 3}, Contents: "<rosencrantz> Both your majesties"},
})

type Backend struct {
	Receiver backend.FilteredStateReceiver
	Events   data.EventList

	Sent []string
}

func (b *Backend) Subscribe(_ context.Context, _ backend.StateReceiver) {
	panic(fmt.Errorf("not implemented"))
}
func (b *Backend) SubscribeFiltered(r backend.FilteredStateReceiver) {
	b.Receiver = r
}

func (b *Backend) EventsBefore(n int, last data.EventID) []data.Event {
	return b.Events.SelectSizeMax(uint(n), last)
}

func (b *Backend) Send(_, _ string, message string) {
	b.Sent = append(b.Sent, message)
}

func NewBackend() *Backend {
	return &Backend{
		Events: Events,
	}

}
