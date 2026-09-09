# Developer Documentation & Architecture Guide

This document details the architectural design, package responsibilities, concurrency patterns, and testing strategy of the `go-content-pipeline`.

---

## 1. Architectural Overview

The application is structured into decoupled, single-responsibility layers:

```
                  ┌───────────────────────────────┐
                  │          cmd/pipeline         │
                  └───────────────┬───────────────┘
                                  │
                  ┌───────────────▼───────────────┐
                  │          internal/cli         │
                  └───────────────┬───────────────┘
                                  │
                  ┌───────────────▼───────────────┐
                  │     internal/orchestrator     │
                  └───────┬───────────────┬───────┘
                          │               │
            ┌─────────────▼─────┐   ┌─────▼─────────────┐
            │ internal/processor│   │  internal/client  │
            └──────┬────────────┘   └───────────────────┘
                   │
    ┌──────────────┼──────────────┬──────────────┐
    ▼              ▼              ▼              ▼
internal/cache  internal/parser internal/models internal/errors
```

---

## 2. Package Responsibilities

### `cmd/pipeline/`
- Application binary entrypoint.
- Instantiates the CLI application and passes os.Args and context.

### `internal/cli/`
- Command-line parsing, flag extraction, validation, and subcommand routing.
- Signal interception (`SIGINT`, `SIGTERM`) for graceful teardown and progress checkpointing.
- Console reporting and summary display.

### `internal/orchestrator/`
- Central workflow coordinator.
- Sequences stages: Post Discovery → Embedding → Candidate Selection → Scoring → Merge → Staging/Front Matter Update → Hook Generation.
- Coordinates error collection and memory management between phases.

### `internal/processor/`
- **`PostManager`**: Discovers Hugo posts (`content/p/*/index.md`), extracts metadata, parses front matter, preserves markdown body byte-for-byte, and writes atomic updates.
- **`Similarity`**: Implements 0-allocation `CosineSimilarity` using float64 accumulation and min-heap `GetTopCandidates` for efficient top-N candidate selection.
- **`Scoring`**: Calculates recency exponential decay (`100 * exp(-days/365)`), length score (400 chars = 100), category boosts, title penalties, and calls LLM completion endpoint for relevance scores.
- **`Merge`**: Merges new articles with existing recommendations, deduplicating and maintaining sequential ranking up to 10 items.
- **`Backstory`**: Handles backstory/podcast episodes by ranking them by descending release date with simple descending scores.
- **`Compressor`**: Implements 60/40 head-tail truncation, token estimation, and LLM-driven article summarization with cache persistence.
- **`BatchProcessor`**: Packs multiple scoring or hook generation requests into batched prompts to minimize network roundtrips, with automatic fallback to single-item generation on error.
- **`HookGenerator`**: Formats contextual hook prompts and writes `sidebar-hooks.yaml`.
- **`BackupManager`**: Creates timestamped MD5-verified file backups before any destructive modification.
- **`StagingManager`**: Manages `{slug}/related-articles.staged.yaml` lifecycle (save, read, apply, clear).
- **`WorkerPool` & `MemoryManager`**: Bounded goroutine worker pool with error queues and periodic runtime GC invocation (every 10 posts).

### `internal/client/`
- **`LlamaClient`**: High-performance HTTP client for `llama.cpp` server.
- Supports keep-alive connection pooling, token-bucket rate limiting (`TokenBucketRateLimiter`), exponential backoff retry (1s, 2s, 4s) on 5xx/connection errors, immediate fast-fail on 4xx, and request statistics tracking.

### `internal/cache/`
- **`EmbeddingCache`**: In-memory thread-safe map protected by `sync.RWMutex`, validated by post MD5 content hash, invalidated on model name change, and persisted atomically to JSON.

### `internal/parser/`
- **`FrontMatter`**: Extracts YAML front matter delimited by `---`, preserves arbitrary extra fields, handles polymorphic categories (scalar or slice), and produces consistent formatted YAML.

### `internal/models/`
- Strongly-typed domain models (`Post`, `FrontMatter`, `RelatedArticle`, `Hook`, `ScoreWeights`, `ScoreFactors`, etc.) with validation methods.

### `internal/errors/`
- **`ErrorCollector`**: Thread-safe error aggregator that categorizes pipeline failures (`Network`, `FileIO`, `Validation`, `Processing`) and formats pipeline summary reports.

---

## 3. Key Design Decisions

1. **Atomic File Persistence**:
   - All disk writes (front matter updates, caches, staging files, hooks) write first to a `.tmp-*` file in the destination directory, sync to disk via `file.Sync()`, close, and then atomically rename via `os.Rename`. This guarantees zero corrupted or partially written files if the process is killed.

2. **Body Byte Preservation**:
   - When updating front matter, the markdown body is preserved byte-for-byte without altering newlines. This guarantees that `post.ContentHash` (MD5 of body) remains identical, preventing unnecessary cache invalidation on subsequent runs.

3. **Zero-Allocation Hot Path Optimization**:
   - Vector operations in `CosineSimilarity` do not allocate temporary slices. Magnitude sums and dot products are accumulated in a single pass using `float64` for maximum numerical accuracy and performance.

4. **Two-Tier Staging Workflow**:
   - Content authors can review recommendations in staging YAML files before modifying live publication posts.

---

## 4. Testing Strategy

The test suite is organized into three tiers:

1. **Unit Tests**:
   - Every package contains comprehensive unit tests testing edge cases, boundary conditions, invalid inputs, and error handling.
   - Run: `go test -v ./internal/...`

2. **Integration & Workflow Tests**:
   - Located in `test/`, utilizing `testutil.MockLlamaServer` and temporary Hugo directory fixtures:
     - `embedding_workflow_test.go`: Cache persistence, cache hit validation, model invalidation.
     - `ranking_workflow_test.go`: Front matter writing, backstory separation, vector-only mode.
     - `hook_workflow_test.go`: Sidebar hook YAML structure, compression caching.
     - `staging_workflow_test.go`: Staging lifecycle (stage -> hooks -> apply -> backups -> clear).
     - `refresh_recent_test.go`: Merge algorithm on live posts.
     - `error_recovery_test.go`: Interrupted resume, API retries, corrupted cache recovery, invalid post tolerance.
     - `pipeline_e2e_test.go`: Full pipeline end-to-end execution.
     - `large_dataset_test.go`: 100 posts execution validating heap memory strictly under 500MB.
   - Run: `go test -v ./test/...`

3. **Benchmarks**:
   - Located in `internal/processor/bench_test.go` and `internal/parser/bench_test.go`.
   - Measure cosine similarity, top-N min-heap selection, composite scoring, and YAML parsing/serialization.
   - Run: `go test -bench=. -run=^$ ./internal/processor ./internal/parser`
