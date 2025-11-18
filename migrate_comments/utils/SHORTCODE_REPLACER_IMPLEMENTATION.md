# Shortcode Replacer Implementation

## Overview

The `ShortcodeReplacer` class provides functionality to scan Hugo content files for legacy `chat` shortcodes and replace them with the new `aicomments` shortcodes. It includes comprehensive backup, verification, and rollback capabilities.

## Features

### 1. Hugo File Scanner (Task 11.1)
- ✅ Recursively scans content directory for markdown files
- ✅ Identifies files containing chat shortcodes
- ✅ Parses shortcode parameters (with or without quotes)
- ✅ Tracks files to be modified with line numbers

### 2. Shortcode Replacement (Task 11.2)
- ✅ Replaces `{{< chat "id" >}}` with `{{< aicomments "id" >}}`
- ✅ Preserves shortcode parameters exactly
- ✅ Handles multiple shortcodes per file
- ✅ Validates replacement correctness

### 3. Backup Functionality (Task 11.3)
- ✅ Creates timestamped backup directories
- ✅ Copies original files before modification
- ✅ Stores backup metadata (timestamp, file paths)
- ✅ Provides rollback capability

### 4. Replacement Verification (Task 11.4)
- ✅ Verifies shortcodes were replaced correctly
- ✅ Checks files are still valid markdown
- ✅ Ensures no duplicate replacements
- ✅ Logs all replacements and issues

## Usage

### Basic Usage

```python
from utils.shortcode_replacer import ShortcodeReplacer

# Initialize replacer
replacer = ShortcodeReplacer(
    content_dir="content",
    backup_dir="backups/shortcode-replacement"
)

# Scan for files with chat shortcodes
files = replacer.scan_files()
print(f"Found {len(files)} files with chat shortcodes")

# Dry run to preview changes
stats = replacer.replace_all(dry_run=True)
print(f"Would replace {stats['total_replacements']} shortcodes")

# Actually replace shortcodes
stats = replacer.replace_all(dry_run=False)
print(f"Replaced {stats['total_replacements']} shortcodes in {stats['successful']} files")
```

### Rollback Changes

```python
from pathlib import Path

# Rollback using backup metadata
metadata_file = Path("backups/shortcode-replacement/20241117-120000/backup_metadata.json")
stats = replacer.rollback(metadata_file)
print(f"Restored {stats['successful']} files")
```

### Command Line Demo

```bash
# Preview changes (dry run)
python migrate_comments/utils/demo_shortcode_replacer.py

# Actually replace shortcodes (add --replace flag when implemented)
python migrate_comments/utils/demo_shortcode_replacer.py --replace
```

## Implementation Details

### Shortcode Pattern Matching

The replacer uses a flexible regex pattern that matches various shortcode formats:

```python
CHAT_SHORTCODE_PATTERN = re.compile(
    r'{{\s*<\s*chat\s+["\']?([^"\'\s>]+)["\']?\s*>}}',
    re.IGNORECASE
)
```

This matches:
- `{{< chat "id" >}}` - with quotes
- `{{< chat id >}}` - without quotes
- `{{<chat "id">}}` - compact format
- `{{< CHAT "id" >}}` - case insensitive

### Backup Structure

Backups are organized by timestamp:

```
backups/shortcode-replacement/
└── 20241117-120000/
    ├── backup_metadata.json
    └── content/
        └── p/
            └── post-name/
                └── index.md
```

### Metadata Format

```json
{
  "timestamp": "2024-11-17T12:00:00",
  "files": [
    {
      "original": "/path/to/content/p/post/index.md",
      "backup": "/path/to/backups/.../content/p/post/index.md",
      "relative_path": "p/post/index.md"
    }
  ]
}
```

## API Reference

### ShortcodeReplacer Class

#### `__init__(content_dir: str, backup_dir: str = "backups/shortcode-replacement")`

Initialize the replacer.

**Parameters:**
- `content_dir`: Path to Hugo content directory
- `backup_dir`: Path to backup directory (default: "backups/shortcode-replacement")

**Raises:**
- `ValueError`: If content directory doesn't exist

#### `scan_files() -> List[Dict]`

Scan content directory for files with chat shortcodes.

**Returns:**
List of file info dictionaries:
```python
{
    "path": Path,
    "relative_path": str,
    "shortcodes": [
        {
            "id": str,
            "match": str,
            "line_number": int
        }
    ]
}
```

#### `replace_shortcodes(file_path: Path, dry_run: bool = False) -> Tuple[bool, int]`

Replace chat shortcodes in a single file.

**Parameters:**
- `file_path`: Path to the file
- `dry_run`: If True, don't modify the file

