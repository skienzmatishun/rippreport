package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/config"
	"github.com/rippreport/go-content-pipeline/internal/models"
)

// Standard timeouts according to specification
const (
	defaultEmbeddingTimeout  = 120 * time.Second
	defaultCompletionTimeout = 300 * time.Second
	defaultHealthTimeout     = 5 * time.Second
	maxEmbeddingChars        = 16384 // Safeguard truncation
)

// LlamaClient defines the interface for interacting with llama.cpp server.
type LlamaClient interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
	Complete(ctx context.Context, req models.CompletionRequest) (*models.CompletionResponse, error)
	HealthCheck(ctx context.Context) error
	GetStats() models.ClientStats
}

// TokenBucketRateLimiter enforces a minimum delay between outgoing requests.
//
// Requirements: 2.5, 5.5
type TokenBucketRateLimiter struct {
	mu           sync.Mutex
	minDelay     time.Duration
	lastRequest  time.Time
	stats        models.RateLimitStats
}

// NewTokenBucketRateLimiter creates a new rate limiter with the specified minimum delay.
func NewTokenBucketRateLimiter(minDelay time.Duration) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		minDelay: minDelay,
	}
}

// Wait blocks until the minimum delay has elapsed since the last request,
// or until the context is cancelled.
func (rl *TokenBucketRateLimiter) Wait(ctx context.Context) error {
	if rl.minDelay <= 0 {
		return nil
	}

	rl.mu.Lock()
	now := time.Now()
	elapsed := now.Sub(rl.lastRequest)
	var waitDuration time.Duration
	if elapsed < rl.minDelay && !rl.lastRequest.IsZero() {
		waitDuration = rl.minDelay - elapsed
	}
	rl.lastRequest = now.Add(waitDuration)
	if waitDuration > 0 {
		rl.stats.TotalWaits++
		rl.stats.TotalWaitTime += waitDuration
		rl.stats.CalculateAvgWaitTime()
	}
	rl.mu.Unlock()

	if waitDuration > 0 {
		select {
		case <-time.After(waitDuration):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// GetStats returns current rate limit statistics.
func (rl *TokenBucketRateLimiter) GetStats() models.RateLimitStats {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.stats
}

// Client implements LlamaClient for communicating with llama.cpp server.
//
// Requirements: 2.1 through 2.8, 5.1, 5.2, 5.5, 5.8, 24.1 through 24.8
type Client struct {
	cfg         config.LlamaServerConfig
	httpClient  *http.Client
	rateLimiter *TokenBucketRateLimiter
	mu          sync.Mutex
	stats       models.ClientStats
	totalRespTime time.Duration
}

// NewClient creates a new Client configured with connection pooling, keep-alive, and rate limiting.
func NewClient(cfg config.LlamaServerConfig) *Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	minDelay := cfg.MinRequestDelay
	if minDelay == 0 {
		minDelay = 50 * time.Millisecond
	}

	return &Client{
		cfg:         cfg,
		httpClient:  httpClient,
		rateLimiter: NewTokenBucketRateLimiter(minDelay),
	}
}

type embeddingRequestBody struct {
	Content string `json:"content"`
	Input   string `json:"input,omitempty"`
	Model   string `json:"model,omitempty"`
}

type embeddingResponseBody struct {
	Embedding []float32 `json:"embedding"`
	Data      []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

// GenerateEmbedding calls the /embedding endpoint with automatic retries and rate limiting.
//
// Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 2.8
func (c *Client) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Truncate if exceeds safe context window
	if len(text) > maxEmbeddingChars {
		text = text[:maxEmbeddingChars]
	}

	timeout := c.cfg.Timeout
	if timeout <= 0 {
		timeout = defaultEmbeddingTimeout
	}

	reqBody := embeddingRequestBody{
		Content: text,
		Input:   text,
		Model:   c.cfg.EmbeddingModel,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal embedding request: %w", err)
	}

	endpoint := strings.TrimRight(c.cfg.BaseURL, "/") + "/embedding"

	respBytes, err := c.requestWithRetry(ctx, http.MethodPost, endpoint, data, timeout)
	if err != nil {
		return nil, err
	}

	var resp embeddingResponseBody
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse embedding response: %w", err)
	}

	if len(resp.Embedding) > 0 {
		return resp.Embedding, nil
	}
	if len(resp.Data) > 0 && len(resp.Data[0].Embedding) > 0 {
		return resp.Data[0].Embedding, nil
	}

	return nil, fmt.Errorf("empty embedding vector returned from server")
}

// Complete calls the /completion endpoint with given parameters.
//
// Requirements: 5.1, 5.2, 5.8
func (c *Client) Complete(ctx context.Context, req models.CompletionRequest) (*models.CompletionResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid completion request: %w", err)
	}

	timeout := defaultCompletionTimeout

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal completion request: %w", err)
	}

	endpoint := strings.TrimRight(c.cfg.BaseURL, "/") + "/completion"

	start := time.Now()
	respBytes, err := c.requestWithRetry(ctx, http.MethodPost, endpoint, data, timeout)
	if err != nil {
		return nil, err
	}
	duration := time.Since(start)

	var resp models.CompletionResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse completion response: %w", err)
	}
	resp.Duration = duration

	if err := resp.Validate(); err != nil {
		return nil, fmt.Errorf("invalid completion response from server: %w", err)
	}

	return &resp, nil
}

