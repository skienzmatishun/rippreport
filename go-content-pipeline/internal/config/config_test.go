package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadConfig_ValidConfig(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	validConfig := `
llama_server:
  base_url: "http://localhost:8080"
  embedding_model: "nomic-embed-text-v1.5"
  llm_model: "qwen3-14b"
  timeout: 120s
  min_request_delay: 100ms
  max_retries: 3

processing:
  posts_directory: "` + tmpDir + `"
  top_candidates: 20
  batch_size: 5
  max_content_length: 100000

storage:
  backup_directory: "backups"
  cache_file: ".embedding_cache.json"
  progress_file: ".progress.json"
  compression_cache: ".compression_cache.json"

scoring:
  weights:
    relevance: 0.7
    recency: 0.15
    length: 0.1
    category: 0.05
  category_boost:
    investigation: 1.2

filtering:
  exclude_categories:
    - holiday
  exclude_titles:
    - "backstory podcast"
  election_cutoff_days: 365

logging:
  level: "INFO"
  file: "pipeline.log"
  console: true

concurrency:
  max_workers: 4
  max_api_requests: 2
  gc_interval_posts: 10
`

	if err := os.WriteFile(configPath, []byte(validConfig), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify values
	if cfg.LlamaServer.BaseURL != "http://localhost:8080" {
		t.Errorf("Expected base_url 'http://localhost:8080', got '%s'", cfg.LlamaServer.BaseURL)
	}
	if cfg.LlamaServer.Timeout != 120*time.Second {
		t.Errorf("Expected timeout 120s, got %v", cfg.LlamaServer.Timeout)
	}
	if cfg.Processing.TopCandidates != 20 {
		t.Errorf("Expected top_candidates 20, got %d", cfg.Processing.TopCandidates)
	}
	if cfg.Scoring.Weights.Relevance != 0.7 {
		t.Errorf("Expected relevance weight 0.7, got %f", cfg.Scoring.Weights.Relevance)
	}
	if len(cfg.Filtering.ExcludeCategories) != 1 {
		t.Errorf("Expected 1 excluded category, got %d", len(cfg.Filtering.ExcludeCategories))
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("Expected error for missing config file, got nil")
	}
	if !strings.Contains(err.Error(), "failed to read config file") {
		t.Errorf("Expected 'failed to read config file' error, got: %v", err)
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	invalidYAML := `
llama_server:
  base_url: "http://localhost:8080"
  invalid yaml structure here [[[
`

	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	_, err := LoadConfig(configPath)
	if err == nil {
		t.Fatal("Expected error for invalid YAML, got nil")
	}
	if !strings.Contains(err.Error(), "failed to parse config file") {
		t.Errorf("Expected 'failed to parse config file' error, got: %v", err)
	}
}

func TestValidate_MissingRequiredFields(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectedErr string
	}{
		{
			name: "missing base_url",
			config: Config{
				LlamaServer: LlamaServerConfig{
					EmbeddingModel:  "model",
					LLMModel:        "model",
					Timeout:         time.Second,
					MinRequestDelay: 0,
					MaxRetries:      1,
				},
			},
			expectedErr: "llama_server.base_url is required",
		},
		{
			name: "missing embedding_model",
			config: Config{
				LlamaServer: LlamaServerConfig{
					BaseURL:         "http://localhost:8080",
					LLMModel:        "model",
					Timeout:         time.Second,
					MinRequestDelay: 0,
					MaxRetries:      1,
				},
			},
			expectedErr: "llama_server.embedding_model is required",
		},
		{
			name: "missing llm_model",
			config: Config{
				LlamaServer: LlamaServerConfig{
					BaseURL:         "http://localhost:8080",
					EmbeddingModel:  "model",
					Timeout:         time.Second,
					MinRequestDelay: 0,
					MaxRetries:      1,
				},
			},
			expectedErr: "llama_server.llm_model is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if err == nil {
				t.Fatal("Expected validation error, got nil")
			}
			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("Expected error containing '%s', got: %v", tt.expectedErr, err)
			}
		})
	}
}

