package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/rippreport/go-content-pipeline/internal/cache"
	"github.com/rippreport/go-content-pipeline/internal/client"
	"github.com/rippreport/go-content-pipeline/internal/config"
	"github.com/rippreport/go-content-pipeline/internal/errors"
	"github.com/rippreport/go-content-pipeline/internal/orchestrator"
	"github.com/rippreport/go-content-pipeline/internal/processor"
	"github.com/rippreport/go-content-pipeline/internal/refiner"
	"github.com/rippreport/go-content-pipeline/internal/server"
	"github.com/rippreport/go-content-pipeline/internal/transcript"
)

// App manages the CLI subcommands and orchestrator lifecycle.
//
// Requirements: 16.1 through 16.8, 10.3, 10.4
type App struct {
	stdout io.Writer
	stderr io.Writer
}

// NewApp creates a new CLI App.
func NewApp(stdout, stderr io.Writer) *App {
	return &App{
		stdout: stdout,
		stderr: stderr,
	}
}

// Run parses arguments and executes the requested subcommand with signal handling.
func (a *App) Run(args []string) int {
	if len(args) < 2 {
		a.printUsage()
		return 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Signal handling for graceful shutdown (SIGINT, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Fprintln(a.stderr, "\nReceived shutdown signal, terminating gracefully...")
		cancel()
	}()

	cmd := args[1]
	cmdArgs := args[2:]

	switch cmd {
	case "help", "-h", "--help":
		a.printUsage()
		return 0
	case "generate-related":
		return a.runGenerateRelated(ctx, cmdArgs)
	case "generate-hooks":
		return a.runGenerateHooks(ctx, cmdArgs)
	case "refine-hooks":
		return a.runRefineHooks(ctx, cmdArgs)
	case "fetch-transcript":
		return a.runFetchTranscript(ctx, cmdArgs)
	case "summarize-transcript":
		return a.runSummarizeTranscript(ctx, cmdArgs)
	case "pipeline":
		return a.runPipeline(ctx, cmdArgs)
	case "apply-staged":
		return a.runApplyStaged(ctx, cmdArgs)
	case "clear-staged":
		return a.runClearStaged(ctx, cmdArgs)
	default:
		fmt.Fprintf(a.stderr, "Unknown command: %s\n\n", cmd)
		a.printUsage()
		return 1
	}
}

func (a *App) printUsage() {
	fmt.Fprintln(a.stdout, "Go Content Pipeline - High Performance Hugo Recommendation System")
	fmt.Fprintln(a.stdout, "\nUsage: pipeline <command> [options]")
	fmt.Fprintln(a.stdout, "\nCommands:")
	fmt.Fprintln(a.stdout, "  generate-related     Generate related article recommendations")
	fmt.Fprintln(a.stdout, "  generate-hooks       Generate contextual sidebar hooks")
	fmt.Fprintln(a.stdout, "  refine-hooks         Refine existing hooks using Gemini API")
	fmt.Fprintln(a.stdout, "  fetch-transcript     Download YouTube transcript for podcast posts")
	fmt.Fprintln(a.stdout, "  summarize-transcript Generate factual summary from transcript")
	fmt.Fprintln(a.stdout, "  pipeline             Run complete pipeline (related + hooks)")
	fmt.Fprintln(a.stdout, "  apply-staged         Apply staged rankings or hooks to posts (use --hooks flag)")
	fmt.Fprintln(a.stdout, "  clear-staged      Clear staged ranking or hook files (use --hooks flag)")
	fmt.Fprintln(a.stdout, "\nRun 'pipeline <command> --help' for options on a specific command.")
}

func (a *App) buildOrchestrator(cfgPath string, managedOverride ...bool) (*orchestrator.Orchestrator, error) {
	var cfg *config.Config
	var err error

	if cfgPath != "" {
		if _, statErr := os.Stat(cfgPath); statErr == nil {
			cfg, err = config.LoadConfig(cfgPath)
			if err != nil {
				return nil, fmt.Errorf("failed loading config %s: %w", cfgPath, err)
			}
		}
	}

	if cfg == nil {
		cfg = config.DefaultConfig()
		cfg.Processing.PostsDirectory = "/Volumes/1tb/rippreport/content/p"
	}

	if len(managedOverride) > 0 && managedOverride[0] {
		m := true
		cfg.LlamaServer.Managed = &m
	}

	llamaClient := client.NewClient(cfg.LlamaServer)
	postManager := processor.NewPostManager(cfg.Processing.PostsDirectory)
	embeddingCache := cache.NewEmbeddingCache(cfg.Storage.CacheFile, cfg.LlamaServer.EmbeddingModel)
	_ = embeddingCache.Load()
	progressTracker := processor.NewProgressTracker(cfg.Storage.ProgressFile)
	_ = progressTracker.Load()
	backupManager := processor.NewBackupManager(cfg.Storage.BackupDir)
	stagingManager := processor.NewStagingManager(cfg.Processing.PostsDirectory, postManager, backupManager)
	compressionCache := processor.NewCompressionCache(cfg.Storage.CompressionCache, cfg.LlamaServer.GetHookModel())
	_ = compressionCache.Load()

	return orchestrator.NewOrchestrator(
		cfg,
		llamaClient,
		postManager,
		embeddingCache,
		progressTracker,
		backupManager,
		stagingManager,
		compressionCache,
	), nil
}

func parseSlugs(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	var cleaned []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}

func (a *App) runGenerateRelated(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("generate-related", flag.ContinueOnError)
	fs.SetOutput(a.stdout)

	configPath := fs.String("config", "config.yaml", "Path to configuration file")
	resume := fs.Bool("resume", false, "Skip already-processed posts")
	posts := fs.String("posts", "", "Comma-separated list of specific post slugs")
	stage := fs.Bool("stage", false, "Save rankings to staging files instead of updating posts")
	refreshRecent := fs.Bool("refresh-recent", false, "Merge rankings preserving older articles")
	vectorOnly := fs.Bool("vector-only", false, "Skip LLM scoring and use vector similarity only")
	workers := fs.Int("workers", 4, "Number of concurrent workers")
	dryRun := fs.Bool("dry-run", false, "Simulate operations without modifying files")
	scoringModel := fs.String("scoring-model", "", "Model to use for scoring (overrides config)")
	useReranker := fs.Bool("use-reranker", false, "Use cross-encoder reranker model (/v1/rerank) for scoring")
	managed := fs.Bool("managed", false, "Auto-manage llama-server process lifecycle")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	orch, err := a.buildOrchestrator(*configPath, *managed)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error initializing pipeline: %v\n", err)
		return 1
	}
	defer orch.Close()

	res, err := orch.GenerateRelatedArticles(ctx, orchestrator.RankingOptions{
		Resume:        *resume,
		Stage:         *stage,
		RefreshRecent: *refreshRecent,
		VectorOnly:    *vectorOnly,
		SpecificSlugs: parseSlugs(*posts),
		Workers:       *workers,
		DryRun:        *dryRun,
		Model:         *scoringModel,
		UseReranker:   *useReranker,
	})
	if err != nil {
		fmt.Fprintf(a.stderr, "Error generating related articles: %v\n", err)
		return 1
	}

	// Print error summary if there were any issues
	if res.FailedPosts > 0 || orch.ErrorCollector().Count() > 0 {
		summaryText := orch.ErrorCollector().FormatSummary(errors.PipelineSummary{
			TotalDuration:      res.Duration,
			TotalPosts:         res.ProcessedPosts + res.FailedPosts + res.SkippedPosts,
			ProcessedPosts:     res.ProcessedPosts,
			FailedPosts:        res.FailedPosts,
			SkippedPosts:       res.SkippedPosts,
			EmbeddingsComputed: res.CacheHits, // Approximate
			CacheHits:          res.CacheHits,
			TotalAPICalls:      0, // Not tracked in RankingResult
		})
		fmt.Fprintln(a.stdout, summaryText)
	}

	fmt.Fprintf(a.stdout, "\nCompleted: %d processed, %d failed, %d skipped in %v\n",
		res.ProcessedPosts, res.FailedPosts, res.SkippedPosts, res.Duration)
	return 0
}

