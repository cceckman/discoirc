// Package model provides models for chat and connection updates.
package model

import (
	"context"
	"sync"
)

type Channel interface {
	Name() string

	// GetMessages requests a range of messages.
	// 'size' gives how many messages to return; "start" gives a starting offset, where an offset of 0
	// begins with the most recent message.
	// The returned slice is sorted from least (0) to most (len-1) recent.
	GetMessages(offset, size uint) []string
	SendMessage(string)

	GetTopic() string
	SetTopic(string)

	// Await awaits for changes to this channel.
	Await(context.Context) <-chan *Notification
}

// Notification represents an update to the channel.
type Notification struct {
	// Messages represents the total count of messages available.
	// A Notification receiver should always request the maximum number of messages it is able
	// to display, as more messages may have been received by the Channel since the Notification
	// arrived.
	Messages int
	// Topic indicates the topic of the channel.
	Topic string

	// Next is a channel to listen on for the next notification.
	Next chan *Notification
}

// MockChannel implements the Channel model, for testing.
type MockChannel struct {
	name string

	messages []string
	topic    string

	mu sync.RWMutex

	notification chan *Notification
}

// Notify asyncrhonously notifies listeners that the Channel has been updated.
// It is thread-safe (it spawns its own goroutine.)
func (m *MockChannel) Notify() {
	go func() {
		next := make(chan *Notification, 1)

		m.mu.RLock()
		new := &Notification{
			Messages: len(m.messages),
			Topic:    m.topic,
			Next:     next,
		}
		m.mu.RUnlock()

		m.mu.Lock()
		m.notification <- new
		m.notification = next
		m.mu.Unlock()
	}()
}

func (m *MockChannel) Name() string {
	return m.name
}

func (m *MockChannel) GetMessages(offset, size uint) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// The 'messages' dict is indexed oldest to newest, but offset here is from newest.
	// i                                offset
	// 0 hi                             3
	// 1 how are you                    2
	// 2 i'm well thank you very much   1
	// 3 ok bye                         0
	// But "end" of a slice is one beyond the end.
	end := len(m.messages) - int(offset)
	if end < 1 {
		return []string{}
	}
	start := end - int(size)
	if start < 0 {
		start = 0
	}

	return m.messages[start:end]
}

func (m *MockChannel) SendMessage(msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	defer m.Notify()

	m.messages = append(m.messages, msg)
}

func (m *MockChannel) GetTopic() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.topic
}

func (m *MockChannel) SetTopic(topic string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Don't notify if it's a spurious update.
	if topic != m.topic {
		defer m.Notify()
		m.topic = topic
	}
}

func (m *MockChannel) Await(ctx context.Context) <-chan *Notification {
	c := make(chan *Notification)

	go func() {
		defer close(c)

		await := m.notification

		select {
			case <-ctx.Done():
				return
			case notification := <-await:
				// Send on to the next listener; non-blocking
				await <- notification
				// Update our listener...
				await = notification.Next
				// And notify our own consumer. This is blocking.
				c <- notification
		}
	}()

	return c
}

func NewMockChannel(name string) Channel {
	r := &MockChannel{
		name:          name,
		notification: make(chan *Notification, 1),
	}
	return r
}
