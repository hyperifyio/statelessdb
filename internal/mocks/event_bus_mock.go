package mocks

import (
	"github.com/google/uuid"
	"statelessdb/pkg/events"
	"sync"
)

// MockEventBus is a mock implementation of EventBus
type MockEventBus[T comparable, D interface{}] struct {
	subscribers map[T][]chan *events.Event[T, D]
	mu          sync.Mutex
}

var _ events.EventBus[uuid.UUID, interface{}] = &MockEventBus[uuid.UUID, interface{}]{}

// NewMockEventBus initializes a new MockEventBus
func NewMockEventBus[T comparable, D interface{}]() *MockEventBus[T, D] {
	return &MockEventBus[T, D]{
		subscribers: make(map[T][]chan *events.Event[T, D]),
	}
}

// Subscribe adds a subscriber channel for a given type
func (m *MockEventBus[T, D]) Subscribe(id T, ch chan *events.Event[T, D]) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscribers[id] = append(m.subscribers[id], ch)
}

// Unsubscribe removes a subscriber channel for a given type
func (m *MockEventBus[T, D]) Unsubscribe(id T, ch chan *events.Event[T, D]) {
	m.mu.Lock()
	defer m.mu.Unlock()
	channels, found := m.subscribers[id]
	if !found {
		return
	}
	for i, subscriber := range channels {
		if subscriber == ch {
			m.subscribers[id] = append(channels[:i], channels[i+1:]...)
			break
		}
	}
	if len(m.subscribers[id]) == 0 {
		delete(m.subscribers, id)
	}
}

// Publish sends an event to all subscribers of the event's type
func (m *MockEventBus[T, D]) Publish(event *events.Event[T, D]) {
	m.mu.Lock()
	defer m.mu.Unlock()
	subscribers, found := m.subscribers[event.Type]
	if !found {
		return
	}
	for _, ch := range subscribers {
		// Non-blocking send to prevent goroutine leaks
		select {
		case ch <- event:
		default:
		}
	}
}
