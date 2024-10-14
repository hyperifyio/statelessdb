// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package main_test

import (
	"errors"
	"github.com/hyperifyio/statelessdb/pkg/helpers"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hyperifyio/statelessdb/pkg/dtos"
	"github.com/hyperifyio/statelessdb/pkg/events"
	"github.com/hyperifyio/statelessdb/pkg/requests"
	"github.com/hyperifyio/statelessdb/pkg/states"

	"github.com/hyperifyio/statelessdb/cmd/statelessdb"
)

const (
	eventTimeoutTime         = 100 * time.Millisecond // Default request timeout time
	eventExpirationTime      = 200 * time.Millisecond // Time until events expire
	eventCleanupIntervalTime = 300 * time.Millisecond // Interval to clean up expired events
)

func TestApiEventHandler_NilState(t *testing.T) {
	// Initialize MockEventBus
	//var mockBus events.EventBus[uuid.UUID, interface{}] = mocks.NewMockEventBus[uuid.UUID, interface{}]()

	// Initialize EventManager with mockBus
	//manager := events.NewEventManager[uuid.UUID, interface{}](mockBus, 20*time.Second, 30*time.Second)

	// Create main.ApiEventHandler
	handler := main.ApiEventHandler(events.NewLocalEventBus[uuid.UUID, interface{}](), eventTimeoutTime, eventExpirationTime, eventCleanupIntervalTime) // Pass actual LocalEventBus or mock as needed

	// Create a ComputeRequest
	req := &requests.ComputeRequest{
		PrivateData: "dummy_private_data",
	}

	// Call handler with nil state
	_, err := handler(req, nil)
	assert.Error(t, err, "Expected error when state is nil")
	assert.Equal(t, main.ErrNoStateProvided, err, "Expected ErrNoStateProvided")
}

// CustomComputeState that fails initialization
type FailingComputeState struct {
	states.ComputeState
}

func (s *FailingComputeState) Initialize() error {
	return errors.New("initialization failed")
}

//func TestApiEventHandler_InitializationError(t *testing.T) {
//	// Initialize MockEventBus
//	//mockBus := mocks.NewMockEventBus[uuid.UUID, interface{}]()
//
//	// Initialize EventManager with mockBus
//	//manager := events.NewEventManager(mockBus, 20*time.Second, 30*time.Second)
//
//	// Create main.ApiEventHandler
//	handler := main.ApiEventHandler(events.NewLocalEventBus[uuid.UUID, interface{}](), eventTimeoutTime, eventExpirationTime, eventCleanupIntervalTime) // Pass actual LocalEventBus or mock as needed
//
//	// Create a failing ComputeState
//	state := &FailingComputeState{}
//
//	// Create a ComputeRequest
//	req := &requests.ComputeRequest{
//		PrivateData: "dummy_private_data",
//	}
//
//	// Call handler with failing state
//	_, err := handler(req, state)
//	assert.Error(t, err, "Expected error during state initialization")
//	assert.Equal(t, "initialization failed", err.Error(), "Expected initialization failure error")
//}

