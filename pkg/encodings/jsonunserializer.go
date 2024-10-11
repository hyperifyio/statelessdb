// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encodings

import (
	"bytes"
	"statelessdb/pkg/errors"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var jsonDecoderPoolState = sync.Pool{
	New: func() interface{} {
		buf := getBytesBuffer()
		return &JsonDecoderState{
			buf,
			json.NewDecoder(buf),
		}
	},
}

func GetJsonDecoderState() *JsonDecoderState {
	return jsonDecoderPoolState.Get().(*JsonDecoderState)
}

type JsonDecoderState struct {
	buffer  *bytes.Buffer
	Decoder *jsoniter.Decoder
}

func (e *JsonDecoderState) Release() {
	e.buffer.Reset()
	jsonDecoderPoolState.Put(e)
}

// JsonUnserializer manages a pool of reusable json.Decoder instances.
type JsonUnserializer[T interface{}] struct {
}

var _ Unserializer[string] = &JsonUnserializer[string]{}

// NewJsonUnserializer initializes and returns a new JsonUnserializer.
func NewJsonUnserializer[T interface{}](name string) *JsonUnserializer[T] {
	return &JsonUnserializer[T]{}
}

// Unserialize decodes serialized data
func (dp *JsonUnserializer[T]) Unserialize(serialized []byte, out T) error {
	state := GetJsonDecoderState()
	defer state.Release()
	buf := state.buffer
	buf.Write(serialized)
	if err := state.Decoder.Decode(out); err != nil {
		log.Errorf("[JsonUnserializer.Unserialize]: json decode failed: %v", err)
		return errors.ErrDecryptDecodingJsonSerializationFailed
	}
	return nil
}
