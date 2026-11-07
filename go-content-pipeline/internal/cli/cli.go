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
	"time"

	"github.com/rippreport/go-content-pipeline/internal/cache"
	"github.com/rippreport/go-content-pipeline/internal/client"
	"github.com/rippreport/go-content-pipeline/internal/config"
	"github.com/rippreport/go-content-pipeline/internal/email"
	"github.com/rippreport/go-content-pipeline/internal/errors"
	"github.com/rippreport/go-content-pipeline/internal/models"
	"github.com/rippreport/go-content-pipeline/internal/newsletter"
	"github.com/rippreport/go-content-pipeline/internal/orchestrator"
	"github.com/rippreport/go-content-pipeline/internal/processor"
	"github.com/rippreport/go-content-pipeline/internal/refiner"
	"github.com/rippreport/go-content-pipeline/internal/server"
	"github.com/rippreport/go-content-pipeline/internal/transcript"
	"github.com/rippreport/go-content-pipeline/internal/youtube"
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
	case "check-links", "link-checker":
		return a.runCheckLinks(ctx, cmdArgs)
	case "analyze-broken-links":
		return a.runAnalyzeBrokenLinks(ctx, cmdArgs)
	case "archive-links":
		return a.runLinkRepairs(ctx, cmdArgs, false)
	case "update-redirected-links":
		return a.runLinkRepairs(ctx, cmdArgs, true)
	case "audit-content", "content-audit":
		return a.runAuditContent(ctx, cmdArgs)
	case "verify-archives", "archive-link-verifier":
		return a.runVerifyArchives(ctx, cmdArgs)
	case "missing-pdfs":
		return a.runMissingPDFs(cmdArgs)
	case "help", "-h", "--help":
		a.printUsage()
		return 0
	case "generate-related":
		return a.runGenerateRelated(ctx, cmdArgs)
	case "generate-hooks":
		return a.runGenerateHooks(ctx, cmdArgs)
	case "generate-frontpage", "frontpage-hooks":
		return a.runGenerateFrontpage(ctx, cmdArgs)
	case "scan-hook-articles", "scan-hook-authors":
		return a.runHookScan(ctx, cmdArgs, cmd)
	case "scan-hook-phrases":
		return a.runScanHookPhrases(ctx, cmdArgs)
	case "refine-hooks", "refine":
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
	case "send-newsletter":
		return a.runSendNewsletter(ctx, cmdArgs)
	case "social-post":
		return a.runSocialPost(ctx, cmdArgs)
	case "new-post":
		return a.runNewPost(ctx, cmdArgs)
	case "fetch-email", "email-post":
		return a.runNewPost(ctx, append([]string{"--email"}, cmdArgs...))
	case "generate-alt", "alt-post":
		return a.runGenerateAlt(ctx, cmdArgs)
	case "repost":
		return a.runRepost(ctx, cmdArgs)
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
	fmt.Fprintln(a.stdout, "  check-links          Audit links in rendered Hugo posts with resumable JSON results")
	fmt.Fprintln(a.stdout, "  analyze-broken-links Group broken links by domain, status, and error")
	fmt.Fprintln(a.stdout, "  archive-links        Find and install verified Wayback replacements (JSON + CSV)")
	fmt.Fprintln(a.stdout, "  update-redirected-links Update successful redirected destinations (previews by default)")
	fmt.Fprintln(a.stdout, "  audit-content        Audit WordPress shortcodes using a local model or offline checks")
	fmt.Fprintln(a.stdout, "  verify-archives      Verify and repair audit links with archives and fallbacks")
	fmt.Fprintln(a.stdout, "  missing-pdfs         Export unresolved PDF links from an audit report to JSON")
	fmt.Fprintln(a.stdout, "  generate-related     Generate related article recommendations")
	fmt.Fprintln(a.stdout, "  generate-hooks       Generate contextual sidebar hooks")
	fmt.Fprintln(a.stdout, "  generate-frontpage   Generate frontpage narrative bridges connecting recent posts")
	fmt.Fprintln(a.stdout, "  scan-hook-articles   Export article-number references in hook text to JSON")
	fmt.Fprintln(a.stdout, "  scan-hook-authors    Export Ripp names and author references in hook text to JSON")
	fmt.Fprintln(a.stdout, "  scan-hook-phrases    Export repetitive contrast phrasing in hooks to JSON")
	fmt.Fprintln(a.stdout, "  refine-hooks         Refine existing hooks using Gemini API")
	fmt.Fprintln(a.stdout, "  fetch-transcript     Download YouTube transcript for podcast posts")
	fmt.Fprintln(a.stdout, "  summarize-transcript Generate factual summary from transcript")
	fmt.Fprintln(a.stdout, "  pipeline             Run complete pipeline (related + hooks)")
	fmt.Fprintln(a.stdout, "  apply-staged         Apply staged rankings or hooks to posts (use --hooks flag)")
	fmt.Fprintln(a.stdout, "  clear-staged         Clear staged rankings, hooks (--hooks), or both (--all)")
	fmt.Fprintln(a.stdout, "  send-newsletter      Send newsletter via Brevo for a specific post")
	fmt.Fprintln(a.stdout, "  social-post          Print post summary and URL for social media")
	fmt.Fprintln(a.stdout, "  new-post             Create a new post with folder structure and front matter (use --email for Apple Mail)")
	fmt.Fprintln(a.stdout, "  email-post           Create a new post by retrieving the latest email from Apple Mail")
	fmt.Fprintln(a.stdout, "  generate-alt         Generate alt tag for thumbnail using local vision LLM and update front matter")
	fmt.Fprintln(a.stdout, "  repost               Duplicate a post folder as a repost with today's date")
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
		if err == flag.ErrHelp {
			return 0
		}
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
	resume := fs.Bool("resume", false, "Skip posts with complete hooks in the selected output location")
	dryRun := fs.Bool("dry-run", false, "Simulate operations without writing files")
	cloudflare := fs.Bool("cloudflare", false, "Use Cloudflare Workers AI alongside other hook workers (requires CLOUDFLARE_API_TOKEN)")
	cloudflareAccountID := fs.String("cloudflare-account-id", os.Getenv("CLOUDFLARE_ACCOUNT_ID"), "Cloudflare account ID (defaults to CLOUDFLARE_ACCOUNT_ID)")
	cloudflareModel := fs.String("cloudflare-model", client.DefaultCloudflareModel, "Cloudflare Workers AI model ID")
	cloudflareMaxTokens := fs.Int("cloudflare-max-tokens", client.DefaultCloudflareMaxTokens, "Output token limit for Cloudflare requests")
	gemini := fs.Bool("gemini", false, "Use Gemini alongside other hook workers (requires GEMINI_API_KEY)")
	geminiModel := fs.String("gemini-model", client.DefaultGeminiModel, "Gemini API model ID")
	geminiMaxTokens := fs.Int("gemini-max-tokens", client.DefaultGeminiMaxTokens, "Output token limit for Gemini requests")
	plugsky := fs.Bool("plugsky", false, "Use Plugsky alongside other hook workers (requires PLUGSKY_API_KEY)")
	plugskyModel := fs.String("plugsky-model", client.DefaultPlugskyModel, "Plugsky API model ID")
	plugskyMaxTokens := fs.Int("plugsky-max-tokens", client.DefaultPlugskyMaxTokens, "Output token limit for Plugsky requests")
	openRouter := fs.Bool("openrouter", false, "Use OpenRouter alongside llama-server for hook generation (requires OPENROUTER_API_KEY)")
	openRouterModel := fs.String("openrouter-model", client.DefaultOpenRouterModel, "OpenRouter helper model")
	openRouterMaxTokens := fs.Int("openrouter-max-tokens", client.DefaultOpenRouterMaxTokens, "Output token limit for OpenRouter requests (independent of llama-server)")
	hookModel := fs.String("hook-model", "", "Model to use for hook generation (overrides config)")
	managed := fs.Bool("managed", false, "Auto-manage llama-server process lifecycle")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	orch, err := a.buildOrchestrator(*configPath, *managed)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error initializing pipeline: %v\n", err)
		return 1
	}
	defer orch.Close()

	res, err := orch.GenerateHooks(ctx, orchestrator.HookOptions{
		SpecificSlugs:       parseSlugs(*posts),
		UseCompressed:       *useCompressed,
		BatchSize:           *batchSize,
		FromStaging:         *fromStaging,
		Stage:               *stage,
		Resume:              *resume,
		DryRun:              *dryRun,
		Model:               *hookModel,
		Cloudflare:          *cloudflare,
		CloudflareAccountID: *cloudflareAccountID,
		CloudflareModel:     *cloudflareModel,
		CloudflareMaxTokens: *cloudflareMaxTokens,
		Gemini:              *gemini,
		GeminiModel:         *geminiModel,
		GeminiMaxTokens:     *geminiMaxTokens,
		Plugsky:             *plugsky,
		PlugskyModel:        *plugskyModel,
		PlugskyMaxTokens:    *plugskyMaxTokens,
		OpenRouter:          *openRouter,
		OpenRouterModel:     *openRouterModel,
		OpenRouterMaxTokens: *openRouterMaxTokens,
	})
	if err != nil {
		fmt.Fprintf(a.stderr, "Error generating hooks: %v\n", err)
		return 1
	}

	fmt.Fprintf(a.stdout, "\nGenerated %d hooks (%d failed, %d posts skipped) in %v\n",
		res.TotalHooks, res.FailedHooks, res.SkippedPosts, res.Duration)
	return 0
}

