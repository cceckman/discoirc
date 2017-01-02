// Provides for a set of buffered output channels on an input.
package bufchan

import (
	"context"
	"sync"
)

// A Broadcaster listens on a string channel, and writes to all of several listeners.
// It is non-blocking; slower Listen()ers do not block faster ones.
type Broadcaster interface {
	Listen(context.Context) <-chan string
}

// Creates a new Broadcaster on the given input channel.
func NewBroadcaster(c <-chan string) Broadcaster {
	r := &listBroadcaster{input: c}
	go r.broadcast()
	return r
}

// listBroadcaster implements Broadcaster as a list of Bufchans.
type listBroadcaster struct {
	input <-chan string

	receivers []*Bufchan
	sync.RWMutex
}

// broadcast listens on input, and writes to receivers.
func (lb *listBroadcaster) broadcast() {
	for {
		x, ok := <-lb.input
		lb.RLock()
		// Don't defer; explicitly unlock after loop over receivers.
		for _, receiver := range lb.receivers {
			if ok {
				receiver.In() <- x
			} else {
				close(receiver.In())
			}
		}
		lb.RUnlock()

		if !ok {
			// Input is closed, all receivers closed in above. All done!
			break
		}
	}
}

func (lb *listBroadcaster) Listen(ctx context.Context) <-chan string {
	lb.Lock()
	defer lb.Unlock()

	n := New(ctx)
	lb.receivers = append(lb.receivers, n)
	go func() {
		// Remove from receivers once the context is cancelled.
		<-ctx.Done()

		lb.Lock()
		defer lb.Unlock()
		for i, r := range lb.receivers {
			if r == n {
				// Nice pointer: https://github.com/golang/go/wiki/SliceTricks
				// "delete without preserving order". That's fine.
				newLen := len(lb.receivers) - 1
				lb.receivers[i] = lb.receivers[newLen]
				lb.receivers[newLen] = nil // Allow GC of the *Bufchan (once other threads complete too.)
				lb.receivers = lb.receivers[:newLen]
				break
			}
		}
	}()

	return n.Out()
}
