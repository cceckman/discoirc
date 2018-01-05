package demo

import (
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
	d.RLock()
	defer d.RUnlock()

	// No subscriber? update is a noop.
	if d.subscriber == nil {
		return
	}

	// UpdateNetwork call is synchronous, and we're holding an RLock.
	// TODO: do something else to keep things ordered, and do Update in
	// a goroutine.

	if d.filter != nil {
		net, ch := d.filter()
		netState, ok := d.nets[net]
		if ok {
			d.subscriber.UpdateNetwork(*netState)
		}

		chId := ChanIdent{
			Network: net,
			Channel: ch,
		}

		tgtState, ok := d.chans[chId]
		if ok {
			d.subscriber.UpdateChannel(*tgtState)
		}

		return
	}

	// No filter; update everything.
	for _, v := range d.nets {
		d.subscriber.UpdateNetwork(*v)
	}

	for _, v := range d.chans {
		d.subscriber.UpdateChannel(*v)
	}
}
