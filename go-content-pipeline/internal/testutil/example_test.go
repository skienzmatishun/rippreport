package testutil_test

import (
	"testing"

	"github.com/rippreport/go-content-pipeline/internal/testutil"
)

// Example test demonstrating helper usage
func TestHelpers(t *testing.T) {
	// Create temporary directory
	dir := testutil.TempDir(t)
	testutil.AssertTrue(t, testutil.FileExists(t, dir), "temp dir should exist")

	// Create temporary file
	content := "test content"
	file := testutil.TempFile(t, content)
	testutil.AssertTrue(t, testutil.FileExists(t, file), "temp file should exist")

	// Read and verify file content
	read := testutil.ReadFile(t, file)
	testutil.AssertEqual(t, read, content, "file content should match")

	// Test hash calculation
	hash1 := testutil.CalculateHash(content)
	hash2 := testutil.CalculateHash(content)
	testutil.AssertEqual(t, hash1, hash2, "hashes should be identical for same content")
}

// Example test demonstrating fixture usage
func TestFixtures(t *testing.T) {
	// Create a Hugo site structure
	dir := testutil.TempDir(t)
	testutil.CreateHugoSiteStructure(t, dir)

	// Create a default post
	post := testutil.DefaultPost("test-article")
	postDir := post.CreatePost(t, dir)
	testutil.AssertTrue(t, testutil.FileExists(t, postDir), "post directory should exist")

	// Create a backstory post with transcript
	backstory := testutil.BackstoryPost("episode-1")
	transcript := "This is the transcript of the podcast episode."
	backstory.CreatePostWithTranscript(t, dir, transcript)

	// Create posts with different characteristics
	_ = testutil.ShortPost("short").CreatePost(t, dir)
	_ = testutil.LongPost("long").CreatePost(t, dir)
	_ = testutil.ElectionPost("election-2024").CreatePost(t, dir)

	// Create post with existing related articles
	related := testutil.PostWithRelated("main", []string{"rel-1", "rel-2", "rel-3"})
	_ = related.CreatePost(t, dir)
	testutil.AssertEqual(t, len(related.RelatedArticles), 3, "should have 3 related articles")
}

// Example test demonstrating test data usage
func TestTestData(t *testing.T) {
	// Generate test embeddings
	emb1 := testutil.TestEmbedding(384, 0.5)
	emb2 := testutil.IdenticalEmbedding(emb1)
	_ = testutil.SimilarEmbedding(emb1, 0.9) // Similar embedding
	emb4 := testutil.DifferentEmbedding(emb1)

	testutil.AssertEqual(t, len(emb1), 384, "embedding should have correct dimension")
	testutil.AssertEqual(t, len(emb2), 384, "identical embedding should have same dimension")

	// Test vector operations
	dot1 := testutil.DotProduct(emb1, emb2)
	testutil.AssertFloatEqual(t, float64(dot1), 1.0, 0.001, "identical embeddings should have dot product ~1")

	dot2 := testutil.DotProduct(emb1, emb4)
	testutil.AssertTrue(t, dot2 < 0.5, "different embeddings should have low similarity")

	// Test score factors
	highFactors := testutil.HighScoreFactors()
	testutil.AssertTrue(t, highFactors.LLMScore > 90, "high factors should have high LLM score")

	mediumFactors := testutil.MediumScoreFactors()
	testutil.AssertTrue(t, mediumFactors.LLMScore > 60 && mediumFactors.LLMScore < 80,
		"medium factors should have medium LLM score")
}

