// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package requests

import (
	"github.com/hyperifyio/statelessdb/pkg/encodings"
)

type ApiRequestHandlerFunc[T interface{}, R Request] func(r R, state T) (T, error)

type ApiBytesRequestHandlerFunc func(body []byte) (encodings.SerializerState, error)

type RequestManager[T interface{}, R Request, D interface{}] interface {
	DecodeRequest(body []byte) (R, error)
	DecryptState(private string) (T, error)
	EncryptState(state T) (string, error)
	HandleWith(handleRequest ApiRequestHandlerFunc[T, R]) *RequestResponseManager[T, R, D]
}