func (a *App) runGenerateFrontpage(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("generate-frontpage", flag.ContinueOnError)
	fs.SetOutput(a.stdout)

	configPath := fs.String("config", "config.yaml", "Path to configuration file")
	limit := fs.Int("limit", 20, "Number of recent posts to include in the frontpage narrative chain")
	outputPath := fs.String("output", "data/frontpage-hooks.yaml", "Path to output YAML file")
	useCompressed := fs.Bool("use-compressed", false, "Use compressed article summaries for bridges")
	dryRun := fs.Bool("dry-run", false, "Simulate operations without writing files")
	model := fs.String("model", "", "Model to use for bridge generation (overrides config)")
	managed := fs.Bool("managed", false, "Auto-manage llama-server process lifecycle")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	orch, err := a.buildOrchestrator(*configPath, *managed)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error initializing pipeline: %v\n", err)
		return 1
	}
	defer orch.Close()

	res, err := orch.GenerateFrontpageBridges(ctx, processor.FrontpageOptions{
		Limit:         *limit,
		OutputPath:    *outputPath,
		UseCompressed: *useCompressed,
		DryRun:        *dryRun,
		Model:         *model,
	})
	if err != nil {
		fmt.Fprintf(a.stderr, "Error generating frontpage bridges: %v\n", err)
		return 1
	}

	fmt.Fprintf(a.stdout, "\nFrontpage bridges complete: %d generated (%d failed) in %v -> %s\n",
		res.BridgesGenerated, res.BridgesFailed, res.Duration, res.OutputPath)
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
	threshold := fs.Float64("similarity-threshold", 0.85, "Rewrite hook pairs with cosine similarity above this threshold")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error loading config: %v\n", err)
		return 1
	}

	if *threshold <= 0 || *threshold >= 1 || *batch < 0 {
		fmt.Fprintln(a.stderr, "Error: similarity-threshold must be between 0 and 1 and batch must be nonnegative")
		return 1
	}
	var sm *server.ServerManager
	defer func() {
		if sm != nil {
			_ = sm.Stop()
		}
	}()
	refinerCfg := refiner.Config{
		GeminiAPIKey: geminiAPIKey,
		PostsDir:     cfg.Processing.PostsDirectory,
		PrepareEmbedding: func(ctx context.Context) (refiner.Embedder, error) {
			embeddingCfg := cfg.LlamaServer
			if embeddingCfg.IsManaged() {
				sm = server.NewServerManager(server.ManagerConfig{
					Managed: true, BinaryPath: embeddingCfg.GetBinaryPath(), ModelsDir: embeddingCfg.GetModelsDir(),
					Host: "127.0.0.1", Port: 8080, ContextSize: embeddingCfg.ContextSize,
					GPULayers: embeddingCfg.GPULayers, BatchSize: embeddingCfg.BatchSize,
					UBatchSize: embeddingCfg.UBatchSize, StartupTimeout: embeddingCfg.StartupTimeout,
				})
				if err := sm.EnsureModel(ctx, embeddingCfg.EmbeddingModel, server.ModeEmbedding); err != nil {
					return nil, err
				}
				embeddingCfg.BaseURL = sm.GetBaseURL()
			}
			return client.NewClient(embeddingCfg), nil
		},
	}

	r, err := refiner.New(refinerCfg)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error creating refiner: %v\n", err)
		return 1
	}

	res, err := r.RefineHooks(ctx, refiner.RefineOptions{
		SpecificSlugs:       parseSlugs(*posts),
		Force:               *force,
		DryRun:              *dryRun,
		Model:               *geminiModel,
		BatchSize:           *batch,
		SimilarityThreshold: *threshold,
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
		if err == flag.ErrHelp {
			return 0
		}
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
		if err == flag.ErrHelp {
			return 0
		}
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
	cloudflare := fs.Bool("cloudflare", false, "Use Cloudflare Workers AI alongside other hook workers (requires CLOUDFLARE_API_TOKEN)")
	cloudflareAccountID := fs.String("cloudflare-account-id", os.Getenv("CLOUDFLARE_ACCOUNT_ID"), "Cloudflare account ID (defaults to CLOUDFLARE_ACCOUNT_ID)")
	cloudflareModel := fs.String("cloudflare-model", client.DefaultCloudflareModel, "Cloudflare Workers AI model ID")
	cloudflareMaxTokens := fs.Int("cloudflare-max-tokens", client.DefaultCloudflareMaxTokens, "Output token limit for Cloudflare requests")
	gemini := fs.Bool("gemini", false, "Use Gemini alongside other hook workers (requires GEMINI_API_KEY)")
	geminiModel := fs.String("gemini-model", client.DefaultGeminiModel, "Gemini API model ID")
	geminiMaxTokens := fs.Int("gemini-max-tokens", client.DefaultGeminiMaxTokens, "Output token limit for Gemini requests")
	plugsky := fs.Bool("plugsky", false, "Use Plugsky alongside other hook workers (requires PLUGSKY_API_KEY)")
	plugskyModel := fs.String("plugsky-model", client.DefaultPlugskyModel, "Plugsky API model ID")
	plugskyMaxTokens := fs.Int("plugsky-max-tokens", client.DefaultPlugskyMaxTokens, "Output token limit for Plugsky requests")
	openRouter := fs.Bool("openrouter", false, "Use OpenRouter alongside llama-server for hook generation (requires OPENROUTER_API_KEY)")
	openRouterModel := fs.String("openrouter-model", client.DefaultOpenRouterModel, "OpenRouter helper model")
	openRouterMaxTokens := fs.Int("openrouter-max-tokens", client.DefaultOpenRouterMaxTokens, "Output token limit for OpenRouter requests (independent of llama-server)")
	hookModel := fs.String("hook-model", "", "Model to use for hook generation (overrides config)")
	managed := fs.Bool("managed", false, "Auto-manage llama-server process lifecycle")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
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
			SpecificSlugs:       slugList,
			UseCompressed:       *useCompressed,
			BatchSize:           *batchSize,
			DryRun:              *dryRun,
			Model:               *hookModel,
			Cloudflare:          *cloudflare,
			CloudflareAccountID: *cloudflareAccountID,
			CloudflareModel:     *cloudflareModel,
			CloudflareMaxTokens: *cloudflareMaxTokens,
			Gemini:              *gemini,
			GeminiModel:         *geminiModel,
			GeminiMaxTokens:     *geminiMaxTokens,
			Plugsky:             *plugsky,
			PlugskyModel:        *plugskyModel,
			PlugskyMaxTokens:    *plugskyMaxTokens,
			OpenRouter:          *openRouter,
			OpenRouterModel:     *openRouterModel,
			OpenRouterMaxTokens: *openRouterMaxTokens,
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
		if err == flag.ErrHelp {
			return 0
		}
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
	all := fs.Bool("all", false, "Clear both staged rankings and hooks (all posts unless --posts is set)")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	orch, err := a.buildOrchestrator(*configPath)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error initializing pipeline: %v\n", err)
		return 1
	}

	slugs := parseSlugs(*posts)

	if *hooks || *all {
		if err := orch.ClearStagedHooks(ctx, slugs); err != nil {
			fmt.Fprintf(a.stderr, "Error clearing staged hooks: %v\n", err)
			return 1
		}
		fmt.Fprintln(a.stdout, "Successfully cleared staged hooks files.")
	}
	if !*hooks || *all {
		if err := orch.ClearStaged(ctx, slugs); err != nil {
			fmt.Fprintf(a.stderr, "Error clearing staged rankings: %v\n", err)
			return 1
		}
		fmt.Fprintln(a.stdout, "Successfully cleared staged ranking files.")
	}

	return 0
}

