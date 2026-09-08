# Test Utilities Quick Reference

## Most Common Operations

### File Setup
```go
dir := testutil.TempDir(t)                              // Auto-cleanup temp dir
testutil.WriteFile(t, path, content)                    // Write with mkdir
content := testutil.ReadFile(t, path)                   // Read file
testutil.AssertTrue(t, testutil.FileExists(t, path), msg)
```

### Hugo Posts
```go
// Create test site
dir := testutil.TempDir(t)
testutil.CreateHugoSiteStructure(t, dir)

// Create posts
post := testutil.DefaultPost("my-slug")
backstory := testutil.BackstoryPost("episode-1")
short := testutil.ShortPost("brief")
election := testutil.ElectionPost("vote-2024")

// Write to disk
postDir := post.CreatePost(t, dir)
```

### Embeddings
```go
emb := testutil.TestEmbedding(384, 0.5)                 // Basic vector
similar := testutil.SimilarEmbedding(emb, 0.9)          // 90% similar
different := testutil.DifferentEmbedding(emb)           // Low similarity
dot := testutil.DotProduct(emb1, emb2)                  // Similarity calc
```

### Assertions
```go
testutil.AssertNoError(t, err, "should succeed")
testutil.AssertEqual(t, got, want, "values match")
testutil.AssertFloatEqual(t, got, want, 0.001, "floats")
testutil.AssertContains(t, str, "substring", "contains")
```

### Table Tests
```go
tests := []testutil.StringTest{
    {Name: "case", Input: "in", Expected: "out"},
}
testutil.RunStringTests(t, tests, myFunc)

tests := []testutil.FloatTest{
    {Name: "case", Input: 1.0, Expected: 2.0, Epsilon: 0.01},
}
testutil.RunFloatTests(t, tests, myFunc)
```

### Mock Server
```go
server := testutil.NewMockLlamaServer(t)                // Auto-cleanup
url := server.URL()                                     // For client config

// Track requests
server.AssertEmbeddingRequestCount(t, 5, "5 calls")
reqs := server.GetEmbeddingRequests()

// Simulate errors
server.SimulateServerError = true
server.SimulateTimeout = true
```

### Round-Trip Tests (Req 23.8)
```go
tests := []testutil.RoundTripTest{
    {Name: "case", Input: data},
}
testutil.RunRoundTripTests(t, tests, serialize, deserialize)
```

## Test Patterns by Requirement

| Requirement | Pattern | Example |
|-------------|---------|---------|
| 23.1 | Embedding cache round-trip | `TestEmbedding` + cache |
| 23.2 | Identical vector similarity | `VectorTest` with identical vecs |
| 23.3 | Front matter preservation | `RoundTripTest` with YAML |
| 23.4 | Score bounds | `FloatTest` with 0-100 range |
| 23.5 | Staging equivalence | Compare direct vs staged |
| 23.6 | Cache hash validation | `CalculateHash` + cache test |
| 23.7 | Ranking properties | Verify ranks 1-10 |
| 23.8 | YAML round-trip | `RoundTripTest` |

## Common Test Structure

```go
func TestMyComponent(t *testing.T) {
    // Setup
    dir := testutil.TempDir(t)
    server := testutil.NewMockLlamaServer(t)
    
    // Test data
    post := testutil.DefaultPost("test")
    emb := testutil.TestEmbedding(384, 0.5)
    
    // Execute
    result, err := MyComponent(input)
    
    // Verify
    testutil.AssertNoError(t, err, "should succeed")
    testutil.AssertEqual(t, result, expected, "result matches")
    
    // Cleanup is automatic
}
```

## File Locations

- Helpers: `internal/testutil/helpers.go`
- Fixtures: `internal/testutil/fixtures.go`
- Test Data: `internal/testutil/testdata.go`
- Table Patterns: `internal/testutil/tables.go`
- Mock Server: `internal/testutil/mock_server.go`
- Examples: `internal/testutil/example_test.go`
- Full Docs: `internal/testutil/README.md`