func TestApiEventHandler_BufferedEvents(t *testing.T) {
	// Initialize MockEventBus
	//mockBus := mocks.NewMockEventBus[uuid.UUID, interface{}]()

	// Initialize EventManager with mockBus
	//manager := events.NewEventManager(mockBus, 20*time.Second, 30*time.Second)

	// Create main.ApiEventHandler
	handler := main.ApiEventHandler(events.NewLocalEventBus[uuid.UUID, interface{}](), eventTimeoutTime, eventExpirationTime, eventCleanupIntervalTime) // Pass actual LocalEventBus or mock as needed

	// Create a ComputeState with buffered events
	stateID := uuid.New()
	state := states.NewComputeState(
		stateID,
		uuid.New(),
		time.Now().UnixMilli(),
		time.Now().UnixMilli(),
		nil,
		nil,
		[]*events.Event[uuid.UUID, interface{}]{
			{
				Type:    stateID,
				Data:    "event1",
				Created: time.Now().UnixMilli(),
			},
			{
				Type:    stateID,
				Data:    "event2",
				Created: time.Now().UnixMilli(),
			},
		},
	)

	// Create a ComputeRequest
	req := &requests.ComputeRequest{
		PrivateData: "dummy_private_data",
	}

	// Call handler
	updatedState, err := handler(req, state)
	assert.NoError(t, err, "Expected no error when processing buffered events")
	assert.Equal(t, state.Id, updatedState.Id, "State ID should remain the same")
	assert.Equal(t, 2, len(updatedState.Events()), "State should have 2 events")
	assert.Equal(t, "event1", updatedState.Events()[0].Data, "First event should be 'event1'")
	assert.Equal(t, "event2", updatedState.Events()[1].Data, "Second event should be 'event2'")
}

func TestApiEventHandler_SubscribeAndTimeout(t *testing.T) {
	// Initialize MockEventBus
	//mockBus := mocks.NewMockEventBus[uuid.UUID, interface{}]()

	// Initialize EventManager with mockBus
	//bufferExpiration := 20 * time.Second
	//cleanupInterval := 30 * time.Second
	//manager := events.NewEventManager(mockBus, bufferExpiration, cleanupInterval)

	eventBus := events.NewLocalEventBus[uuid.UUID, interface{}]()

	// Create main.ApiEventHandler
	handler := main.ApiEventHandler(eventBus, eventTimeoutTime, eventExpirationTime, eventCleanupIntervalTime)

	// Create a ComputeState with no buffered events
	stateID := uuid.New()
	state := states.NewComputeState(
		stateID,
		uuid.New(),
		time.Now().UnixMilli(),
		time.Now().UnixMilli(),
		nil,
		nil,
		nil,
	)

	// Create a ComputeRequest
	req := &requests.ComputeRequest{
		PrivateData: "dummy_private_data",
	}

	timeInterval := eventTimeoutTime / 10

	// Start a goroutine to publish an event after 2 seconds
	go func() {
		time.Sleep(timeInterval)
		event := &events.Event[uuid.UUID, interface{}]{
			Type:    stateID,
			Data:    "new_event",
			Created: time.Now().UnixMilli(),
		}
		eventBus.Publish(event)
	}()

	// Call handler (timeout is set to 10 seconds)
	start := time.Now()
	updatedState, err := handler(req, state)
	duration := time.Since(start)

	assert.NoError(t, err, "Expected no error when event is published before timeout")
	assert.True(t, duration >= timeInterval, "Handler should wait for the event")
	assert.Equal(t, state.Id, updatedState.Id, "State ID should remain the same")
	assert.Equal(t, 1, len(updatedState.Events()), "State should have 1 event")
	assert.Equal(t, "new_event", updatedState.Events()[0].Data, "Event data should be 'new_event'")
}

func TestApiEventHandler_TimeoutWithoutEvents(t *testing.T) {
	// Initialize MockEventBus
	//mockBus := mocks.NewMockEventBus[uuid.UUID, interface{}]()

	//// Initialize EventManager with mockBus
	//bufferExpiration := 20 * time.Second
	//cleanupInterval := 30 * time.Second
	//manager := events.NewEventManager(mockBus, bufferExpiration, cleanupInterval)

	// Create main.ApiEventHandler
	handler := main.ApiEventHandler(events.NewLocalEventBus[uuid.UUID, interface{}](), eventTimeoutTime, eventExpirationTime, eventCleanupIntervalTime) // Pass actual LocalEventBus or mock as needed

	// Create a ComputeState with no buffered events
	stateID := uuid.New()
	state := states.NewComputeState(
		stateID,
		uuid.New(),
		time.Now().UnixMilli(),
		time.Now().UnixMilli(),
		nil,
		nil,
		nil,
	)

	// Create a ComputeRequest
	req := &requests.ComputeRequest{
		PrivateData: "dummy_private_data",
	}

	// Call handler (timeout is set to 10 seconds)
	start := time.Now()
	updatedState, err := handler(req, state)
	duration := time.Since(start)

	assert.NoError(t, err, "Expected no error when timeout occurs without events")
	assert.True(t, duration >= eventTimeoutTime, "Handler should timeout after 10 seconds")
	assert.Equal(t, state.Id, updatedState.Id, "State ID should remain the same")
	assert.Equal(t, 0, len(updatedState.Events()), "State should have no events")
}