func (a *App) runSendNewsletter(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("send-newsletter", flag.ContinueOnError)
	fs.SetOutput(a.stdout)

	configPath := fs.String("config", "config.yaml", "Path to configuration file")
	post := fs.String("post", "", "Post slug to send newsletter for (required)")
	recipients := fs.String("recipients", "", "Comma-separated list of recipient email addresses")
	recipientsFile := fs.String("recipients-file", "", "Path to file containing recipient emails (one per line)")
	lists := fs.String("lists", "", "Comma-separated list of Brevo contact list IDs (e.g., '2,5')")
	testMode := fs.Bool("test", false, "Test mode - only print what would be sent without actually sending")
	fromEmail := fs.String("from-email", "news@rippreport.com", "From email address")
	fromName := fs.String("from-name", "Ripp Report", "From name")
	template := fs.String("template", "nextone", "Brevo template name")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	if *post == "" {
		fmt.Fprintf(a.stderr, "Error: --post flag is required\n")
		fmt.Fprintf(a.stderr, "Usage: pipeline send-newsletter --post october-surprise --recipients email@example.com\n")
		return 1
	}

	// Get Brevo API key from environment (not required in test mode)
	brevoAPIKey := os.Getenv("BREVO_API_KEY")
	if brevoAPIKey == "" && !*testMode {
		fmt.Fprintf(a.stderr, "Error: BREVO_API_KEY environment variable is not set\n")
		fmt.Fprintf(a.stderr, "Tip: Use --test flag to test without actually sending\n")
		return 1
	}

	// Load config to get posts directory
	var cfg *config.Config
	if *configPath != "" {
		if _, statErr := os.Stat(*configPath); statErr == nil {
			loadedCfg, err := config.LoadConfig(*configPath)
			if err == nil {
				cfg = loadedCfg
			}
		}
	}
	if cfg == nil {
		cfg = config.DefaultConfig()
		cfg.Processing.PostsDirectory = "/Volumes/1tb/rippreport/content/p"
	}

	// Parse recipients
	var recipientList []string
	if *recipients != "" {
		recipientList = parseSlugs(*recipients) // Reuse parseSlugs for comma-separated parsing
	}
	if *recipientsFile != "" {
		fileRecipients, err := readLinesFromFile(*recipientsFile)
		if err != nil {
			fmt.Fprintf(a.stderr, "Error reading recipients file: %v\n", err)
			return 1
		}
		recipientList = append(recipientList, fileRecipients...)
	}

	if len(recipientList) == 0 {
		fmt.Fprintf(a.stderr, "Error: no recipients specified (use --recipients or --recipients-file)\n")
		return 1
	}

	// Parse list IDs if provided
	var listIDs []int
	if *lists != "" {
		listStrs := parseSlugs(*lists) // Reuse parseSlugs for comma-separated parsing
		for _, listStr := range listStrs {
			var listID int
			if _, err := fmt.Sscanf(listStr, "%d", &listID); err == nil {
				listIDs = append(listIDs, listID)
			} else {
				fmt.Fprintf(a.stderr, "Warning: invalid list ID '%s', skipping\n", listStr)
			}
		}
	}

	if len(recipientList) == 0 && len(listIDs) == 0 {
		fmt.Fprintf(a.stderr, "Error: no recipients or lists specified (use --recipients, --recipients-file, or --lists)\n")
		return 1
	}

	// Create newsletter client
	nlClient := newsletter.NewBrevoClient(newsletter.BrevoConfig{
		APIKey:         brevoAPIKey,
		FromEmail:      *fromEmail,
		FromName:       *fromName,
		TemplateName:   *template,
		PostsDirectory: cfg.Processing.PostsDirectory,
		TestMode:       *testMode,
	})

	// Send newsletter
	fmt.Fprintf(a.stdout, "Preparing newsletter for post: %s\n", *post)
	if len(recipientList) > 0 {
		fmt.Fprintf(a.stdout, "Recipients: %d\n", len(recipientList))
	}
	if len(listIDs) > 0 {
		fmt.Fprintf(a.stdout, "Lists: %v\n", listIDs)
	}
	if *testMode {
		fmt.Fprintf(a.stdout, "TEST MODE: Newsletter will not actually be sent\n")
	}
	fmt.Fprintf(a.stdout, "\n")

	result, err := nlClient.SendNewsletter(ctx, newsletter.NewsletterRequest{
		PostSlug:   *post,
		Recipients: recipientList,
		ListIDs:    listIDs,
	})
	if err != nil {
		fmt.Fprintf(a.stderr, "Error sending newsletter: %v\n", err)
		return 1
	}

	if result.Success {
		if *testMode {
			fmt.Fprintf(a.stdout, "✓ Test successful! Newsletter is ready to send.\n")
		} else {
			fmt.Fprintf(a.stdout, "✓ Newsletter sent successfully!\n")
			fmt.Fprintf(a.stdout, "  Message ID: %s\n", result.MessageID)
		}
		fmt.Fprintf(a.stdout, "  Recipients: %d\n", result.Recipients)
		fmt.Fprintf(a.stdout, "  Sent at: %s\n", result.SentAt.Format(time.RFC3339))
		return 0
	}

	fmt.Fprintf(a.stderr, "✗ Newsletter send failed: %v\n", result.Error)
	return 1
}

