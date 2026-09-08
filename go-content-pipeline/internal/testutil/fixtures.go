package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// HugoPostFixture represents a sample Hugo post for testing.
type HugoPostFixture struct {
	Slug            string
	Title           string
	Date            time.Time
	Description     string
	Categories      []string
	Tags            []string
	Content         string
	RelatedArticles []RelatedArticleFixture
}

// RelatedArticleFixture represents a related article in front matter.
type RelatedArticleFixture struct {
	Slug  string
	Title string
	Score int
	Rank  int
}

// DefaultPost returns a basic Hugo post fixture with sensible defaults.
func DefaultPost(slug string) *HugoPostFixture {
	return &HugoPostFixture{
		Slug:        slug,
		Title:       fmt.Sprintf("Test Post: %s", slug),
		Date:        time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		Description: fmt.Sprintf("This is a test post about %s", slug),
		Categories:  []string{"test", "sample"},
		Tags:        []string{"testing", "example"},
		Content: `This is the main content of the test post.

It contains multiple paragraphs to simulate a real blog post.

## Section Header

Here's some more content with details about the topic. This paragraph
is longer to provide realistic content length for embedding generation.

The post should have enough content to be meaningful for similarity calculations.`,
		RelatedArticles: []RelatedArticleFixture{},
	}
}

// BackstoryPost returns a backstory category post fixture.
func BackstoryPost(slug string) *HugoPostFixture {
	post := DefaultPost(slug)
	post.Categories = []string{"backstory"}
	post.Title = fmt.Sprintf("Backstory Podcast: %s", slug)
	post.Content = `This is the web page content for the podcast episode.

For embedding purposes, the transcript.md file should be used instead.`
	return post
}

// ShortPost returns a post with minimal content (under 400 characters).
func ShortPost(slug string) *HugoPostFixture {
	post := DefaultPost(slug)
	post.Content = "This is a very short post with minimal content."
	return post
}

// LongPost returns a post with extensive content.
func LongPost(slug string) *HugoPostFixture {
	post := DefaultPost(slug)
	post.Content = `This is a comprehensive article with extensive content covering multiple aspects of the topic.

## Introduction

The introduction provides context and sets up the main arguments that will be explored
throughout this detailed analysis. We'll examine various perspectives and provide
evidence-based conclusions.

## Background

Understanding the historical context is crucial for appreciating the current situation.
This section delves into the origins and evolution of the topic over time.

### Early Developments

The early phase was characterized by rapid innovation and experimentation. Multiple
approaches were tried, with varying degrees of success.

### Modern Era

Today's landscape is vastly different from those early days. Technological advances
and changing social norms have transformed the field entirely.

## Main Analysis

The central argument of this article revolves around three key points that we'll
explore in depth. Each point builds on the previous one to create a cohesive narrative.

### Point One

First, we must consider the primary factor that influences all other aspects of this
discussion. The evidence clearly shows a pattern that cannot be ignored.

### Point Two

Secondly, the relationship between various stakeholders reveals important dynamics
that shape outcomes. Multiple case studies demonstrate this principle in action.

### Point Three

Finally, the implications for future developments are significant and far-reaching.
Policy makers and practitioners alike must consider these factors carefully.

## Conclusion

In summary, this analysis has demonstrated the complexity of the issue while providing
clear guidance for moving forward. The evidence supports a nuanced approach that
takes multiple perspectives into account.`
	return post
}

// ElectionPost returns a post in the elections category.
func ElectionPost(slug string) *HugoPostFixture {
	post := DefaultPost(slug)
	post.Categories = []string{"elections", "politics"}
	post.Title = fmt.Sprintf("Election Results: %s", slug)
	post.Content = `Analysis of the recent election results shows interesting patterns.

Voter turnout was higher than expected, and several key races had surprising outcomes.
The implications for future governance will be significant.`
	return post
}

// PostWithRelated returns a post that already has related articles configured.
func PostWithRelated(slug string, relatedSlugs []string) *HugoPostFixture {
	post := DefaultPost(slug)
	post.RelatedArticles = make([]RelatedArticleFixture, len(relatedSlugs))
	for i, relSlug := range relatedSlugs {
		post.RelatedArticles[i] = RelatedArticleFixture{
			Slug:  relSlug,
			Title: fmt.Sprintf("Related: %s", relSlug),
			Score: 95 - (i * 5),
			Rank:  i + 1,
		}
	}
	return post
}

