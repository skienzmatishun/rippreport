package processor

import (
	"container/heap"
	"math"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

// normalize scales a vector to unit length.
// Returns a zero vector if magnitude is zero.
//
// Requirements: 4.3
func normalize(v []float32) []float32 {
	var sum float64
	for _, val := range v {
		sum += float64(val) * float64(val)
	}
	mag := float32(math.Sqrt(sum))
	if mag == 0 {
		return make([]float32, len(v))
	}

	norm := make([]float32, len(v))
	for i, val := range v {
		norm[i] = val / mag
	}
	return norm
}

// CosineSimilarity calculates the cosine similarity between two float32 vectors.
// Normalizes both vectors and computes the dot product, clamped to [-1.0, 1.0].
//
// Requirements: 4.1, 4.4, 23.2
func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	normA := normalize(a)
	normB := normalize(b)

	var dot float32
	for i := range normA {
		dot += normA[i] * normB[i]
	}

	if dot > 1.0 {
		return 1.0
	}
	if dot < -1.0 {
		return -1.0
	}
	return dot
}

// candidateHeap implements a min-heap for Candidate based on similarity.
type candidateHeap []models.Candidate

func (h candidateHeap) Len() int           { return len(h) }
func (h candidateHeap) Less(i, j int) bool {
	if h[i].Similarity == h[j].Similarity {
		return h[i].Path > h[j].Path // Inverted for min-heap so smaller path stays
	}
	return h[i].Similarity < h[j].Similarity
}
func (h candidateHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *candidateHeap) Push(x interface{}) {
	*h = append(*h, x.(models.Candidate))
}
func (h *candidateHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

// GetTopCandidates finds the top N most similar candidates from a corpus of embeddings.
// Posts in the exclude list are omitted. Ties are stably broken by path ascending.
//
// Requirements: 4.2, 4.5, 4.6, 4.7
func GetTopCandidates(query []float32, corpus map[string][]float32, n int, exclude []string) []models.Candidate {
	if n <= 0 || len(query) == 0 || len(corpus) == 0 {
		return nil
	}

	excludeSet := make(map[string]struct{}, len(exclude))
	for _, ex := range exclude {
		excludeSet[ex] = struct{}{}
		// Also normalize slug and path
		excludeSet[filepath.Base(filepath.Dir(ex))] = struct{}{}
		excludeSet[filepath.Base(ex)] = struct{}{}
	}

	h := &candidateHeap{}
	heap.Init(h)

	for path, vec := range corpus {
		slug := extractSlugFromPath(path)
		if slug == "" || slug == "." {
			slug = path
		}

		if _, excluded := excludeSet[path]; excluded {
			continue
		}
		if _, excluded := excludeSet[slug]; excluded {
			continue
		}

		sim := CosineSimilarity(query, vec)
		cand := models.Candidate{
			Path:       path,
			Slug:       slug,
			Similarity: sim,
		}

		if h.Len() < n {
			heap.Push(h, cand)
		} else if sim > (*h)[0].Similarity || (sim == (*h)[0].Similarity && path < (*h)[0].Path) {
			heap.Pop(h)
			heap.Push(h, cand)
		}
	}

	result := make([]models.Candidate, h.Len())
	for i := len(result) - 1; i >= 0; i-- {
		result[i] = heap.Pop(h).(models.Candidate)
	}

	// Sort results descending by similarity, then ascending by path for ties
	sort.Slice(result, func(i, j int) bool {
		if result[i].Similarity == result[j].Similarity {
			return strings.Compare(result[i].Path, result[j].Path) < 0
		}
		return result[i].Similarity > result[j].Similarity
	})

	return result
}
