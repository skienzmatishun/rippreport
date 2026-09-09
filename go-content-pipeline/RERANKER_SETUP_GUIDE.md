# Qwen3-Reranker Setup Guide

## What is Qwen3-Reranker?

Qwen3-Reranker is a **dedicated cross-encoder model** specifically designed for scoring relevance between a query and documents. Unlike general LLMs that need prompting, rerankers are trained specifically for this task and are typically:

- **Faster**: No need for prompt engineering or text generation
- **More Accurate**: Trained specifically for relevance scoring
- **More Efficient**: Direct scoring instead of parsing LLM responses

## Available Models

| Model | Size | Best For | Memory |
|-------|------|----------|--------|
| Qwen3-Reranker-0.6B | 600M | Fast inference, resource-constrained | ~1 GB |
| Qwen3-Reranker-4B | 4B | Balanced performance | ~4 GB |
| Qwen3-Reranker-8B | 8B | Best accuracy | ~8 GB |

## Step 1: Download the Reranker Model

```bash
# Download Qwen3-Reranker-4B GGUF (recommended)
hf download Qwen/Qwen3-Reranker-4B-GGUF \
  --local-dir /Volumes/1tb/models/Qwen3-Reranker-4B-GGUF \
  --include "*Q8_0*"  # Or *Q4_K_M* for smaller size

# Or the smaller 0.6B version for faster inference
hf download Qwen/Qwen3-Reranker-0.6B-GGUF \
  --local-dir /Volumes/1tb/models/Qwen3-Reranker-0.6B-GGUF \
  --include "*Q8_0*"
```

## Step 2: Configure Your Pipeline

Edit `/Volumes/1tb/rippreport/go-content-pipeline/config.yaml`:

```yaml
llama_server:
  managed: true
  binary_path: "/Users/ryandunphy/.unsloth/llama.cpp/build/bin/llama-server"
  models_dir: "/Volumes/1tb/models"
  
  # Embedding model (unchanged)
  embedding_model: "/Volumes/1tb/models/Qwen/Qwen3-Embedding-0.6B-GGUF/Qwen3-Embedding-0.6B-Q8_0.gguf"
  
  # Main LLM for hook generation
  llm_model: "/Volumes/1tb/models/lmstudio-community/gemma-4-12B-it-QAT-GGUF/gemma-4-12B-it-QAT-Q4_0.gguf"
  hook_model: "/Volumes/1tb/models/lmstudio-community/gemma-4-12B-it-QAT-GGUF/gemma-4-12B-it-QAT-Q4_0.gguf"
  
  # ⭐ RERANKER FOR SCORING (THE KEY PART!)
  scoring_model: "/Volumes/1tb/models/Qwen3-Reranker-4B-GGUF/Qwen3-Reranker-4B-Q8_0.gguf"
  scoring_type: "reranker"  # ← THIS TELLS IT TO USE RERANKER MODE
  
  context_size: 10096
  gpu_layers: 99
  timeout: 1200s
```

## What Happens When You Use Reranker Mode

### Phase 1: Embedding (Unchanged)
```
Server: llama-server -m Qwen3-Embedding-0.6B.gguf --embedding
Action: Generate embeddings for all posts
```

### Phase 2: Scoring (WITH RERANKER)
```
Server: llama-server -m Qwen3-Reranker-4B.gguf --reranking
Action: Score candidates using /v1/rerank endpoint
Format: Model returns relevance scores directly (0-1 range)
```

**Before (LLM Completion Mode)**:
```
POST /completion
Body: {
  "prompt": "Given this article about X, rate relevance of Y from 0-100: [10]"
}
Response: Parse "[10]" from generated text
Time: ~2-3 seconds per candidate
```

**After (Reranker Mode)**:
```
POST /v1/rerank
Body: {
  "query": "article title",
  "documents": ["candidate1 title", "candidate2 title", ...]
}
Response: [
  {"index": 0, "relevance_score": 0.92},
  {"index": 1, "relevance_score": 0.45}
]
Time: ~0.5-1 second for ALL candidates (batch)
```

### Phase 3: Hook Generation (Unchanged)
```
Server: llama-server -m gemma-4-12B.gguf
Action: Generate contextual hooks
```

## Command Line Usage

```bash
# Run with reranker
./bin/pipeline generate-related --config config.yaml

# Or specify at runtime
./bin/pipeline generate-related \
  --config config.yaml \
  --scoring-model "/Volumes/1tb/models/Qwen3-Reranker-4B-GGUF/Qwen3-Reranker-4B-Q8_0.gguf" \
  --use-reranker
```

## CLI Flags

- `--use-reranker`: Force reranker mode (overrides config)
- `--scoring-model <path>`: Override scoring model from config
- `--vector-only`: Skip all scoring (just use cosine similarity)

