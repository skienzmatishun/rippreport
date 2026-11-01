package refiner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
	"github.com/rippreport/go-content-pipeline/internal/processor"
	"github.com/rippreport/go-content-pipeline/pkg/gemini"
)

// Refiner handles hook refinement using external APIs
type Refiner struct {
	geminiClient *gemini.Client
	postManager  *processor.PostManager
}

// Config holds refiner configuration
type Config struct {
	GeminiAPIKey string
	PostsDir     string
}

// RefineOptions holds options for refinement
type RefineOptions struct {
	SpecificSlugs []string
	Force         bool
	DryRun        bool
	Model         string
}

// RefineResult holds the results of a refinement operation
type RefineResult struct {
	ProcessedPosts int
	FailedPosts    int
	SkippedPosts   int
	Duration       time.Duration
}

// New creates a new Refiner
func New(cfg Config) (*Refiner, error) {
	var geminiClient *gemini.Client
	if cfg.GeminiAPIKey != "" {
		var err error
		geminiClient, err = gemini.NewClient(gemini.Config{
			APIKey: cfg.GeminiAPIKey,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create Gemini client: %w", err)
		}
	}

	postManager := processor.NewPostManager(cfg.PostsDir)

	return &Refiner{
		geminiClient: geminiClient,
		postManager:  postManager,
	}, nil
}

// RefineHooks refines existing hooks using the Gemini API
func (r *Refiner) RefineHooks(ctx context.Context, opts RefineOptions) (*RefineResult, error) {
	if r.geminiClient == nil {
		return nil, fmt.Errorf("Gemini API key not configured")
	}

	startTime := time.Now()
	result := &RefineResult{}

	// Discover posts to process
	filter := models.PostFilter{
		SpecificSlugs: opts.SpecificSlugs,
	}

	postPaths, err := r.postManager.DiscoverPosts(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to discover posts: %w", err)
	}

	for _, postPath := range postPaths {
		post, err := r.postManager.ReadPost(postPath)
		if err != nil {
			result.FailedPosts++
			continue
		}

		hooksPath := filepath.Join(filepath.Dir(string(postPath)), "sidebar-hooks.yaml")
		
		// Check if hooks file exists
		if _, err := os.Stat(hooksPath); os.IsNotExist(err) {
			result.SkippedPosts++
			continue
		}

		// Load existing hooks
		hooks, err := processor.LoadHooksYAML(hooksPath)
		if err != nil {
			result.FailedPosts++
			continue
		}

		// Check if already refined
		if hooks.Metadata.RefinedBy == "gemini" && !opts.Force {
			result.SkippedPosts++
			continue
		}

		if opts.DryRun {
			fmt.Printf("[DRY RUN] Would refine hooks for %s\n", post.Slug)
			result.SkippedPosts++
			continue
		}

		// Prepare items for refinement
		var items []gemini.HookItem
		for slug, hook := range hooks.Hooks {
			if hook.Brief != "" {
				items = append(items, gemini.HookItem{
					Slug:  slug,
					Brief: hook.Brief,
				})
			}
		}

		if len(items) == 0 {
			result.SkippedPosts++
			continue
		}

		// Refine using Gemini
		refined, err := r.geminiClient.RefineHooks(ctx, items, opts.Model)
		if err != nil {
			fmt.Printf("Failed to refine %s: %v\n", post.Slug, err)
			result.FailedPosts++
			continue
		}

		// Update hooks with refined text
		anyUpdated := false
		for _, refinedItem := range refined {
			if hook, exists := hooks.Hooks[refinedItem.Slug]; exists && refinedItem.Text != "" {
				hook.Text = refinedItem.Text
				hooks.Hooks[refinedItem.Slug] = hook
				anyUpdated = true
			}
		}

		if !anyUpdated {
			result.FailedPosts++
			continue
		}

		// Update metadata
		hooks.Metadata.RefinedBy = "gemini"
		hooks.Metadata.RefinedAt = time.Now()

		// Save updated hooks
		if err := processor.SaveHooksYAML(hooksPath, hooks); err != nil {
			fmt.Printf("Failed to save refined hooks for %s: %v\n", post.Slug, err)
			result.FailedPosts++
			continue
		}

		fmt.Printf("Successfully refined hooks for %s\n", post.Slug)
		result.ProcessedPosts++
	}

	result.Duration = time.Since(startTime)
	return result, nil
}
