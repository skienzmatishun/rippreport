package processor

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool_ExecutionAndShutdown(t *testing.T) {
	wp := NewWorkerPool(4, 20)
	wp.Start()

	var counter atomic.Int32
	numTasks := 20

	for i := 0; i < numTasks; i++ {
		err := wp.Submit(func(ctx context.Context) error {
			time.Sleep(5 * time.Millisecond)
			counter.Add(1)
			return nil
		})
		if err != nil {
			t.Fatalf("Submit failed: %v", err)
		}
	}

	wp.Shutdown()

	if counter.Load() != int32(numTasks) {
		t.Errorf("expected %d completed tasks, got %d", numTasks, counter.Load())
	}
}

func TestWorkerPool_ErrorCollection(t *testing.T) {
	wp := NewWorkerPool(2, 10)
	wp.Start()

	_ = wp.Submit(func(ctx context.Context) error {
		return errors.New("task error 1")
	})
	_ = wp.Submit(func(ctx context.Context) error {
		return errors.New("task error 2")
	})

	wp.Shutdown()

	var collected []error
	for err := range wp.Errors() {
		collected = append(collected, err)
	}

	if len(collected) != 2 {
		t.Errorf("expected 2 errors, got %d", len(collected))
	}
}

func TestMemoryManager(t *testing.T) {
	mm := NewMemoryManager(5)

	// Should not trigger on 4
	if mm.CheckAndGC(4) {
		t.Errorf("should not GC before interval")
	}

	// Should trigger on 5
	if !mm.CheckAndGC(5) {
		t.Errorf("expected GC to trigger at count 5")
	}

	// Should not trigger again immediately
	if mm.CheckAndGC(6) {
		t.Errorf("should not GC again on count 6")
	}

	stats := mm.Stats()
	if stats.SysBytes == 0 {
		t.Errorf("expected non-zero memory stats")
	}
}
