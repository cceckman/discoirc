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

type ChanIdent struct {
	Network, Channel string
}

// Demo provides data and updates to discoirc UI components.
type Demo struct {
	sync.RWMutex

	subscriber backend.StateReceiver
	filter     func() (string, string)

	nets     map[string]*data.NetworkState
	chans    map[ChanIdent]*data.ChannelState
	contents map[ChanIdent][]data.Event
}

func New() *Demo {
	d := &Demo{
		nets:     make(map[string]*data.NetworkState),
		chans:    make(map[ChanIdent]*data.ChannelState),
		contents: make(map[ChanIdent][]data.Event),
	}
	return d
}

func (d *Demo) Send(network, channel string, message string) {
	d.ensureChannel(network, channel)

	d.Lock()
	defer d.Unlock()
	nick := d.nets[network].Nick
	d.appendMessage(network, channel, nick, message)
}

// appendMessage must be called under the write lock.
func (d *Demo) appendMessage(network, channel string, speaker, message string) {
	id := ChanIdent{
		Network: network,
		Channel: channel,
	}
	last := d.chans[id].LastMessage.EventID
	var next data.Event

	// Construct a new message.
	next.Seq = last.Seq + 1
	next.Epoch = last.Epoch
	next.Contents = fmt.Sprintf(
		"<%s> %s",
		speaker, message,
	)

	// Doesn't update unread; 'send' doesn't count as unread.
	d.contents[id] = append(d.contents[id], next)
	d.chans[id].LastMessage = next

	go d.updateAll()
}

func (d *Demo) EventsBefore(network, channel string, n int, last data.EventID) []data.Event {
	d.RLock()
	defer d.RUnlock()

	id := ChanIdent{
		Network: network,
		Channel: channel,
	}

	evs := data.NewEvents(d.contents[id])

	v := evs.SelectSizeMax(uint(n), last)

	// Update unread; How many messages have been read, as of this one?
	readToIdx := sort.Search(len(evs), func(i int) bool {
		if evs[i].Epoch == last.Epoch {
			return evs[i].Seq >= last.Seq
		}
		return evs[i].Epoch > last.Epoch
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
