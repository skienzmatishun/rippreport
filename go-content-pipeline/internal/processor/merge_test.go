package processor

import (
	"testing"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

func TestMergeWithExisting(t *testing.T) {
	newArticles := []models.RelatedArticle{
		{Slug: "new-1", Title: "New 1", Score: 95, Rank: 1},
		{Slug: "new-2", Title: "New 2", Score: 90, Rank: 2},
		{Slug: "shared-dup", Title: "Shared Dup New", Score: 85, Rank: 3},
	}

	existingArticles := []models.RelatedArticle{
		{Slug: "shared-dup", Title: "Shared Dup Old", Score: 80, Rank: 1},
		{Slug: "old-1", Title: "Old 1", Score: 75, Rank: 2},
		{Slug: "old-2", Title: "Old 2", Score: 70, Rank: 3},
		{Slug: "old-3", Title: "Old 3", Score: 65, Rank: 4},
	}

	merged := MergeWithExisting(newArticles, existingArticles, 5)

	// Max 5 articles
	if len(merged) != 5 {
		t.Fatalf("expected 5 merged articles, got %d", len(merged))
	}

	// First 3 should be new articles
	if merged[0].Slug != "new-1" || merged[0].Rank != 1 {
		t.Errorf("expected new-1 at rank 1, got %+v", merged[0])
	}
	if merged[1].Slug != "new-2" || merged[1].Rank != 2 {
		t.Errorf("expected new-2 at rank 2, got %+v", merged[1])
	}
	if merged[2].Slug != "shared-dup" || merged[2].Score != 85 || merged[2].Rank != 3 {
		t.Errorf("expected new version of shared-dup at rank 3, got %+v", merged[2])
	}

	// Next 2 should be old-1 and old-2 (shared-dup deduplicated, old-3 truncated)
	if merged[3].Slug != "old-1" || merged[3].Rank != 4 {
		t.Errorf("expected old-1 at rank 4, got %+v", merged[3])
	}
	if merged[4].Slug != "old-2" || merged[4].Rank != 5 {
		t.Errorf("expected old-2 at rank 5, got %+v", merged[4])
	}
}
