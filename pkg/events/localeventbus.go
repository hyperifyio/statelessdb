// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package events

import "sync"

type LocalEventBus[T comparable, D interface{}] struct {
	subscribers map[T][]chan *Event[T, D]
	mu          sync.RWMutex
}

func NewLocalEventBus[T comparable, D interface{}](
	bufferSize int,
) *LocalEventBus[T, D] {
	return &LocalEventBus[T, D]{
		subscribers: make(map[T][]chan *Event[T, D], bufferSize),
	}
}

// Subscribe to a specific event
func (bus *LocalEventBus[T, D]) Subscribe(eventType T, ch chan *Event[T, D]) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	bus.subscribers[eventType] = append(bus.subscribers[eventType], ch)
}

// Unsubscribe from a specific event
func (bus *LocalEventBus[T, D]) Unsubscribe(eventType T, ch chan *Event[T, D]) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if channels, found := bus.subscribers[eventType]; found {

		for i := range channels {
			if channels[i] == ch {
				// Remove the channel by slicing out the matching entry
				bus.subscribers[eventType] = append(channels[:i], channels[i+1:]...)
				break
			}
		}

		// Clean up if no subscribers left for this event type
		if len(bus.subscribers[eventType]) == 0 {
			delete(bus.subscribers, eventType)
		}
	}
}

// Publish an event to all subscribers
func (bus *LocalEventBus[T, D]) Publish(event *Event[T, D]) {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	if channels, found := bus.subscribers[event.Type]; found {
		for _, ch := range channels {
			go func(c chan *Event[T, D]) {
				c <- event
			}(ch)
		}
	} else {
		log.Warnf("Nothing listening events by: %s", event.Type)
	}
}
