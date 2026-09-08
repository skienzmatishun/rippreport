package models

import (
	"testing"
	"time"
)

func TestPost_UpdateContentHash(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantLen int
	}{
		{
			name:    "empty body",
			body:    "",
			wantLen: 32,
		},
		{
			name:    "simple body",
			body:    "This is a blog post body.",
			wantLen: 32,
		},
		{
			name: "multiline body",
			body: `# Article Title

This is the body of a blog post.

## Section 1
Content here.

## Section 2
More content.`,
			wantLen: 32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &Post{
				Path: "content/p/test-post/index.md",
				Slug: "test-post",
				Body: tt.body,
			}

			// Initially content hash should be empty
			if post.ContentHash != "" {
				t.Errorf("Initial content hash should be empty, got %v", post.ContentHash)
			}

			// Update content hash
			post.UpdateContentHash()

			// Verify hash was set
			if post.ContentHash == "" {
				t.Error("Content hash should be set after UpdateContentHash()")
			}

			// Verify hash length
			if len(post.ContentHash) != tt.wantLen {
				t.Errorf("Content hash length = %d, want %d", len(post.ContentHash), tt.wantLen)
			}

			// Verify hash is consistent
			expectedHash := CalculateContentHash(tt.body)
			if post.ContentHash != expectedHash {
				t.Errorf("Content hash = %v, want %v", post.ContentHash, expectedHash)
			}
		})
	}
}

func TestPost_UpdateContentHashConsistency(t *testing.T) {
	post := &Post{
		Path: "content/p/test-post/index.md",
		Slug: "test-post",
		Body: "Original content of the blog post.",
	}

	// Update hash multiple times with same content
	post.UpdateContentHash()
	hash1 := post.ContentHash

	post.UpdateContentHash()
	hash2 := post.ContentHash

	if hash1 != hash2 {
		t.Errorf("Hash inconsistency: %v != %v", hash1, hash2)
	}

	// Change content and verify hash changes
	post.Body = "Modified content of the blog post."
	post.UpdateContentHash()
	hash3 := post.ContentHash

	if hash1 == hash3 {
		t.Errorf("Hash should change when content changes, but remained %v", hash1)
	}
}

func TestPost_HasCategory(t *testing.T) {
	tests := []struct {
		name       string
		categories []string
		checkCat   string
		want       bool
	}{
		{
			name:       "exact match",
			categories: []string{"news", "politics"},
			checkCat:   "news",
			want:       true,
		},
		{
			name:       "case insensitive match",
			categories: []string{"News", "Politics"},
			checkCat:   "news",
			want:       true,
		},
		{
			name:       "no match",
			categories: []string{"news", "politics"},
			checkCat:   "sports",
			want:       false,
		},
		{
			name:       "empty categories",
			categories: []string{},
			checkCat:   "news",
			want:       false,
		},
		{
			name:       "nil front matter",
			categories: nil,
			checkCat:   "news",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &Post{
				Path: "content/p/test-post/index.md",
				Slug: "test-post",
			}

			if tt.categories != nil {
				post.FrontMatter = &FrontMatter{
					Title:      "Test Post",
					Date:       time.Now(),
					Categories: tt.categories,
				}
			}

			got := post.HasCategory(tt.checkCat)
			if got != tt.want {
				t.Errorf("HasCategory(%q) = %v, want %v", tt.checkCat, got, tt.want)
			}
		})
	}
}

