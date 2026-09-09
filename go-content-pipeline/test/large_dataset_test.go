package test

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/orchestrator"
	"github.com/rippreport/go-content-pipeline/internal/testutil"
)

// Requirement 23.3: Test with large dataset
// - Run pipeline on large dataset (100+ posts)
// - Verify memory stays under 500MB
// - Verify total time is acceptable
// - Profile for bottlenecks
func TestLargeDataset_100Posts_PerformanceAndMemory(t *testing.T) {
	env := SetupTestEnv(t)

	// Create 100 posts
	postCount := 100
	t.Logf("Generating %d test posts...", postCount)
	posts := env.CreateSamplePosts(postCount)
	testutil.AssertEqual(t, len(posts), postCount, "Created 100 posts")

	// Trigger GC before starting measurement
	runtime.GC()
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	startTime := time.Now()
	ctx := context.Background()

	res, err := env.Orchestrator.RunPipeline(ctx, orchestrator.PipelineOptions{
		RankingOptions: orchestrator.RankingOptions{
			VectorOnly: true, // Focus on large-scale candidate ranking & front matter I/O
		},
		HookOptions: orchestrator.HookOptions{},
	})
	elapsed := time.Since(startTime)

	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	testutil.AssertNoError(t, err, "RunPipeline on 100 posts should succeed")
	testutil.AssertEqual(t, res.Ranking.ProcessedPosts, postCount, "All 100 posts ranked")
	testutil.AssertEqual(t, res.Ranking.FailedPosts, 0, "Zero failures")

	// Memory analysis
	heapAllocMB := float64(memAfter.Alloc) / (1024 * 1024)
	totalAllocMB := float64(memAfter.TotalAlloc-memBefore.TotalAlloc) / (1024 * 1024)

	t.Logf("=== Large Dataset Performance Report ===")
	t.Logf("Posts Processed:   %d", postCount)
	t.Logf("Total Time:        %v (%.2f ms/post)", elapsed, float64(elapsed.Milliseconds())/float64(postCount))
	t.Logf("Live Heap Memory:  %.2f MB (Limit: 500 MB)", heapAllocMB)
	t.Logf("Total Cumulative:  %.2f MB", totalAllocMB)
	t.Logf("GC Runs:           %d", memAfter.NumGC-memBefore.NumGC)
	t.Logf("========================================")

	// Target constraints from Requirements:
	// 1. Memory stays strictly under 500MB
	testutil.AssertTrue(t, heapAllocMB < 500.0, "Live heap allocation must stay under 500MB")

	// 2. Performance target: 100 posts well under target 15 minutes (under 60s in mock)
	testutil.AssertTrue(t, elapsed < 15*time.Minute, "Execution time must be under 15 minutes")
}
