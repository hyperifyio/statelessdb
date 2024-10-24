// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package events

import (
	"sync"
	"time"
)

// EventManager is responsible for handling event buffering, subscribing, and
// unsubscribing for each client
type EventManager[T comparable, D interface{}] struct {
	subscribers        map[T][]chan int64   // Notification channels per resource (state.Id)
	buffers            map[T][]*Event[T, D] // Event buffers per resource (state.Id)
	mu                 sync.Mutex           // Thread safety lock
	eventBus           EventBus[T, D]       // Global event bus. Safe to use from threads, immutable.
	eventChannel       chan *Event[T, D]    // Internal event channel. Safe to use from threads, immutable.
	bufferExpiration   time.Duration        // Duration after which events expire from the buffer. Safe to use from threads, immutable.
	cleanupInterval    time.Duration        // Interval to clean up events. Safe to use from threads, immutable.
	retryEventInterval time.Duration        // Interval to try sending event again. Safe to use from threads, immutable.
	maxRetries         int                  // Max retry times for triggering events. Safe to use from threads, immutable.
}

func NewEventManager[T comparable, D interface{}](
	bus EventBus[T, D],
	bufferExpiration time.Duration,
	cleanupInterval time.Duration,
	subscribersBufferSize int,
	internalBufferSize int,
	retryEventInterval time.Duration,
	maxRetries int,
) *EventManager[T, D] {

	m := &EventManager[T, D]{
		subscribers:        make(map[T][]chan int64, subscribersBufferSize),
		buffers:            make(map[T][]*Event[T, D]),
		eventBus:           bus,
		eventChannel:       make(chan *Event[T, D], internalBufferSize),
		bufferExpiration:   bufferExpiration,
		cleanupInterval:    cleanupInterval,
		retryEventInterval: retryEventInterval,
		maxRetries:         maxRetries,
	}

	// Start the event processing goroutine
	go m.processEvents()

	// Start the periodic buffer cleanup goroutine
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				m.cleanExpiredEvents()
			}
		}
	}()

	return m
}

// processEvents listens to the internal event channel and processes incoming events
func (m *EventManager[T, D]) processEvents() {
	for event := range m.eventChannel {
		m.mu.Lock()

		log.Debugf("[processEvents]: Event received %v %v", event.Type, event.Created)

		// Store the event in the buffer
		m.buffers[event.Type] = append(m.buffers[event.Type], event)

		// Notify all subscribers of this event type
		if subscribers, found := m.subscribers[event.Type]; found {
			for _, ch := range subscribers {
				go m.triggerEvent(ch, event.Type, event.Created)
			}
		}

		m.mu.Unlock()
	}
}

// Notify the subscriber
func (m *EventManager[T, D]) triggerEvent(ch chan int64, eventType T, value int64) {

	ticker := time.NewTicker(m.retryEventInterval)
	defer ticker.Stop()

	retries := m.maxRetries
	for i := 0; i < retries; i++ {
		select {
		case ch <- value:
			log.Debugf("[triggerEvent]: Event sent successfully to %v", eventType)
			return
		default:
			log.Warnf("[triggerEvent]: Subscriber was not ready -- waiting a moment (%d/%d): %v", i, retries, eventType)
			<-ticker.C
			continue
		}
	}

	log.Errorf("[triggerEvent]: Failed to send event: Subscriber was not ready -- skipped: %v", eventType)
}

// Subscribe adds a new subscription for a given id
func (m *EventManager[T, D]) Subscribe(stateId T, notificationChannel chan int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, found := m.buffers[stateId]; !found {
		log.Debugf("[Subscribe]: Subscribed for parent events: %v", stateId)
		m.eventBus.Subscribe(stateId, m.eventChannel)
	}

	log.Debugf("[Subscribe]: Client subscribed for: %v", stateId)
	m.subscribers[stateId] = append(m.subscribers[stateId], notificationChannel)

}

// Unsubscribe removes a subscription
func (m *EventManager[T, D]) Unsubscribe(stateId T, notificationChannel chan int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Remove the subscriber's notification channel
	if channels, found := m.subscribers[stateId]; found {
		for i, ch := range channels {
			if ch == notificationChannel {
				m.subscribers[stateId] = append(channels[:i], channels[i+1:]...)
				break
			}
		}
		if len(m.subscribers[stateId]) == 0 {
			delete(m.subscribers, stateId)
			log.Debugf("[Unsubscribe]: Last client unsubscribed for: %v", stateId)
		} else {
			log.Debugf("[Unsubscribe]: Client unsubscribed for: %v", stateId)
		}
	} else {
		log.Warnf("[Unsubscribe]: Warning: Client was not subscribed: %v", stateId)
	}
}

// GetBufferedEvents returns any buffered events for a specific client
func (m *EventManager[T, D]) GetBufferedEvents(stateId T, since int64) []*Event[T, D] {
	m.mu.Lock()
	defer m.mu.Unlock()

	if bufferedEvents, found := m.buffers[stateId]; found && len(bufferedEvents) > 0 {
		log.Debugf("[GetBufferedEvents]: Client requesting buffer for: %v since %d (%d events found)", stateId, since, len(bufferedEvents))

		var filteredEvents []*Event[T, D]
		for _, event := range bufferedEvents {
			if event.Created >= since {
				filteredEvents = append(filteredEvents, event)
			}
		}
		return filteredEvents
	} else {
		log.Debugf("[GetBufferedEvents]: Client requesting buffer for: %v since %d (no events found)", stateId, since)
	}

	return nil
}

// cleanExpiredEvents removes events from the buffer that have expired
func (m *EventManager[T, D]) cleanExpiredEvents() {
	m.mu.Lock()
	defer m.mu.Unlock()

	cutoffTime := time.Now().Add(-m.bufferExpiration).UnixMilli()

	for stateId, events := range m.buffers {

		var newEvents []*Event[T, D]
		for _, e := range events {
			if e.Created >= cutoffTime {
				newEvents = append(newEvents, e)
			}
		}

		totalCount := len(events)
		leftCount := len(newEvents)
		removedCount := totalCount - leftCount
		log.Debugf("[cleanExpiredEvents]: Cleaning expired events for %v since %d: removing %d of %d events", stateId, cutoffTime, removedCount, totalCount)

		if leftCount == 0 {

			if _, found := m.subscribers[stateId]; !found {
				log.Debugf("[cleanExpiredEvents]: Unsubscribed for parent events: %v", stateId)
				m.eventBus.Unsubscribe(stateId, m.eventChannel)
			}
			delete(m.buffers, stateId)

		} else {
			m.buffers[stateId] = newEvents
		}
	}
}
