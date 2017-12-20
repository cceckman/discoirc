package channel

import (
	"github.com/cceckman/discoirc/data"
)

var evs data.EventList

func init() {
	evs = NewEvents([]Event{
		Event{EventId{Epoch: 1, Seq: 1}, Contents: "TOPIC Act I, Scene 1"},
		Event{EventId{Epoch: 1, Seq: 2}, Contents: "JOIN barnardo"},
		Event{EventId{Epoch: 1, Seq: 3}, Contents: "JOIN francisco"},
		Event{EventId{Epoch: 1, Seq: 4}, Contents: "<barnardo> Who's there?"},
		Event{EventId{Epoch: 1, Seq: 5}, Contents: "<francisco> Nay answer me: Stand & vnfold your selfe"},
		Event{EventId{Epoch: 1, Seq: 6}, Contents: "<barnardo> Long liue the King"},
			Event{Epoch: 2, Seq: 1}, Contents: "<claudius> Welcome, dear Rosencrantz and Guildenstern!"},
		Event{EventId{Epoch: 2, Seq: 2}, Contents: "<gertrude> Good gentlemen, he hath much talk'd of you;"},
		Event{EventId{Epoch: 2, Seq: 3}, Contents: "<rosencrantz> Both your majesties"},
	})
}
