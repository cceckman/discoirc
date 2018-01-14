package demo

import (
	"sync"

	"github.com/cceckman/discoirc/backend"
)

// Subscribe attaches the receiver.
func (d *Demo) Subscribe(recv backend.StateReceiver) {
	d.Lock()
	defer d.Unlock()

	d.subscriber = recv

	// Release the lock before running an update.
	go d.updateAll()
}

func (d *Demo) updateAll() {
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
	for _, v := range d.nets {
		v := *v
		// Filter doesn't entirely express what's of interest to the
		// channel view; a real backend has to do some amount of
		// duplication to the channel.
		/// Do a more specific match here.
		if filter.MatchNet && v.Scope.Net != filter.Net {
			continue
		}

		wg.Add(1)
		go func() {
			recv.UpdateNetwork(v)
			wg.Done()
		}()
	}

	for _, v := range d.chans {
		v := *v
		if !filter.Match(v.Scope) {
			continue
		}

		wg.Add(1)
		go func() {
			recv.UpdateChannel(v)
			wg.Done()
		}()
	}
}
