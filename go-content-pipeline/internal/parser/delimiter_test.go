package parser

import (
	"strings"
	"testing"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

func TestDetectDelimiter(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		want        models.Delimiter
		wantErr     bool
		errContains string
	}{
		{
			name: "valid YAML delimiter",
			content: `---
title: Test Post
date: 2025-01-20
---
Post content here`,
			want:    models.DelimiterYAML,
			wantErr: false,
		},
		{
			name: "valid TOML delimiter",
			content: `+++
title = "Test Post"
date = 2025-01-20
+++
Post content here`,
			want:    models.DelimiterTOML,
			wantErr: false,
		},
		{
			name: "YAML delimiter with leading whitespace",
			content: `   ---
title: Test Post
---`,
			want:    models.DelimiterYAML,
			wantErr: false,
		},
		{
			name:    "YAML delimiter with leading newlines",
			content: "\n\n---\ntitle: Test Post\n---",
			want:    models.DelimiterYAML,
			wantErr: false,
		},
		{
			name:    "YAML delimiter with Windows line endings",
			content: "---\r\ntitle: Test Post\r\n---\r\nContent",
			want:    models.DelimiterYAML,
			wantErr: false,
		},
		{
			name:    "YAML delimiter with trailing spaces",
			content: "---   \ntitle: Test Post\n---",
			want:    models.DelimiterYAML,
			wantErr: false,
		},
		{
			name:    "TOML delimiter with leading whitespace",
			content: "\t  +++\ntitle = \"Test\"\n+++",
			want:    models.DelimiterTOML,
			wantErr: false,
		},
		{
			name:        "empty content",
			content:     "",
			want:        0,
			wantErr:     true,
			errContains: "content is empty",
		},
		{
			name:        "whitespace only content",
			content:     "   \n\t\n  \r\n  ",
			want:        0,
			wantErr:     true,
			errContains: "content contains only whitespace",
		},
		{
			name:        "missing delimiter - plain text",
			content:     "This is just plain text with no front matter",
			want:        0,
			wantErr:     true,
			errContains: "invalid or missing front matter delimiter",
		},
		{
			name:        "wrong delimiter - single dash",
			content:     "-\ntitle: Test",
			want:        0,
			wantErr:     true,
			errContains: "invalid or missing front matter delimiter",
		},
		{
			name:        "wrong delimiter - double dash",
			content:     "--\ntitle: Test",
			want:        0,
			wantErr:     true,
			errContains: "invalid or missing front matter delimiter",
		},
		{
			name:        "wrong delimiter - four dashes",
			content:     "----\ntitle: Test",
			want:        0,
			wantErr:     true,
			errContains: "invalid or missing front matter delimiter",
		},
		{
			name:        "wrong delimiter - single plus",
			content:     "+\ntitle: Test",
			want:        0,
			wantErr:     true,
			errContains: "invalid or missing front matter delimiter",
		},
		{
			name:        "wrong delimiter - double plus",
			content:     "++\ntitle: Test",
			want:        0,
			wantErr:     true,
			errContains: "invalid or missing front matter delimiter",
		},
		{
			name:        "wrong delimiter - four plus",
			content:     "++++\ntitle: Test",
			want:        0,
			wantErr:     true,
			errContains: "invalid or missing front matter delimiter",
		},
		{
			name:        "delimiter with extra characters",
			content:     "---extra\ntitle: Test",
			want:        0,
			wantErr:     true,
			errContains: "invalid or missing front matter delimiter",
		},
		{
			name:    "YAML delimiter - single line content",
			content: "---",
			want:    models.DelimiterYAML,
			wantErr: false,
		},
		{
			name:    "TOML delimiter - single line content",
			content: "+++",
			want:    models.DelimiterTOML,
			wantErr: false,
		},
		{
			name:    "delimiter with mixed line endings",
			content: "---\r\ntitle: Test\ndate: 2025-01-20\r\n---",
			want:    models.DelimiterYAML,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectDelimiter(tt.content)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DetectDelimiter() expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("DetectDelimiter() error = %v, should contain %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("DetectDelimiter() unexpected error = %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("DetectDelimiter() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDetectDelimiter_RealWorldExamples tests with actual Hugo post examples
func TestDetectDelimiter_RealWorldExamples(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    models.Delimiter
	}{
		{
			name: "typical Hugo post with YAML",
			content: `---
title: "My Blog Post"
date: 2025-01-20T10:30:00Z
categories:
  - news
  - politics
tags:
  - fairhope
  - corruption
---

This is the post content with multiple paragraphs.

More content here.`,
			want: models.DelimiterYAML,
		},
		{
			name: "Hugo post with TOML",
			content: `+++
title = "My Blog Post"
date = 2025-01-20T10:30:00Z
categories = ["news", "politics"]
tags = ["fairhope", "corruption"]
+++

This is the post content.`,
			want: models.DelimiterTOML,
		},
		{
			name: "Hugo post with minimal front matter",
			content: `---
title: Post
date: 2025-01-20
---
Content`,
			want: models.DelimiterYAML,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectDelimiter(tt.content)
			if err != nil {
				t.Errorf("DetectDelimiter() unexpected error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("DetectDelimiter() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDetectDelimiter_EdgeCases tests boundary conditions
func TestDetectDelimiter_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantErr     bool
		errContains string
	}{
		{
			name:        "null character in content",
			content:     "\x00---\ntitle: Test",
			wantErr:     true,
			errContains: "invalid or missing front matter delimiter",
		},
		{
			name:    "very long leading whitespace",
			content: strings.Repeat(" ", 1000) + "---\ntitle: Test",
			wantErr: false,
		},
		{
			name:    "very long leading newlines",
			content: strings.Repeat("\n", 1000) + "---\ntitle: Test",
			wantErr: false,
		},
		{
			name:        "delimiter with unicode spaces (non-breaking space not trimmed)",
			content:     "\u00A0---\ntitle: Test", // non-breaking space
			wantErr:     true,
			errContains: "invalid or missing front matter delimiter",
		},
		{
			name:    "delimiter at EOF",
			content: "---",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DetectDelimiter(tt.content)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DetectDelimiter() expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("DetectDelimiter() error = %v, should contain %q", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("DetectDelimiter() unexpected error = %v", err)
				}
			}
		})
	}
}
