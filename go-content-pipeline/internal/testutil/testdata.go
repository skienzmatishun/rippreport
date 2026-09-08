package testutil

import (
	"encoding/json"
	"math"
	"testing"
)

// TestEmbedding generates a simple test embedding vector.
// The vector is normalized to unit length.
func TestEmbedding(dim int, seed float32) []float32 {
	vec := make([]float32, dim)
	for i := 0; i < dim; i++ {
		vec[i] = seed + float32(i)*0.01
	}
	return NormalizeVector(vec)
}

// IdenticalEmbedding returns an embedding identical to the given one.
func IdenticalEmbedding(embedding []float32) []float32 {
	result := make([]float32, len(embedding))
	copy(result, embedding)
	return result
}

// SimilarEmbedding returns an embedding similar to the given one (high cosine similarity).
func SimilarEmbedding(embedding []float32, similarity float32) []float32 {
	// Create a vector that will have the desired cosine similarity
	// by mixing the original with a small orthogonal component
	dim := len(embedding)
	result := make([]float32, dim)
	
	// Start with the original vector scaled by the similarity
	for i := 0; i < dim; i++ {
		result[i] = embedding[i] * similarity
	}
	
	// Add orthogonal component
	orthScale := float32(math.Sqrt(float64(1 - similarity*similarity)))
	for i := 0; i < dim; i++ {
		if i%2 == 0 {
			result[i] += orthScale * 0.1 * float32(i+1)
		}
	}
	
	return NormalizeVector(result)
}

// DifferentEmbedding returns an embedding with low similarity to the given one.
func DifferentEmbedding(embedding []float32) []float32 {
	dim := len(embedding)
	result := make([]float32, dim)
	
	// Create a roughly orthogonal vector
	for i := 0; i < dim; i++ {
		if i < dim/2 {
			result[i] = -embedding[dim-1-i]
		} else {
			result[i] = embedding[dim-1-i] * 0.5
		}
	}
	
	return NormalizeVector(result)
}

// NormalizeVector normalizes a vector to unit length.
func NormalizeVector(vec []float32) []float32 {
	var magnitude float32
	for _, v := range vec {
		magnitude += v * v
	}
	magnitude = float32(math.Sqrt(float64(magnitude)))
	
	if magnitude == 0 {
		return vec
	}
	
	normalized := make([]float32, len(vec))
	for i, v := range vec {
		normalized[i] = v / magnitude
	}
	return normalized
}

// DotProduct calculates the dot product of two vectors.
func DotProduct(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}
	var sum float32
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

// EmbeddingCacheJSON generates JSON for an embedding cache file.
func EmbeddingCacheJSON(entries map[string]EmbeddingCacheEntry) string {
	cache := map[string]interface{}{
		"version":    "1.0",
		"model_name": "test-embedding-model",
		"updated_at": "2024-01-15T10:30:00Z",
		"entries":    entries,
	}
	data, _ := json.MarshalIndent(cache, "", "  ")
	return string(data)
}

// EmbeddingCacheEntry represents a cache entry for testing.
type EmbeddingCacheEntry struct {
	Embedding   []float32 `json:"embedding"`
	ContentHash string    `json:"content_hash"`
	ModelName   string    `json:"model_name"`
	Timestamp   string    `json:"timestamp"`
}

// NewEmbeddingCacheEntry creates a test cache entry.
func NewEmbeddingCacheEntry(embedding []float32, contentHash string) EmbeddingCacheEntry {
	return EmbeddingCacheEntry{
		Embedding:   embedding,
		ContentHash: contentHash,
		ModelName:   "test-embedding-model",
		Timestamp:   "2024-01-15T10:30:00Z",
	}
}

// ProgressFileJSON generates JSON for a progress tracker file.
func ProgressFileJSON(processed map[string]ProcessedPost) string {
	progress := map[string]interface{}{
		"version":    "1.0",
		"started_at": "2024-01-15T10:00:00Z",
		"updated_at": "2024-01-15T10:30:00Z",
		"processed":  processed,
	}
	data, _ := json.MarshalIndent(progress, "", "  ")
	return string(data)
}

