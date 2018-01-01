package mocks

import (
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

type target struct {
	Net string
	Tgt string
}

type Backend struct {
	Receiver backend.StateReceiver

	events map[target]data.EventList

	Sent []string
}

func (b *Backend) Subscribe(r backend.StateReceiver) {
	b.Receiver = r
}
func (b *Backend) SubscribeFiltered(r backend.FilteredStateReceiver) {
	b.Receiver = r
}

func (b *Backend) EventsBefore(network, tgt string, n int, last data.EventID) []data.Event {
	if v, ok := b.events[target{
		Net: network,
		Tgt: tgt,
	}]; ok {
		return v.SelectSizeMax(uint(n), last)
	}

	return nil
}

func (b *Backend) Send(_, _ string, message string) {
	b.Sent = append(b.Sent, message)
}

func NewBackend() *Backend {
	contents := map[target]data.EventList{
		target{
			Net: "HamNet",
			Tgt: "#hamlet",
		}: Events,
	}
	return &Backend{
		events: contents,
	}

}
