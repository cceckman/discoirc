package demo

import (
	"sync"

	"github.com/cceckman/discoirc/backend"
)

func (d *Demo) Subscribe(recv backend.StateReceiver) {
	d.subscribe(recv, nil)
}

func (d *Demo) SubscribeFiltered(recv backend.FilteredStateReceiver) {
	d.subscribe(recv, recv.Filter)
}

func (d *Demo) subscribe(recv backend.StateReceiver, filter func() (string, string)) {
	d.Lock()
	defer d.Unlock()

	d.subscriber = recv
	d.filter = filter

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

	recv := d.subscriber

	if recv == nil {
		// Nothing to receive our updates.
		return
	}

	if d.filter != nil {
		net, ch := d.filter()
		netState, ok := d.nets[net]
		netV := *netState // pass by value
		if ok {
			wg.Add(1)
			go func() {
				recv.UpdateNetwork(netV)
				wg.Done()
			}()
		}

		chId := chanIdent{
			Network: net,
			Channel: ch,
		}

		tgtState, ok := d.chans[chId]
		tgtV := *tgtState
		if ok {
			wg.Add(1)
			go func() {
				recv.UpdateChannel(tgtV)
				wg.Done()
			}()
		}

		return
	}

	// No filter; update everything.
	for _, v := range d.nets {
		v := *v
		wg.Add(1)
		go func() {
			recv.UpdateNetwork(v)
			wg.Done()
		}()
	}

	for _, v := range d.chans {
		v := *v
		wg.Add(1)
		go func() {
			recv.UpdateChannel(v)
			wg.Done()
		}()
	}
}
