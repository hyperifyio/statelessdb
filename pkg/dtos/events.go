// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package dtos

import (
	"github.com/google/uuid"
	"github.com/hyperifyio/statelessdb/internal/helpers"
)

// EventDTO struct defines DTO for event list
type EventDTO struct {
	Id      string      `json:"id"`      // Id identifies the resource which was listened
	Data    interface{} `json:"data"`    // Data is information provided with the event
	Created string      `json:"created"` // Created is the time when this event was received
}

func NewEventDTO(
	id uuid.UUID,
	data interface{},
	created int64,
) *EventDTO {
	return &EventDTO{
		Created: helpers.MillisToISO(created),
		Id:      id.String(),
		Data:    data,
	}
}

// EventListDTO struct defines DTO for event list
type EventListDTO struct {
	Created string      `json:"created"` // Created is the time when this event list was sent. You can use this to request more events after this time.
	Payload []*EventDTO `json:"payload"` // Payload contains all events received
	Private string      `json:"private"` // Private can be used to request next set of events. It contains information required to know when to
}

func NewEventListDTO(
	created int64,
	payload []*EventDTO,
	private string,
) *EventListDTO {
	return &EventListDTO{
		Created: helpers.MillisToISO(created),
		Payload: payload,
		Private: private,
	}
}
