// Provides for a set of buffered output channels on an input.
package bufchan

import (
	"context"
)

// A Broadcaster listens on a string channel, and writes to all of several listeners.
// It is non-blocking; slower Listen()ers do not block faster ones.
type Broadcaster interface {
	Listen(context.Context) <-chan string
}

// Creates a new Broadcaster on the given input channel.
func NewBroadcaster(c <-chan string) Broadcaster {
	r := &listBroadcaster{
		input:        c,
		rmReceivers:  make(chan *Bufchan),
		addReceivers: make(chan *Bufchan),
	}
	go r.broadcast()
	return r
}

// listBroadcaster implements Broadcaster as a list of Bufchans.
type listBroadcaster struct {
	input <-chan string

	receivers    []*Bufchan
	rmReceivers  chan *Bufchan
	addReceivers chan *Bufchan
}

// broadcast is the main channel of a listBroadcaster; it listens for input or add / remove requests,
// and writes to its receivers.
func (lb *listBroadcaster) broadcast() {
	closed := false
	for !closed {
		select {
		case rm := <-lb.rmReceivers:
			for i, receiver := range lb.receivers {
				if receiver == rm {
					// Nice pointer: https://github.com/golang/go/wiki/SliceTricks
					// "delete without preserving order". That's fine.
					newLen := len(lb.receivers) - 1
					lb.receivers[i] = lb.receivers[newLen]
					lb.receivers[newLen] = nil // Allow GC of the *Bufchan (once other threads complete too.)
					lb.receivers = lb.receivers[:newLen]
					break // from the loop over receivers.
				}
			}
		case add := <-lb.addReceivers:
			lb.receivers = append(lb.receivers, add)
		case x, ok := <-lb.input:
			if !ok {
				closed = true
			} else {
				for _, receiver := range lb.receivers {
					receiver.In() <- x
				}
			}
		}
	}

	// Input channel is closed. Close all current receivers.
	for i, r := range lb.receivers {
		close(r.In())
		lb.receivers[i] = nil // allow GC, when other threads complete.
	}
	lb.receivers = []*Bufchan{}

	// Hang around indefinitely to handle the add / remove channels.
	// TODO: I don't like this. This tail-end leaks the entire listBroadcaster, keeping it around
	// for the lifetime of the program.
	// I think there's a way to do it with a done-channel alongside Input,
	// s.t. Listen can return immediately, but I"m too tired to think of it at the moment; seems like
	// everything's coming up races.
	for {
		select {
		case <-lb.rmReceivers:
			// pass; we've already cleared lb.receivers.
		case add := <-lb.addReceivers:
			// Input is closed, so close the receiver's input too.
			close(add.In())
		}
	}
}

func (lb *listBroadcaster) Listen(ctx context.Context) <-chan string {
	n := New(ctx)
	lb.addReceivers <- n
	go func() {
		// Remove from receivers once the context is cancelled.
		// The case of "if input is closed" is handled by the main broadcast loop.
		<-ctx.Done()
		lb.rmReceivers <- n
	}()

	return n.Out()
}
