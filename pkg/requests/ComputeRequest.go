// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package requests

// ComputeRequest defines a structure of the request body to the compute server
type ComputeRequest struct {
	Received    int64                  `json:"received,omitempty"` // Received is time when this request was received.
	Public      map[string]interface{} `json:"public,omitempty"`   // Public contains public properties for a new resource
	PrivateData string                 `json:"private,omitempty"`  // Private contains the private property from previous request. If omitted, a new resource is initialized.
}

var _ Request = &ComputeRequest{}

func NewComputeRequest(
	received int64,
	public map[string]interface{},
	private string,
) *ComputeRequest {
	return &ComputeRequest{
		Received:    received,
		Public:      public,
		PrivateData: private,
	}
}

// Private returns encrypted state data
func (r *ComputeRequest) Private() string {
	return r.PrivateData
}