// ProcessedPost represents a processed post for testing.
type ProcessedPost struct {
	Slug      string `json:"slug"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Error     string `json:"error,omitempty"`
}

// NewProcessedPost creates a test processed post entry.
func NewProcessedPost(slug, status string, errorMsg ...string) ProcessedPost {
	entry := ProcessedPost{
		Slug:      slug,
		Status:    status,
		Timestamp: "2024-01-15T10:30:00Z",
	}
	if len(errorMsg) > 0 && errorMsg[0] != "" {
		entry.Error = errorMsg[0]
	}
	return entry
}

// CompressionCacheJSON generates JSON for an article compression cache file.
func CompressionCacheJSON(entries map[string]CompressionCacheEntry) string {
	cache := map[string]interface{}{
		"version":    "1.0",
		"updated_at": "2024-01-15T10:30:00Z",
		"entries":    entries,
	}
	data, _ := json.MarshalIndent(cache, "", "  ")
	return string(data)
}

// CompressionCacheEntry represents a compression cache entry for testing.
type CompressionCacheEntry struct {
	OriginalHash  string  `json:"original_hash"`
	Compressed    string  `json:"compressed"`
	OriginalLen   int     `json:"original_len"`
	CompressedLen int     `json:"compressed_len"`
	Ratio         float32 `json:"ratio"`
	Timestamp     string  `json:"timestamp"`
	ModelName     string  `json:"model_name"`
}

// NewCompressionCacheEntry creates a test compression cache entry.
func NewCompressionCacheEntry(originalHash, compressed string, originalLen, compressedLen int) CompressionCacheEntry {
	return CompressionCacheEntry{
		OriginalHash:  originalHash,
		Compressed:    compressed,
		OriginalLen:   originalLen,
		CompressedLen: compressedLen,
		Ratio:         float32(compressedLen) / float32(originalLen),
		Timestamp:     "2024-01-15T10:30:00Z",
		ModelName:     "test-llm-model",
	}
}

// SampleScores returns a set of sample LLM scores for testing.
type SampleScores struct {
	HighRelevance   int // 90-100
	MediumRelevance int // 60-80
	LowRelevance    int // 20-40
	Irrelevant      int // 0-10
}

// DefaultScores returns default test scores.
func DefaultScores() SampleScores {
	return SampleScores{
		HighRelevance:   95,
		MediumRelevance: 72,
		LowRelevance:    35,
		Irrelevant:      5,
	}
}

// ScoreFactorsData represents sample data for composite score calculation testing.
type ScoreFactorsData struct {
	LLMScore       int
	VectorScore    float32
	RecencyDays    int
	ContentLength  int
	CategoryMatch  bool
	TitleHasVote   bool
	IsElection     bool
}

// HighScoreFactors returns factors that should produce a high composite score.
func HighScoreFactors() ScoreFactorsData {
	return ScoreFactorsData{
		LLMScore:      95,
		VectorScore:   0.92,
		RecencyDays:   30,
		ContentLength: 800,
		CategoryMatch: true,
		TitleHasVote:  false,
		IsElection:    false,
	}
}

// MediumScoreFactors returns factors that should produce a medium composite score.
func MediumScoreFactors() ScoreFactorsData {
	return ScoreFactorsData{
		LLMScore:      70,
		VectorScore:   0.75,
		RecencyDays:   180,
		ContentLength: 450,
		CategoryMatch: false,
		TitleHasVote:  false,
		IsElection:    false,
	}
}

// LowScoreFactors returns factors that should produce a low composite score.
func LowScoreFactors() ScoreFactorsData {
	return ScoreFactorsData{
		LLMScore:      40,
		VectorScore:   0.45,
		RecencyDays:   500,
		ContentLength: 200,
		CategoryMatch: false,
		TitleHasVote:  false,
		IsElection:    false,
	}
}

// PenalizedScoreFactors returns factors with title penalty applied.
func PenalizedScoreFactors() ScoreFactorsData {
	return ScoreFactorsData{
		LLMScore:      85,
		VectorScore:   0.88,
		RecencyDays:   60,
		ContentLength: 600,
		CategoryMatch: false,
		TitleHasVote:  true,
		IsElection:    false,
	}
}

// CreateTestConfig writes a sample config YAML file for testing.
func CreateTestConfig(t *testing.T, path string) {
	t.Helper()
	config := `llama_server:
  base_url: http://localhost:8080
  embedding_model: test-embedding-model
  llm_model: test-llm-model
  timeout: 120s
  min_request_delay: 100ms
  max_retries: 3

processing:
  posts_directory: content/p
  top_candidates: 10
  batch_size: 5
  max_content_length: 100000

storage:
  backup_directory: backups
  cache_file: .cache/embeddings.json
  progress_file: .cache/progress.json
  compression_cache: .cache/compression.json

scoring:
  weights:
    relevance: 0.6
    recency: 0.2
    length: 0.1
    category: 0.1
  category_boost:
    elections: 1.5
    backstory: 1.0

filtering:
  exclude_categories:
    - holiday
  exclude_titles:
    - "backstory podcast"
    - "wonderful wednesday"
    - "freaky friday"
  start_date: 2020-01-01

logging:
  level: info
  file: pipeline.log

concurrency:
  max_workers: 4
  max_concurrent_api_calls: 2
`
	WriteFile(t, path, config)
}

// MockLlamaResponse provides sample responses from llama.cpp server.
type MockLlamaResponse struct {
	EmbeddingResponse string
	CompletionResponse string
	HealthCheckResponse string
}

// DefaultMockLlamaResponse returns default mock responses.
func DefaultMockLlamaResponse() MockLlamaResponse {
	return MockLlamaResponse{
		EmbeddingResponse: `{"embedding": [0.1, 0.2, 0.3, 0.4, 0.5]}`,
		CompletionResponse: `{"content": "85", "tokens_evaluated": 150, "tokens_predicted": 5}`,
		HealthCheckResponse: `{"status": "ok", "model_loaded": true}`,
	}
}