// readLinesFromFile reads a file and returns non-empty lines
func readLinesFromFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") { // Skip empty lines and comments
			result = append(result, trimmed)
		}
	}
	return result, nil
}

func (a *App) runSocialPost(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("social-post", flag.ContinueOnError)
	fs.SetOutput(a.stdout)

	configPath := fs.String("config", "config.yaml", "Path to configuration file")
	post := fs.String("post", "", "Post slug to generate social media text for (required)")
	format := fs.String("format", "default", "Output format: default, twitter, facebook, linkedin")
	maxLength := fs.Int("max-length", 0, "Maximum length for summary (0 = no limit)")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	if *post == "" {
		fmt.Fprintf(a.stderr, "Error: --post flag is required\n")
		fmt.Fprintf(a.stderr, "Usage: pipeline social-post --post october-surprise\n")
		return 1
	}

	// Load config to get posts directory
	var cfg *config.Config
	if *configPath != "" {
		if _, statErr := os.Stat(*configPath); statErr == nil {
			loadedCfg, err := config.LoadConfig(*configPath)
			if err == nil {
				cfg = loadedCfg
			}
		}
	}
	if cfg == nil {
		cfg = config.DefaultConfig()
		cfg.Processing.PostsDirectory = "/Volumes/1tb/rippreport/content/p"
	}

	// Read the post
	postManager := processor.NewPostManager(cfg.Processing.PostsDirectory)
	postPath := models.PostPath(filepath.Join(cfg.Processing.PostsDirectory, *post, "index.md"))
	postData, err := postManager.ReadPost(postPath)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error reading post: %v\n", err)
		return 1
	}

	// Build the post URL
	postURL := fmt.Sprintf("https://rippreport.com/p/%s/", *post)

	// Get the description/summary
	summary := ""
	if postData.FrontMatter != nil && postData.FrontMatter.Description != "" {
		summary = postData.FrontMatter.Description
	} else if len(postData.Body) > 0 {
		// Truncate body if no description
		summary = postData.Body
		if len(summary) > 300 {
			summary = summary[:297] + "..."
		}
	}

	// Apply max length if specified
	if *maxLength > 0 && len(summary) > *maxLength {
		summary = summary[:*maxLength-3] + "..."
	}

	// Format output based on platform
	switch strings.ToLower(*format) {
	case "twitter", "x":
		// Twitter/X format - concise with link
		fmt.Fprintf(a.stdout, "%s\n\n%s\n", summary, postURL)
	case "facebook", "fb":
		// Facebook format - more descriptive
		fmt.Fprintf(a.stdout, "%s\n\n%s\n", summary, postURL)
	case "linkedin":
		// LinkedIn format - professional
		fmt.Fprintf(a.stdout, "%s\n\n%s\n", summary, postURL)
	default:
		// Default format - clean and simple
		fmt.Fprintf(a.stdout, "%s\n\n%s\n", summary, postURL)
	}

	return 0
}

