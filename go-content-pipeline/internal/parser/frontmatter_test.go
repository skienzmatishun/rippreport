package parser

import (
	"strings"
	"testing"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

func TestParseFrontMatter(t *testing.T) {
	t.Run("valid complete front matter", func(t *testing.T) {
		content := `---
title: "Sample Post Title"
date: 2025-02-14T15:04:05Z
description: "A short sample description"
categories:
  - news
  - politics
tags:
  - alabama
  - fairhope
related_articles:
  - slug: related-one
    title: "Related One"
    score: 90
    rank: 1
custom_author: "Jane Doe"
---

This is the post content.
Line 2 of content.`

		fm, body, delim, err := ParseFrontMatter(content)
		if err != nil {
			t.Fatalf("ParseFrontMatter failed: %v", err)
		}

		if delim != models.DelimiterYAML {
			t.Errorf("expected DelimiterYAML, got %v", delim)
		}
		if fm.Title != "Sample Post Title" {
			t.Errorf("expected title 'Sample Post Title', got %q", fm.Title)
		}
		expectedDate := time.Date(2025, 2, 14, 15, 4, 5, 0, time.UTC)
		if !fm.Date.Equal(expectedDate) {
			t.Errorf("expected date %v, got %v", expectedDate, fm.Date)
		}
		if len(fm.Categories) != 2 || fm.Categories[0] != "news" || fm.Categories[1] != "politics" {
			t.Errorf("unexpected categories: %v", fm.Categories)
		}
		if len(fm.RelatedArticles) != 1 || fm.RelatedArticles[0].Slug != "related-one" {
			t.Errorf("unexpected related articles: %v", fm.RelatedArticles)
		}
		if fm.Extra == nil || fm.Extra["custom_author"] != "Jane Doe" {
			t.Errorf("expected custom_author preserved in Extra, got: %v", fm.Extra)
		}
		if !strings.Contains(body, "This is the post content.") {
			t.Errorf("body missing expected content: %q", body)
		}
	})

	t.Run("categories as string scalar", func(t *testing.T) {
		content := `---
title: "Single Category"
date: 2025-01-01
categories: opinion
---

Body text`

		fm, _, _, err := ParseFrontMatter(content)
		if err != nil {
			t.Fatalf("ParseFrontMatter failed: %v", err)
		}
		if len(fm.Categories) != 1 || fm.Categories[0] != "opinion" {
			t.Errorf("expected categories [opinion], got %v", fm.Categories)
		}
	})

	t.Run("date-only format", func(t *testing.T) {
		content := `---
title: "Date Only"
date: 2024-11-05
---

Body`

		fm, _, _, err := ParseFrontMatter(content)
		if err != nil {
			t.Fatalf("ParseFrontMatter failed: %v", err)
		}
		expected := time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC)
		if !fm.Date.Equal(expected) {
			t.Errorf("expected %v, got %v", expected, fm.Date)
		}
	})

	t.Run("error missing delimiter", func(t *testing.T) {
		content := `title: No Delimiter`
		_, _, _, err := ParseFrontMatter(content)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("error unclosed delimiter", func(t *testing.T) {
		content := `---
title: Unclosed
`
		_, _, _, err := ParseFrontMatter(content)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestSerializeFrontMatterAndRoundTrip(t *testing.T) {
	fm := &models.FrontMatter{
		Title:       "Round Trip Test",
		Date:        time.Date(2025, 3, 10, 12, 0, 0, 0, time.UTC),
		Description: "Testing round-trip serialization",
		Categories:  []string{"news", "investigation"},
		Tags:        []string{"tag1", "tag2"},
		RelatedArticles: []models.RelatedArticle{
			{
				Slug:  "article-1",
				Title: "Article 1",
				Score: 95,
				Rank:  1,
			},
		},
		Extra: map[string]interface{}{
			"custom_id": 42,
		},
	}

	serialized, err := SerializeFrontMatter(fm, models.DelimiterYAML)
	if err != nil {
		t.Fatalf("SerializeFrontMatter failed: %v", err)
	}

	if !strings.HasPrefix(serialized, "---\n") {
		t.Errorf("serialized front matter should start with ---\\n, got %q", serialized)
	}
	if !strings.HasSuffix(serialized, "---\n") {
		t.Errorf("serialized front matter should end with ---\\n, got %q", serialized)
	}

	// Parse it back
	fullDoc := serialized + "\nSome markdown body content."
	parsedFm, body, delim, err := ParseFrontMatter(fullDoc)
	if err != nil {
		t.Fatalf("ParseFrontMatter on serialized content failed: %v", err)
	}

	if delim != models.DelimiterYAML {
		t.Errorf("expected DelimiterYAML, got %v", delim)
	}
	if parsedFm.Title != fm.Title {
		t.Errorf("title mismatch: %q vs %q", parsedFm.Title, fm.Title)
	}
	if !parsedFm.Date.Equal(fm.Date) {
		t.Errorf("date mismatch: %v vs %v", parsedFm.Date, fm.Date)
	}
	if len(parsedFm.RelatedArticles) != len(fm.RelatedArticles) {
		t.Errorf("related articles length mismatch: %d vs %d", len(parsedFm.RelatedArticles), len(fm.RelatedArticles))
	}
	if !strings.Contains(body, "Some markdown body content.") {
		t.Errorf("body content mismatch: %q", body)
	}
}

func TestFormatYAML(t *testing.T) {
	data := map[string]interface{}{
		"name":  "test",
		"items": []string{"a", "b"},
	}

	yamlStr, err := FormatYAML(data)
	if err != nil {
		t.Fatalf("FormatYAML failed: %v", err)
	}

	// Verify 2-space indentation
	if !strings.Contains(yamlStr, "items:\n  - a\n  - b\n") {
		t.Errorf("expected 2-space indentation in YAML output, got:\n%s", yamlStr)
	}
}
