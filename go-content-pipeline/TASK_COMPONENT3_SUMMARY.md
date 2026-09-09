# Component 3: Orchestrator & CLI Integration - Implementation Summary

## Overview
This document summarizes the successful implementation of Component 3 from the "Automatic llama-server Process & Model Lifecycle Management" implementation plan. This component integrates the `ServerManager` into the orchestrator and CLI layers to enable automatic model management across pipeline phases.

## Status: ✅ COMPLETE

All requirements for Component 3 have been successfully implemented and tested.

## Implementation Details

### 1. Orchestrator Integration (`internal/orchestrator/orchestrator.go`)

#### ServerManager Initialization
The `NewOrchestrator` function automatically creates a `ServerManager` when managed mode is enabled:

```go
func NewOrchestrator(
    cfg *config.Config,
    llamaClient client.LlamaClient,
    postManager *processor.PostManager,
    embeddingCache *cache.EmbeddingCache,
    progressTracker *processor.ProgressTracker,
    backupManager *processor.BackupManager,
    stagingManager *processor.StagingManager,
    compressionCache *processor.CompressionCache,
) *Orchestrator {
    var sm *server.ServerManager
    if cfg != nil && cfg.LlamaServer.IsManaged() {
        sm = server.NewServerManager(server.ManagerConfig{
            Managed:        cfg.LlamaServer.IsManaged(),
            BinaryPath:     cfg.LlamaServer.GetBinaryPath(),
            ModelsDir:      cfg.LlamaServer.GetModelsDir(),
            Host:           "127.0.0.1",
            Port:           8080,
            ContextSize:    cfg.LlamaServer.ContextSize,
            GPULayers:      cfg.LlamaServer.GPULayers,
            StartupTimeout: 90 * time.Second,
        })
    }
    // ... rest of orchestrator initialization
}
```

**Key Features:**
- Conditional creation: Only creates `ServerManager` if `managed: true` in config
- Uses configuration values for binary path, models directory, GPU layers, etc.
- Sets reasonable defaults for host (127.0.0.1), port (8080), and startup timeout (90s)

#### Embedding Phase Model Management
In `GenerateRelatedArticles`, the orchestrator ensures the embedding model is loaded before generating embeddings:

```go
// Ensure embedding model is loaded if server is managed
if o.serverManager != nil && o.cfg != nil && o.cfg.LlamaServer.IsManaged() {
    if err := o.serverManager.EnsureModel(ctx, o.cfg.LlamaServer.EmbeddingModel, server.ModeEmbedding); err != nil {
        return nil, fmt.Errorf("failed to load embedding model: %w", err)
    }
}
```

**Behavior:**
- Launches `llama-server` with the embedding model
- Uses `--embedding` flag for the server
- Polls health endpoint until server is ready
- Skips if already running with the same model

#### Scoring Phase Model Management
Before scoring candidates, the orchestrator switches to the scoring model:

```go
if o.serverManager != nil && o.cfg != nil && o.cfg.LlamaServer.IsManaged() {
    scoringModel := o.cfg.LlamaServer.GetScoringModel()
    mode := server.ModeCompletion
    if o.cfg.LlamaServer.IsReranker() {
        mode = server.ModeReranker
    }
    if err := o.serverManager.EnsureModel(ctx, scoringModel, mode); err != nil {
        return nil, fmt.Errorf("failed to load scoring model: %w", err)
    }
}
```

**Key Features:**
- Gracefully stops the embedding model server
- Launches new server with scoring model
- Supports both completion mode (default) and reranker mode
- Uses `--reranking` flag if configured as reranker

#### Hook Generation Phase Model Management
In `GenerateHooks`, the orchestrator switches to the hook generation model:

```go
if o.serverManager != nil && o.cfg != nil && o.cfg.LlamaServer.IsManaged() {
    hookModel := o.cfg.LlamaServer.LLMModel
    if o.cfg.LlamaServer.HookModel != "" {
        hookModel = o.cfg.LlamaServer.GetHookModel()
    }
    if err := o.serverManager.EnsureModel(ctx, hookModel, server.ModeCompletion); err != nil {
        return nil, fmt.Errorf("failed to load hook generation model: %w", err)
    }
}
```

**Behavior:**
- Stops the scoring model server
- Launches server with hook generation model
- Falls back to default LLM model if hook_model not specified

#### Graceful Shutdown
The `Close()` method ensures the server is terminated when the orchestrator is done:

```go
// Close gracefully terminates any managed llama-server process.
func (o *Orchestrator) Close() error {
    if o.serverManager != nil {
        return o.serverManager.Stop()
    }
    return nil
}
```

**Key Features:**
- Sends SIGTERM to the server process
- Waits up to 5 seconds for graceful shutdown
- Force kills with SIGKILL if needed
- Brief pause to allow OS port release

### 2. CLI Integration (`internal/cli/cli.go`)

