// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package events

type Event[T comparable, D interface{}] struct {
	Type    T
	Data    D
	Created int64
}

func NewEvent[T comparable, D interface{}](t T, d D, c int64) *Event[T, D] {
	return &Event[T, D]{t, d, c}
}
