package processor

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

func TestGenerateHook_Basic(t *testing.T) {
	main := &models.Post{
		Path: "/p/main/index.md",
		Slug: "main",
		FrontMatter: &models.FrontMatter{
			Title: "Grand Hotel Renovation",
			Date:  time.Now(),
		},
		Body: "Long article describing the history and preservation of the historic Grand Hotel.",
	}
	widget := &models.Post{
		Path: "/p/widget/index.md",
		Slug: "widget",
		FrontMatter: &models.FrontMatter{
			Title: "Point Clear Coastal Heritage",
			Date:  time.Now(),
		},
		Body: "Historic overview of Point Clear, Alabama.",
	}

	req := models.HookRequest{
		MainPost:   main,
		WidgetPost: widget,
	}

	hook, err := GenerateHook(context.Background(), nil, "test-model", nil, req)
	if err != nil {
		t.Fatalf("GenerateHook failed: %v", err)
	}

	if hook.WidgetSlug != "widget" {
		t.Errorf("expected widget slug 'widget', got %s", hook.WidgetSlug)
	}
	if hook.WidgetTitle != "Point Clear Coastal Heritage" {
		t.Errorf("expected widget title, got %s", hook.WidgetTitle)
	}
	if hook.Brief == "" {
		t.Errorf("brief should not be empty")
	}
}

func TestSaveAndLoadHooks(t *testing.T) {
	tmpDir := t.TempDir()
	postDir := filepath.Join(tmpDir, "sample-post")
	if err := os.MkdirAll(postDir, 0755); err != nil {
		t.Fatalf("failed to create post dir: %v", err)
	}

	hooks := []*models.Hook{
		{
			WidgetSlug:  "related-one",
			WidgetTitle: "Related One",
			Brief:       "Discover how this connects to related topic one.",
			Generated:   time.Now(),
			ModelName:   "test-model",
		},
		{
			WidgetSlug:  "related-two",
			WidgetTitle: "Related Two",
			Brief:       "Deep dive into the financial records.",
			Generated:   time.Now(),
			ModelName:   "test-model",
		},
	}

	meta := models.HookMetadata{
		GeneratedAt:      time.Now().UTC(),
		GeneratorVersion: models.HookGeneratorVersion,
		PostSlug:         "sample-post",
	}
	if err := SaveHooks(postDir, meta, hooks); err != nil {
		t.Fatalf("SaveHooks failed: %v", err)
	}

	loaded, err := LoadHooks(postDir)
	if err != nil {
		t.Fatalf("LoadHooks failed: %v", err)
	}

	if len(loaded) != 2 {
		t.Fatalf("expected 2 loaded hooks, got %d", len(loaded))
	}
	if loaded[0].WidgetSlug != "related-one" || loaded[1].WidgetSlug != "related-two" {
		t.Errorf("unexpected loaded hooks: %+v", loaded)
	}
}
