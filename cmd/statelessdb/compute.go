// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package main

import (
	"github.com/google/uuid"
	"github.com/hyperifyio/statelessdb/pkg/dtos"
	"github.com/hyperifyio/statelessdb/pkg/events"
	"github.com/hyperifyio/statelessdb/pkg/requests"
	"github.com/hyperifyio/statelessdb/pkg/states"
)

// ApiRequestHandler is called to implement POST /api/v1 which implements compute actions on a state
func ApiRequestHandler(bus events.EventBus[uuid.UUID, interface{}]) requests.ApiRequestHandlerFunc[*states.ComputeState, *requests.ComputeRequest] {
	return func(r *requests.ComputeRequest, state *states.ComputeState) (*states.ComputeState, error) {

		now := states.NewTimeNow()
		r.Received = now

		if state == nil {
			var private map[string]interface{}
			//private = make(map[string]interface{})
			state = states.NewComputeState(uuid.New(), uuid.New(), now, now, r.Public, private, nil)
		}

		if err := state.Initialize(); err != nil {
			return nil, err
		}

		state.Updated = now

		return state, nil
	}
}

func NewComputeResponseDTO(bus events.EventBus[uuid.UUID, interface{}]) requests.CreateResponseFunc[*states.ComputeState] {
	return func(state *states.ComputeState, private string) interface{} {
		dto := dtos.NewComputeResponseDTO(
			state.Id,
			state.Owner,
			state.Created,
			state.Updated,
			state.Public,
			private,
		)
		bus.Publish(events.NewEvent[uuid.UUID, interface{}](state.Id, dto, state.Updated))
		return dto
	}
}
