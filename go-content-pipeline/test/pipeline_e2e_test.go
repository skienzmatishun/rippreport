package test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
	"github.com/rippreport/go-content-pipeline/internal/orchestrator"
	"github.com/rippreport/go-content-pipeline/internal/processor"
	"github.com/rippreport/go-content-pipeline/internal/testutil"
)

// Requirement 22.8: Test full pipeline end-to-end
// - Run complete pipeline on test dataset
// - Verify all outputs correct
// - Verify performance acceptable
// - Verify memory usage acceptable
func TestFullPipeline_EndToEnd(t *testing.T) {
	env := SetupTestEnv(t)

	// Create diverse test dataset
	p1 := testutil.DefaultPost("e2e-county-politics")
	p1.Categories = []string{"government", "politics"}
	p1.Content = "Detailed investigation into county council contracts and spending oversight."
	env.CreatePost(p1)

	p2 := testutil.DefaultPost("e2e-budget-audit")
	p2.Categories = []string{"government", "finance"}
	p2.Content = "An audit of the annual county budget reveals discrepancies in infrastructure allocations."
	env.CreatePost(p2)

	p3 := testutil.ElectionPost("e2e-commissioner-race")
	p3.Content = "Voters head to the polls to elect new county commissioners in closely watched race."
	env.CreatePost(p3)

	p4 := testutil.ShortPost("e2e-quick-notice")
	p4.Content = "Public notice: The town hall meeting has been moved to Thursday evening at 7 PM."
	env.CreatePost(p4)

	p5 := testutil.BackstoryPost("e2e-backstory-podcast")
	p5.Date = time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	env.CreatePost(p5)

	slugs := []string{
		"e2e-county-politics",
		"e2e-budget-audit",
		"e2e-commissioner-race",
		"e2e-quick-notice",
		"e2e-backstory-podcast",
	}

	startAlloc := getHeapAlloc()
	startTime := time.Now()

	ctx := context.Background()
	pipelineRes, err := env.Orchestrator.RunPipeline(ctx, orchestrator.PipelineOptions{
		RankingOptions: orchestrator.RankingOptions{
			VectorOnly: false, // Exercises similarity + LLM scoring + composite
		},
		HookOptions: orchestrator.HookOptions{
			UseCompressed: true,
		},
	})

	duration := time.Since(startTime)
	endAlloc := getHeapAlloc()

	testutil.AssertNoError(t, err, "RunPipeline should execute successfully")
	testutil.AssertEqual(t, pipelineRes.Ranking.ProcessedPosts, 5, "All 5 posts should be processed in ranking")
	testutil.AssertEqual(t, pipelineRes.Ranking.FailedPosts, 0, "No posts should fail")
	testutil.AssertTrue(t, pipelineRes.Hooks.TotalHooks > 0, "Hooks should be generated")

	// Verify Summary
	testutil.AssertEqual(t, pipelineRes.Summary.TotalPosts, 5, "Total posts in summary")
	testutil.AssertEqual(t, pipelineRes.Summary.ProcessedPosts, 5, "Processed posts in summary")
	testutil.AssertEqual(t, pipelineRes.Summary.FailedPosts, 0, "Failed posts in summary")
	testutil.AssertEqual(t, pipelineRes.Summary.EmbeddingsComputed, 5, "Embeddings computed")

	// Verify Front Matter in modified posts
	for _, slug := range slugs {
		postPath := models.PostPath(filepath.Join(env.PostsDir, slug, "index.md"))
		post, readErr := env.PostManager.ReadPost(postPath)
		testutil.AssertNoError(t, readErr, "Read updated post: "+slug)

		// Non-backstory posts have related articles; backstory only has related if other backstory posts exist
		if !post.IsBackstory() {
			testutil.AssertTrue(t, len(post.FrontMatter.RelatedArticles) > 0,
				"Post should have related articles: "+slug)
			for _, ra := range post.FrontMatter.RelatedArticles {
				testutil.AssertNotEqual(t, ra.Slug, slug, "Cannot recommend self")
			}
		}
	}

	// Verify sidebar hooks generated for posts that have related articles
	hookFile := filepath.Join(env.PostsDir, "e2e-county-politics", "sidebar-hooks.yaml")
	testutil.AssertTrue(t, testutil.FileExists(t, hookFile), "sidebar-hooks.yaml should exist for e2e-county-politics")

	hooks, err := processor.LoadHooks(filepath.Join(env.PostsDir, "e2e-county-politics"))
	testutil.AssertNoError(t, err, "LoadHooks")
	testutil.AssertTrue(t, len(hooks) > 0, "Loaded hooks should not be empty")

	// Verify storage persistence
	testutil.AssertTrue(t, testutil.FileExists(t, env.Config.Storage.CacheFile), "Embedding cache saved")
	testutil.AssertTrue(t, testutil.FileExists(t, env.Config.Storage.ProgressFile), "Progress file saved")

	// Performance and memory verification
	// Target: < 500MB heap memory usage (Requirement 19.1, Design performance)
	var memUsedMB float64
	if endAlloc >= startAlloc {
		memUsedMB = float64(endAlloc-startAlloc) / (1024 * 1024)
	} else {
		memUsedMB = float64(endAlloc) / (1024 * 1024)
	}
	t.Logf("Full pipeline completed in %v, heap memory: %.2f MB", duration, memUsedMB)
	testutil.AssertTrue(t, memUsedMB < 500.0, "Memory usage should stay well under 500MB limit")
}

func getHeapAlloc() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}
