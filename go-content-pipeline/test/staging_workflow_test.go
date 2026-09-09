package test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/rippreport/go-content-pipeline/internal/models"
	"github.com/rippreport/go-content-pipeline/internal/orchestrator"
	"github.com/rippreport/go-content-pipeline/internal/processor"
	"github.com/rippreport/go-content-pipeline/internal/testutil"
)

// Requirement 22.5: Test staging workflow
// - Generate staged rankings
// - Generate hooks from staged rankings
// - Apply staged rankings
// - Verify backups created
// - Clear staged files
func TestStagingWorkflow_FullLifecycle(t *testing.T) {
	env := SetupTestEnv(t)
	posts := env.CreateSamplePosts(3)

	ctx := context.Background()

	// 1. Generate staged rankings
	res, err := env.Orchestrator.GenerateRelatedArticles(ctx, orchestrator.RankingOptions{
		Stage:      true,
		VectorOnly: true,
	})
	testutil.AssertNoError(t, err, "GenerateRelatedArticles with Stage=true should succeed")
	testutil.AssertEqual(t, res.ProcessedPosts, 3, "All posts processed to staging")

	// Verify post front matter is NOT modified yet
	for _, p := range posts {
		postPath := models.PostPath(filepath.Join(env.PostsDir, p.Slug, "index.md"))
		post, err := env.PostManager.ReadPost(postPath)
		testutil.AssertNoError(t, err, "Read post")
		testutil.AssertEqual(t, len(post.FrontMatter.RelatedArticles), 0,
			"Post front matter should have no related articles while staged: "+p.Slug)

		// Verify staging file exists
		stagedFile := filepath.Join(env.PostsDir, p.Slug, "related-articles.staged.yaml")
		testutil.AssertTrue(t, testutil.FileExists(t, stagedFile), "Staged file should exist: "+stagedFile)

		stagedArticles, err := env.StagingManager.ReadStaged(p.Slug)
		testutil.AssertNoError(t, err, "ReadStaged")
		testutil.AssertTrue(t, len(stagedArticles) > 0, "Staged articles should be non-empty")
	}

	// 2. Generate hooks from staged rankings
	hookRes, err := env.Orchestrator.GenerateHooks(ctx, orchestrator.HookOptions{
		FromStaging: true,
	})
	testutil.AssertNoError(t, err, "GenerateHooks from staging should succeed")
	testutil.AssertTrue(t, hookRes.TotalHooks >= 3, "Hooks should be generated from staged rankings")

	for _, p := range posts {
		hooks, err := processor.LoadHooks(filepath.Join(env.PostsDir, p.Slug))
		testutil.AssertNoError(t, err, "LoadHooks from post dir")
		testutil.AssertTrue(t, len(hooks) > 0, "Hooks should exist for "+p.Slug)
	}

	// 3. Apply staged rankings
	err = env.Orchestrator.ApplyStaged(ctx, nil)
	testutil.AssertNoError(t, err, "ApplyStaged for all posts should succeed")

	// Verify backups created
	backupEntries, err := os.ReadDir(env.BackupsDir)
	testutil.AssertNoError(t, err, "Read backups directory")
	testutil.AssertTrue(t, len(backupEntries) > 0, "At least one backup directory should be created")

	// Verify post front matter now contains the rankings and staging files are cleaned up
	for _, p := range posts {
		postPath := models.PostPath(filepath.Join(env.PostsDir, p.Slug, "index.md"))
		updatedPost, err := env.PostManager.ReadPost(postPath)
		testutil.AssertNoError(t, err, "Read post after ApplyStaged")
		testutil.AssertTrue(t, len(updatedPost.FrontMatter.RelatedArticles) > 0,
			"Post front matter must now have related articles: "+p.Slug)

		// Staging file should be removed after applying
		stagedFile := filepath.Join(env.PostsDir, p.Slug, "related-articles.staged.yaml")
		testutil.AssertFalse(t, testutil.FileExists(t, stagedFile),
			"Staged file should be deleted after applying: "+stagedFile)
	}
}

func TestStagingWorkflow_ClearStaged(t *testing.T) {
	env := SetupTestEnv(t)
	posts := env.CreateSamplePosts(2)

	ctx := context.Background()
	_, err := env.Orchestrator.GenerateRelatedArticles(ctx, orchestrator.RankingOptions{
		Stage:      true,
		VectorOnly: true,
	})
	testutil.AssertNoError(t, err, "Generate staged rankings")

	for _, p := range posts {
		stagedFile := filepath.Join(env.PostsDir, p.Slug, "related-articles.staged.yaml")
		testutil.AssertTrue(t, testutil.FileExists(t, stagedFile), "Staged file exists before clear")
	}

	// Clear only the first post's staged rankings
	err = env.Orchestrator.ClearStaged(ctx, []string{posts[0].Slug})
	testutil.AssertNoError(t, err, "ClearStaged for single post")

	testutil.AssertFalse(t, testutil.FileExists(t, filepath.Join(env.PostsDir, posts[0].Slug, "related-articles.staged.yaml")),
		"Cleared post staging file should be deleted")
	testutil.AssertTrue(t, testutil.FileExists(t, filepath.Join(env.PostsDir, posts[1].Slug, "related-articles.staged.yaml")),
		"Non-cleared post staging file should remain")
}
