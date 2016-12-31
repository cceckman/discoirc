//
// bufchan.go
// Copyright (C) 2016 cceckman <charles@cceckman.com>
//
// Distributed under terms of the MIT license.
//

// Package bufchan provides channels with unlimited buffering.
package bufchan

import (
	"context"
)

// Bufchan is a channel of strings with an unlimited internal buffer.
// In cases where a UI must interact with a background process (e.g. the network),
// it is necessary that writes between the UI and the background are non-blocking.
// Channels - buffered channels in particular - can help alleviate this; however, channels in Go
// still have a limited size, after which time writes will be blocking operations.
// Bufchan allows channel writes to happen in an approximately always non-blocking manner.
type Bufchan struct {
	in, out chan string
	buf     []string
}

func (x *Bufchan) In() chan<- string {
	return x.in
}

func (x *Bufchan) Out() <-chan string {
	return x.out
}

// mirror reads input as available, and flushes the buffer as the output is available.
// When cancelled (via Context), it closes the output but does not flush it.
// It should be invoked as a goroutine (e.g. go foo.mirror(ctx))
func (x *Bufchan) mirror(ctx context.Context) {
	defer close(x.out)
	for {
		if len(x.buf) == 0 {
			// Select on only "input" and "cancelled".
			select {
			case <-ctx.Done():
				return
			case s := <-x.in:
				x.buf = append(x.buf, s)
			}
		} else {
			// Select on input, output, and cancelled.
			select {
			case <-ctx.Done():
				return
			case s := <-x.in:
				x.buf = append(x.buf, s)
			case out <- x.buf[0]:
				x.buf = x.buf[1:]
			}
		}
	}
}

func New(ctx context.Context) *Bufchan {
	r := &Bufchan{
		in:  make(chan string),
		out: make(chan string),
		buf: make([]string),
	}
	go r.mirror(ctx)
}
