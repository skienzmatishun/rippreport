# Configuration Guide

This guide explains all configuration options for the Comment Migration Tool.

## Configuration File

The tool uses a YAML configuration file (`config.yaml`) to control its behavior. All settings can also be overridden using environment variables.

## Quick Start

1. Copy the example configuration:
   ```bash
   cp config.yaml.example config.yaml
   ```

2. Edit the required fields (marked with ⚠️ below)

3. Verify your configuration:
   ```bash
   python -c "from utils.config import Config; c = Config('config.yaml'); print('Valid!')"
   ```

## Configuration Sections

### Matrix Configuration

Settings for connecting to the Cactus Chat (Matrix) backend to extract existing comments.

```yaml
matrix:
  homeserver_url: "https://matrix.cactus.chat:8448"
  server_name: "cactus.chat"
  site_name: "your-site.com"  # ⚠️ REQUIRED
```

#### `homeserver_url`
- **Type**: String (URL)
- **Default**: `"https://matrix.cactus.chat:8448"`
- **Description**: The Matrix homeserver URL for Cactus Chat
- **Required**: No (default is usually correct)
- **Example**: `"https://matrix.cactus.chat:8448"`

#### `server_name`
- **Type**: String
- **Default**: `"cactus.chat"`
- **Description**: The Matrix server name
- **Required**: No (default is usually correct)
- **Example**: `"cactus.chat"`

#### `site_name`
- **Type**: String
- **Default**: None
- **Description**: Your site's domain name (used to construct Matrix room IDs)
- **Required**: ⚠️ **YES**
- **Example**: `"rippreport.com"` or `"myblog.com"`
- **Notes**: This must match the site name you used when setting up Cactus Chat

**Environment Variable Overrides:**
- `MATRIX_HOMESERVER_URL`
- `MATRIX_SERVER_NAME`
- `MATRIX_SITE_NAME`

---

### Cloudflare Configuration

Settings for connecting to your Cloudflare Workers comment system.

```yaml
cloudflare:
  api_base: "https://comments.your-site.com/api"  # ⚠️ REQUIRED
  database_id: "your-database-id"                  # ⚠️ REQUIRED
  wrangler_path: "wrangler"
  database_name: "comments_db"
```

#### `api_base`
- **Type**: String (URL)
- **Default**: None
- **Description**: The base URL for your Cloudflare Workers comment API
- **Required**: ⚠️ **YES**
- **Example**: `"https://comments.rippreport.com/api"`
- **Notes**: This is the URL where your comment worker is deployed

#### `database_id`
- **Type**: String (UUID)
- **Default**: None
- **Description**: Your Cloudflare D1 database ID
- **Required**: ⚠️ **YES**
- **Example**: `"4c94fdd6-f883-439a-944c-a63a5cffac9c"`
- **How to find**: Run `wrangler d1 list` to see your databases

#### `wrangler_path`
- **Type**: String (path)
- **Default**: `"wrangler"`
- **Description**: Path to the Wrangler CLI executable
- **Required**: No (uses system PATH by default)
- **Example**: `"wrangler"` or `"/usr/local/bin/wrangler"`
- **Notes**: Only change if Wrangler is not in your PATH

#### `database_name`
- **Type**: String
- **Default**: `"comments_db"`
- **Description**: The name of your D1 database
- **Required**: No (used for logging only)
- **Example**: `"comments_db"` or `"prod_comments"`

**Environment Variable Overrides:**
- `CLOUDFLARE_API_BASE`
- `CLOUDFLARE_DATABASE_ID`
- `CLOUDFLARE_DATABASE_NAME`
- `WRANGLER_PATH`

---

### Hugo Configuration

Settings for your Hugo static site.

```yaml
hugo:
  content_dir: "content"
  base_url: "https://your-site.com"  # ⚠️ REQUIRED
```

#### `content_dir`
- **Type**: String (path)
- **Default**: `"content"`
- **Description**: Path to your Hugo content directory (relative to project root)
- **Required**: No (default works for standard Hugo sites)
- **Example**: `"content"` or `"site/content"`
- **Notes**: The tool scans this directory for pages with chat shortcodes

#### `base_url`
- **Type**: String (URL)
- **Default**: None
- **Description**: Your site's base URL
- **Required**: ⚠️ **YES**
- **Example**: `"https://rippreport.com"` or `"https://myblog.com"`
- **Notes**: Used to construct page IDs and verify permalinks

