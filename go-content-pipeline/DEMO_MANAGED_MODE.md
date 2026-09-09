# Demonstration: Automatic llama-server Management

This document demonstrates the automatic llama-server process and model lifecycle management feature.

## Quick Start Demo

### 1. Create a test configuration

```bash
cat > demo-config.yaml << 'EOF'
llama_server:
  managed: true
  binary_path: "/Users/ryandunphy/.unsloth/llama.cpp/build/bin/llama-server"
  models_dir: "/Volumes/1tb/models"
  base_url: "http://localhost:8080"
  embedding_model: "nomic-embed-text-v1.5"
  llm_model: "qwen3-14b"
  scoring_model: "qwen3-14b"
  hook_model: "llama-3-8b-instruct"
  timeout: 120s
  min_request_delay: 100ms
  max_retries: 3

processing:
  posts_directory: "/Volumes/1tb/rippreport/content/p"
  top_candidates: 20
  batch_size: 5
  max_content_length: 100000

storage:
  backup_directory: "backups"
  cache_file: ".embedding_cache.json"
  progress_file: ".progress.json"
  compression_cache: "internal/cache/compressed_articles.json"

scoring:
  weights:
    relevance: 0.7
    recency: 0.15
    length: 0.1
    category: 0.05

filtering:
  exclude_categories:
    - holiday

logging:
  level: "INFO"
  file: "pipeline.log"
  console: true

concurrency:
  max_workers: 4
  max_api_requests: 2
  gc_interval_posts: 10
EOF
```

### 2. Run with managed mode (via config)

```bash
# The pipeline will automatically:
# 1. Start llama-server with embedding model
# 2. Generate embeddings
# 3. Stop server, restart with scoring model
# 4. Score candidates
# 5. Stop server, restart with hook model
# 6. Generate hooks
# 7. Clean up and stop server

./bin/pipeline pipeline --config demo-config.yaml --posts "example-post" --dry-run
```

### 3. Run with managed mode (via flag override)

```bash
# Override config setting with CLI flag
./bin/pipeline generate-related --managed --posts "example-post" --dry-run
```

### 4. Test graceful shutdown

```bash
# Start a longer-running operation
./bin/pipeline pipeline --managed --posts "post1,post2,post3"

# Press Ctrl+C during execution
# You should see: "Received shutdown signal, terminating gracefully..."
# The server will be stopped automatically
```

## What Happens Behind the Scenes

### Phase 1: Embedding Generation
```
[CLI] User runs: ./bin/pipeline pipeline --managed
  ↓
[Orchestrator] NewOrchestrator creates ServerManager (managed=true)
  ↓
[Orchestrator] GenerateRelatedArticles starts
  ↓
[ServerManager] EnsureModel("nomic-embed-text-v1.5", ModeEmbedding)
  ↓
[ServerManager] Launches: llama-server -m /Volumes/1tb/models/nomic-embed-text-v1.5.gguf \
                          --host 127.0.0.1 --port 8080 -c 4096 -ngl 99 --embedding
  ↓
[ServerManager] Polls GET http://127.0.0.1:8080/health until status 200
  ↓
[Orchestrator] Generates embeddings for all posts via API
```

### Phase 2: Scoring
```
[Orchestrator] Top candidates selected via cosine similarity
  ↓
[ServerManager] EnsureModel("qwen3-14b", ModeCompletion)
  ↓
[ServerManager] Stops embedding server (SIGTERM → wait → SIGKILL)
  ↓
[ServerManager] Launches: llama-server -m /Volumes/1tb/models/qwen3-14b.gguf \
                          --host 127.0.0.1 --port 8080 -c 4096 -ngl 99
  ↓
[ServerManager] Polls health endpoint until ready
  ↓
[Orchestrator] Scores candidates via completion API
```

### Phase 3: Hook Generation
```
[Orchestrator] GenerateHooks starts
  ↓
[ServerManager] EnsureModel("llama-3-8b-instruct", ModeCompletion)
  ↓
[ServerManager] Stops scoring server
  ↓
[ServerManager] Launches: llama-server -m /Volumes/1tb/models/llama-3-8b-instruct.gguf \
                          --host 127.0.0.1 --port 8080 -c 4096 -ngl 99
  ↓
[ServerManager] Polls health endpoint until ready
  ↓
[Orchestrator] Generates hooks via completion API
```

### Cleanup
```
[CLI] Pipeline completes (or Ctrl+C pressed)
  ↓
[CLI] defer orch.Close() executes
  ↓
[Orchestrator] Close() → ServerManager.Stop()
  ↓
[ServerManager] Sends SIGTERM to process
  ↓
[ServerManager] Waits up to 5 seconds
  ↓
[ServerManager] SIGKILL if still running
  ↓
[ServerManager] Process terminated, port 8080 released
```