func (a *App) runNewPost(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("new-post", flag.ContinueOnError)
	fs.SetOutput(a.stdout)

	configPath := fs.String("config", "config.yaml", "Path to configuration file")
	title := fs.String("title", "", "Post title (required)")
	slug := fs.String("slug", "", "Post slug (auto-generated from title if not provided)")
	fromMail := fs.Bool("email", false, "Retrieve latest post from local Apple Mail")
	fs.BoolVar(fromMail, "mail", false, "Alias for --email")
	fromEmail := fs.String("from-email", "fpripp1@gmail.com", "Sender email address to search for in Apple Mail")
	emailDays := fs.Int("email-days", 2, "Number of days back to search in Apple Mail (default 2)")
	runPipeline := fs.Bool("pipeline", false, "Run recommendation pipeline (generate-related, generate-hooks, refine-hooks) on the new post")
	skipRefine := fs.Bool("skip-refine", false, "Skip hook refinement with Gemini when --pipeline is used")
	podcast := fs.Bool("podcast", false, "Create as podcast post (includes youtube field)")
	youtubeURL := fs.String("youtube", "", "YouTube video URL or ID (fetches summary, thumbnail, and embeds shortcodes)")
	categories := fs.String("categories", "", "Comma-separated list of categories")
	tags := fs.String("tags", "", "Comma-separated list of tags")
	description := fs.String("description", "", "Post description/summary")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	if *title == "" {
		fmt.Fprintf(a.stderr, "Error: --title flag is required\n")
		fmt.Fprintf(a.stderr, "Usage: pipeline new-post --title \"My Post Title\"\n")
		return 1
	}

	// Generate slug from title if not provided (keeps folder naming accurate)
	postSlug := *slug
	if postSlug == "" {
		postSlug = generateSlug(*title)
	}

	isPodcast := *podcast
	var youtubeID string
	bodyContent := "Write your post content here.\n"

	// If email flag is enabled, fetch latest email from Apple Mail
	if *fromMail {
		fmt.Fprintf(a.stdout, "Searching Apple Mail for latest post matching %q from %s (last %d days)...\n", *title, *fromEmail, *emailDays)
		opts := email.FetchOptions{
			Title:      *title,
			FromEmail:  *fromEmail,
			MaxAgeDays: *emailDays,
		}
		emailRes, err := email.FetchFromAppleMail(ctx, opts)
		if err != nil {
			fmt.Fprintf(a.stderr, "Error fetching from Apple Mail: %v\n", err)
			return 1
		}

		fmt.Fprintf(a.stdout, "✓ Found email: %s\n", emailRes.Subject)

		// If title flag was provided as a slug, prefer the email subject for display
		if strings.Contains(*title, "-") && emailRes.Subject != "" {
			*title = emailRes.Subject
		} else if strings.Contains(*title, "-") {
			*title = email.SlugToSearchQuery(postSlug)
		}

		// Body content with {{< aicomments >}} at the end
		bodyContent = fmt.Sprintf("%s\n\n{{< aicomments >}}\n", emailRes.Body)

		if *description == "" && emailRes.Body != "" {
			lines := strings.Split(emailRes.Body, "\n")
			for _, l := range lines {
				trimmed := strings.TrimSpace(l)
				if trimmed != "" && !strings.HasPrefix(trimmed, "{{<") {
					*description = trimmed
					if len(*description) > 200 {
						*description = (*description)[:197] + "..."
					}
					break
				}
			}
		}
	}

	// If youtube URL or ID is provided, retrieve video info and thumbnail
	if *youtubeURL != "" {
		isPodcast = true
		fmt.Fprintf(a.stdout, "Fetching YouTube video info for: %s...\n", *youtubeURL)
		info, err := youtube.FetchVideoInfo(*youtubeURL)
		if err != nil {
			fmt.Fprintf(a.stderr, "Warning: failed to fetch YouTube info via yt-dlp: %v\n", err)
			// Try extracting video ID as fallback
			if extractedID, extractErr := youtube.ExtractVideoID(*youtubeURL); extractErr == nil {
				youtubeID = extractedID
			}
		} else {
			youtubeID = info.VideoID
			summary := strings.TrimSpace(info.Description)
			if summary != "" {
				bodyContent = fmt.Sprintf("%s\n\n{{< youtubecomments >}}\n{{< aicomments >}}\n", summary)
				if *description == "" {
					lines := strings.Split(summary, "\n")
					for _, line := range lines {
						trimmed := strings.TrimSpace(line)
						if trimmed != "" {
							*description = trimmed
							if len(*description) > 200 {
								*description = (*description)[:197] + "..."
							}
							break
						}
					}
				}
			} else {
				bodyContent = "{{< youtubecomments >}}\n{{< aicomments >}}\n"
			}
		}
	} else if isPodcast {
		bodyContent = "Write your post content here.\n\n{{< youtubecomments >}}\n{{< aicomments >}}\n"
	}

	// Load config to get posts directory
	var cfg *config.Config
	if *configPath != "" {
		if _, statErr := os.Stat(*configPath); statErr == nil {
			loadedCfg, err := config.LoadConfig(*configPath)
			if err == nil {
				cfg = loadedCfg
			}
		}
	}
	if cfg == nil {
		cfg = config.DefaultConfig()
		cfg.Processing.PostsDirectory = "/Volumes/1tb/rippreport/content/p"
	}

	// Create post directory
	postDir := filepath.Join(cfg.Processing.PostsDirectory, postSlug)
	if _, err := os.Stat(postDir); err == nil {
		fmt.Fprintf(a.stderr, "Error: post directory already exists: %s\n", postDir)
		return 1
	}

	if err := os.MkdirAll(postDir, 0755); err != nil {
		fmt.Fprintf(a.stderr, "Error creating post directory: %v\n", err)
		return 1
	}

	// Download thumbnail if YouTube video ID is available
	var downloadedThumbPath string
	if youtubeID != "" {
		contentImgDir := filepath.Join(filepath.Dir(cfg.Processing.PostsDirectory), "img")
		thumbFilename := fmt.Sprintf("%s.png", postSlug)
		thumbPath := filepath.Join(contentImgDir, thumbFilename)

		fmt.Fprintf(a.stdout, "Downloading YouTube thumbnail to: %s...\n", thumbPath)
		if err := youtube.DownloadThumbnail(youtubeID, thumbPath); err != nil {
			fmt.Fprintf(a.stderr, "Warning: could not download thumbnail: %v\n", err)
		} else {
			downloadedThumbPath = thumbPath
			fmt.Fprintf(a.stdout, "✓ Saved thumbnail: %s\n", thumbPath)

			// Copy to assets/img if directory exists
			repoRoot := filepath.Dir(filepath.Dir(cfg.Processing.PostsDirectory))
			assetsImgDir := filepath.Join(repoRoot, "assets", "img")
			if stat, err := os.Stat(assetsImgDir); err == nil && stat.IsDir() {
				_ = copyFile(thumbPath, filepath.Join(assetsImgDir, thumbFilename))
			}

			// Also copy into the post directory
			_ = copyFile(thumbPath, filepath.Join(postDir, thumbFilename))
		}
	}

	// Parse categories and tags
	var categoryList []string
	if *categories != "" {
		categoryList = parseSlugs(*categories)
	}

	var tagList []string
	if *tags != "" {
		tagList = parseSlugs(*tags)
	}

	// Create front matter
	frontMatter := buildFrontMatter(*title, postSlug, *description, youtubeID, categoryList, tagList, isPodcast, bodyContent)

	// Create index.md
	indexPath := filepath.Join(postDir, "index.md")
	if err := os.WriteFile(indexPath, []byte(frontMatter), 0644); err != nil {
		fmt.Fprintf(a.stderr, "Error creating index.md: %v\n", err)
		return 1
	}

	fmt.Fprintf(a.stdout, "✓ Created new post: %s\n", postSlug)
	fmt.Fprintf(a.stdout, "  Directory: %s\n", postDir)
	fmt.Fprintf(a.stdout, "  Index: %s\n", indexPath)
	if downloadedThumbPath != "" {
		fmt.Fprintf(a.stdout, "  Thumbnail: %s (downloaded)\n", downloadedThumbPath)
	} else {
		fmt.Fprintf(a.stdout, "  Thumbnail: img/%s.png (you need to create this)\n", postSlug)
	}
	if isPodcast {
		if youtubeID != "" {
			fmt.Fprintf(a.stdout, "  Type: Podcast (YouTube ID: %s)\n", youtubeID)
		} else {
			fmt.Fprintf(a.stdout, "  Type: Podcast (remember to add YouTube ID)\n")
		}
	}
	fmt.Fprintf(a.stdout, "\nNext steps:\n")
	fmt.Fprintf(a.stdout, "  1. Review your content in: %s\n", indexPath)
	if downloadedThumbPath == "" {
		fmt.Fprintf(a.stdout, "  2. Create thumbnail image: /Volumes/1tb/rippreport/content/img/%s.png\n", postSlug)
	}
	if isPodcast && youtubeID == "" {
		fmt.Fprintf(a.stdout, "  3. Add YouTube video ID to the 'youtube:' field\n")
	}

	// If --pipeline flag is set, run the full pipeline sequence on the new post
	if *runPipeline {
		fmt.Fprintf(a.stdout, "\n========================================\n")
		fmt.Fprintf(a.stdout, "Running Recommendation Pipeline for: %s\n", postSlug)
		fmt.Fprintf(a.stdout, "========================================\n\n")

		// 1. Generate related articles
		fmt.Fprintf(a.stdout, "--> Step 1/3: Generating related articles for %s...\n", postSlug)
		relatedArgs := []string{"--posts", postSlug}
		if *configPath != "" {
			relatedArgs = append(relatedArgs, "--config", *configPath)
		}
		if code := a.runGenerateRelated(ctx, relatedArgs); code != 0 {
			fmt.Fprintf(a.stderr, "Warning: generate-related completed with exit code %d\n", code)
		}

		// 2. Generate hooks
		fmt.Fprintf(a.stdout, "\n--> Step 2/3: Generating contextual hooks for %s...\n", postSlug)
		hookArgs := []string{"--posts", postSlug}
		if *configPath != "" {
			hookArgs = append(hookArgs, "--config", *configPath)
		}
		if code := a.runGenerateHooks(ctx, hookArgs); code != 0 {
			fmt.Fprintf(a.stderr, "Warning: generate-hooks completed with exit code %d\n", code)
		}

		// 3. Refine hooks with Gemini (unless skipped or no API key)
		if !*skipRefine {
			if os.Getenv("GEMINI_API_KEY") != "" {
				fmt.Fprintf(a.stdout, "\n--> Step 3/4: Refining hooks with Gemini for %s...\n", postSlug)
				refineArgs := []string{"--posts", postSlug}
				if *configPath != "" {
					refineArgs = append(refineArgs, "--config", *configPath)
				}
				if code := a.runRefineHooks(ctx, refineArgs); code != 0 {
					fmt.Fprintf(a.stderr, "Warning: refine-hooks completed with exit code %d\n", code)
				}
			} else {
				fmt.Fprintf(a.stdout, "\n--> Step 3/4: Skipping hook refinement (GEMINI_API_KEY not set)\n")
			}
		}

		// 4. Print Facebook social media post
		fmt.Fprintf(a.stdout, "\n========================================\n")
		fmt.Fprintf(a.stdout, "--> Step 4/4: Facebook Social Media Post:\n")
		fmt.Fprintf(a.stdout, "========================================\n\n")
		_ = a.runSocialPost(ctx, []string{"--post", postSlug, "--format", "facebook"})
		fmt.Fprintf(a.stdout, "\n========================================\n")

		fmt.Fprintf(a.stdout, "\n✓ Pipeline sequence completed for %s!\n", postSlug)
	}

	return 0
}