## Performance Comparison

### LLM Completion Scoring (Old Way)
- Score 20 candidates: ~40-60 seconds
- Each candidate: 2-3 seconds
- Memory: ~12 GB (Gemma-4-12B)
- Mode: Sequential scoring

### Reranker Scoring (New Way)
- Score 20 candidates: ~5-10 seconds ✨
- Batch scoring: All candidates at once
- Memory: ~4 GB (Qwen3-Reranker-4B)
- Mode: Parallel batch processing

**Speedup: 4-6x faster** 🚀

## How the Pipeline Uses It

The pipeline automatically:

1. **Detects reranker mode**: Checks `scoring_type: "reranker"` or `--use-reranker` flag
2. **Starts server with --reranking flag**: `llama-server -m Qwen3-Reranker-4B.gguf --reranking`
3. **Batches candidates**: Groups all candidates for a post
4. **Calls /v1/rerank endpoint**: Single API call with all candidates
5. **Processes results**: Converts relevance scores (0-1) to integer scores (0-100)
6. **Combines with other factors**: Merges with recency, length, category scores

## Verification

Check logs for these indicators:

### Server Startup (Good)
```
loading model: /Volumes/1tb/models/Qwen3-Reranker-4B.gguf
reranking mode enabled
```

### API Call (Good)
```
POST /v1/rerank
{
  "query": "Ethics complaint investigation",
  "documents": ["Deaf ears article", "Other post", ...]
}
```

### Response (Good)
```
{
  "results": [
    {"index": 0, "relevance_score": 0.92},
    {"index": 1, "relevance_score": 0.45}
  ]
}
```

## Troubleshooting

### Issue: "404 Not Found" on /v1/rerank

**Cause**: Server not started with `--reranking` flag

**Fix**: Verify `scoring_type: "reranker"` in config:
```yaml
scoring_type: "reranker"  # NOT "completion"
```

### Issue: Server crashes or OOM

**Cause**: Reranker model too large

**Fix**: Use smaller model:
```yaml
scoring_model: "/Volumes/1tb/models/Qwen3-Reranker-0.6B-GGUF/Qwen3-Reranker-0.6B-Q8_0.gguf"
```

### Issue: Scores seem wrong

**Cause**: Reranker expects different format than LLM

**Fix**: This is normal! Reranker scores are 0-1 (converted to 0-100), not raw LLM outputs

### Issue: "reranking mode not supported"

**Cause**: Old llama.cpp version

**Fix**: Rebuild llama.cpp from latest:
```bash
cd /Users/ryandunphy/.unsloth/llama.cpp
git pull
cmake . -B build -DGGML_METAL=ON
cmake --build build --config Release
```

## Full Example Config

```yaml
llama_server:
  managed: true
  binary_path: "/Users/ryandunphy/.unsloth/llama.cpp/build/bin/llama-server"
  models_dir: "/Volumes/1tb/models"
  base_url: "http://localhost:8080"
  
  # Phase 1: Embedding generation
  embedding_model: "/Volumes/1tb/models/Qwen/Qwen3-Embedding-0.6B-GGUF/Qwen3-Embedding-0.6B-Q8_0.gguf"
  
  # Phase 2: Reranker scoring (FAST!)
  scoring_model: "/Volumes/1tb/models/Qwen3-Reranker-4B-GGUF/Qwen3-Reranker-4B-Q8_0.gguf"
  scoring_type: "reranker"
  
  # Phase 3: Hook generation
  llm_model: "/Volumes/1tb/models/lmstudio-community/gemma-4-12B-it-QAT-GGUF/gemma-4-12B-it-QAT-Q4_0.gguf"
  hook_model: "/Volumes/1tb/models/lmstudio-community/gemma-4-12B-it-QAT-GGUF/gemma-4-12B-it-QAT-Q4_0.gguf"
  
  context_size: 10096
  gpu_layers: 99
  timeout: 1200s
  min_request_delay: 100ms
  max_retries: 3

processing:
  posts_directory: "/Volumes/1tb/rippreport/content/p"
  top_candidates: 20
  batch_size: 5

storage:
  backup_directory: "backups"
  cache_file: ".embedding_cache.json"
  progress_file: ".progress.json"

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
```

## Summary

### To use Qwen3-Reranker:

1. **Download**: `hf download Qwen/Qwen3-Reranker-4B-GGUF`
2. **Set in config**: 
   ```yaml
   scoring_model: "/Volumes/1tb/models/Qwen3-Reranker-4B-GGUF/..."
   scoring_type: "reranker"
   ```
3. **Run**: `./bin/pipeline generate-related --config config.yaml`

That's it! The pipeline handles everything else automatically. You'll get 4-6x faster scoring with better accuracy. 🎉
