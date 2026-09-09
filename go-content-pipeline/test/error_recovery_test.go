package test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/rippreport/go-content-pipeline/internal/cache"
	internalerrors "github.com/rippreport/go-content-pipeline/internal/errors"
	"github.com/rippreport/go-content-pipeline/internal/orchestrator"
	"github.com/rippreport/go-content-pipeline/internal/processor"
	"github.com/rippreport/go-content-pipeline/internal/testutil"
)

// Requirement 22.7: Test error recovery
// - Test resume after interruption
// - Test retry logic for API failures
// - Test handling of corrupted cache files
// - Test handling of invalid posts

func TestErrorRecovery_ResumeAfterInterruption(t *testing.T) {
	env := SetupTestEnv(t)
	posts := env.CreateSamplePosts(4)

	// Simulate post-1 and post-2 already being processed from a prior interrupted run
	_ = env.ProgressTracker.MarkProcessed(posts[0].Slug)
	_ = env.ProgressTracker.MarkProcessed(posts[1].Slug)
	_ = env.ProgressTracker.Persist()

	ctx := context.Background()
	res, err := env.Orchestrator.GenerateRelatedArticles(ctx, orchestrator.RankingOptions{
		Resume:     true,
		VectorOnly: true,
	})
	testutil.AssertNoError(t, err, "GenerateRelatedArticles with resume should succeed")

	// 2 posts skipped, 2 processed
	testutil.AssertEqual(t, res.SkippedPosts, 2, "Previously completed posts should be skipped")
	testutil.AssertEqual(t, res.ProcessedPosts, 2, "Remaining posts should be processed")

	// All 4 should now be processed
	for _, p := range posts {
		testutil.AssertTrue(t, env.ProgressTracker.IsProcessed(p.Slug), "Post should be marked processed: "+p.Slug)
	}
}

func TestErrorRecovery_APIRetryOnServerError(t *testing.T) {
	env := SetupTestEnv(t)

	// Configure mock server to fail once with 500, then succeed on retry
	failCount := 0
	env.MockServer.EmbeddingResponse = func(req testutil.EmbeddingRequest) ([]float32, error) {
		if failCount < 1 {
			failCount++
			return nil, os.ErrDeadlineExceeded // Will trigger HTTP 500 in mock server
		}
		return testutil.TestEmbedding(64, 0.42), nil
	}

	ctx := context.Background()
	vec, err := env.LlamaClient.GenerateEmbedding(ctx, "Test content with transient error")
	testutil.AssertNoError(t, err, "GenerateEmbedding should succeed after retry")
	testutil.AssertTrue(t, len(vec) == 64, "Embedding vector length")
	testutil.AssertTrue(t, failCount >= 1, "Failed at least once before retry succeeded")
}

func TestErrorRecovery_CorruptedCacheFiles(t *testing.T) {
	env := SetupTestEnv(t)

	// Write garbage JSON to cache and progress files
	testutil.WriteFile(t, env.Config.Storage.CacheFile, "{invalid json content!@#$")
	testutil.WriteFile(t, env.Config.Storage.ProgressFile, "THIS IS NOT JSON")

	// Verify embedding cache recovers cleanly
	c := cache.NewEmbeddingCache(env.Config.Storage.CacheFile, env.Config.LlamaServer.EmbeddingModel)
	err := c.Load()
	testutil.AssertNoError(t, err, "EmbeddingCache Load should recover from corruption without error")
	testutil.AssertEqual(t, c.Count(), 0, "Corrupted cache should reset to 0 entries")

	// Verify progress tracker recovers cleanly
	pt := processor.NewProgressTracker(env.Config.Storage.ProgressFile)
	err = pt.Load()
	testutil.AssertNoError(t, err, "ProgressTracker Load should recover from corruption without error")
	testutil.AssertFalse(t, pt.IsProcessed("sample-post-1"), "Corrupted progress should have no processed posts")

	// Pipeline execution should still succeed after corruption recovery
	env.CreateSamplePosts(2)
	res, err := env.Orchestrator.GenerateRelatedArticles(context.Background(), orchestrator.RankingOptions{
		VectorOnly: true,
	})
	testutil.AssertNoError(t, err, "Pipeline should execute cleanly after cache corruption recovery")
	testutil.AssertEqual(t, res.ProcessedPosts, 2, "Posts processed")
}

func TestErrorRecovery_InvalidPostsHandling(t *testing.T) {
	env := SetupTestEnv(t)

	// Create 2 valid posts
	env.CreatePost(testutil.DefaultPost("valid-post-1"))
	env.CreatePost(testutil.DefaultPost("valid-post-2"))

	// Create 1 invalid post with malformed front matter
	invalidDir := filepath.Join(env.PostsDir, "invalid-post")
	testutil.WriteFile(t, filepath.Join(invalidDir, "index.md"), "---\ntitle: [unclosed\n---\nBroken post body")

	ctx := context.Background()
	res, err := env.Orchestrator.GenerateRelatedArticles(ctx, orchestrator.RankingOptions{
		VectorOnly: true,
	})
	testutil.AssertNoError(t, err, "Orchestrator should not crash when encountering invalid post")

	// Valid posts should be processed
	testutil.AssertEqual(t, res.ProcessedPosts, 2, "Valid posts processed")

	// Error should be collected in ErrorCollector
	errCol := env.Orchestrator.ErrorCollector()
	testutil.AssertTrue(t, errCol.Count() > 0, "Error collector should have recorded error for invalid post")
	grouped := errCol.GroupByCategory()
	testutil.AssertTrue(t, len(grouped[internalerrors.ErrorCategoryFileIO]) > 0, "File I/O error should be recorded")
}
