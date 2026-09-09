package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/client"
	"github.com/rippreport/go-content-pipeline/internal/models"
)

const (
	defaultCompressionVersion = "1.0"
	minCharsForCompression    = 500
)

// EstimateTokens returns an approximate token count based on 1 token ≈ 4 characters.
//
// Requirements: 13.6
func EstimateTokens(text string) int {
	if len(text) == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(text)) / 4.0))
}

// Truncate preserves the first 60% and last 40% of text within maxTokens with a truncation marker.
//
// Requirements: Design algorithm section
func Truncate(text string, maxTokens int) string {
	if maxTokens <= 0 || EstimateTokens(text) <= maxTokens {
		return text
	}

	maxChars := maxTokens * 4
	headChars := int(float64(maxChars) * 0.6)
	tailChars := maxChars - headChars

	if headChars+tailChars >= len(text) {
		return text
	}

	head := text[:headChars]
	tail := text[len(text)-tailChars:]

	return head + "\n[... content truncated ...]\n" + tail
}

// CompressionCache manages caching of LLM-compressed post summaries.
//
// Requirements: 13.2, 13.3, 13.8
type CompressionCache struct {
	mu        sync.RWMutex
	filePath  string
	modelName string
	version   string
	entries   map[string]*models.CompressionEntry
}

// NewCompressionCache creates a new CompressionCache.
func NewCompressionCache(filePath, modelName string) *CompressionCache {
	return &CompressionCache{
		filePath:  filePath,
		modelName: modelName,
		version:   defaultCompressionVersion,
		entries:   make(map[string]*models.CompressionEntry),
	}
}

// GetCached returns the cached compressed content if present and valid.
func (cc *CompressionCache) GetCached(postSlug string) (string, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	entry, exists := cc.entries[postSlug]
	if !exists {
		return "", false
	}
	return entry.Compressed, true
}

// PutCached stores a compressed version of a post in the cache.
func (cc *CompressionCache) PutCached(postSlug, contentHash, compressed string, origLen int) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	compLen := len(compressed)
	ratio := float32(0.5)
	if origLen > 0 {
		ratio = float32(compLen) / float32(origLen)
		if ratio > 1.0 {
			ratio = 1.0
		}
	}

	cc.entries[postSlug] = &models.CompressionEntry{
		OriginalHash:  contentHash,
		Compressed:    compressed,
		OriginalLen:   origLen,
		CompressedLen: compLen,
		Ratio:         ratio,
		Timestamp:     time.Now(),
		ModelName:     cc.modelName,
	}
}

// Persist saves the compression cache to disk atomically.
func (cc *CompressionCache) Persist() error {
	cc.mu.RLock()
	if cc.filePath == "" {
		cc.mu.RUnlock()
		return nil
	}

	snapshot := make(map[string]*models.CompressionEntry, len(cc.entries))
	for k, v := range cc.entries {
		snapshot[k] = v
	}
	version := cc.version
	filePath := cc.filePath
	cc.mu.RUnlock()

	cacheFile := models.CompressionCache{
		Version:   version,
		Entries:   snapshot,
		UpdatedAt: time.Now(),
	}

	data, err := json.MarshalIndent(cacheFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal compression cache: %w", err)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for compression cache: %w", err)
	}

	tempFile, err := os.CreateTemp(dir, ".compression_cache_*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file for compression cache: %w", err)
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
		return fmt.Errorf("failed to write compression cache: %w", err)
	}
	if err := tempFile.Sync(); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to sync compression cache: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp compression cache: %w", err)
	}
	if err := os.Rename(tempPath, filePath); err != nil {
		return fmt.Errorf("failed to atomically replace compression cache: %w", err)
	}

	success = true
	return nil
}

// Len returns the number of cached entries in the compression cache.
func (cc *CompressionCache) Len() int {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	return len(cc.entries)
}

