package processor

import (
	"context"
	"testing"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

func TestBatchAssembler(t *testing.T) {
	ba := NewBatchAssembler(3, 100)

	req1 := models.HookRequest{
		MainPost:   &models.Post{Slug: "main"},
		WidgetPost: &models.Post{Slug: "w1"},
	}
	req2 := models.HookRequest{
		MainPost:   &models.Post{Slug: "main"},
		WidgetPost: &models.Post{Slug: "w2"},
	}
	req3 := models.HookRequest{
		MainPost:   &models.Post{Slug: "main"},
		WidgetPost: &models.Post{Slug: "w3"},
	}
	req4 := models.HookRequest{
		MainPost:   &models.Post{Slug: "main"},
		WidgetPost: &models.Post{Slug: "w4"},
	}

	if !ba.Add(req1, 30) {
		t.Errorf("expected req1 to be added")
	}
	if !ba.Add(req2, 30) {
		t.Errorf("expected req2 to be added")
	}
	if !ba.Add(req3, 30) {
		t.Errorf("expected req3 to be added")
	}
	// Exceeds maxBatchSize of 3
	if ba.Add(req4, 10) {
		t.Errorf("expected req4 to be rejected because batch is full")
	}

	flushed := ba.Flush()
	if len(flushed) != 3 {
		t.Errorf("expected 3 flushed items, got %d", len(flushed))
	}
	if ba.Len() != 0 {
		t.Errorf("expected empty batch after flush")
	}
}

func TestScoreBatch_Fallback(t *testing.T) {
	currentPost := &models.Post{
		Slug: "main",
		FrontMatter: &models.FrontMatter{
			Title: "Main Article",
		},
	}
	candidates := []*models.Post{
		{Slug: "c1", FrontMatter: &models.FrontMatter{Title: "Cand 1"}},
		{Slug: "c2", FrontMatter: &models.FrontMatter{Title: "Cand 2"}},
	}

	scores, err := ScoreBatch(context.Background(), nil, "test-model", currentPost, candidates)
	if err != nil {
		t.Fatalf("ScoreBatch failed: %v", err)
	}
	if len(scores) != 2 {
		t.Fatalf("expected 2 scores, got %d", len(scores))
	}
}

func TestGenerateHooksBatch(t *testing.T) {
	main := &models.Post{
		Path: "/p/main/index.md",
		Slug: "main",
		FrontMatter: &models.FrontMatter{
			Title: "Fairhope Council Meeting",
			Date:  time.Now(),
		},
		Body: "Body of main article",
	}
	w1 := &models.Post{
		Path: "/p/w1/index.md",
		Slug: "w1",
		FrontMatter: &models.FrontMatter{
			Title: "Municipal Budget 2025",
			Date:  time.Now(),
		},
		Body: "Body of budget",
	}

	reqs := []models.HookRequest{
		{MainPost: main, WidgetPost: w1},
	}

	hooks, err := GenerateHooksBatch(context.Background(), nil, "model-test", nil, reqs)
	if err != nil {
		t.Fatalf("GenerateHooksBatch failed: %v", err)
	}
	if len(hooks) != 1 {
		t.Fatalf("expected 1 hook, got %d", len(hooks))
	}
	if hooks[0].WidgetSlug != "w1" || hooks[0].WidgetTitle != "Municipal Budget 2025" {
		t.Errorf("unexpected hook: %+v", hooks[0])
	}
}
