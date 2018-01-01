package demo_test

import (
	"testing"

	"github.com/cceckman/discoirc/backend"
	"github.com/cceckman/discoirc/backend/demo"
	"github.com/cceckman/discoirc/backend/mocks"
	"github.com/cceckman/discoirc/data"
)

var (
	_ backend.DataPublisher = demo.New()
)

// TODO: Test filtered version

func TestClientView(t *testing.T) {
	view := mocks.NewClient()

	// Prepare some state before the channel attaches.
	d := demo.New()
	// Network tick changes nick, cycles connection state
	d.TickNetwork("Barnet")
	// Channel tick changes mode, sends a message, changes topic
	d.TickChannel("Barnet", "#discoirc")
	d.TickChannel("Barnet", "&somethingelse")

	// Register the subscription.
	d.Subscribe(view)
	// Add some more events.
	d.TickNetwork("Baznet")
	d.TickNetwork("Baznet")
	d.TickChannel("Slacknet", "#tuigo")
	d.TickChannel("Slacknet", "#tuigo")

	// Cancel out subscription.
	d.Subscribe(nil)

	// Add some more events; shouldn't be received, after cancellation.
	d.TickNetwork("Network9")
	d.TickChannel("Network10", "#nochan")

	// Check that the networks are right first; everything else depends on that.
	wantNets := []string{"Barnet", "Baznet", "Slacknet"}
	if len(view.Nets) != len(wantNets) {
		t.Fatalf("incorrect number of networks: got: %v want: %v", view.Nets, wantNets)
	}
	mismatch := false
	for _, v := range wantNets {
		if _, ok := view.Nets[v]; !ok {
			mismatch = true
			t.Errorf("missing network: want: %s", v)
		}
	}
	if mismatch {
		return
	}

	// Check channel list.
	wantChannels := map[string][]string{
		"Barnet":   []string{"#discoirc", "&somethingelse"},
		"Baznet":   []string{},
		"Slacknet": []string{"#tuigo"},
	}

	for net, channels := range wantChannels {
		for _, channel := range channels {
			index := mocks.ChannelIdent{
				Network: net,
				Channel: channel,
			}

			if _, ok := view.Chans[index]; !ok {
				t.Errorf("missing channel: want: %v", index)
			}
		}

		// Lists are correct. Check properties.

		//////////////////////////

		barnet := view.Nets["Barnet"]
		if barnet.State != data.Connecting {
			t.Errorf("unexpected state for Barnet: got: %v want: %v", barnet.State, data.Connecting)
		}
		if barnet.Nick == "" {
			t.Errorf("unexpected nick for Barnet: %v want: nonempty", barnet.Nick)
		}

		//////////////////////////

		discoirc := view.Chans[mocks.ChannelIdent{
			Network: "Barnet",
			Channel: "#discoirc",
		}]
		if discoirc.Unread == 0 {
			t.Errorf("no unread messages for #discoirc")
		}
		if discoirc.Members == 0 {
			t.Errorf("no members for #discoirc")
		}
		if discoirc.ChannelMode == "" {
			t.Errorf("no mode for #discoirc")
		}
		if discoirc.UserMode == "" {
			t.Errorf("no user mode for #discoirc")
		}
		if discoirc.Presence != data.Joined {
			t.Errorf("not joined for #discoirc")
		}

		//////////////////////////

		somethingelse := view.Chans[mocks.ChannelIdent{
			Network: "Barnet",
			Channel: "&somethingelse",
		}]
		if somethingelse.Unread == 0 {
			t.Errorf("no unread messages for &somethingelse")
		}
		if somethingelse.Members == 0 {
			t.Errorf("no members for &somethingelse")
		}
		if somethingelse.ChannelMode == "" {
			t.Errorf("no mode for #somethingelse")
		}
		if somethingelse.UserMode == "" {
			t.Errorf("no user mode for #somethingelse")
		}
		if somethingelse.Presence != data.Joined {
			t.Errorf("not joined for &somethingelse")
		}

		//////////////////////////

		baznet := view.Nets["Baznet"]
		if baznet.State != data.Connected {
			t.Errorf("unexpected state for Baznet: got: %v want: %v", baznet.State, data.Connected)
		}
		if baznet.Nick == "" {
			t.Errorf("unexpected nick for Baznet: %v want: nonempty", baznet.Nick)
		}

		//////////////////////////

		tuigo := view.Chans[mocks.ChannelIdent{
			Network: "Slacknet",
			Channel: "#tuigo",
		}]
		if tuigo.Unread != 4 {
			t.Errorf("unexpected unread messages for #tuigo: got: %d want: %d", tuigo.Unread, 2)
		}
		if tuigo.Members != 2 {
			t.Errorf("unexpected members for #tuigo: got: %d want: %d", tuigo.Members, 3)
		}
		if tuigo.ChannelMode == "" {
			t.Errorf("no mode for #tuigo")
		}
		if tuigo.UserMode == "" {
			t.Errorf("no user mode for #tuigo")
		}
		if tuigo.Presence != data.NotPresent {
			t.Errorf("not joined for #tuigo")
		}
	}
}
