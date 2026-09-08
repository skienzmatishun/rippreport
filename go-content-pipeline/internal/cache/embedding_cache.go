package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

const (
	defaultCacheVersion = "1.0"
)

// EmbeddingCache provides thread-safe caching and persistence of embedding vectors.
//
// Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7, 3.8
type EmbeddingCache struct {
	mu        sync.RWMutex
	filePath  string
	modelName string
	version   string
	entries   map[string]*models.CacheEntry
	hits      atomic.Int64
	misses    atomic.Int64
}

// NewEmbeddingCache creates a new embedding cache instance.
func NewEmbeddingCache(filePath, modelName string) *EmbeddingCache {
	return &EmbeddingCache{
		filePath:  filePath,
		modelName: modelName,
		version:   defaultCacheVersion,
		entries:   make(map[string]*models.CacheEntry),
	}
}

// Get retrieves an embedding from the cache if it exists and matches the content hash.
//
// Requirements: 3.1, 3.2, 3.3
func (ec *EmbeddingCache) Get(postPath, contentHash string) ([]float32, bool) {
	ec.mu.RLock()
	entry, exists := ec.entries[postPath]
	ec.mu.RUnlock()

	if !exists {
		ec.misses.Add(1)
		return nil, false
	}

	if entry.ContentHash != contentHash {
		ec.misses.Add(1)
		return nil, false
	}

	ec.hits.Add(1)
	return entry.Embedding, true
}

// Put stores or updates an embedding in the cache.
//
// Requirements: 3.1, 3.2
func (ec *EmbeddingCache) Put(postPath, contentHash string, embedding []float32) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	ec.entries[postPath] = &models.CacheEntry{
		Embedding:   embedding,
		ContentHash: contentHash,
		ModelName:   ec.modelName,
		Timestamp:   time.Now(),
	}
}

// GetStats returns current cache hit/miss statistics.
//
// Requirements: 3.6
func (ec *EmbeddingCache) GetStats() models.CacheStats {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	hits := ec.hits.Load()
	misses := ec.misses.Load()
	stats := models.CacheStats{
		TotalEntries: len(ec.entries),
		Hits:         hits,
		Misses:       misses,
	}
	stats.CalculateHitRate()
	return stats
}

// InvalidateModel clears all entries if the model changes and returns the count of invalidated entries.
//
// Requirements: 3.4, 3.5
func (ec *EmbeddingCache) InvalidateModel(oldModel, newModel string) int {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if oldModel == newModel && newModel == ec.modelName {
		return 0
	}

	count := len(ec.entries)
	ec.entries = make(map[string]*models.CacheEntry)
	ec.modelName = newModel
	return count
}

// Persist saves the cache to disk atomically using temporary file + rename.
//
// Requirements: 3.7
func (ec *EmbeddingCache) Persist() error {
	ec.mu.RLock()
	if ec.filePath == "" {
		ec.mu.RUnlock()
		return fmt.Errorf("cache file path is empty")
	}

	// Create snapshot of entries
	snapshot := make(map[string]*models.CacheEntry, len(ec.entries))
	for k, v := range ec.entries {
		snapshot[k] = v
	}
	modelName := ec.modelName
	version := ec.version
	filePath := ec.filePath
	ec.mu.RUnlock()

	cacheFile := models.CacheFile{
		Version:   version,
		ModelName: modelName,
		Entries:   snapshot,
		UpdatedAt: time.Now(),
	}

	data, err := json.MarshalIndent(cacheFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory %s: %w", dir, err)
	}

	tempFile, err := os.CreateTemp(dir, ".embedding_cache_*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temporary cache file: %w", err)
	}
	tempPath := tempFile.Name()

	success := false
	defer func() {
		if !success {
			_ = os.Remove(tempPath)
		}
	}()

	if _, err := tempFile.Write(data); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to write cache to temp file: %w", err)
	}

	if err := tempFile.Sync(); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to sync temp cache file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp cache file: %w", err)
	}

	if err := os.Rename(tempPath, filePath); err != nil {
		return fmt.Errorf("failed to atomically replace cache file %s: %w", filePath, err)
	}

	success = true
	return nil
}

// Load loads the cache from disk. If the cache file is corrupted, it resets to an empty cache.
// If the stored model does not match the configured model, all entries are invalidated.
//
// Requirements: 3.4, 3.5, 3.8
func (ec *EmbeddingCache) Load() error {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if ec.filePath == "" {
		return fmt.Errorf("cache file path is empty")
	}

	data, err := os.ReadFile(ec.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Clean start when cache does not yet exist
			ec.entries = make(map[string]*models.CacheEntry)
			return nil
		}
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	var cacheFile models.CacheFile
	if err := json.Unmarshal(data, &cacheFile); err != nil {
		// Requirement 3.8: Corrupted cache file starts fresh with empty cache
		ec.entries = make(map[string]*models.CacheEntry)
		return nil
	}

	// Requirement 3.5: Model mismatch invalidates all entries
	if cacheFile.ModelName != "" && ec.modelName != "" && cacheFile.ModelName != ec.modelName {
		ec.entries = make(map[string]*models.CacheEntry)
		return nil
	}

	if cacheFile.Entries == nil {
		ec.entries = make(map[string]*models.CacheEntry)
	} else {
		ec.entries = cacheFile.Entries
	}

	if cacheFile.ModelName != "" {
		ec.modelName = cacheFile.ModelName
	}
	if cacheFile.Version != "" {
		ec.version = cacheFile.Version
	}

	return nil
}

// Count returns the number of entries in the cache.
func (ec *EmbeddingCache) Count() int {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	return len(ec.entries)
}
