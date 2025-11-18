# Comment Migration Tool - Usage Guide

## Quick Start

### 1. Configure

Copy the example configuration file and fill in your values:

```bash
cp migrate_comments/config.yaml.example migrate_comments/config.yaml
```

Edit `config.yaml` with your settings:
- Matrix homeserver URL and site name
- Cloudflare API base URL and database ID
- Hugo content directory and base URL

### 2. Test with Dry Run

Preview the migration without importing any data:

```bash
python migrate_comments/migrate.py --dry-run
```

This will:
- Extract comments from Matrix
- Transform to new format
- Validate all data
- Check for issues
- Generate preview reports
- Provide a go/no-go recommendation

### 3. Run Migration

If the dry run looks good, run the actual migration:

```bash
python migrate_comments/migrate.py
```

The migration will:
- Process all pages with chat shortcodes
- Import comments with preserved timestamps
- Show real-time progress
- Save checkpoints after each page
- Generate comprehensive reports

### 4. Review Results

Check the generated reports in `migration_reports/`:
- `summary_report_*.txt` - High-level statistics
- `detailed_report_*.txt` - Per-page breakdown
- `id_mapping_*.txt` - Matrix event ID to comment ID mapping
- `migration_report_*.json` - Complete data in JSON format

## Command-Line Options

### Basic Usage

```bash
python migrate_comments/migrate.py [options]
```

### Options

#### `--dry-run`
Run migration without importing data (preview mode)

```bash
python migrate_comments/migrate.py --dry-run
```

#### `--resume`
Resume migration from previous checkpoint

```bash
python migrate_comments/migrate.py --resume
```

Use this if migration was interrupted or failed partway through.

#### `--pages PAGE_ID [PAGE_ID ...]`
Migrate only specific pages

```bash
python migrate_comments/migrate.py --pages /p/article-1 /p/article-2
```

#### `--config FILE`
Use custom configuration file

```bash
python migrate_comments/migrate.py --config my_config.yaml
```

#### `--verbose` or `-v`
Enable verbose logging (DEBUG level)

```bash
python migrate_comments/migrate.py --verbose
```

#### `--quiet` or `-q`
Suppress console output (log to file only)

```bash
python migrate_comments/migrate.py --quiet
```

#### `--help` or `-h`
Show help message

```bash
python migrate_comments/migrate.py --help
```

#### `--version`
Show version information

```bash
python migrate_comments/migrate.py --version
```

## Common Workflows

### Preview Before Migrating

```bash
# 1. Test with dry run
python migrate_comments/migrate.py --dry-run

# 2. Review the preview report
cat dry_run_reports/dry_run_preview_*.txt

# 3. If everything looks good, run the migration
python migrate_comments/migrate.py
```

### Migrate Specific Pages First

```bash
# Test with a few pages first
python migrate_comments/migrate.py --pages /p/test-1 /p/test-2

# If successful, migrate all pages
python migrate_comments/migrate.py
```

### Resume After Failure

```bash
# If migration fails or is interrupted
python migrate_comments/migrate.py --resume
```

The tool will:
- Skip already completed pages
- Load ID mappings from checkpoint
- Continue from where it left off

### Debug Issues

```bash
# Run with verbose logging
python migrate_comments/migrate.py --verbose

# Check the log file
tail -f migration.log
```

## Progress Tracking

During migration, you'll see a real-time progress bar:

```
[████████████████████████████████░░░░░░░░] 80.0% | 8/10 pages | ✓ 7 ✗ 1 | ETA: 2m 15s
```

This shows:
- Visual progress bar
- Percentage complete
- Pages processed / total pages
- Success count (✓)
- Failure count (✗)
- Estimated time remaining

## Output Files

### Migration Reports

Generated in `migration_reports/` directory:

1. **Summary Report** (`summary_report_*.txt`)
   - Total pages and comments
   - Success rate
   - Failed pages summary
   - Overall statistics

2. **Detailed Report** (`detailed_report_*.txt`)
   - Per-page breakdown
   - Timing for each operation
   - Error messages for failures

3. **ID Mapping Report** (`id_mapping_*.txt`)
   - Matrix event ID to new comment ID mapping
   - Useful for debugging parent-child relationships

4. **JSON Report** (`migration_report_*.json`)
   - Complete migration data in JSON format
   - Can be processed programmatically

5. **Failed Comments** (`failed_comments_*.json`)
   - Details of any failed imports
   - Only generated if there are failures

### Dry Run Reports

Generated in `dry_run_reports/` directory:

1. **Preview Report** (`dry_run_preview_*.txt`)
   - What would be migrated
   - Potential issues and warnings
   - Sample transformed comments
   - Go/no-go recommendation

2. **JSON Results** (`dry_run_results_*.json`)
   - Complete dry run data in JSON format

### Checkpoint File

`.migration_checkpoint.json` - Saved progress for resume functionality

This file is automatically:
- Created when migration starts
- Updated after each page completes
- Loaded when using `--resume`
- Deleted when migration completes successfully

### Log File

`migration.log` - Detailed log of all operations

