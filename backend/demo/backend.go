// Package demo implements the discoirc non-UI portions with demo data.
package demo

import (
	"fmt"
	"sort"
	"sync"

	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/data"
)

var _ backend.Backend = &Demo{}

type message string

func (m message) String() string    { return string(m) }

// Demo provides data and updates to discoirc UI components.
type Demo struct {
	sync.RWMutex

	subscriber backend.StateReceiver

	nets     map[string]*data.NetworkState
	chans    map[data.Scope]*data.ChannelState
	contents map[data.Scope]data.EventList
}

// New returns a new demonstration backend
func New() *Demo {
	d := &Demo{
		nets:     make(map[string]*data.NetworkState),
		chans:    make(map[data.Scope]*data.ChannelState),
		contents: make(map[data.Scope]data.EventList),
	}
	return d
}

// Send sends the given message to the target.
func (d *Demo) Send(scope data.Scope, message string) {
	d.ensureChannel(scope.Net, scope.Name)

	d.Lock()
	defer d.Unlock()
	nick := d.nets[scope.Net].Nick
	d.appendMessage(scope, nick, message)
}

// appendMessage must be called under the write lock.
func (d *Demo) appendMessage(id data.Scope, speaker, msg string) {
	last := d.chans[id].LastMessage
	next := data.Event{
		Scope: id,
		Seq:   last + 1,
		EventContents: message(fmt.Sprintf(
			"<%s> %s",
			speaker, msg,
		)),
	}

	// Doesn't update unread; 'send' doesn't count as unread.
	d.contents[id] = append(d.contents[id], next)
	d.chans[id].LastMessage = next.Seq

	go d.updateAll()
}

// EventsBefore returns N events preceding the given event in the given channel.
func (d *Demo) EventsBefore(id data.Scope, n int, last data.Seq) data.EventList {
	d.RLock()
	defer d.RUnlock()

	evs := d.contents[id]

	v := evs.SelectSizeMax(n, last)

	// Update unread; How many messages have been read, as of this one?
	readToIdx := sort.Search(len(evs), func(i int) bool {
		return evs[i].Seq >= last
	})
	// Use a separate thread to update the number unread-
	// it requires the write-lock, and we don't want to block on that.
	go func() {
		d.Lock()
		defer d.Unlock()

		ch, ok := d.chans[id]
		if !ok {
			// No channel metadata to update;
			// just return.
			return
		}

		// We know at least 'read' messages have been read, but
		// it's possible that more have been read before this goroutine
		// runs.
		// Recompute, and take the min- newly-arriving messages increase
		// it, this should only decrease.
		unread := len(d.contents[id]) - (readToIdx + 1)
		// +1 accounts for array of size 1- we've read through index 0.
		if unread < 0 {
			unread = 0
		}
		if unread < ch.Unread {
			ch.Unread = unread
		}
		go d.updateAll()
	}()
	// TODO: Handle num-unread better in a non-demo backend.

	return v
}