func TestApiEventHandler_ConcurrentAccess(t *testing.T) {
	// Initialize MockEventBus
	//mockBus := mocks.NewMockEventBus[uuid.UUID, interface{}]()

	// Initialize EventManager with mockBus
	//bufferExpiration := 20 * time.Second
	//cleanupInterval := 30 * time.Second
	//manager := events.NewEventManager(mockBus, bufferExpiration, cleanupInterval)

	// Create main.ApiEventHandler
	handler := main.ApiEventHandler(events.NewLocalEventBus[uuid.UUID, interface{}](), eventTimeoutTime, eventExpirationTime, eventCleanupIntervalTime) // Pass actual LocalEventBus or mock as needed

	// Define number of concurrent goroutines
	const goroutines = 50

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Channel to collect errors
	errCh := make(chan error, goroutines)

	// Create a ComputeState with no buffered events
	stateID := uuid.New()
	state := states.NewComputeState(
		stateID,
		uuid.New(),
		time.Now().UnixMilli(),
		time.Now().UnixMilli(),
		nil,
		nil,
		nil,
	)

	// Launch concurrent goroutines
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// Each goroutine creates its own ComputeRequest
			req := &requests.ComputeRequest{
				PrivateData: "dummy_private_data",
			}

			// Call handler
			_, err := handler(req, state)
			if err != nil {
				errCh <- err
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errCh)

	// Check for errors
	for err := range errCh {
		if err != nil {
			t.Errorf("Concurrent main.ApiEventHandler failed: %v", err)
		}
	}

	// No assertions needed if no errors are reported
}

func TestNewEventResponseDTO_NoEvents(t *testing.T) {
	// Initialize MockEventBus (not used in this test)
	//mockBus := mocks.NewMockEventBus[uuid.UUID, interface{}]()

	// Create main.NewEventResponseDTO function
	createResponse := main.NewEventResponseDTO(events.NewLocalEventBus[uuid.UUID, interface{}]()) // Pass actual LocalEventBus or mock as needed

	// Create a ComputeState with no events
	state := states.NewComputeState(
		uuid.New(),
		uuid.New(),
		time.Now().UnixMilli(),
		time.Now().UnixMilli(),
		nil,
		nil,
		nil,
	)

	privateData := "encrypted_private_data"

	// Call CreateResponseFunc
	response := createResponse(state, privateData)

	// Assert response type
	eventListDTO, ok := response.(*dtos.EventListDTO)
	assert.True(t, ok, "Response should be of type *dtos.EventListDTO")

	// Assert fields
	assert.Equal(t, helpers.MillisToISO(state.Updated), eventListDTO.Created, "Updated time should match")
	assert.Equal(t, 0, len(eventListDTO.Payload), "Events list should be empty")
	assert.Equal(t, privateData, eventListDTO.Private, "Private data should match")
}