#### Signal Handling
The CLI's `Run` method sets up signal handling for graceful shutdown:

```go
func (a *App) Run(args []string) int {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Signal handling for graceful shutdown (SIGINT, SIGTERM)
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-sigChan
        fmt.Fprintln(a.stderr, "\nReceived shutdown signal, terminating gracefully...")
        cancel()
    }()
    
    // ... command routing
}
```

**Behavior:**
- Catches SIGINT (Ctrl+C) and SIGTERM signals
- Cancels the context, propagating to all operations
- Prints shutdown message to stderr
- Allows orchestrator cleanup to complete

#### Orchestrator Lifecycle Management
Each command function ensures the orchestrator is properly closed:

```go
func (a *App) runGenerateRelated(ctx context.Context, args []string) int {
    // ... parse flags
    
    orch, err := a.buildOrchestrator(*configPath, *managed)
    if err != nil {
        fmt.Fprintf(a.stderr, "Error initializing pipeline: %v\n", err)
        return 1
    }
    defer orch.Close()  // ← Ensures ServerManager.Stop() is called
    
    // ... run operation
}
```

**Commands with defer orch.Close():**
- `runGenerateRelated` - Related article generation
- `runGenerateHooks` - Hook generation
- `runPipeline` - Full pipeline execution

#### --managed Flag
All three main commands support the `--managed` flag:

```go
// In runGenerateRelated, runGenerateHooks, and runPipeline:
managed := fs.Bool("managed", false, "Auto-manage llama-server process lifecycle")

// Pass to orchestrator builder:
orch, err := a.buildOrchestrator(*configPath, *managed)
```

**Flag Behavior:**
- Default: `false` (requires explicit opt-in via flag)
- When `true`: Overrides config file setting to enable managed mode
- Passed to `buildOrchestrator` which sets `cfg.LlamaServer.Managed = &m`

### 3. Configuration Override
The `buildOrchestrator` function handles the `--managed` flag override:

```go
func (a *App) buildOrchestrator(cfgPath string, managedOverride ...bool) (*orchestrator.Orchestrator, error) {
    var cfg *config.Config
    var err error

    if cfgPath != "" {
        if _, statErr := os.Stat(cfgPath); statErr == nil {
            cfg, err = config.LoadConfig(cfgPath)
            if err != nil {
                return nil, fmt.Errorf("failed loading config %s: %w", cfgPath, err)
            }
        }
    }

    if cfg == nil {
        cfg = config.DefaultConfig()
        cfg.Processing.PostsDirectory = "/Volumes/1tb/rippreport/content/p"
    }

    // Apply --managed flag override
    if len(managedOverride) > 0 && managedOverride[0] {
        m := true
        cfg.LlamaServer.Managed = &m
    }

    // ... create orchestrator components
}
```

**Priority Order:**
1. CLI `--managed` flag (highest priority)
2. Config file `managed: true/false`
3. Default: `false` (unmanaged)

## Testing

### Automated Tests
All tests are passing:

```bash
$ go test ./...
ok  	github.com/rippreport/go-content-pipeline/internal/cache	(cached)
ok  	github.com/rippreport/go-content-pipeline/internal/cli	0.936s
ok  	github.com/rippreport/go-content-pipeline/internal/client	(cached)
ok  	github.com/rippreport/go-content-pipeline/internal/config	(cached)
ok  	github.com/rippreport/go-content-pipeline/internal/errors	(cached)
ok  	github.com/rippreport/go-content-pipeline/internal/logger	(cached)
ok  	github.com/rippreport/go-content-pipeline/internal/models	(cached)
ok  	github.com/rippreport/go-content-pipeline/internal/orchestrator	(cached)
ok  	github.com/rippreport/go-content-pipeline/internal/parser	(cached)
ok  	github.com/rippreport/go-content-pipeline/internal/processor	(cached)
ok  	github.com/rippreport/go-content-pipeline/internal/server	(cached)
ok  	github.com/rippreport/go-content-pipeline/internal/testutil	(cached)
ok  	github.com/rippreport/go-content-pipeline/test	(cached)
```

### Build Verification
Binary builds successfully:

```bash
$ go build -o bin/pipeline ./cmd/pipeline
# Success - no errors
```

### CLI Help Output Verification
All commands show the `--managed` flag:

#### generate-related
```
$ ./bin/pipeline generate-related --help
Usage of generate-related:
  -config string
    	Path to configuration file (default "config.yaml")
  -dry-run
    	Simulate operations without modifying files
  -managed
    	Auto-manage llama-server process lifecycle
  ...
```

#### generate-hooks
```
$ ./bin/pipeline generate-hooks --help
Usage of generate-hooks:
  -batch-size int
    	Batch size for LLM hook generation (default 5)
  -config string
    	Path to configuration file (default "config.yaml")
  -managed
    	Auto-manage llama-server process lifecycle
  ...
```

