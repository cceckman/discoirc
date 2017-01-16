// 2017-01-15 user1 <charles@user1.com>
package stream_test

import (
	"github.com/cceckman/discoirc/api/stream"
	"testing"
)

// filters are filters that events may or may not match against.
var filters = map[string]*stream.Filter{
	"Nothing": &stream.Filter{},
	"Network2": &stream.Filter{
		Matches: []*stream.Match{
			&stream.Match{
				Id: &stream.Id{
					Network: "Network2",
					Channel: "#user2",
				},
				MatchNetwork: true,
				// Ignore specified channel; catch everything.
			},
		},
	},
	"Network1": &stream.Filter{
		Matches: []*stream.Match{
			&stream.Match{
				Id: &stream.Id{
					Network: "Network1",
				},
				MatchNetwork: true,
			},
		},
	},
	"Network1#channel1": &stream.Filter{
		Matches: []*stream.Match{
			&stream.Match{
				Id: &stream.Id{
					Network: "Network1",
					Channel: "#channel1",
				},
				MatchNetwork: true,
				MatchChannel: true,
			},
		},
	},
	"plugin": &stream.Filter{
		Matches: []*stream.Match{
			&stream.Match{
				Id: &stream.Id{
					Plugin:  "plugin1",
					Network: "Network1",
					Channel: "#channel1",
				},
				MatchPlugin: true,
			},
		},
	},
}

// events are expected to match or not match certain filters.
// This is the main test table.
var events = []struct {
	// The event to match
	Event *stream.Event
	// For the given filter in filters, whether this event is expected to match or not.
	Want map[string]bool
}{
	// Case 1: a message on Network1#channel1.
	{
		&stream.Event{
			Stream: &stream.Id{
				Plugin:  "system",
				Network: "Network1",
				Channel: "#channel1",
			},
			Text: "user1 joined #channel1",
		},
		map[string]bool{
			"Network2":          false,
			"Network1":          true,
			"Network1#channel1": true,
			"plugin":            false,
			"Nothing":           false,
		},
	},
	// Case 2: a message on Network1#channel2.
	{
		&stream.Event{
			Stream: &stream.Id{
				Plugin:  "system",
				Network: "Network1",
				Channel: "#channel2",
			},
			Text: "user1 joined #channel2",
		},
		map[string]bool{
			"Network2":          false,
			"Network1":          true,
			"Network1#channel1": false,
			"plugin":            false,
			"Nothing":           false,
		},
	},
	// Case 3: a message on Network2, no channel.
	{
		&stream.Event{
			Stream: &stream.Id{
				Plugin:  "system",
				Network: "Network2",
			},
			Text: "user1 messaged user2",
		},
		map[string]bool{
			"Network2":          true,
			"Network1":          false,
			"Network1#channel1": false,
			"plugin":            false,
			"Nothing":           false,
		},
	},
	// Case 4: a message to a plugin.
	{
		&stream.Event{
			Stream: &stream.Id{
				Plugin: "plugin1",
			},
			Text: "user1 messaged user2",
		},
		map[string]bool{
			"Network2":          false,
			"Network1":          false,
			"Network1#channel1": false,
			"plugin":            true,
			"Nothing":           false,
		},
	},
}

func TestFilter(t *testing.T) {
	for _, cs := range events {
		for name, filter := range filters {
			want, ok := cs.Want[name]
			if !ok {
				t.Errorf("no expectation for filter '%s' on event: %v", name, cs.Event)
				continue
			}
			got := filter.Exec(cs.Event)
			if got != want {
				t.Errorf("for filter '%s' got: %t want: %t for event: %v", name, got, want, cs.Event)
			}
		}
	}
}
