package test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/rippreport/go-content-pipeline/internal/orchestrator"
	"github.com/rippreport/go-content-pipeline/internal/processor"
	"github.com/rippreport/go-content-pipeline/internal/testutil"
)

// Requirement 22.4: Test hook generation workflow
// - Generate hooks for test posts
// - Verify hook quality
// - Verify hook file format
// - Verify compressed content handling
func TestHookWorkflow_GenerationAndFormatting(t *testing.T) {
	env := SetupTestEnv(t)

	// Create 2 posts with mutual related articles
	p1 := testutil.PostWithRelated("article-alpha", []string{"article-beta"})
	p1.Content = "Comprehensive discussion about municipal water policy and infrastructure upgrades."
	env.CreatePost(p1)

	p2 := testutil.PostWithRelated("article-beta", []string{"article-alpha"})
	p2.Content = "Environmental impact report on local watershed management and protection."
	env.CreatePost(p2)

	ctx := context.Background()
	res, err := env.Orchestrator.GenerateHooks(ctx, orchestrator.HookOptions{})
	testutil.AssertNoError(t, err, "GenerateHooks should succeed")
	testutil.AssertEqual(t, res.FailedHooks, 0, "No hooks should fail")
	testutil.AssertTrue(t, res.TotalHooks >= 2, "Expected at least 2 hooks generated")

	// Verify sidebar-hooks.yaml file exists and has correct format
	for _, slug := range []string{"article-alpha", "article-beta"} {
		postDir := filepath.Join(env.PostsDir, slug)
		hookFile := filepath.Join(postDir, "sidebar-hooks.yaml")
		testutil.AssertTrue(t, testutil.FileExists(t, hookFile), "sidebar-hooks.yaml should exist for "+slug)

		hooks, err := processor.LoadHooks(postDir)
		testutil.AssertNoError(t, err, "LoadHooks should succeed for "+slug)
		testutil.AssertTrue(t, len(hooks) > 0, "Loaded hooks should not be empty")

		for _, h := range hooks {
			testutil.AssertTrue(t, h.WidgetSlug != "", "WidgetSlug must not be empty")
			testutil.AssertTrue(t, h.WidgetTitle != "", "WidgetTitle must not be empty")
			testutil.AssertTrue(t, h.Brief != "", "Hook brief must not be empty")
			testutil.AssertTrue(t, h.ModelName != "", "Hook model name must be set")
		}
	}
}

func TestHookWorkflow_CompressedContent(t *testing.T) {
	env := SetupTestEnv(t)

	p1 := testutil.PostWithRelated("post-comp-1", []string{"post-comp-2"})
	p1.Content = "A longer article that discusses civic planning and transportation systems in detail."
	env.CreatePost(p1)

	p2 := testutil.PostWithRelated("post-comp-2", []string{"post-comp-1"})
	p2.Content = "A complementary article examining green transit options and urban development."
	env.CreatePost(p2)

	ctx := context.Background()
	res, err := env.Orchestrator.GenerateHooks(ctx, orchestrator.HookOptions{
		UseCompressed: true,
	})
	testutil.AssertNoError(t, err, "GenerateHooks with compression should succeed")
	testutil.AssertTrue(t, res.TotalHooks >= 2, "Total hooks generated")

	// Verify hooks generated and written
	hooks, err := processor.LoadHooks(filepath.Join(env.PostsDir, "post-comp-1"))
	testutil.AssertNoError(t, err, "LoadHooks")
	testutil.AssertTrue(t, len(hooks) > 0, "Hooks should be generated using compressed content")
}

func TestHookWorkflow_DryRun(t *testing.T) {
	env := SetupTestEnv(t)

	p1 := testutil.PostWithRelated("dry-hook-1", []string{"dry-hook-2"})
	env.CreatePost(p1)
	p2 := testutil.PostWithRelated("dry-hook-2", []string{"dry-hook-1"})
	env.CreatePost(p2)

	ctx := context.Background()
	res, err := env.Orchestrator.GenerateHooks(ctx, orchestrator.HookOptions{
		DryRun: true,
	})
	testutil.AssertNoError(t, err, "GenerateHooks dry run should succeed")
	testutil.AssertTrue(t, res.TotalHooks >= 2, "Hooks computed in memory")

	// Verify no files written to disk
	hookFile := filepath.Join(env.PostsDir, "dry-hook-1", "sidebar-hooks.yaml")
	testutil.AssertFalse(t, testutil.FileExists(t, hookFile), "No sidebar-hooks.yaml should be written in dry-run mode")
}
