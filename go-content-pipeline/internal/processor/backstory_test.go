package processor

import (
	"testing"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

func TestRankBackstoryRelated(t *testing.T) {
	tNow := time.Now()

	currentPost := &models.Post{
		Slug: "backstory-ep-10",
		FrontMatter: &models.FrontMatter{
			Title:      "Backstory Episode 10",
			Date:       tNow,
			Categories: []string{"Backstory Podcast"},
		},
	}

	allPosts := []*models.Post{
		currentPost,
		{
			Slug: "backstory-ep-9",
			FrontMatter: &models.FrontMatter{
				Title:      "Backstory Episode 9",
				Date:       tNow.Add(-24 * time.Hour),
				Categories: []string{"Backstory Podcast"},
			},
		},
		{
			Slug: "backstory-ep-8",
			FrontMatter: &models.FrontMatter{
				Title:      "Backstory Episode 8",
				Date:       tNow.Add(-48 * time.Hour),
				Categories: []string{"Backstory Podcast"},
			},
		},
		{
			Slug: "regular-news",
			FrontMatter: &models.FrontMatter{
				Title:      "Regular News",
				Date:       tNow.Add(-12 * time.Hour),
				Categories: []string{"News"},
			},
		},
	}

	related := RankBackstoryRelated(currentPost, allPosts)

	// Should only include ep-9 and ep-8 (excluding self and regular news)
	if len(related) != 2 {
		t.Fatalf("expected 2 related backstory posts, got %d", len(related))
	}

	// First should be ep-9 (newest), score 100
	if related[0].Slug != "backstory-ep-9" || related[0].Score != 100 || related[0].Rank != 1 {
		t.Errorf("unexpected first backstory: %+v", related[0])
	}

	// Second should be ep-8, score 99
	if related[1].Slug != "backstory-ep-8" || related[1].Score != 99 || related[1].Rank != 2 {
		t.Errorf("unexpected second backstory: %+v", related[1])
	}
}

func TestFilterBackstoryFromNonBackstory(t *testing.T) {
	posts := []*models.Post{
		{
			Slug: "backstory-1",
			FrontMatter: &models.FrontMatter{
				Categories: []string{"backstory"},
			},
		},
		{
			Slug: "regular-1",
			FrontMatter: &models.FrontMatter{
				Categories: []string{"investigation"},
			},
		},
	}

	candidates := []models.Candidate{
		{Slug: "backstory-1", Similarity: 0.9},
		{Slug: "regular-1", Similarity: 0.8},
		{Slug: "unknown-1", Similarity: 0.7},
	}

	filtered := FilterBackstoryFromNonBackstory(candidates, posts)

	if len(filtered) != 2 {
		t.Fatalf("expected 2 candidates after filtering, got %d", len(filtered))
	}

	for _, cand := range filtered {
		if cand.Slug == "backstory-1" {
			t.Errorf("backstory candidate was not filtered out")
		}
	}
}
