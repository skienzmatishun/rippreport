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
- `--stage`: Save to staging files instead of writing directly
- `--from-staging`: Read rankings from staging files instead of front matter
- `--dry-run`: Simulate without writing any files
- `--use-compressed`: Use compressed article summaries for context
- `--hook-model`: Override the model used for generation

### apply-staged
- `--hooks`: Apply hooks instead of rankings
- `--posts`: Comma-separated list of post slugs (omit to apply all)

### clear-staged
- `--hooks`: Clear hooks instead of rankings  
- `--posts`: Comma-separated list of post slugs (omit to clear all)
