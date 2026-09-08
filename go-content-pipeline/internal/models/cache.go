package models

import (
	"fmt"
	"time"
)

// Embedding represents a cached embedding vector
type Embedding struct {
	Vector      []float32 `json:"embedding"`
	ContentHash string    `json:"content_hash"`
	ModelName   string    `json:"model_name"`
	Timestamp   time.Time `json:"timestamp"`
}

// CacheEntry represents a single cache entry for embeddings
type CacheEntry struct {
	Embedding   []float32 `json:"embedding"`
	ContentHash string    `json:"content_hash"`
	ModelName   string    `json:"model_name"`
	Timestamp   time.Time `json:"timestamp"`
}

// CacheFile represents the entire embedding cache file structure
type CacheFile struct {
	Version   string                 `json:"version"`
	ModelName string                 `json:"model_name"`
	Entries   map[string]*CacheEntry `json:"entries"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// CompressionEntry represents a cached compressed article
type CompressionEntry struct {
	OriginalHash  string    `json:"original_hash"`
	Compressed    string    `json:"compressed"`
	OriginalLen   int       `json:"original_len"`
	CompressedLen int       `json:"compressed_len"`
	Ratio         float32   `json:"ratio"`
	Timestamp     time.Time `json:"timestamp"`
	ModelName     string    `json:"model_name"`
}

// CompressionCache represents the compression cache file structure
type CompressionCache struct {
	Version   string                       `json:"version"`
	Entries   map[string]*CompressionEntry `json:"entries"`
	UpdatedAt time.Time                    `json:"updated_at"`
}

// ProcessedPost represents a post that has been processed
type ProcessedPost struct {
	Slug      string    `json:"slug"`
	Status    Status    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"`
}

// Status represents the processing status of a post
type Status string

const (
	// StatusSuccess indicates successful processing
	StatusSuccess Status = "success"
	// StatusFailed indicates failed processing
	StatusFailed Status = "failed"
)

