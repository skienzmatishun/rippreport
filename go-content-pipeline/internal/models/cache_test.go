package models

import (
	"encoding/json"
	"testing"
	"time"
)

// TestCacheEntry_JSONMarshaling tests JSON serialization of CacheEntry
func TestCacheEntry_JSONMarshaling(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)
	
	entry := &CacheEntry{
		Embedding:   []float32{0.1, 0.2, 0.3, 0.4},
		ContentHash: "abc123def456",
		ModelName:   "nomic-embed-text-v1.5",
		Timestamp:   testTime,
	}
	
	// Marshal to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("Failed to marshal CacheEntry: %v", err)
	}
	
	// Unmarshal back
	var decoded CacheEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal CacheEntry: %v", err)
	}
	
	// Verify all fields
	if len(decoded.Embedding) != len(entry.Embedding) {
		t.Errorf("Embedding length mismatch: got %d, want %d", len(decoded.Embedding), len(entry.Embedding))
	}
	for i, val := range decoded.Embedding {
		if val != entry.Embedding[i] {
			t.Errorf("Embedding[%d] mismatch: got %f, want %f", i, val, entry.Embedding[i])
		}
	}
	if decoded.ContentHash != entry.ContentHash {
		t.Errorf("ContentHash mismatch: got %s, want %s", decoded.ContentHash, entry.ContentHash)
	}
	if decoded.ModelName != entry.ModelName {
		t.Errorf("ModelName mismatch: got %s, want %s", decoded.ModelName, entry.ModelName)
	}
	if !decoded.Timestamp.Equal(entry.Timestamp) {
		t.Errorf("Timestamp mismatch: got %v, want %v", decoded.Timestamp, entry.Timestamp)
	}
}

// TestCacheFile_JSONMarshaling tests JSON serialization of CacheFile
func TestCacheFile_JSONMarshaling(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)
	
	cacheFile := &CacheFile{
		Version:   "1.0",
		ModelName: "nomic-embed-text-v1.5",
		UpdatedAt: testTime,
		Entries: map[string]*CacheEntry{
			"content/p/test-post/index.md": {
				Embedding:   []float32{0.1, 0.2, 0.3},
				ContentHash: "hash1",
				ModelName:   "nomic-embed-text-v1.5",
				Timestamp:   testTime,
			},
			"content/p/another-post/index.md": {
				Embedding:   []float32{0.4, 0.5, 0.6},
				ContentHash: "hash2",
				ModelName:   "nomic-embed-text-v1.5",
				Timestamp:   testTime,
			},
		},
	}
	
	// Marshal to JSON
	data, err := json.Marshal(cacheFile)
	if err != nil {
		t.Fatalf("Failed to marshal CacheFile: %v", err)
	}
	
	// Unmarshal back
	var decoded CacheFile
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal CacheFile: %v", err)
	}
	
	// Verify top-level fields
	if decoded.Version != cacheFile.Version {
		t.Errorf("Version mismatch: got %s, want %s", decoded.Version, cacheFile.Version)
	}
	if decoded.ModelName != cacheFile.ModelName {
		t.Errorf("ModelName mismatch: got %s, want %s", decoded.ModelName, cacheFile.ModelName)
	}
	if !decoded.UpdatedAt.Equal(cacheFile.UpdatedAt) {
		t.Errorf("UpdatedAt mismatch: got %v, want %v", decoded.UpdatedAt, cacheFile.UpdatedAt)
	}
	
	// Verify entries
	if len(decoded.Entries) != len(cacheFile.Entries) {
		t.Errorf("Entries count mismatch: got %d, want %d", len(decoded.Entries), len(cacheFile.Entries))
	}
	
	for path, originalEntry := range cacheFile.Entries {
		decodedEntry, exists := decoded.Entries[path]
		if !exists {
			t.Errorf("Missing entry for path: %s", path)
			continue
		}
		if decodedEntry.ContentHash != originalEntry.ContentHash {
			t.Errorf("Entry[%s].ContentHash mismatch: got %s, want %s", path, decodedEntry.ContentHash, originalEntry.ContentHash)
		}
	}
}

