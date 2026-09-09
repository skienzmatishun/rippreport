package processor

import (
	"context"
	"testing"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/client"
	"github.com/rippreport/go-content-pipeline/internal/models"
)

func TestCalculateRecencyScore(t *testing.T) {
	// 0 days diff should be 100
	score0 := calculateRecencyScore(0)
	if score0 != 100.0 {
		t.Errorf("expected 100 for 0 days, got %f", score0)
	}

	// 365 days diff should be 100 / e ~ 36.78
	score365 := calculateRecencyScore(365)
	if score365 < 36.0 || score365 > 37.5 {
		t.Errorf("expected ~36.8 for 365 days, got %f", score365)
	}

	// 730 days should be even lower
	score730 := calculateRecencyScore(730)
	if score730 >= score365 {
		t.Errorf("recency score should decay over time: %f vs %f", score730, score365)
	}
}

func TestCalculateLengthScore(t *testing.T) {
	if s := calculateLengthScore(0); s != 0.0 {
		t.Errorf("expected 0 for 0 len, got %f", s)
	}
	if s := calculateLengthScore(200); s != 50.0 {
		t.Errorf("expected 50 for 200 len, got %f", s)
	}
	if s := calculateLengthScore(400); s != 100.0 {
		t.Errorf("expected 100 for 400 len, got %f", s)
	}
	if s := calculateLengthScore(800); s != 100.0 {
		t.Errorf("expected cap at 100 for 800 len, got %f", s)
	}
}

func TestCalculateCompositeScore(t *testing.T) {
	weights := models.ScoreWeights{
		Relevance: 0.5,
		Recency:   0.2,
		Length:    0.2,
		Category:  0.1,
	}

	factors := models.ScoreFactors{
		LLMScore:      90,
		VectorScore:   0.9,
		RecencyDays:   0,   // recency = 100
		ContentLength: 400, // length = 100
		CategoryMatch: true,
		Weights:       weights,
	}

	score := CalculateCompositeScore(factors)
	// 0.5 * 90 + 0.2 * 100 + 0.2 * 100 + 0.1 * 100 = 45 + 20 + 20 + 10 = 95
	if score != 95 {
		t.Errorf("expected composite score 95, got %d", score)
	}

	// Test clamping
	factorsMax := models.ScoreFactors{
		LLMScore:      100,
		VectorScore:   1.0,
		RecencyDays:   0,
		ContentLength: 500,
		CategoryMatch: true,
		Weights:       weights,
	}
	if s := CalculateCompositeScore(factorsMax); s != 100 {
		t.Errorf("expected max score 100, got %d", s)
	}

	// Test title penalty
	factorsPenalty := factors
	factorsPenalty.TitlePenalty = 0.2
	scorePenalty := CalculateCompositeScore(factorsPenalty)
	// 95 * 0.8 = 76
	if scorePenalty != 76 {
		t.Errorf("expected penalized score 76, got %d", scorePenalty)
	}
}

func TestScoreCandidates_VectorOnly(t *testing.T) {
	currentPost := &models.Post{
		Path: "/content/p/main-post/index.md",
		Slug: "main-post",
		FrontMatter: &models.FrontMatter{
			Title: "Election Transparency In Baldwin County",
			Date:  time.Now(),
		},
		Body: "This is a body longer than 400 characters to ensure full length score is given. " +
			"We are discussing local elections, accountability, civic oversight, ballot verification, " +
			"ethics laws, investigative reporting, public records requests, and campaign contributions. " +
			"Public trust depends on truthful verification.",
	}

	candidates := []models.Candidate{
		{Path: "/p/cand-1", Slug: "cand-1", Similarity: 0.85},
		{Path: "/p/cand-2", Slug: "cand-2", Similarity: 0.95},
	}

	req := models.ScoringRequest{
		CurrentPost: currentPost,
		Candidates:  candidates,
		VectorOnly:  true,
	}

	weights := models.ScoreWeights{
		Relevance: 0.7,
		Recency:   0.1,
		Length:    0.2,
	}

	scored, err := ScoreCandidates(context.Background(), nil, req, weights)
	if err != nil {
		t.Fatalf("ScoreCandidates failed: %v", err)
	}

	if len(scored) != 2 {
		t.Fatalf("expected 2 scored articles, got %d", len(scored))
	}

	// cand-2 had higher similarity (0.95), should be rank 1
	if scored[0].Slug != "cand-2" || scored[0].Rank != 1 {
		t.Errorf("expected cand-2 as rank 1, got %+v", scored[0])
	}
	if scored[1].Slug != "cand-1" || scored[1].Rank != 2 {
		t.Errorf("expected cand-1 as rank 2, got %+v", scored[1])
	}
}

type mockRerankClient struct {
	client.LlamaClient
	rerankCalled bool
	lastReq      models.RerankRequest
	response     *models.RerankResponse
}

func (m *mockRerankClient) Rerank(ctx context.Context, req models.RerankRequest) (*models.RerankResponse, error) {
	m.rerankCalled = true
	m.lastReq = req
	return m.response, nil
}

func TestScoreCandidates_Reranker(t *testing.T) {
	currentPost := &models.Post{
		Path: "/content/p/main-post/index.md",
		Slug: "main-post",
		FrontMatter: &models.FrontMatter{
			Title: "Ethics Violations Investigation",
			Date:  time.Now(),
		},
		Body: "Body of main post with sufficient length for testing.",
	}

	candidates := []models.Candidate{
		{Path: "/p/cand-1", Slug: "cand-1", Similarity: 0.70},
		{Path: "/p/cand-2", Slug: "cand-2", Similarity: 0.85},
	}

	req := models.ScoringRequest{
		CurrentPost: currentPost,
		Candidates:  candidates,
		UseReranker: true,
		Model:       "bge-reranker-large",
	}

	// Reranker gives cand-1 (index 0) higher score (0.98) than cand-2 (index 1, 0.40)
	mockClient := &mockRerankClient{
		response: &models.RerankResponse{
			Model: "bge-reranker-large",
			Results: []models.RerankResult{
				{Index: 0, RelevanceScore: 0.98},
				{Index: 1, RelevanceScore: 0.40},
			},
		},
	}

	weights := models.ScoreWeights{
		Relevance: 0.8,
		Recency:   0.1,
		Length:    0.1,
	}

	scored, err := ScoreCandidates(context.Background(), mockClient, req, weights)
	if err != nil {
		t.Fatalf("ScoreCandidates failed: %v", err)
	}

	if !mockClient.rerankCalled {
		t.Errorf("expected LlamaClient.Rerank to be called")
	}

	if len(scored) != 2 {
		t.Fatalf("expected 2 scored articles, got %d", len(scored))
	}

	// cand-1 should have LLMScore 98, cand-2 should have LLMScore 40
	var cand1, cand2 *models.ScoredArticle
	for i := range scored {
		if scored[i].Slug == "cand-1" {
			cand1 = &scored[i]
		}
		if scored[i].Slug == "cand-2" {
			cand2 = &scored[i]
		}
	}

	if cand1 == nil || cand1.LLMScore != 98 {
		t.Errorf("expected cand-1 LLMScore 98, got %+v", cand1)
	}
	if cand2 == nil || cand2.LLMScore != 40 {
		t.Errorf("expected cand-2 LLMScore 40, got %+v", cand2)
	}

	// cand-1 should be rank 1 due to high reranker score despite lower initial similarity
	if scored[0].Slug != "cand-1" || scored[0].Rank != 1 {
		t.Errorf("expected cand-1 rank 1, got %+v", scored[0])
	}
}

