package orchestrator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rippreport/go-content-pipeline/internal/cache"
	"github.com/rippreport/go-content-pipeline/internal/config"
	"github.com/rippreport/go-content-pipeline/internal/processor"
)

func setupTestOrchestrator(t *testing.T) (*Orchestrator, string) {
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	backupDir := filepath.Join(tmpDir, "backups")
	cacheFile := filepath.Join(tmpDir, "embeddings.json")
	progFile := filepath.Join(tmpDir, "progress.json")
	compFile := filepath.Join(tmpDir, "compression.json")

	posts := []struct {
		slug  string
		title string
		body  string
	}{
		{
			slug:  "post-1",
			title: "First Post Title",
			body:  "Body content of first post with interesting discussion.",
		},
		{
			slug:  "post-2",
			title: "Second Post Title",
			body:  "Body content of second post relating to first post.",
		},
		{
			slug:  "post-3",
			title: "Third Post Title",
			body:  "Body content of third post discussing ethics and transparency.",
		},
	}

	for _, p := range posts {
		dir := filepath.Join(contentDir, p.slug)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		content := fmt.Sprintf("---\ntitle: %q\ndate: 2025-01-20\ncategories:\n  - news\n---\n\n%s\n", p.title, p.body)
		if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte(content), 0644); err != nil {
			t.Fatalf("failed to write post: %v", err)
		}
	}

	cfg := &config.Config{
		Processing: config.ProcessingConfig{
			PostsDirectory: contentDir,
			TopCandidates:  2,
			BatchSize:      2,
		},
		Scoring: config.ScoringConfig{
			Weights: config.ScoreWeights{
				Relevance: 0.5,
				Recency:   0.3,
				Length:    0.2,
			},
		},
	}

	pm := processor.NewPostManager(contentDir)
	ec := cache.NewEmbeddingCache(cacheFile, "test-model")
	pt := processor.NewProgressTracker(progFile)
	bm := processor.NewBackupManager(backupDir)
	sm := processor.NewStagingManager(contentDir, pm, bm)
	cc := processor.NewCompressionCache(compFile, "test-model")

	orch := NewOrchestrator(cfg, nil, pm, ec, pt, bm, sm, cc)
	return orch, contentDir
}

func TestOrchestrator_GenerateRelatedArticles_Direct(t *testing.T) {
	orch, _ := setupTestOrchestrator(t)

	res, err := orch.GenerateRelatedArticles(context.Background(), RankingOptions{
		VectorOnly: true,
	})
	if err != nil {
		t.Fatalf("GenerateRelatedArticles failed: %v", err)
	}

	if res.ProcessedPosts != 3 {
		t.Errorf("expected 3 processed posts, got %d", res.ProcessedPosts)
	}
}

func TestOrchestrator_GenerateRelatedArticles_Staged(t *testing.T) {
	orch, contentDir := setupTestOrchestrator(t)

	res, err := orch.GenerateRelatedArticles(context.Background(), RankingOptions{
		Stage:      true,
		VectorOnly: true,
	})
	if err != nil {
		t.Fatalf("GenerateRelatedArticles with staging failed: %v", err)
	}

	if res.ProcessedPosts != 3 {
		t.Errorf("expected 3 processed posts, got %d", res.ProcessedPosts)
	}

	// Verify staging file exists for post-1
	stagedFile := filepath.Join(contentDir, "post-1", "related-articles.staged.yaml")
	if _, err := os.Stat(stagedFile); os.IsNotExist(err) {
		t.Errorf("expected staging file to exist at %s", stagedFile)
	}

	// Test ApplyStaged
	if err := orch.ApplyStaged(context.Background(), []string{"post-1"}); err != nil {
		t.Fatalf("ApplyStaged failed: %v", err)
	}

	// Staging file should now be removed
	if _, err := os.Stat(stagedFile); !os.IsNotExist(err) {
		t.Errorf("staging file should be removed after ApplyStaged")
	}
}

func TestOrchestrator_RunPipeline(t *testing.T) {
	orch, _ := setupTestOrchestrator(t)

	res, err := orch.RunPipeline(context.Background(), PipelineOptions{
		RankingOptions: RankingOptions{
			VectorOnly: true,
		},
		HookOptions: HookOptions{
			BatchSize: 2,
		},
	})
	if err != nil {
		t.Fatalf("RunPipeline failed: %v", err)
	}

	if res.Ranking.ProcessedPosts != 3 {
		t.Errorf("expected 3 ranking processed posts, got %d", res.Ranking.ProcessedPosts)
	}
	if res.Hooks.TotalHooks == 0 {
		t.Errorf("expected generated hooks > 0, got %d", res.Hooks.TotalHooks)
	}
}