// TestCompressionEntry_JSONMarshaling tests JSON serialization of CompressionEntry
func TestCompressionEntry_JSONMarshaling(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)
	
	entry := &CompressionEntry{
		OriginalHash:  "hash123",
		Compressed:    "This is compressed text.",
		OriginalLen:   1000,
		CompressedLen: 250,
		Ratio:         0.25,
		Timestamp:     testTime,
		ModelName:     "llama-3.2-3b",
	}
	
	// Marshal to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("Failed to marshal CompressionEntry: %v", err)
	}
	
	// Unmarshal back
	var decoded CompressionEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal CompressionEntry: %v", err)
	}
	
	// Verify all fields
	if decoded.OriginalHash != entry.OriginalHash {
		t.Errorf("OriginalHash mismatch: got %s, want %s", decoded.OriginalHash, entry.OriginalHash)
	}
	if decoded.Compressed != entry.Compressed {
		t.Errorf("Compressed mismatch: got %s, want %s", decoded.Compressed, entry.Compressed)
	}
	if decoded.OriginalLen != entry.OriginalLen {
		t.Errorf("OriginalLen mismatch: got %d, want %d", decoded.OriginalLen, entry.OriginalLen)
	}
	if decoded.CompressedLen != entry.CompressedLen {
		t.Errorf("CompressedLen mismatch: got %d, want %d", decoded.CompressedLen, entry.CompressedLen)
	}
	if decoded.Ratio != entry.Ratio {
		t.Errorf("Ratio mismatch: got %f, want %f", decoded.Ratio, entry.Ratio)
	}
	if !decoded.Timestamp.Equal(entry.Timestamp) {
		t.Errorf("Timestamp mismatch: got %v, want %v", decoded.Timestamp, entry.Timestamp)
	}
	if decoded.ModelName != entry.ModelName {
		t.Errorf("ModelName mismatch: got %s, want %s", decoded.ModelName, entry.ModelName)
	}
}

// TestCompressionCache_JSONMarshaling tests JSON serialization of CompressionCache
func TestCompressionCache_JSONMarshaling(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)
	
	cache := &CompressionCache{
		Version:   "1.0",
		UpdatedAt: testTime,
		Entries: map[string]*CompressionEntry{
			"test-post": {
				OriginalHash:  "hash1",
				Compressed:    "Compressed version 1",
				OriginalLen:   1000,
				CompressedLen: 300,
				Ratio:         0.30,
				Timestamp:     testTime,
				ModelName:     "llama-3.2-3b",
			},
			"another-post": {
				OriginalHash:  "hash2",
				Compressed:    "Compressed version 2",
				OriginalLen:   2000,
				CompressedLen: 500,
				Ratio:         0.25,
				Timestamp:     testTime,
				ModelName:     "llama-3.2-3b",
			},
		},
	}
	
	// Marshal to JSON
	data, err := json.Marshal(cache)
	if err != nil {
		t.Fatalf("Failed to marshal CompressionCache: %v", err)
	}
	
	// Unmarshal back
	var decoded CompressionCache
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal CompressionCache: %v", err)
	}
	
	// Verify fields
	if decoded.Version != cache.Version {
		t.Errorf("Version mismatch: got %s, want %s", decoded.Version, cache.Version)
	}
	if !decoded.UpdatedAt.Equal(cache.UpdatedAt) {
		t.Errorf("UpdatedAt mismatch: got %v, want %v", decoded.UpdatedAt, cache.UpdatedAt)
	}
	if len(decoded.Entries) != len(cache.Entries) {
		t.Errorf("Entries count mismatch: got %d, want %d", len(decoded.Entries), len(cache.Entries))
	}
}

// TestProcessedPost_JSONMarshaling tests JSON serialization of ProcessedPost
func TestProcessedPost_JSONMarshaling(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)
	
	tests := []struct {
		name string
		post *ProcessedPost
	}{
		{
			name: "Success status",
			post: &ProcessedPost{
				Slug:      "test-post",
				Status:    StatusSuccess,
				Timestamp: testTime,
				Error:     "",
			},
		},
		{
			name: "Failed status",
			post: &ProcessedPost{
				Slug:      "failed-post",
				Status:    StatusFailed,
				Timestamp: testTime,
				Error:     "embedding generation failed",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.post)
			if err != nil {
				t.Fatalf("Failed to marshal ProcessedPost: %v", err)
			}
			
			// Unmarshal back
			var decoded ProcessedPost
			if err := json.Unmarshal(data, &decoded); err != nil {
				t.Fatalf("Failed to unmarshal ProcessedPost: %v", err)
			}
			
			// Verify all fields
			if decoded.Slug != tt.post.Slug {
				t.Errorf("Slug mismatch: got %s, want %s", decoded.Slug, tt.post.Slug)
			}
			if decoded.Status != tt.post.Status {
				t.Errorf("Status mismatch: got %s, want %s", decoded.Status, tt.post.Status)
			}
			if !decoded.Timestamp.Equal(tt.post.Timestamp) {
				t.Errorf("Timestamp mismatch: got %v, want %v", decoded.Timestamp, tt.post.Timestamp)
			}
			if decoded.Error != tt.post.Error {
				t.Errorf("Error mismatch: got %s, want %s", decoded.Error, tt.post.Error)
			}
		})
	}
}

