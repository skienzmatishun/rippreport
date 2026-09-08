package models

import (
	"testing"
	"time"
)

// TestEmbedding_Validate tests the Embedding validation method
func TestEmbedding_Validate(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name      string
		embedding *Embedding
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid embedding",
			embedding: &Embedding{
				Vector:      []float32{0.1, 0.2, 0.3},
				ContentHash: "hash123",
				ModelName:   "nomic-embed-text-v1.5",
				Timestamp:   testTime,
			},
			wantErr: false,
		},
		{
			name: "empty vector",
			embedding: &Embedding{
				Vector:      []float32{},
				ContentHash: "hash123",
				ModelName:   "nomic-embed-text-v1.5",
				Timestamp:   testTime,
			},
			wantErr: true,
			errMsg:  "embedding vector cannot be empty",
		},
		{
			name: "missing content hash",
			embedding: &Embedding{
				Vector:      []float32{0.1, 0.2, 0.3},
				ContentHash: "",
				ModelName:   "nomic-embed-text-v1.5",
				Timestamp:   testTime,
			},
			wantErr: true,
			errMsg:  "embedding content_hash cannot be empty",
		},
		{
			name: "missing model name",
			embedding: &Embedding{
				Vector:      []float32{0.1, 0.2, 0.3},
				ContentHash: "hash123",
				ModelName:   "",
				Timestamp:   testTime,
			},
			wantErr: true,
			errMsg:  "embedding model_name cannot be empty",
		},
		{
			name: "zero timestamp",
			embedding: &Embedding{
				Vector:      []float32{0.1, 0.2, 0.3},
				ContentHash: "hash123",
				ModelName:   "nomic-embed-text-v1.5",
				Timestamp:   time.Time{},
			},
			wantErr: true,
			errMsg:  "embedding timestamp cannot be zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.embedding.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Embedding.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Embedding.Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestCacheEntry_Validate tests the CacheEntry validation method
func TestCacheEntry_Validate(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name    string
		entry   *CacheEntry
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid cache entry",
			entry: &CacheEntry{
				Embedding:   []float32{0.1, 0.2, 0.3},
				ContentHash: "hash123",
				ModelName:   "nomic-embed-text-v1.5",
				Timestamp:   testTime,
			},
			wantErr: false,
		},
		{
			name: "empty embedding",
			entry: &CacheEntry{
				Embedding:   []float32{},
				ContentHash: "hash123",
				ModelName:   "nomic-embed-text-v1.5",
				Timestamp:   testTime,
			},
			wantErr: true,
			errMsg:  "cache entry embedding cannot be empty",
		},
		{
			name: "missing content hash",
			entry: &CacheEntry{
				Embedding:   []float32{0.1, 0.2, 0.3},
				ContentHash: "",
				ModelName:   "nomic-embed-text-v1.5",
				Timestamp:   testTime,
			},
			wantErr: true,
			errMsg:  "cache entry content_hash cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CacheEntry.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("CacheEntry.Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestCacheFile_Validate tests the CacheFile validation method
func TestCacheFile_Validate(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)

	validEntry := &CacheEntry{
		Embedding:   []float32{0.1, 0.2, 0.3},
		ContentHash: "hash123",
		ModelName:   "nomic-embed-text-v1.5",
		Timestamp:   testTime,
	}

	invalidEntry := &CacheEntry{
		Embedding:   []float32{},
		ContentHash: "hash123",
		ModelName:   "nomic-embed-text-v1.5",
		Timestamp:   testTime,
	}

	tests := []struct {
		name      string
		cacheFile *CacheFile
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid cache file",
			cacheFile: &CacheFile{
				Version:   "1.0",
				ModelName: "nomic-embed-text-v1.5",
				UpdatedAt: testTime,
				Entries: map[string]*CacheEntry{
					"post1": validEntry,
				},
			},
			wantErr: false,
		},
		{
			name: "missing version",
			cacheFile: &CacheFile{
				Version:   "",
				ModelName: "nomic-embed-text-v1.5",
				UpdatedAt: testTime,
				Entries:   map[string]*CacheEntry{},
			},
			wantErr: true,
			errMsg:  "cache file version cannot be empty",
		},
		{
			name: "missing model name",
			cacheFile: &CacheFile{
				Version:   "1.0",
				ModelName: "",
				UpdatedAt: testTime,
				Entries:   map[string]*CacheEntry{},
			},
			wantErr: true,
			errMsg:  "cache file model_name cannot be empty",
		},
		{
			name: "nil entries",
			cacheFile: &CacheFile{
				Version:   "1.0",
				ModelName: "nomic-embed-text-v1.5",
				UpdatedAt: testTime,
				Entries:   nil,
			},
			wantErr: true,
			errMsg:  "cache file entries cannot be nil",
		},
		{
			name: "zero updated_at",
			cacheFile: &CacheFile{
				Version:   "1.0",
				ModelName: "nomic-embed-text-v1.5",
				UpdatedAt: time.Time{},
				Entries:   map[string]*CacheEntry{},
			},
			wantErr: true,
			errMsg:  "cache file updated_at cannot be zero",
		},
		{
			name: "invalid entry",
			cacheFile: &CacheFile{
				Version:   "1.0",
				ModelName: "nomic-embed-text-v1.5",
				UpdatedAt: testTime,
				Entries: map[string]*CacheEntry{
					"post1": invalidEntry,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cacheFile.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CacheFile.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("CacheFile.Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestCompressionEntry_Validate tests the CompressionEntry validation method
func TestCompressionEntry_Validate(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name    string
		entry   *CompressionEntry
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid compression entry",
			entry: &CompressionEntry{
				OriginalHash:  "hash123",
				Compressed:    "Compressed text",
				OriginalLen:   1000,
				CompressedLen: 300,
				Ratio:         0.30,
				Timestamp:     testTime,
				ModelName:     "llama-3.2-3b",
			},
			wantErr: false,
		},
		{
			name: "missing original hash",
			entry: &CompressionEntry{
				OriginalHash:  "",
				Compressed:    "Compressed text",
				OriginalLen:   1000,
				CompressedLen: 300,
				Ratio:         0.30,
				Timestamp:     testTime,
				ModelName:     "llama-3.2-3b",
			},
			wantErr: true,
			errMsg:  "compression entry original_hash cannot be empty",
		},
		{
			name: "missing compressed text",
			entry: &CompressionEntry{
				OriginalHash:  "hash123",
				Compressed:    "",
				OriginalLen:   1000,
				CompressedLen: 300,
				Ratio:         0.30,
				Timestamp:     testTime,
				ModelName:     "llama-3.2-3b",
			},
			wantErr: true,
			errMsg:  "compression entry compressed cannot be empty",
		},
		{
			name: "zero original length",
			entry: &CompressionEntry{
				OriginalHash:  "hash123",
				Compressed:    "Compressed text",
				OriginalLen:   0,
				CompressedLen: 300,
				Ratio:         0.30,
				Timestamp:     testTime,
				ModelName:     "llama-3.2-3b",
			},
			wantErr: true,
			errMsg:  "compression entry original_len must be positive",
		},
		{
			name: "negative compressed length",
			entry: &CompressionEntry{
				OriginalHash:  "hash123",
				Compressed:    "Compressed text",
				OriginalLen:   1000,
				CompressedLen: -1,
				Ratio:         0.30,
				Timestamp:     testTime,
				ModelName:     "llama-3.2-3b",
			},
			wantErr: true,
			errMsg:  "compression entry compressed_len must be positive",
		},
		{
			name: "ratio too high",
			entry: &CompressionEntry{
				OriginalHash:  "hash123",
				Compressed:    "Compressed text",
				OriginalLen:   1000,
				CompressedLen: 300,
				Ratio:         1.5,
				Timestamp:     testTime,
				ModelName:     "llama-3.2-3b",
			},
			wantErr: true,
			errMsg:  "compression entry ratio must be between 0 and 1, got 1.500000",
		},
		{
			name: "zero ratio",
			entry: &CompressionEntry{
				OriginalHash:  "hash123",
				Compressed:    "Compressed text",
				OriginalLen:   1000,
				CompressedLen: 300,
				Ratio:         0,
				Timestamp:     testTime,
				ModelName:     "llama-3.2-3b",
			},
			wantErr: true,
			errMsg:  "compression entry ratio must be between 0 and 1, got 0.000000",
		},
		{
			name: "missing model name",
			entry: &CompressionEntry{
				OriginalHash:  "hash123",
				Compressed:    "Compressed text",
				OriginalLen:   1000,
				CompressedLen: 300,
				Ratio:         0.30,
				Timestamp:     testTime,
				ModelName:     "",
			},
			wantErr: true,
			errMsg:  "compression entry model_name cannot be empty",
		},
		{
			name: "zero timestamp",
			entry: &CompressionEntry{
				OriginalHash:  "hash123",
				Compressed:    "Compressed text",
				OriginalLen:   1000,
				CompressedLen: 300,
				Ratio:         0.30,
				Timestamp:     time.Time{},
				ModelName:     "llama-3.2-3b",
			},
			wantErr: true,
			errMsg:  "compression entry timestamp cannot be zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CompressionEntry.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("CompressionEntry.Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestCompressionCache_Validate tests the CompressionCache validation method
func TestCompressionCache_Validate(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)

	validEntry := &CompressionEntry{
		OriginalHash:  "hash123",
		Compressed:    "Compressed text",
		OriginalLen:   1000,
		CompressedLen: 300,
		Ratio:         0.30,
		Timestamp:     testTime,
		ModelName:     "llama-3.2-3b",
	}

	invalidEntry := &CompressionEntry{
		OriginalHash:  "",
		Compressed:    "Compressed text",
		OriginalLen:   1000,
		CompressedLen: 300,
		Ratio:         0.30,
		Timestamp:     testTime,
		ModelName:     "llama-3.2-3b",
	}

	tests := []struct {
		name    string
		cache   *CompressionCache
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid compression cache",
			cache: &CompressionCache{
				Version:   "1.0",
				UpdatedAt: testTime,
				Entries: map[string]*CompressionEntry{
					"post1": validEntry,
				},
			},
			wantErr: false,
		},
		{
			name: "missing version",
			cache: &CompressionCache{
				Version:   "",
				UpdatedAt: testTime,
				Entries:   map[string]*CompressionEntry{},
			},
			wantErr: true,
			errMsg:  "compression cache version cannot be empty",
		},
		{
			name: "nil entries",
			cache: &CompressionCache{
				Version:   "1.0",
				UpdatedAt: testTime,
				Entries:   nil,
			},
			wantErr: true,
			errMsg:  "compression cache entries cannot be nil",
		},
		{
			name: "zero updated_at",
			cache: &CompressionCache{
				Version:   "1.0",
				UpdatedAt: time.Time{},
				Entries:   map[string]*CompressionEntry{},
			},
			wantErr: true,
			errMsg:  "compression cache updated_at cannot be zero",
		},
		{
			name: "invalid entry",
			cache: &CompressionCache{
				Version:   "1.0",
				UpdatedAt: testTime,
				Entries: map[string]*CompressionEntry{
					"post1": invalidEntry,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cache.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CompressionCache.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("CompressionCache.Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestProcessedPost_Validate tests the ProcessedPost validation method
func TestProcessedPost_Validate(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name    string
		post    *ProcessedPost
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid success post",
			post: &ProcessedPost{
				Slug:      "test-post",
				Status:    StatusSuccess,
				Timestamp: testTime,
				Error:     "",
			},
			wantErr: false,
		},
		{
			name: "valid failed post",
			post: &ProcessedPost{
				Slug:      "test-post",
				Status:    StatusFailed,
				Timestamp: testTime,
				Error:     "processing error",
			},
			wantErr: false,
		},
		{
			name: "missing slug",
			post: &ProcessedPost{
				Slug:      "",
				Status:    StatusSuccess,
				Timestamp: testTime,
				Error:     "",
			},
			wantErr: true,
			errMsg:  "processed post slug cannot be empty",
		},
		{
			name: "invalid status",
			post: &ProcessedPost{
				Slug:      "test-post",
				Status:    "invalid",
				Timestamp: testTime,
				Error:     "",
			},
			wantErr: true,
			errMsg:  "processed post status must be 'success' or 'failed', got invalid",
		},
		{
			name: "zero timestamp",
			post: &ProcessedPost{
				Slug:      "test-post",
				Status:    StatusSuccess,
				Timestamp: time.Time{},
				Error:     "",
			},
			wantErr: true,
			errMsg:  "processed post timestamp cannot be zero",
		},
		{
			name: "failed status without error",
			post: &ProcessedPost{
				Slug:      "test-post",
				Status:    StatusFailed,
				Timestamp: testTime,
				Error:     "",
			},
			wantErr: true,
			errMsg:  "processed post with failed status must have an error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.post.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessedPost.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("ProcessedPost.Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestProgressFile_Validate tests the ProgressFile validation method
func TestProgressFile_Validate(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)

	validPost := &ProcessedPost{
		Slug:      "test-post",
		Status:    StatusSuccess,
		Timestamp: testTime,
		Error:     "",
	}

	invalidPost := &ProcessedPost{
		Slug:      "",
		Status:    StatusSuccess,
		Timestamp: testTime,
		Error:     "",
	}

	mismatchPost := &ProcessedPost{
		Slug:      "wrong-slug",
		Status:    StatusSuccess,
		Timestamp: testTime,
		Error:     "",
	}

	tests := []struct {
		name     string
		progress *ProgressFile
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid progress file",
			progress: &ProgressFile{
				Version:   "1.0",
				StartedAt: testTime,
				UpdatedAt: testTime,
				Processed: map[string]*ProcessedPost{
					"test-post": validPost,
				},
			},
			wantErr: false,
		},
		{
			name: "missing version",
			progress: &ProgressFile{
				Version:   "",
				StartedAt: testTime,
				UpdatedAt: testTime,
				Processed: map[string]*ProcessedPost{},
			},
			wantErr: true,
			errMsg:  "progress file version cannot be empty",
		},
		{
			name: "zero started_at",
			progress: &ProgressFile{
				Version:   "1.0",
				StartedAt: time.Time{},
				UpdatedAt: testTime,
				Processed: map[string]*ProcessedPost{},
			},
			wantErr: true,
			errMsg:  "progress file started_at cannot be zero",
		},
		{
			name: "zero updated_at",
			progress: &ProgressFile{
				Version:   "1.0",
				StartedAt: testTime,
				UpdatedAt: time.Time{},
				Processed: map[string]*ProcessedPost{},
			},
			wantErr: true,
			errMsg:  "progress file updated_at cannot be zero",
		},
		{
			name: "nil processed",
			progress: &ProgressFile{
				Version:   "1.0",
				StartedAt: testTime,
				UpdatedAt: testTime,
				Processed: nil,
			},
			wantErr: true,
			errMsg:  "progress file processed cannot be nil",
		},
		{
			name: "invalid post",
			progress: &ProgressFile{
				Version:   "1.0",
				StartedAt: testTime,
				UpdatedAt: testTime,
				Processed: map[string]*ProcessedPost{
					"test-post": invalidPost,
				},
			},
			wantErr: true,
		},
		{
			name: "mismatched slug",
			progress: &ProgressFile{
				Version:   "1.0",
				StartedAt: testTime,
				UpdatedAt: testTime,
				Processed: map[string]*ProcessedPost{
					"test-post": mismatchPost,
				},
			},
			wantErr: true,
			errMsg:  "processed post key test-post does not match slug wrong-slug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.progress.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ProgressFile.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("ProgressFile.Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestStagedRankings_Validate tests the StagedRankings validation method
func TestStagedRankings_Validate(t *testing.T) {
	testTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)

	validArticle := RelatedArticle{
		Slug:  "related-1",
		Title: "Related Article 1",
		Score: 95,
		Rank:  1,
	}

	invalidArticle := RelatedArticle{
		Slug:  "",
		Title: "Related Article",
		Score: 95,
		Rank:  1,
	}

	tests := []struct {
		name    string
		staged  *StagedRankings
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid staged rankings",
			staged: &StagedRankings{
				PostSlug:  "test-post",
				Generated: testTime,
				Articles:  []RelatedArticle{validArticle},
			},
			wantErr: false,
		},
		{
			name: "missing post slug",
			staged: &StagedRankings{
				PostSlug:  "",
				Generated: testTime,
				Articles:  []RelatedArticle{validArticle},
			},
			wantErr: true,
			errMsg:  "staged rankings post_slug cannot be empty",
		},
		{
			name: "zero generated",
			staged: &StagedRankings{
				PostSlug:  "test-post",
				Generated: time.Time{},
				Articles:  []RelatedArticle{validArticle},
			},
			wantErr: true,
			errMsg:  "staged rankings generated cannot be zero",
		},
		{
			name: "no articles",
			staged: &StagedRankings{
				PostSlug:  "test-post",
				Generated: testTime,
				Articles:  []RelatedArticle{},
			},
			wantErr: true,
			errMsg:  "staged rankings must contain at least one article",
		},
		{
			name: "invalid article",
			staged: &StagedRankings{
				PostSlug:  "test-post",
				Generated: testTime,
				Articles:  []RelatedArticle{invalidArticle},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.staged.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("StagedRankings.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("StagedRankings.Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestCacheStats_CalculateHitRate tests the CacheStats hit rate calculation
func TestCacheStats_CalculateHitRate(t *testing.T) {
	tests := []struct {
		name    string
		stats   *CacheStats
		wantHit float64
	}{
		{
			name: "50% hit rate",
			stats: &CacheStats{
				Hits:   50,
				Misses: 50,
			},
			wantHit: 0.5,
		},
		{
			name: "100% hit rate",
			stats: &CacheStats{
				Hits:   100,
				Misses: 0,
			},
			wantHit: 1.0,
		},
		{
			name: "0% hit rate",
			stats: &CacheStats{
				Hits:   0,
				Misses: 100,
			},
			wantHit: 0.0,
		},
		{
			name: "no requests",
			stats: &CacheStats{
				Hits:   0,
				Misses: 0,
			},
			wantHit: 0.0,
		},
		{
			name: "75% hit rate",
			stats: &CacheStats{
				Hits:   75,
				Misses: 25,
			},
			wantHit: 0.75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stats.CalculateHitRate()
			if tt.stats.HitRate != tt.wantHit {
				t.Errorf("CacheStats.CalculateHitRate() = %v, want %v", tt.stats.HitRate, tt.wantHit)
			}
		})
	}
}
