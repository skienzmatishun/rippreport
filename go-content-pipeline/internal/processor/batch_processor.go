package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rippreport/go-content-pipeline/internal/client"
	"github.com/rippreport/go-content-pipeline/internal/models"
)

// BatchAssembler groups requests into batches respecting size and token constraints.
//
// Requirements: 14.1, 14.2, 14.5
type BatchAssembler struct {
	maxBatchSize  int
	maxTokens     int
	currentTokens int
	items         []models.HookRequest
}

// NewBatchAssembler creates a new BatchAssembler.
func NewBatchAssembler(maxBatchSize, maxTokens int) *BatchAssembler {
	if maxBatchSize <= 0 {
		maxBatchSize = 5
	}
	if maxTokens <= 0 {
		maxTokens = 4096
	}
	return &BatchAssembler{
		maxBatchSize: maxBatchSize,
		maxTokens:    maxTokens,
	}
}

// Add attempts to add a request to the current batch. Returns false if the batch is full.
func (ba *BatchAssembler) Add(req models.HookRequest, tokens int) bool {
	if len(ba.items) >= ba.maxBatchSize || (ba.currentTokens+tokens > ba.maxTokens && len(ba.items) > 0) {
		return false
	}
	ba.items = append(ba.items, req)
	ba.currentTokens += tokens
	return true
}

// Flush clears the current batch and returns the queued requests.
func (ba *BatchAssembler) Flush() []models.HookRequest {
	flushed := ba.items
	ba.items = nil
	ba.currentTokens = 0
	return flushed
}

// Len returns the current number of items in the batch.
func (ba *BatchAssembler) Len() int {
	return len(ba.items)
}

// ScoreBatch attempts to score multiple candidates in a single LLM completion,
// with automatic fallback to one-by-one scoring on failure.
//
// Requirements: 5.6, 5.7, 14.3, 14.4
func ScoreBatch(ctx context.Context, llamaClient client.LlamaClient, modelName string, currentPost *models.Post, candidates []*models.Post) ([]int, error) {
	if len(candidates) == 0 {
		return nil, nil
	}

	if llamaClient == nil {
		// Vector-only fallback
		scores := make([]int, len(candidates))
		for i := range candidates {
			scores[i] = 80
		}
		return scores, nil
	}

	var candidateList strings.Builder
	for i, c := range candidates {
		candidateList.WriteString(fmt.Sprintf("%d. Title: %q (Slug: %s)\n", i+1, c.FrontMatter.Title, c.Slug))
	}

	prompt := fmt.Sprintf(`Given the current article %q, rate the topical relevance of each candidate article below on a scale from 0 to 100.
Respond ONLY with a JSON array of %d integer scores, for example: [85, 92, 70].

Candidates:
%s`, currentPost.FrontMatter.Title, len(candidates), candidateList.String())

	resp, err := llamaClient.Complete(ctx, models.CompletionRequest{
		Prompt:      prompt,
		Model:       modelName,
		Temperature: 0.1,
		MaxTokens:   100,
	})

	if err == nil && resp != nil {
		var scores []int
		cleaned := strings.TrimSpace(resp.Text)
		startIdx := strings.Index(cleaned, "[")
		endIdx := strings.LastIndex(cleaned, "]")
		if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
			if jsonErr := json.Unmarshal([]byte(cleaned[startIdx:endIdx+1]), &scores); jsonErr == nil && len(scores) == len(candidates) {
				return scores, nil
			}
		}
	}

	// Fallback to one-by-one scoring
	scores := make([]int, len(candidates))
	for i, cand := range candidates {
		singlePrompt := fmt.Sprintf("Rate relevance of %q to %q (0-100). Respond with single number.", cand.FrontMatter.Title, currentPost.FrontMatter.Title)
		singleResp, sErr := llamaClient.Complete(ctx, models.CompletionRequest{
			Prompt:      singlePrompt,
			Model:       modelName,
			Temperature: 0.1,
			MaxTokens:   10,
		})
		if sErr == nil && singleResp != nil {
			if sc, ok := parseScoreFromText(singleResp.Text); ok {
				scores[i] = sc
				continue
			}
		}
		scores[i] = 50 // Safe fallback
	}

	return scores, nil
}

// GenerateHooksBatch generates hooks for multiple article pairs in batch,
// utilizing the compression cache and falling back to individual generation on error.
//
// Requirements: 12.7, 14.3, 14.4
func GenerateHooksBatch(ctx context.Context, llamaClient client.LlamaClient, modelName string, cache *CompressionCache, reqs []models.HookRequest) ([]*models.Hook, error) {
	if len(reqs) == 0 {
		return nil, nil
	}

	hooks := make([]*models.Hook, 0, len(reqs))
	var firstErr error
	for i, req := range reqs {
		fmt.Printf("      [%d/%d] → %s...", i+1, len(reqs), req.WidgetPost.Slug)
		hook, hErr := GenerateHook(ctx, llamaClient, modelName, cache, req)
		if hErr != nil {
			fmt.Printf(" ✗ %v\n", hErr)
			if firstErr == nil {
				firstErr = hErr
			}
			continue
		}
		fmt.Printf(" ✓\n")
		hooks = append(hooks, hook)
	}
	if len(hooks) == 0 {
		if firstErr != nil {
			return nil, firstErr
		}
		return nil, fmt.Errorf("no hooks generated")
	}
	return hooks, nil
}
