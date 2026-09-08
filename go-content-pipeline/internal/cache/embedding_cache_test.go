package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestEmbeddingCache_HitAndMiss(t *testing.T) {
	cache := NewEmbeddingCache("/tmp/test_cache.json", "bge-base")

	vec := []float32{0.1, 0.2, 0.3}
	cache.Put("post1", "hash1", vec)

	// Hit
	got, ok := cache.Get("post1", "hash1")
	if !ok || len(got) != 3 || got[0] != 0.1 {
		t.Fatalf("expected hit, got %v, ok=%v", got, ok)
	}

	// Miss - stale hash
	got, ok = cache.Get("post1", "hash2")
	if ok || got != nil {
		t.Fatalf("expected miss for different hash, got %v, ok=%v", got, ok)
	}

	// Miss - nonexistent post
	got, ok = cache.Get("post2", "hash1")
	if ok || got != nil {
		t.Fatalf("expected miss for nonexistent post, got %v, ok=%v", got, ok)
	}

	stats := cache.GetStats()
	if stats.Hits != 1 {
		t.Errorf("expected 1 hit, got %d", stats.Hits)
	}
	if stats.Misses != 2 {
		t.Errorf("expected 2 misses, got %d", stats.Misses)
	}
	if stats.TotalEntries != 1 {
		t.Errorf("expected 1 total entry, got %d", stats.TotalEntries)
	}
}

func TestEmbeddingCache_PersistenceAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "embeddings.json")

	cache := NewEmbeddingCache(cachePath, "bge-base")
	cache.Put("post-a", "hash-a", []float32{1.0, 2.0})
	cache.Put("post-b", "hash-b", []float32{3.0, 4.0})

	if err := cache.Persist(); err != nil {
		t.Fatalf("Persist failed: %v", err)
	}

	// New cache instance reading same file
	cache2 := NewEmbeddingCache(cachePath, "bge-base")
	if err := cache2.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cache2.Count() != 2 {
		t.Fatalf("expected 2 entries in loaded cache, got %d", cache2.Count())
	}

	vec, ok := cache2.Get("post-a", "hash-a")
	if !ok || len(vec) != 2 || vec[0] != 1.0 {
		t.Fatalf("unexpected value for post-a: %v, ok=%v", vec, ok)
	}
}

func TestEmbeddingCache_ModelInvalidation(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "embeddings.json")

	cache := NewEmbeddingCache(cachePath, "model-v1")
	cache.Put("post-a", "hash-a", []float32{1.0})
	cache.Put("post-b", "hash-b", []float32{2.0})
	if err := cache.Persist(); err != nil {
		t.Fatalf("Persist failed: %v", err)
	}

	// Loading with different model name should invalidate
	cacheNewModel := NewEmbeddingCache(cachePath, "model-v2")
	if err := cacheNewModel.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cacheNewModel.Count() != 0 {
		t.Errorf("expected 0 entries after model mismatch, got %d", cacheNewModel.Count())
	}

	// Explicit InvalidateModel
	cache.Put("post-c", "hash-c", []float32{3.0})
	invalidated := cache.InvalidateModel("model-v1", "model-v3")
	if invalidated != 3 {
		t.Errorf("expected 3 invalidated entries, got %d", invalidated)
	}
	if cache.Count() != 0 {
		t.Errorf("expected count 0 after InvalidateModel, got %d", cache.Count())
	}
}

func TestEmbeddingCache_CorruptedFileRecovery(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "corrupted.json")

	// Write invalid JSON
	if err := os.WriteFile(cachePath, []byte("{invalid-json-content..."), 0644); err != nil {
		t.Fatalf("failed to write corrupted file: %v", err)
	}

	cache := NewEmbeddingCache(cachePath, "bge-base")
	if err := cache.Load(); err != nil {
		t.Fatalf("Load should recover from corrupted file, but returned error: %v", err)
	}

	if cache.Count() != 0 {
		t.Errorf("expected 0 entries for corrupted file, got %d", cache.Count())
	}

	// Should be able to put and persist cleanly
	cache.Put("post-1", "hash-1", []float32{1.0})
	if err := cache.Persist(); err != nil {
		t.Fatalf("Persist after corruption recovery failed: %v", err)
	}
}

func TestEmbeddingCache_Concurrency(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "concurrent_cache.json")
	cache := NewEmbeddingCache(cachePath, "bge-base")

	var wg sync.WaitGroup
	numWorkers := 10
	opsPerWorker := 50

	for i := 0; i < numWorkers; i++ {
		workerID := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerWorker; j++ {
				post := fmt.Sprintf("post-%d-%d", workerID, j)
				hash := fmt.Sprintf("hash-%d", j)
				cache.Put(post, hash, []float32{float32(j)})
				cache.Get(post, hash)
				cache.Get(post, "wrong-hash")
				_ = cache.GetStats()
			}
		}()
	}

	wg.Wait()

	if cache.Count() != numWorkers*opsPerWorker {
		t.Errorf("expected %d entries, got %d", numWorkers*opsPerWorker, cache.Count())
	}

	if err := cache.Persist(); err != nil {
		t.Fatalf("concurrent persist failed: %v", err)
	}
}
