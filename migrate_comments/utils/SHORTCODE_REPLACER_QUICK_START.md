# Shortcode Replacer Quick Start Guide

## What It Does

The Shortcode Replacer automatically updates your Hugo content files to use the new `aicomments` shortcode instead of the legacy `chat` shortcode.

**Before:**
```markdown
{{< chat "my-post-id" >}}
```

**After:**
```markdown
{{< aicomments "my-post-id" >}}
```

## Quick Start

### 1. Preview Changes (Dry Run)

```bash
python migrate_comments/utils/demo_shortcode_replacer.py
```

This will:
- Scan all markdown files in `content/`
- Show you which files will be modified
- Display how many shortcodes will be replaced
- **NOT modify any files**

### 2. Review the Output

Look for:
- Number of files found
- Number of shortcodes to replace
- Any errors or warnings

Example output:
```
Found 118 file(s) with chat shortcodes:

1. p/my-post/index.md
   Line 97: {{< chat my-id >}}
   → Will replace with: {{< aicomments "my-id" >}}

Total: 118 files, 118 shortcodes
```

### 3. Run the Replacement

```python
from utils.shortcode_replacer import ShortcodeReplacer

replacer = ShortcodeReplacer("content")
replacer.scan_files()
stats = replacer.replace_all(dry_run=False)

print(f"✓ Replaced {stats['total_replacements']} shortcodes")
print(f"✓ Modified {stats['successful']} files")
```

### 4. Verify Changes

```bash
# Check that no chat shortcodes remain
grep -r "{{< chat" content/

# Should return no results if successful
```

### 5. Test Your Site

```bash
# Build Hugo site
hugo

# Check for any errors
# Test a few pages with comments
```

## Rollback (If Needed)

If something goes wrong, you can restore from backup:

```python
from pathlib import Path
from utils.shortcode_replacer import ShortcodeReplacer

replacer = ShortcodeReplacer("content")

# Find the backup metadata file
backup_file = Path("backups/shortcode-replacement/20241117-120000/backup_metadata.json")

# Rollback
stats = replacer.rollback(backup_file)
print(f"✓ Restored {stats['successful']} files")
```

## Common Use Cases

### Replace Shortcodes in Specific Directory

```python
replacer = ShortcodeReplacer("content/p")  # Only posts
replacer.scan_files()
replacer.replace_all(dry_run=False)
```

### Check What Would Be Changed

```python
replacer = ShortcodeReplacer("content")
files = replacer.scan_files()

for file_info in files:
    print(f"\n{file_info['relative_path']}:")
    for sc in file_info['shortcodes']:
        print(f"  Line {sc['line_number']}: {sc['id']}")
```

### Verify Specific File

```python
from pathlib import Path

file_path = Path("content/p/my-post/index.md")
success, issues = replacer.verify_replacement(file_path)

if success:
    print("✓ Replacement verified")
else:
    print("✗ Issues found:")
    for issue in issues:
        print(f"  - {issue}")
```

## Safety Features

1. **Backups**: Original files are backed up before modification
2. **Dry Run**: Preview changes without modifying files
3. **Verification**: Automatic verification after replacement
4. **Rollback**: Easy restoration from backups
5. **Logging**: All operations are logged

## Backup Location

Backups are stored in:
```
backups/shortcode-replacement/
└── YYYYMMDD-HHMMSS/
    ├── backup_metadata.json
    └── content/
        └── [original directory structure]
```

## Troubleshooting

### No files found
- Check that `content/` directory exists
- Verify files contain `{{< chat` shortcodes

### Replacement failed
- Check file permissions
- Ensure files are not open in editor
- Review error messages in output

### Verification issues
- Check `verification_issues` in stats
- Manually inspect affected files
- Run verification again after fixes

## Next Steps

After successful replacement:

1. ✅ Commit changes to git
2. ✅ Deploy updated site
3. ✅ Test comments on live site
4. ✅ Keep backups for a while
5. ✅ Delete old backups after verification

## Need Help?

Check the full documentation:
- `SHORTCODE_REPLACER_IMPLEMENTATION.md` - Complete API reference
- `test_shortcode_replacer.py` - Usage examples
- `demo_shortcode_replacer.py` - Demo script
