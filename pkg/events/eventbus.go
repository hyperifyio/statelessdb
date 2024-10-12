// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package events

type EventBus[T comparable, D interface{}] interface {
	Subscribe(eventType T, ch chan *Event[T, D])
	Unsubscribe(eventType T, ch chan *Event[T, D])
	Publish(event *Event[T, D])
}
