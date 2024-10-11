// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encodings

import (
	"bytes"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var jsonEncoderPoolState = sync.Pool{
	New: func() interface{} {
		buf := getBytesBuffer()
		return &JsonEncoderState{
			buf,
			json.NewEncoder(buf),
		}
	},
}

func GetJsonEncoderState() *JsonEncoderState {
	return jsonEncoderPoolState.Get().(*JsonEncoderState)
}

type JsonEncoderState struct {
	buffer  *bytes.Buffer
	Encoder *jsoniter.Encoder
}

var _ SerializerState = &JsonEncoderState{}

func (e *JsonEncoderState) Release() {
	e.buffer.Reset()
	jsonEncoderPoolState.Put(e)
}

func (e *JsonEncoderState) Bytes() []byte {
	return e.buffer.Bytes()
}

// JsonSerializer manages a pool of json.Encoder instances.
type JsonSerializer[T interface{}] struct {
}

var _ Serializer[string] = &JsonSerializer[string]{}

// NewJsonSerializer initializes a new JsonSerializer with a json.Encoder pool.
func NewJsonSerializer[T interface{}](name string) *JsonSerializer[T] {
	return &JsonSerializer[T]{}
}

// Serialize serializes the given data using a reusable json.Encoder.
// It returns the serialized bytes or an error.
func (s *JsonSerializer[T]) Serialize(data T) (SerializerState, error) {
	state := GetJsonEncoderState()
	if err := state.Encoder.Encode(data); err != nil {
		state.Release()
		return nil, err
	}
	return state, nil
}
