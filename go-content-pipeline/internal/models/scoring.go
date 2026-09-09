package models

import (
	"fmt"
)

// Candidate represents a post candidate from similarity search
type Candidate struct {
	Path       string  `json:"path"`
	Slug       string  `json:"slug"`
	Similarity float32 `json:"similarity"`
}

// ScoredArticle represents a scored and ranked article recommendation
type ScoredArticle struct {
	Slug        string  `json:"slug"`
	Title       string  `json:"title"`
	VectorScore float32 `json:"vector_score"`
	LLMScore    int     `json:"llm_score"`
	FinalScore  int     `json:"final_score"`
	Rank        int     `json:"rank"`
}

// ScoreFactors contains all factors used in composite scoring
type ScoreFactors struct {
	LLMScore      int
	VectorScore   float32
	RecencyDays   int
	ContentLength int
	CategoryMatch bool
	TitlePenalty  float32
	Weights       ScoreWeights
}

// ScoreWeights defines the weights for different scoring factors
type ScoreWeights struct {
	Relevance float32 `yaml:"relevance"`
	Recency   float32 `yaml:"recency"`
	Length    float32 `yaml:"length"`
	Category  float32 `yaml:"category"`
}

// ScoringRequest represents a request to score candidates
type ScoringRequest struct {
	CurrentPost   *Post
	Candidates    []Candidate
	BatchMode     bool
	VectorOnly    bool
	Model         string
	UseReranker   bool
	RerankBaseURL string
}

// Validate checks if the Candidate has valid fields
func (c *Candidate) Validate() error {
	if c.Path == "" {
		return fmt.Errorf("candidate path cannot be empty")
	}
	if c.Slug == "" {
		return fmt.Errorf("candidate slug cannot be empty")
	}
	if c.Similarity < -1.0 || c.Similarity > 1.0 {
		return fmt.Errorf("candidate similarity must be between -1 and 1, got %f", c.Similarity)
	}
	return nil
}

// Validate checks if the ScoredArticle has valid fields
func (sa *ScoredArticle) Validate() error {
	if sa.Slug == "" {
		return fmt.Errorf("scored article slug cannot be empty")
	}
	if sa.Title == "" {
		return fmt.Errorf("scored article title cannot be empty")
	}
	if sa.VectorScore < -1.0 || sa.VectorScore > 1.0 {
		return fmt.Errorf("scored article vector_score must be between -1 and 1, got %f", sa.VectorScore)
	}
	if sa.LLMScore < 0 || sa.LLMScore > 100 {
		return fmt.Errorf("scored article llm_score must be between 0 and 100, got %d", sa.LLMScore)
	}
	if sa.FinalScore < 0 || sa.FinalScore > 100 {
		return fmt.Errorf("scored article final_score must be between 0 and 100, got %d", sa.FinalScore)
	}
	if sa.Rank < 1 {
		return fmt.Errorf("scored article rank must be at least 1, got %d", sa.Rank)
	}
	return nil
}

// Normalize ensures weights sum to 1.0 and returns normalized weights
func (sw *ScoreWeights) Normalize() ScoreWeights {
	sum := sw.Relevance + sw.Recency + sw.Length + sw.Category
	if sum == 0 {
		// Default to equal weights if all are zero
		return ScoreWeights{
			Relevance: 0.25,
			Recency:   0.25,
			Length:    0.25,
			Category:  0.25,
		}
	}
	return ScoreWeights{
		Relevance: sw.Relevance / sum,
		Recency:   sw.Recency / sum,
		Length:    sw.Length / sum,
		Category:  sw.Category / sum,
	}
}

// Validate checks if the ScoreWeights are valid
func (sw *ScoreWeights) Validate() error {
	if sw.Relevance < 0 {
		return fmt.Errorf("relevance weight cannot be negative")
	}
	if sw.Recency < 0 {
		return fmt.Errorf("recency weight cannot be negative")
	}
	if sw.Length < 0 {
		return fmt.Errorf("length weight cannot be negative")
	}
	if sw.Category < 0 {
		return fmt.Errorf("category weight cannot be negative")
	}
	
	sum := sw.Relevance + sw.Recency + sw.Length + sw.Category
	if sum == 0 {
		return fmt.Errorf("at least one weight must be positive")
	}
	
	return nil
}

// Validate checks if the ScoreFactors are valid
func (sf *ScoreFactors) Validate() error {
	if sf.LLMScore < 0 || sf.LLMScore > 100 {
		return fmt.Errorf("llm_score must be between 0 and 100, got %d", sf.LLMScore)
	}
	if sf.VectorScore < -1.0 || sf.VectorScore > 1.0 {
		return fmt.Errorf("vector_score must be between -1 and 1, got %f", sf.VectorScore)
	}
	if sf.RecencyDays < 0 {
		return fmt.Errorf("recency_days cannot be negative, got %d", sf.RecencyDays)
	}
	if sf.ContentLength < 0 {
		return fmt.Errorf("content_length cannot be negative, got %d", sf.ContentLength)
	}
	if sf.TitlePenalty < 0 || sf.TitlePenalty > 1.0 {
		return fmt.Errorf("title_penalty must be between 0 and 1, got %f", sf.TitlePenalty)
	}
	return sf.Weights.Validate()
}

// Validate checks if the ScoringRequest is valid
func (sr *ScoringRequest) Validate() error {
	if sr.CurrentPost == nil {
		return fmt.Errorf("scoring request current_post cannot be nil")
	}
	if err := sr.CurrentPost.Validate(); err != nil {
		return fmt.Errorf("current_post validation failed: %w", err)
	}
	if len(sr.Candidates) == 0 {
		return fmt.Errorf("scoring request must have at least one candidate")
	}
	for i, candidate := range sr.Candidates {
		if err := candidate.Validate(); err != nil {
			return fmt.Errorf("candidate %d validation failed: %w", i, err)
		}
	}
	return nil
}

// ToRelatedArticle converts a ScoredArticle to a RelatedArticle
func (sa *ScoredArticle) ToRelatedArticle() RelatedArticle {
	return RelatedArticle{
		Slug:  sa.Slug,
		Title: sa.Title,
		Score: sa.FinalScore,
		Rank:  sa.Rank,
	}
}
