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
	ErrPoolClosed                   = errors.New("cannot publish job: pool is stopped")
	ErrCannotStartPool              = errors.New("cannot start the pool: it is already running or shutting down")
	ErrCannotStopPool               = errors.New("cannot stop the pool: it is not running")
	ErrWorkerFuncAlreadyInitialized = errors.New("worker function was already initialized")
)

const (
	StateStopped      uint32 = iota // Initial state, workers not started
	StateRunning                    // Workers are running
	StateShuttingDown               // Workers are shutting down
)

// Pool defines the interface for a worker pool.
type Pool[T any] interface {

	// Start starts workers on pool.
	Start(workers int, f func(T)) error

	// Stop stops the worker pool.
	Stop() error

	// Publish is a potentially blocking thread-safe method to add a job to the
	// worker pool. It will block if no workers are available.
	Publish(job T) error

	// TryPublish is a non-blocking thread-safe method to add a job to the
	// worker pool. Instead of blocking, it will fail instantly, if no workers
	// are available and return value will be ErrWorkersBusy in this case.
	TryPublish(job T) (bool, error)

	// TryStealWork allows an idle worker to attempt processing a single job from
	// the pool. It returns true if a job was processed, false otherwise. Returns
	// an error ErrPoolClosed if pool is not started.
	TryStealWork() (bool, error)
}

// WorkerPool is an implementation of the Pool interface.
// It manages a pool of worker goroutines to process jobs of type T.
type WorkerPool[T any] struct {
	jobs       chan T          // jobs is a channel for sending jobs to workers.
	wg         sync.WaitGroup  // wg is WaitGroup to wait for all workers to finish.
	ctx        context.Context // ctx is context for cancellation.
	state      uint32          // state is atomic variable representing the pool's state.
	bufferSize int             // bufferSize is the buffer size for jobs channel
	workerFunc func(T)         // workerFunc is a function which processes a job from the pool
}

var _ Pool[int] = &WorkerPool[int]{} // Implements Pool[int]

func NewPool[T any](
	ctx context.Context,
	bufferSize int,
) *WorkerPool[T] {
	log.Debugf("NewPool: Creating a worker pool with buffer %d", bufferSize)
	return &WorkerPool[T]{
		ctx:        ctx,
		bufferSize: bufferSize,
	}
}

// Start starts some workers to process jobs
func (m *WorkerPool[T]) Start(workers int, f func(T)) error {
	if !atomic.CompareAndSwapUint32(&m.state, StateStopped, StateRunning) {
		return ErrCannotStartPool
	}
	if m.workerFunc != nil {
		return ErrWorkerFuncAlreadyInitialized
	}

	m.workerFunc = f
	m.jobs = make(chan T, m.bufferSize)
	m.wg = sync.WaitGroup{}

	log.Debugf("Start: Starting workers on the pool (%d workers)", workers)
	for i := 0; i < workers; i++ {
		m.wg.Add(1)
		go m.worker()
	}
	go func() {
		<-m.ctx.Done()
		err := m.Stop()
		if err != nil {
			log.Warnf("Start: Warning! Pool stop failed: %v", err)
		}
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
		log.Debugf("Stop: Waiting for workers to stop")
		m.wg.Wait()
		atomic.StoreUint32(&m.state, StateStopped)
		log.Debugf("Stop: All workers have stopped")
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

// Publish is a potentially blocking thread-safe method to add a job to the
// worker pool. It will block if no workers are available.
func (m *WorkerPool[T]) Publish(job T) error {
	if atomic.LoadUint32(&m.state) != StateRunning {
		return ErrPoolClosed
	}
	select {
	case m.jobs <- job:
		log.Debugf("Publish: Published a job")
		return nil
	case <-m.ctx.Done():
		return ErrPoolClosed
	}
}

// TryPublish is a non-blocking thread-safe method to add a job to the
// worker pool. Instead of blocking, it will fail instantly, if no workers
// are available and return value will be ErrWorkersBusy in this case.
func (m *WorkerPool[T]) TryPublish(job T) (bool, error) {
	if atomic.LoadUint32(&m.state) != StateRunning {
		return false, ErrPoolClosed
	}
	select {
	case m.jobs <- job:
		log.Debugf("TryPublish: Published a job")
		return true, nil
	default:
		log.Debugf("TryPublish: All workers busy and queue full")
		return false, nil
	}
}

func (m *WorkerPool[T]) TryStealWork() (bool, error) {
	if atomic.LoadUint32(&m.state) != StateRunning {
		return false, ErrPoolClosed
	}
	select {
	case job, ok := <-m.jobs:
		if !ok {
			log.Debugf("StealWork: Worker job channel closed.")
			return false, ErrPoolClosed
		}
		log.Debugf("StealWork: Stole and started working on a job.")
		m.workerFunc(job)
		log.Debugf("StealWork: Stole and processed a job.")
		return true, nil
	default:
		log.Debugf("StealWork: No work available to steal.")
		return false, nil
	}
}

func (m *WorkerPool[T]) worker() {
	defer m.wg.Done()
	if m.workerFunc == nil {
		log.Errorf("Worker: ERROR: No worker function initialized. Worker stopped.")
		return
	}
	for {
		select {
		case job, ok := <-m.jobs:
			if !ok {
				log.Debugf("Worker: Shutting down. Worker job channel closed.")
				return
			}
			log.Debugf("Worker: Started working on a job.")
			m.workerFunc(job)
			log.Debugf("Worker: Processed a job.")
		case <-m.ctx.Done():
			log.Debugf("Worker: Shutting down. Worker pool context closed.")
			return
		}
	}
}