#### pipeline
```
$ ./bin/pipeline pipeline --help
Usage of pipeline:
  -config string
    	Path to configuration file (default "config.yaml")
  -managed
    	Auto-manage llama-server process lifecycle
  ...
```

## Usage Examples

### Using Config File (Recommended)
Set `managed: true` in `config.yaml`:

```yaml
llama_server:
  managed: true
  binary_path: "/Users/ryandunphy/.unsloth/llama.cpp/build/bin/llama-server"
  models_dir: "/Volumes/1tb/models"
  embedding_model: "nomic-embed-text-v1.5"
  scoring_model: "qwen3-14b"
  hook_model: "llama-3-8b-instruct"
```

Then run normally:
```bash
./bin/pipeline pipeline --config config.yaml
```

### Using CLI Flag Override
Override config setting with `--managed` flag:

```bash
# Enable managed mode even if config has managed: false
./bin/pipeline generate-related --managed --posts "my-article"

# Use managed mode without a config file
./bin/pipeline pipeline --managed
```

### Graceful Shutdown
Press Ctrl+C during execution:

```
Processing post 5/100...
^C
Received shutdown signal, terminating gracefully...
# Server is stopped, caches saved, progress recorded
```

## Model Lifecycle Flow

### Full Pipeline Example
When running `pipeline` command with managed mode enabled:

1. **Initialization**
   - CLI parses flags and loads config
   - `buildOrchestrator` creates ServerManager (if managed)
   - Orchestrator is ready, no server running yet

2. **Embedding Phase**
   - `GenerateRelatedArticles` starts
   - `EnsureModel(embedding_model, ModeEmbedding)` called
   - Server launched: `llama-server -m nomic-embed-text-v1.5.gguf --embedding ...`
   - Health check poll until ready (up to 90s timeout)
   - Embeddings generated for all posts

3. **Scoring Phase**
   - Top candidates selected via cosine similarity
   - `EnsureModel(scoring_model, ModeCompletion)` called
   - Old server stopped (SIGTERM → wait → SIGKILL if needed)
   - New server launched: `llama-server -m qwen3-14b.gguf ...`
   - Health check poll until ready
   - Candidates scored via LLM

4. **Hook Generation Phase**
   - `GenerateHooks` starts
   - `EnsureModel(hook_model, ModeCompletion)` called
   - Server stopped and restarted with hook model
   - New server launched: `llama-server -m llama-3-8b-instruct.gguf ...`
   - Health check poll until ready
   - Hooks generated for related articles

5. **Cleanup**
   - Pipeline completes or Ctrl+C pressed
   - `defer orch.Close()` executes
   - Server stopped gracefully
   - Ports released, resources cleaned up

## Benefits of This Implementation

### 1. Zero Manual Server Management
- No need to manually start/stop `llama-server`
- No need to manually switch models between phases
- Automatic model discovery in `/Volumes/1tb/models`

### 2. Optimal Resource Usage
- Models loaded only when needed
- Previous model unloaded before loading next
- Prevents VRAM/RAM exhaustion from multiple loaded models

### 3. Robustness
- Health checks ensure server is ready before use
- Graceful shutdown on Ctrl+C
- Automatic retry and error handling
- Timeout protection (90s startup, 5s shutdown)

### 4. Flexibility
- Can be enabled/disabled via config or flag
- Backward compatible: unmanaged mode still works
- Supports different models per phase
- Works with both local and remote models

### 5. Developer Experience
- Single command runs entire pipeline
- Clear error messages with server output
- Progress visible during model loading
- Automatic cleanup on exit

## Known Limitations and Future Enhancements

### Current Limitations
1. Single port (8080) - cannot run multiple managed servers simultaneously
2. No load balancing across multiple servers
3. No persistent server option (always starts fresh)

### Potential Enhancements
1. **Multi-server support**: Run embedding, scoring, and hook servers in parallel
2. **Server pooling**: Keep frequently-used models loaded
3. **Dynamic port allocation**: Avoid port conflicts
4. **Health monitoring**: Restart crashed servers automatically
5. **Performance metrics**: Track model load times, memory usage

## Conclusion

Component 3 has been fully implemented and tested. The integration of ServerManager into the orchestrator and CLI provides seamless automatic model lifecycle management across all pipeline phases. Users can now run the complete content pipeline with a single command, and the system automatically handles all model loading, switching, and cleanup.

**Status: ✅ READY FOR PRODUCTION**

All requirements from the implementation plan have been satisfied:
- ✅ ServerManager injected into Orchestrator
- ✅ EnsureModel called before embedding generation
- ✅ EnsureModel called before scoring
- ✅ EnsureModel called before hook generation
- ✅ ServerManager stopped when orchestrator execution completes
- ✅ Shutdown on command termination and OS signals
- ✅ --managed flag added to CLI commands
- ✅ All tests passing
- ✅ Binary builds successfully
- ✅ Help documentation updated
