package processor

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

// Requirement 23.1: Benchmark cosine similarity calculation
func BenchmarkCosineSimilarity(b *testing.B) {
	dimensions := []int{384, 768, 1536}

	for _, dim := range dimensions {
		b.Run(fmt.Sprintf("dim-%d", dim), func(b *testing.B) {
			v1 := make([]float32, dim)
			v2 := make([]float32, dim)
			for i := 0; i < dim; i++ {
				v1[i] = float32(i) * 0.01
				v2[i] = float32(dim-i) * 0.01
			}

			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = CosineSimilarity(v1, v2)
			}
		})
	}
}

// Requirement 23.1: Benchmark top-N selection with min-heap
func BenchmarkGetTopCandidates(b *testing.B) {
	corpusSizes := []int{100, 500, 1000}
	dim := 384

	for _, size := range corpusSizes {
		b.Run(fmt.Sprintf("corpus-%d", size), func(b *testing.B) {
			corpus := make(map[string][]float32, size)
			for i := 0; i < size; i++ {
				vec := make([]float32, dim)
				for d := 0; d < dim; d++ {
					vec[d] = float32(d+i) * 0.001
				}
				corpus[fmt.Sprintf("post-%d", i)] = vec
			}

			query := make([]float32, dim)
			for d := 0; d < dim; d++ {
				query[d] = 0.5
			}

			exclude := []string{"post-0", "post-1"}

			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = GetTopCandidates(query, corpus, 10, exclude)
			}
		})
	}
}

// Requirement 23.1: Benchmark composite scoring
func BenchmarkCalculateCompositeScore(b *testing.B) {
	weights := models.ScoreWeights{
		Relevance: 0.5,
		Recency:   0.2,
		Length:    0.15,
		Category:  0.15,
	}

	factors := models.ScoreFactors{
		LLMScore:      85,
		RecencyDays:   45,
		ContentLength: 850,
		CategoryMatch: true,
		TitlePenalty:  0.0,
		Weights:       weights,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = CalculateCompositeScore(factors)
	}
}

// Benchmark recency score exponential decay
func BenchmarkCalculateRecencyScore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = calculateRecencyScore(i % 1000)
	}
}

// Benchmark length score calculation
func BenchmarkCalculateLengthScore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = calculateLengthScore(i % 2000)
	}
}

// Benchmark candidate scoring inner loop with 50 candidates
func BenchmarkCandidateScoringPipeline(b *testing.B) {
	rng := rand.New(rand.NewSource(42))
	dim := 384
	query := make([]float32, dim)
	for i := 0; i < dim; i++ {
		query[i] = rng.Float32()
	}

	corpus := make(map[string][]float32, 500)
	for i := 0; i < 500; i++ {
		v := make([]float32, dim)
		for d := 0; d < dim; d++ {
			v[d] = rng.Float32()
		}
		corpus[fmt.Sprintf("article-%d", i)] = v
	}

	weights := models.ScoreWeights{
		Relevance: 0.6,
		Recency:   0.2,
		Length:    0.1,
		Category:  0.1,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		candidates := GetTopCandidates(query, corpus, 10, nil)
		for _, c := range candidates {
			_ = CalculateCompositeScore(models.ScoreFactors{
				LLMScore:      int(c.Similarity * 100),
				RecencyDays:   30,
				ContentLength: 600,
				CategoryMatch: false,
				Weights:       weights,
			})
		}
	}
}
