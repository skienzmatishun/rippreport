package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/config"
	"github.com/rippreport/go-content-pipeline/internal/models"
)

func TestClient_EmbeddingSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/embedding" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"embedding": []float32{0.1, 0.2, 0.3},
		})
	}))
	defer ts.Close()

	cfg := config.LlamaServerConfig{
		BaseURL:         ts.URL,
		EmbeddingModel:  "test-embed",
		Timeout:         5 * time.Second,
		MinRequestDelay: 10 * time.Millisecond,
		MaxRetries:      2,
	}

	c := NewClient(cfg)
	vec, err := c.GenerateEmbedding(context.Background(), "Hello world text")
	if err != nil {
		t.Fatalf("GenerateEmbedding failed: %v", err)
	}

	if len(vec) != 3 || vec[0] != 0.1 {
		t.Errorf("unexpected vector: %v", vec)
	}

	stats := c.GetStats()
	if stats.SuccessfulReqs != 1 {
		t.Errorf("expected 1 successful request, got %d", stats.SuccessfulReqs)
	}
}

func TestClient_CompletionSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/completion" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"content":          "Generated LLM text response",
			"tokens_evaluated": 15,
			"tokens_predicted": 25,
		})
	}))
	defer ts.Close()

	cfg := config.LlamaServerConfig{
		BaseURL:         ts.URL,
		LLMModel:        "test-llm",
		Timeout:         5 * time.Second,
		MinRequestDelay: 10 * time.Millisecond,
		MaxRetries:      2,
	}

	c := NewClient(cfg)
	resp, err := c.Complete(context.Background(), models.CompletionRequest{
		Prompt:      "Test prompt",
		Temperature: 0.7,
	})
	if err != nil {
		t.Fatalf("Complete failed: %v", err)
	}

	if resp.Text != "Generated LLM text response" {
		t.Errorf("unexpected text: %s", resp.Text)
	}
	if resp.TokensPrompt != 15 || resp.TokensGen != 25 {
		t.Errorf("unexpected tokens: prompt=%d, gen=%d", resp.TokensPrompt, resp.TokensGen)
	}
}

func TestClient_HealthCheck(t *testing.T) {
	var isHealthy atomic.Bool
	isHealthy.Store(true)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isHealthy.Load() {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}))
	defer ts.Close()

	cfg := config.LlamaServerConfig{
		BaseURL: ts.URL,
	}
	c := NewClient(cfg)

	// Healthy
	if err := c.HealthCheck(context.Background()); err != nil {
		t.Errorf("HealthCheck should succeed: %v", err)
	}

	// Unhealthy
	isHealthy.Store(false)
	if err := c.HealthCheck(context.Background()); err == nil {
		t.Errorf("HealthCheck should fail when service unavailable")
	}
}

func TestClient_RetryOn5xx(t *testing.T) {
	var attempts atomic.Int32

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		att := attempts.Add(1)
		if att < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("temporary error"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"embedding": []float32{0.5, 0.6},
		})
	}))
	defer ts.Close()

	cfg := config.LlamaServerConfig{
		BaseURL:         ts.URL,
		MinRequestDelay: 5 * time.Millisecond,
		MaxRetries:      3,
	}
	c := NewClient(cfg)

	vec, err := c.GenerateEmbedding(context.Background(), "retry test")
	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if len(vec) != 2 {
		t.Errorf("unexpected vector length: %d", len(vec))
	}
	if attempts.Load() != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts.Load())
	}
}

func TestClient_FailFastOn4xx(t *testing.T) {
	var attempts atomic.Int32

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
	}))
	defer ts.Close()

	cfg := config.LlamaServerConfig{
		BaseURL:         ts.URL,
		MinRequestDelay: 5 * time.Millisecond,
		MaxRetries:      3,
	}
	c := NewClient(cfg)

	_, err := c.GenerateEmbedding(context.Background(), "fail fast test")
	if err == nil {
		t.Fatalf("expected error on 400 Bad Request, got nil")
	}
	// Should fail on first attempt without retrying
	if attempts.Load() != 1 {
		t.Errorf("expected exactly 1 attempt for 4xx error, got %d", attempts.Load())
	}
}

func TestTokenBucketRateLimiter(t *testing.T) {
	minDelay := 20 * time.Millisecond
	limiter := NewTokenBucketRateLimiter(minDelay)

	start := time.Now()
	for i := 0; i < 4; i++ {
		if err := limiter.Wait(context.Background()); err != nil {
			t.Fatalf("Wait failed: %v", err)
		}
	}
	elapsed := time.Since(start)

	// 4 requests with 20ms delay should take at least ~50ms
	if elapsed < 50*time.Millisecond {
		t.Errorf("rate limiter did not enforce delay, elapsed %v", elapsed)
	}

	stats := limiter.GetStats()
	if stats.TotalWaits == 0 {
		t.Errorf("expected TotalWaits > 0, got %d", stats.TotalWaits)
	}
}