**Environment Variable Overrides:**
- `HUGO_CONTENT_DIR`
- `HUGO_BASE_URL`

---

### Migration Settings

Control the migration process behavior.

```yaml
migration:
  dry_run: false
  batch_size: 10
  delay_between_batches: 1.0
  max_retries: 3
  checkpoint_file: ".migration_checkpoint.json"
  report_dir: "migration_reports"
```

#### `dry_run`
- **Type**: Boolean
- **Default**: `false`
- **Description**: If true, extract and transform comments without importing
- **Required**: No
- **Example**: `true` or `false`
- **Notes**: Use dry-run mode to preview the migration before committing

#### `batch_size`
- **Type**: Integer
- **Default**: `10`
- **Description**: Number of comments to process in each batch
- **Required**: No
- **Example**: `10` (recommended), `5` (slower, safer), `20` (faster, riskier)
- **Notes**: Smaller batches are safer but slower; larger batches are faster but may hit rate limits

#### `delay_between_batches`
- **Type**: Float (seconds)
- **Default**: `1.0`
- **Description**: Delay in seconds between processing batches
- **Required**: No
- **Example**: `1.0` (recommended), `0.5` (faster), `2.0` (safer)
- **Notes**: Helps avoid rate limiting; increase if you encounter 429 errors

#### `max_retries`
- **Type**: Integer
- **Default**: `3`
- **Description**: Maximum number of retry attempts for failed operations
- **Required**: No
- **Example**: `3` (recommended), `5` (more resilient), `1` (fail fast)
- **Notes**: Uses exponential backoff (1s, 2s, 4s, 8s)

#### `checkpoint_file`
- **Type**: String (path)
- **Default**: `".migration_checkpoint.json"`
- **Description**: Path to the checkpoint file for recovery
- **Required**: No
- **Example**: `".migration_checkpoint.json"` or `"checkpoints/migration.json"`
- **Notes**: Allows resuming from failures; automatically deleted on success

#### `report_dir`
- **Type**: String (path)
- **Default**: `"migration_reports"`
- **Description**: Directory for migration reports
- **Required**: No
- **Example**: `"migration_reports"` or `"reports"`
- **Notes**: Created automatically if it doesn't exist

**Environment Variable Overrides:**
- `MIGRATION_DRY_RUN` (set to "true" or "false")
- `MIGRATION_BATCH_SIZE`
- `MIGRATION_MAX_RETRIES`

---

### Logging Configuration

Control logging behavior.

```yaml
logging:
  level: "INFO"
  file: "migration.log"
  console: true
```

#### `level`
- **Type**: String
- **Default**: `"INFO"`
- **Description**: Minimum log level to record
- **Required**: No
- **Options**: `"DEBUG"`, `"INFO"`, `"WARNING"`, `"ERROR"`, `"CRITICAL"`
- **Example**: `"INFO"` (normal), `"DEBUG"` (troubleshooting), `"WARNING"` (quiet)
- **Notes**: Use DEBUG for troubleshooting, INFO for normal operation

#### `file`
- **Type**: String (path)
- **Default**: `"migration.log"`
- **Description**: Path to the log file
- **Required**: No
- **Example**: `"migration.log"` or `"logs/migration.log"`
- **Notes**: Automatically rotated when it reaches 10MB

#### `console`
- **Type**: Boolean
- **Default**: `true`
- **Description**: Whether to also log to console (stdout)
- **Required**: No
- **Example**: `true` (see logs in terminal), `false` (file only)
- **Notes**: Console logs are color-coded if colorlog is installed

**Environment Variable Overrides:**
- None (logging must be configured in YAML)

---

## Environment Variables

All configuration options can be overridden using environment variables. This is useful for:
- CI/CD pipelines
- Different environments (dev, staging, prod)
- Keeping secrets out of config files

### Setting Environment Variables

**On macOS/Linux:**
```bash
export MATRIX_SITE_NAME="myblog.com"
export CLOUDFLARE_API_BASE="https://comments.myblog.com/api"
export CLOUDFLARE_DATABASE_ID="your-database-id"
export HUGO_BASE_URL="https://myblog.com"
```

