# Speculative Decoding & Qwen 3.8 Support Guide

## Overview

The Go Content Pipeline now supports advanced llama.cpp features:
- **Speculative Decoding**: 2-4x faster inference using draft models
- **Qwen 3.8 Reasoning**: Optimized settings for Qwen 3.8's thinking mode
- **Preserve Thinking**: Maintain reasoning traces across conversation turns

## Speculative Decoding

### What is Speculative Decoding?

Speculative decoding uses a small, fast draft model to predict several tokens ahead. The main model then verifies all predictions in a single batch. When predictions are correct (which they often are), generation becomes 2-4x faster while maintaining mathematically identical output.

### Supported Methods

#### 1. DFlash (Recommended for Qwen Models)
Block-diffusion drafting that generates an entire block of tokens per step.

```yaml
llama_server:
  managed: true
  llm_model: "Qwen3-4B"
  draft_model: "Qwen3-4B-DFlash"  # z-lab/Qwen3-4B-DFlash
  spec_type: "draft-dflash"
  spec_draft_n_max: 15  # DFlash can draft 7-15 tokens
```

**Models Available**:
- `Qwen3-4B-DFlash` (for Qwen3-4B)
- `Qwen3-8B-DFlash` (for Qwen3-8B)

**Performance**: ~2-3x speedup with 15-token drafts

#### 2. EAGLE-3
Single-layer transformer that reads the target model's hidden states.

```yaml
llama_server:
  managed: true
  llm_model: "Qwen3-4B"
  draft_model: "Qwen3-4B-eagle3"  # AngelSlim/Qwen3-4B_eagle3
  spec_type: "draft-eagle3"
  spec_draft_n_max: 4
```

**Models Available**:
- `Qwen3-1.7B_eagle3`, `Qwen3-4B_eagle3`, `Qwen3-8B_eagle3`
- `Qwen3-14B_eagle3`, `Qwen3-32B_eagle3`
- LLaMA, Gemma variants also available

**Performance**: Higher acceptance rate than standalone draft models

#### 3. DSpark
DFlash backbone + semi-autoregressive Markov head.

```yaml
llama_server:
  managed: true
  llm_model: "Qwen3-4B"
  draft_model: "dspark_qwen3_4b_block7"  # deepseek-ai/dspark_qwen3_4b_block7
  spec_type: "draft-dspark"
  spec_draft_n_max: 7
```

**Performance**: Balanced speed and accuracy

#### 4. n-gram (No Draft Model Required)
Pattern matching without additional models.

```yaml
llama_server:
  managed: true
  spec_type: "ngram-mod"
  spec_draft_n_max: 64
```

**Best For**: Repetitive tasks (code refactoring, structured output)

### How Draft Models Work

When you configure a draft model:

1. **Startup**: The managed server loads BOTH the main model AND the draft model
   ```bash
   llama-server -m Qwen3-4B.gguf -md Qwen3-4B-DFlash.gguf \
     --spec-type draft-dflash --spec-draft-n-max 15
   ```

2. **Generation**: Draft model predicts 15 tokens → main model verifies in one batch

3. **Model Switching**: When switching phases (embedding → scoring → hooks), both main and draft models are unloaded and reloaded

### Configuration Options

```yaml
llama_server:
  # Main model
  llm_model: "Qwen3-4B"
  
  # Draft model (path/name relative to models_dir or absolute)
  draft_model: "Qwen3-4B-DFlash"
  
  # Speculative decoding type
  spec_type: "draft-dflash"  # or "draft-eagle3", "draft-dspark", "ngram-mod"
  
  # Max tokens to draft per step
  spec_draft_n_max: 15  # DFlash: 7-15, EAGLE-3: 3-5, DSpark: 7, ngram: up to 64
```

### Model Resolution

Draft models are resolved the same way as main models:

```yaml
# Full path
draft_model: "/Volumes/1tb/models/Qwen3-4B-DFlash.gguf"

# Relative to models_dir
draft_model: "Qwen3-4B-DFlash.gguf"

# Pattern search (finds matching .gguf recursively)
draft_model: "Qwen3-4B-DFlash"

# Subdirectory
draft_model: "draft/Qwen3-4B-DFlash.gguf"
```

## Qwen 3.8 Reasoning Support

### What is Qwen 3.8?

Qwen 3.8 is a hybrid thinking model with two modes:
- **Thinking Mode**: Extended reasoning with internal thought process (like o1)
- **Instruct Mode**: Direct responses without visible reasoning

### Reasoning Effort

Control how deeply Qwen 3.8 thinks before responding:

```yaml
llama_server:
  llm_model: "Qwen3.8-27B"
  reasoning_effort: "medium"  # xhigh, high, medium, low, none
```

**Options**:
- `xhigh` (default): Complex tasks demanding thorough analysis
- `high`: Standard reasoning depth
- `medium`: Balance between accuracy and speed ⭐ **Recommended**
- `low`: Quick reasoning for simple tasks
- `none`: Disable thinking mode completely

**Command-line equivalent**:
```bash
llama-server -m Qwen3.8-27B.gguf \
  --chat-template-kwargs '{"reasoning_effort":"medium"}'
```

### Preserve Thinking

Keep reasoning traces from previous conversation turns:

```yaml
llama_server:
  llm_model: "Qwen3.8-27B"
  reasoning_effort: "medium"
  preserve_thinking: true
```

**Trade-offs**:
- ✅ May improve accuracy in multi-turn conversations
- ❌ Increases token usage (context includes all thinking traces)
- ❌ Slower as context grows

**Command-line equivalent**:
```bash
llama-server -m Qwen3.8-27B.gguf \
  --chat-template-kwargs '{"reasoning_effort":"medium","preserve_thinking":true}'
```

### Qwen 3.8 Sampling Parameters

Qwen 3.8 uses different settings for thinking vs non-thinking modes:

**Thinking Mode** (default for Qwen 3.8):
```
temperature: 1.0
top_p: 0.95
top_k: 20
min_p: 0.0
presence_penalty: 0.0
repetition_penalty: 1.0
```

**Instruct Mode** (non-thinking):
```
temperature: 0.7
top_p: 0.80
top_k: 20
min_p: 0.0
presence_penalty: 1.5
repetition_penalty: 1.0
```

These are set automatically by llama.cpp's chat template when you use Qwen 3.8 models.

## Complete Configuration Examples

### Example 1: Qwen 3.8 with DFlash Speculative Decoding

```yaml
llama_server:
  managed: true
  binary_path: "/Users/ryandunphy/.unsloth/llama.cpp/build/bin/llama-server"
  models_dir: "/Volumes/1tb/models"
  base_url: "http://localhost:8080"
  
  # Main models
  embedding_model: "nomic-embed-text-v1.5"
  scoring_model: "Qwen3.8-27B"
  hook_model: "Qwen3.8-27B"
  
  # Speculative decoding for 2-3x speedup
  draft_model: "Qwen3.8-27B-DFlash"
  spec_type: "draft-dflash"
  spec_draft_n_max: 15
  
  # Qwen 3.8 reasoning configuration
  reasoning_effort: "medium"
  preserve_thinking: false  # Set true for multi-turn accuracy
  
  context_size: 32768  # Qwen 3.8 supports up to 262K
  gpu_layers: 99
  timeout: 180s  # Longer timeout for thinking mode
```

### Example 2: Qwen 3 with EAGLE-3 Speculative Decoding

```yaml
llama_server:
  managed: true
  models_dir: "/Volumes/1tb/models"
  
  llm_model: "Qwen3-4B"
  scoring_model: "Qwen3-4B"
  hook_model: "Qwen3-4B"
  
  # EAGLE-3 for high acceptance rate
  draft_model: "Qwen3-4B-eagle3"
  spec_type: "draft-eagle3"
  spec_draft_n_max: 4
  
  context_size: 4096
  gpu_layers: 99
```