**Returns:**
- Tuple of (success: bool, replacement_count: int)

#### `replace_all(dry_run: bool = False) -> Dict`

Replace shortcodes in all scanned files.

**Parameters:**
- `dry_run`: If True, don't modify files

**Returns:**
Statistics dictionary:
```python
{
    "total_files": int,
    "successful": int,
    "failed": int,
    "total_replacements": int,
    "errors": List[str],
    "verification_issues": List[Dict]
}
```

#### `verify_replacement(file_path: Path) -> Tuple[bool, List[str]]`

Verify that replacement was successful.

**Parameters:**
- `file_path`: Path to the file to verify

**Returns:**
- Tuple of (success: bool, issues: List[str])

#### `create_backup(file_path: Path) -> Path`

Create a backup of a file.

**Parameters:**
- `file_path`: Path to the file to backup

**Returns:**
- Path to the backup file

#### `rollback(backup_metadata_file: Path) -> Dict`

Rollback files from a backup.

**Parameters:**
- `backup_metadata_file`: Path to backup_metadata.json

**Returns:**
Statistics dictionary:
```python
{
    "total_files": int,
    "successful": int,
    "failed": int,
    "errors": List[str]
}
```

## Testing

### Run Unit Tests

```bash
python -m pytest migrate_comments/utils/test_shortcode_replacer.py -v
```

### Test Coverage

The test suite includes:
- ✅ Initialization with valid/invalid directories
- ✅ Scanning for files with various shortcode formats
- ✅ Extracting shortcodes with line numbers
- ✅ Replacing shortcodes (dry run and actual)
- ✅ Handling multiple shortcodes per file
- ✅ Creating backups with metadata
- ✅ Verifying replacements (success and failure cases)
- ✅ Replacing all files in batch
- ✅ Rolling back changes

All 15 tests pass successfully.

## Integration with Migration Tool

The shortcode replacer can be integrated into the main migration orchestrator:

```python
from utils.shortcode_replacer import ShortcodeReplacer

# After successful comment migration
if migration_successful:
    replacer = ShortcodeReplacer(config['hugo']['content_dir'])
    replacer.scan_files()
    
    # Preview changes
    stats = replacer.replace_all(dry_run=True)
    logger.info(f"Will replace {stats['total_replacements']} shortcodes")
    
    # Confirm with user
    if confirm_replacement():
        stats = replacer.replace_all(dry_run=False)
        logger.info(f"Replaced {stats['total_replacements']} shortcodes")
```

## Error Handling

The replacer includes comprehensive error handling:

1. **File Access Errors**: Logged and skipped, doesn't stop batch processing
2. **Backup Failures**: Prevents modification if backup fails
3. **Verification Issues**: Logged but doesn't rollback automatically
4. **Rollback Errors**: Individual file failures logged, continues with remaining files

## Best Practices

1. **Always run dry run first**: Preview changes before modifying files
2. **Verify backups exist**: Check backup directory before proceeding
3. **Test on a few files first**: Use a subset of files for initial testing
4. **Keep backups**: Don't delete backup directories until verified
5. **Use version control**: Commit changes before running replacer

## Verification Checklist

After running the replacer:

- [ ] All chat shortcodes replaced with aicomments
- [ ] No remaining `{{< chat` patterns in content
- [ ] Shortcode IDs preserved correctly
- [ ] Files are valid markdown
- [ ] Hugo site builds successfully
- [ ] Comments display correctly on site
- [ ] Backup metadata saved
- [ ] No duplicate replacements

## Troubleshooting

### Issue: Shortcodes not found

**Solution**: Check that content_dir path is correct and files contain `{{< chat` patterns.

### Issue: Replacement verification fails

**Solution**: Check verification issues in stats['verification_issues'] for specific problems.

### Issue: Backup directory full

**Solution**: Clean up old backup directories or specify a different backup location.

### Issue: Rollback fails

**Solution**: Ensure backup_metadata.json exists and backup files are intact.

## Requirements Satisfied

This implementation satisfies all requirements from Requirement 8:

- ✅ 8.1: Scan Hugo content files for chat shortcode
- ✅ 8.2: Replace with aicomments shortcode
- ✅ 8.3: Preserve page section ID parameter
- ✅ 8.4: Create backup of original file
- ✅ 8.5: Skip files already updated

## Future Enhancements

Potential improvements:

1. Add command-line interface for standalone usage
2. Support for other shortcode types
3. Parallel processing for large sites
4. Interactive mode with file-by-file confirmation
5. Integration with git for automatic commits
6. Automatic Hugo site rebuild and validation
