// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encryptions

import (
	"bytes"
	"encoding/gob"
	"sync"
)

var gobEncoderPoolState = sync.Pool{
	New: func() interface{} {
		var buf *bytes.Buffer = new(bytes.Buffer)
		return &GobEncoderState{
			buf,
		}
	},
}

type GobEncoderState struct {
	buffer *bytes.Buffer
}

var _ SerializerState = &GobEncoderState{}

func (e *GobEncoderState) Release() {
	e.buffer.Reset()
	gobEncoderPoolState.Put(e)
}

func (e *GobEncoderState) Bytes() []byte {
	return e.buffer.Bytes()
}

func getGobEncoderState() *GobEncoderState {
	return gobEncoderPoolState.Get().(*GobEncoderState)
}

// GobSerializer manages a pool of gob.Encoder instances.
type GobSerializer[T interface{}] struct {
}

var _ Serializer[string] = &GobSerializer[string]{}

// NewGobSerializer initializes a new GobSerializer with a gob.Encoder pool.
func NewGobSerializer[T interface{}](name string) *GobSerializer[T] {
	RegisterGobTypeOnce[T](name)
	return &GobSerializer[T]{}
}

// Serialize serializes the given data using a reusable gob.Encoder.
// It returns the serialized bytes or an error.
func (s *GobSerializer[T]) Serialize(data T) (SerializerState, error) {
	state := getGobEncoderState()
	encoder := gob.NewEncoder(state.buffer)
	if err := encoder.Encode(data); err != nil {
		state.Release()
		return nil, err
	}
	return state, nil
}