func (a *App) runGenerateHooks(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("generate-hooks", flag.ContinueOnError)
	fs.SetOutput(a.stdout)

	configPath := fs.String("config", "config.yaml", "Path to configuration file")
	posts := fs.String("posts", "", "Comma-separated list of specific post slugs")
	useCompressed := fs.Bool("use-compressed", false, "Use compressed article summaries for hooks")
	batchSize := fs.Int("batch-size", 5, "Batch size for LLM hook generation")
	fromStaging := fs.Bool("from-staging", false, "Read rankings from staging files")
	stage := fs.Bool("stage", false, "Save hooks to staging files instead of writing directly")
	dryRun := fs.Bool("dry-run", false, "Simulate operations without writing files")
	hookModel := fs.String("hook-model", "", "Model to use for hook generation (overrides config)")
	managed := fs.Bool("managed", false, "Auto-manage llama-server process lifecycle")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	orch, err := a.buildOrchestrator(*configPath, *managed)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error initializing pipeline: %v\n", err)
		return 1
	}
	defer orch.Close()

	res, err := orch.GenerateHooks(ctx, orchestrator.HookOptions{
		SpecificSlugs: parseSlugs(*posts),
		UseCompressed: *useCompressed,
		BatchSize:     *batchSize,
		FromStaging:   *fromStaging,
		Stage:         *stage,
		DryRun:        *dryRun,
		Model:         *hookModel,
	})
	if err != nil {
		fmt.Fprintf(a.stderr, "Error generating hooks: %v\n", err)
		return 1
	}

	fmt.Fprintf(a.stdout, "\nGenerated %d hooks (%d failed) in %v\n",
		res.TotalHooks, res.FailedHooks, res.Duration)
	return 0
}