func (a *App) runGenerateAlt(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("generate-alt", flag.ContinueOnError)
	fs.SetOutput(a.stdout)

	configPath := fs.String("config", "config.yaml", "Path to configuration file")
	post := fs.String("post", "", "Post slug to generate alt tag for (required)")
	visionModel := fs.String("vision-model", "", "Local visual-capable model path or name (required or set in config)")
	refillCache := fs.Bool("refill-cache", true, "Clear and refill embedding cache after updating alt tag")
	fs.BoolVar(refillCache, "refresh-cache", true, "Alias for --refill-cache")
	managed := fs.Bool("managed", false, "Auto-manage llama-server process lifecycle")
	dryRun := fs.Bool("dry-run", false, "Simulate without modifying post front matter")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	if *post == "" {
		fmt.Fprintf(a.stderr, "Error: --post flag is required\n")
		fmt.Fprintf(a.stderr, "Usage: pipeline generate-alt --post october-surprise --vision-model <model>\n")
		return 1
	}

	orch, err := a.buildOrchestrator(*configPath, *managed)
	if err != nil {
		fmt.Fprintf(a.stderr, "Error initializing orchestrator: %v\n", err)
		return 1
	}
	defer orch.Close()

	res, err := orch.GenerateAltTag(ctx, orchestrator.AltTagOptions{
		PostSlug:    *post,
		VisionModel: *visionModel,
		DryRun:      *dryRun,
	})
	if err != nil {
		fmt.Fprintf(a.stderr, "Error generating alt tag: %v\n", err)
		return 1
	}

	fmt.Fprintf(a.stdout, "\n✓ Successfully generated alt tag for %s in %v:\n", res.PostSlug, res.Duration)
	fmt.Fprintf(a.stdout, "  Thumbnail: %s\n", res.ThumbnailPath)
	fmt.Fprintf(a.stdout, "  Alt tag:   %s\n", res.AltTag)

	if *refillCache && !*dryRun {
		orch.Close()

		// 1. Clear the embeddings cache file
		cacheFile := ".embedding_cache.json"
		if *configPath != "" {
			if cfg, err := config.LoadConfig(*configPath); err == nil && cfg.Storage.CacheFile != "" {
				cacheFile = cfg.Storage.CacheFile
			}
		}
		fmt.Fprintf(a.stdout, "\nClearing embeddings cache (%s)...\n", cacheFile)
		_ = os.Remove(cacheFile)

		// 2. Refill the cache by running generate-related with --dry-run --vector-only
		fmt.Fprintf(a.stdout, "Refilling embeddings cache via generate-related --dry-run --vector-only --posts %s...\n\n", *post)
		relatedArgs := []string{"--dry-run", "--vector-only", "--posts", *post}
		if *configPath != "" {
			relatedArgs = append(relatedArgs, "--config", *configPath)
		}
		if *managed {
			relatedArgs = append(relatedArgs, "--managed")
		}
		if code := a.runGenerateRelated(ctx, relatedArgs); code != 0 {
			fmt.Fprintf(a.stderr, "Warning: generate-related exited with code %d\n", code)
		} else {
			fmt.Fprintf(a.stdout, "\n✓ Embeddings cache cleared and refilled for next post!\n")
		}
	}

	return 0
}

