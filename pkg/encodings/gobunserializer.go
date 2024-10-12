// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encodings

import (
	"bytes"
	"encoding/gob"
	"statelessdb/pkg/errors"
	"sync"
)

var gobDecoderPoolState = sync.Pool{
	New: func() interface{} {
		buf := new(bytes.Buffer)
		return &GobDecoderState{
			buf,
			gob.NewDecoder(buf),
		}
	},
}

func getGobDecoderState() *GobDecoderState {
	return gobDecoderPoolState.Get().(*GobDecoderState)
}

type GobDecoderState struct {
	buffer  *bytes.Buffer
	decoder *gob.Decoder
}

func (e *GobDecoderState) Release() {
	e.buffer.Reset()
	gobDecoderPoolState.Put(e)
}

var _ Unserializer[string] = &GobUnserializer[string]{}

// GobUnserializer manages a pool of reusable gob.Decoder instances.
type GobUnserializer[T interface{}] struct {
}

// NewGobUnserializer initializes and returns a new GobUnserializer.
func NewGobUnserializer[T interface{}](name string) *GobUnserializer[T] {
	RegisterGobTypeOnce[T](name)
	return &GobUnserializer[T]{}
}

// Unserialize decodes serialized data
func (dp *GobUnserializer[T]) Unserialize(serialized []byte, out T) error {
	state := getGobDecoderState()
	defer state.Release()
	buf := state.buffer
	buf.Write(serialized)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(out); err != nil {
		log.Errorf("[GobUnserializer.Unserialize]: gob decode failed: %v", err)
		return errors.ErrDecryptDecodingGobSerializationFailed
	}
	return nil
}
