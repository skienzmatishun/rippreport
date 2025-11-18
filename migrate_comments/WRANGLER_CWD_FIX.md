# Wrangler Working Directory Fix

## Problem

The migration tool was failing to import comments with the error:
```
ERROR | Failed to create comment after 1 attempts: Wrangler command failed (exit 1):
ERROR | Could not parse comment ID from wrangler output
```

## Root Cause

The wrangler CLI was being executed without specifying a working directory. When run from the `migrate_comments/` directory, wrangler couldn't find the `wrangler.toml` configuration file, which is located in `cloudflare_comments/`.

Without the configuration file, wrangler couldn't resolve the database name `comments_db` to the actual database ID.

## Solution

### Fix 1: Add Working Directory Parameter

Added `wrangler_cwd` parameter to `CloudflareCommentClient` to specify the working directory where `wrangler.toml` is located:

```python
def __init__(
    self,
    api_base: str,
    database_id: str,
    database_name: str = "comments_db",
    wrangler_path: str = "wrangler",
    wrangler_cwd: Optional[str] = None,  # NEW
    max_retries: int = 3,
    backoff_factor: float = 2.0,
    logger: Optional[logging.Logger] = None
):
    ...
    self.wrangler_cwd = wrangler_cwd
```

### Fix 2: Pass Working Directory to subprocess

Updated `_execute_d1_command` to use the working directory:

```python
result = subprocess.run(
    cmd,
    capture_output=True,
    text=True,
    timeout=30,
    cwd=self.wrangler_cwd  # Run from directory with wrangler.toml
)
```

### Fix 3: Add RETURNING Clause

Added `RETURNING id` to the INSERT statement to get the new comment ID:

```python
sql = f"""
INSERT INTO comments (page_id, parent_id, author_name, content, created_at, ip_address, content_hash)
VALUES ('{page_id_escaped}', {parent_sql}, '{author_escaped}', '{content_escaped}', '{created_at}', '{ip_escaped}', {hash_sql})
RETURNING id;
""".strip()
```

### Fix 4: Update Configuration

Added `wrangler_cwd` to the configuration:

```yaml
cloudflare:
  api_base: "https://comments.rippreport.com/api"
  wrangler_path: "wrangler"
  wrangler_cwd: "../cloudflare_comments"  # Directory with wrangler.toml
  database_name: "comments_db"
  database_id: "4c94fdd6-f883-439a-944c-a63a5cffac9c"
```

### Fix 5: Update Orchestrator

Updated the orchestrator to pass `wrangler_cwd` to the client:

```python
self.importer = CloudflareCommentClient(
    api_base=cloudflare_config['api_base'],
    database_id=cloudflare_config['database_id'],
    database_name=cloudflare_config.get('database_name', 'comments_db'),
    wrangler_path=cloudflare_config.get('wrangler_path', 'wrangler'),
    wrangler_cwd=cloudflare_config.get('wrangler_cwd'),  # NEW
    max_retries=migration_config.get('max_retries', 3),
    logger=self.logger
)
```

## Changes Made

### Files Modified

1. **migrate_comments/importers/cloudflare_client.py**
   - Added `wrangler_cwd` parameter to `__init__`
   - Updated `_execute_d1_command` to use `cwd` parameter
   - Added `RETURNING id` to INSERT statement

2. **migrate_comments/config.yaml**
   - Added `wrangler_cwd: "../cloudflare_comments"`

3. **migrate_comments/orchestrator.py**
   - Added `wrangler_cwd` parameter when creating `CloudflareCommentClient`

## Verification

### Before Fix
```
Total Pages:      1
Successful Pages: 0
Failed Pages:     1
Comments Imported: 0
Failed Imports:    5
```

### After Fix
```
Total Pages:      1
Successful Pages: 1
Failed Pages:     0
Comments Imported: 5
Failed Imports:    0
```

### Database Verification

```bash
wrangler d1 execute comments_db --command "SELECT COUNT(*) FROM comments" --json
# Result: 5 comments successfully imported
```

## Technical Details

### Why Working Directory Matters

Wrangler CLI looks for `wrangler.toml` in the current working directory to:
1. Resolve database names to database IDs
2. Get environment configuration
3. Access bindings and other settings

Without the correct working directory, wrangler can't find the configuration and fails with:
```json
{
  "error": {
    "text": "Couldn't find a D1 DB with the name or binding 'comments_db' in your Wrangler configuration file."
  }
}
```

### RETURNING Clause

SQLite supports the `RETURNING` clause to get values from inserted rows. This is more reliable than using `last_insert_rowid()` because:
1. It's atomic with the INSERT
2. It works correctly with triggers
3. It returns the actual inserted ID

The wrangler output format with RETURNING is:
```json
[
  {
    "results": [
      {
        "id": 6
      }
    ],
    "success": true,
    "meta": {
      "duration": 1
    }
  }
]
```

## Impact

This fix enables the migration tool to:
- ✅ Successfully execute wrangler commands
- ✅ Import comments to the D1 database
- ✅ Retrieve the new comment IDs
- ✅ Map Matrix event IDs to new comment IDs
- ✅ Complete the full migration workflow

## Date

Fixed: November 17, 2024

## Status

✅ **RESOLVED** - Migration tool now successfully imports comments to Cloudflare D1.
