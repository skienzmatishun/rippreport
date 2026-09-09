# Quick Start: Qwen 3.8 with DFlash Speculative Decoding

## Your Setup

This guide is tailored for your environment:
- **System**: macOS with Metal GPU
- **Binary**: `/Users/ryandunphy/.unsloth/llama.cpp/build/bin/llama-server`
- **Models**: `/Volumes/1tb/models/`

## Step 1: Update Your Config

Edit `/Volumes/1tb/rippreport/go-content-pipeline/config.yaml`:

```yaml
llama_server:
  # Managed server
  managed: true
  binary_path: "/Users/ryandunphy/.unsloth/llama.cpp/build/bin/llama-server"
  models_dir: "/Volumes/1tb/models"
  base_url: "http://localhost:8080"
  
  # Models
  embedding_model: "nomic-embed-text-v1.5"
  scoring_model: "Qwen3.8-27B"  # Or "Qwen3-14B", "Qwen3-8B"
  hook_model: "Qwen3.8-27B"
  
  # ⚡ DFlash Speculative Decoding (2-3x faster)
  draft_model: "Qwen3.8-27B-DFlash"  # Or "Qwen3-8B-DFlash", "Qwen3-4B-DFlash"
  spec_type: "draft-dflash"
  spec_draft_n_max: 15  # Draft 15 tokens per step
  
  # 🧠 Qwen 3.8 Reasoning (optional)
  reasoning_effort: "medium"  # Options: xhigh, high, medium, low, none
  preserve_thinking: false    # Set true for multi-turn conversations
  
  # Settings
  context_size: 32768  # Qwen 3.8 supports up to 262K
  gpu_layers: 99       # All layers on Metal
  timeout: 180s        # Longer timeout for thinking mode
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
```

## Step 2: Download Required Models

### Main Models (if you don't have them)

```bash
# Embedding model
hf download unsloth/nomic-embed-text-v1.5-GGUF \
  --local-dir /Volumes/1tb/models/nomic-embed-text-v1.5 \
  --include "*Q8_0*"

# Qwen 3.8-27B (main model)
hf download unsloth/Qwen3.8-27B-GGUF \
  --local-dir /Volumes/1tb/models/Qwen3.8-27B-GGUF \
  --include "*UD-Q4_K_XL*"  # Or *Q4_K_M* for smaller size
```

### DFlash Draft Model (NEW - for speedup)

```bash
# Check if DFlash model exists for Qwen 3.8-27B
# If available:
hf download unsloth/Qwen3.8-27B-DFlash-GGUF \
  --local-dir /Volumes/1tb/models/Qwen3.8-27B-DFlash

# If not available yet, use smaller model:
hf download z-lab/Qwen3-8B-DFlash \
  --local-dir /Volumes/1tb/models/Qwen3-8B-DFlash

# Then update config to use Qwen3-8B for both main and draft:
# scoring_model: "Qwen3-8B"
# draft_model: "Qwen3-8B-DFlash"
```

**Important**: Draft model must match main model version!

## Step 3: Verify Configuration

```bash
cd /Volumes/1tb/rippreport/go-content-pipeline

# Test config loads correctly
./bin/pipeline --help

# Dry run to verify server would start correctly
./bin/pipeline generate-related --config config.yaml --dry-run --posts "test"
```

## Step 4: Run Pipeline

```bash
cd /Volumes/1tb/rippreport/go-content-pipeline

# Generate related articles for all posts
./bin/pipeline generate-related --config config.yaml

# Or for specific posts
./bin/pipeline generate-related --config config.yaml --posts "my-article-slug"

# Full pipeline (related + hooks)
./bin/pipeline pipeline --config config.yaml
```

## What Happens Behind the Scenes

### Without DFlash (Standard)
```
1. Server starts with Qwen3.8-27B
2. Generate 1 token → GPU forward pass → 1 token
3. Generate 1 token → GPU forward pass → 1 token
   ...repeat for all tokens (slow)
```

### With DFlash (2-3x faster)
```
1. Server starts with Qwen3.8-27B + Qwen3.8-27B-DFlash
2. Draft model: predict 15 tokens (fast, small model)
3. Main model: verify all 15 tokens in ONE batch (fast parallel check)
4. Accept correct predictions (typically 10-12 out of 15)
5. Repeat...
   Result: 10-12 tokens per main model forward pass vs 1 token
```

## Expected Performance

### Timing Estimates (Qwen3.8-27B on M1/M2 Mac with 32GB RAM)

**Without DFlash**:
- Embedding phase: ~5 min for 100 posts
- Scoring phase: ~15 min for 100 posts (5 candidates each)
- Hook phase: ~10 min for 100 posts
- **Total**: ~30 min

**With DFlash**:
- Embedding phase: ~5 min (embedding model unchanged)
- Scoring phase: ~6 min (2.5x faster)
- Hook phase: ~4 min (2.5x faster)
- **Total**: ~15 min (50% time savings)

### Memory Usage

```
Qwen3.8-27B Q4_K_M:      ~16 GB
Qwen3.8-27B-DFlash:      ~4 GB
KV Cache (32K context):  ~2 GB
System overhead:         ~2 GB
------------------------
Total:                   ~24 GB (fits in 32GB unified memory)
```

