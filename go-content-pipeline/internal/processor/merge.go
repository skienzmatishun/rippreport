package processor

import (
	"github.com/rippreport/go-content-pipeline/internal/models"
)

const defaultMaxArticles = 10

// MergeWithExisting merges newly generated article recommendations with existing ones.
// New articles are placed first, followed by existing articles not present in the new set,
// capped at maxArticles (default 10), with ranks recalculated sequentially (1 to N).
//
// Requirements: 11.1, 11.2, 11.3, 11.4, 11.5, 11.6, 11.8
func MergeWithExisting(newArticles, existingArticles []models.RelatedArticle, maxArticles int) []models.RelatedArticle {
	if maxArticles <= 0 {
		maxArticles = defaultMaxArticles
	}

	newSlugs := make(map[string]struct{}, len(newArticles))
	var merged []models.RelatedArticle

	// Place new articles first
	for _, article := range newArticles {
		if article.Slug == "" {
			continue
		}
		if _, exists := newSlugs[article.Slug]; !exists {
			newSlugs[article.Slug] = struct{}{}
			merged = append(merged, article)
		}
	}

	// Preserve existing articles not in new set
	for _, article := range existingArticles {
		if article.Slug == "" {
			continue
		}
		if _, exists := newSlugs[article.Slug]; !exists {
			newSlugs[article.Slug] = struct{}{}
			merged = append(merged, article)
		}
	}

	// Limit to maxArticles
	if len(merged) > maxArticles {
		merged = merged[:maxArticles]
	}

	// Update ranks 1 to N
	for i := range merged {
		merged[i].Rank = i + 1
	}

	return merged
}