// TestProgressFile_JSONMarshaling tests JSON serialization of ProgressFile
func TestProgressFile_JSONMarshaling(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)
	
	progress := &ProgressFile{
		Version:   "1.0",
		StartedAt: testTime,
		UpdatedAt: testTime.Add(1 * time.Hour),
		Processed: map[string]*ProcessedPost{
			"post-1": {
				Slug:      "post-1",
				Status:    StatusSuccess,
				Timestamp: testTime,
				Error:     "",
			},
			"post-2": {
				Slug:      "post-2",
				Status:    StatusFailed,
				Timestamp: testTime.Add(30 * time.Minute),
				Error:     "timeout error",
			},
		},
	}
	
	// Marshal to JSON
	data, err := json.Marshal(progress)
	if err != nil {
		t.Fatalf("Failed to marshal ProgressFile: %v", err)
	}
	
	// Unmarshal back
	var decoded ProgressFile
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ProgressFile: %v", err)
	}
	
	// Verify fields
	if decoded.Version != progress.Version {
		t.Errorf("Version mismatch: got %s, want %s", decoded.Version, progress.Version)
	}
	if !decoded.StartedAt.Equal(progress.StartedAt) {
		t.Errorf("StartedAt mismatch: got %v, want %v", decoded.StartedAt, progress.StartedAt)
	}
	if !decoded.UpdatedAt.Equal(progress.UpdatedAt) {
		t.Errorf("UpdatedAt mismatch: got %v, want %v", decoded.UpdatedAt, progress.UpdatedAt)
	}
	if len(decoded.Processed) != len(progress.Processed) {
		t.Errorf("Processed count mismatch: got %d, want %d", len(decoded.Processed), len(progress.Processed))
	}
	
	for slug, originalPost := range progress.Processed {
		decodedPost, exists := decoded.Processed[slug]
		if !exists {
			t.Errorf("Missing processed post: %s", slug)
			continue
		}
		if decodedPost.Status != originalPost.Status {
			t.Errorf("Post[%s].Status mismatch: got %s, want %s", slug, decodedPost.Status, originalPost.Status)
		}
	}
}

// TestStatus_Constants tests the Status type constants
func TestStatus_Constants(t *testing.T) {
	if StatusSuccess != "success" {
		t.Errorf("StatusSuccess constant mismatch: got %s, want 'success'", StatusSuccess)
	}
	if StatusFailed != "failed" {
		t.Errorf("StatusFailed constant mismatch: got %s, want 'failed'", StatusFailed)
	}
}

// TestCacheEntry_EmptyEmbedding tests handling of empty embeddings
func TestCacheEntry_EmptyEmbedding(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)
	
	entry := &CacheEntry{
		Embedding:   []float32{},
		ContentHash: "hash123",
		ModelName:   "model-v1",
		Timestamp:   testTime,
	}
	
	// Marshal and unmarshal
	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("Failed to marshal CacheEntry with empty embedding: %v", err)
	}
	
	var decoded CacheEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal CacheEntry with empty embedding: %v", err)
	}
	
	// Verify empty embedding is preserved
	if decoded.Embedding == nil {
		t.Error("Expected empty slice, got nil")
	}
	if len(decoded.Embedding) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(decoded.Embedding))
	}
}

// TestProgressFile_EmptyProcessed tests handling of empty processed map
func TestProgressFile_EmptyProcessed(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)
	
	progress := &ProgressFile{
		Version:   "1.0",
		StartedAt: testTime,
		UpdatedAt: testTime,
		Processed: map[string]*ProcessedPost{},
	}
	
	// Marshal and unmarshal
	data, err := json.Marshal(progress)
	if err != nil {
		t.Fatalf("Failed to marshal ProgressFile with empty processed: %v", err)
	}
	
	var decoded ProgressFile
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ProgressFile with empty processed: %v", err)
	}
	
	// Verify empty map is preserved
	if decoded.Processed == nil {
		t.Error("Expected empty map, got nil")
	}
	if len(decoded.Processed) != 0 {
		t.Errorf("Expected empty map, got length %d", len(decoded.Processed))
	}
}