### Example 3: Standard Model (No Speculative Decoding)

```yaml
llama_server:
  managed: true
  models_dir: "/Volumes/1tb/models"
  
  embedding_model: "nomic-embed-text-v1.5"
  llm_model: "qwen3-14b"
  
  # No draft model - standard generation
  # draft_model: ""
  # spec_type: ""
  
  context_size: 4096
  gpu_layers: 99
```

### Example 4: n-gram Speculative Decoding (No Draft Model)

```yaml
llama_server:
  managed: true
  models_dir: "/Volumes/1tb/models"
  
  llm_model: "qwen3-14b"
  
  # Pattern-based drafting without separate model
  spec_type: "ngram-mod"
  spec_draft_n_max: 64
  
  context_size: 4096
  gpu_layers: 99
```

## Download Draft Models

### Using Hugging Face Hub

```bash
pip install -U "huggingface_hub[cli]"

# DFlash models
hf download z-lab/Qwen3-4B-DFlash --local-dir /Volumes/1tb/models/Qwen3-4B-DFlash
hf download z-lab/Qwen3-8B-DFlash --local-dir /Volumes/1tb/models/Qwen3-8B-DFlash

# EAGLE-3 models
hf download AngelSlim/Qwen3-4B_eagle3 --local-dir /Volumes/1tb/models/Qwen3-4B-eagle3
hf download AngelSlim/Qwen3-8B_eagle3 --local-dir /Volumes/1tb/models/Qwen3-8B-eagle3
hf download AngelSlim/Qwen3-14B_eagle3 --local-dir /Volumes/1tb/models/Qwen3-14B-eagle3

# DSpark models
hf download deepseek-ai/dspark_qwen3_4b_block7 \
  --local-dir /Volumes/1tb/models/dspark_qwen3_4b_block7
```

### Using Unsloth (Qwen 3.8 GGUFs)

```bash
# Qwen 3.8-27B main model
hf download unsloth/Qwen3.8-27B-GGUF \
  --local-dir /Volumes/1tb/models/Qwen3.8-27B-GGUF \
  --include "*UD-Q4_K_XL*"

# Qwen 3.8 DFlash draft model (when available)
hf download unsloth/Qwen3.8-27B-DFlash-GGUF \
  --local-dir /Volumes/1tb/models/Qwen3.8-27B-DFlash
```

## Model Compatibility

### Target-Draft Pairs

Draft models must be trained for specific target models:

| Target Model | Compatible Draft Models |
|-------------|------------------------|
| Qwen3-4B | Qwen3-4B-DFlash, Qwen3-4B-eagle3, dspark_qwen3_4b_block7 |
| Qwen3-8B | Qwen3-8B-DFlash, Qwen3-8B-eagle3 |
| Qwen3-14B | Qwen3-14B-eagle3 |
| Qwen3.8-27B | Qwen3.8-27B-DFlash (when available) |
| LLaMA 3.1-8B | yuhuili/EAGLE3-LLaMA3.1-Instruct-8B |
| Gemma 4-31B | RedHatAI/gemma-4-31B-it-speculator.eagle3 |

**Important**: Using a mismatched draft model will either fail to load or produce incorrect output.

## Performance Expectations

### Speedup Factors

Actual speedup depends on:
- Draft model acceptance rate (how often predictions are correct)
- Hardware (memory bandwidth, GPU compute)
- Context length (longer contexts → more verification overhead)

**Typical Results**:

| Method | Draft Tokens | Acceptance Rate | Speedup |
|--------|--------------|----------------|---------|
| DFlash | 15 | ~70% | 2.0-3.0x |
| EAGLE-3 | 4 | ~80% | 1.8-2.5x |
| DSpark | 7 | ~75% | 2.0-2.8x |
| ngram-mod | 64 | varies | 1.5-4.0x (task-dependent) |

