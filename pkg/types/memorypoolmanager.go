// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package types

import "sync"

// MemoryPoolManager manages pools of different sizes
type MemoryPoolManager[T any] struct {
	pools   map[int]*MemoryPool[T]  // Pools keyed by slice length (e.g., 64, 128)
	mu      sync.Mutex              // Protects the map for concurrent access
	newFunc func(size int) func() T // Factory function for object creation function
}

// NewMemoryPoolManager creates a new manager for different pool sizes
func NewMemoryPoolManager[T any](
	newFunc func(size int) func() T,
) *MemoryPoolManager[T] {
	return &MemoryPoolManager[T]{
		pools:   make(map[int]*MemoryPool[T]),
		newFunc: newFunc,
	}
}

// Pool returns a memory pool for the specified size
func (m *MemoryPoolManager[T]) Pool(size int) *MemoryPool[T] {

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if a pool for this size already exists
	if pool, exists := m.pools[size]; exists {
		return pool
	}

	// Create a new pool for this size if it doesn't exist
	pool := NewMemoryPool(m.newFunc(size))
	m.pools[size] = pool
	return pool
}