func TestValidate_InvalidValues(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		config      Config
		expectedErr string
	}{
		{
			name: "negative timeout",
			config: Config{
				LlamaServer: LlamaServerConfig{
					BaseURL:        "http://localhost:8080",
					EmbeddingModel: "model",
					LLMModel:       "model",
					Timeout:        -1 * time.Second,
				},
				Processing: ProcessingConfig{
					PostsDirectory: tmpDir,
					TopCandidates:  1,
					BatchSize:      1,
					MaxContentLen:  1,
				},
			},
			expectedErr: "llama_server.timeout must be positive",
		},
		{
			name: "negative min_request_delay",
			config: Config{
				LlamaServer: LlamaServerConfig{
					BaseURL:         "http://localhost:8080",
					EmbeddingModel:  "model",
					LLMModel:        "model",
					Timeout:         time.Second,
					MinRequestDelay: -1 * time.Second,
				},
				Processing: ProcessingConfig{
					PostsDirectory: tmpDir,
					TopCandidates:  1,
					BatchSize:      1,
					MaxContentLen:  1,
				},
			},
			expectedErr: "llama_server.min_request_delay must be non-negative",
		},
		{
			name: "negative max_retries",
			config: Config{
				LlamaServer: LlamaServerConfig{
					BaseURL:         "http://localhost:8080",
					EmbeddingModel:  "model",
					LLMModel:        "model",
					Timeout:         time.Second,
					MinRequestDelay: 0,
					MaxRetries:      -1,
				},
				Processing: ProcessingConfig{
					PostsDirectory: tmpDir,
					TopCandidates:  1,
					BatchSize:      1,
					MaxContentLen:  1,
				},
			},
			expectedErr: "llama_server.max_retries must be non-negative",
		},
		{
			name: "non-positive top_candidates",
			config: Config{
				LlamaServer: LlamaServerConfig{
					BaseURL:        "http://localhost:8080",
					EmbeddingModel: "model",
					LLMModel:       "model",
					Timeout:        time.Second,
					MaxRetries:     0,
				},
				Processing: ProcessingConfig{
					PostsDirectory: tmpDir,
					TopCandidates:  0,
				},
			},
			expectedErr: "processing.top_candidates must be positive",
		},
		{
			name: "non-positive batch_size",
			config: Config{
				LlamaServer: LlamaServerConfig{
					BaseURL:        "http://localhost:8080",
					EmbeddingModel: "model",
					LLMModel:       "model",
					Timeout:        time.Second,
					MaxRetries:     0,
				},
				Processing: ProcessingConfig{
					PostsDirectory: tmpDir,
					TopCandidates:  1,
					BatchSize:      0,
				},
			},
			expectedErr: "processing.batch_size must be positive",
		},
		{
			name: "non-positive max_content_length",
			config: Config{
				LlamaServer: LlamaServerConfig{
					BaseURL:        "http://localhost:8080",
					EmbeddingModel: "model",
					LLMModel:       "model",
					Timeout:        time.Second,
					MaxRetries:     0,
				},
				Processing: ProcessingConfig{
					PostsDirectory: tmpDir,
					TopCandidates:  1,
					BatchSize:      1,
					MaxContentLen:  0,
				},
			},
			expectedErr: "processing.max_content_length must be positive",
		},
		{
			name: "non-positive max_workers",
			config: Config{
				LlamaServer: LlamaServerConfig{
					BaseURL:        "http://localhost:8080",
					EmbeddingModel: "model",
					LLMModel:       "model",
					Timeout:        time.Second,
					MaxRetries:     0,
				},
				Processing: ProcessingConfig{
					PostsDirectory: tmpDir,
					TopCandidates:  1,
					BatchSize:      1,
					MaxContentLen:  1,
				},
				Storage: StorageConfig{
					BackupDir:        "backups",
					CacheFile:        "cache.json",
					ProgressFile:     "progress.json",
					CompressionCache: "compression.json",
				},
				Scoring: ScoringConfig{
					Weights: ScoreWeights{Relevance: 1.0},
				},
				Logging: LoggingConfig{Level: "INFO"},
				Concurrency: ConcurrencyConfig{
					MaxWorkers:     0,
					MaxAPIRequests: 1,
					GCInterval:     1,
				},
			},
			expectedErr: "concurrency.max_workers must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if err == nil {
				t.Fatal("Expected validation error, got nil")
			}
			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("Expected error containing '%s', got: %v", tt.expectedErr, err)
			}
		})
	}
}