func TestNewEventResponseDTO_WithEvents(t *testing.T) {
	// Initialize MockEventBus (not used in this test)
	//mockBus := mocks.NewMockEventBus[uuid.UUID, interface{}]()

	// Create main.NewEventResponseDTO function
	createResponse := main.NewEventResponseDTO(events.NewLocalEventBus[uuid.UUID, interface{}]()) // Pass actual LocalEventBus or mock as needed

	// Create a ComputeState with multiple events
	event1 := events.NewEvent[uuid.UUID, interface{}](uuid.New(), "event_data_1", time.Now().UnixMilli())
	event2 := events.NewEvent[uuid.UUID, interface{}](uuid.New(), "event_data_2", time.Now().UnixMilli())

	state := states.NewComputeState(
		uuid.New(),
		uuid.New(),
		time.Now().UnixMilli(),
		time.Now().UnixMilli(),
		nil,
		nil,
		[]*events.Event[uuid.UUID, interface{}]{event1, event2},
	)

	privateData := "encrypted_private_data"

	// Call CreateResponseFunc
	response := createResponse(state, privateData)

	// Assert response type
	eventListDTO, ok := response.(*dtos.EventListDTO)
	assert.True(t, ok, "Response should be of type *dtos.EventListDTO")

	// Assert fields
	assert.Equal(t, helpers.MillisToISO(state.Updated), eventListDTO.Created, "Updated time should match")
	assert.Equal(t, 2, len(eventListDTO.Payload), "Events list should contain 2 events")
	assert.Equal(t, "event_data_1", eventListDTO.Payload[0].Data, "First event data should match")
	assert.Equal(t, "event_data_2", eventListDTO.Payload[1].Data, "Second event data should match")
	assert.Equal(t, privateData, eventListDTO.Private, "Private data should match")
}

func TestNewEventResponseDTO_ConcurrentAccess(t *testing.T) {
	// Initialize MockEventBus (not used in this test)
	//mockBus := mocks.NewMockEventBus[uuid.UUID, interface{}]()

	// Create main.NewEventResponseDTO function
	createResponse := main.NewEventResponseDTO(events.NewLocalEventBus[uuid.UUID, interface{}]()) // Pass actual LocalEventBus or mock as needed

	// Create a ComputeState with multiple events
	event1 := events.NewEvent[uuid.UUID, interface{}](uuid.New(), "event_data_1", time.Now().UnixMilli())
	event2 := events.NewEvent[uuid.UUID, interface{}](uuid.New(), "event_data_2", time.Now().UnixMilli())

	state := states.NewComputeState(
		uuid.New(),
		uuid.New(),
		time.Now().UnixMilli(),
		time.Now().UnixMilli(),
		nil,
		nil,
		[]*events.Event[uuid.UUID, interface{}]{event1, event2},
	)

	privateData := "encrypted_private_data"

	// Define number of concurrent goroutines
	const goroutines = 50

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Channel to collect errors
	errCh := make(chan error, goroutines)

	// Launch concurrent goroutines
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// Call CreateResponseFunc
			response := createResponse(state, privateData)

			// Assert response type
			eventListDTO, ok := response.(*dtos.EventListDTO)
			if !ok {
				errCh <- errors.New("Response is not of type *dtos.EventListDTO")
				return
			}

			// Assert fields
			if eventListDTO.Created != helpers.MillisToISO(state.Updated) {
				errCh <- errors.New("Updated time does not match")
				return
			}

			if len(eventListDTO.Payload) != 2 {
				errCh <- errors.New("Events list does not contain 2 events")
				return
			}

			if eventListDTO.Payload[0].Data != "event_data_1" || eventListDTO.Payload[1].Data != "event_data_2" {
				errCh <- errors.New("Event data does not match expected values")
				return
			}

			if eventListDTO.Private != privateData {
				errCh <- errors.New("Private data does not match")
				return
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errCh)

	// Check for errors
	for err := range errCh {
		if err != nil {
			t.Errorf("Concurrent main.NewEventResponseDTO failed: %v", err)
		}
	}

	// No assertions needed if no errors are reported
}

// main_test.go (continued)

type DecodeError struct {
	GoroutineID int
	Expected    string
	Actual      string
}

func (e *DecodeError) Error() string {
	return "DecodeRequest mismatch in goroutine " +
		string(rune(e.GoroutineID)) + ": expected '" + e.Expected + "', got '" + e.Actual + "'"
}
