package processor

import (
	"math"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	t.Run("identical vectors return 1.0", func(t *testing.T) {
		a := []float32{1.0, 2.0, 3.0}
		b := []float32{1.0, 2.0, 3.0}
		sim := CosineSimilarity(a, b)
		if math.Abs(float64(sim-1.0)) > 1e-5 {
			t.Errorf("expected 1.0, got %f", sim)
		}
	})

	t.Run("orthogonal vectors return 0.0", func(t *testing.T) {
		a := []float32{1.0, 0.0}
		b := []float32{0.0, 1.0}
		sim := CosineSimilarity(a, b)
		if math.Abs(float64(sim)) > 1e-5 {
			t.Errorf("expected 0.0, got %f", sim)
		}
	})

	t.Run("opposite vectors return -1.0", func(t *testing.T) {
		a := []float32{1.0, 2.0}
		b := []float32{-1.0, -2.0}
		sim := CosineSimilarity(a, b)
		if math.Abs(float64(sim-(-1.0))) > 1e-5 {
			t.Errorf("expected -1.0, got %f", sim)
		}
	})

	t.Run("zero magnitude vector handles gracefully", func(t *testing.T) {
		a := []float32{0.0, 0.0}
		b := []float32{1.0, 2.0}
		sim := CosineSimilarity(a, b)
		if sim != 0.0 {
			t.Errorf("expected 0.0 for zero vector, got %f", sim)
		}
	})
}

func TestGetTopCandidates(t *testing.T) {
	query := []float32{1.0, 0.0, 0.0}
	corpus := map[string][]float32{
		"/content/p/post-exact/index.md":  {1.0, 0.0, 0.0}, // sim = 1.0
		"/content/p/post-close/index.md":  {0.9, 0.1, 0.0}, // sim ~ 0.99
		"/content/p/post-ortho/index.md":  {0.0, 1.0, 0.0}, // sim = 0.0
		"/content/p/post-opp/index.md":    {-1.0, 0.0, 0.0}, // sim = -1.0
		"/content/p/post-exclude/index.md": {1.0, 0.0, 0.0}, // should be excluded
	}

	exclude := []string{"post-exclude"}
	candidates := GetTopCandidates(query, corpus, 3, exclude)

	if len(candidates) != 3 {
		t.Fatalf("expected 3 candidates, got %d", len(candidates))
	}

	for _, c := range candidates {
		if c.Slug == "post-exclude" {
			t.Errorf("excluded post should not be returned")
		}
	}

	// First should be post-exact (1.0)
	if candidates[0].Slug != "post-exact" {
		t.Errorf("expected first candidate to be post-exact, got %s", candidates[0].Slug)
	}

	// Tie breaking stability
	corpusTies := map[string][]float32{
		"/content/p/post-z/index.md": {1.0, 0.0},
		"/content/p/post-a/index.md": {1.0, 0.0},
	}
	ties := GetTopCandidates([]float32{1.0, 0.0}, corpusTies, 2, nil)
	if len(ties) != 2 {
		t.Fatalf("expected 2 ties, got %d", len(ties))
	}
	if ties[0].Slug != "post-a" || ties[1].Slug != "post-z" {
		t.Errorf("expected stable alphabetical tie breaking [post-a, post-z], got [%s, %s]", ties[0].Slug, ties[1].Slug)
	}
}