func (a *App) runRepost(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("repost", flag.ContinueOnError)
	fs.SetOutput(a.stdout)

	configPath := fs.String("config", "config.yaml", "Path to configuration file")
	post := fs.String("post", "", "Post slug to repost (required)")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	if *post == "" {
		fmt.Fprintf(a.stderr, "Error: --post flag is required\n")
		fmt.Fprintf(a.stderr, "Usage: pipeline repost --post alabama-special-election\n")
		return 1
	}

	// Load config to get posts directory
	var cfg *config.Config
	if *configPath != "" {
		if _, statErr := os.Stat(*configPath); statErr == nil {
			loadedCfg, err := config.LoadConfig(*configPath)
			if err == nil {
				cfg = loadedCfg
			}
		}
	}
	if cfg == nil {
		cfg = config.DefaultConfig()
		cfg.Processing.PostsDirectory = "/Volumes/1tb/rippreport/content/p"
	}

	sourceSlug := *post
	sourceDir := filepath.Join(cfg.Processing.PostsDirectory, sourceSlug)

	// Verify source post exists
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		fmt.Fprintf(a.stderr, "Error: source post directory does not exist: %s\n", sourceDir)
		return 1
	}

	// Determine the next available repost suffix (-2, -3, -4, ...)
	destSlug := nextRepostSlug(cfg.Processing.PostsDirectory, sourceSlug)
	destDir := filepath.Join(cfg.Processing.PostsDirectory, destSlug)

	fmt.Fprintf(a.stdout, "Reposting %s → %s\n", sourceSlug, destSlug)

	// Copy entire directory
	if err := copyDir(sourceDir, destDir); err != nil {
		fmt.Fprintf(a.stderr, "Error copying post directory: %v\n", err)
		return 1
	}

	// Update the index.md in the destination: set date to today, categories to repost
	indexPath := filepath.Join(destDir, "index.md")
	if err := updateRepostFrontMatter(indexPath); err != nil {
		fmt.Fprintf(a.stderr, "Error updating front matter: %v\n", err)
		// Clean up on failure
		_ = os.RemoveAll(destDir)
		return 1
	}

	fmt.Fprintf(a.stdout, "✓ Created repost: %s\n", destSlug)
	fmt.Fprintf(a.stdout, "  Directory: %s\n", destDir)
	fmt.Fprintf(a.stdout, "  Date set to: %s\n", time.Now().Format("2006-01-02"))
	fmt.Fprintf(a.stdout, "  Category set to: repost\n")

	return 0
}