func (a *App) runRefineHooks(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("refine-hooks", flag.ContinueOnError)
	fs.SetOutput(a.stdout)

	configPath := fs.String("config", "config.yaml", "Path to configuration file")
	posts := fs.String("posts", "", "Comma-separated list of specific post slugs")
	force := fs.Bool("force", false, "Refine even if hooks are already refined")
	dryRun := fs.Bool("dry-run", false, "Simulate operations without modifying files")
	geminiModel := fs.String("gemini-model", "", "Gemini model to use (e.g., gemini-3.7-flash)")
	batch := fs.Int("batch", 0, "Batch size for cross-post refinement (0 = per-post)")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		fmt.Fprintf(a.stderr, "Error: GEMINI_API_KEY environment variable is not set\n")
		return 1
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error loading config: %v\n", err)
		return 1
	}

	refinerCfg := refiner.Config{
		GeminiAPIKey: geminiAPIKey,
		PostsDir:     cfg.Processing.PostsDirectory,
	}

	r, err := refiner.New(refinerCfg)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error creating refiner: %v\n", err)
		return 1
	}

	res, err := r.RefineHooks(ctx, refiner.RefineOptions{
		SpecificSlugs: parseSlugs(*posts),
		Force:         *force,
		DryRun:        *dryRun,
		Model:         *geminiModel,
		BatchSize:     *batch,
	})
	if err != nil {
		fmt.Fprintf(a.stderr, "Error refining hooks: %v\n", err)
		return 1
	}

	fmt.Fprintf(a.stdout, "\nRefined: %d processed, %d failed, %d skipped in %v\n",
		res.ProcessedPosts, res.FailedPosts, res.SkippedPosts, res.Duration)
	return 0
}

