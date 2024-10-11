// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package states

import (
	"github.com/google/uuid"
	"statelessdb/pkg/events"

	"statelessdb/internal/helpers"
)

// ComputeState is the actual state of computation, which is encrypted to
// the Private field of dtos.ComputeResponseDTO. This is the state used by
// StatelessDB by default, but users may implement their own states.
type ComputeState struct {
	Id      uuid.UUID                               `json:"id"`      // ID identifies the object
	Owner   uuid.UUID                               `json:"owner"`   // Owner identifies the owner of the object
	Created int64                                   `json:"created"` // Created is the time when the object was created
	Updated int64                                   `json:"updated"` // Updated is the time when the object was updated
	Public  map[string]interface{}                  `json:"data"`    // Public contains public properties of the object
	Private map[string]interface{}                  `json:"private"` // Private contains unencrypted private properties of the object
	Events  []*events.Event[uuid.UUID, interface{}] // Events are special internal property for event handler
}

func NewComputeState(
	id, owner uuid.UUID,
	created, updated int64,
	public, private map[string]interface{},
	evList []*events.Event[uuid.UUID, interface{}],
) *ComputeState {
	return &ComputeState{
		Id:      id,
		Owner:   owner,
		Created: created,
		Updated: updated,
		Public:  public,
		Private: private,
		Events:  evList,
	}
}

// Equals returns true if both ComputeStates are equal state. This is not counting
// calculated data like internal caches, attackMaps, etc. which should be same
// anyway if this data is equal.
func (b *ComputeState) Equals(other *ComputeState) bool {
	if b == other {
		return true
	}
	if other == nil {
		return false
	}
	if b.Id != other.Id ||
		b.Owner != other.Owner ||
		b.Created != other.Created ||
		b.Updated != other.Updated {
		return false
	}
	if !helpers.CompareMaps(b.Public, other.Public) {
		return false
	}
	if !helpers.CompareMaps(b.Private, other.Private) {
		return false
	}
	return true
}

// Initialize initializes internal state. This may allocate internal memory!
func (g *ComputeState) Initialize() error {
	return nil
}
