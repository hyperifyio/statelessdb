// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package requests

type ResponseManager interface {
	ProcessBytes(body []byte) (interface{}, error)
}

type CreateResponseFunc[T interface{}, D interface{}] func(state T, private string) D

type RequestResponseManager[T interface{}, R Request, D interface{}] struct {
	parent         *EncryptedRequestManager[T, R, D]
	handleRequest  ApiRequestHandlerFunc[T, R]
	handleResponse CreateResponseFunc[T, D]
}

var _ ResponseManager = &RequestResponseManager[any, Request, any]{}

// ProcessBytes decodes, decrypts, processes, and encrypts results for a request
func (r *RequestResponseManager[T, R, D]) ProcessBytes(body []byte) (interface{}, error) {

	req, err := r.parent.DecodeRequest(body)
	if err != nil {
		var dto interface{}
		return dto, err
	}

	state, err := r.parent.DecryptState(req.Private())
	if err != nil {
		var dto interface{}
		return dto, err
	}

	state, err = r.handleRequest(req, state)
	if err != nil {
		var dto interface{}
		return dto, err
	}

	private, err := r.parent.EncryptState(state)
	if err != nil {
		var dto interface{}
		return dto, err
	}

	if r.handleResponse != nil {
		return r.handleResponse(state, private), nil
	}

	var dto interface{}
	return dto, nil
}

// WithResponse configures a response DTO handler
func (r *RequestResponseManager[T, R, D]) WithResponse(handler CreateResponseFunc[T, D]) *RequestResponseManager[T, R, D] {
	r.handleResponse = handler
	return r
}
