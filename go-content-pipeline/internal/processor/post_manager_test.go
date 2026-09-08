package processor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

func TestDiscoverPosts(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	// Create test posts
	testPosts := []struct {
		slug       string
		title      string
		date       string
		categories []string
	}{
		{
			slug:       "post-1",
			title:      "First Post",
			date:       "2025-01-15",
			categories: []string{"news"},
		},
		{
			slug:       "post-2",
			title:      "Second Post",
			date:       "2025-01-14",
			categories: []string{"investigation"},
		},
		{
			slug:       "holiday-post",
			title:      "Holiday Post",
			date:       "2025-01-13",
			categories: []string{"holiday"},
		},
		{
			slug:       "backstory-podcast-ep1",
			title:      "Backstory Podcast Episode 1",
			date:       "2025-01-12",
			categories: []string{"backstory"},
		},
		{
			slug:       "wonderful-wednesday",
			title:      "Wonderful Wednesday Special",
			date:       "2025-01-11",
			categories: []string{"news"},
		},
	}

	for _, tp := range testPosts {
		postDir := filepath.Join(tmpDir, tp.slug)
		err := os.MkdirAll(postDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create post directory: %v", err)
		}

		content := createTestPostContent(tp.title, tp.date, tp.categories)
		err = os.WriteFile(filepath.Join(postDir, "index.md"), []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write post file: %v", err)
		}
	}

	pm := NewPostManager(tmpDir)

	t.Run("discover all posts", func(t *testing.T) {
		filter := models.PostFilter{}
		paths, err := pm.DiscoverPosts(filter)
		if err != nil {
			t.Fatalf("DiscoverPosts failed: %v", err)
		}

		// Should exclude holiday by default
		expectedCount := 4
		if len(paths) != expectedCount {
			t.Errorf("Expected %d posts, got %d", expectedCount, len(paths))
		}

		// Check sorting (newest first)
		if len(paths) >= 2 {
			// First post should be post-1 (2025-01-15)
			if !containsSlugInPath(string(paths[0]), "post-1") {
				t.Errorf("Expected first post to be post-1, got %s", paths[0])
			}
		}
	})

	t.Run("exclude categories", func(t *testing.T) {
		filter := models.PostFilter{
			ExcludeCategories: []string{"backstory"},
		}
		paths, err := pm.DiscoverPosts(filter)
		if err != nil {
			t.Fatalf("DiscoverPosts failed: %v", err)
		}

		// Should exclude holiday (default) and backstory
		expectedCount := 3
		if len(paths) != expectedCount {
			t.Errorf("Expected %d posts, got %d", expectedCount, len(paths))
		}

		// Verify backstory post is not included
		for _, path := range paths {
			if containsSlugInPath(string(path), "backstory-podcast-ep1") {
				t.Errorf("Backstory post should be excluded")
			}
		}
	})

	t.Run("exclude titles", func(t *testing.T) {
		filter := models.PostFilter{
			ExcludeTitles: []string{"wonderful wednesday"},
		}
		paths, err := pm.DiscoverPosts(filter)
		if err != nil {
			t.Fatalf("DiscoverPosts failed: %v", err)
		}

		// Should exclude holiday (default) and wonderful wednesday
		expectedCount := 3
		if len(paths) != expectedCount {
			t.Errorf("Expected %d posts, got %d", expectedCount, len(paths))
		}

		// Verify wonderful wednesday post is not included
		for _, path := range paths {
			if containsSlugInPath(string(path), "wonderful-wednesday") {
				t.Errorf("Wonderful Wednesday post should be excluded")
			}
		}
	})

	t.Run("specific slugs", func(t *testing.T) {
		filter := models.PostFilter{
			SpecificSlugs: []string{"post-1", "post-2"},
		}
		paths, err := pm.DiscoverPosts(filter)
		if err != nil {
			t.Fatalf("DiscoverPosts failed: %v", err)
		}

		expectedCount := 2
		if len(paths) != expectedCount {
			t.Errorf("Expected %d posts, got %d", expectedCount, len(paths))
		}

		// Verify only requested slugs are included
		for _, path := range paths {
			pathStr := string(path)
			if !containsSlugInPath(pathStr, "post-1") && !containsSlugInPath(pathStr, "post-2") {
				t.Errorf("Unexpected post in results: %s", path)
			}
		}
	})

	t.Run("date range filter", func(t *testing.T) {
		startDate := time.Date(2025, 1, 14, 0, 0, 0, 0, time.UTC)
		filter := models.PostFilter{
			StartDate: &startDate,
		}
		paths, err := pm.DiscoverPosts(filter)
		if err != nil {
			t.Fatalf("DiscoverPosts failed: %v", err)
		}

		// Should include post-1 (Jan 15) and post-2 (Jan 14), exclude earlier posts and holiday
		expectedCount := 2
		if len(paths) != expectedCount {
			t.Errorf("Expected %d posts, got %d", expectedCount, len(paths))
		}
	})

	t.Run("combined filters", func(t *testing.T) {
		filter := models.PostFilter{
			ExcludeCategories: []string{"investigation"},
			ExcludeTitles:     []string{"backstory podcast"},
		}
		paths, err := pm.DiscoverPosts(filter)
		if err != nil {
			t.Fatalf("DiscoverPosts failed: %v", err)
		}

		// Should have post-1 and wonderful-wednesday (excluding holiday, investigation, backstory)
		expectedCount := 2
		if len(paths) != expectedCount {
			t.Errorf("Expected %d posts, got %d", expectedCount, len(paths))
		}
	})
}

func TestExtractSlugFromPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "unix path",
			path:     "/content/p/my-post/index.md",
			expected: "my-post",
		},
		{
			name:     "relative path",
			path:     "content/p/another-post/index.md",
			expected: "another-post",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSlugFromPath(tt.path)
			if result != tt.expected {
				t.Errorf("Expected slug %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestQuickParseFrontMatter(t *testing.T) {
	pm := NewPostManager("")

	t.Run("parse complete front matter", func(t *testing.T) {
		content := `---
title: "Test Post"
date: 2025-01-15T10:30:00Z
categories:
  - news
  - investigation
---

Post body content here.`

		fm := &models.FrontMatter{}
		err := pm.quickParseFrontMatter(content, fm)
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		if fm.Title != "Test Post" {
			t.Errorf("Expected title 'Test Post', got '%s'", fm.Title)
		}

		expectedDate := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
		if !fm.Date.Equal(expectedDate) {
			t.Errorf("Expected date %v, got %v", expectedDate, fm.Date)
		}

		if len(fm.Categories) != 2 {
			t.Errorf("Expected 2 categories, got %d", len(fm.Categories))
		}

		if !containsIgnoreCase(fm.Categories, "news") {
			t.Errorf("Expected category 'news' not found")
		}
	})

	t.Run("parse date-only format", func(t *testing.T) {
		content := `---
title: "Test Post"
date: 2025-01-15
---

Body`

		fm := &models.FrontMatter{}
		err := pm.quickParseFrontMatter(content, fm)
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		expectedDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
		if !fm.Date.Equal(expectedDate) {
			t.Errorf("Expected date %v, got %v", expectedDate, fm.Date)
		}
	})

	t.Run("parse inline categories", func(t *testing.T) {
		content := `---
title: "Test Post"
date: 2025-01-15
categories: [news, investigation]
---

Body`

		fm := &models.FrontMatter{}
		err := pm.quickParseFrontMatter(content, fm)
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		if len(fm.Categories) != 2 {
			t.Errorf("Expected 2 categories, got %d", len(fm.Categories))
		}
	})
}

// Helper function to create test post content
func createTestPostContent(title, date string, categories []string) string {
	content := "---\n"
	content += "title: \"" + title + "\"\n"
	content += "date: " + date + "\n"
	if len(categories) > 0 {
		content += "categories:\n"
		for _, cat := range categories {
			content += "  - " + cat + "\n"
		}
	}
	content += "---\n\n"
	content += "This is the post body content.\n"
	return content
}

// Helper function to check if a path contains a slug
func containsSlugInPath(path, slug string) bool {
	// Normalize path separators for cross-platform compatibility
	normalizedPath := filepath.ToSlash(path)
	return filepath.Base(filepath.Dir(normalizedPath)) == slug
}

func TestGetEmbeddingContent(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	pm := NewPostManager(tmpDir)

	t.Run("non-backstory post returns body", func(t *testing.T) {
		// Create a non-backstory post
		postDir := filepath.Join(tmpDir, "regular-post")
		err := os.MkdirAll(postDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create post directory: %v", err)
		}

		bodyContent := "This is the regular post body content."
		postContent := createTestPostContent("Regular Post", "2025-01-15", []string{"news"})
		indexPath := filepath.Join(postDir, "index.md")
		err = os.WriteFile(indexPath, []byte(postContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write index.md: %v", err)
		}

		post := &models.Post{
			Path: models.PostPath(indexPath),
			Slug: "regular-post",
			Body: bodyContent,
			FrontMatter: &models.FrontMatter{
				Title:      "Regular Post",
				Date:       time.Now(),
				Categories: []string{"news"},
			},
		}

		content, err := pm.GetEmbeddingContent(post)
		if err != nil {
			t.Fatalf("GetEmbeddingContent failed: %v", err)
		}

		if content != bodyContent {
			t.Errorf("Expected body content '%s', got '%s'", bodyContent, content)
		}
	})

	t.Run("backstory post with transcript returns transcript", func(t *testing.T) {
		// Create a backstory post with transcript.md
		postDir := filepath.Join(tmpDir, "backstory-post-with-transcript")
		err := os.MkdirAll(postDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create post directory: %v", err)
		}

		bodyContent := "This is the index.md body content."
		transcriptContent := "This is the transcript content with full podcast details."
		postContent := createTestPostContent("Backstory Episode", "2025-01-15", []string{"Backstory Podcast"})
		indexPath := filepath.Join(postDir, "index.md")
		transcriptPath := filepath.Join(postDir, "transcript.md")

		err = os.WriteFile(indexPath, []byte(postContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write index.md: %v", err)
		}

		err = os.WriteFile(transcriptPath, []byte(transcriptContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write transcript.md: %v", err)
		}

		post := &models.Post{
			Path: models.PostPath(indexPath),
			Slug: "backstory-post-with-transcript",
			Body: bodyContent,
			FrontMatter: &models.FrontMatter{
				Title:      "Backstory Episode",
				Date:       time.Now(),
				Categories: []string{"Backstory Podcast"},
			},
		}

		content, err := pm.GetEmbeddingContent(post)
		if err != nil {
			t.Fatalf("GetEmbeddingContent failed: %v", err)
		}

		if content != transcriptContent {
			t.Errorf("Expected transcript content '%s', got '%s'", transcriptContent, content)
		}
	})

	t.Run("backstory post without transcript returns body", func(t *testing.T) {
		// Create a backstory post without transcript.md
		postDir := filepath.Join(tmpDir, "backstory-post-no-transcript")
		err := os.MkdirAll(postDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create post directory: %v", err)
		}

		bodyContent := "This is the backstory post body without transcript."
		postContent := createTestPostContent("Backstory Episode No Transcript", "2025-01-15", []string{"Backstory Podcast"})
		indexPath := filepath.Join(postDir, "index.md")

		err = os.WriteFile(indexPath, []byte(postContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write index.md: %v", err)
		}

		post := &models.Post{
			Path: models.PostPath(indexPath),
			Slug: "backstory-post-no-transcript",
			Body: bodyContent,
			FrontMatter: &models.FrontMatter{
				Title:      "Backstory Episode No Transcript",
				Date:       time.Now(),
				Categories: []string{"Backstory Podcast"},
			},
		}

		content, err := pm.GetEmbeddingContent(post)
		if err != nil {
			t.Fatalf("GetEmbeddingContent failed: %v", err)
		}

		if content != bodyContent {
			t.Errorf("Expected body content '%s', got '%s'", bodyContent, content)
		}
	})

	t.Run("backstory post with case-insensitive category match", func(t *testing.T) {
		// Test that category matching is case-insensitive
		postDir := filepath.Join(tmpDir, "backstory-lowercase")
		err := os.MkdirAll(postDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create post directory: %v", err)
		}

		bodyContent := "Body content"
		transcriptContent := "Transcript with lowercase category."
		postContent := createTestPostContent("Lowercase Backstory", "2025-01-15", []string{"backstory"})
		indexPath := filepath.Join(postDir, "index.md")
		transcriptPath := filepath.Join(postDir, "transcript.md")

		err = os.WriteFile(indexPath, []byte(postContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write index.md: %v", err)
		}

		err = os.WriteFile(transcriptPath, []byte(transcriptContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write transcript.md: %v", err)
		}

		post := &models.Post{
			Path: models.PostPath(indexPath),
			Slug: "backstory-lowercase",
			Body: bodyContent,
			FrontMatter: &models.FrontMatter{
				Title:      "Lowercase Backstory",
				Date:       time.Now(),
				Categories: []string{"backstory"},
			},
		}

		content, err := pm.GetEmbeddingContent(post)
		if err != nil {
			t.Fatalf("GetEmbeddingContent failed: %v", err)
		}

		if content != transcriptContent {
			t.Errorf("Expected transcript content with lowercase category, got body instead")
		}
	})

	t.Run("backstory post with empty transcript returns empty string", func(t *testing.T) {
		// Create a backstory post with empty transcript.md
		postDir := filepath.Join(tmpDir, "backstory-empty-transcript")
		err := os.MkdirAll(postDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create post directory: %v", err)
		}

		bodyContent := "This is the body content."
		postContent := createTestPostContent("Empty Transcript", "2025-01-15", []string{"Backstory Podcast"})
		indexPath := filepath.Join(postDir, "index.md")
		transcriptPath := filepath.Join(postDir, "transcript.md")

		err = os.WriteFile(indexPath, []byte(postContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write index.md: %v", err)
		}

		// Create empty transcript
		err = os.WriteFile(transcriptPath, []byte(""), 0644)
		if err != nil {
			t.Fatalf("Failed to write transcript.md: %v", err)
		}

		post := &models.Post{
			Path: models.PostPath(indexPath),
			Slug: "backstory-empty-transcript",
			Body: bodyContent,
			FrontMatter: &models.FrontMatter{
				Title:      "Empty Transcript",
				Date:       time.Now(),
				Categories: []string{"Backstory Podcast"},
			},
		}

		content, err := pm.GetEmbeddingContent(post)
		if err != nil {
			t.Fatalf("GetEmbeddingContent failed: %v", err)
		}

		if content != "" {
			t.Errorf("Expected empty transcript content, got '%s'", content)
		}
	})

	t.Run("post with multiple categories including backstory", func(t *testing.T) {
		// Create a post with multiple categories, one being backstory
		postDir := filepath.Join(tmpDir, "multi-category-backstory")
		err := os.MkdirAll(postDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create post directory: %v", err)
		}

		bodyContent := "Body content"
		transcriptContent := "Transcript for multi-category post."
		postContent := `---
title: "Multi Category Post"
date: 2025-01-15
categories:
  - news
  - Backstory Podcast
  - investigation
---

Body content`
		indexPath := filepath.Join(postDir, "index.md")
		transcriptPath := filepath.Join(postDir, "transcript.md")

		err = os.WriteFile(indexPath, []byte(postContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write index.md: %v", err)
		}

		err = os.WriteFile(transcriptPath, []byte(transcriptContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write transcript.md: %v", err)
		}

		post := &models.Post{
			Path: models.PostPath(indexPath),
			Slug: "multi-category-backstory",
			Body: bodyContent,
			FrontMatter: &models.FrontMatter{
				Title:      "Multi Category Post",
				Date:       time.Now(),
				Categories: []string{"news", "Backstory Podcast", "investigation"},
			},
		}

		content, err := pm.GetEmbeddingContent(post)
		if err != nil {
			t.Fatalf("GetEmbeddingContent failed: %v", err)
		}

		if content != transcriptContent {
			t.Errorf("Expected transcript content for multi-category post, got body instead")
		}
	})
}

func TestReadPost(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	t.Run("read valid post with complete front matter", func(t *testing.T) {
		// Create a test post
		postDir := filepath.Join(tmpDir, "test-post-1")
		err := os.MkdirAll(postDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create post directory: %v", err)
		}

		content := `---
title: "Test Post Title"
date: 2025-01-15T10:30:00Z
description: "A test post description"
categories:
  - news
  - investigation
tags:
  - fairhope
  - corruption
related_articles:
  - slug: related-post
    title: "Related Post Title"
    score: 95
    rank: 1
---

This is the post body content.
It has multiple lines.

## Heading

More content here.`

		postPath := filepath.Join(postDir, "index.md")
		err = os.WriteFile(postPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write post file: %v", err)
		}

		pm := NewPostManager(tmpDir)
		post, err := pm.ReadPost(models.PostPath(postPath))
		if err != nil {
			t.Fatalf("ReadPost failed: %v", err)
		}

		// Verify post fields
		if post.Slug != "test-post-1" {
			t.Errorf("Expected slug 'test-post-1', got '%s'", post.Slug)
		}

		if post.Path != models.PostPath(postPath) {
			t.Errorf("Expected path '%s', got '%s'", postPath, post.Path)
		}

		// Verify front matter
		if post.FrontMatter.Title != "Test Post Title" {
			t.Errorf("Expected title 'Test Post Title', got '%s'", post.FrontMatter.Title)
		}

		expectedDate := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
		if !post.FrontMatter.Date.Equal(expectedDate) {
			t.Errorf("Expected date %v, got %v", expectedDate, post.FrontMatter.Date)
		}

		if post.FrontMatter.Description != "A test post description" {
			t.Errorf("Expected description 'A test post description', got '%s'", post.FrontMatter.Description)
		}

		if len(post.FrontMatter.Categories) != 2 {
			t.Errorf("Expected 2 categories, got %d", len(post.FrontMatter.Categories))
		}

		if len(post.FrontMatter.Tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(post.FrontMatter.Tags))
		}

		if len(post.FrontMatter.RelatedArticles) != 1 {
			t.Errorf("Expected 1 related article, got %d", len(post.FrontMatter.RelatedArticles))
		}

		if len(post.FrontMatter.RelatedArticles) > 0 {
			ra := post.FrontMatter.RelatedArticles[0]
			if ra.Slug != "related-post" {
				t.Errorf("Expected related article slug 'related-post', got '%s'", ra.Slug)
			}
			if ra.Score != 95 {
				t.Errorf("Expected related article score 95, got %d", ra.Score)
			}
		}

		// Verify body
		expectedBody := "\nThis is the post body content.\nIt has multiple lines.\n\n## Heading\n\nMore content here."
		if post.Body != expectedBody {
			t.Errorf("Body mismatch.\nExpected:\n%s\nGot:\n%s", expectedBody, post.Body)
		}

		// Verify content hash was calculated
		if post.ContentHash == "" {
			t.Errorf("Content hash should not be empty")
		}

		// Verify content hash is consistent
		expectedHash := models.CalculateContentHash(post.Body)
		if post.ContentHash != expectedHash {
			t.Errorf("Content hash mismatch. Expected %s, got %s", expectedHash, post.ContentHash)
		}
	})

	t.Run("read post with date-only format", func(t *testing.T) {
		postDir := filepath.Join(tmpDir, "test-post-2")
		err := os.MkdirAll(postDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create post directory: %v", err)
		}

		content := `---
title: "Post with Date Only"
date: 2025-01-20
---

Body content.`

		postPath := filepath.Join(postDir, "index.md")
		err = os.WriteFile(postPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write post file: %v", err)
		}

		pm := NewPostManager(tmpDir)
		post, err := pm.ReadPost(models.PostPath(postPath))
		if err != nil {
			t.Fatalf("ReadPost failed: %v", err)
		}

		// Verify date parsing
		expectedDate := time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC)
		if !post.FrontMatter.Date.Equal(expectedDate) {
			t.Errorf("Expected date %v, got %v", expectedDate, post.FrontMatter.Date)
		}
	})

	t.Run("read post with minimal front matter", func(t *testing.T) {
		postDir := filepath.Join(tmpDir, "test-post-3")
		err := os.MkdirAll(postDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create post directory: %v", err)
		}

		content := `---
title: "Minimal Post"
date: 2025-01-15
---

Body.`

		postPath := filepath.Join(postDir, "index.md")
		err = os.WriteFile(postPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write post file: %v", err)
		}

		pm := NewPostManager(tmpDir)
		post, err := pm.ReadPost(models.PostPath(postPath))
		if err != nil {
			t.Fatalf("ReadPost failed: %v", err)
		}

		if post.FrontMatter.Title != "Minimal Post" {
			t.Errorf("Expected title 'Minimal Post', got '%s'", post.FrontMatter.Title)
		}
	})

	t.Run("error on missing title", func(t *testing.T) {
		postDir := filepath.Join(tmpDir, "test-post-4")
		err := os.MkdirAll(postDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create post directory: %v", err)
		}

		content := `---
date: 2025-01-15
---

Body.`

		postPath := filepath.Join(postDir, "index.md")
		err = os.WriteFile(postPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write post file: %v", err)
		}

		pm := NewPostManager(tmpDir)
		_, err = pm.ReadPost(models.PostPath(postPath))
		if err == nil {
			t.Errorf("Expected error for missing title, got none")
		}
	})

	t.Run("error on missing date", func(t *testing.T) {
		postDir := filepath.Join(tmpDir, "test-post-5")
		err := os.MkdirAll(postDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create post directory: %v", err)
		}

		content := `---
title: "No Date Post"
---

Body.`

		postPath := filepath.Join(postDir, "index.md")
		err = os.WriteFile(postPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write post file: %v", err)
		}

		pm := NewPostManager(tmpDir)
		_, err = pm.ReadPost(models.PostPath(postPath))
		if err == nil {
			t.Errorf("Expected error for missing date, got none")
		}
	})

	t.Run("error on invalid yaml", func(t *testing.T) {
		postDir := filepath.Join(tmpDir, "test-post-6")
		err := os.MkdirAll(postDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create post directory: %v", err)
		}

		content := `---
title: "Invalid YAML
date: 2025-01-15
---

Body.`

		postPath := filepath.Join(postDir, "index.md")
		err = os.WriteFile(postPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write post file: %v", err)
		}

		pm := NewPostManager(tmpDir)
		_, err = pm.ReadPost(models.PostPath(postPath))
		if err == nil {
			t.Errorf("Expected error for invalid YAML, got none")
		}
	})

	t.Run("error on missing front matter delimiter", func(t *testing.T) {
		postDir := filepath.Join(tmpDir, "test-post-7")
		err := os.MkdirAll(postDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create post directory: %v", err)
		}

		content := `title: "No Delimiter"
date: 2025-01-15

Body.`

		postPath := filepath.Join(postDir, "index.md")
		err = os.WriteFile(postPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write post file: %v", err)
		}

		pm := NewPostManager(tmpDir)
		_, err = pm.ReadPost(models.PostPath(postPath))
		if err == nil {
			t.Errorf("Expected error for missing delimiter, got none")
		}
	})

	t.Run("error on file not found", func(t *testing.T) {
		pm := NewPostManager(tmpDir)
		_, err := pm.ReadPost(models.PostPath("/nonexistent/path/index.md"))
		if err == nil {
			t.Errorf("Expected error for nonexistent file, got none")
		}
	})
}

func TestUpdateRelatedArticles(t *testing.T) {
	tmpDir := t.TempDir()
	pm := NewPostManager(tmpDir)

	t.Run("successful update preserving all fields and body", func(t *testing.T) {
		postDir := filepath.Join(tmpDir, "update-test")
		if err := os.MkdirAll(postDir, 0755); err != nil {
			t.Fatalf("failed to create post dir: %v", err)
		}

		initialContent := `---
title: "Article to Update"
date: 2025-01-20T12:00:00Z
description: "Original description"
categories:
  - news
  - politics
tags:
  - alabama
custom_extra_key: "custom_value_123"
---

# Heading 1

This is the original body text that must be preserved exactly.`

		postFile := filepath.Join(postDir, "index.md")
		if err := os.WriteFile(postFile, []byte(initialContent), 0644); err != nil {
			t.Fatalf("failed to write post: %v", err)
		}

		newArticles := []models.RelatedArticle{
			{
				Slug:  "related-one",
				Title: "Related Post One",
				Score: 98,
				Rank:  1,
			},
			{
				Slug:  "related-two",
				Title: "Related Post Two",
				Score: 88,
				Rank:  2,
			},
		}

		err := pm.UpdateRelatedArticles(models.PostPath(postFile), newArticles)
		if err != nil {
			t.Fatalf("UpdateRelatedArticles failed: %v", err)
		}

		// Read updated post back
		updatedPost, err := pm.ReadPost(models.PostPath(postFile))
		if err != nil {
			t.Fatalf("ReadPost on updated post failed: %v", err)
		}

		// Verify related articles updated
		if len(updatedPost.FrontMatter.RelatedArticles) != 2 {
			t.Fatalf("expected 2 related articles, got %d", len(updatedPost.FrontMatter.RelatedArticles))
		}
		if updatedPost.FrontMatter.RelatedArticles[0].Slug != "related-one" || updatedPost.FrontMatter.RelatedArticles[0].Score != 98 {
			t.Errorf("unexpected first article: %+v", updatedPost.FrontMatter.RelatedArticles[0])
		}
		if updatedPost.FrontMatter.RelatedArticles[1].Slug != "related-two" || updatedPost.FrontMatter.RelatedArticles[1].Rank != 2 {
			t.Errorf("unexpected second article: %+v", updatedPost.FrontMatter.RelatedArticles[1])
		}

		// Verify all other fields preserved
		if updatedPost.FrontMatter.Title != "Article to Update" {
			t.Errorf("expected title preserved, got %q", updatedPost.FrontMatter.Title)
		}
		if updatedPost.FrontMatter.Description != "Original description" {
			t.Errorf("expected description preserved, got %q", updatedPost.FrontMatter.Description)
		}
		if len(updatedPost.FrontMatter.Categories) != 2 || updatedPost.FrontMatter.Categories[0] != "news" {
			t.Errorf("expected categories preserved, got %v", updatedPost.FrontMatter.Categories)
		}
		if len(updatedPost.FrontMatter.Tags) != 1 || updatedPost.FrontMatter.Tags[0] != "alabama" {
			t.Errorf("expected tags preserved, got %v", updatedPost.FrontMatter.Tags)
		}
		if updatedPost.FrontMatter.Extra == nil || updatedPost.FrontMatter.Extra["custom_extra_key"] != "custom_value_123" {
			t.Errorf("expected custom_extra_key preserved, got %v", updatedPost.FrontMatter.Extra)
		}

		// Verify body preserved
		if !strings.Contains(updatedPost.Body, "This is the original body text that must be preserved exactly.") {
			t.Errorf("body text was not preserved, got: %q", updatedPost.Body)
		}
	})

	t.Run("validation error on invalid article", func(t *testing.T) {
		postDir := filepath.Join(tmpDir, "invalid-article-test")
		if err := os.MkdirAll(postDir, 0755); err != nil {
			t.Fatalf("failed to create post dir: %v", err)
		}
		postFile := filepath.Join(postDir, "index.md")
		initialContent := `---
title: "Valid Post"
date: 2025-01-20
---

Body`
		if err := os.WriteFile(postFile, []byte(initialContent), 0644); err != nil {
			t.Fatalf("failed to write post: %v", err)
		}

		invalidArticles := []models.RelatedArticle{
			{
				Slug:  "", // empty slug is invalid
				Title: "Invalid",
				Score: 50,
				Rank:  1,
			},
		}

		err := pm.UpdateRelatedArticles(models.PostPath(postFile), invalidArticles)
		if err == nil {
			t.Error("expected error for invalid article, got nil")
		}
	})

	t.Run("error on nonexistent post file", func(t *testing.T) {
		err := pm.UpdateRelatedArticles(models.PostPath("/nonexistent/file/index.md"), nil)
		if err == nil {
			t.Error("expected error for nonexistent post file, got nil")
		}
	})
}

