//
// stream.go
// Copyright (C) 2017 cceckman <charles@cceckman.com>
//
// Distributed under terms of the MIT license.
//

package stream

// Aliases, for brevity.
type Client EventProviderClient
type Server EventProviderServer

// Exec determines if the Event matches anything in this Filter.
func (f *Filter) Exec(e *Event) bool {
	for _, m := range f.Matches {
		if m.Exec(e) {
			return true
		}
	}
	return false
}

// Exec determines if the given Event matches this Match.
func (f *Match) Exec(e *Event) bool {
	if f.MatchPlugin && f.Id.Plugin != e.Stream.Plugin {
		return false
	}
	if f.MatchNetwork && f.Id.Network != e.Stream.Network {
		return false
	}
	if f.MatchChannel && f.Id.Channel != e.Stream.Channel {
		return false
	}

	return true
}
