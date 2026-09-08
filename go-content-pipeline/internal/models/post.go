package models

import (
	"fmt"
	"strings"
	"time"
)

// PostPath represents the file path to a Hugo post
type PostPath string

// Delimiter represents the front matter delimiter type
type Delimiter int

const (
	// DelimiterYAML represents YAML front matter (---)
	DelimiterYAML Delimiter = iota
	// DelimiterTOML represents TOML front matter (+++)
	DelimiterTOML
)

// Post represents a Hugo blog post with front matter and content
type Post struct {
	Path        PostPath     `json:"path"`
	Slug        string       `json:"slug"`
	FrontMatter *FrontMatter `json:"front_matter"`
	Body        string       `json:"body"`
	ContentHash string       `json:"content_hash"`
}

// FrontMatter contains YAML metadata from a Hugo post
type FrontMatter struct {
	Title           string                 `yaml:"title"`
	Date            time.Time              `yaml:"date"`
	Description     string                 `yaml:"description,omitempty"`
	Categories      []string               `yaml:"categories,omitempty"`
	Tags            []string               `yaml:"tags,omitempty"`
	AltTags         string                 `yaml:"alttags,omitempty"`
	RelatedArticles []RelatedArticle       `yaml:"related_articles,omitempty"`
	Draft           bool                   `yaml:"draft,omitempty"`
	Extra           map[string]interface{} `yaml:",inline"`
}

// RelatedArticle represents a related article recommendation
type RelatedArticle struct {
	Slug  string `yaml:"slug"`
	Title string `yaml:"title"`
	Score int    `yaml:"score"`
	Rank  int    `yaml:"rank"`
}

// Validate checks if the Post has all required fields
func (p *Post) Validate() error {
	if p.Path == "" {
		return fmt.Errorf("post path cannot be empty")
	}
	if p.Slug == "" {
		return fmt.Errorf("post slug cannot be empty")
	}
	if p.FrontMatter == nil {
		return fmt.Errorf("post front matter cannot be nil")
	}
	return p.FrontMatter.Validate()
}

// HasCategory checks if the post has a specific category
func (p *Post) HasCategory(category string) bool {
	if p.FrontMatter == nil {
		return false
	}
	for _, cat := range p.FrontMatter.Categories {
		if strings.EqualFold(cat, category) {
			return true
		}
	}
	return false
}

// IsBackstory checks if the post is a backstory podcast post.
// Returns true if any category contains "backstory" (case-insensitive).
func (p *Post) IsBackstory() bool {
	if p.FrontMatter == nil {
		return false
	}
	for _, cat := range p.FrontMatter.Categories {
		if strings.Contains(strings.ToLower(cat), "backstory") {
			return true
		}
	}
	return false
}

// UpdateContentHash calculates and updates the content hash based on the post body.
// This should be called whenever the post content is loaded or modified.
func (p *Post) UpdateContentHash() {
	p.ContentHash = CalculateContentHash(p.Body)
}

// Validate checks if the FrontMatter has all required fields
func (fm *FrontMatter) Validate() error {
	if fm.Title == "" {
		return fmt.Errorf("front matter title cannot be empty")
	}
	if fm.Date.IsZero() {
		return fmt.Errorf("front matter date cannot be zero")
	}

	// Validate related articles if present
	for i, article := range fm.RelatedArticles {
		if err := article.Validate(); err != nil {
			return fmt.Errorf("related article %d: %w", i, err)
		}
	}

	return nil
}

// Validate checks if the RelatedArticle has all required fields
func (ra *RelatedArticle) Validate() error {
	if ra.Slug == "" {
		return fmt.Errorf("related article slug cannot be empty")
	}
	if ra.Title == "" {
		return fmt.Errorf("related article title cannot be empty")
	}
	if ra.Score < 0 || ra.Score > 100 {
		return fmt.Errorf("related article score must be between 0 and 100, got %d", ra.Score)
	}
	if ra.Rank < 1 {
		return fmt.Errorf("related article rank must be at least 1, got %d", ra.Rank)
	}
	return nil
}

// PostFilter defines criteria for filtering posts
type PostFilter struct {
	ExcludeCategories []string
	ExcludeTitles     []string
	StartDate         *time.Time
	EndDate           *time.Time
	SpecificSlugs     []string
}

// Matches checks if a post matches the filter criteria
func (pf *PostFilter) Matches(post *Post) bool {
	// Check specific slugs first (allowlist)
	if len(pf.SpecificSlugs) > 0 {
		found := false
		for _, slug := range pf.SpecificSlugs {
			if post.Slug == slug {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check category exclusions
	for _, excludeCat := range pf.ExcludeCategories {
		if post.HasCategory(excludeCat) {
			return false
		}
	}

	// Check title exclusions (case-insensitive)
	if post.FrontMatter != nil {
		lowerTitle := strings.ToLower(post.FrontMatter.Title)
		for _, excludePhrase := range pf.ExcludeTitles {
			if strings.Contains(lowerTitle, strings.ToLower(excludePhrase)) {
				return false
			}
		}
	}

	// Check date range
	if post.FrontMatter != nil {
		postDate := post.FrontMatter.Date
		if pf.StartDate != nil && postDate.Before(*pf.StartDate) {
			return false
		}
		if pf.EndDate != nil && postDate.After(*pf.EndDate) {
			return false
		}
	}

	return true
}
