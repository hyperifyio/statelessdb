// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package workers

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrPoolClosed      = errors.New("cannot publish job: pool is stopped")
	ErrCannotStartPool = errors.New("cannot start the pool: it is already running or shutting down")
	ErrCannotStopPool  = errors.New("cannot stop the pool: it is not running")
)

const (
	StateStopped      uint32 = iota // Initial state, workers not started
	StateRunning                    // Workers are running
	StateShuttingDown               // Workers are shutting down
)

// Pool defines the interface for a worker pool.
type Pool[T any] interface {
	Start(workers int, f func(T)) error // Start starts workers on pool.
	Stop() error                        // Stop stops the worker pool.
	Publish(job T) error                // Publish is a thread-safe method to add a job to the worker pool.
}

// WorkerPool is an implementation of the Pool interface.
// It manages a pool of worker goroutines to process jobs of type T.
type WorkerPool[T any] struct {
	jobs       chan T          // jobs is a channel for sending jobs to workers.
	wg         sync.WaitGroup  // wg is WaitGroup to wait for all workers to finish.
	ctx        context.Context // ctx is context for cancellation.
	state      uint32          // state is atomic variable representing the pool's state.
	bufferSize int             // bufferSize is the buffer size for jobs channel
}

var _ Pool[int] = &WorkerPool[int]{} // Implements Pool[int]

func NewPool[T any](
	ctx context.Context,
	bufferSize int,
) *WorkerPool[T] {
	log.Debugf("Creating a worker pool with buffer %d", bufferSize)
	return &WorkerPool[T]{
		ctx:        ctx,
		bufferSize: bufferSize,
	}
}

// Start starts some workers to process jobs
func (m *WorkerPool[T]) Start(
	workers int,
	f func(T),
) error {
	if !atomic.CompareAndSwapUint32(&m.state, StateStopped, StateRunning) {
		return ErrCannotStartPool
	}

	m.jobs = make(chan T, m.bufferSize)
	m.wg = sync.WaitGroup{}

	log.Debugf("Starting workers on the pool (%d workers)", workers)
	for i := 0; i < workers; i++ {
		m.wg.Add(1)
		go m.worker(f)
	}
	go func() {
		<-m.ctx.Done()
		m.Stop()
	}()
	return nil
}

// Stop stops the worker pool
func (m *WorkerPool[T]) Stop() error {
	state := atomic.LoadUint32(&m.state)
	if state == StateStopped {
		// Already stopped, nothing to do
		return nil
	}
	if atomic.CompareAndSwapUint32(&m.state, StateRunning, StateShuttingDown) {
		close(m.jobs)
		log.Debugf("Waiting for workers to stop")
		m.wg.Wait()
		atomic.StoreUint32(&m.state, StateStopped)
		log.Debugf("All workers have stopped")
		return nil
	} else if state == StateShuttingDown {
		// Already shutting down, wait for it to finish
		m.wg.Wait()
		return nil
	} else {
		// State is not running or shutting down
		return ErrCannotStopPool
	}
}

// Publish is a thread safe method to add jobs to the worker pool. If the job
// queue is full, this method will block.
func (m *WorkerPool[T]) Publish(job T) error {
	if atomic.LoadUint32(&m.state) != StateRunning {
		return ErrPoolClosed
	}
	select {
	case m.jobs <- job:
		log.Debugf("Published a job")
		return nil
	case <-m.ctx.Done():
		return ErrPoolClosed
	}
}

func (m *WorkerPool[T]) worker(f func(T)) {
	defer m.wg.Done()
	for {
		select {
		case job, ok := <-m.jobs:
			if !ok {
				log.Debugf("Shutting down. Worker job channel closed.")
				return
			}
			f(job)
		case <-m.ctx.Done():
			log.Debugf("Shutting down. Worker pool context closed.")
			return
		}
	}
}
