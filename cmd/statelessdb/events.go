// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package main

import (
	"time"

	"github.com/google/uuid"

	"github.com/hyperifyio/statelessdb/pkg/dtos"
	"github.com/hyperifyio/statelessdb/pkg/events"
	"github.com/hyperifyio/statelessdb/pkg/requests"
	"github.com/hyperifyio/statelessdb/pkg/states"
)

// ApiEventHandler is called to implement GET /api/v1/events which implements an HTTP long polling end point
func ApiEventHandler(
	bus events.EventBus[uuid.UUID, interface{}],
	timeoutTime,
	eventExpirationTime,
	intervalTime time.Duration,
) requests.ApiRequestHandlerFunc[*states.ComputeState, *requests.ComputeRequest] {

	manager := events.NewEventManager(
		bus,
		eventExpirationTime,
		intervalTime,
		EventSubscribersBufferSize,
		InternalEventManagerBufferSize,
		EventTriggerIntervalTimeMs*time.Second,
		EventTriggerRetryTime,
	)

	return func(r *requests.ComputeRequest, state *states.ComputeState) (*states.ComputeState, error) {

		if state == nil {
			return nil, ErrNoStateProvided
		}

		if err := state.Initialize(); err != nil {
			return nil, err
		}

		now := states.NewTimeNow()
		r.Received = now

		if bufferedEvents := manager.GetBufferedEvents(state.Id, state.Updated); len(bufferedEvents) > 0 {
			state.AddEvent(bufferedEvents...)
			state.Updated = states.NewTimeNow()
			return state, nil
		}

		eventChannel := make(chan int64, EventBufferSize)
		manager.Subscribe(state.Id, eventChannel)
		defer manager.Unsubscribe(state.Id, eventChannel)

		timeout := time.After(timeoutTime)

	EventLoop:
		for {
			select {
			case <-eventChannel:
				break EventLoop

			case <-timeout:
				break EventLoop
			}
		}

		if bufferedEvents := manager.GetBufferedEvents(state.Id, r.Received); len(bufferedEvents) > 0 {
			state.AddEvent(bufferedEvents...)
		}
		state.Updated = states.NewTimeNow()
		return state, nil
	}
}

// NewEventResponseDTO handles internal events to public DTO
func NewEventResponseDTO(bus events.EventBus[uuid.UUID, interface{}]) requests.CreateResponseFunc[*states.ComputeState] {
	return func(state *states.ComputeState, private string) interface{} {
		now := state.Updated
		evList := state.Events()
		list := make([]*dtos.EventDTO, len(evList))
		if evList == nil {
			return dtos.NewEventListDTO(now, list, private)
		}

		for i, v := range evList {
			list[i] = dtos.NewEventDTO(
				v.Type,
				v.Data,
				v.Created,
			)
		}

		return dtos.NewEventListDTO(now, list, private)
	}
}
