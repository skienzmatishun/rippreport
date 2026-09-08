package models

import (
	"fmt"
	"time"
)

// CompletionRequest represents a request for LLM completion
type CompletionRequest struct {
	Prompt     string   `json:"prompt"`
	Temperature float32  `json:"temperature,omitempty"`
	TopP       float32  `json:"top_p,omitempty"`
	TopK       int      `json:"top_k,omitempty"`
	MaxTokens  int      `json:"n_predict,omitempty"`
	StopTokens []string `json:"stop,omitempty"`
}

// CompletionResponse represents the response from LLM completion
type CompletionResponse struct {
	Text         string        `json:"content"`
	TokensPrompt int           `json:"tokens_evaluated"`
	TokensGen    int           `json:"tokens_predicted"`
	Duration     time.Duration `json:"-"`
}

// ModelInfo represents information about a loaded model
type ModelInfo struct {
	Name         string `json:"name"`
	ContextSize  int    `json:"n_ctx"`
	EmbeddingDim int    `json:"n_embd"`
}

// ClientStats tracks statistics for the HTTP client
type ClientStats struct {
	TotalRequests   int64         `json:"total_requests"`
	SuccessfulReqs  int64         `json:"successful_requests"`
	FailedReqs      int64         `json:"failed_requests"`
	RetriedReqs     int64         `json:"retried_requests"`
	TotalWaitTime   time.Duration `json:"total_wait_time"`
	AvgResponseTime time.Duration `json:"avg_response_time"`
}

// RateLimitStats tracks statistics for rate limiting
type RateLimitStats struct {
	TotalWaits    int64         `json:"total_waits"`
	TotalWaitTime time.Duration `json:"total_wait_time"`
	AvgWaitTime   time.Duration `json:"avg_wait_time"`
}

// Validate checks if the CompletionRequest is valid
func (cr *CompletionRequest) Validate() error {
	if cr.Prompt == "" {
		return fmt.Errorf("completion request prompt cannot be empty")
	}
	if cr.Temperature < 0 || cr.Temperature > 2.0 {
		return fmt.Errorf("completion request temperature must be between 0 and 2, got %f", cr.Temperature)
	}
	if cr.TopP < 0 || cr.TopP > 1.0 {
		return fmt.Errorf("completion request top_p must be between 0 and 1, got %f", cr.TopP)
	}
	if cr.TopK < 0 {
		return fmt.Errorf("completion request top_k cannot be negative, got %d", cr.TopK)
	}
	if cr.MaxTokens < 0 {
		return fmt.Errorf("completion request max_tokens cannot be negative, got %d", cr.MaxTokens)
	}
	return nil
}

// Validate checks if the CompletionResponse is valid
func (cr *CompletionResponse) Validate() error {
	if cr.Text == "" {
		return fmt.Errorf("completion response text cannot be empty")
	}
	if cr.TokensPrompt < 0 {
		return fmt.Errorf("completion response tokens_prompt cannot be negative, got %d", cr.TokensPrompt)
	}
	if cr.TokensGen < 0 {
		return fmt.Errorf("completion response tokens_gen cannot be negative, got %d", cr.TokensGen)
	}
	return nil
}

// Validate checks if the ModelInfo is valid
func (mi *ModelInfo) Validate() error {
	if mi.Name == "" {
		return fmt.Errorf("model info name cannot be empty")
	}
	if mi.ContextSize <= 0 {
		return fmt.Errorf("model info context_size must be positive, got %d", mi.ContextSize)
	}
	if mi.EmbeddingDim <= 0 {
		return fmt.Errorf("model info embedding_dim must be positive, got %d", mi.EmbeddingDim)
	}
	return nil
}

// CalculateAvgResponseTime calculates the average response time
func (cs *ClientStats) CalculateAvgResponseTime(totalDuration time.Duration) {
	if cs.SuccessfulReqs == 0 {
		cs.AvgResponseTime = 0
	} else {
		cs.AvgResponseTime = totalDuration / time.Duration(cs.SuccessfulReqs)
	}
}

// CalculateAvgWaitTime calculates the average wait time
func (rls *RateLimitStats) CalculateAvgWaitTime() {
	if rls.TotalWaits == 0 {
		rls.AvgWaitTime = 0
	} else {
		rls.AvgWaitTime = rls.TotalWaitTime / time.Duration(rls.TotalWaits)
	}
}
