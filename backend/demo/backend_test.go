package demo_test

import (
	"testing"

	"github.com/cceckman/discoirc/backend/demo"
	"github.com/cceckman/discoirc/backend/mocks"
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
	b.SubscribeFiltered(ch)

	// First portion of test: Got initial state-fill
	_, ok := ch.Nets["sonnet"]
	if !ok || len(ch.Nets) != 1 {
		t.Errorf("unexpected networks: got: %v wanted: %q", ch.Nets, "sonnet")
	}
	fstChanState, ok := ch.Chans[mocks.ChannelIdent{
		Network: "sonnet",
		Channel: "#eighteen",
	}]

	if !ok || len(ch.Chans) != 1 {
		t.Errorf("unexpected channels: got: %v wanted: %q", ch.Chans, "sonnet #discoirc")
	}


	// Send message updats
	b.TickMessages("sonnet", "#eighteen")

	sndChanState, ok := ch.Chans[mocks.ChannelIdent{
		Network: "sonnet",
		Channel: "#eighteen",
	}]
	if !ok || len(ch.Chans) != 1 {
		t.Errorf("unexpected channels: got: %v wanted: %q", ch.Chans, "sonnet #discoirc")
	}
	if fstChanState.LastMessage == sndChanState.LastMessage {
		t.Errorf("didn't receive new messages: got: %v / %v", fstChanState.LastMessage, sndChanState.LastMessage)
	}

	// Send some more dummy updates
	b.TickNetwork("botnet")
	b.TickChannel("botnet", "#t3000")
	b.TickChannel("sonnet", "#one90one")

	if _, ok := ch.Nets["sonnet"]; !ok || len(ch.Nets) != 1 {
		t.Errorf("unexpected networks: got: %v wanted: %q", ch.Nets, "sonnet")
	}




}
