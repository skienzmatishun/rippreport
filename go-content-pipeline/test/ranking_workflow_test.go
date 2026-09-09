package test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
	"github.com/rippreport/go-content-pipeline/internal/orchestrator"
	"github.com/rippreport/go-content-pipeline/internal/testutil"
)

// Requirement 22.3: Test ranking workflow
// - Generate rankings for test posts
// - Verify similarity calculation
// - Verify LLM scoring
// - Verify composite scoring
// - Verify rankings written to front matter
func TestRankingWorkflow_FrontMatterUpdate(t *testing.T) {
	env := SetupTestEnv(t)
	posts := env.CreateSamplePosts(4)

	ctx := context.Background()
	res, err := env.Orchestrator.GenerateRelatedArticles(ctx, orchestrator.RankingOptions{
		VectorOnly: false, // Exercises LLM scoring + composite score calculation
	})
	testutil.AssertNoError(t, err, "GenerateRelatedArticles with LLM scoring should succeed")
	testutil.AssertEqual(t, res.ProcessedPosts, 4, "Should process all 4 posts")

	// Verify rankings written to front matter
	for _, p := range posts {
		postPath := models.PostPath(filepath.Join(env.PostsDir, p.Slug, "index.md"))
		updatedPost, err := env.PostManager.ReadPost(postPath)
		testutil.AssertNoError(t, err, "Read updated post: "+p.Slug)

		testutil.AssertTrue(t, len(updatedPost.FrontMatter.RelatedArticles) > 0,
			"Post should have related articles: "+p.Slug)
		testutil.AssertTrue(t, len(updatedPost.FrontMatter.RelatedArticles) <= env.Config.Processing.TopCandidates,
			"Post should have at most top_candidates related articles")

		// Verify fields, sequential ranks, and self-exclusion
		for rankIdx, rel := range updatedPost.FrontMatter.RelatedArticles {
			testutil.AssertNotEqual(t, rel.Slug, p.Slug, "Article must not recommend itself")
			testutil.AssertTrue(t, rel.Title != "", "Related article title must not be empty")
			testutil.AssertTrue(t, rel.Score >= 0 && rel.Score <= 100, "Related article score in [0, 100]")
			testutil.AssertEqual(t, rel.Rank, rankIdx+1, "Rank must be 1-indexed sequential")
		}
	}
}

func TestRankingWorkflow_BackstorySeparation(t *testing.T) {
	env := SetupTestEnv(t)

	// Create regular posts
	reg1 := testutil.DefaultPost("news-article-1")
	reg1.Content = "Local municipal news and zoning updates."
	env.CreatePost(reg1)

	reg2 := testutil.DefaultPost("news-article-2")
	reg2.Content = "County budget hearings and financial forecasts."
	env.CreatePost(reg2)

	// Create backstory posts with different dates
	bs1 := testutil.BackstoryPost("backstory-ep-1")
	bs1.Date = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	env.CreatePost(bs1)

	bs2 := testutil.BackstoryPost("backstory-ep-2")
	bs2.Date = time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC) // Newer
	env.CreatePost(bs2)

	ctx := context.Background()
	_, err := env.Orchestrator.GenerateRelatedArticles(ctx, orchestrator.RankingOptions{})
	testutil.AssertNoError(t, err, "GenerateRelatedArticles should succeed")

	// 1. Regular posts should NOT have backstory posts in their recommendations
	updatedReg1, err := env.PostManager.ReadPost(models.PostPath(filepath.Join(env.PostsDir, "news-article-1", "index.md")))
	testutil.AssertNoError(t, err, "Read regular post")
	for _, rel := range updatedReg1.FrontMatter.RelatedArticles {
		testutil.AssertNotEqual(t, rel.Slug, "backstory-ep-1", "Regular post should not recommend backstory post")
		testutil.AssertNotEqual(t, rel.Slug, "backstory-ep-2", "Regular post should not recommend backstory post")
	}

	// 2. Backstory posts should rank other backstory posts by descending date
	updatedBs1, err := env.PostManager.ReadPost(models.PostPath(filepath.Join(env.PostsDir, "backstory-ep-1", "index.md")))
	testutil.AssertNoError(t, err, "Read backstory post")
	testutil.AssertTrue(t, len(updatedBs1.FrontMatter.RelatedArticles) > 0, "Backstory post should have related backstory post")
	testutil.AssertEqual(t, updatedBs1.FrontMatter.RelatedArticles[0].Slug, "backstory-ep-2", "Newer backstory post ranked first")
	testutil.AssertEqual(t, updatedBs1.FrontMatter.RelatedArticles[0].Score, 100, "Backstory simple score starts at 100")
}

func TestRankingWorkflow_VectorOnlyMode(t *testing.T) {
	env := SetupTestEnv(t)
	env.CreateSamplePosts(3)

	initialCompletions := len(env.MockServer.CompletionRequests)

	ctx := context.Background()
	res, err := env.Orchestrator.GenerateRelatedArticles(ctx, orchestrator.RankingOptions{
		VectorOnly: true,
	})
	testutil.AssertNoError(t, err, "GenerateRelatedArticles vector only should succeed")
	testutil.AssertEqual(t, res.ProcessedPosts, 3, "Processed posts")

	// In vector-only mode, no LLM completions are used for scoring
	testutil.AssertEqual(t, len(env.MockServer.CompletionRequests), initialCompletions,
		"No completion calls should be made in vector-only mode")
}
