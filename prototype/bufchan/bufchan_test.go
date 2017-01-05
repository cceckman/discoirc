// 2016-12-31 cceckman <charles@cceckman.com>
package bufchan_test

import (
	"context"
	"fmt"
	"github.com/cceckman/discoirc/prototype/bufchan"
	"strings"
	"sync"
	"testing"
	"time"
)

// testAtRates generates a test case for correct behavior when reading and writing at the respective rates.
func testAtRates(r, w time.Duration) func(*testing.T) {
	return func(t *testing.T) {
		c := bufchan.New()

		var wg sync.WaitGroup
		wg.Add(2)

		// Start writer...
		go func() {
			ticker := time.NewTicker(w)
			defer ticker.Stop()
			for i := 0; i < 100; i++ {
				tm := <-ticker.C
				// Assert that the channel blocks for at most this amount of time.
				// Should be pretty small.
				timeout := time.After(time.Millisecond)
				select {
				case c.In() <- tm:
					continue
				case <-timeout:
					t.Errorf("iteration %d timed out", i)
				}
			}
			wg.Done()
		}()
		// ...and reader.
		go func() {
			ticker := time.NewTicker(r)
			defer ticker.Stop()

			lastTime := time.Time{}
			var zero time.Duration

			for i := 0; i < 100; i++ {
				<-ticker.C
				// Can't make any assertions about how long it blocks for;
				// may be indefinitely, if the writer is slower than the reader.
				v, ok := <-c.Out()
				if !ok {
					t.Errorf("Output closed when not expected")
				}
				tm, ok := v.(time.Time)
				if !ok {
					t.Errorf("type assertion failed: %v", tm)
					continue
				}
				since := tm.Sub(lastTime)
				if since < w {
					t.Logf("read timestamps were faster than expected: got %s want %s",
						since,
						w,
					)
				}

				if since < zero {
					t.Errorf("went backwards in time: %s to %s (%s)",
						lastTime.Format(time.RFC3339Nano),
						tm.Format(time.RFC3339Nano),
						since,
					)
				}
				lastTime = tm
			}
			wg.Done()
		}()

		// Wait for reader and writer to be done.
		wg.Wait()
	}
}

// Test with various read-and-write rates.
func TestRwRates(t *testing.T) {
	for _, rates := range []struct {
		r, w time.Duration
	}{
		{time.Millisecond, time.Millisecond},
		{time.Microsecond, time.Microsecond},
		{time.Microsecond * 2, time.Microsecond},
		{time.Microsecond * 100, time.Microsecond},
		{time.Microsecond, time.Microsecond * 2},
		{time.Microsecond, time.Microsecond * 100},
	} {
		name := fmt.Sprintf("r=%s/w=%s", rates.r, rates.w)
		t.Run(name, testAtRates(rates.r, rates.w))
	}
}

// Test that closing input closes output, after N items.
func testClose(n int) func(*testing.T) {
	return func(t *testing.T) {
		c := bufchan.New()

		var wg sync.WaitGroup
		wg.Add(2)

		// Start writer...
		go func() {
			for i := 0; i < n; i++ {
				// Assert that the channel blocks for at most this amount of time.
				// Should be pretty small.
				timeout := time.After(time.Microsecond * 100)
				select {
				case c.In() <- i:
					t.Logf("wrote %d", i)
					continue
				case <-timeout:
					t.Errorf("writer iteration %d timed out", i)
				}
			}
			close(c.In())
			wg.Done()
		}()
		// ...and reader.
		go func() {
			i := 0
			done := false
			for !done {
				select {
				case v, ok := <-c.Out():
					if !ok {
						done = true
						continue
					}

					x, ok := v.(int)
					if !ok {
						t.Errorf("type assert failed, got %v", x)
					}

					if x != i {
						t.Errorf("reader got: %d want: %d", x, i)
					}
					i++
				}
			}
			if i != n {
				t.Errorf("reader got: %d want: %d", i, n)
			}
			wg.Done()
		}()
		wg.Wait()
	}
}

// Test that closing input closes output.
func TestClose(t *testing.T) {
	for _, n := range []int{
		0, 1, 10, 20,
	} {
		name := fmt.Sprintf("n=%d", n)
		t.Run(name, testClose(n))
	}

}

// testReceivers tests a broadcaster with a writer at rate w, and receivers at rates rs.
func testReceivers(w time.Duration, rs []time.Duration) func(*testing.T) {
	return func(t *testing.T) {
		b := bufchan.NewBroadcaster()

		count := 100

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel() // Cleanup.

		var wg sync.WaitGroup

		// Start readers before writer, so they get all of the results.
		for n, r := range rs {
			n, r := n, r
			wg.Add(1)
			l := b.Listen(ctx)
			go func(c <-chan interface{}) {
				ticker := time.NewTicker(r)
				defer ticker.Stop()

				for i := 0; i < count; i++ {
					<-ticker.C

					// Can't make any assertions about how long reader blocks for;
					// may be for a long time, if the writer is slower than the reader.
					x, ok := <-c
					if !ok {
						t.Errorf("channel closed too early, iteration %d", i)
					}
					n, ok := x.(int)
					if !ok {
						t.Errorf("type assertion failed, count %d: %v", i, n)
					}
					if n != i {
						t.Errorf("got: %d want: %d", i, n)
					}
				}
				// After a short sync delay, channel should be closed; have written, and read, count timestamps.
				time.Sleep(w)

				if _, ok := <-c; ok {
					t.Errorf("expected input channel %d to be closed, was open", n)
				}

				wg.Done()
			}(l)
		}

		// Start the writer.
		wg.Add(1)
		go func() {
			ticker := time.NewTicker(w)
			defer ticker.Stop()

			ichan := b.Send()

			for i := 0; i < count; i++ {
				<-ticker.C
				// We aren't asserting anything about the timings here,
				// though we expect them to be not very variable.
				ichan <- i
			}

			// Close the channel at the end.
			close(ichan)

			wg.Done()
		}()

		// Start readers


		// Wait for reader and writer to be done...
		wg.Wait()
	}
}

var ratesTable = []struct {
	w  time.Duration
	rs []time.Duration
}{
	{time.Millisecond, []time.Duration{time.Millisecond}},
	{time.Microsecond * 10, []time.Duration{time.Microsecond * 10, time.Microsecond * 10}},
	{time.Microsecond * 10, []time.Duration{time.Microsecond * 5, time.Microsecond * 100}},
	{time.Microsecond * 10, []time.Duration{
		time.Microsecond * 1, time.Microsecond * 2, time.Microsecond * 5,
		time.Microsecond * 10, time.Microsecond * 20, time.Microsecond * 50,
		time.Microsecond * 100, time.Microsecond * 200, time.Microsecond * 500,
	}},
}

func TestBroadcaster(t *testing.T) {
	for _, rates := range ratesTable {
		ss := make([]string, len(rates.rs))
		for i, r := range rates.rs {
			ss[i] = r.String()
		}
		name := fmt.Sprintf("w=%s/r=[%s]", rates.w, strings.Join(ss, ","))
		t.Run(name, testReceivers(rates.w, rates.rs))
	}
}

