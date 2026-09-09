package test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rippreport/go-content-pipeline/internal/orchestrator"
	"github.com/rippreport/go-content-pipeline/internal/processor"
	"github.com/rippreport/go-content-pipeline/internal/testutil"
)

// TestSeparateModels_ScoringAndHooks verifies that different models can be configured
// and used for relevance scoring vs contextual hook generation.
func TestSeparateModels_ScoringAndHooks(t *testing.T) {
	env := SetupTestEnv(t)

	// Configure distinct models
	scoringModel := "specialized-scorer-v1"
	hookModel := "creative-hook-writer-v2"

	env.Config.LlamaServer.ScoringModel = scoringModel
	env.Config.LlamaServer.HookModel = hookModel

	// Create 2 mutual test posts
	p1 := testutil.DefaultPost("model-post-1")
	p1.Content = "Article one discussing municipal planning and coastal zoning."
	env.CreatePost(p1)

	p2 := testutil.DefaultPost("model-post-2")
	p2.Content = "Article two covering environmental assessments of municipal shoreline projects."
	env.CreatePost(p2)

	ctx := context.Background()
	res, err := env.Orchestrator.RunPipeline(ctx, orchestrator.PipelineOptions{
		RankingOptions: orchestrator.RankingOptions{
			VectorOnly: false, // Exercises LLM scoring
		},
		HookOptions: orchestrator.HookOptions{},
	})
	testutil.AssertNoError(t, err, "RunPipeline with distinct models should succeed")
	testutil.AssertEqual(t, res.Ranking.ProcessedPosts, 2, "Processed posts in ranking")
	testutil.AssertTrue(t, res.Hooks.TotalHooks >= 2, "Hooks generated")

	// Inspect all completion requests sent to mock server
	scoringReqCount := 0
	hookReqCount := 0

	for _, req := range env.MockServer.CompletionRequests {
		if strings.Contains(req.Prompt, "topical relevance") {
			scoringReqCount++
			testutil.AssertEqual(t, req.Model, scoringModel, "Scoring completion request should use scoring_model")
		} else if strings.Contains(req.Prompt, "contextual sidebar hook") || strings.Contains(req.Prompt, "narrative teaser") || strings.Contains(req.Prompt, "BRIDGE BRIEF") || strings.Contains(req.Prompt, "THE CONFLICT") {
			hookReqCount++
			testutil.AssertEqual(t, req.Model, hookModel, "Hook generation completion request should use hook_model")
		}
	}

	testutil.AssertTrue(t, scoringReqCount > 0, "Expected at least one scoring completion request")
	testutil.AssertTrue(t, hookReqCount > 0, "Expected at least one hook completion request")

	// Verify sidebar-hooks.yaml has ModelName set to hookModel
	hooks, err := processor.LoadHooks(filepath.Join(env.PostsDir, "model-post-1"))
	testutil.AssertNoError(t, err, "LoadHooks")
	testutil.AssertTrue(t, len(hooks) > 0, "Loaded hooks should not be empty")
	for _, h := range hooks {
		testutil.AssertEqual(t, h.ModelName, hookModel, "Saved hook should record hook_model name")
	}
}

// TestSeparateModels_OptionsOverride verifies that runtime options override the configured models.
func TestSeparateModels_OptionsOverride(t *testing.T) {
	env := SetupTestEnv(t)

	env.Config.LlamaServer.ScoringModel = "default-scorer"
	env.Config.LlamaServer.HookModel = "default-hooker"

	runtimeScoringModel := "runtime-scorer-override"
	runtimeHookModel := "runtime-hook-override"

	p1 := testutil.DefaultPost("override-post-1")
	env.CreatePost(p1)
	p2 := testutil.DefaultPost("override-post-2")
	env.CreatePost(p2)

	ctx := context.Background()
	_, err := env.Orchestrator.RunPipeline(ctx, orchestrator.PipelineOptions{
		RankingOptions: orchestrator.RankingOptions{
			VectorOnly: false,
			Model:      runtimeScoringModel,
		},
		HookOptions: orchestrator.HookOptions{
			Model: runtimeHookModel,
		},
	})
	testutil.AssertNoError(t, err, "RunPipeline with runtime model overrides should succeed")

	foundScorer := false
	foundHooker := false
	for _, req := range env.MockServer.CompletionRequests {
		if req.Model == runtimeScoringModel {
			foundScorer = true
		}
		if req.Model == runtimeHookModel {
			foundHooker = true
		}
	}

	testutil.AssertTrue(t, foundScorer, "Should find completion request with runtime-scorer-override")
	testutil.AssertTrue(t, foundHooker, "Should find completion request with runtime-hook-override")
}