func (a *App) runFetchTranscript(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("fetch-transcript", flag.ContinueOnError)
	fs.SetOutput(a.stdout)

	configPath := fs.String("config", "config.yaml", "Path to configuration file")
	posts := fs.String("posts", "", "Comma-separated list of post slugs")
	force := fs.Bool("force", false, "Overwrite existing transcript.md files")
	autoOnly := fs.Bool("auto-only", false, "Only fetch auto-generated transcripts")
	language := fs.String("language", "en", "Transcript language code")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *posts == "" {
		fmt.Fprintf(a.stderr, "Error: --posts flag is required\n")
		fmt.Fprintf(a.stderr, "Usage: pipeline fetch-transcript --posts october-surprise\n")
		return 1
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error loading config: %v\n", err)
		return 1
	}

	fetcher := transcript.NewFetcher(cfg.Processing.PostsDirectory)
	slugs := parseSlugs(*posts)

	if len(slugs) == 0 {
		fmt.Fprintf(a.stderr, "Error: no post slugs provided\n")
		return 1
	}

	fmt.Fprintf(a.stdout, "Fetching transcripts for %d post(s)...\n\n", len(slugs))

	successCount := 0
	failCount := 0

	for i, slug := range slugs {
		fmt.Fprintf(a.stdout, "[%d/%d] Processing %s...\n", i+1, len(slugs), slug)

		opts := transcript.FetchOptions{
			PostSlug: slug,
			Force:    *force,
			Language: *language,
			AutoOnly: *autoOnly,
		}

		result, err := fetcher.FetchTranscript(opts)
		if err != nil {
			fmt.Fprintf(a.stderr, "  ✗ Failed: %v\n\n", err)
			failCount++
			continue
		}

		if result.Success {
			fmt.Fprintf(a.stdout, "  ✓ Success! Video ID: %s\n", result.VideoID)
			fmt.Fprintf(a.stdout, "  ✓ Saved to: %s\n", result.TranscriptPath)
			fmt.Fprintf(a.stdout, "  ✓ Word count: %d\n\n", result.WordCount)
			successCount++
		} else {
			fmt.Fprintf(a.stderr, "  ✗ Failed: %v\n\n", result.Error)
			failCount++
		}
	}

	fmt.Fprintf(a.stdout, "\n")
	fmt.Fprintf(a.stdout, "Summary: %d succeeded, %d failed\n", successCount, failCount)

	if successCount > 0 {
		fmt.Fprintf(a.stdout, "\nNext steps:\n")
		fmt.Fprintf(a.stdout, "  1. Review the transcript(s) for accuracy\n")
		fmt.Fprintf(a.stdout, "  2. Generate embeddings: ./bin/pipeline generate-related --posts %s\n", *posts)
		fmt.Fprintf(a.stdout, "  3. Generate hooks: ./bin/pipeline generate-hooks --posts %s\n", *posts)
	}

	if failCount > 0 {
		return 1
	}
	return 0
}