**On Windows (Command Prompt):**
```cmd
set MATRIX_SITE_NAME=myblog.com
set CLOUDFLARE_API_BASE=https://comments.myblog.com/api
set CLOUDFLARE_DATABASE_ID=your-database-id
set HUGO_BASE_URL=https://myblog.com
```

**On Windows (PowerShell):**
```powershell
$env:MATRIX_SITE_NAME="myblog.com"
$env:CLOUDFLARE_API_BASE="https://comments.myblog.com/api"
$env:CLOUDFLARE_DATABASE_ID="your-database-id"
$env:HUGO_BASE_URL="https://myblog.com"
```

### Environment Variable Priority

Environment variables take precedence over config file values:

1. **Environment variable** (highest priority)
2. **config.yaml value**
3. **Default value** (lowest priority)

---

## Complete Example Configuration

Here's a complete, annotated example configuration:

```yaml
# Comment Migration Configuration

# Matrix (Cactus Chat) Configuration
matrix:
  # Default Cactus Chat homeserver (usually don't need to change)
  homeserver_url: "https://matrix.cactus.chat:8448"
  
  # Default Cactus Chat server name (usually don't need to change)
  server_name: "cactus.chat"
  
  # YOUR SITE NAME - must match what you used in Cactus Chat setup
  site_name: "rippreport.com"
  
# Cloudflare Configuration
cloudflare:
  # Your Cloudflare Workers comment API endpoint
  api_base: "https://comments.rippreport.com/api"
  
  # Path to wrangler CLI (usually just "wrangler")
  wrangler_path: "wrangler"
  
  # Your D1 database name (for logging)
  database_name: "comments_db"
  
  # Your D1 database ID (get from: wrangler d1 list)
  database_id: "4c94fdd6-f883-439a-944c-a63a5cffac9c"
  
# Hugo Configuration
hugo:
  # Path to Hugo content directory (relative to project root)
  content_dir: "content"
  
  # Your site's base URL
  base_url: "https://rippreport.com"
  
# Migration Settings
migration:
  # Set to true to preview without importing
  dry_run: false
  
  # Number of comments per batch (10 is safe)
  batch_size: 10
  
  # Delay between batches in seconds (helps avoid rate limits)
  delay_between_batches: 1.0
  
  # Maximum retry attempts for failed operations
  max_retries: 3
  
  # Checkpoint file for recovery (auto-deleted on success)
  checkpoint_file: ".migration_checkpoint.json"
  
  # Directory for migration reports
  report_dir: "migration_reports"
  
# Logging
logging:
  # Log level: DEBUG, INFO, WARNING, ERROR, CRITICAL
  level: "INFO"
  
  # Log file path
  file: "migration.log"
  
  # Also log to console
  console: true
```

---

## Configuration Validation

The tool validates your configuration on startup. Common validation errors:

### Missing Required Fields

**Error**: `Configuration error: matrix.site_name is required`

**Solution**: Add the missing field to your config.yaml

### Invalid URLs

**Error**: `Configuration error: cloudflare.api_base must be a valid URL`

**Solution**: Ensure URLs start with `http://` or `https://`

### Invalid Types

**Error**: `Configuration error: migration.batch_size must be an integer`

**Solution**: Remove quotes from numeric values

### Invalid Paths

**Error**: `Configuration error: hugo.content_dir does not exist`

**Solution**: Ensure the path exists and is relative to the project root

---

## Testing Your Configuration

After configuring the tool, test it:

```bash
# Test configuration loading
python -c "from utils.config import Config; c = Config('config.yaml'); print('Valid!')"

# Test Matrix connection (dry-run)
python migrate.py --dry-run --pages test-page

# Test Wrangler access
wrangler d1 execute comments_db --command "SELECT COUNT(*) FROM comments"
```

---

## Security Best Practices

1. **Don't commit secrets**: Add `config.yaml` to `.gitignore`
2. **Use environment variables**: For sensitive values in CI/CD
3. **Restrict file permissions**: `chmod 600 config.yaml`
4. **Use separate configs**: Different files for dev/staging/prod

---

## Next Steps

After configuration:

1. Read [USAGE.md](USAGE.md) to learn how to run migrations
2. Run a dry-run to preview: `python migrate.py --dry-run`
3. Check [TROUBLESHOOTING.md](TROUBLESHOOTING.md) if you encounter issues