// ToMarkdown converts a HugoPostFixture to markdown format with front matter.
func (p *HugoPostFixture) ToMarkdown() string {
	markdown := "---\n"
	markdown += fmt.Sprintf("title: %q\n", p.Title)
	markdown += fmt.Sprintf("date: %s\n", p.Date.Format(time.RFC3339))

	if p.Description != "" {
		markdown += fmt.Sprintf("description: %q\n", p.Description)
	}

	if len(p.Categories) > 0 {
		markdown += "categories:\n"
		for _, cat := range p.Categories {
			markdown += fmt.Sprintf("  - %s\n", cat)
		}
	}

	if len(p.Tags) > 0 {
		markdown += "tags:\n"
		for _, tag := range p.Tags {
			markdown += fmt.Sprintf("  - %s\n", tag)
		}
	}

	if len(p.RelatedArticles) > 0 {
		markdown += "related_articles:\n"
		for _, rel := range p.RelatedArticles {
			markdown += fmt.Sprintf("  - slug: %s\n", rel.Slug)
			markdown += fmt.Sprintf("    title: %q\n", rel.Title)
			markdown += fmt.Sprintf("    score: %d\n", rel.Score)
			markdown += fmt.Sprintf("    rank: %d\n", rel.Rank)
		}
	}

	markdown += "---\n\n"
	markdown += p.Content

	return markdown
}

// CreatePost creates a Hugo post fixture as a file in the test directory.
func (p *HugoPostFixture) CreatePost(t *testing.T, contentDir string) string {
	t.Helper()
	postDir := filepath.Join(contentDir, "p", p.Slug)
	if err := mkdir(postDir); err != nil {
		t.Fatalf("failed to create post directory %s: %v", postDir, err)
	}

	indexPath := filepath.Join(postDir, "index.md")
	WriteFile(t, indexPath, p.ToMarkdown())

	return postDir
}

// CreatePostWithTranscript creates a backstory post with a transcript.md file.
func (p *HugoPostFixture) CreatePostWithTranscript(t *testing.T, contentDir, transcript string) string {
	t.Helper()
	postDir := p.CreatePost(t, contentDir)
	transcriptPath := filepath.Join(postDir, "transcript.md")
	WriteFile(t, transcriptPath, transcript)
	return postDir
}

func mkdir(path string) error {
	return mkdirAll(path, 0755)
}

func mkdirAll(path string, perm int) error {
	return osMkdirAll(path, perm)
}

func osMkdirAll(path string, perm int) error {
	return finallyCallOsMkdirAll(path, perm)
}

func finallyCallOsMkdirAll(path string, perm int) error {
	return os.MkdirAll(path, os.FileMode(perm))
}

// SampleHookYAML returns sample sidebar hooks YAML content.
func SampleHookYAML(postSlug string) string {
	return fmt.Sprintf(`post_slug: %s
generated: 2024-01-15T10:30:00Z
model_name: test-model
hooks:
  - widget_slug: related-post-1
    widget_title: "Related Article One"
    brief: "This article provides important context about the same topic."
    generated: 2024-01-15T10:30:00Z
    model_name: test-model
  - widget_slug: related-post-2
    widget_title: "Related Article Two"
    brief: "This piece explores a connected issue with relevant insights."
    generated: 2024-01-15T10:30:00Z
    model_name: test-model
`, postSlug)
}

// SampleStagedYAML returns sample staged rankings YAML content.
func SampleStagedYAML(postSlug string) string {
	return fmt.Sprintf(`post_slug: %s
generated: 2024-01-15T10:30:00Z
articles:
  - slug: related-1
    title: "First Related Post"
    score: 95
    rank: 1
  - slug: related-2
    title: "Second Related Post"
    score: 88
    rank: 2
  - slug: related-3
    title: "Third Related Post"
    score: 82
    rank: 3
`, postSlug)
}
