package models

import (
	"fmt"
	"time"
)

// Hook represents a contextual sidebar hook linking related articles
type Hook struct {
	WidgetSlug  string    `yaml:"widget_slug"`
	WidgetTitle string    `yaml:"widget_title"`
	Brief       string    `yaml:"brief"`
	Generated   time.Time `yaml:"generated"`
	ModelName   string    `yaml:"model_name"`
}

// HookFile represents the YAML file containing hooks for a post
type HookFile struct {
	PostSlug  string    `yaml:"post_slug"`
	Generated time.Time `yaml:"generated"`
	ModelName string    `yaml:"model_name"`
	Hooks     []*Hook   `yaml:"hooks"`
}

// Validate checks if the Hook has all required fields
func (h *Hook) Validate() error {
	if h.WidgetSlug == "" {
		return fmt.Errorf("hook widget_slug cannot be empty")
	}
	if h.WidgetTitle == "" {
		return fmt.Errorf("hook widget_title cannot be empty")
	}
	if h.Brief == "" {
		return fmt.Errorf("hook brief cannot be empty")
	}
	if h.ModelName == "" {
		return fmt.Errorf("hook model_name cannot be empty")
	}
	if h.Generated.IsZero() {
		return fmt.Errorf("hook generated timestamp cannot be zero")
	}
	return nil
}

// Validate checks if the HookFile has all required fields
func (hf *HookFile) Validate() error {
	if hf.PostSlug == "" {
		return fmt.Errorf("hook file post_slug cannot be empty")
	}
	if hf.ModelName == "" {
		return fmt.Errorf("hook file model_name cannot be empty")
	}
	if hf.Generated.IsZero() {
		return fmt.Errorf("hook file generated timestamp cannot be zero")
	}
	if len(hf.Hooks) == 0 {
		return fmt.Errorf("hook file must contain at least one hook")
	}
	
	// Validate each hook
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
