// 2016-12-31 cceckman <charles@cceckman.com>
package bufchan_test

import (
	"context"
	"fmt"
	"github.com/cceckman/discoirc/prototype/bufchan"
	"sync"
	"testing"
	"time"
)

// testAtRates generates a test case for correct behavior when reading and writing at the respective rates.
func testAtRates(r, w time.Duration) func(*testing.T) {
	return func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		c := bufchan.New(ctx)

		timeFmt := time.RFC3339

		var wg sync.WaitGroup
		wg.Add(2)

		// Start writer...
		go func() {
			ticker := time.NewTicker(w)
			defer ticker.Stop()
			for i := 0; i < 100; i++ {
				tm := <-ticker.C
				s := tm.Format(timeFmt)
				// Assert that the channel blocks for at most this amount of time.
				// Should be pretty small.
				timeout := time.After(time.Millisecond)
				select {
				case c.In() <- s:
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
				tm, err := time.Parse(timeFmt, <-c.Out())
				if err != nil {
					t.Errorf("error in reader: %v", err)
					continue
				}
				since := tm.Sub(lastTime)
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

		// Wait for reader and writer to be done...
		wg.Wait()
		// Then clean up.
		cancel()
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
