// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package dtos

// ComputeRequestDTO defines the structure of the request body to the compute server
type ComputeRequestDTO struct {
	Public  map[string]interface{} `json:"public,omitempty"`  // Public contains public properties for a new resource
	Private string                 `json:"private,omitempty"` // Private contains the private property from previous request. If omitted, a new resource is initialized.
}

func NewComputeRequestDTO(
	public map[string]interface{},
	private string,
) *ComputeRequestDTO {
	return &ComputeRequestDTO{
		Public:  public,
		Private: private,
	}
}
