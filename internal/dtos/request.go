// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package dtos

// ComputeRequestDTO defines the structure of the request body to the compute server
type ComputeRequestDTO struct {

	// Payload The state of compute from previous request.
	// If not defined, a new compute resource is initialized.
	Payload *ComputeStateDTO `json:"payload,omitempty"`
}