// Load loads the compression cache from disk, starting fresh on corruption.
// It supports polymorphic JSON formats:
// 1. Structured models.CompressionCache ({"version": "...", "entries": {...}})
// 2. Flat map[string]string ({"slug": "summary..."}) - as in compressed_articles.json
// 3. Flat map[string]*models.CompressionEntry
// If cc.filePath does not exist, it checks fallback locations (e.g. internal/cache/compressed_articles.json).
func (cc *CompressionCache) Load() error {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	targetPath := cc.filePath
	if targetPath == "" {
		targetPath = filepath.Join("internal", "cache", "compressed_articles.json")
	}

	data, err := os.ReadFile(targetPath)
	if err != nil && os.IsNotExist(err) {
		// Try fallback locations
		candidates := []string{
			filepath.Join("internal", "cache", "compressed_articles.json"),
			"compressed_articles.json",
			".compression_cache.json",
		}
		if cc.filePath != "" {
			dir := filepath.Dir(cc.filePath)
			candidates = append([]string{
				filepath.Join(dir, "compressed_articles.json"),
				filepath.Join(dir, "internal", "cache", "compressed_articles.json"),
			}, candidates...)
		}
		for _, cand := range candidates {
			if cand == targetPath {
				continue
			}
			if d, rErr := os.ReadFile(cand); rErr == nil {
				data = d
				err = nil
				targetPath = cand
				break
			}
		}
	}

	if err != nil {
		if os.IsNotExist(err) {
			cc.entries = make(map[string]*models.CompressionEntry)
			return nil
		}
		return fmt.Errorf("failed to read compression cache file: %w", err)
	}

	// 1. Try structured models.CompressionCache
	var cacheFile models.CompressionCache
	if uErr := json.Unmarshal(data, &cacheFile); uErr == nil && len(cacheFile.Entries) > 0 {
		cc.entries = cacheFile.Entries
		return nil
	}

	// 2. Try raw map[string]string (format of compressed_articles.json)
	var rawStrings map[string]string
	if uErr := json.Unmarshal(data, &rawStrings); uErr == nil && len(rawStrings) > 0 {
		cc.entries = make(map[string]*models.CompressionEntry, len(rawStrings))
		now := time.Now()
		for slug, text := range rawStrings {
			textLen := len(text)
			cc.entries[slug] = &models.CompressionEntry{
				Compressed:    text,
				OriginalLen:   textLen,
				CompressedLen: textLen,
				Ratio:         0.5,
				Timestamp:     now,
				ModelName:     cc.modelName,
			}
		}
		return nil
	}

	// 3. Try raw map[string]*models.CompressionEntry
	var rawEntries map[string]*models.CompressionEntry
	if uErr := json.Unmarshal(data, &rawEntries); uErr == nil && len(rawEntries) > 0 {
		cc.entries = rawEntries
		return nil
	}

	cc.entries = make(map[string]*models.CompressionEntry)
	return nil
}

// Compress returns a compressed summary of an article, checking cache first.
// Skips compression if article is under 500 characters.
//
// Requirements: 13.1, 13.3, 13.4, 13.5, 13.7
func Compress(ctx context.Context, llamaClient client.LlamaClient, cache *CompressionCache, article *models.Post, targetTokens int) (string, error) {
	if article == nil {
		return "", fmt.Errorf("article cannot be nil")
	}

	// Skip compression if article is under 500 characters
	if len(article.Body) < minCharsForCompression {
		return article.Body, nil
	}

	// Check cache
	if cache != nil {
		if cached, ok := cache.GetCached(article.Slug); ok && cached != "" {
			return cached, nil
		}
	}

	// Fallback if no LLM client is available
	if llamaClient == nil {
		return Truncate(article.Body, targetTokens), nil
	}

	prompt := fmt.Sprintf("Summarize the following article in about 2-3 concise paragraphs (25-50%% of original length), preserving key factual details, names, and context:\n\nTitle: %s\n\n%s", article.FrontMatter.Title, Truncate(article.Body, 1500))

	resp, err := llamaClient.Complete(ctx, models.CompletionRequest{
		Prompt:      prompt,
		Temperature: 0.3,
		MaxTokens:   targetTokens,
	})
	if err != nil {
		// Fallback to truncation
		return Truncate(article.Body, targetTokens), nil
	}

	compressed := strings.TrimSpace(resp.Text)
	if compressed == "" {
		compressed = Truncate(article.Body, targetTokens)
	}

	if cache != nil {
		cache.PutCached(article.Slug, article.ContentHash, compressed, len(article.Body))
	}

	return compressed, nil
}
