package data_test

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/cceckman/discoirc/data"
)

var events = []data.Event{
	&data.NetworkStateEvent{
		EventID: data.EventID{Seq: 1},
	},
	&data.NetworkStateEvent{
		EventID: data.EventID{Seq: 2},
	},
	&data.ChannelStateEvent{
		EventID: data.EventID{Seq: 3},
	},
	&data.ChannelStateEvent{
		EventID: data.EventID{Seq: 4},
	},
	&data.NetworkStateEvent{
		EventID: data.EventID{Seq: 5},
	},
}

func TestSelectSizeMax(t *testing.T) {
	eventList := data.EventList(events)
	for _, tt := range []struct {
		name  string
		count int
		max   data.Seq
		want  []data.Event
	}{
		{
			name:  "underflow",
			count: 10000000,
			max:   10000000,
			want:  events,
		},
		{
			name:  "zero size",
			count: 0,
			max:   5,
			want:  []data.Event{},
		},
		{
			name:  "last",
			count: 1,
			max:   math.MaxInt64,
			want:  events[len(events)-1:],
		},
		{
			name:  "precursor",
			count: 2,
			max:   0,
			want:  []data.Event{},
		},
		{
			name:  "mid",
			count: 2,
			max:   3,
			want:  events[1:3],
		},
		{
			name:  "midstart",
			count: 10,
			max:   3,
			want:  events[0:3],
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := eventList.SelectSizeMax(tt.count, tt.max)
			want := data.EventList(tt.want)
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("contents differ: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestStringify(t *testing.T) {
	// This is really just to satisfy the coverage counters.
	hello := "Hello"
	for i, ev := range []data.Event{
		&data.NetworkStateEvent{
			Line: hello,
		},
		&data.ChannelStateEvent{
			Line: hello,
		},
	} {
		got := ev.String()
		if diff := cmp.Diff(got, hello); diff != "" {
			t.Errorf("case %d: contents differ: (-got +want)\n%s", i, diff)
		}
	}
}
