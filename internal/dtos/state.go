// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package dtos

// ComputeStateDTO struct defines the structure of the response DTO
type ComputeStateDTO struct {
	Id      string                 `json:"id"`      // Id identifies the resource
	Owner   string                 `json:"owner"`   // Owner is the owner of the resource
	Created string                 `json:"created"` // Created is the time this resource was created
	Updated string                 `json:"updated"` // Updated is the time this resource was updated last time
	Public  map[string]interface{} `json:"public"`  // Public is public properties of the resource
	Private string                 `json:"private"` // Private is the internal encrypted types.ComputeState
}
