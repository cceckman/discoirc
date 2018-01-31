package demo

import (
	"sync"
	"sync/atomic"

	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/data"
)

// Subscribe attaches the receiver.
func (d *Demo) Subscribe(recv backend.Receiver) {
	d.Lock()
	defer d.Unlock()

	d.subscriber = recv

	// Release the lock before running an update.
	go d.updateAll()
}

func (d *Demo) updateAll() {
	// Assign a unique sequence to each push, since this isn't actually
	// tracking logs.
	seq := data.Seq(atomic.AddInt64(&d.seq, 1))

	// UpdateNetwork and UpdateChannels are presumed synchronous, and may
	// need the RLock back (if they call to EventsBefore). We don't want to
	// hold any locks while we have them.

	// We'll issue zero or more updates as goroutines during this run.
	// To maintain that updateAll is synchronous - completes all its work
	// before returning - we want to put this last, after any locks we may
	// hold are released.
	// in the defer stack so they complete before it returns.
	var wg sync.WaitGroup
	defer wg.Wait()

	d.RLock()
	defer d.RUnlock()

	if d.subscriber == nil {
		return
	}

	recv := d.subscriber
	filter := recv.Filter()

	// Walk through everything; skip if it doesn't match the scope.
	for scope, v := range d.nets {
		// Filter doesn't entirely express what's of interest to the
		// channel view; a real backend has to do some amount of
		// duplication to the channel.
		/// Do a more specific match here.
		if filter.MatchNet && scope.Net != filter.Net {
			continue
		}

		event := &data.NetworkStateEvent{
			EventID: data.EventID{
				Scope: scope,
				Seq:   seq,
			},
			NetworkState: *v,
		}

		wg.Add(1)
		// TODO and WARNING: Running this in a goroutine, without blocking
		// future updates (e.g. ticks) on it, means that events can be
		///sent out-of-order with respect to their seq IDs.
		// I'm not fixing this at the moment,
		// but a real backend should be stricter about sending updates
		// forward in order.
		go func() {
			recv.Receive(event)
			wg.Done()
		}()
	}

	for scope, v := range d.chans {
		if !filter.Match(scope) {
			continue
		}

		event := &data.ChannelStateEvent{
			EventID: data.EventID{
				Scope: scope,
				Seq:   seq,
			},
			ChannelState: *v,
		}

		wg.Add(1)
		// TODO and WARNING: Running this in a goroutine, without blocking
		// future updates (e.g. ticks) on it, means that events can be
		///sent out-of-order with respect to their seq IDs.
		// I'm not fixing this at the moment,
		// but a real backend should be stricter about sending updates
		// forward in order.
		go func() {
			recv.Receive(event)
			wg.Done()
		}()
	}
}