### Memory Requirements

Speculative decoding requires loading two models:

```
Total VRAM = Main Model + Draft Model + KV Cache
```

**Example** (Qwen3-4B + DFlash):
- Main model (Q4_K_M): ~2.5 GB
- Draft model (Q4_K_M): ~1.0 GB
- KV cache (4K context): ~0.5 GB
- **Total**: ~4 GB (fits on 8GB VRAM)

**Larger Models** (Qwen3.8-27B + DFlash):
- Main model (Q4_K_M): ~16 GB
- Draft model: ~4 GB
- KV cache (32K context): ~2 GB
- **Total**: ~22 GB (needs 24GB+ VRAM or Metal unified memory)

## Troubleshooting

### Draft Model Not Loading

**Error**: `failed to resolve draft model: model "Qwen3-4B-DFlash" not found`

**Solutions**:
1. Check model file exists: `ls -lh /Volumes/1tb/models/*DFlash*`
2. Use full path in config: `draft_model: "/Volumes/1tb/models/Qwen3-4B-DFlash.gguf"`
3. Verify models_dir setting: `models_dir: "/Volumes/1tb/models"`

### Wrong Draft Model for Target

**Error**: Server starts but output is gibberish or generation fails

**Solution**: Ensure draft model matches target model:
```yaml
# CORRECT
llm_model: "Qwen3-4B"
draft_model: "Qwen3-4B-DFlash"  # Trained for Qwen3-4B

# WRONG
llm_model: "Qwen3-8B"
draft_model: "Qwen3-4B-DFlash"  # Trained for different model!
```

### Out of Memory

**Error**: `failed to allocate tensor` or server crashes

**Solutions**:
1. Reduce draft length: `spec_draft_n_max: 7` (instead of 15)
2. Use smaller quant: Download Q3_K_M instead of Q4_K_M
3. Reduce context: `context_size: 2048` (instead of 4096)
4. Use ngram-mod: No draft model needed
   ```yaml
   spec_type: "ngram-mod"
   # draft_model: ""  # Comment out
   ```

### Reasoning Effort Not Working

**Error**: Qwen 3.8 not using thinking mode

**Check**:
1. Verify model is actually Qwen 3.8: `ls -lh /Volumes/1tb/models/*Qwen3.8*`
2. Check server args in logs - should see: `--chat-template-kwargs '{"reasoning_effort":"medium"}'`
3. Try explicit reasoning_effort: `reasoning_effort: "xhigh"`

### Slower with Speculative Decoding

**Possible Causes**:
1. Draft acceptance rate too low (<50%)
2. CPU bottleneck (draft model not on GPU)
3. Disk offloading (not enough RAM/VRAM)

**Solutions**:
1. Try different draft model (EAGLE-3 vs DFlash)
2. Reduce draft length: `spec_draft_n_max: 7`
3. Check GPU layers: `gpu_layers: 99` (offload everything)
4. Monitor acceptance rate in logs

## Testing Your Configuration

### 1. Verify Configuration Loads

```bash
cd /Volumes/1tb/rippreport/go-content-pipeline

# Check config parses correctly
./bin/pipeline pipeline --config config.yaml --dry-run --help
```

### 2. Test Server Startup

Watch for these log messages (from llama-server stderr):

```
loading model from: /Volumes/1tb/models/Qwen3-4B.gguf
loading draft model from: /Volumes/1tb/models/Qwen3-4B-DFlash.gguf
spec: draft_model = Qwen3-4B-DFlash.gguf
spec: n_draft = 15
```

### 3. Monitor Performance

Check llama-server output for acceptance statistics:

```
draft acceptance rate = 0.72 (720 accepted / 1000 generated)
```

Higher acceptance rate = better speedup.

### 4. Benchmark

