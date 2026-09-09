package orchestrator

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/cache"
	"github.com/rippreport/go-content-pipeline/internal/client"
	"github.com/rippreport/go-content-pipeline/internal/config"
	"github.com/rippreport/go-content-pipeline/internal/errors"
	"github.com/rippreport/go-content-pipeline/internal/models"
	"github.com/rippreport/go-content-pipeline/internal/processor"
	"github.com/rippreport/go-content-pipeline/internal/server"
)

// RankingOptions configures article recommendation generation.
type RankingOptions struct {
	Resume        bool
	Stage         bool
	RefreshRecent bool
	VectorOnly    bool
	UseReranker   bool
	SpecificSlugs []string
	Workers       int
	DryRun        bool
	Model         string
}

// RankingResult summarizes related article generation.
type RankingResult struct {
	ProcessedPosts int
	FailedPosts    int
	SkippedPosts   int
	CacheHits      int
	Duration       time.Duration
}

// HookOptions configures contextual sidebar hook generation.
type HookOptions struct {
	SpecificSlugs []string
	UseCompressed bool
	BatchSize     int
	DryRun        bool
	FromStaging   bool
	Stage         bool // Save to staging files instead of writing directly
	Model         string
}

// HookResult summarizes hook generation.
type HookResult struct {
	TotalHooks  int
	FailedHooks int
	Duration    time.Duration
}

// PipelineOptions configures full pipeline execution.
type PipelineOptions struct {
	RankingOptions
	HookOptions
}

// PipelineResult summarizes end-to-end pipeline execution.
type PipelineResult struct {
	Ranking *RankingResult
	Hooks   *HookResult
	Summary errors.PipelineSummary
}

// Orchestrator coordinates the content recommendation and hook generation workflow.
type Orchestrator struct {
	cfg              *config.Config
	postManager      *processor.PostManager
	embeddingCache   *cache.EmbeddingCache
	llamaClient      client.LlamaClient
	progressTracker  *processor.ProgressTracker
	backupManager    *processor.BackupManager
	stagingManager   *processor.StagingManager
	compressionCache *processor.CompressionCache
	serverManager    *server.ServerManager
	errorCollector   *errors.ErrorCollector
	memoryManager    *processor.MemoryManager
	mu               sync.Mutex
}

// NewOrchestrator creates a fully configured Orchestrator instance.
func NewOrchestrator(
	cfg *config.Config,
	llamaClient client.LlamaClient,
	postManager *processor.PostManager,
	embeddingCache *cache.EmbeddingCache,
	progressTracker *processor.ProgressTracker,
	backupManager *processor.BackupManager,
	stagingManager *processor.StagingManager,
	compressionCache *processor.CompressionCache,
) *Orchestrator {
	var sm *server.ServerManager
	if cfg != nil && cfg.LlamaServer.IsManaged() {
		sm = server.NewServerManager(server.ManagerConfig{
			Managed:          cfg.LlamaServer.IsManaged(),
			BinaryPath:       cfg.LlamaServer.GetBinaryPath(),
			ModelsDir:        cfg.LlamaServer.GetModelsDir(),
			Host:             "127.0.0.1",
			Port:             8080,
			ContextSize:      cfg.LlamaServer.ContextSize,
			GPULayers:        cfg.LlamaServer.GPULayers,
			StartupTimeout:   cfg.LlamaServer.StartupTimeout,
			DraftModel:       cfg.LlamaServer.DraftModel,
			SpecType:         cfg.LlamaServer.SpecType,
			SpecDraftNMax:    cfg.LlamaServer.SpecDraftNMax,
			ReasoningEffort:  cfg.LlamaServer.ReasoningEffort,
			PreserveThinking: cfg.LlamaServer.PreserveThinking,
		})
	}

	return &Orchestrator{
		cfg:              cfg,
		llamaClient:      llamaClient,
		postManager:      postManager,
		embeddingCache:   embeddingCache,
		progressTracker:  progressTracker,
		backupManager:    backupManager,
		stagingManager:   stagingManager,
		compressionCache: compressionCache,
		serverManager:    sm,
		errorCollector:   errors.NewErrorCollector(),
		memoryManager:    processor.NewMemoryManager(10),
	}
}

// SetServerManager configures a custom ServerManager instance.
func (o *Orchestrator) SetServerManager(sm *server.ServerManager) {
	o.serverManager = sm
}

// GetServerManager returns the current ServerManager instance.
func (o *Orchestrator) GetServerManager() *server.ServerManager {
	return o.serverManager
}

