package processor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

const (
	defaultProgressVersion = "1.0"
	defaultSaveInterval    = 10
)

// ProgressTracker tracks processing status of posts, supporting persistence, resume, and failure reporting.
//
// Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7, 7.8
type ProgressTracker struct {
	mu                      sync.Mutex
	filePath                string
	version                 string
	startedAt               time.Time
	processed               map[string]*models.ProcessedPost
	processedSinceLastSave  int
	saveInterval            int
}

// NewProgressTracker creates a new ProgressTracker instance.
func NewProgressTracker(filePath string) *ProgressTracker {
	now := time.Now()
	return &ProgressTracker{
		filePath:     filePath,
		version:      defaultProgressVersion,
		startedAt:    now,
		processed:    make(map[string]*models.ProcessedPost),
		saveInterval: defaultSaveInterval,
	}
}

// MarkProcessed marks a post as successfully processed.
// Automatically persists progress after every 10 processed posts.
//
// Requirements: 7.1, 7.5
func (pt *ProgressTracker) MarkProcessed(postSlug string) error {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	pt.processed[postSlug] = &models.ProcessedPost{
		Slug:      postSlug,
		Status:    models.StatusSuccess,
		Timestamp: time.Now(),
	}

	pt.processedSinceLastSave++
	if pt.processedSinceLastSave >= pt.saveInterval && pt.filePath != "" {
		pt.processedSinceLastSave = 0
		return pt.persistLocked()
	}

	return nil
}

// MarkFailed marks a post as failed with an associated error message.
// Automatically persists progress after every 10 processed posts.
//
// Requirements: 7.2, 7.5
func (pt *ProgressTracker) MarkFailed(postSlug string, err error) error {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	pt.processed[postSlug] = &models.ProcessedPost{
		Slug:      postSlug,
		Status:    models.StatusFailed,
		Timestamp: time.Now(),
		Error:     errMsg,
	}

	pt.processedSinceLastSave++
	if pt.processedSinceLastSave >= pt.saveInterval && pt.filePath != "" {
		pt.processedSinceLastSave = 0
		return pt.persistLocked()
	}

	return nil
}

// IsProcessed checks whether a post has been processed (either success or failure).
//
// Requirements: 7.3
func (pt *ProgressTracker) IsProcessed(postSlug string) bool {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	entry, exists := pt.processed[postSlug]
	return exists && entry.Status == models.StatusSuccess
}

// GetFailed returns a list of posts that failed processing.
//
// Requirements: 7.4
func (pt *ProgressTracker) GetFailed() []models.FailedPost {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	var failed []models.FailedPost
	for slug, post := range pt.processed {
		if post.Status == models.StatusFailed {
			failed = append(failed, models.FailedPost{
				Slug:  slug,
				Error: post.Error,
			})
		}
	}
	return failed
}

// FilterUnprocessedSlugs returns only slugs that have not yet been successfully processed.
//
// Requirements: 7.7
func (pt *ProgressTracker) FilterUnprocessedSlugs(slugs []string) []string {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	var unprocessed []string
	for _, slug := range slugs {
		entry, exists := pt.processed[slug]
		if !exists || entry.Status != models.StatusSuccess {
			unprocessed = append(unprocessed, slug)
		}
	}
	return unprocessed
}

// Persist saves the progress state to disk atomically using temporary file + rename.
//
// Requirements: 7.5, 7.6
func (pt *ProgressTracker) Persist() error {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	return pt.persistLocked()
}

func (pt *ProgressTracker) persistLocked() error {
	if pt.filePath == "" {
		return nil
	}

	// Create snapshot
	snapshot := make(map[string]*models.ProcessedPost, len(pt.processed))
	for k, v := range pt.processed {
		snapshot[k] = v
	}

	progressFile := models.ProgressFile{
		Version:   pt.version,
		StartedAt: pt.startedAt,
		UpdatedAt: time.Now(),
		Processed: snapshot,
	}

	data, err := json.MarshalIndent(progressFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal progress: %w", err)
	}

	dir := filepath.Dir(pt.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create progress directory: %w", err)
	}

	tempFile, err := os.CreateTemp(dir, ".progress_*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp progress file: %w", err)
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
		return fmt.Errorf("failed to write progress to temp file: %w", err)
	}

	if err := tempFile.Sync(); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to sync temp progress file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp progress file: %w", err)
	}

	if err := os.Rename(tempPath, pt.filePath); err != nil {
		return fmt.Errorf("failed to atomically replace progress file: %w", err)
	}

	success = true
	return nil
}

// Load loads the progress state from disk. Starts fresh if corrupted.
//
// Requirements: 7.6, 7.8
func (pt *ProgressTracker) Load() error {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	if pt.filePath == "" {
		return nil
	}

	data, err := os.ReadFile(pt.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			pt.processed = make(map[string]*models.ProcessedPost)
			return nil
		}
		return fmt.Errorf("failed to read progress file: %w", err)
	}

	var pf models.ProgressFile
	if err := json.Unmarshal(data, &pf); err != nil {
		// Corrupted progress file: start fresh
		pt.processed = make(map[string]*models.ProcessedPost)
		return nil
	}

	if pf.Processed == nil {
		pt.processed = make(map[string]*models.ProcessedPost)
	} else {
		pt.processed = pf.Processed
	}

	if !pf.StartedAt.IsZero() {
		pt.startedAt = pf.StartedAt
	}
	if pf.Version != "" {
		pt.version = pf.Version
	}

	return nil
}

// ProcessedCount returns the total number of successfully processed posts.
func (pt *ProgressTracker) ProcessedCount() int {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	count := 0
	for _, post := range pt.processed {
		if post.Status == models.StatusSuccess {
			count++
		}
	}
	return count
}
