# Gemini Hook Refinement

The Go pipeline now supports refining hooks using the Gemini API, matching the functionality previously only available in the Python version.

## Overview

The `refine-hooks` command reads existing `sidebar-hooks.yaml` files and uses Google's Gemini API to polish the hook text for publication. This provides editorial refinement that transforms raw "bridge briefs" into aggressive, publication-ready contextual bridges.

## Setup

### 1. Get a Gemini API Key

1. Go to [Google AI Studio](https://aistudio.google.com/app/apikey)
2. Create a new API key
3. Export it as an environment variable:

```bash
export GEMINI_API_KEY="your-api-key-here"
```

### 2. Generate Hooks First

Before refining, you need to generate the initial hooks:

```bash
./bin/pipeline generate-hooks --posts october-surprise
```

This creates `content/p/october-surprise/sidebar-hooks.yaml` with briefs.

### 3. Refine the Hooks

```bash
./bin/pipeline refine-hooks --posts october-surprise
```

This sends the briefs to Gemini and updates the `text` field with refined hooks.

## Usage Examples

### Basic Refinement

```bash
# Refine a single post
./bin/pipeline refine-hooks --posts october-surprise

# Refine multiple posts
./bin/pipeline refine-hooks --posts october-surprise,alabama-ag

# Refine all posts with hooks
./bin/pipeline refine-hooks
```

### Advanced Options

```bash
# Use a specific Gemini model
./bin/pipeline refine-hooks --gemini-model gemini-3.6-flash

# Force refinement even if already refined
./bin/pipeline refine-hooks --force --posts october-surprise

# Dry run to preview what would be refined
./bin/pipeline refine-hooks --dry-run
```

## How It Works

1. **Loads Existing Hooks**: Reads `sidebar-hooks.yaml` files from post directories
2. **Checks Refinement Status**: Skips posts that are already refined (unless `--force`)
3. **Prepares Briefs**: Extracts the `brief` field from each hook
4. **Calls Gemini API**: Sends briefs to Gemini with editorial instructions
5. **Updates Text**: Replaces the `text` field with refined output
6. **Updates Metadata**: Adds `refined_by: gemini` and `refined_at` timestamp

## Output Format

Before refinement:
```yaml
hooks:
  related-post-slug:
    brief: |
      THE CONFLICT:
      Main article discusses corruption in Baldwin County.
      
      THE PIVOT:
      This widget article exposes Sheriff Lowery's abuse of power.
      
      THE DISCOVERY:
      Harry Still was unlawfully arrested for criticizing the sheriff.
    text: ""
    widget_date: "2025-01-15"
```

After refinement:
```yaml
metadata:
  refined_by: gemini
  refined_at: "2026-11-01T12:34:56.789012Z"
  
hooks:
  related-post-slug:
    brief: |
      THE CONFLICT:
      Main article discusses corruption in Baldwin County.
      
      THE PIVOT:
      This widget article exposes Sheriff Lowery's abuse of power.
      
      THE DISCOVERY:
      Harry Still was unlawfully arrested for criticizing the sheriff.
    text: "Sheriff Lowery weaponizes his badge to silence critics who dare question Baldwin County's political machine. Harry Still's unlawful arrest proves dissent comes with a price."
    widget_date: "2025-01-15"
```

## Comparison with Python Version

The Go implementation is functionally equivalent to:

```bash
python3 generate_all_hooks.py --refine --gemini-model gemini-3.6-flash --post october-surprise
```

### Key Differences

- **Go**: `./bin/pipeline refine-hooks --posts october-surprise --gemini-model gemini-3.6-flash`
- **Python**: `python3 generate_all_hooks.py --refine --gemini-model gemini-3.6-flash --post october-surprise`

The Go version:
- Is faster (native concurrency)
- Has simpler setup (no Python dependencies)
- Uses the same Gemini API and prompts
- Produces identical output format

## Editorial Guidelines

The Gemini refinement follows these editorial principles:

### Style
- **Cynical, hard-hitting, investigative tone**
- Use specific names (Wilters, Lowery, etc.)
- NO passive voice
- 6th-grade vocabulary
- Variety in sentence structure

### Structure
- **Sentence 1**: The setup. Connect to the main article's theme.
- **Sentence 2**: The payoff. Use the PIVOT to show why this is the "smoking gun."

### Constraints
- 220-320 characters total
- 2 sentences exactly
- Focus on the PIVOT section of the brief

## Cost Considerations

Gemini API pricing (as of Nov 2026):
- **gemini-3.7-flash**: Very low cost, suitable for most refinement tasks
- **gemini-3.6-flash**: Older model, still effective

Each hook refinement costs a fraction of a cent. Processing 100 posts with 5 hooks each typically costs less than $0.50.

## Troubleshooting

### "GEMINI_API_KEY is required"

Make sure you've exported the API key:
```bash
export GEMINI_API_KEY="your-key-here"
echo $GEMINI_API_KEY  # Verify it's set
```

### "No hooks file found"

Generate hooks first:
```bash
./bin/pipeline generate-hooks --posts your-post-slug
```

### "Already refined"

Use `--force` to re-refine:
```bash
./bin/pipeline refine-hooks --force --posts your-post-slug
```

### Rate Limits

Gemini has generous rate limits, but if you hit them:
- Wait a few seconds between batches
- Process posts in smaller batches
- Use `--posts` to target specific posts

## Integration with Pipeline

You can run the full pipeline with automatic refinement:

```bash
# Generate related articles, hooks, and refine in one command
./bin/pipeline generate-related --posts october-surprise
./bin/pipeline generate-hooks --posts october-surprise
./bin/pipeline refine-hooks --posts october-surprise
```

Or use the Python pipeline mode for now:
```bash
python3 generate_all_hooks.py --pipeline --posts october-surprise
```

(Future enhancement: Add `--refine` flag to `pipeline` command for one-step processing)
