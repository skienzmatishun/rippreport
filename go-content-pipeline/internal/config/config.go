// Package config provides configuration management for the go-content-pipeline.
//
// The configuration system supports loading settings from YAML files with
// environment variable overrides for sensitive values. All configuration
// is validated on load to catch errors early.
//
// # Configuration Structure
//
// The configuration is organized into logical sections:
//   - LlamaServer: Connection settings for llama.cpp server
//   - Processing: Post processing behavior
//   - Storage: File paths for caches and backups
//   - Scoring: Algorithm weights and boosts
//   - Filtering: Post filtering criteria
//   - Logging: Log output settings
//   - Concurrency: Worker pool and rate limiting
//
// # Example Usage
//
// Loading configuration from a file:
//
//	cfg, err := config.LoadConfig("config.yaml")
//	if err != nil {
//	    log.Fatalf("Failed to load config: %v", err)
//	}
//
// Using default configuration:
//
//	cfg := config.DefaultConfig()
//	cfg.LlamaServer.BaseURL = "http://remote-server:8080"
//	if err := cfg.Validate(); err != nil {
//	    log.Fatalf("Invalid config: %v", err)
//	}
//
// Converting relative paths to absolute:
//
//	cfg.MakeAbsolutePaths("/path/to/project/root")
//
// # Environment Variable Overrides
//
// The following environment variables can override config file values:
//
// Llama Server:
//   - LLAMA_SERVER_URL: Override base_url
//   - EMBEDDING_MODEL: Override embedding_model
//   - LLM_MODEL: Override llm_model
//
// Processing:
//   - POSTS_DIRECTORY: Override posts_directory
//   - TOP_CANDIDATES: Override top_candidates
//   - BATCH_SIZE: Override batch_size
//
// Storage:
//   - BACKUP_DIRECTORY: Override backup_directory
//   - CACHE_FILE: Override cache_file
//   - PROGRESS_FILE: Override progress_file
//
// Logging:
//   - LOG_LEVEL: Override log level (DEBUG, INFO, WARN, ERROR)
//   - LOG_FILE: Override log file path
//
// Concurrency:
//   - MAX_WORKERS: Override max_workers
//   - MAX_API_REQUESTS: Override max_api_requests
//
// Environment variables are applied after loading the config file,
// allowing runtime customization without modifying config files.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the complete application configuration
type Config struct {
	LlamaServer LlamaServerConfig `yaml:"llama_server"`
	Processing  ProcessingConfig  `yaml:"processing"`
	Storage     StorageConfig     `yaml:"storage"`
	Scoring     ScoringConfig     `yaml:"scoring"`
	Filtering   FilteringConfig   `yaml:"filtering"`
	Logging     LoggingConfig     `yaml:"logging"`
	Concurrency ConcurrencyConfig `yaml:"concurrency"`
}

// LlamaServerConfig contains llama.cpp server connection settings
type LlamaServerConfig struct {
	BaseURL         string        `yaml:"base_url"`
	EmbeddingModel  string        `yaml:"embedding_model"`
	LLMModel        string        `yaml:"llm_model"`
	Timeout         time.Duration `yaml:"timeout"`
	MinRequestDelay time.Duration `yaml:"min_request_delay"`
	MaxRetries      int           `yaml:"max_retries"`
}

// ProcessingConfig contains post processing settings
type ProcessingConfig struct {
	PostsDirectory string `yaml:"posts_directory"`
	TopCandidates  int    `yaml:"top_candidates"`
	BatchSize      int    `yaml:"batch_size"`
	MaxContentLen  int    `yaml:"max_content_length"`
}

// StorageConfig contains file storage paths
type StorageConfig struct {
	BackupDir        string `yaml:"backup_directory"`
	CacheFile        string `yaml:"cache_file"`
	ProgressFile     string `yaml:"progress_file"`
	CompressionCache string `yaml:"compression_cache"`
}

// ScoringConfig contains scoring algorithm settings
type ScoringConfig struct {
	Weights       ScoreWeights       `yaml:"weights"`
	CategoryBoost map[string]float32 `yaml:"category_boost"`
}

// ScoreWeights defines the weights for composite scoring
type ScoreWeights struct {
	Relevance float32 `yaml:"relevance"`
	Recency   float32 `yaml:"recency"`
	Length    float32 `yaml:"length"`
	Category  float32 `yaml:"category"`
}

