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

// Requirement 22.6: Test refresh recent workflow
// - Generate initial rankings
// - Generate new rankings with refresh recent mode
// - Verify merge preserves older articles
func TestRefreshRecentWorkflow(t *testing.T) {
	env := SetupTestEnv(t)
	now := time.Now()

	// 1. Create a recent post that already has pre-existing related articles (from past runs)
	mainPost := testutil.PostWithRelated("recent-post-alpha", []string{"historical-article-1", "historical-article-2"})
	mainPost.Date = now.AddDate(0, -1, 0) // 1 month ago (within 2 years)
	mainPost.Content = "Recent discussion regarding municipal budget adjustments and allocations."
	env.CreatePost(mainPost)

	// Create historical posts (older, e.g. 5 years ago)
	hist1 := testutil.DefaultPost("historical-article-1")
	hist1.Date = now.AddDate(-5, 0, 0)
	hist1.Content = "Historical context about municipal financial history from five years ago."
	env.CreatePost(hist1)

	hist2 := testutil.DefaultPost("historical-article-2")
	hist2.Date = now.AddDate(-5, 0, 0)
	hist2.Content = "Historical background on regional government bonding."
	env.CreatePost(hist2)

	// Create newly published posts (e.g. today and last week)
	newPost1 := testutil.DefaultPost("breaking-news-1")
	newPost1.Date = now
	newPost1.Content = "Breaking news: New allocations for county projects announced today."
	env.CreatePost(newPost1)

	newPost2 := testutil.DefaultPost("breaking-news-2")
	newPost2.Date = now.AddDate(0, 0, -3)
	newPost2.Content = "Updates on municipal appropriations and fiscal year reviews."
	env.CreatePost(newPost2)

	ctx := context.Background()

	// 2. Run with RefreshRecent = true
	res, err := env.Orchestrator.GenerateRelatedArticles(ctx, orchestrator.RankingOptions{
		RefreshRecent: true,
		VectorOnly:    true,
	})
	testutil.AssertNoError(t, err, "GenerateRelatedArticles with RefreshRecent should succeed")
	testutil.AssertTrue(t, res.ProcessedPosts > 0, "Processed posts")

	// 3. Read back recent-post-alpha
	postPath := models.PostPath(filepath.Join(env.PostsDir, "recent-post-alpha", "index.md"))
	updatedPost, err := env.PostManager.ReadPost(postPath)
	testutil.AssertNoError(t, err, "Read updated post")

	related := updatedPost.FrontMatter.RelatedArticles
	testutil.AssertTrue(t, len(related) >= 2, "Should have multiple related articles")

	// Verify sequential ranking
	for i, r := range related {
		testutil.AssertEqual(t, r.Rank, i+1, "Rank must be 1-indexed sequential")
	}

	// Verify older articles were preserved through the merge
	foundHistorical := false
	for _, r := range related {
		if r.Slug == "historical-article-1" || r.Slug == "historical-article-2" {
			foundHistorical = true
			break
		}
	}
	testutil.AssertTrue(t, foundHistorical, "Historical articles should be preserved after merge with recent articles")
}
