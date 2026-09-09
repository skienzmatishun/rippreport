package models

import (
	"fmt"
	"time"
)

const HookGeneratorVersion = "1.0.0"

// Hook is a Bridge Brief linking a main article to one related (widget) article.
// YAML on disk is a map keyed by widget slug; WidgetSlug is filled on load.
type Hook struct {
	WidgetSlug  string    `yaml:"-"`
	WidgetTitle string    `yaml:"-"`
	Brief       string    `yaml:"brief"`
	Text        string    `yaml:"text,omitempty"`
	WidgetDate  string    `yaml:"widget_date,omitempty"`
	Generated   time.Time `yaml:"generated_at"`
	ModelName   string    `yaml:"model_name,omitempty"`
}

// HookMetadata is the production sidebar-hooks.yaml header.
type HookMetadata struct {
	GeneratedAt      time.Time `yaml:"generated_at"`
	GeneratorVersion string    `yaml:"generator_version"`
	PostSlug         string    `yaml:"post_slug"`
	MainDate         string    `yaml:"main_date,omitempty"`
	RefinedBy        string    `yaml:"refined_by,omitempty"`
	RefinedAt        time.Time `yaml:"refined_at,omitempty"`
	ModelName        string    `yaml:"model_name,omitempty"`
}

// HookFile is the production sidebar-hooks.yaml document.
type HookFile struct {
	Metadata HookMetadata `yaml:"metadata"`
	Hooks    []*Hook      `yaml:"hooks"`
}

// Validate checks if the Hook has the fields needed to write a brief.
func (h *Hook) Validate() error {
	if h.WidgetSlug == "" {
		return fmt.Errorf("hook widget_slug cannot be empty")
	}
	if h.Brief == "" {
		return fmt.Errorf("hook brief cannot be empty")
	}
	if h.Generated.IsZero() {
		return fmt.Errorf("hook generated timestamp cannot be zero")
	}
	return nil
}

// Validate checks if the HookFile has all required fields
func (hf *HookFile) Validate() error {
	if hf.Metadata.PostSlug == "" {
		return fmt.Errorf("hook file post_slug cannot be empty")
	}
	if hf.Metadata.GeneratorVersion == "" {
		return fmt.Errorf("hook file generator_version cannot be empty")
	}
	if hf.Metadata.GeneratedAt.IsZero() {
		return fmt.Errorf("hook file generated timestamp cannot be zero")
	}
	if len(hf.Hooks) == 0 {
		return fmt.Errorf("hook file must contain at least one hook")
	}

	for i, hook := range hf.Hooks {
		if err := hook.Validate(); err != nil {
			return fmt.Errorf("hook %d: %w", i, err)
		}
	}

	return nil
}

// HookRequest represents a request to generate a hook
type HookRequest struct {
	MainPost      *Post
	WidgetPost    *Post
	UseCompressed bool
}

// Validate checks if the HookRequest has all required fields
func (hr *HookRequest) Validate() error {
	if hr.MainPost == nil {
		return fmt.Errorf("hook request main_post cannot be nil")
	}
	if hr.WidgetPost == nil {
		return fmt.Errorf("hook request widget_post cannot be nil")
	}
	if err := hr.MainPost.Validate(); err != nil {
		return fmt.Errorf("main_post validation failed: %w", err)
	}
	if err := hr.WidgetPost.Validate(); err != nil {
		return fmt.Errorf("widget_post validation failed: %w", err)
	}
	return nil
}