Run a small test:

```bash
time ./bin/pipeline generate-related \
  --config config.yaml \
  --posts "test-article" \
  --dry-run
```

Compare with and without speculative decoding.

## Reference: llama-server Command Generation

When you configure speculative decoding, the managed server generates commands like:

### DFlash Example
```bash
llama-server \
  -m /Volumes/1tb/models/Qwen3-4B.gguf \
  -md /Volumes/1tb/models/Qwen3-4B-DFlash.gguf \
  --host 127.0.0.1 \
  --port 8080 \
  -c 4096 \
  -ngl 99 \
  --spec-type draft-dflash \
  --spec-draft-n-max 15
```

### Qwen 3.8 with Reasoning
```bash
llama-server \
  -m /Volumes/1tb/models/Qwen3.8-27B.gguf \
  --host 127.0.0.1 \
  --port 8080 \
  -c 32768 \
  -ngl 99 \
  --chat-template-kwargs '{"reasoning_effort":"medium","preserve_thinking":true}'
```

### Combined: DFlash + Qwen 3.8 Reasoning
```bash
llama-server \
  -m /Volumes/1tb/models/Qwen3.8-27B.gguf \
  -md /Volumes/1tb/models/Qwen3.8-27B-DFlash.gguf \
  --host 127.0.0.1 \
  --port 8080 \
  -c 32768 \
  -ngl 99 \
  --spec-type draft-dflash \
  --spec-draft-n-max 15 \
  --chat-template-kwargs '{"reasoning_effort":"medium"}'
```

## Best Practices

### For Development/Testing
- Use `ngram-mod` or no speculative decoding
- Simpler debugging without draft model complexity

### For Production (Speed)
- Use DFlash or EAGLE-3 with properly matched models
- Monitor acceptance rate (aim for >60%)
- Balance draft length with memory constraints

### For Qwen 3.8 Thinking
- Start with `reasoning_effort: "medium"`
- Use `preserve_thinking: false` unless multi-turn accuracy critical
- Increase `timeout` to 180s-300s for complex reasoning

### For Maximum Speed
```yaml
llama_server:
  llm_model: "Qwen3-4B"
  draft_model: "Qwen3-4B-DFlash"
  spec_type: "draft-dflash"
  spec_draft_n_max: 15
  gpu_layers: 99  # All on GPU
  context_size: 4096  # Moderate context
```

### For Maximum Quality (Qwen 3.8)
```yaml
llama_server:
  llm_model: "Qwen3.8-27B"
  reasoning_effort: "xhigh"
  preserve_thinking: true
  context_size: 65536
  timeout: 300s
```

### For Memory-Constrained Systems
```yaml
llama_server:
  llm_model: "Qwen3-4B"
  spec_type: "ngram-mod"  # No draft model
  spec_draft_n_max: 32
  context_size: 2048
```

## Resources

- [llama.cpp Speculative Decoding Docs](https://github.com/ggml-org/llama.cpp/blob/master/docs/speculative.md)
- [Unsloth Qwen 3.8 Guide](https://unsloth.ai/docs/models/qwen3.8)
- [DFlash Paper & Implementation](https://github.com/z-lab-ai/dflash)
- [EAGLE-3 Models on HuggingFace](https://huggingface.co/collections/yuhuili/eagle-3-676e42b52e7e5cc4e9e4deeb)

## Summary

✅ **Speculative decoding** can provide 2-4x faster generation with identical output  
✅ **DFlash** is recommended for Qwen models (15-token drafts)  
✅ **EAGLE-3** offers high acceptance rates (3-5 token drafts)  
✅ **Qwen 3.8** has hybrid thinking mode with configurable depth  
✅ **`reasoning_effort: "medium"`** balances speed and quality  
✅ **`preserve_thinking`** helps multi-turn accuracy at cost of tokens  

Choose configuration based on your priorities: speed, quality, or memory constraints.
