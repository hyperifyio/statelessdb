// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package workers_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/hyperifyio/statelessdb/pkg/workers"
)

func TestWorkerPool_BasicFunctionality(t *testing.T) {
	ctx := context.Background()
	pool := workers.NewPool[int](ctx, 10)

	var mu sync.Mutex
	processedJobs := make([]int, 0)
	var wg sync.WaitGroup

	jobHandler := func(job int) {
		defer wg.Done()
		mu.Lock()
		processedJobs = append(processedJobs, job)
		mu.Unlock()
	}

	pool.Start(3, jobHandler)

	// Publish some jobs
	numJobs := 5
	wg.Add(numJobs)
	for i := 0; i < numJobs; i++ {
		err := pool.Publish(i)
		if err != nil {
			t.Fatalf("Failed to publish job %d: %v", i, err)
		}
	}

	// Wait for all jobs to be processed
	wg.Wait()
	pool.Stop()

	// Check that all jobs were processed
	mu.Lock()
	defer mu.Unlock()
	if len(processedJobs) != numJobs {
		t.Errorf("Expected %d jobs processed, got %d", numJobs, len(processedJobs))
	}

	// Create a map to track which jobs have been processed
	expectedJobs := make(map[int]struct{})
	for i := 0; i < numJobs; i++ {
		expectedJobs[i] = struct{}{}
	}

	for _, job := range processedJobs {
		delete(expectedJobs, job)
	}

	if len(expectedJobs) != 0 {
		t.Errorf("Not all jobs were processed. Missing jobs: %v", expectedJobs)
	}
}

func TestWorkerPool_PublishAfterStop(t *testing.T) {
	ctx := context.Background()
	pool := workers.NewPool[int](ctx, 10)

	jobHandler := func(job int) {
		// Do nothing
	}

	pool.Start(1, jobHandler)
	pool.Stop()

	err := pool.Publish(1)
	if err == nil {
		t.Error("Expected error when publishing to a stopped pool, but got nil")
	}
	if !errors.Is(err, workers.ErrPoolClosed) {
		t.Errorf("Expected ErrPoolClosed when publishing to a stopped pool, but got %v", err)
	}
}

func TestWorkerPool_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	pool := workers.NewPool[int](ctx, 10)

	var wg sync.WaitGroup
	jobHandler := func(job int) {
		defer wg.Done()
		// Simulate some work
		time.Sleep(100 * time.Millisecond)
	}

	pool.Start(2, jobHandler)

	// Publish some jobs
	numJobs := 5
	wg.Add(numJobs)
	for i := 0; i < numJobs; i++ {
		err := pool.Publish(i)
		if err != nil {
			t.Fatalf("Failed to publish job %d: %v", i, err)
		}
	}

	// Cancel the context before all jobs are processed
	time.Sleep(200 * time.Millisecond)
	cancel()

	// Wait for workers to stop
	pool.Stop()

	// Check if pool is stopped
	err := pool.Publish(100)
	if err == nil {
		t.Error("Expected error when publishing after context cancellation, but got nil")
	}
	if !errors.Is(err, workers.ErrPoolClosed) {
		t.Errorf("Expected ErrPoolClosed when publishing after context cancellation, but got %v", err)
	}
}

func TestWorkerPool_StopWaitsForWorkers(t *testing.T) {
	ctx := context.Background()
	pool := workers.NewPool[int](ctx, 10)

	jobStarted := make(chan struct{})
	jobFinished := make(chan struct{})

	jobHandler := func(job int) {
		close(jobStarted) // Signal that job has started
		// Simulate some work
		time.Sleep(500 * time.Millisecond)
		close(jobFinished) // Signal that job has finished
	}

	if err := pool.Start(1, jobHandler); err != nil {
		t.Fatalf("Failed to start the pool: err=%v", err)
	}

	err := pool.Publish(1)
	if err != nil {
		t.Fatalf("Failed to publish job: %v", err)
	}

	// Wait for job to start
	select {
	case <-jobStarted:
		// Job has started
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Job did not start in time")
	}

	stopDone := make(chan struct{})
	go func() {
		pool.Stop()
		close(stopDone)
	}()

	select {
	case <-stopDone:
		t.Error("Stop() returned before worker finished")
	case <-time.After(100 * time.Millisecond):
		// Stop() has not returned yet, which is expected
	}

	// Wait for job to finish
	select {
	case <-jobFinished:
		// Job finished
	case <-time.After(1 * time.Second):
		t.Fatal("Job did not finish in time")
	}

	// Now Stop() should return
	select {
	case <-stopDone:
		// Stop() returned after job finished
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Stop() did not return after job finished")
	}
}

func TestWorkerPool_TryPublish(t *testing.T) {

	ctx := context.Background()
	pool := workers.NewPool[int](ctx, 2) // Small buffer size for testing
	if err := pool.Start(1, func(job int) {
		time.Sleep(100 * time.Millisecond) // Simulate work
	}); err != nil {
		t.Fatalf("Failed to start the pool: err=%v", err)
	}

	defer pool.Stop()

	// Publish jobs to fill the queue
	for i := 0; i < 2; i++ {
		success, err := pool.TryPublish(i)
		if !success || err != nil {
			t.Fatalf("Failed to publish job %d: success=%v, err=%v", i, success, err)
		}
	}

	// Attempt to publish another job; should fail because the queue is full
	success, err := pool.TryPublish(2)
	if success {
		t.Errorf("Expected TryPublish to return false when queue is full, but got success=%v", success)
	}
	if err != nil {
		t.Errorf("Expected no error when queue is full, but got err=%v", err)
	}
}