func TestPost_IsBackstory(t *testing.T) {
	tests := []struct {
		name       string
		categories []string
		want       bool
	}{
		{
			name:       "is backstory",
			categories: []string{"backstory", "podcast"},
			want:       true,
		},
		{
			name:       "is Backstory (case insensitive)",
			categories: []string{"Backstory"},
			want:       true,
		},
		{
			name:       "not backstory",
			categories: []string{"news", "politics"},
			want:       false,
		},
		{
			name:       "empty categories",
			categories: []string{},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &Post{
				Path: "content/p/test-post/index.md",
				Slug: "test-post",
				FrontMatter: &FrontMatter{
					Title:      "Test Post",
					Date:       time.Now(),
					Categories: tt.categories,
				},
			}

			got := post.IsBackstory()
			if got != tt.want {
				t.Errorf("IsBackstory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPost_Validate(t *testing.T) {
	validTime := time.Now()

	tests := []struct {
		name    string
		post    *Post
		wantErr bool
	}{
		{
			name: "valid post",
			post: &Post{
				Path: "content/p/test-post/index.md",
				Slug: "test-post",
				FrontMatter: &FrontMatter{
					Title: "Test Post",
					Date:  validTime,
				},
			},
			wantErr: false,
		},
		{
			name: "missing path",
			post: &Post{
				Path: "",
				Slug: "test-post",
				FrontMatter: &FrontMatter{
					Title: "Test Post",
					Date:  validTime,
				},
			},
			wantErr: true,
		},
		{
			name: "missing slug",
			post: &Post{
				Path: "content/p/test-post/index.md",
				Slug: "",
				FrontMatter: &FrontMatter{
					Title: "Test Post",
					Date:  validTime,
				},
			},
			wantErr: true,
		},
		{
			name: "nil front matter",
			post: &Post{
				Path:        "content/p/test-post/index.md",
				Slug:        "test-post",
				FrontMatter: nil,
			},
			wantErr: true,
		},
		{
			name: "invalid front matter - missing title",
			post: &Post{
				Path: "content/p/test-post/index.md",
				Slug: "test-post",
				FrontMatter: &FrontMatter{
					Title: "",
					Date:  validTime,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid front matter - zero date",
			post: &Post{
				Path: "content/p/test-post/index.md",
				Slug: "test-post",
				FrontMatter: &FrontMatter{
					Title: "Test Post",
					Date:  time.Time{},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.post.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRelatedArticle_Validate(t *testing.T) {
	tests := []struct {
		name    string
		article RelatedArticle
		wantErr bool
	}{
		{
			name: "valid article",
			article: RelatedArticle{
				Slug:  "related-post",
				Title: "Related Post",
				Score: 95,
				Rank:  1,
			},
			wantErr: false,
		},
		{
			name: "missing slug",
			article: RelatedArticle{
				Slug:  "",
				Title: "Related Post",
				Score: 95,
				Rank:  1,
			},
			wantErr: true,
		},
		{
			name: "missing title",
			article: RelatedArticle{
				Slug:  "related-post",
				Title: "",
				Score: 95,
				Rank:  1,
			},
			wantErr: true,
		},
		{
			name: "score too low",
			article: RelatedArticle{
				Slug:  "related-post",
				Title: "Related Post",
				Score: -1,
				Rank:  1,
			},
			wantErr: true,
		},
		{
			name: "score too high",
			article: RelatedArticle{
				Slug:  "related-post",
				Title: "Related Post",
				Score: 101,
				Rank:  1,
			},
			wantErr: true,
		},
		{
			name: "invalid rank",
			article: RelatedArticle{
				Slug:  "related-post",
				Title: "Related Post",
				Score: 95,
				Rank:  0,
			},
			wantErr: true,
		},
		{
			name: "edge case - score 0",
			article: RelatedArticle{
				Slug:  "related-post",
				Title: "Related Post",
				Score: 0,
				Rank:  10,
			},
			wantErr: false,
		},
		{
			name: "edge case - score 100",
			article: RelatedArticle{
				Slug:  "related-post",
				Title: "Related Post",
				Score: 100,
				Rank:  1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.article.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostFilter_Matches(t *testing.T) {
	now := time.Now()
	pastDate := now.AddDate(0, -6, 0) // 6 months ago
	futureDate := now.AddDate(0, 6, 0) // 6 months from now

	tests := []struct {
		name   string
		filter PostFilter
		post   *Post
		want   bool
	}{
		{
			name:   "no filters - should match",
			filter: PostFilter{},
			post: &Post{
				Slug: "test-post",
				FrontMatter: &FrontMatter{
					Title:      "Test Post",
					Date:       now,
					Categories: []string{"news"},
				},
			},
			want: true,
		},
		{
			name: "specific slugs - matches",
			filter: PostFilter{
				SpecificSlugs: []string{"test-post", "another-post"},
			},
			post: &Post{
				Slug: "test-post",
				FrontMatter: &FrontMatter{
					Title: "Test Post",
					Date:  now,
				},
			},
			want: true,
		},
		{
			name: "specific slugs - no match",
			filter: PostFilter{
				SpecificSlugs: []string{"different-post", "another-post"},
			},
			post: &Post{
				Slug: "test-post",
				FrontMatter: &FrontMatter{
					Title: "Test Post",
					Date:  now,
				},
			},
			want: false,
		},
		{
			name: "exclude category - matches",
			filter: PostFilter{
				ExcludeCategories: []string{"holiday"},
			},
			post: &Post{
				Slug: "test-post",
				FrontMatter: &FrontMatter{
					Title:      "Test Post",
					Date:       now,
					Categories: []string{"news"},
				},
			},
			want: true,
		},
		{
			name: "exclude category - excluded",
			filter: PostFilter{
				ExcludeCategories: []string{"holiday", "backstory"},
			},
			post: &Post{
				Slug: "test-post",
				FrontMatter: &FrontMatter{
					Title:      "Test Post",
					Date:       now,
					Categories: []string{"backstory"},
				},
			},
			want: false,
		},
		{
			name: "exclude title - matches",
			filter: PostFilter{
				ExcludeTitles: []string{"backstory podcast", "wonderful wednesday"},
			},
			post: &Post{
				Slug: "test-post",
				FrontMatter: &FrontMatter{
					Title: "Important News Article",
					Date:  now,
				},
			},
			want: true,
		},
		{
			name: "exclude title - excluded",
			filter: PostFilter{
				ExcludeTitles: []string{"backstory podcast", "wonderful wednesday"},
			},
			post: &Post{
				Slug: "test-post",
				FrontMatter: &FrontMatter{
					Title: "Backstory Podcast: Episode 42",
					Date:  now,
				},
			},
			want: false,
		},
		{
			name: "date range - within range",
			filter: PostFilter{
				StartDate: &pastDate,
				EndDate:   &futureDate,
			},
			post: &Post{
				Slug: "test-post",
				FrontMatter: &FrontMatter{
					Title: "Test Post",
					Date:  now,
				},
			},
			want: true,
		},
		{
			name: "date range - before start",
			filter: PostFilter{
				StartDate: &now,
			},
			post: &Post{
				Slug: "test-post",
				FrontMatter: &FrontMatter{
					Title: "Test Post",
					Date:  pastDate,
				},
			},
			want: false,
		},
		{
			name: "date range - after end",
			filter: PostFilter{
				EndDate: &now,
			},
			post: &Post{
				Slug: "test-post",
				FrontMatter: &FrontMatter{
					Title: "Test Post",
					Date:  futureDate,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.post)
			if got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
