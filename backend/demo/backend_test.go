package demo_test

import (
	"testing"

	"github.com/cceckman/discoirc/backend/demo"
	"github.com/cceckman/discoirc/backend/mocks"
	"github.com/cceckman/discoirc/data"
)

func TestSubscribeFiltered(t *testing.T) {
	b := demo.New()

	// Initialize data: two lines in to sonnet 18
	b.TickNetwork("sonnet")
	b.TickChannel("sonnet", "#eighteen")
	b.TickMessages("sonnet", "#eighteen")

	// And a couple dummy messages
	b.TickNetwork("botnet")
	b.TickChannel("botnet", "#t3000")
	b.TickChannel("sonnet", "#one90one")
	b.TickMessages("sonnet", "#one90one")

	ch := mocks.NewChannel("sonnet", "#eighteen")
	defer ch.Close()

	b.SubscribeFiltered(ch)

	var fst, snd data.ChannelState

	// First portion of test: Got initial state-fill
	ch.Join(func() {
		_, ok := ch.Nets["sonnet"]
		if !ok || len(ch.Nets) != 1 {
			t.Errorf("unexpected networks: got: %v wanted: %q", ch.Nets, "sonnet")
		}
		fst, ok = ch.Chans[mocks.ChannelIdent{
			Network: "sonnet",
			Channel: "#eighteen",
		}]

		if !ok || len(ch.Chans) != 1 {
			t.Errorf("unexpected channels: got: %v wanted: %q", ch.Chans, "sonnet #discoirc")
		}
	})

	// Send message updats
	b.TickMessages("sonnet", "#eighteen")

	ch.Join(func() {
		var ok bool
		snd, ok = ch.Chans[mocks.ChannelIdent{
			Network: "sonnet",
			Channel: "#eighteen",
		}]
		if !ok || len(ch.Chans) != 1 {
			t.Errorf("unexpected channels: got: %v wanted: %q", ch.Chans, "sonnet #discoirc")
		}
		if fst.LastMessage == snd.LastMessage {
			t.Errorf("didn't receive new messages: got: %v / %v", fst.LastMessage, snd.LastMessage)
		}
	})
}

func TestSubscribe_FromUI(t *testing.T) {
	b := demo.New()

	// Initialize data: two lines in to sonnet 18
	b.TickNetwork("sonnet")
	b.TickChannel("sonnet", "#eighteen")
	b.TickMessages("sonnet", "#eighteen")

	// And a couple dummy messages
	b.TickNetwork("botnet")
	b.TickChannel("botnet", "#t3000")
	b.TickChannel("sonnet", "#one90one")
	b.TickMessages("sonnet", "#one90one")

	c := mocks.NewClient()
	defer c.Close()

	c.Join(func() {
		b.Subscribe(c)
	})

	// First portion of test: Got initial state-fill
	c.Join(func() {
		_, ok := c.Nets["sonnet"]
		if !ok || len(c.Nets) != 1 {
			t.Errorf("unexpected networks: got: %v wanted: %q", c.Nets, "sonnet")
		}
	})
}