## Monitoring

Watch the terminal output for these indicators:

### Server Startup (Good)
```
loading model from: /Volumes/1tb/models/Qwen3.8-27B.gguf
loading draft model from: /Volumes/1tb/models/Qwen3.8-27B-DFlash.gguf
spec: draft_model = Qwen3.8-27B-DFlash.gguf
spec: n_draft = 15
spec: type = draft-dflash
llama_model_load: total tensors = 765
```

### Generation Statistics (Good)
```
draft acceptance rate = 0.72 (720 accepted / 1000 generated)
```
Higher acceptance rate (>0.60) = good speedup

### Reasoning Mode Active (Good, if using Qwen 3.8)
```
chat_template_kwargs: {"reasoning_effort":"medium"}
```

## Troubleshooting

### "draft model not found"
```bash
# Check if file exists
ls -lh /Volumes/1tb/models/*DFlash*

# If missing, download:
hf download z-lab/Qwen3-8B-DFlash \
  --local-dir /Volumes/1tb/models/Qwen3-8B-DFlash

# Update config with exact path:
draft_model: "/Volumes/1tb/models/Qwen3-8B-DFlash/Qwen3-8B-DFlash-Q4_K_M.gguf"
```

### Out of Memory
```yaml
# Option 1: Reduce draft length
spec_draft_n_max: 7  # Instead of 15

# Option 2: Use smaller quantization
# Download Q3_K_M instead of Q4_K_M

# Option 3: Disable DFlash temporarily
# draft_model: ""  # Comment out
# spec_type: ""    # Comment out
```

### Slower Than Expected
```bash
# Check acceptance rate in logs
# Look for: "draft acceptance rate = X.XX"

# If < 0.50, try:
spec_draft_n_max: 10  # Reduce draft length

# Or switch to EAGLE-3 (if available):
draft_model: "Qwen3.8-27B-eagle3"
spec_type: "draft-eagle3"
spec_draft_n_max: 4
```

### Gibberish Output
```
Wrong draft model for target model!

Fix: Ensure versions match
✅ scoring_model: "Qwen3-8B" + draft_model: "Qwen3-8B-DFlash"
❌ scoring_model: "Qwen3-8B" + draft_model: "Qwen3-4B-DFlash"
```

## Optional: Use n-gram Instead (No Draft Model)

If you want speedup without downloading draft models:

```yaml
llama_server:
  # ... other settings ...
  
  # n-gram pattern matching (no draft model needed)
  spec_type: "ngram-mod"
  spec_draft_n_max: 64
  
  # Comment out draft_model
  # draft_model: ""
```

Works well for repetitive content but lower speedup (~1.5-2x instead of 2-3x).

## Configuration Presets

### Maximum Speed
```yaml
llama_server:
  scoring_model: "Qwen3-8B"  # Smaller = faster
  draft_model: "Qwen3-8B-DFlash"
  spec_type: "draft-dflash"
  spec_draft_n_max: 15
  reasoning_effort: "low"  # Fast reasoning
  context_size: 4096  # Smaller context = faster
```

### Maximum Quality
```yaml
llama_server:
  scoring_model: "Qwen3.8-27B"
  draft_model: "Qwen3.8-27B-DFlash"
  spec_type: "draft-dflash"
  spec_draft_n_max: 15
  reasoning_effort: "xhigh"  # Deep thinking
  preserve_thinking: true
  context_size: 65536  # Long context
  timeout: 300s
```

### Balanced (Recommended)
```yaml
llama_server:
  scoring_model: "Qwen3-14B"  # Or Qwen3.8-27B
  draft_model: "Qwen3-14B-DFlash"  # Match version
  spec_type: "draft-dflash"
  spec_draft_n_max: 15
  reasoning_effort: "medium"  # ⭐ Good balance
  preserve_thinking: false
  context_size: 16384
  timeout: 180s
```

## Verifying It's Working

Run a test and watch for these signs:

```bash
# Terminal 1: Run pipeline
cd /Volumes/1tb/rippreport/go-content-pipeline
./bin/pipeline generate-related --config config.yaml --posts "test-article"

# Look for in output:
✅ "loading draft model from: .../DFlash.gguf"
✅ "spec: n_draft = 15"
✅ "draft acceptance rate = 0.XX" (where XX > 0.50)
✅ Faster completion time than usual
```

## Getting Help

Check these files for more details:
- **`SPECULATIVE_DECODING_GUIDE.md`** - Complete technical guide
- **`FEATURE_SUMMARY_SPECULATIVE_DECODING.md`** - Feature overview
- **`config.example.yaml`** - Annotated configuration

Or review llama.cpp docs:
- https://github.com/ggml-org/llama.cpp/blob/master/docs/speculative.md

## Summary Checklist

- [ ] Update `config.yaml` with draft model settings
- [ ] Download draft model matching your main model
- [ ] Verify config loads: `./bin/pipeline --help`
- [ ] Test with dry run: `--dry-run --posts "test"`
- [ ] Run full pipeline
- [ ] Monitor acceptance rate (should be >0.60)
- [ ] Enjoy 2-3x faster generation! 🚀

---

**Ready to go!** Your pipeline now supports state-of-the-art inference optimization.
