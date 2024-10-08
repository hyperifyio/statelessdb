// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package types

import (
	"sync"
)

type MemoryPool[T any] struct {
	pool sync.Pool
}

func NewMemoryPool[T any](newFunc func() T) *MemoryPool[T] {
	return &MemoryPool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return newFunc()
			},
		},
	}
}

func (p *MemoryPool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put adds an instance of T back into the pool
func (p *MemoryPool[T]) Put(item T) {
	p.pool.Put(item)
}