func TestValidate_NonExistentPostsDirectory(t *testing.T) {
	cfg := Config{
		LlamaServer: LlamaServerConfig{
			BaseURL:        "http://localhost:8080",
			EmbeddingModel: "model",
			LLMModel:       "model",
			Timeout:        time.Second,
			MaxRetries:     0,
		},
		Processing: ProcessingConfig{
			PostsDirectory: "/nonexistent/directory",
			TopCandidates:  1,
			BatchSize:      1,
			MaxContentLen:  1,
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("Expected validation error for non-existent directory, got nil")
	}
	if !strings.Contains(err.Error(), "posts_directory does not exist") {
		t.Errorf("Expected 'posts_directory does not exist' error, got: %v", err)
	}
}

func TestValidate_InvalidLogLevel(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := Config{
		LlamaServer: LlamaServerConfig{
			BaseURL:        "http://localhost:8080",
			EmbeddingModel: "model",
			LLMModel:       "model",
			Timeout:        time.Second,
			MaxRetries:     0,
		},
		Processing: ProcessingConfig{
			PostsDirectory: tmpDir,
			TopCandidates:  1,
			BatchSize:      1,
			MaxContentLen:  1,
		},
		Storage: StorageConfig{
			BackupDir:        "backups",
			CacheFile:        "cache.json",
			ProgressFile:     "progress.json",
			CompressionCache: "compression.json",
		},
		Scoring: ScoringConfig{
			Weights: ScoreWeights{Relevance: 1.0},
		},
		Logging: LoggingConfig{
			Level: "INVALID",
		},
		Concurrency: ConcurrencyConfig{
			MaxWorkers:     1,
			MaxAPIRequests: 1,
			GCInterval:     1,
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("Expected validation error for invalid log level, got nil")
	}
	if !strings.Contains(err.Error(), "logging.level must be one of") {
		t.Errorf("Expected log level validation error, got: %v", err)
	}
}

func TestValidate_NegativeScoreWeights(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := Config{
		LlamaServer: LlamaServerConfig{
			BaseURL:        "http://localhost:8080",
			EmbeddingModel: "model",
			LLMModel:       "model",
			Timeout:        time.Second,
			MaxRetries:     0,
		},
		Processing: ProcessingConfig{
			PostsDirectory: tmpDir,
			TopCandidates:  1,
			BatchSize:      1,
			MaxContentLen:  1,
		},
		Storage: StorageConfig{
			BackupDir:        "backups",
			CacheFile:        "cache.json",
			ProgressFile:     "progress.json",
			CompressionCache: "compression.json",
		},
		Scoring: ScoringConfig{
			Weights: ScoreWeights{
				Relevance: -0.5,
				Recency:   0.3,
				Length:    0.2,
			},
		},
		Logging: LoggingConfig{Level: "INFO"},
		Concurrency: ConcurrencyConfig{
			MaxWorkers:     1,
			MaxAPIRequests: 1,
			GCInterval:     1,
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("Expected validation error for negative weight, got nil")
	}
	if !strings.Contains(err.Error(), "scoring.weights.relevance must be non-negative") {
		t.Errorf("Expected negative weight error, got: %v", err)
	}
}

func TestValidate_AllZeroWeights(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := Config{
		LlamaServer: LlamaServerConfig{
			BaseURL:        "http://localhost:8080",
			EmbeddingModel: "model",
			LLMModel:       "model",
			Timeout:        time.Second,
			MaxRetries:     0,
		},
		Processing: ProcessingConfig{
			PostsDirectory: tmpDir,
			TopCandidates:  1,
			BatchSize:      1,
			MaxContentLen:  1,
		},
		Storage: StorageConfig{
			BackupDir:        "backups",
			CacheFile:        "cache.json",
			ProgressFile:     "progress.json",
			CompressionCache: "compression.json",
		},
		Scoring: ScoringConfig{
			Weights: ScoreWeights{
				Relevance: 0,
				Recency:   0,
				Length:    0,
				Category:  0,
			},
		},
		Logging: LoggingConfig{Level: "INFO"},
		Concurrency: ConcurrencyConfig{
			MaxWorkers:     1,
			MaxAPIRequests: 1,
			GCInterval:     1,
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("Expected validation error for all-zero weights, got nil")
	}
	if !strings.Contains(err.Error(), "at least one scoring weight must be non-zero") {
		t.Errorf("Expected all-zero weights error, got: %v", err)
	}
}

func TestValidate_NegativeCategoryBoost(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := Config{
		LlamaServer: LlamaServerConfig{
			BaseURL:        "http://localhost:8080",
			EmbeddingModel: "model",
			LLMModel:       "model",
			Timeout:        time.Second,
			MaxRetries:     0,
		},
		Processing: ProcessingConfig{
			PostsDirectory: tmpDir,
			TopCandidates:  1,
			BatchSize:      1,
			MaxContentLen:  1,
		},
		Storage: StorageConfig{
			BackupDir:        "backups",
			CacheFile:        "cache.json",
			ProgressFile:     "progress.json",
			CompressionCache: "compression.json",
		},
		Scoring: ScoringConfig{
			Weights: ScoreWeights{Relevance: 1.0},
			CategoryBoost: map[string]float32{
				"investigation": -0.5,
			},
		},
		Logging: LoggingConfig{Level: "INFO"},
		Concurrency: ConcurrencyConfig{
			MaxWorkers:     1,
			MaxAPIRequests: 1,
			GCInterval:     1,
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("Expected validation error for negative category boost, got nil")
	}
	if !strings.Contains(err.Error(), "scoring.category_boost") {
		t.Errorf("Expected category boost error, got: %v", err)
	}
}

func TestApplyEnvOverrides(t *testing.T) {
	// Save original env vars
	envVars := []string{
		"LLAMA_SERVER_URL", "EMBEDDING_MODEL", "LLM_MODEL",
		"POSTS_DIRECTORY", "LOG_FILE", "LOG_LEVEL",
		"TOP_CANDIDATES", "BATCH_SIZE",
		"BACKUP_DIRECTORY", "CACHE_FILE", "PROGRESS_FILE",
		"MAX_WORKERS", "MAX_API_REQUESTS",
	}
	origValues := make(map[string]string)
	for _, key := range envVars {
		origValues[key] = os.Getenv(key)
	}
	defer func() {
		for key, val := range origValues {
			os.Setenv(key, val)
		}
	}()

	// Set test env vars
	os.Setenv("LLAMA_SERVER_URL", "http://override:9090")
	os.Setenv("EMBEDDING_MODEL", "override-embed-model")
	os.Setenv("LLM_MODEL", "override-llm-model")
	os.Setenv("POSTS_DIRECTORY", "/override/posts")
	os.Setenv("LOG_FILE", "/override/log.txt")
	os.Setenv("LOG_LEVEL", "DEBUG")
	os.Setenv("TOP_CANDIDATES", "15")
	os.Setenv("BATCH_SIZE", "8")
	os.Setenv("BACKUP_DIRECTORY", "/override/backups")
	os.Setenv("CACHE_FILE", "/override/cache.json")
	os.Setenv("PROGRESS_FILE", "/override/progress.json")
	os.Setenv("MAX_WORKERS", "8")
	os.Setenv("MAX_API_REQUESTS", "4")

	cfg := &Config{
		LlamaServer: LlamaServerConfig{
			BaseURL:        "http://localhost:8080",
			EmbeddingModel: "original-embed",
			LLMModel:       "original-llm",
		},
		Processing: ProcessingConfig{
			PostsDirectory: "/original/posts",
			TopCandidates:  10,
			BatchSize:      5,
		},
		Storage: StorageConfig{
			BackupDir:    "/original/backups",
			CacheFile:    "/original/cache.json",
			ProgressFile: "/original/progress.json",
		},
		Logging: LoggingConfig{
			File:  "/original/log.txt",
			Level: "INFO",
		},
		Concurrency: ConcurrencyConfig{
			MaxWorkers:     4,
			MaxAPIRequests: 2,
		},
	}

	cfg.applyEnvOverrides()

	// Verify LlamaServer overrides
	if cfg.LlamaServer.BaseURL != "http://override:9090" {
		t.Errorf("Expected base_url override, got %s", cfg.LlamaServer.BaseURL)
	}
	if cfg.LlamaServer.EmbeddingModel != "override-embed-model" {
		t.Errorf("Expected embedding_model override, got %s", cfg.LlamaServer.EmbeddingModel)
	}
	if cfg.LlamaServer.LLMModel != "override-llm-model" {
		t.Errorf("Expected llm_model override, got %s", cfg.LlamaServer.LLMModel)
	}

	// Verify Processing overrides
	if cfg.Processing.PostsDirectory != "/override/posts" {
		t.Errorf("Expected posts_directory override, got %s", cfg.Processing.PostsDirectory)
	}
	if cfg.Processing.TopCandidates != 15 {
		t.Errorf("Expected top_candidates override to 15, got %d", cfg.Processing.TopCandidates)
	}
	if cfg.Processing.BatchSize != 8 {
		t.Errorf("Expected batch_size override to 8, got %d", cfg.Processing.BatchSize)
	}

	// Verify Storage overrides
	if cfg.Storage.BackupDir != "/override/backups" {
		t.Errorf("Expected backup_directory override, got %s", cfg.Storage.BackupDir)
	}
	if cfg.Storage.CacheFile != "/override/cache.json" {
		t.Errorf("Expected cache_file override, got %s", cfg.Storage.CacheFile)
	}
	if cfg.Storage.ProgressFile != "/override/progress.json" {
		t.Errorf("Expected progress_file override, got %s", cfg.Storage.ProgressFile)
	}

	// Verify Logging overrides
	if cfg.Logging.File != "/override/log.txt" {
		t.Errorf("Expected log file override, got %s", cfg.Logging.File)
	}
	if cfg.Logging.Level != "DEBUG" {
		t.Errorf("Expected log level override to DEBUG, got %s", cfg.Logging.Level)
	}

	// Verify Concurrency overrides
	if cfg.Concurrency.MaxWorkers != 8 {
		t.Errorf("Expected max_workers override to 8, got %d", cfg.Concurrency.MaxWorkers)
	}
	if cfg.Concurrency.MaxAPIRequests != 4 {
		t.Errorf("Expected max_api_requests override to 4, got %d", cfg.Concurrency.MaxAPIRequests)
	}
}

func TestApplyEnvOverrides_InvalidIntegers(t *testing.T) {
	// Save and defer restore
	origTopCandidates := os.Getenv("TOP_CANDIDATES")
	defer os.Setenv("TOP_CANDIDATES", origTopCandidates)

	// Set invalid integer value
	os.Setenv("TOP_CANDIDATES", "not-a-number")

	cfg := &Config{
		Processing: ProcessingConfig{
			TopCandidates: 10,
		},
	}

	cfg.applyEnvOverrides()

	// Should remain unchanged when parse fails
	if cfg.Processing.TopCandidates != 10 {
		t.Errorf("Expected top_candidates to remain 10 after invalid override, got %d", cfg.Processing.TopCandidates)
	}
}

func TestApplyEnvOverrides_NegativeIntegers(t *testing.T) {
	// Save and defer restore
	origBatchSize := os.Getenv("BATCH_SIZE")
	defer os.Setenv("BATCH_SIZE", origBatchSize)

	// Set negative value
	os.Setenv("BATCH_SIZE", "-5")

	cfg := &Config{
		Processing: ProcessingConfig{
			BatchSize: 5,
		},
	}

	cfg.applyEnvOverrides()

	// Should remain unchanged when value is not positive
	if cfg.Processing.BatchSize != 5 {
		t.Errorf("Expected batch_size to remain 5 after negative override, got %d", cfg.Processing.BatchSize)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.LlamaServer.BaseURL != "http://localhost:8080" {
		t.Errorf("Unexpected default base_url: %s", cfg.LlamaServer.BaseURL)
	}
	if cfg.LlamaServer.Timeout != 120*time.Second {
		t.Errorf("Unexpected default timeout: %v", cfg.LlamaServer.Timeout)
	}
	if cfg.Processing.TopCandidates != 20 {
		t.Errorf("Unexpected default top_candidates: %d", cfg.Processing.TopCandidates)
	}
	if cfg.Scoring.Weights.Relevance != 0.7 {
		t.Errorf("Unexpected default relevance weight: %f", cfg.Scoring.Weights.Relevance)
	}
	if cfg.Concurrency.MaxWorkers != 4 {
		t.Errorf("Unexpected default max_workers: %d", cfg.Concurrency.MaxWorkers)
	}
}

func TestMakeAbsolutePaths(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &Config{
		Processing: ProcessingConfig{
			PostsDirectory: "content/p",
		},
		Storage: StorageConfig{
			BackupDir:        "backups",
			CacheFile:        ".cache.json",
			ProgressFile:     ".progress.json",
			CompressionCache: ".compression.json",
		},
		Logging: LoggingConfig{
			File: "pipeline.log",
		},
	}

	err := cfg.MakeAbsolutePaths(tmpDir)
	if err != nil {
		t.Fatalf("MakeAbsolutePaths failed: %v", err)
	}

	// Check that paths are now absolute
	if !filepath.IsAbs(cfg.Processing.PostsDirectory) {
		t.Error("PostsDirectory should be absolute")
	}
	if !filepath.IsAbs(cfg.Storage.BackupDir) {
		t.Error("BackupDir should be absolute")
	}
	if !filepath.IsAbs(cfg.Storage.CacheFile) {
		t.Error("CacheFile should be absolute")
	}
	if !filepath.IsAbs(cfg.Storage.ProgressFile) {
		t.Error("ProgressFile should be absolute")
	}
	if !filepath.IsAbs(cfg.Storage.CompressionCache) {
		t.Error("CompressionCache should be absolute")
	}
	if !filepath.IsAbs(cfg.Logging.File) {
		t.Error("Log file should be absolute")
	}

	// Check that paths are based on tmpDir
	if !strings.HasPrefix(cfg.Processing.PostsDirectory, tmpDir) {
		t.Errorf("PostsDirectory should be under %s, got %s", tmpDir, cfg.Processing.PostsDirectory)
	}
}

func TestMakeAbsolutePaths_AlreadyAbsolute(t *testing.T) {
	cfg := &Config{
		Processing: ProcessingConfig{
			PostsDirectory: "/absolute/path/content/p",
		},
		Storage: StorageConfig{
			BackupDir:        "/absolute/path/backups",
			CacheFile:        "/absolute/path/.cache.json",
			ProgressFile:     "/absolute/path/.progress.json",
			CompressionCache: "/absolute/path/.compression.json",
		},
		Logging: LoggingConfig{
			File: "/absolute/path/pipeline.log",
		},
	}

	originalPaths := map[string]string{
		"posts":       cfg.Processing.PostsDirectory,
		"backup":      cfg.Storage.BackupDir,
		"cache":       cfg.Storage.CacheFile,
		"progress":    cfg.Storage.ProgressFile,
		"compression": cfg.Storage.CompressionCache,
		"log":         cfg.Logging.File,
	}

	err := cfg.MakeAbsolutePaths("/some/base/dir")
	if err != nil {
		t.Fatalf("MakeAbsolutePaths failed: %v", err)
	}

	// Paths should be unchanged since they were already absolute
	if cfg.Processing.PostsDirectory != originalPaths["posts"] {
		t.Error("Absolute PostsDirectory was modified")
	}
	if cfg.Storage.BackupDir != originalPaths["backup"] {
		t.Error("Absolute BackupDir was modified")
	}
	if cfg.Storage.CacheFile != originalPaths["cache"] {
		t.Error("Absolute CacheFile was modified")
	}
	if cfg.Storage.ProgressFile != originalPaths["progress"] {
		t.Error("Absolute ProgressFile was modified")
	}
	if cfg.Storage.CompressionCache != originalPaths["compression"] {
		t.Error("Absolute CompressionCache was modified")
	}
	if cfg.Logging.File != originalPaths["log"] {
		t.Error("Absolute log file was modified")
	}
}
