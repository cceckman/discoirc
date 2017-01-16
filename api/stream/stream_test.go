// 2017-01-15 cceckman <charles@cceckman.com>
package stream_test

import(
	"github.com/cceckman/discoirc/api/stream"
	"testing"
)

// filters are filters that events may or may not match against.
var filters = map[string]&stream.Filter{
	"foonetic": &stream.Filter{
		Matches: []*stream.Match{
			&stream.Match{
				Id: &stream.Id {
					Network: "Foonetic",
				},
				MatchNetwork: true,
			},
		},
	},
}

// events are expected to match or not match certain filters.
// This is the main test table.
var events = map[string]struct{
	// The event to match
	event *stream.Event
	// for the given filter in filters, whether this event is expected to match or not
	want map[string]bool
} {

}


func TestFilter(t *testing.T) {
	for _, event := range events {

	}
}
