package processor

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/rippreport/go-content-pipeline/internal/client"
	"github.com/rippreport/go-content-pipeline/internal/models"
)

var scoreRegex = regexp.MustCompile(`\b(100|[1-9]?[0-9])\b`)

// calculateRecencyScore returns a recency score [0, 100] using exponential decay.
// Formula: 100 * exp(-days/365)
//
// Requirements: 6.2
func calculateRecencyScore(daysDiff int) float64 {
	if daysDiff < 0 {
		daysDiff = 0
	}
	score := 100.0 * math.Exp(-float64(daysDiff)/365.0)
	if score > 100.0 {
		return 100.0
	}
	if score < 0.0 {
		return 0.0
	}
	return score
}

// calculateLengthScore returns a content length score [0, 100].
// Target: 400 characters = 100 score. Cap at 100 for >= 400.
//
// Requirements: 6.3, 6.4
func calculateLengthScore(contentLength int) float64 {
	if contentLength <= 0 {
		return 0.0
	}
	if contentLength >= 400 {
		return 100.0
	}
	return (float64(contentLength) / 400.0) * 100.0
}

// CalculateCompositeScore calculates the weighted composite score clamped to [0, 100].
//
// Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 6.7, 6.8
func CalculateCompositeScore(factors models.ScoreFactors) int {
	w := factors.Weights.Normalize()

	llmPart := w.Relevance * float32(factors.LLMScore)
	recencyPart := w.Recency * float32(calculateRecencyScore(factors.RecencyDays))
	lengthPart := w.Length * float32(calculateLengthScore(factors.ContentLength))

	categoryPart := float32(0.0)
	if factors.CategoryMatch {
		categoryPart = w.Category * 100.0
	}

	composite := llmPart + recencyPart + lengthPart + categoryPart

	// Apply title penalty if specified
	if factors.TitlePenalty > 0 && factors.TitlePenalty <= 1.0 {
		composite *= (1.0 - factors.TitlePenalty)
	}

	rounded := int(math.Round(float64(composite)))
	if rounded > 100 {
		return 100
	}
	if rounded < 0 {
		return 0
	}
	return rounded
}

// ScoreCandidates scores candidates using LLM or vector-only scores and ranks them.
//
// Requirements: 5.1, 5.2, 5.3, 5.4, 6.1 through 6.8
func ScoreCandidates(ctx context.Context, llamaClient client.LlamaClient, req models.ScoringRequest, weights models.ScoreWeights) ([]models.ScoredArticle, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid scoring request: %w", err)
	}

	// If reranker model is requested, perform batch cross-encoder reranking
	rerankerScores := make(map[string]int)
	if req.UseReranker && !req.VectorOnly && llamaClient != nil {
		docs := make([]string, len(req.Candidates))
		for i, cand := range req.Candidates {
			docs[i] = cand.Slug
		}
		query := req.CurrentPost.FrontMatter.Title
		rerankResp, rErr := llamaClient.Rerank(ctx, models.RerankRequest{
			Model:     req.Model,
			Query:     query,
			Documents: docs,
		})
		if rErr == nil && rerankResp != nil {
			for _, res := range rerankResp.Results {
				if res.Index >= 0 && res.Index < len(req.Candidates) {
					rawScore := res.RelevanceScore
					var intScore int
					if rawScore <= 1.0 && rawScore >= 0.0 {
						intScore = int(math.Round(rawScore * 100.0))
					} else if rawScore > 1.0 {
						intScore = int(math.Min(100, math.Round(rawScore)))
					} else {
						// Logit format (negative values): apply sigmoid 1 / (1 + exp(-x))
						prob := 1.0 / (1.0 + math.Exp(-rawScore))
						intScore = int(math.Round(prob * 100.0))
					}
					if intScore < 0 {
						intScore = 0
					}
					if intScore > 100 {
						intScore = 100
					}
					rerankerScores[req.Candidates[res.Index].Slug] = intScore
				}
			}
		}
	}

	var scoredArticles []models.ScoredArticle

	for _, cand := range req.Candidates {
		llmScore := int(math.Round(float64(cand.Similarity * 100.0)))
		if llmScore < 0 {
			llmScore = 0
		}
		if llmScore > 100 {
			llmScore = 100
		}

		if req.UseReranker {
			if s, ok := rerankerScores[cand.Slug]; ok {
				llmScore = s
			}
		} else if !req.VectorOnly && llamaClient != nil {
			prompt := fmt.Sprintf("Rate the topical relevance of article %q to the current article %q on a scale from 0 to 100. Respond with only a single integer score.", cand.Slug, req.CurrentPost.FrontMatter.Title)
			resp, err := llamaClient.Complete(ctx, models.CompletionRequest{
				Prompt:      prompt,
				Model:       req.Model,
				Temperature: 0.1,
				MaxTokens:   10,
			})
			if err == nil && resp != nil {
				if parsedScore, ok := parseScoreFromText(resp.Text); ok {
					llmScore = parsedScore
				}
			}
		}

		// Calculate composite factors
		factors := models.ScoreFactors{
			LLMScore:      llmScore,
			VectorScore:   cand.Similarity,
			RecencyDays:   30, // Default recency baseline if dates not passed
			ContentLength: len(req.CurrentPost.Body),
			CategoryMatch: false,
			Weights:       weights,
		}

		// Check title penalty: penalty for "vote" if current post not election-related
		if strings.Contains(strings.ToLower(cand.Slug), "vote") && !strings.Contains(strings.ToLower(req.CurrentPost.FrontMatter.Title), "election") {
			factors.TitlePenalty = 0.2
		}

		finalScore := CalculateCompositeScore(factors)

		scoredArticles = append(scoredArticles, models.ScoredArticle{
			Slug:        cand.Slug,
			Title:       cand.Slug, // Will be enriched with actual post title by caller if needed
			VectorScore: cand.Similarity,
			LLMScore:    llmScore,
			FinalScore:  finalScore,
		})
	}

	// Sort descending by FinalScore, tie break by VectorScore, then Slug
	sort.Slice(scoredArticles, func(i, j int) bool {
		if scoredArticles[i].FinalScore == scoredArticles[j].FinalScore {
			if scoredArticles[i].VectorScore == scoredArticles[j].VectorScore {
				return scoredArticles[i].Slug < scoredArticles[j].Slug
			}
			return scoredArticles[i].VectorScore > scoredArticles[j].VectorScore
		}
		return scoredArticles[i].FinalScore > scoredArticles[j].FinalScore
	})

	// Assign ranks 1 to N
	for i := range scoredArticles {
		scoredArticles[i].Rank = i + 1
	}

	return scoredArticles, nil
}

func parseScoreFromText(text string) (int, bool) {
	match := scoreRegex.FindString(strings.TrimSpace(text))
	if match == "" {
		return 0, false
	}
	score, err := strconv.Atoi(match)
	if err != nil || score < 0 || score > 100 {
		return 0, false
	}
	return score, true
}