// Close gracefully terminates any managed llama-server process.
func (o *Orchestrator) Close() error {
	if o.serverManager != nil {
		return o.serverManager.Stop()
	}
	return nil
}

// GenerateRelatedArticles discovers posts, computes embeddings, scores candidates,
// and saves rankings either directly to post front matter or to staging files.
//
// Requirements: 1.1 through 1.8, 2.1 through 2.8, 3.1 through 3.8, 4.1 through 4.8
func (o *Orchestrator) GenerateRelatedArticles(ctx context.Context, opts RankingOptions) (*RankingResult, error) {
	startTime := time.Now()

	// Discover posts to process (specific slugs or all with filters)
	filter := models.PostFilter{
		SpecificSlugs: opts.SpecificSlugs,
	}
	if opts.RefreshRecent {
		// Only rank posts from the last 2 years for refresh recent
		cutoff := time.Now().AddDate(-2, 0, 0)
		filter.StartDate = &cutoff
	}

	postPathsToProcess, err := o.postManager.DiscoverPosts(filter)
	if err != nil {
		return nil, fmt.Errorf("failed discovering posts to process: %w", err)
	}

	// For the corpus (candidates), always load ALL posts unless we're doing a refresh of recent posts
	var corpusPaths []models.PostPath
	if len(opts.SpecificSlugs) > 0 && !opts.RefreshRecent {
		// When processing specific posts, load all posts as potential candidates
		allPostsFilter := models.PostFilter{} // No filters = all posts
		corpusPaths, err = o.postManager.DiscoverPosts(allPostsFilter)
		if err != nil {
			return nil, fmt.Errorf("failed discovering corpus posts: %w", err)
		}
	} else {
		// Otherwise corpus = posts to process
		corpusPaths = postPathsToProcess
	}

	total := len(postPathsToProcess)
	result := &RankingResult{}

	// Ensure embedding model is loaded if server is managed
	if o.serverManager != nil && o.cfg != nil && o.cfg.LlamaServer.IsManaged() {
		fmt.Printf("[DEBUG] Loading embedding model: %s\n", o.cfg.LlamaServer.EmbeddingModel)
		if err := o.serverManager.EnsureModel(ctx, o.cfg.LlamaServer.EmbeddingModel, server.ModeEmbedding); err != nil {
			return nil, fmt.Errorf("failed to load embedding model: %w", err)
		}
		fmt.Printf("[DEBUG] Embedding model loaded successfully\n")
	} else {
		fmt.Printf("[DEBUG] ServerManager: %v, Config: %v, IsManaged: %v\n",
			o.serverManager != nil, o.cfg != nil, o.cfg != nil && o.cfg.LlamaServer.IsManaged())
	}

	if o.embeddingCache != nil {
		_ = o.embeddingCache.Load()
	}
	if o.progressTracker != nil {
		_ = o.progressTracker.Load()
	}

	// Load posts and calculate/fetch embeddings for the entire corpus
	corpusEmbeddings := make(map[string][]float32)
	postsByPath := make(map[string]*models.Post)
	var allPosts []*models.Post

	for _, pPath := range corpusPaths {
		post, readErr := o.postManager.ReadPost(pPath)
		if readErr != nil {
			o.errorCollector.Add(errors.NewFileIOError("read_post", readErr.Error(), readErr, false))
			continue
		}
		postsByPath[string(pPath)] = post
		allPosts = append(allPosts, post)

		// Check cache for embedding
		if vec, hit := o.embeddingCache.Get(string(pPath), post.ContentHash); hit {
			corpusEmbeddings[string(pPath)] = vec
			result.CacheHits++
		} else {
			// Generate embedding
			embedContent, cErr := o.postManager.GetEmbeddingContent(post)
			if cErr != nil {
				embedContent = post.Body
			}
			var vec []float32
			if o.llamaClient != nil {
				vec, err = o.llamaClient.GenerateEmbedding(ctx, embedContent)
			} else {
				// Fallback dummy vector for testing
				vec = []float32{1.0, 0.0, 0.0}
				err = nil
			}

			if err != nil {
				o.errorCollector.Add(errors.NewNetworkError("generate_embedding", err.Error(), err, true))
				continue
			}
			corpusEmbeddings[string(pPath)] = vec
			o.embeddingCache.Put(string(pPath), post.ContentHash, vec)
		}
	}

	// Persist updated cache
	_ = o.embeddingCache.Persist()

	topN := o.cfg.Processing.TopCandidates
	if topN <= 0 {
		topN = 5
	}

	// Ensure scoring/reranker model is loaded if server is managed
	if o.serverManager != nil && o.cfg != nil && o.cfg.LlamaServer.IsManaged() && !opts.VectorOnly {
		scoringModel := opts.Model
		if scoringModel == "" {
			scoringModel = o.cfg.LlamaServer.GetScoringModel()
		}
		mode := server.ModeCompletion
		if opts.UseReranker || o.cfg.ShouldUseReranker() {
			mode = server.ModeReranker
		}
		if err := o.serverManager.EnsureModel(ctx, scoringModel, mode); err != nil {
			return nil, fmt.Errorf("failed to load scoring model: %w", err)
		}
	}

	// Process rankings for each target post
	for idx, pPath := range postPathsToProcess {
		post, ok := postsByPath[string(pPath)]
		if !ok {
			result.FailedPosts++
			continue
		}

		if opts.Resume && o.progressTracker.IsProcessed(post.Slug) {
			result.SkippedPosts++
			continue
		}

		// Progress notification
		fmt.Printf("[%d/%d] Processing %s\n", idx+1, total, post.Slug)

		var finalArticles []models.RelatedArticle

		if post.IsBackstory() {
			finalArticles = processor.RankBackstoryRelated(post, allPosts)
		} else {
			queryVec := corpusEmbeddings[string(pPath)]
			if len(queryVec) == 0 {
				result.FailedPosts++
				_ = o.progressTracker.MarkFailed(post.Slug, fmt.Errorf("no embedding available"))
				continue
			}

			candidates := processor.GetTopCandidates(queryVec, corpusEmbeddings, topN*2, []string{post.Slug})
			candidates = processor.FilterBackstoryFromNonBackstory(candidates, allPosts)

			scoringModel := opts.Model
			if scoringModel == "" {
				scoringModel = o.cfg.LlamaServer.GetScoringModel()
			}

			useReranker := opts.UseReranker || o.cfg.ShouldUseReranker()

			scoringReq := models.ScoringRequest{
				CurrentPost: post,
				Candidates:  candidates,
				VectorOnly:  opts.VectorOnly,
				Model:       scoringModel,
				UseReranker: useReranker,
			}

			weights := models.ScoreWeights{
				Relevance: o.cfg.Scoring.Weights.Relevance,
				Recency:   o.cfg.Scoring.Weights.Recency,
				Length:    o.cfg.Scoring.Weights.Length,
				Category:  o.cfg.Scoring.Weights.Category,
			}

			scoredArticles, sErr := processor.ScoreCandidates(ctx, o.llamaClient, scoringReq, weights)
			if sErr != nil {
				o.errorCollector.Add(errors.NewProcessingError("score_candidates", sErr.Error(), sErr, false))
				result.FailedPosts++
				_ = o.progressTracker.MarkFailed(post.Slug, sErr)
				continue
			}

			for i := 0; i < len(scoredArticles) && i < topN; i++ {
				// Enrich title from post if known
				title := scoredArticles[i].Slug
				for _, p := range allPosts {
					if p.Slug == scoredArticles[i].Slug {
						title = p.FrontMatter.Title
						break
					}
				}
				scoredArticles[i].Title = title
				finalArticles = append(finalArticles, scoredArticles[i].ToRelatedArticle())
			}

			if opts.RefreshRecent && len(post.FrontMatter.RelatedArticles) > 0 {
				finalArticles = processor.MergeWithExisting(finalArticles, post.FrontMatter.RelatedArticles, 10)
			}
		}

		if !opts.DryRun {
			if opts.Stage {
				if sErr := o.stagingManager.SaveStaged(post.Slug, finalArticles); sErr != nil {
					o.errorCollector.Add(errors.NewFileIOError("save_staged", sErr.Error(), sErr, false))
					result.FailedPosts++
					_ = o.progressTracker.MarkFailed(post.Slug, sErr)
					continue
				}
			} else {
				// Backup before updating
				if o.backupManager != nil {
					_, _ = o.backupManager.Backup(string(pPath))
				}
				if uErr := o.postManager.UpdateRelatedArticles(pPath, finalArticles); uErr != nil {
					o.errorCollector.Add(errors.NewFileIOError("update_post", uErr.Error(), uErr, false))
					result.FailedPosts++
					_ = o.progressTracker.MarkFailed(post.Slug, uErr)
					continue
				}
			}
		}

		_ = o.progressTracker.MarkProcessed(post.Slug)
		result.ProcessedPosts++

		o.memoryManager.CheckAndGC(result.ProcessedPosts)
	}

	_ = o.progressTracker.Persist()
	result.Duration = time.Since(startTime)
	return result, nil
}

