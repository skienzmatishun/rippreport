# Hooks Staging Workflow

Sidebar hooks now support the same staging workflow as related articles, allowing you to generate, review, and apply hooks in separate steps.

## Commands

### Generate Hooks to Staging

Generate hooks and save to staging files instead of writing directly:

```bash
./bin/pipeline generate-hooks --config config.yaml --posts "your-slug" --stage
```

This creates `{slug}/sidebar-hooks.staged.yaml` files that you can review before applying.

### Generate from Staged Rankings

If you've staged rankings, you can generate hooks from those staged rankings:

```bash
./bin/pipeline generate-hooks --config config.yaml --posts "your-slug" --from-staging --stage
```

### Resume Interrupted Hook Generation

Rerun your command with `--resume` to skip posts with valid, complete staged hooks:

```bash
./bin/pipeline generate-hooks --config config.yaml --from-staging --stage --resume
```

Missing, invalid, or incomplete hook files are regenerated. Progress is saved after each post; an interrupted post is retried. Without `--stage`, resume checks production hooks instead. Omit `--resume` to regenerate all selected posts. The full `pipeline --stage --resume` command also resumes hooks.

### Apply Staged Hooks

After reviewing, apply the staged hooks to production:

```bash
# Apply specific posts
./bin/pipeline apply-staged --hooks --posts "your-slug"

# Apply all staged hooks
./bin/pipeline apply-staged --hooks
```

This will:
1. Back up existing `sidebar-hooks.yaml` (if it exists)
2. Write the staged hooks to production
3. Delete the staging file

### Clear Staged Hooks

Remove staged hook files without applying:

```bash
# Clear specific posts
./bin/pipeline clear-staged --hooks --posts "your-slug"

# Clear all staged hooks
./bin/pipeline clear-staged --hooks

# Clear all staged rankings and hooks together
./bin/pipeline clear-staged --all
```

## File Structure

```
content/p/your-post/
├── index.md
├── related-articles.staged.yaml   # Staged rankings
├── sidebar-hooks.staged.yaml      # Staged hooks
├── related-articles.yaml          # (not used by Hugo, for reference)
└── sidebar-hooks.yaml            # Production hooks (used by Hugo)
```

## Workflow Examples

### Full Pipeline with Staging

```bash
# 1. Generate and stage both rankings and hooks
./bin/pipeline generate-related --config config.yaml --posts "fraud" --stage

# 2. Generate hooks from staged rankings
./bin/pipeline generate-hooks --config config.yaml --posts "fraud" --from-staging --stage

# 3. Review the staged files
cat content/p/fraud/related-articles.staged.yaml
cat content/p/fraud/sidebar-hooks.staged.yaml

# 4. Apply both when satisfied
./bin/pipeline apply-staged --posts "fraud"
./bin/pipeline apply-staged --hooks --posts "fraud"
```

### Regenerate Just Hooks

If you want to regenerate hooks for existing rankings:

```bash
# Generate hooks from current front matter (already applied rankings)
./bin/pipeline generate-hooks --config config.yaml --posts "fraud" --stage

# Review
cat content/p/fraud/sidebar-hooks.staged.yaml

# Apply
./bin/pipeline apply-staged --hooks --posts "fraud"
```

## Flags Summary

### generate-hooks
- `--resume`: Skip posts whose selected hook output already covers all current related articles
- `--stage`: Save to staging files instead of writing directly
- `--from-staging`: Read rankings from staging files instead of front matter
- `--dry-run`: Simulate without writing any files
- `--use-compressed`: Use compressed article summaries for context
- `--hook-model`: Override the model used for generation

### apply-staged
- `--hooks`: Apply hooks instead of rankings
- `--posts`: Comma-separated list of post slugs (omit to apply all)

### clear-staged
- `--all`: Clear both staged rankings and hooks; combine with `--posts` to limit posts
- `--hooks`: Clear hooks instead of rankings  
- `--posts`: Comma-separated list of post slugs (omit to clear all)

