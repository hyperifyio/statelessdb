// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package encodings

type Unserializer[T interface{}] interface {
	Unserialize(serialized []byte, out T) error
}

type NewUnserializer[T interface{}] func(name string) Unserializer[T]