## Model Resolution Examples

The ServerManager can find models in multiple ways:

### 1. Full Absolute Path
```yaml
embedding_model: "/Volumes/1tb/models/nomic-embed-text-v1.5.gguf"
```

### 2. Relative to models_dir
```yaml
models_dir: "/Volumes/1tb/models"
embedding_model: "nomic-embed-text-v1.5.gguf"
```

### 3. Model Name Search
```yaml
models_dir: "/Volumes/1tb/models"
embedding_model: "nomic-embed"  # Recursively finds matching .gguf file
```

### 4. Subdirectory
```yaml
models_dir: "/Volumes/1tb/models"
embedding_model: "embedding/nomic-embed-text-v1.5.gguf"
```

## CLI Flag Reference

### --managed Flag
Available in all three main commands:

- `generate-related --managed`: Enable managed mode for related article generation
- `generate-hooks --managed`: Enable managed mode for hook generation  
- `pipeline --managed`: Enable managed mode for full pipeline

Default: `false` (must be explicitly enabled)

Priority: CLI flag > config file > default

### Other Relevant Flags

- `--config <path>`: Path to configuration file (default: config.yaml)
- `--dry-run`: Simulate operations without starting actual server
- `--scoring-model <name>`: Override scoring model from config
- `--hook-model <name>`: Override hook model from config

## Debugging Tips

### Check if Server is Running
```bash
# During pipeline execution, in another terminal:
curl http://localhost:8080/health

# Expected: {"status":"ok"} (or similar)
```

### View Server Output
Server stderr is captured by ServerManager. On errors, the last 500 characters are included in error messages.

### Test Model Resolution
```bash
# Create a small test program:
cat > test_resolve.go << 'EOF'
package main

import (
    "fmt"
    "github.com/rippreport/go-content-pipeline/internal/server"
)

func main() {
    sm := server.NewServerManager(server.ManagerConfig{
        ModelsDir: "/Volumes/1tb/models",
    })
    
    path, err := sm.ResolveModelPath("nomic-embed")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Resolved: %s\n", path)
    }
}
EOF

go run test_resolve.go
# Output: Resolved: /Volumes/1tb/models/nomic-embed-text-v1.5.gguf
```

### Check Configuration
```bash
# Verify managed mode is enabled
grep -A 10 "llama_server:" config.yaml | grep "managed:"

# Expected: managed: true
```

## Troubleshooting

### Error: "timeout waiting for llama-server to load model"
**Cause**: Server took longer than 90 seconds to start

**Solutions**:
1. Check if model file exists: `ls -lh /Volumes/1tb/models/*.gguf`
2. Verify binary path: `ls -l /Users/ryandunphy/.unsloth/llama.cpp/build/bin/llama-server`
3. Check available memory: Large models need sufficient RAM/VRAM
4. Increase timeout in config (future enhancement)

### Error: "model not found"
**Cause**: Model file could not be located

**Solutions**:
1. Use full path: `/Volumes/1tb/models/model-name.gguf`
2. Check models_dir setting: `ls -R /Volumes/1tb/models`
3. Verify file has `.gguf` extension

### Error: "address already in use"
**Cause**: Port 8080 is occupied by another process

**Solutions**:
1. Check for existing llama-server: `ps aux | grep llama-server`
2. Kill stray processes: `pkill -9 llama-server`
3. Wait a few seconds for port to be released

### Server Keeps Restarting
**Cause**: Health check failing immediately after start

**Solutions**:
1. Check server logs (last 500 chars shown in error)
2. Verify model is compatible with llama-server version
3. Ensure sufficient disk space

## Performance Characteristics

### Model Loading Times (Approximate)
- Small models (< 5GB): 2-5 seconds
- Medium models (5-15GB): 5-15 seconds  
- Large models (> 15GB): 15-30 seconds

### Memory Usage
- One model loaded at a time
- VRAM/RAM freed when switching models
- GPU layers (-ngl 99) offload to Metal on Apple Silicon

### Startup Overhead
- Initial server launch: ~5-10 seconds
- Model switches: ~10-20 seconds total (stop + start + health check)
- Total overhead for 3 phases: ~30-60 seconds

Compare to manual management: Several minutes of human time + risk of errors

## Conclusion

The automatic llama-server management system provides:
- ✅ Zero manual intervention required
- ✅ Optimal resource usage (one model at a time)
- ✅ Graceful error handling and recovery
- ✅ Simple configuration and CLI flags
- ✅ Backward compatible with unmanaged mode

**Ready for production use!**