// FilteringConfig contains post filtering criteria
type FilteringConfig struct {
	ExcludeCategories  []string `yaml:"exclude_categories"`
	ExcludeTitles      []string `yaml:"exclude_titles"`
	ElectionCutoffDays int      `yaml:"election_cutoff_days"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level   string `yaml:"level"`
	File    string `yaml:"file"`
	Console bool   `yaml:"console"`
}

// ConcurrencyConfig contains concurrency limits
type ConcurrencyConfig struct {
	MaxWorkers     int `yaml:"max_workers"`
	MaxAPIRequests int `yaml:"max_api_requests"`
	GCInterval     int `yaml:"gc_interval_posts"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	// Parse YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	// Apply environment variable overrides
	cfg.applyEnvOverrides()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// applyEnvOverrides applies environment variable overrides to sensitive values
func (c *Config) applyEnvOverrides() {
	// Llama server configuration overrides
	if url := os.Getenv("LLAMA_SERVER_URL"); url != "" {
		c.LlamaServer.BaseURL = url
	}
	if model := os.Getenv("EMBEDDING_MODEL"); model != "" {
		c.LlamaServer.EmbeddingModel = model
	}
	if model := os.Getenv("LLM_MODEL"); model != "" {
		c.LlamaServer.LLMModel = model
	}

	// Processing configuration overrides
	if dir := os.Getenv("POSTS_DIRECTORY"); dir != "" {
		c.Processing.PostsDirectory = dir
	}
	if val := os.Getenv("TOP_CANDIDATES"); val != "" {
		if n, err := parseIntEnv(val); err == nil && n > 0 {
			c.Processing.TopCandidates = n
		}
	}
	if val := os.Getenv("BATCH_SIZE"); val != "" {
		if n, err := parseIntEnv(val); err == nil && n > 0 {
			c.Processing.BatchSize = n
		}
	}

	// Storage configuration overrides
	if dir := os.Getenv("BACKUP_DIRECTORY"); dir != "" {
		c.Storage.BackupDir = dir
	}
	if file := os.Getenv("CACHE_FILE"); file != "" {
		c.Storage.CacheFile = file
	}
	if file := os.Getenv("PROGRESS_FILE"); file != "" {
		c.Storage.ProgressFile = file
	}

	// Logging configuration overrides
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		c.Logging.Level = level
	}
	if file := os.Getenv("LOG_FILE"); file != "" {
		c.Logging.File = file
	}

	// Concurrency configuration overrides
	if val := os.Getenv("MAX_WORKERS"); val != "" {
		if n, err := parseIntEnv(val); err == nil && n > 0 {
			c.Concurrency.MaxWorkers = n
		}
	}
	if val := os.Getenv("MAX_API_REQUESTS"); val != "" {
		if n, err := parseIntEnv(val); err == nil && n > 0 {
			c.Concurrency.MaxAPIRequests = n
		}
	}
}

// parseIntEnv is a helper function to parse integer environment variables
func parseIntEnv(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}

// Validate checks that all required fields are present and valid
func (c *Config) Validate() error {
	// Validate LlamaServer config
	if c.LlamaServer.BaseURL == "" {
		return fmt.Errorf("llama_server.base_url is required")
	}
	if c.LlamaServer.EmbeddingModel == "" {
		return fmt.Errorf("llama_server.embedding_model is required")
	}
	if c.LlamaServer.LLMModel == "" {
		return fmt.Errorf("llama_server.llm_model is required")
	}
	if c.LlamaServer.Timeout <= 0 {
		return fmt.Errorf("llama_server.timeout must be positive, got %v", c.LlamaServer.Timeout)
	}
	if c.LlamaServer.MinRequestDelay < 0 {
		return fmt.Errorf("llama_server.min_request_delay must be non-negative, got %v", c.LlamaServer.MinRequestDelay)
	}
	if c.LlamaServer.MaxRetries < 0 {
		return fmt.Errorf("llama_server.max_retries must be non-negative, got %d", c.LlamaServer.MaxRetries)
	}

	// Validate Processing config
	if c.Processing.PostsDirectory == "" {
		return fmt.Errorf("processing.posts_directory is required")
	}
	// Check if posts directory exists
	if _, err := os.Stat(c.Processing.PostsDirectory); os.IsNotExist(err) {
		return fmt.Errorf("processing.posts_directory does not exist: %s", c.Processing.PostsDirectory)
	}
	if c.Processing.TopCandidates <= 0 {
		return fmt.Errorf("processing.top_candidates must be positive, got %d", c.Processing.TopCandidates)
	}
	if c.Processing.BatchSize <= 0 {
		return fmt.Errorf("processing.batch_size must be positive, got %d", c.Processing.BatchSize)
	}
	if c.Processing.MaxContentLen <= 0 {
		return fmt.Errorf("processing.max_content_length must be positive, got %d", c.Processing.MaxContentLen)
	}

	// Validate Storage config
	if c.Storage.BackupDir == "" {
		return fmt.Errorf("storage.backup_directory is required")
	}
	if c.Storage.CacheFile == "" {
		return fmt.Errorf("storage.cache_file is required")
	}
	if c.Storage.ProgressFile == "" {
		return fmt.Errorf("storage.progress_file is required")
	}
	if c.Storage.CompressionCache == "" {
		return fmt.Errorf("storage.compression_cache is required")
	}

	// Validate Scoring config
	if err := c.validateScoreWeights(); err != nil {
		return err
	}

	// Validate Logging config
	if c.Logging.Level == "" {
		return fmt.Errorf("logging.level is required")
	}
	validLevels := map[string]bool{"DEBUG": true, "INFO": true, "WARN": true, "ERROR": true}
	if !validLevels[c.Logging.Level] {
		return fmt.Errorf("logging.level must be one of DEBUG, INFO, WARN, ERROR, got %s", c.Logging.Level)
	}

	// Validate Concurrency config
	if c.Concurrency.MaxWorkers <= 0 {
		return fmt.Errorf("concurrency.max_workers must be positive, got %d", c.Concurrency.MaxWorkers)
	}
	if c.Concurrency.MaxAPIRequests <= 0 {
		return fmt.Errorf("concurrency.max_api_requests must be positive, got %d", c.Concurrency.MaxAPIRequests)
	}
	if c.Concurrency.GCInterval <= 0 {
		return fmt.Errorf("concurrency.gc_interval_posts must be positive, got %d", c.Concurrency.GCInterval)
	}

	return nil
}

