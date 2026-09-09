package test

import (
	"context"
	"os"
	"testing"

	"github.com/rippreport/go-content-pipeline/internal/cache"
	"github.com/rippreport/go-content-pipeline/internal/orchestrator"
	"github.com/rippreport/go-content-pipeline/internal/processor"
	"github.com/rippreport/go-content-pipeline/internal/testutil"
)

// Requirement 22.2: Test embedding generation workflow
// - Generate embeddings for test posts
// - Verify cache persistence
// - Verify cache invalidation on model change
// - Verify progress tracking
func TestEmbeddingWorkflow_GenerationAndPersistence(t *testing.T) {
	env := SetupTestEnv(t)
	posts := env.CreateSamplePosts(3)

	ctx := context.Background()
	res, err := env.Orchestrator.GenerateRelatedArticles(ctx, orchestrator.RankingOptions{
		VectorOnly: true,
	})
	testutil.AssertNoError(t, err, "GenerateRelatedArticles should succeed")
	testutil.AssertEqual(t, res.ProcessedPosts, 3, "Processed posts count")
	testutil.AssertEqual(t, res.FailedPosts, 0, "Failed posts count")

	// 1. Verify cache persistence
	testutil.AssertTrue(t, testutil.FileExists(t, env.Config.Storage.CacheFile), "Embedding cache file should exist on disk")
	cacheData, err := os.ReadFile(env.Config.Storage.CacheFile)
	testutil.AssertNoError(t, err, "Read cache file")
	testutil.AssertTrue(t, len(cacheData) > 0, "Cache file should not be empty")

	// 2. Verify cache hits on subsequent run
	initialEmbedReqCount := len(env.MockServer.EmbeddingRequests)
	testutil.AssertTrue(t, initialEmbedReqCount >= 3, "Expected at least 3 embedding requests sent to server")

	// Create new orchestrator reusing the cache
	reloadedCache := cache.NewEmbeddingCache(env.Config.Storage.CacheFile, env.Config.LlamaServer.EmbeddingModel)
	reloadedTracker := processor.NewProgressTracker(env.Config.Storage.ProgressFile)
	orch2 := orchestrator.NewOrchestrator(
		env.Config,
		env.LlamaClient,
		env.PostManager,
		reloadedCache,
		reloadedTracker,
		env.BackupManager,
		env.StagingManager,
		env.CompressionCache,
	)

	res2, err := orch2.GenerateRelatedArticles(ctx, orchestrator.RankingOptions{
		VectorOnly: true,
	})
	testutil.AssertNoError(t, err, "Second run should succeed")
	testutil.AssertEqual(t, res2.CacheHits, 3, "All 3 posts should hit the embedding cache")
	testutil.AssertEqual(t, len(env.MockServer.EmbeddingRequests), initialEmbedReqCount, "No new embedding requests should be sent to server")

	// 3. Verify progress tracking
	for _, p := range posts {
		testutil.AssertTrue(t, env.ProgressTracker.IsProcessed(p.Slug), "Post should be marked processed in tracker: "+p.Slug)
		testutil.AssertTrue(t, reloadedTracker.IsProcessed(p.Slug), "Post should be marked processed in reloaded tracker: "+p.Slug)
	}
	testutil.AssertTrue(t, testutil.FileExists(t, env.Config.Storage.ProgressFile), "Progress tracker file should exist on disk")
}

func TestEmbeddingWorkflow_ModelChangeInvalidation(t *testing.T) {
	env := SetupTestEnv(t)
	env.CreateSamplePosts(2)

	ctx := context.Background()
	_, err := env.Orchestrator.GenerateRelatedArticles(ctx, orchestrator.RankingOptions{
		VectorOnly: true,
	})
	testutil.AssertNoError(t, err, "Initial generation should succeed")
	testutil.AssertTrue(t, env.EmbeddingCache.Count() >= 2, "Cache should have at least 2 entries")

	// Simulate changing model name
	newModelCache := cache.NewEmbeddingCache(env.Config.Storage.CacheFile, "brand-new-different-model")
	// Cache should invalidate old entries because model name does not match
	testutil.AssertEqual(t, newModelCache.Count(), 0, "Cache should be empty/invalidated when initialized with a different model name")
}