// ProgressFile represents the progress tracking file structure
type ProgressFile struct {
	Version   string                    `json:"version"`
	StartedAt time.Time                 `json:"started_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
	Processed map[string]*ProcessedPost `json:"processed"`
}

// FailedPost represents a post that failed processing
type FailedPost struct {
	Slug  string `json:"slug"`
	Error string `json:"error"`
}

// StagedRankings represents staged rankings before application
type StagedRankings struct {
	PostSlug  string           `yaml:"post_slug"`
	Generated time.Time        `yaml:"generated"`
	Articles  []RelatedArticle `yaml:"articles"`
}

// CacheStats provides statistics about cache performance
type CacheStats struct {
	TotalEntries int     `json:"total_entries"`
	Hits         int64   `json:"hits"`
	Misses       int64   `json:"misses"`
	HitRate      float64 `json:"hit_rate"`
}

// Validate checks if the Embedding is valid
func (e *Embedding) Validate() error {
	if len(e.Vector) == 0 {
		return fmt.Errorf("embedding vector cannot be empty")
	}
	if e.ContentHash == "" {
		return fmt.Errorf("embedding content_hash cannot be empty")
	}
	if e.ModelName == "" {
		return fmt.Errorf("embedding model_name cannot be empty")
	}
	if e.Timestamp.IsZero() {
		return fmt.Errorf("embedding timestamp cannot be zero")
	}
	return nil
}

// Validate checks if the CacheEntry is valid
func (ce *CacheEntry) Validate() error {
	if len(ce.Embedding) == 0 {
		return fmt.Errorf("cache entry embedding cannot be empty")
	}
	if ce.ContentHash == "" {
		return fmt.Errorf("cache entry content_hash cannot be empty")
	}
	if ce.ModelName == "" {
		return fmt.Errorf("cache entry model_name cannot be empty")
	}
	if ce.Timestamp.IsZero() {
		return fmt.Errorf("cache entry timestamp cannot be zero")
	}
	return nil
}

// Validate checks if the CacheFile is valid
func (cf *CacheFile) Validate() error {
	if cf.Version == "" {
		return fmt.Errorf("cache file version cannot be empty")
	}
	if cf.ModelName == "" {
		return fmt.Errorf("cache file model_name cannot be empty")
	}
	if cf.Entries == nil {
		return fmt.Errorf("cache file entries cannot be nil")
	}
	if cf.UpdatedAt.IsZero() {
		return fmt.Errorf("cache file updated_at cannot be zero")
	}

	// Validate each entry
	for path, entry := range cf.Entries {
		if err := entry.Validate(); err != nil {
			return fmt.Errorf("cache entry %s: %w", path, err)
		}
	}

	return nil
}

// Validate checks if the CompressionEntry is valid
func (ce *CompressionEntry) Validate() error {
	if ce.OriginalHash == "" {
		return fmt.Errorf("compression entry original_hash cannot be empty")
	}
	if ce.Compressed == "" {
		return fmt.Errorf("compression entry compressed cannot be empty")
	}
	if ce.OriginalLen <= 0 {
		return fmt.Errorf("compression entry original_len must be positive")
	}
	if ce.CompressedLen <= 0 {
		return fmt.Errorf("compression entry compressed_len must be positive")
	}
	if ce.Ratio <= 0 || ce.Ratio > 1 {
		return fmt.Errorf("compression entry ratio must be between 0 and 1, got %f", ce.Ratio)
	}
	if ce.ModelName == "" {
		return fmt.Errorf("compression entry model_name cannot be empty")
	}
	if ce.Timestamp.IsZero() {
		return fmt.Errorf("compression entry timestamp cannot be zero")
	}
	return nil
}

// Validate checks if the CompressionCache is valid
func (cc *CompressionCache) Validate() error {
	if cc.Version == "" {
		return fmt.Errorf("compression cache version cannot be empty")
	}
	if cc.Entries == nil {
		return fmt.Errorf("compression cache entries cannot be nil")
	}
	if cc.UpdatedAt.IsZero() {
		return fmt.Errorf("compression cache updated_at cannot be zero")
	}

	// Validate each entry
	for slug, entry := range cc.Entries {
		if err := entry.Validate(); err != nil {
			return fmt.Errorf("compression entry %s: %w", slug, err)
		}
	}

	return nil
}

// Validate checks if the ProcessedPost is valid
func (pp *ProcessedPost) Validate() error {
	if pp.Slug == "" {
		return fmt.Errorf("processed post slug cannot be empty")
	}
	if pp.Status != StatusSuccess && pp.Status != StatusFailed {
		return fmt.Errorf("processed post status must be 'success' or 'failed', got %s", pp.Status)
	}
	if pp.Timestamp.IsZero() {
		return fmt.Errorf("processed post timestamp cannot be zero")
	}
	if pp.Status == StatusFailed && pp.Error == "" {
		return fmt.Errorf("processed post with failed status must have an error message")
	}
	return nil
}

// Validate checks if the ProgressFile is valid
func (pf *ProgressFile) Validate() error {
	if pf.Version == "" {
		return fmt.Errorf("progress file version cannot be empty")
	}
	if pf.StartedAt.IsZero() {
		return fmt.Errorf("progress file started_at cannot be zero")
	}
	if pf.UpdatedAt.IsZero() {
		return fmt.Errorf("progress file updated_at cannot be zero")
	}
	if pf.Processed == nil {
		return fmt.Errorf("progress file processed cannot be nil")
	}

	// Validate each processed post
	for slug, post := range pf.Processed {
		if slug != post.Slug {
			return fmt.Errorf("processed post key %s does not match slug %s", slug, post.Slug)
		}
		if err := post.Validate(); err != nil {
			return fmt.Errorf("processed post %s: %w", slug, err)
		}
	}

	return nil
}

// Validate checks if the StagedRankings is valid
func (sr *StagedRankings) Validate() error {
	if sr.PostSlug == "" {
		return fmt.Errorf("staged rankings post_slug cannot be empty")
	}
	if sr.Generated.IsZero() {
		return fmt.Errorf("staged rankings generated cannot be zero")
	}
	if len(sr.Articles) == 0 {
		return fmt.Errorf("staged rankings must contain at least one article")
	}

	// Validate each article
	for i, article := range sr.Articles {
		if err := article.Validate(); err != nil {
			return fmt.Errorf("article %d: %w", i, err)
		}
	}

	return nil
}

// CalculateHitRate calculates the cache hit rate
func (cs *CacheStats) CalculateHitRate() {
	total := cs.Hits + cs.Misses
	if total == 0 {
		cs.HitRate = 0.0
	} else {
		cs.HitRate = float64(cs.Hits) / float64(total)
	}
}