// validateScoreWeights checks that scoring weights are valid
func (c *Config) validateScoreWeights() error {
	w := c.Scoring.Weights

	if w.Relevance < 0 {
		return fmt.Errorf("scoring.weights.relevance must be non-negative, got %f", w.Relevance)
	}
	if w.Recency < 0 {
		return fmt.Errorf("scoring.weights.recency must be non-negative, got %f", w.Recency)
	}
	if w.Length < 0 {
		return fmt.Errorf("scoring.weights.length must be non-negative, got %f", w.Length)
	}
	if w.Category < 0 {
		return fmt.Errorf("scoring.weights.category must be non-negative, got %f", w.Category)
	}

	// Check that at least one weight is non-zero
	if w.Relevance == 0 && w.Recency == 0 && w.Length == 0 && w.Category == 0 {
		return fmt.Errorf("at least one scoring weight must be non-zero")
	}

	// Validate category boost values
	for category, boost := range c.Scoring.CategoryBoost {
		if boost < 0 {
			return fmt.Errorf("scoring.category_boost.%s must be non-negative, got %f", category, boost)
		}
	}

	return nil
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		LlamaServer: LlamaServerConfig{
			BaseURL:         "http://localhost:8080",
			EmbeddingModel:  "nomic-embed-text-v1.5",
			LLMModel:        "qwen3-14b",
			Timeout:         120 * time.Second,
			MinRequestDelay: 100 * time.Millisecond,
			MaxRetries:      3,
		},
		Processing: ProcessingConfig{
			PostsDirectory: "content/p",
			TopCandidates:  20,
			BatchSize:      5,
			MaxContentLen:  100000,
		},
		Storage: StorageConfig{
			BackupDir:        "backups",
			CacheFile:        ".embedding_cache.json",
			ProgressFile:     ".progress.json",
			CompressionCache: ".compression_cache.json",
		},
		Scoring: ScoringConfig{
			Weights: ScoreWeights{
				Relevance: 0.7,
				Recency:   0.15,
				Length:    0.1,
				Category:  0.05,
			},
			CategoryBoost: map[string]float32{
				"investigation": 1.2,
				"analysis":      1.1,
			},
		},
		Filtering: FilteringConfig{
			ExcludeCategories: []string{"holiday"},
			ExcludeTitles: []string{
				"backstory podcast",
				"wonderful wednesday",
				"freaky friday",
			},
			ElectionCutoffDays: 365,
		},
		Logging: LoggingConfig{
			Level:   "INFO",
			File:    "pipeline.log",
			Console: true,
		},
		Concurrency: ConcurrencyConfig{
			MaxWorkers:     4,
			MaxAPIRequests: 2,
			GCInterval:     10,
		},
	}
}

// MakeAbsolutePaths converts relative paths in the config to absolute paths
// based on the given base directory (typically the project root)
func (c *Config) MakeAbsolutePaths(baseDir string) error {
	// Convert posts directory to absolute path if relative
	if !filepath.IsAbs(c.Processing.PostsDirectory) {
		c.Processing.PostsDirectory = filepath.Join(baseDir, c.Processing.PostsDirectory)
	}

	// Convert storage paths to absolute paths if relative
	if !filepath.IsAbs(c.Storage.BackupDir) {
		c.Storage.BackupDir = filepath.Join(baseDir, c.Storage.BackupDir)
	}
	if !filepath.IsAbs(c.Storage.CacheFile) {
		c.Storage.CacheFile = filepath.Join(baseDir, c.Storage.CacheFile)
	}
	if !filepath.IsAbs(c.Storage.ProgressFile) {
		c.Storage.ProgressFile = filepath.Join(baseDir, c.Storage.ProgressFile)
	}
	if !filepath.IsAbs(c.Storage.CompressionCache) {
		c.Storage.CompressionCache = filepath.Join(baseDir, c.Storage.CompressionCache)
	}
	if !filepath.IsAbs(c.Logging.File) {
		c.Logging.File = filepath.Join(baseDir, c.Logging.File)
	}

	return nil
}
