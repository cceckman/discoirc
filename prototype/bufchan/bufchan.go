//
// bufchan.go
// Copyright (C) 2016 cceckman <charles@cceckman.com>
//
// Distributed under terms of the MIT license.
//

// Package bufchan provides channels with unlimited buffering.
package bufchan

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
// When its input channel is closed, it closes the output channel once the buffer is drained.
// It should be invoked as a goroutine (e.g. go foo.mirror())
func (x *Bufchan) mirror() {
	defer close(x.out)

	inputLoop:
	for {
		if len(x.buf) == 0 {
			s, ok := <-x.in
			if !ok {
				// Channel closed.
				break inputLoop
			}
			x.buf = append(x.buf, s)
		} else {
			select {
			case s, ok := <-x.in:
				if !ok {
					break inputLoop
				}
				x.buf = append(x.buf, s)
			case x.out <- x.buf[0]:
				x.buf = x.buf[1:]
			}
		}
	}

	// Flush output.
	for _, p := range x.buf {
		x.out <- p
	}
}

func New() *Bufchan {
	r := &Bufchan{
		in:  make(chan string),
		out: make(chan string),
		buf: make([]string, 0),
	}
	go r.mirror()
	return r
}