// Example test demonstrating table-driven tests
func TestTableDriven(t *testing.T) {
	// String test example
	stringTests := []testutil.StringTest{
		{
			Name:     "empty string",
			Input:    "",
			Expected: "",
		},
		{
			Name:     "simple string",
			Input:    "hello",
			Expected: "HELLO",
		},
		{
			Name:     "with spaces",
			Input:    "hello world",
			Expected: "HELLO WORLD",
		},
	}

	// Example function to test
	toUpper := func(s string) (string, error) {
		result := ""
		for _, c := range s {
			if c >= 'a' && c <= 'z' {
				result += string(c - 32)
			} else {
				result += string(c)
			}
		}
		return result, nil
	}

	testutil.RunStringTests(t, stringTests, toUpper)

	// Float test example
	floatTests := []testutil.FloatTest{
		{
			Name:     "zero",
			Input:    0.0,
			Expected: 0.0,
			Epsilon:  0.0001,
		},
		{
			Name:     "positive",
			Input:    5.5,
			Expected: 30.25,
			Epsilon:  0.01,
		},
		{
			Name:     "negative",
			Input:    -3.0,
			Expected: 9.0,
			Epsilon:  0.01,
		},
	}

	square := func(x float64) (float64, error) {
		return x * x, nil
	}

	testutil.RunFloatTests(t, floatTests, square)
}

// Example test demonstrating mock server usage
func TestMockServer(t *testing.T) {
	// Create mock llama server
	server := testutil.NewMockLlamaServer(t)

	// Configure custom responses
	server.SetEmbeddingHandler(func(req testutil.EmbeddingRequest) ([]float32, error) {
		// Return embedding based on content length
		dim := len(req.Content)
		if dim > 1000 {
			dim = 384
		}
		return testutil.TestEmbedding(dim, 0.5), nil
	})

	server.SetCompletionHandler(func(req testutil.CompletionRequest) (string, error) {
		// Return score based on prompt content
		if len(req.Prompt) > 500 {
			return "95", nil
		}
		return "75", nil
	})

	// Simulate some requests (in real tests, this would be done by the client)
	// server.URL() provides the base URL for the mock server

	// Verify server behavior
	server.AssertEmbeddingRequestCount(t, 0, "no embedding requests yet")
	server.AssertCompletionRequestCount(t, 0, "no completion requests yet")

	// Test error simulation
	server.SimulateServerError = true
	// In real tests, client requests would fail here

	server.Reset()
	server.AssertEmbeddingRequestCount(t, 0, "requests should be cleared after reset")
}

// Example test demonstrating validation tests
func TestValidation(t *testing.T) {
	tests := []testutil.ValidationTest{
		{
			Name:    "valid positive",
			Input:   10,
			IsValid: true,
		},
		{
			Name:     "invalid negative",
			Input:    -5,
			IsValid:  false,
			ErrorMsg: "must be positive",
		},
		{
			Name:    "valid zero",
			Input:   0,
			IsValid: true,
		},
	}

	validatePositive := func(val interface{}) error {
		num, ok := val.(int)
		if !ok {
			return newSimpleError("not an integer")
		}
		if num < 0 {
			return newSimpleError("must be positive")
		}
		return nil
	}

	testutil.RunValidationTests(t, tests, validatePositive)
}

// Example test demonstrating concurrency tests
func TestConcurrency(t *testing.T) {
	// Shared counter (would use sync.Mutex in real implementation)
	counter := 0

	tests := []testutil.ConcurrencyTest{
		{
			Name:       "concurrent increments",
			Goroutines: 10,
			Operations: 100,
			Run: func(goroutineID, operationID int) error {
				// In real test, would use proper synchronization
				counter++
				return nil
			},
		},
	}

	testutil.RunConcurrencyTests(t, tests)

	// In real test, would verify counter equals Goroutines * Operations
	// This example shows the pattern
}

// Helper types and functions for examples

type simpleError struct {
	message string
}

func (e *simpleError) Error() string {
	return e.message
}

func newSimpleError(msg string) error {
	return &simpleError{message: msg}
}

// Example benchmark using benchmark test utilities
func BenchmarkExample(b *testing.B) {
	tests := []testutil.BenchmarkTest{
		{
			Name: "hash calculation",
			Run: func(b *testing.B) {
				content := "sample content for hashing"
				for i := 0; i < b.N; i++ {
					_ = testutil.CalculateHash(content)
				}
			},
		},
		{
			Name: "vector normalization",
			Run: func(b *testing.B) {
				vec := testutil.TestEmbedding(384, 0.5)
				for i := 0; i < b.N; i++ {
					_ = testutil.NormalizeVector(vec)
				}
			},
		},
	}

	testutil.RunBenchmarks(b, tests)
}