func (a *App) runSummarizeTranscript(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("summarize-transcript", flag.ContinueOnError)
	fs.SetOutput(a.stdout)

	configPath := fs.String("config", "config.yaml", "Path to configuration file")
	posts := fs.String("posts", "", "Comma-separated list of post slugs")
	force := fs.Bool("force", false, "Overwrite existing summary.md files")
	model := fs.String("model", "", "LLM model to use for summarization")
	managed := fs.Bool("managed", false, "Auto-manage llama-server process lifecycle")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *posts == "" {
		fmt.Fprintf(a.stderr, "Error: --posts flag is required\n")
		fmt.Fprintf(a.stderr, "Usage: pipeline summarize-transcript --posts midterm-election-special\n")
		return 1
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error loading config: %v\n", err)
		return 1
	}

	// Determine which model to use
	summaryModel := *model
	if summaryModel == "" {
		summaryModel = cfg.LlamaServer.LLMModel
	}

	// Build orchestrator with managed mode if requested
	orch, err := a.buildOrchestrator(*configPath, *managed)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error initializing pipeline: %v\n", err)
		return 1
	}
	defer orch.Close()

	// If managed mode, the server manager will handle model loading
	if *managed || cfg.LlamaServer.IsManaged() {
		sm := orch.GetServerManager()
		if sm != nil {
			fmt.Fprintf(a.stdout, "Starting managed llama-server with model %s...\n", summaryModel)
			if err := sm.EnsureModel(ctx, summaryModel, server.ModeCompletion); err != nil {
				fmt.Fprintf(a.stderr, "Failed to start llama-server: %v\n", err)
				return 1
			}
			defer sm.Stop()
			fmt.Fprintf(a.stdout, "✓ llama-server started\n\n")
		}
	}

	summarizer := transcript.NewSummarizer(transcript.SummarizerConfig{
		LlamaServerURL: cfg.LlamaServer.BaseURL,
		Model:          summaryModel,
	})

	slugs := parseSlugs(*posts)
	if len(slugs) == 0 {
		fmt.Fprintf(a.stderr, "Error: no post slugs provided\n")
		return 1
	}

	fmt.Fprintf(a.stdout, "Generating summaries for %d post(s)...\n\n", len(slugs))

	successCount := 0
	failCount := 0

	for i, slug := range slugs {
		fmt.Fprintf(a.stdout, "[%d/%d] Processing %s...\n", i+1, len(slugs), slug)

		postDir := filepath.Join(cfg.Processing.PostsDirectory, slug)
		transcriptPath := filepath.Join(postDir, "transcript.md")

		opts := transcript.SummarizeOptions{
			PostSlug:       slug,
			TranscriptPath: transcriptPath,
			Force:          *force,
		}

		result, err := summarizer.Summarize(ctx, opts)
		if err != nil {
			fmt.Fprintf(a.stderr, "  ✗ Failed: %v\n\n", err)
			failCount++
			continue
		}

		if result.Success {
			fmt.Fprintf(a.stdout, "  ✓ Success!\n")
			fmt.Fprintf(a.stdout, "  ✓ Saved to: %s\n", result.SummaryPath)
			fmt.Fprintf(a.stdout, "  ✓ Estimated word count: %d\n\n", result.WordCount)
			successCount++
		} else {
			fmt.Fprintf(a.stderr, "  ✗ Failed: %v\n\n", result.Error)
			failCount++
		}
	}

	fmt.Fprintf(a.stdout, "\n")
	fmt.Fprintf(a.stdout, "Summary: %d succeeded, %d failed\n", successCount, failCount)

	if successCount > 0 {
		fmt.Fprintf(a.stdout, "\nNext steps:\n")
		fmt.Fprintf(a.stdout, "  1. Review the summaries for accuracy\n")
		fmt.Fprintf(a.stdout, "  2. Generate embeddings: ./bin/pipeline generate-related --posts %s\n", *posts)
		fmt.Fprintf(a.stdout, "     (Will automatically use summary.md for embeddings)\n")
	}

	if failCount > 0 {
		return 1
	}
	return 0
}

