package test

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/cache"
	"github.com/rippreport/go-content-pipeline/internal/client"
	"github.com/rippreport/go-content-pipeline/internal/config"
	"github.com/rippreport/go-content-pipeline/internal/orchestrator"
	"github.com/rippreport/go-content-pipeline/internal/processor"
	"github.com/rippreport/go-content-pipeline/internal/testutil"
)

// TestEnv provides a fully assembled test harness with temporary folders,
// mock llama server, and all pipeline components wired together.
type TestEnv struct {
	T                *testing.T
	RootDir          string
	ContentDir       string
	PostsDir         string
	CacheDir         string
	BackupsDir       string
	MockServer       *testutil.MockLlamaServer
	Config           *config.Config
	LlamaClient      client.LlamaClient
	PostManager      *processor.PostManager
	EmbeddingCache   *cache.EmbeddingCache
	ProgressTracker  *processor.ProgressTracker
	BackupManager    *processor.BackupManager
	StagingManager   *processor.StagingManager
	CompressionCache *processor.CompressionCache
	Orchestrator     *orchestrator.Orchestrator

	mu sync.Mutex
}

// SetupTestEnv creates and initializes a complete test environment.
func SetupTestEnv(t *testing.T) *TestEnv {
	t.Helper()

	rootDir := testutil.TempDir(t)
	contentDir := filepath.Join(rootDir, "content")
	postsDir := filepath.Join(contentDir, "p")
	cacheDir := filepath.Join(rootDir, "cache")
	backupsDir := filepath.Join(rootDir, "backups")

	mockServer := testutil.NewMockLlamaServer(t)

	// Configure mock server with intelligent deterministic handlers
	mockServer.EmbeddingResponse = func(req testutil.EmbeddingRequest) ([]float32, error) {
		dim := 64
		// Generate embedding based on content hash/length to give varied similarity
		seed := float32(len(req.Content)%100) / 100.0
		return testutil.TestEmbedding(dim, seed), nil
	}

	mockServer.CompletionResponse = func(req testutil.CompletionRequest) (string, error) {
		prompt := req.Prompt
		if strings.Contains(prompt, "THE CONFLICT") || strings.Contains(prompt, "Article 1") {
			return "THE CONFLICT:\n- High stakes county zoning and coastal planning battle.\n\nTHE DISCOVERY:\n- Secret financial audits revealed conflicts of interest.\n\nTHE PIVOT:\n- The investigation uncovers deep ties to local officials.\n\nKEY ENTITIES:\n- County board, developers", nil
		}
		// Check if this is a scoring request
		if strings.Contains(prompt, "Relevance Score:") || strings.Contains(prompt, "score") {
			return "85", nil
		}
		// Check if this is a hook batch request
		if strings.Contains(prompt, "JSON") || strings.Contains(prompt, "hook") {
			return `[{"widget_slug":"test-widget","widget_title":"Test Title","brief":"Insightful hook text linking the two articles together."}]`, nil
		}
		return "THE CONFLICT:\n- High stakes county zoning and coastal planning battle.\n\nTHE DISCOVERY:\n- Secret financial audits revealed conflicts of interest.\n\nTHE PIVOT:\n- The investigation uncovers deep ties to local officials.", nil
	}

	cfg := config.DefaultConfig()
	cfg.LlamaServer.BaseURL = mockServer.Server.URL
	cfg.LlamaServer.MinRequestDelay = 1 * time.Millisecond // Fast for tests
	cfg.LlamaServer.EmbeddingModel = "test-embed-model"
	cfg.LlamaServer.LLMModel = "test-llm-model"

	cfg.Processing.PostsDirectory = postsDir
	cfg.Processing.TopCandidates = 3
	cfg.Processing.MaxContentLen = 4000

	cfg.Storage.BackupDir = backupsDir
	cfg.Storage.CacheFile = filepath.Join(cacheDir, "embedding_cache.json")
	cfg.Storage.ProgressFile = filepath.Join(cacheDir, "progress.json")
	cfg.Storage.CompressionCache = filepath.Join(cacheDir, "compression_cache.json")

	llamaClient := client.NewClient(cfg.LlamaServer)
	postManager := processor.NewPostManager(postsDir)
	embeddingCache := cache.NewEmbeddingCache(cfg.Storage.CacheFile, cfg.LlamaServer.EmbeddingModel)
	progressTracker := processor.NewProgressTracker(cfg.Storage.ProgressFile)
	backupManager := processor.NewBackupManager(cfg.Storage.BackupDir)
	stagingManager := processor.NewStagingManager(postsDir, postManager, backupManager)
	compressionCache := processor.NewCompressionCache(cfg.Storage.CompressionCache, cfg.LlamaServer.LLMModel)

	orch := orchestrator.NewOrchestrator(
		cfg,
		llamaClient,
		postManager,
		embeddingCache,
		progressTracker,
		backupManager,
		stagingManager,
		compressionCache,
	)

	return &TestEnv{
		T:                t,
		RootDir:          rootDir,
		ContentDir:       contentDir,
		PostsDir:         postsDir,
		CacheDir:         cacheDir,
		BackupsDir:       backupsDir,
		MockServer:       mockServer,
		Config:           cfg,
		LlamaClient:      llamaClient,
		PostManager:      postManager,
		EmbeddingCache:   embeddingCache,
		ProgressTracker:  progressTracker,
		BackupManager:    backupManager,
		StagingManager:   stagingManager,
		CompressionCache: compressionCache,
		Orchestrator:     orch,
	}
}

// CreatePost creates a single post in the test environment.
func (env *TestEnv) CreatePost(post *testutil.HugoPostFixture) string {
	env.T.Helper()
	return post.CreatePost(env.T, env.ContentDir)
}

// CreateSamplePosts creates multiple standard posts for pipeline testing.
func (env *TestEnv) CreateSamplePosts(count int) []*testutil.HugoPostFixture {
	env.T.Helper()
	posts := make([]*testutil.HugoPostFixture, count)
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	for i := 0; i < count; i++ {
		slug := "sample-post-" + strconv.Itoa(i+1)
		p := testutil.DefaultPost(slug)
		p.Title = fmt.Sprintf("Sample Post Number %d", i+1)
		p.Date = baseTime.Add(time.Duration(i*24) * time.Hour)
		p.Content = fmt.Sprintf("Article %d discusses topic number %d in depth with interesting observations and facts.\n\nParagraph two adds more substantive details about subject %d.", i+1, i+1, i+1)
		env.CreatePost(p)
		posts[i] = p
	}
	return posts
}
