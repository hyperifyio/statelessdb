// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package main

import (
	"time"

	"github.com/google/uuid"

	"statelessdb/pkg/dtos"
	"statelessdb/pkg/events"
	"statelessdb/pkg/requests"
	"statelessdb/pkg/states"
)

// ApiEventHandler is called to implement GET /api/v1/events which implements an HTTP long polling end point
func ApiEventHandler(bus *events.EventBus[uuid.UUID, interface{}]) requests.ApiRequestHandlerFunc[*states.ComputeState, *requests.ComputeRequest] {
	return func(r *requests.ComputeRequest, state *states.ComputeState) (*states.ComputeState, error) {

		if state == nil {
			return nil, ErrNoStateProvided
		}

		if err := state.Initialize(); err != nil {
			return nil, err
		}

		now := states.NewTimeNow()
		r.Received = now

		eventChannel := make(chan *events.Event[uuid.UUID, interface{}])

		bus.Subscribe(state.Id, eventChannel)
		defer bus.Unsubscribe(state.Id, eventChannel)

		received := make([]*events.Event[uuid.UUID, interface{}], 0)
		timeout := time.After(time.Second * 10)

		for {
			select {

			case event := <-eventChannel:
				received = append(received, event)
				state.Events = received
				state.Updated = r.Received
				return state, nil

			case <-timeout:
				state.Events = received
				state.Updated = r.Received
				return state, nil

			}
		}
	}
}

// NewEventResponseDTO handles internal events to public DTO
func NewEventResponseDTO(bus *events.EventBus[uuid.UUID, interface{}]) requests.CreateResponseFunc[*states.ComputeState] {
	return func(state *states.ComputeState, private string) interface{} {
		now := state.Updated
		evList := state.Events
		if evList == nil {
			return dtos.NewEventListDTO(now, nil)
		}

		list := make([]*dtos.EventDTO, len(evList))
		for i, v := range evList {
			list[i] = dtos.NewEventDTO(
				v.Type,
				v.Data,
				v.Created,
			)
		}
		return dtos.NewEventListDTO(now, list)
	}
}