// nextRepostSlug finds the next available repost slug.
// Given base slug "foo", it checks for "foo-2", "foo-3", etc.
// If the slug already ends with a number suffix (e.g. "foo-2"), it strips it
// to use the original base for numbering.
func nextRepostSlug(postsDir, slug string) string {
	// Strip any existing trailing -N suffix to find the true base
	base := slug
	for i := len(base) - 1; i >= 0; i-- {
		if base[i] == '-' {
			suffix := base[i+1:]
			allDigits := len(suffix) > 0
			for _, ch := range suffix {
				if ch < '0' || ch > '9' {
					allDigits = false
					break
				}
			}
			if allDigits {
				base = base[:i]
			}
			break
		}
	}

	// Find the next available number starting from 2
	for n := 2; ; n++ {
		candidate := fmt.Sprintf("%s-%d", base, n)
		candidateDir := filepath.Join(postsDir, candidate)
		if _, err := os.Stat(candidateDir); os.IsNotExist(err) {
			return candidate
		}
	}
}

// copyDir recursively copies a directory tree from src to dst.
func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// updateRepostFrontMatter reads the index.md file and replaces the date with today
// and the categories with [repost], preserving all other content.
func updateRepostFrontMatter(indexPath string) error {
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", indexPath, err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	// Find the front matter boundaries (--- delimiters)
	startIdx := -1
	endIdx := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "---" {
			if startIdx == -1 {
				startIdx = i
			} else {
				endIdx = i
				break
			}
		}
	}

	if startIdx == -1 || endIdx == -1 {
		return fmt.Errorf("could not find front matter delimiters in %s", indexPath)
	}

	today := time.Now().Format("2006-01-02")

	// Process front matter lines: replace date and categories
	var newLines []string
	newLines = append(newLines, lines[:startIdx+1]...) // include opening ---

	inCategories := false
	dateReplaced := false
	categoriesReplaced := false

	for i := startIdx + 1; i < endIdx; i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Handle multi-line categories block
		if inCategories {
			if strings.HasPrefix(trimmed, "- ") {
				// Skip old category list items
				continue
			}
			// Not a list item, so categories block ended
			inCategories = false
		}

		// Replace date line
		if strings.HasPrefix(trimmed, "date:") {
			newLines = append(newLines, fmt.Sprintf("date: %s", today))
			dateReplaced = true
			continue
		}

		// Replace categories
		if strings.HasPrefix(trimmed, "categories:") {
			newLines = append(newLines, "categories:")
			newLines = append(newLines, "- repost")
			categoriesReplaced = true

			// Check if categories were inline (e.g. "categories: news") vs block
			afterColon := strings.TrimSpace(strings.TrimPrefix(trimmed, "categories:"))
			if afterColon == "" {
				// Block style - mark that we need to skip subsequent "- item" lines
				inCategories = true
			}
			// If inline value, we already replaced it
			continue
		}

		newLines = append(newLines, line)
	}

	// If there was no date line, add one
	if !dateReplaced {
		// Insert date right after opening ---
		insertAt := startIdx + 1
		newLines = append(newLines[:insertAt+1], newLines[insertAt:]...)
		newLines[insertAt] = fmt.Sprintf("date: %s", today)
	}

	// If there was no categories line, add one before the closing ---
	if !categoriesReplaced {
		newLines = append(newLines, "categories:")
		newLines = append(newLines, "- repost")
	}

	// Add closing --- and the rest of the file
	newLines = append(newLines, lines[endIdx:]...)

	result := strings.Join(newLines, "\n")
	return os.WriteFile(indexPath, []byte(result), 0644)
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// generateSlug converts a title to a URL-friendly slug
func generateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces and special characters with hyphens
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '-'
	}, slug)

	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// buildFrontMatter creates the YAML front matter and content for a new post
func buildFrontMatter(title, slug, description, youtubeID string, categories, tags []string, podcast bool, bodyContent string) string {
	var builder strings.Builder

	builder.WriteString("---\n")
	builder.WriteString(fmt.Sprintf("title: %s\n", title))
	builder.WriteString(fmt.Sprintf("date: %s\n", time.Now().Format("2006-01-02T15:04:05Z07:00")))

	if description != "" {
		builder.WriteString(fmt.Sprintf("description: %s\n", description))
	}

	if len(categories) > 0 {
		builder.WriteString("categories:\n")
		for _, cat := range categories {
			builder.WriteString(fmt.Sprintf("  - %s\n", cat))
		}
	}

	if len(tags) > 0 {
		builder.WriteString("tags:\n")
		for _, tag := range tags {
			builder.WriteString(fmt.Sprintf("  - %s\n", tag))
		}
	}

	if podcast {
		if youtubeID != "" {
			builder.WriteString(fmt.Sprintf("youtube: %s\n", youtubeID))
		} else {
			builder.WriteString("youtube: \n")
		}
	}

	builder.WriteString(fmt.Sprintf("thumbnail: img/%s.png\n", slug))
	builder.WriteString("---\n\n")
	if bodyContent != "" {
		builder.WriteString(bodyContent)
		if !strings.HasSuffix(bodyContent, "\n") {
			builder.WriteString("\n")
		}
	} else {
		builder.WriteString("Write your post content here.\n")
	}

	return builder.String()
}
