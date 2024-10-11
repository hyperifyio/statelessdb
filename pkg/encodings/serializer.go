// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encodings

type SerializerState interface {
	Release()
	Bytes() []byte
}

type Serializer[T interface{}] interface {
	Serialize(data T) (SerializerState, error)
}

type NewSerializer[T interface{}] func(name string) Serializer[T]
