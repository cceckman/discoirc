package demo_test

import (
	"testing"
	"time"

	"github.com/cceckman/discoirc/backend/demo"
	"github.com/cceckman/discoirc/backend/mocks"
	"github.com/cceckman/discoirc/data"
)

var (
	sonnet   = data.Scope{Net: "sonnet"}
	eighteen = data.Scope{Net: "sonnet", Name: "#eighteen"}
)

func delay(n int) int {
	time.Sleep(200 * time.Millisecond)
	return n + 1
}

func TestSubscribeFiltered(t *testing.T) {
	b := demo.New()

	// Initialize data: two lines in to sonnet 18
	b.TickNetwork(sonnet.Net)
	b.TickChannel(eighteen.Net, eighteen.Name)
	b.TickMessages(eighteen.Net, eighteen.Name)

	// And a couple dummy messages
	b.TickNetwork("botnet")
	b.TickChannel("botnet", "#t3000")
	b.TickChannel(sonnet.Net, "#one90one")
	b.TickMessages(sonnet.Net, "#one90one")

	ch := mocks.NewChannel(eighteen.Net, eighteen.Name)

	b.Subscribe(ch)

	var fst, snd, thd data.ChannelState
	var ok bool

	attempts := 4

	// Poll for 1s for up-to-date-ness.
	for i, done := 0, false; !(done || i > attempts); i = delay(i) {
		// First portion of test: Got initial state-fill
		ch.Join(func() {
			_, ok = ch.Nets[sonnet]
			fst, ok = ch.Chans[eighteen]

			expected_network := ok && len(ch.Nets) == 1
			expected_channel := ok && len(ch.Chans) == 1
			done = expected_network && expected_channel

			if !expected_network && i == attempts {
				t.Errorf("unexpected networks: got: %v wanted: %q", ch.Nets, eighteen.Net)
			}

			if !expected_channel && i == attempts {
				t.Errorf("unexpected channels: got: %v wanted: %q", ch.Chans, "sonnet #discoirc")
			}
		})
	}

	// Send message updates
	b.TickMessages(eighteen.Net, eighteen.Name)

	for i, done := 0, false; !(done || i > attempts); i = delay(i) {
		// Second portion: expect update to messages and unread.
		ch.Join(func() {
			snd, ok = ch.Chans[eighteen]

			expected_channel := ok && len(ch.Chans) == 1
			expected_message := fst.LastMessage != snd.LastMessage
			expected_unread := fst.Unread < snd.Unread
			done = expected_channel && expected_message && expected_unread

			if !expected_channel && i == attempts {
				t.Errorf("unexpected channels: got: %v wanted: %q", ch.Chans, "sonnet #discoirc")
			}
			if !expected_message && i == attempts {
				t.Errorf("didn't receive new messages: got %v, then %v", fst.LastMessage, snd.LastMessage)
			}

			if !expected_unread && i == attempts {
				t.Errorf("didn't see unread updated: got %v, then %v", fst.Unread, snd.Unread)
			}
		})
	}

	// Test unread clearing
	go func() {
		_ = b.EventsBefore(eighteen, 1000, snd.LastMessage)
	}()

	for i, done := 0, false; !(done || i > attempts); i = delay(i) {
		// Second portion: expect update to messages and unread.
		ch.Join(func() {
			thd, ok = ch.Chans[eighteen]

			expected_unread := thd.Unread < snd.Unread
			expected_zero := thd.Unread == 0
			done = expected_unread && expected_zero

			if !expected_unread && i == attempts {
				t.Errorf("didn't see unread messages cleared: got %d, then %d", snd.Unread, thd.Unread)
			}

			if !expected_zero && i == attempts {
				t.Errorf("unread did not reset to zero: got %d want %d", thd.Unread, 0)
			}
		})
	}
}

func TestSubscribe_FromUI(t *testing.T) {
	b := demo.New()

	// Initialize data: two networks
	b.TickMessages(eighteen.Net, eighteen.Name)
	b.TickNetwork("botnet")

	c := mocks.NewClient()

	// Run subscribe in the UI thread, make sure we don't get a race.
	c.Join(func() {
		b.Subscribe(c)
	})

	attempts := 4

	for i, done := 0, false; !(done || i > attempts); i = delay(i) {
		// First portion of test: Got initial state-fill
		c.Join(func() {
			expect_networks := len(c.Nets) == 2
			done = expect_networks
			if !expect_networks && i == attempts {
				t.Errorf("unexpected networks: got: %v wanted: %d", c.Nets, 2)
			}
		})
	}
}

func TestChannelCallback(t *testing.T) {
	attempts := 4
	b := demo.New()
	c := mocks.NewChannel(eighteen.Net, eighteen.Name)

	c.Archive = b
	b.Subscribe(c)

	// Send 3 messages in two channels
	for i := 0; i < 3; i++ {
		b.TickMessages(eighteen.Net, eighteen.Name)
		b.TickMessages(sonnet.Net, "#globe")
	}

	for i, done := 0, false; !(done || i > attempts); i = delay(i) {
		c.Join(func() {
			expect_channels := len(c.Contents) == 1
			expect_contents := len(c.Contents[eighteen]) == 3
			done = expect_channels && expect_contents

			if !expect_channels && i == attempts {
				t.Errorf("unexpected channels: got: %v want: %v", c.Contents, eighteen)
			}

			if !expect_contents && i == attempts {
				t.Errorf("unexpected contents: got: %v want: %v", c.Contents[eighteen], "three lines of Shakespare's sonnet 18")
			}
		})
	}
}

func TestSubscribe_Resubscribe(t *testing.T) {
	attempts := 4
	b := demo.New()

	// Initialize data: one network
	b.TickMessages(eighteen.Net, eighteen.Name)

	c1 := mocks.NewClient()
	go b.Subscribe(c1)

	for i, done := 0, false; !(done || i > attempts); i = delay(i) {
		// First portion of test: Got initial state-fill
		c1.Join(func() {
			expect_networks := len(c1.Nets) == 1
			done = expect_networks
			if !expect_networks && i == attempts {
				t.Errorf("unexpected networks: got: %v wanted: %d", c1.Nets, 1)
			}
		})
	}

	c2 := mocks.NewClient()
	go b.Subscribe(c2)

	b.TickMessages("botnet", "#eighteen")

	for i, done := 0, false; !(done || i > attempts); i = delay(i) {
		// First portion of test: Got initial state-fill
		c2.Join(func() {
			expect_networks := len(c2.Nets) == 2
			done = expect_networks
			if !expect_networks && i == attempts {
				t.Errorf("unexpected networks: got: %v wanted: %d", c2.Nets, 2)
			}
		})
	}

	b.Subscribe(nil)

	c2.Join(func() {
		// Do nothing; just verifying closure.
	})
}
