package processor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestProgressTracker_BasicOperations(t *testing.T) {
	pt := NewProgressTracker("")

	if err := pt.MarkProcessed("post-1"); err != nil {
		t.Fatalf("MarkProcessed failed: %v", err)
	}
	if !pt.IsProcessed("post-1") {
		t.Errorf("expected post-1 to be processed")
	}

	if err := pt.MarkFailed("post-2", errors.New("timeout error")); err != nil {
		t.Fatalf("MarkFailed failed: %v", err)
	}
	if pt.IsProcessed("post-2") {
		t.Errorf("failed post should not report IsProcessed=true")
	}

	failed := pt.GetFailed()
	if len(failed) != 1 || failed[0].Slug != "post-2" || failed[0].Error != "timeout error" {
		t.Errorf("unexpected failed posts: %+v", failed)
	}

	if pt.ProcessedCount() != 1 {
		t.Errorf("expected ProcessedCount=1, got %d", pt.ProcessedCount())
	}
}

func TestProgressTracker_ResumeFiltering(t *testing.T) {
	pt := NewProgressTracker("")

	_ = pt.MarkProcessed("done-1")
	_ = pt.MarkProcessed("done-2")
	_ = pt.MarkFailed("failed-1", errors.New("temporary failure"))

	allSlugs := []string{"done-1", "todo-1", "failed-1", "done-2", "todo-2"}
	unprocessed := pt.FilterUnprocessedSlugs(allSlugs)

	expected := []string{"todo-1", "failed-1", "todo-2"}
	if len(unprocessed) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, unprocessed)
	}
	for i, s := range expected {
		if unprocessed[i] != s {
			t.Errorf("at index %d: expected %s, got %s", i, s, unprocessed[i])
		}
	}
}

func TestProgressTracker_PersistenceAndRecovery(t *testing.T) {
	tmpDir := t.TempDir()
	progressPath := filepath.Join(tmpDir, "progress.json")

	pt := NewProgressTracker(progressPath)
	for i := 0; i < 15; i++ {
		_ = pt.MarkProcessed(fmt.Sprintf("post-%d", i))
	}

	// Should have auto-persisted at least once (at 10 items)
	if _, err := os.Stat(progressPath); os.IsNotExist(err) {
		t.Fatalf("progress file should have been auto-persisted")
	}

	// Explicitly persist remaining items
	if err := pt.Persist(); err != nil {
		t.Fatalf("Persist failed: %v", err)
	}

	// Load into new instance
	pt2 := NewProgressTracker(progressPath)
	if err := pt2.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if pt2.ProcessedCount() != 15 {
		t.Errorf("expected 15 processed posts in loaded tracker, got %d", pt2.ProcessedCount())
	}
}

func TestProgressTracker_CorruptedFileRecovery(t *testing.T) {
	tmpDir := t.TempDir()
	progressPath := filepath.Join(tmpDir, "corrupted_progress.json")

	if err := os.WriteFile(progressPath, []byte("NOT_VALID_JSON"), 0644); err != nil {
		t.Fatalf("failed to write corrupted file: %v", err)
	}

	pt := NewProgressTracker(progressPath)
	if err := pt.Load(); err != nil {
		t.Fatalf("Load should recover from corrupted file, but got: %v", err)
	}

	if pt.ProcessedCount() != 0 {
		t.Errorf("expected 0 entries for corrupted progress file, got %d", pt.ProcessedCount())
	}
}

func TestProgressTracker_Concurrency(t *testing.T) {
	pt := NewProgressTracker("")

	var wg sync.WaitGroup
	workers := 10
	postsPerWorker := 30

	for i := 0; i < workers; i++ {
		workerID := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < postsPerWorker; j++ {
				slug := fmt.Sprintf("post-%d-%d", workerID, j)
				if j%2 == 0 {
					_ = pt.MarkProcessed(slug)
				} else {
					_ = pt.MarkFailed(slug, fmt.Errorf("error %d", j))
				}
				_ = pt.IsProcessed(slug)
				_ = pt.GetFailed()
			}
		}()
	}

	wg.Wait()

	expectedProcessed := workers * (postsPerWorker / 2)
	if pt.ProcessedCount() != expectedProcessed {
		t.Errorf("expected %d processed posts, got %d", expectedProcessed, pt.ProcessedCount())
	}
}
