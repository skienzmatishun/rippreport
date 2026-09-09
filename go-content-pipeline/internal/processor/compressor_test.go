package processor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

func TestEstimateTokens(t *testing.T) {
	if EstimateTokens("") != 0 {
		t.Errorf("expected 0 for empty text")
	}
	if EstimateTokens("1234") != 1 {
		t.Errorf("expected 1 token for 4 chars, got %d", EstimateTokens("1234"))
	}
	if EstimateTokens("12345678") != 2 {
		t.Errorf("expected 2 tokens for 8 chars, got %d", EstimateTokens("12345678"))
	}
}

func TestTruncate(t *testing.T) {
	longText := strings.Repeat("A", 100) + strings.Repeat("Z", 100)
	// 200 chars ≈ 50 tokens. Truncate to 20 tokens (~80 chars: 48 head, 32 tail)
	truncated := Truncate(longText, 20)

	if !strings.Contains(truncated, "[... content truncated ...]") {
		t.Errorf("expected truncation marker in %s", truncated)
	}
	if !strings.HasPrefix(truncated, "AAAA") {
		t.Errorf("expected prefix preserved")
	}
	if !strings.HasSuffix(truncated, "ZZZZ") {
		t.Errorf("expected suffix preserved")
	}
}

func TestCompressionCache_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "compression_cache.json")

	cache := NewCompressionCache(cachePath, "test-model")
	cache.PutCached("post-slug", "hash123", "Compressed text", 1000)

	val, ok := cache.GetCached("post-slug")
	if !ok || val != "Compressed text" {
		t.Fatalf("expected hit, got %s, ok=%v", val, ok)
	}

	if err := cache.Persist(); err != nil {
		t.Fatalf("Persist failed: %v", err)
	}

	cache2 := NewCompressionCache(cachePath, "test-model")
	if err := cache2.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	val, ok = cache2.GetCached("post-slug")
	if !ok || val != "Compressed text" {
		t.Fatalf("expected hit from loaded cache, got %s, ok=%v", val, ok)
	}
}

func TestCompress_SkipsShortArticle(t *testing.T) {
	article := &models.Post{
		Slug: "short-one",
		Body: "Short body under 500 characters.",
		FrontMatter: &models.FrontMatter{
			Title: "Short",
		},
	}

	compressed, err := Compress(context.Background(), nil, nil, article, 100)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}
	if compressed != article.Body {
		t.Errorf("expected original body for short article, got %s", compressed)
	}
}

func TestCompressionCache_PolymorphicLoad(t *testing.T) {
	tmpDir := t.TempDir()
	rawMapPath := filepath.Join(tmpDir, "raw_articles.json")

	// Write raw map[string]string JSON
	rawJSON := `{"article-1": "Summary 1", "article-2": "Summary 2"}`
	if err := os.WriteFile(rawMapPath, []byte(rawJSON), 0644); err != nil {
		t.Fatalf("failed to write raw JSON: %v", err)
	}

	cache := NewCompressionCache(rawMapPath, "test-model")
	if err := cache.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cache.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", cache.Len())
	}

	s1, ok1 := cache.GetCached("article-1")
	if !ok1 || s1 != "Summary 1" {
		t.Errorf("expected 'Summary 1', got %q (ok=%v)", s1, ok1)
	}

	s2, ok2 := cache.GetCached("article-2")
	if !ok2 || s2 != "Summary 2" {
		t.Errorf("expected 'Summary 2', got %q (ok=%v)", s2, ok2)
	}
}

func TestCompressionCache_ActualCompressedArticlesJSON(t *testing.T) {
	// Points to the real compressed_articles.json in internal/cache
	actualPath := filepath.Join("..", "cache", "compressed_articles.json")
	if _, err := os.Stat(actualPath); os.IsNotExist(err) {
		t.Skip("internal/cache/compressed_articles.json not found from test dir")
	}

	cache := NewCompressionCache(actualPath, "hook-model")
	if err := cache.Load(); err != nil {
		t.Fatalf("Load failed on actual compressed_articles.json: %v", err)
	}

	if cache.Len() < 700 {
		t.Errorf("expected at least 700 articles in compressed_articles.json, got %d", cache.Len())
	}

	// Verify known slug
	deafEars, ok := cache.GetCached("deaf-ears")
	if !ok || !strings.Contains(deafEars, "KEY ENTITIES") {
		t.Errorf("expected 'deaf-ears' summary with KEY ENTITIES, got %q", deafEars)
	}
}

