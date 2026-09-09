package processor

import (
	"context"
	"fmt"
	"sync"
)

// Task represents a unit of concurrent work.
type Task func(ctx context.Context) error

// WorkerPool manages a pool of worker goroutines executing submitted tasks.
//
// Requirements: 22.1, 22.2, 22.3, 22.4
type WorkerPool struct {
	numWorkers int
	tasks      chan Task
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	errors     chan error
	once       sync.Once
	closed     bool
	mu         sync.Mutex
}

// NewWorkerPool creates a new WorkerPool with the given number of workers and queue capacity.
func NewWorkerPool(numWorkers, queueSize int) *WorkerPool {
	if numWorkers <= 0 {
		numWorkers = 4
	}
	if queueSize <= 0 {
		queueSize = numWorkers * 2
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		numWorkers: numWorkers,
		tasks:      make(chan Task, queueSize),
		errors:     make(chan error, queueSize*2),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start launches worker goroutines.
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case task, ok := <-wp.tasks:
			if !ok {
				return
			}
			if err := task(wp.ctx); err != nil {
				select {
				case wp.errors <- err:
				default:
					// Drop or buffer if error channel full
				}
			}
		}
	}
}

// Submit queues a task for execution. Returns error if pool is closed or context cancelled.
func (wp *WorkerPool) Submit(task Task) error {
	wp.mu.Lock()
	if wp.closed {
		wp.mu.Unlock()
		return fmt.Errorf("worker pool is closed")
	}
	wp.mu.Unlock()

	select {
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	case wp.tasks <- task:
		return nil
	}
}

// Shutdown gracefully terminates all workers, waiting for queued tasks to complete.
func (wp *WorkerPool) Shutdown() {
	wp.once.Do(func() {
		wp.mu.Lock()
		wp.closed = true
		wp.mu.Unlock()

		close(wp.tasks)
		wp.wg.Wait()
		wp.cancel()
		close(wp.errors)
	})
}

// Errors returns the channel of collected task execution errors.
func (wp *WorkerPool) Errors() <-chan error {
	return wp.errors
}