// HealthCheck verifies that the server is reachable and responsive.
//
// Requirements: 24.1, 24.2, 25.1, 25.2, 25.3
func (c *Client) HealthCheck(ctx context.Context) error {
	endpoint := strings.TrimRight(c.cfg.BaseURL, "/") + "/health"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	client := &http.Client{Timeout: defaultHealthTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("llama.cpp server unreachable at %s: %w", c.cfg.BaseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusServiceUnavailable {
		return fmt.Errorf("llama.cpp server is loading models (503 Service Unavailable)")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("llama.cpp server returned unhealthy status: %d", resp.StatusCode)
	}

	return nil
}

// requestWithRetry handles executing an HTTP request with rate limiting and exponential backoff retry.
//
// Requirements: 2.6, 24.3, 24.4, 24.5, 24.7
func (c *Client) requestWithRetry(ctx context.Context, method, url string, body []byte, timeout time.Duration) ([]byte, error) {
	maxRetries := c.cfg.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	backoffs := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}

	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Respect rate limit before each attempt
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limit wait cancelled: %w", err)
		}

		c.recordRequest()

		// Context with per-attempt timeout
		attemptTimeout := timeout
		if attempt > 0 && errors.Is(lastErr, context.DeadlineExceeded) {
			// Increase timeout for retry after timeout
			attemptTimeout += 30 * time.Second
		}

		attemptCtx, cancel := context.WithTimeout(ctx, attemptTimeout)
		req, err := http.NewRequestWithContext(attemptCtx, method, url, bytes.NewReader(body))
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		startTime := time.Now()
		resp, err := c.httpClient.Do(req)
		callDuration := time.Since(startTime)

		if err != nil {
			cancel()
			lastErr = err

			// Check if context cancelled by parent
			if ctx.Err() != nil {
				c.recordFailure()
				return nil, ctx.Err()
			}

			// Network error or timeout -> retry
			if attempt < maxRetries {
				c.recordRetry()
				backoff := backoffs[min(attempt, len(backoffs)-1)]
				select {
				case <-time.After(backoff):
					continue
				case <-ctx.Done():
					c.recordFailure()
					return nil, ctx.Err()
				}
			}
			c.recordFailure()
			return nil, fmt.Errorf("request to %s failed after %d attempts: %w", url, attempt+1, err)
		}

		respBody, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		cancel()

		if readErr != nil {
			lastErr = readErr
			if attempt < maxRetries {
				c.recordRetry()
				continue
			}
			c.recordFailure()
			return nil, fmt.Errorf("failed to read response body: %w", readErr)
		}

		// Success
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			c.recordSuccess(callDuration)
			return respBody, nil
		}

		// 4xx errors: fail immediately (do not retry)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			c.recordFailure()
			if resp.StatusCode == http.StatusRequestEntityTooLarge || resp.StatusCode == 413 {
				return nil, fmt.Errorf("server context window exceeded (status 413): %s", string(respBody))
			}
			return nil, fmt.Errorf("client error status %d from %s: %s", resp.StatusCode, url, string(respBody))
		}

		// 503 Model not loaded
		if resp.StatusCode == http.StatusServiceUnavailable {
			lastErr = fmt.Errorf("model not loaded (503 Service Unavailable): %s", string(respBody))
		} else {
			lastErr = fmt.Errorf("server error status %d from %s: %s", resp.StatusCode, url, string(respBody))
		}

		// 5xx errors: retry with backoff
		if attempt < maxRetries {
			c.recordRetry()
			backoff := backoffs[min(attempt, len(backoffs)-1)]
			select {
			case <-time.After(backoff):
				continue
			case <-ctx.Done():
				c.recordFailure()
				return nil, ctx.Err()
			}
		}
	}

	c.recordFailure()
	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

func (c *Client) recordRequest() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stats.TotalRequests++
}

func (c *Client) recordSuccess(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stats.SuccessfulReqs++
	c.totalRespTime += d
	c.stats.CalculateAvgResponseTime(c.totalRespTime)
}

func (c *Client) recordFailure() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stats.FailedReqs++
}

func (c *Client) recordRetry() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stats.RetriedReqs++
}

// GetStats returns the HTTP client statistics.
//
// Requirements: 2.7, 17.7, 24.8
func (c *Client) GetStats() models.ClientStats {
	c.mu.Lock()
	defer c.mu.Unlock()
	rlStats := c.rateLimiter.GetStats()
	c.stats.TotalWaitTime = rlStats.TotalWaitTime
	return c.stats
}