// GenerateHooks produces contextual sidebar hooks for posts.
//
// Requirements: 12.1 through 12.8, 13.1 through 13.8
func (o *Orchestrator) GenerateHooks(ctx context.Context, opts HookOptions) (*HookResult, error) {
	startTime := time.Now()
	result := &HookResult{}

	postPaths, err := o.postManager.DiscoverPosts(models.PostFilter{
		SpecificSlugs: opts.SpecificSlugs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed discovering posts for hooks: %w", err)
	}

	modelName := opts.Model
	if modelName == "" {
		modelName = o.cfg.LlamaServer.GetHookModel()
	}
	if modelName == "" {
		modelName = "llama-model"
	}

	// Ensure hook generation model is loaded if server is managed
	if o.serverManager != nil && o.cfg != nil && o.cfg.LlamaServer.IsManaged() {
		hookModel := opts.Model
		if hookModel == "" {
			hookModel = o.cfg.LlamaServer.GetHookModel()
		}
		if err := o.serverManager.EnsureModel(ctx, hookModel, server.ModeCompletion); err != nil {
			return nil, fmt.Errorf("failed to load hook generation model: %w", err)
		}
	}

	if o.compressionCache != nil {
		_ = o.compressionCache.Load()
	}

	totalPosts := len(postPaths)
	fmt.Printf("\n🪝 Hook generation: %d posts to process (model: %s)\n", totalPosts, modelName)

	for postIdx, pPath := range postPaths {
		post, readErr := o.postManager.ReadPost(pPath)
		if readErr != nil {
			o.errorCollector.Add(errors.NewFileIOError("read_post", readErr.Error(), readErr, false))
			result.FailedHooks++
			continue
		}

		var relatedArticles []models.RelatedArticle
		if opts.FromStaging {
			staged, sErr := o.stagingManager.ReadStaged(post.Slug)
			if sErr == nil {
				relatedArticles = staged
			}
		}
		if len(relatedArticles) == 0 && post.FrontMatter != nil {
			relatedArticles = post.FrontMatter.RelatedArticles
		}

		if len(relatedArticles) == 0 {
			fmt.Printf("  [%d/%d] %s — skipped (no related articles)\n", postIdx+1, totalPosts, post.Slug)
			continue
		}

		fmt.Printf("  [%d/%d] %s — generating %d hooks...\n", postIdx+1, totalPosts, post.Slug, len(relatedArticles))

		var hookReqs []models.HookRequest
		for _, ra := range relatedArticles {
			// Find widget post
			widgetPostPath := filepath.Join(filepath.Dir(string(pPath)), "..", ra.Slug, "index.md")
			widgetPost, wErr := o.postManager.ReadPost(models.PostPath(widgetPostPath))
			if wErr != nil {
				widgetPost = &models.Post{
					Path: models.PostPath(widgetPostPath),
					Slug: ra.Slug,
					FrontMatter: &models.FrontMatter{
						Title: ra.Title,
						Date:  time.Now(),
					},
					Body: ra.Title,
				}
			}

			hookReqs = append(hookReqs, models.HookRequest{
				MainPost:      post,
				WidgetPost:    widgetPost,
				UseCompressed: opts.UseCompressed,
			})
		}

		hooks, hErr := processor.GenerateHooksBatch(ctx, o.llamaClient, modelName, o.compressionCache, hookReqs)
		if hErr != nil {
			fmt.Printf("    ✗ batch failed: %v\n", hErr)
			o.errorCollector.Add(errors.NewProcessingError("generate_hooks", hErr.Error(), hErr, false))
			result.FailedHooks++
			continue
		}
		fmt.Printf("    ✓ %d/%d hooks generated\n", len(hooks), len(hookReqs))

		if !opts.DryRun {
			meta := models.HookMetadata{
				GeneratedAt:      time.Now().UTC(),
				GeneratorVersion: models.HookGeneratorVersion,
				PostSlug:         post.Slug,
			}
			if post.FrontMatter != nil {
				meta.MainDate = post.FrontMatter.Date.Format("2006-01-02")
				if post.FrontMatter.Date.IsZero() {
					meta.MainDate = ""
				}
			}
			if opts.Stage {
				if sErr := o.stagingManager.SaveStagedHooks(post.Slug, meta, hooks); sErr != nil {
					o.errorCollector.Add(errors.NewFileIOError("save_staged_hooks", sErr.Error(), sErr, false))
					result.FailedHooks++
					continue
				}
			} else {
				postDir := filepath.Dir(string(pPath))
				if sErr := processor.SaveHooks(postDir, meta, hooks); sErr != nil {
					o.errorCollector.Add(errors.NewFileIOError("save_hooks", sErr.Error(), sErr, false))
					result.FailedHooks++
					continue
				}
			}
		}

		result.TotalHooks += len(hooks)
	}

	fmt.Printf("\n✅ Hook generation complete: %d hooks generated, %d failed\n", result.TotalHooks, result.FailedHooks)

	if o.compressionCache != nil && !opts.DryRun {
		_ = o.compressionCache.Persist()
	}

	result.Duration = time.Since(startTime)
	return result, nil
}

// RunPipeline runs both related article generation and hook generation sequentially.
//
// Requirements: All pipeline requirements
func (o *Orchestrator) RunPipeline(ctx context.Context, opts PipelineOptions) (*PipelineResult, error) {
	defer func() {
		if o.serverManager != nil && o.cfg != nil && o.cfg.LlamaServer.IsManaged() {
			_ = o.serverManager.Stop()
		}
	}()

	rankResult, err := o.GenerateRelatedArticles(ctx, opts.RankingOptions)
	if err != nil {
		return nil, fmt.Errorf("ranking stage failed: %w", err)
	}

	hookOpts := opts.HookOptions
	if opts.RankingOptions.Stage {
		hookOpts.FromStaging = true
	}

	hookResult, err := o.GenerateHooks(ctx, hookOpts)
	if err != nil {
		return nil, fmt.Errorf("hook stage failed: %w", err)
	}

	clientStats := models.ClientStats{}
	if o.llamaClient != nil {
		clientStats = o.llamaClient.GetStats()
	}

	summary := errors.PipelineSummary{
		TotalPosts:         rankResult.ProcessedPosts + rankResult.FailedPosts + rankResult.SkippedPosts,
		ProcessedPosts:     rankResult.ProcessedPosts,
		FailedPosts:        rankResult.FailedPosts,
		SkippedPosts:       rankResult.SkippedPosts,
		EmbeddingsComputed: o.embeddingCache.Count(),
		CacheHits:          rankResult.CacheHits,
		TotalAPICalls:      int(clientStats.TotalRequests),
		TotalDuration:      rankResult.Duration + hookResult.Duration,
	}

	return &PipelineResult{
		Ranking: rankResult,
		Hooks:   hookResult,
		Summary: summary,
	}, nil
}

// ApplyStaged applies staged rankings for specified post slugs or all discovered staged files.
func (o *Orchestrator) ApplyStaged(ctx context.Context, postSlugs []string) error {
	if len(postSlugs) == 0 {
		var err error
		postSlugs, err = o.stagingManager.ListStaged()
		if err != nil {
			return err
		}
	}

	for _, slug := range postSlugs {
		if err := o.stagingManager.ApplyStaged(slug); err != nil {
			o.errorCollector.Add(errors.NewFileIOError("apply_staged", err.Error(), err, false))
		}
	}
	return nil
}

// ClearStaged removes staged ranking files.
func (o *Orchestrator) ClearStaged(ctx context.Context, postSlugs []string) error {
	return o.stagingManager.ClearStaged(postSlugs)
}

// ApplyStagedHooks applies staged hooks for specified post slugs or all discovered staged hook files.
func (o *Orchestrator) ApplyStagedHooks(ctx context.Context, postSlugs []string) error {
	if len(postSlugs) == 0 {
		var err error
		postSlugs, err = o.stagingManager.ListStagedHooks()
		if err != nil {
			return err
		}
	}

	for _, slug := range postSlugs {
		if err := o.stagingManager.ApplyStagedHooks(slug); err != nil {
			o.errorCollector.Add(errors.NewFileIOError("apply_staged_hooks", err.Error(), err, false))
		}
	}
	return nil
}

// ClearStagedHooks removes staged hook files.
func (o *Orchestrator) ClearStagedHooks(ctx context.Context, postSlugs []string) error {
	return o.stagingManager.ClearStagedHooks(postSlugs)
}

// ErrorCollector returns the active error collector.
func (o *Orchestrator) ErrorCollector() *errors.ErrorCollector {
	return o.errorCollector
}