### Use OpenRouter alongside llama-server

Set `OPENROUTER_API_KEY` in your environment, then add `--openrouter`:

```sh
./bin/pipeline generate-hooks --from-staging --stage --resume --openrouter
```

The default remote model is `google/gemma-4-26b-a4b-it:free`. Override it with
`--openrouter-model MODEL_ID`. These flags also work with the `pipeline` command.
Rebuild the binary after updating the source (`go build -o bin/pipeline ./cmd/pipeline`).

With this option, one local worker and one OpenRouter worker share each post's
article pairs. Each worker generates both the brief and sidebar text for its pair;
the saved hook records that worker's model. This mode uses individual pairs rather
than the usual single-prompt batch. Without `--openrouter`, generation is unchanged.
OpenRouter receives article content (or cached summaries). Requests are spaced at
least three seconds apart. Connection resets, timeouts, and interrupted response
reads retry up to three times with 1-, 2-, and 4-second backoff (also respecting
request spacing). If retries are exhausted, or the helper encounters another
failure such as a rate or daily quota limit, it is disabled for the run and its
work falls back to llama-server. Cancellation interrupts the retry wait.
Speed gains depend on model availability, request quotas, and response times.

The local bridge brief output budget is 30,000 tokens (including batched briefs).
OpenRouter uses a separate 100,000-token output budget for its requests, including
briefs and sidebar text. Override it with `--openrouter-max-tokens N` on either
`generate-hooks` or `pipeline`. This is an output budget, not a context-window size;
the selected provider/model must support the requested value.

### Add Gemini as another worker

Set `GEMINI_API_KEY` to your Google AI Studio API key, then enable Gemini:

```sh
export GEMINI_API_KEY='your-gemini-key'
./bin/pipeline generate-hooks --from-staging --stage --resume --openrouter --gemini
```

`--gemini` runs alongside local and OpenRouter workers. It also works without
`--openrouter`. Both `generate-hooks` and `pipeline` support:

- `--gemini-model`: defaults to `gemini-3.5-flash-lite`.
- `--gemini-max-tokens`: independent output budget, default 30,000 tokens.

The helper uses Google's OpenAI-compatible Gemini API endpoint at
`https://generativelanguage.googleapis.com/v1beta/openai/chat/completions`.
Network retries, error diagnostics, model provenance, and local fallback work
independently for each remote worker. Gemini receives the article content or
summaries for its assigned pairs. API access and quotas depend on your Google
project; AI Studio access does not imply unlimited API usage.

Model documentation: https://ai.google.dev/gemini-api/docs/models/gemini-3.5-flash-lite
API documentation: https://ai.google.dev/gemini-api/docs/openai

### Add Cloudflare Workers AI

Cloudflare can run alongside any combination of local, OpenRouter, and Gemini.
Set the account and a Cloudflare API token with Workers AI permission:

```sh
export CLOUDFLARE_ACCOUNT_ID='your-32-character-account-id'
export CLOUDFLARE_API_TOKEN='your-token'
./bin/pipeline generate-hooks --from-staging --stage --resume --openrouter --gemini --cloudflare
```

Both `generate-hooks` and `pipeline` accept:

- `--cloudflare`: enable the additional worker.
- `--cloudflare-model`: defaults to `@cf/google/gemma-4-26b-a4b-it`.
- `--cloudflare-max-tokens`: independent output budget, default 30,000.
- `--cloudflare-account-id`: override `CLOUDFLARE_ACCOUNT_ID`.

The client calls the account's `/ai/v1/chat/completions` REST endpoint directly;
no Worker deployment or Wrangler login is required. The token must authorize
Workers AI inference for that account. It uses the same network retries and
independent local fallback as the other remote workers, and sends the content
or summaries for its assigned pairs to Cloudflare. Account quotas and Workers
AI pricing apply.

Documentation: https://developers.cloudflare.com/workers-ai/configuration/open-ai-compatibility/
Model: https://developers.cloudflare.com/workers-ai/models/gemma-4-26b-a4b-it/
