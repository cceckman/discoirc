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

type message struct {
	data.EventID
	message string
}

var _ data.Event = &message{}

func (m *message) String() string    { return m.message }
func (m *message) ID() *data.EventID { return &m.EventID }

// Demo provides data and updates to discoirc UI components.
type Demo struct {
	sync.RWMutex

	subscriber backend.Receiver

	nets     map[data.Scope]*data.NetworkState
	chans    map[data.Scope]*data.ChannelState
	contents map[data.Scope]data.EventList

	seq int64
}

// New returns a new demonstration backend
func New() *Demo {
	d := &Demo{
		nets:     make(map[data.Scope]*data.NetworkState),
		chans:    make(map[data.Scope]*data.ChannelState),
		contents: make(map[data.Scope]data.EventList),
	}
	return d
}

// Send sends the given message to the target.
func (d *Demo) Send(scope data.Scope, message string) {
	d.ensureChannel(scope)

	d.Lock()
	defer d.Unlock()
	netScope := data.Scope{ Net: scope.Net }

	nick := d.nets[netScope].Nick
	d.appendMessage(scope, nick, message)
}

// appendMessage must be called under the write lock.
func (d *Demo) appendMessage(scope data.Scope, speaker, contents string) {
	last := d.chans[scope].LastMessage
	next := &message{
		EventID: data.EventID{
			Scope: scope,
			Seq:   last + 1,
		},
		message: fmt.Sprintf(
			"<%s> %s",
			speaker, contents,
		),
	}

	// Doesn't update unread; 'send' doesn't count as unread.
	d.contents[scope] = append(d.contents[scope], next)
	d.chans[scope].LastMessage = next.ID().Seq

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
		return evs[i].ID().Seq >= last
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
