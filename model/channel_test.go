package model

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// TestMockChannel creates a MockChannel and exercises it.
func TestMockChannel(t *testing.T) {
	var m Channel = NewMockChannel(log.New(&bytes.Buffer{}, "", 0), "#foobar", "")
	if m.Name() != "#foobar" {
		t.Error("bad Name")
	}

	messages := []string{
		"0 Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		"1 Donec bibendum consequat nibh vitae vestibulum.",
		"2 Praesent rutrum massa ac lorem bibendum, a facilisis est vehicula.",
		"3 Vivamus efficitur vehicula eros id mollis.",
		"4 In hac habitasse platea dictumst.",
		"5 Donec id est scelerisque, mollis nibh ut, tempus ligula.",
		"6 Donec interdum faucibus leo ac rutrum.",
		"7 Vestibulum gravida tempor dui, vitae vulputate arcu ullamcorper ac.",
		"8 Ut porttitor libero at ipsum mattis elementum.",
		"9 Nullam odio odio, lacinia non venenatis non, sagittis nec purus.",
	}
	topics := []string{
		"",
		"no topic like the present",
		"no present like a topic",
	}

	var wg sync.WaitGroup

	// Consumers.
	for i := range []int{0, 1} {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			lastCount := 0
			lastTopic := 0
			notifications := 0
			// Receive all updates, in order.
			for n := range m.Await(ctx) {
				updateTopic := false
				updateCount := false

				// Topic updates are inlined.
				if n.Topic != topics[lastTopic] {
					t.Logf("listener [%d]: got topic '%s' after topic '%s', bumping up counter to %d'",
						i, n.Topic, topics[lastTopic])
					lastTopic += 1
					updateTopic = true
				}
				if lastTopic == len(topics) {
					t.Fatalf("listener [%d]: unexpected topic notification: got: '%s' after running off the list", i, n.Topic)
				} else if n.Topic != topics[lastTopic] {
					t.Errorf("listener [%d]: unexpected topic notification: got: '%s' want: '%s'", i,
						n.Topic, topics[lastTopic])
				}

				// Should receive an update for every message sent.
				if n.Messages != lastCount {
					lastCount = n.Messages
					updateCount = true
				}

				if !updateTopic && !updateCount {
					t.Errorf("listener [%d]: notification %d had no update to messages or to topic",
						i, notifications)
				}

				// Message updates require a second call to get the contents.
				received := m.GetMessages(0, uint(len(messages)))
				if len(received) < n.Messages {
					t.Errorf("listener [%d]: received notification for %d messages, but only drew %d back",
						i, n.Messages, len(received))
				}
				if n.Messages > len(messages) {
					t.Errorf("listener [%d]: received notification for %d messages, but expected at most %d",
						i, n.Messages, len(messages))
				}
				if len(received) > len(messages) {
					t.Errorf("listener [%d]: received %d messages, but expected at most %d",
						i, len(received), len(messages))
				}

				// "received" should be all of the messages from 0 to len(received) -1.
				for j := range received {
					if messages[j] != received[j] {
						t.Errorf("listener [%d]: unexpected message %d: got: '%s' want: '%s'", i, j, received[j], messages[j])
					}
				}

				notifications++

				if lastCount == len(messages) && lastTopic == len(topics)-1 {
					t.Logf("listener [%d]: complete after %d notifications", i, notifications)
					cancel()
				}
			}
		}()
	}

	// Producers.
	wg.Add(2)
	go func() {
		defer wg.Done()
		for _, msg := range messages {
			sleep := time.Duration(rand.NormFloat64()*20 + 200)
			if sleep < 0 {
				sleep = 0
			}
			time.Sleep(time.Millisecond * sleep)
			m.MessageInput() <- msg
		}
	}()
	go func() {
		defer wg.Done()
		for _, msg := range topics {
			sleep := time.Duration(rand.NormFloat64()*20 + 200)
			if sleep < 0 {
				sleep = 0
			}
			time.Sleep(time.Millisecond * sleep)
			m.SetTopic(msg)
		}
	}()

	wg.Wait()
}

// TestMockChannelGetMessages makes sure I've done my indexing right.
func TestMockChannelGetMessages(t *testing.T) {
	m := &MockChannel{
		messages: []string{
			"0 hello muddah",
			"1 hello faddah",
			"2 here I am at",
			"3 Camp Grenada",
		},
	}

	for i, cs := range []struct {
		Size, Offset uint
		Want         []string
	}{
		{0, 100, m.messages},
		{0, 1, []string{"3 Camp Grenada"}},
		{1, 1, []string{"2 here I am at"}},
		{4, 1, []string{}},
		{2, 3, []string{"0 hello muddah", "1 hello faddah"}},
	} {
		got := m.GetMessages(cs.Offset, cs.Size)
		err := fmt.Sprintf("unexpected result for test case %d: got: %v want: %v", i, got, cs.Want)
		if len(got) != len(cs.Want) {
			t.Errorf(err)
		}
		for i := range got {
			if got[i] != cs.Want[i] {
				t.Errorf(err)
			}
		}
	}
}