func (a *App) runPipeline(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("pipeline", flag.ContinueOnError)
	fs.SetOutput(a.stdout)

	configPath := fs.String("config", "config.yaml", "Path to configuration file")
	resume := fs.Bool("resume", false, "Skip already-processed posts")
	posts := fs.String("posts", "", "Comma-separated list of specific post slugs")
	stage := fs.Bool("stage", false, "Stage rankings instead of writing directly to posts")
	refreshRecent := fs.Bool("refresh-recent", false, "Refresh recent mode")
	vectorOnly := fs.Bool("vector-only", false, "Vector only scoring")
	useCompressed := fs.Bool("use-compressed", false, "Use compressed summaries for hooks")
	batchSize := fs.Int("batch-size", 5, "Batch size for hooks")
	dryRun := fs.Bool("dry-run", false, "Dry run mode")
	scoringModel := fs.String("scoring-model", "", "Model to use for scoring (overrides config)")
	useReranker := fs.Bool("use-reranker", false, "Use cross-encoder reranker model (/v1/rerank) for scoring")
	hookModel := fs.String("hook-model", "", "Model to use for hook generation (overrides config)")
	managed := fs.Bool("managed", false, "Auto-manage llama-server process lifecycle")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	orch, err := a.buildOrchestrator(*configPath, *managed)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error initializing pipeline: %v\n", err)
		return 1
	}
	defer orch.Close()

	slugList := parseSlugs(*posts)
	opts := orchestrator.PipelineOptions{
		RankingOptions: orchestrator.RankingOptions{
			Resume:        *resume,
			Stage:         *stage,
			RefreshRecent: *refreshRecent,
			VectorOnly:    *vectorOnly,
			SpecificSlugs: slugList,
			DryRun:        *dryRun,
			Model:         *scoringModel,
			UseReranker:   *useReranker,
		},
		HookOptions: orchestrator.HookOptions{
			SpecificSlugs: slugList,
			UseCompressed: *useCompressed,
			BatchSize:     *batchSize,
			DryRun:        *dryRun,
			Model:         *hookModel,
		},
	}

	res, err := orch.RunPipeline(ctx, opts)
	if err != nil {
		fmt.Fprintf(a.stderr, "Pipeline execution error: %v\n", err)
		return 1
	}

	summaryText := orch.ErrorCollector().FormatSummary(res.Summary)
	fmt.Fprintln(a.stdout, summaryText)
	return 0
}

func (a *App) runApplyStaged(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("apply-staged", flag.ContinueOnError)
	fs.SetOutput(a.stdout)

	configPath := fs.String("config", "config.yaml", "Path to configuration file")
	posts := fs.String("posts", "", "Comma-separated list of post slugs to apply")
	hooks := fs.Bool("hooks", false, "Apply staged hooks instead of rankings")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	orch, err := a.buildOrchestrator(*configPath)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error initializing pipeline: %v\n", err)
		return 1
	}

	slugs := parseSlugs(*posts)

	if *hooks {
		if err := orch.ApplyStagedHooks(ctx, slugs); err != nil {
			fmt.Fprintf(a.stderr, "Error applying staged hooks: %v\n", err)
			return 1
		}
		fmt.Fprintln(a.stdout, "Successfully applied staged hooks.")
	} else {
		if err := orch.ApplyStaged(ctx, slugs); err != nil {
			fmt.Fprintf(a.stderr, "Error applying staged rankings: %v\n", err)
			return 1
		}
		fmt.Fprintln(a.stdout, "Successfully applied staged rankings.")
	}

	return 0
}

func (a *App) runClearStaged(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("clear-staged", flag.ContinueOnError)
	fs.SetOutput(a.stdout)

	configPath := fs.String("config", "config.yaml", "Path to configuration file")
	posts := fs.String("posts", "", "Comma-separated list of post slugs to clear")
	hooks := fs.Bool("hooks", false, "Clear staged hooks instead of rankings")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	orch, err := a.buildOrchestrator(*configPath)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error initializing pipeline: %v\n", err)
		return 1
	}

	slugs := parseSlugs(*posts)

	if *hooks {
		if err := orch.ClearStagedHooks(ctx, slugs); err != nil {
			fmt.Fprintf(a.stderr, "Error clearing staged hooks: %v\n", err)
			return 1
		}
		fmt.Fprintln(a.stdout, "Successfully cleared staged hooks files.")
	} else {
		if err := orch.ClearStaged(ctx, slugs); err != nil {
			fmt.Fprintf(a.stderr, "Error clearing staged rankings: %v\n", err)
			return 1
		}
		fmt.Fprintln(a.stdout, "Successfully cleared staged files.")
	}

	return 0
}