Contains:
- Timestamped log entries
- Debug information (if `--verbose` used)
- Error messages and stack traces
- Operation timing

## Troubleshooting

### Configuration Errors

If you see "Configuration error":
1. Check that `config.yaml` exists
2. Verify all required fields are filled in
3. Check YAML syntax is valid

### Matrix Connection Errors

If extraction fails:
1. Verify Matrix homeserver URL is correct
2. Check that site name matches your Cactus Chat setup
3. Ensure Matrix rooms exist for the pages

### Cloudflare Import Errors

If imports fail:
1. Verify wrangler is installed and authenticated
2. Check database ID is correct
3. Ensure API base URL is accessible
4. Check for rate limiting (tool will retry automatically)

### Validation Errors

If validation fails:
1. Review the detailed error messages
2. Check for data integrity issues
3. Consider fixing source data before migrating

### Resume Not Working

If `--resume` doesn't work:
1. Check that `.migration_checkpoint.json` exists
2. Verify checkpoint file is not corrupted
3. Try deleting checkpoint and starting fresh

## Best Practices

### Before Migration

1. **Backup your data**
   - Export current D1 database
   - Save Hugo content files
   - Document current state

2. **Test thoroughly**
   - Run dry-run mode first
   - Test with a few pages
   - Review all reports

3. **Plan for downtime**
   - Migration can take time for large sites
   - Consider running during low-traffic periods

### During Migration

1. **Monitor progress**
   - Watch the progress bar
   - Check log file for errors
   - Be ready to intervene if needed

2. **Don't interrupt**
   - Let migration complete if possible
   - If you must stop, use Ctrl+C (checkpoint will be saved)

### After Migration

1. **Verify results**
   - Check reports for any failures
   - Spot-check migrated comments on site
   - Verify parent-child relationships

2. **Handle failures**
   - Review failed comments report
   - Fix issues and use `--resume`
   - Consider manual intervention for edge cases

3. **Clean up**
   - Archive migration reports
   - Remove checkpoint file (if not auto-deleted)
   - Update Hugo shortcodes (see task 11)

## Performance Tips

### Speed Up Migration

1. **Reduce delay between batches**
   ```yaml
   migration:
     delay_between_batches: 0.5  # Default is 1.0
   ```

2. **Increase batch size**
   ```yaml
   migration:
     batch_size: 20  # Default is 10
   ```

3. **Run on faster network**
   - Closer to Matrix homeserver
   - Closer to Cloudflare edge

### Handle Large Migrations

1. **Migrate in batches**
   ```bash
   # Migrate 10 pages at a time
   python migrate_comments/migrate.py --pages /p/page-1 /p/page-2 ... /p/page-10
   ```

2. **Use resume functionality**
   - Migration will checkpoint after each page
   - Can stop and resume anytime

3. **Monitor resources**
   - Watch memory usage
   - Check disk space for logs
   - Monitor network bandwidth

## Exit Codes

- `0` - Success
- `1` - Error (check logs)
- `130` - Interrupted by user (Ctrl+C)

## Getting Help

If you encounter issues:

1. Check this usage guide
2. Review the implementation summary (`TASK_10_IMPLEMENTATION_SUMMARY.md`)
3. Check the log file for detailed errors
4. Run with `--verbose` for more information
5. Review the design document (`.kiro/specs/comment-migration/design.md`)

## Examples

### Complete Migration Workflow

```bash
# 1. Configure
cp migrate_comments/config.yaml.example migrate_comments/config.yaml
# Edit config.yaml with your settings

# 2. Test with dry run
python migrate_comments/migrate.py --dry-run

# 3. Review preview
cat dry_run_reports/dry_run_preview_*.txt

# 4. Test with one page
python migrate_comments/migrate.py --pages /p/test-article

# 5. Verify the test page
# Check the comment on your site

# 6. Run full migration
python migrate_comments/migrate.py

# 7. Review results
cat migration_reports/summary_report_*.txt

# 8. Handle any failures
python migrate_comments/migrate.py --resume
```

### Debug a Specific Page

```bash
# Run with verbose logging for one page
python migrate_comments/migrate.py --verbose --pages /p/problematic-page

# Check the log
tail -100 migration.log
```

### Migrate in Stages

```bash
# Stage 1: Recent posts (test)
python migrate_comments/migrate.py --pages /p/recent-1 /p/recent-2

# Stage 2: Popular posts
python migrate_comments/migrate.py --pages /p/popular-1 /p/popular-2

# Stage 3: All remaining posts
python migrate_comments/migrate.py
```

## Next Steps

After successful migration:

1. **Update Hugo shortcodes** (Task 11)
   - Replace `{{< chat "id" >}}` with `{{< aicomments "id" >}}`
   - Use the shortcode replacer tool

2. **Verify on site**
   - Check comments display correctly
   - Test reply functionality
   - Verify timestamps are preserved

3. **Archive old system**
   - Keep Matrix data as backup
   - Document the migration
   - Update documentation

4. **Monitor new system**
   - Watch for any issues
   - Check comment submission works
   - Verify AI features function correctly
