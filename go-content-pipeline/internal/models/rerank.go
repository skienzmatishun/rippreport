package models

import "fmt"

// RerankRequest represents a request to score relevance between a query and multiple documents using a cross-encoder reranker.
type RerankRequest struct {
	Model     string   `json:"model,omitempty"`
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
	TopN      int      `json:"top_n,omitempty"`
}

// Validate checks that Query and Documents are populated.
func (r *RerankRequest) Validate() error {
	if r.Query == "" {
		return fmt.Errorf("rerank query cannot be empty")
	}
	if len(r.Documents) == 0 {
		return fmt.Errorf("rerank documents cannot be empty")
	}
	return nil
}

// RerankResult represents a single scored document from the reranker.
type RerankResult struct {
	Index          int     `json:"index"`
	RelevanceScore float64 `json:"relevance_score"`
	Document       string  `json:"document,omitempty"`
}

// RerankResponse represents the response from the reranking endpoint.
type RerankResponse struct {
	Model   string         `json:"model,omitempty"`
	Results []RerankResult `json:"results"`
}
