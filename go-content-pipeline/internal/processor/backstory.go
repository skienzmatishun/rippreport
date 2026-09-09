package processor

import (
	"sort"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

// RankBackstoryRelated ranks related articles for a backstory post.
// It selects only recent backstory posts sorted by date descending,
// assigns simple descending scores (100 down to 91), and skips LLM scoring.
//
// Requirements: 18.4, 18.5, 18.6, 18.7
func RankBackstoryRelated(currentPost *models.Post, allPosts []*models.Post) []models.RelatedArticle {
	var backstories []*models.Post
	for _, p := range allPosts {
		if p == nil || p.Slug == currentPost.Slug {
			continue
		}
		if p.IsBackstory() {
			backstories = append(backstories, p)
		}
	}

	// Sort by date descending (newest first)
	sort.Slice(backstories, func(i, j int) bool {
		return backstories[i].FrontMatter.Date.After(backstories[j].FrontMatter.Date)
	})

	limit := 10
	if len(backstories) < limit {
		limit = len(backstories)
	}

	related := make([]models.RelatedArticle, limit)
	for i := 0; i < limit; i++ {
		score := 100 - i
		if score < 91 {
			score = 91
		}
		related[i] = models.RelatedArticle{
			Slug:  backstories[i].Slug,
			Title: backstories[i].FrontMatter.Title,
			Score: score,
			Rank:  i + 1,
		}
	}

	return related
}

// FilterBackstoryFromNonBackstory removes backstory posts from candidate recommendations
// for non-backstory articles.
//
// Requirements: 18.8
func FilterBackstoryFromNonBackstory(candidates []models.Candidate, posts []*models.Post) []models.Candidate {
	backstorySlugs := make(map[string]struct{})
	for _, p := range posts {
		if p != nil && p.IsBackstory() {
			backstorySlugs[p.Slug] = struct{}{}
		}
	}

	var filtered []models.Candidate
	for _, cand := range candidates {
		if _, isBackstory := backstorySlugs[cand.Slug]; !isBackstory {
			filtered = append(filtered, cand)
		}
	}

	return filtered
}
