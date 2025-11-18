# Troubleshooting Guide

This guide helps you diagnose and fix common issues with the Comment Migration Tool.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Configuration Issues](#configuration-issues)
- [Matrix Extraction Issues](#matrix-extraction-issues)
- [Transformation Issues](#transformation-issues)
- [Import Issues](#import-issues)
- [Recovery Issues](#recovery-issues)
- [Performance Issues](#performance-issues)
- [Data Validation Issues](#data-validation-issues)
- [Debugging Tips](#debugging-tips)
- [FAQ](#faq)

---

## Installation Issues

### Python Version Problems

**Problem**: `python: command not found`

**Cause**: Python is not installed or not in PATH

**Solution**:
```bash
# Check if Python 3 is available
python3 --version

# If not installed:
# macOS: brew install python@3.9
# Ubuntu: sudo apt install python3.9
# Windows: Download from python.org
```

---

**Problem**: `SyntaxError` or `invalid syntax` errors

**Cause**: Using Python 2.x instead of Python 3.x

**Solution**:
```bash
# Always use python3 explicitly
python3 migrate.py

# Or create an alias
alias python=python3
```

---

### Virtual Environment Issues

**Problem**: `No module named 'venv'`

**Cause**: venv module not installed

**Solution**:
```bash
# Ubuntu/Debian
sudo apt install python3-venv

# Or use virtualenv instead
pip3 install virtualenv
python3 -m virtualenv venv
```

---

**Problem**: Virtual environment not activating

**Cause**: Wrong activation command for your shell

**Solution**:
```bash
# Bash/Zsh (macOS/Linux)
source venv/bin/activate

# Fish shell
source venv/bin/activate.fish

# Windows Command Prompt
venv\Scripts\activate.bat

# Windows PowerShell
venv\Scripts\Activate.ps1
```

---

### Dependency Installation Issues

**Problem**: `pip install` fails with SSL errors

**Cause**: SSL certificate verification issues

**Solution**:
```bash
pip install --trusted-host pypi.org --trusted-host files.pythonhosted.org -r requirements.txt
```

---

**Problem**: `pip install` fails with permission errors

**Cause**: Trying to install system-wide without sudo

**Solution**:
```bash
# Ensure virtual environment is activated
source venv/bin/activate

# Verify you're in venv
which python  # Should show path to venv/bin/python

# Then install
pip install -r requirements.txt
```

---

**Problem**: `ModuleNotFoundError` when running script

**Cause**: Dependencies not installed or wrong Python interpreter

**Solution**:
```bash
# Ensure venv is activated
source venv/bin/activate

# Reinstall dependencies
pip install -r requirements.txt

# Verify installation
pip list
```

---

### Wrangler Issues

**Problem**: `wrangler: command not found`

**Cause**: Wrangler not installed or not in PATH

**Solution**:
```bash
# Install wrangler
npm install -g wrangler

# Verify installation
wrangler --version

# If still not found, add npm global bin to PATH
export PATH="$(npm config get prefix)/bin:$PATH"
```

---

**Problem**: `wrangler login` fails

**Cause**: Browser not opening or authentication issues

**Solution**:
```bash
# Try manual authentication
wrangler login

# If browser doesn't open, copy the URL and open manually

# Verify authentication
wrangler whoami
```

---

**Problem**: `Error: Not authenticated`

**Cause**: Wrangler not logged in

**Solution**:
```bash
# Login to Cloudflare
wrangler login

# Verify
wrangler whoami
```

---

## Configuration Issues

### Configuration File Not Found

**Problem**: `FileNotFoundError: config.yaml not found`

**Cause**: Configuration file doesn't exist or wrong path

**Solution**:
```bash
# Create from example
cp config.yaml.example config.yaml

# Verify file exists
ls -la config.yaml

# If running from wrong directory
cd migrate_comments
python migrate.py
```

---

### YAML Syntax Errors

**Problem**: `yaml.scanner.ScannerError` or `yaml.parser.ParserError`

**Cause**: Invalid YAML syntax

**Solution**:
1. Check indentation (must be spaces, not tabs)
2. Ensure strings with special characters are quoted
3. Validate YAML online: https://www.yamllint.com/

Common issues:
```yaml
# WRONG - tabs instead of spaces
matrix:
	site_name: "example.com"

# CORRECT - spaces for indentation
matrix:
  site_name: "example.com"

# WRONG - unquoted URL with special characters
api_base: https://api.example.com/path?key=value

# CORRECT - quoted URL
api_base: "https://api.example.com/path?key=value"
```

---

### Missing Required Fields

**Problem**: `Configuration error: matrix.site_name is required`

**Cause**: Required configuration field is missing

**Solution**:
```yaml
# Add the required field to config.yaml
matrix:
  site_name: "your-site.com"  # Add this

cloudflare:
  api_base: "https://comments.your-site.com/api"  # Add this
  database_id: "your-database-id"  # Add this

hugo:
  base_url: "https://your-site.com"  # Add this
```

---

### Invalid Configuration Values

**Problem**: `Configuration error: cloudflare.api_base must be a valid URL`

**Cause**: Invalid URL format

**Solution**:
```yaml
# WRONG
cloudflare:
  api_base: comments.example.com

# CORRECT
cloudflare:
  api_base: "https://comments.example.com/api"
```

---

### Environment Variable Issues

**Problem**: Environment variables not being recognized

**Cause**: Variables not exported or wrong naming

**Solution**:
```bash
# Ensure variables are exported
export MATRIX_SITE_NAME="example.com"

# Verify they're set
echo $MATRIX_SITE_NAME

# Check exact variable names (case-sensitive)
# Correct: MATRIX_SITE_NAME
# Wrong: matrix_site_name
```

---

## Matrix Extraction Issues

### Connection Errors

**Problem**: `ConnectionError: Failed to connect to Matrix homeserver`

**Cause**: Network issues or wrong homeserver URL

**Solution**:
```bash
# Test connection manually
curl https://matrix.cactus.chat:8448/_matrix/client/versions

# If that fails, check:
# 1. Network connectivity
# 2. Firewall settings
# 3. Homeserver URL in config.yaml
```

---

**Problem**: `Timeout error` when connecting to Matrix

**Cause**: Slow network or homeserver issues

**Solution**:
1. Check your internet connection
2. Try again later (homeserver might be busy)
3. Increase timeout in code if persistent

---

### Room Not Found

**Problem**: `Matrix room not found for section ID: xyz`

**Cause**: Room doesn't exist or wrong site name

**Solution**:
```bash
# Verify site name matches Cactus Chat setup
# Room ID format: #comments_<site_name>_<section_id>:cactus.chat

# Check config.yaml
matrix:
  site_name: "example.com"  # Must match exactly

# Test room ID manually
# Visit: https://matrix.to/#/#comments_example.com_xyz:cactus.chat
```

---

**Problem**: `No comments found for page`

**Cause**: Page has no comments or wrong section ID

**Solution**:
1. Verify the page actually has comments in Cactus Chat
2. Check the section ID in the Hugo shortcode
3. Ensure the Matrix room exists

---

### Authentication Errors

**Problem**: `403 Forbidden` when accessing Matrix

**Cause**: Matrix API requires authentication for this room

**Solution**:
- Cactus Chat rooms are usually public
- If private, you'll need to add authentication to the extractor
- Contact Cactus Chat support if rooms should be public

---

### Data Parsing Errors

**Problem**: `KeyError` or `AttributeError` when parsing Matrix events

**Cause**: Unexpected Matrix event format

**Solution**:
```bash
# Run with debug logging
python migrate.py --verbose

# Check the log for the problematic event
grep "ERROR" migration.log

# Report the issue with the event data
```

---

## Transformation Issues

### Page ID Mapping Errors

**Problem**: `Page not found for section ID: xyz`

**Cause**: Hugo page doesn't exist or shortcode missing

**Solution**:
```bash
# Verify the page exists
ls content/p/xyz/index.md

# Check for chat shortcode in the file
grep "chat" content/p/xyz/index.md

# Ensure frontmatter has permalink
grep "permalink" content/p/xyz/index.md
```

---

**Problem**: `Multiple pages found for section ID`

**Cause**: Duplicate section IDs in Hugo content

**Solution**:
1. Search for duplicate shortcodes:
   ```bash
   grep -r "chat \"xyz\"" content/
   ```
2. Ensure each section ID is unique
3. Update duplicate shortcodes with unique IDs

---

### Timestamp Conversion Errors

**Problem**: `ValueError: Invalid timestamp format`

**Cause**: Unexpected timestamp format from Matrix

**Solution**:
```bash
# Run with debug logging to see the problematic timestamp
python migrate.py --verbose

# Check migration.log for the error
grep "timestamp" migration.log
```

---

### Content Hash Errors

**Problem**: `Duplicate content hash detected`

**Cause**: Identical comment content (might be spam or duplicate)

**Solution**:
- This is usually intentional (prevents duplicates)
- Check if the comment is actually a duplicate
- If legitimate, the tool will log a warning and skip

---

### Parent ID Resolution Errors

**Problem**: `Parent comment not found for reply`

**Cause**: Parent comment failed to import or wrong event ID

**Solution**:
1. Check if parent comment import failed
2. Review the detailed report for parent comment errors
3. Use `--resume` to retry failed imports
4. Orphaned replies will be imported as root comments

---

## Import Issues

### Wrangler Execution Errors

**Problem**: `Error executing wrangler command`

**Cause**: Wrangler not found, not authenticated, or wrong database ID

**Solution**:
```bash
# Verify wrangler works
wrangler --version

# Verify authentication
wrangler whoami

# Test database access
wrangler d1 execute <database-id> --command "SELECT COUNT(*) FROM comments"

# Check database ID in config.yaml
cloudflare:
  database_id: "correct-database-id"  # Get from: wrangler d1 list
```

---

**Problem**: `Database not found`

**Cause**: Wrong database ID or database doesn't exist

**Solution**:
```bash
# List your databases
wrangler d1 list

# Copy the correct database ID to config.yaml
cloudflare:
  database_id: "4c94fdd6-f883-439a-944c-a63a5cffac9c"
```

---

### API Rate Limiting

**Problem**: `429 Too Many Requests` errors

**Cause**: Hitting Cloudflare API rate limits

**Solution**:
```yaml
# Increase delay between batches
migration:
  delay_between_batches: 2.0  # Increase from 1.0

# Reduce batch size
migration:
  batch_size: 5  # Reduce from 10
```

The tool automatically retries with exponential backoff, but adjusting these settings helps prevent rate limiting.

---

**Problem**: Repeated rate limit errors despite retries

**Cause**: Too aggressive batch settings

**Solution**:
```yaml
# Use very conservative settings
migration:
  batch_size: 1
  delay_between_batches: 3.0
  max_retries: 5
```

---

### SQL Errors

**Problem**: `SQL error: UNIQUE constraint failed`

**Cause**: Trying to import duplicate comments

**Solution**:
- This is usually prevented by content hash checking
- If it occurs, check for duplicate imports
- Use `--resume` to skip already-imported comments

---

**Problem**: `SQL error: FOREIGN KEY constraint failed`

**Cause**: Parent comment doesn't exist for a reply

**Solution**:
- The tool should handle this automatically
- Check if parent comment import failed
- Review detailed report for parent comment errors

---

### Network Errors

**Problem**: `ConnectionError` or `Timeout` during import

**Cause**: Network issues or API unavailable

**Solution**:
1. Check internet connection
2. Verify API endpoint is accessible:
   ```bash
   curl https://comments.your-site.com/api/health
   ```
3. Use `--resume` to continue after network is restored

---

## Recovery Issues

### Checkpoint File Corrupted

**Problem**: `JSONDecodeError` when loading checkpoint

**Cause**: Checkpoint file is corrupted or incomplete

**Solution**:
```bash
# Remove corrupted checkpoint
rm .migration_checkpoint.json

# Start migration fresh
python migrate.py
```

**Note**: You'll lose progress, but the tool will skip already-imported comments.

---

### Resume Not Working

**Problem**: `--resume` flag doesn't skip completed pages

**Cause**: Checkpoint file missing or not being loaded

**Solution**:
```bash
# Verify checkpoint file exists
ls -la .migration_checkpoint.json

# Check checkpoint file is valid JSON
python -c "import json; json.load(open('.migration_checkpoint.json'))"

# If invalid, remove and start fresh
rm .migration_checkpoint.json
python migrate.py
```

---

### ID Mapping Lost

**Problem**: Parent-child relationships broken after resume

**Cause**: ID mapping not persisted in checkpoint

**Solution**:
- This shouldn't happen with the current implementation
- If it does, you may need to re-import affected pages
- Use `--pages` to re-import specific pages

---

## Performance Issues

### Migration Too Slow

**Problem**: Migration taking too long

**Cause**: Conservative batch settings or slow network

**Solution**:
```yaml
# Increase batch size
migration:
  batch_size: 20  # From 10

# Reduce delay
migration:
  delay_between_batches: 0.5  # From 1.0
```

**Warning**: May increase rate limiting risk

---

**Problem**: High memory usage

**Cause**: Processing too many comments at once

**Solution**:
```yaml
# Reduce batch size
migration:
  batch_size: 5  # From 10
```

---

### Disk Space Issues

**Problem**: `No space left on device`

**Cause**: Log files or reports filling disk

**Solution**:
```bash
# Check disk space
df -h

# Remove old log files
rm migration.log.*

# Remove old reports
rm -rf migration_reports/*

# Reduce log level
# In config.yaml:
logging:
  level: "WARNING"  # From "INFO" or "DEBUG"
```

---

## Data Validation Issues

### Required Field Missing

**Problem**: `Validation error: Required field 'content' is missing`

**Cause**: Matrix comment has no content

**Solution**:
- This indicates a problem with the source data
- Check the Matrix event in the log
- The comment will be skipped (logged as failed)

---

### Invalid Timestamp

**Problem**: `Validation error: Invalid timestamp format`

**Cause**: Matrix timestamp is not in expected format

**Solution**:
```bash
# Run with debug logging
python migrate.py --verbose

# Check the problematic timestamp in the log
grep "timestamp" migration.log

# Report the issue with the timestamp value
```

---

### Thread Validation Errors

**Problem**: `Validation error: Circular reference detected`

**Cause**: Comment thread has circular parent-child relationship

**Solution**:
- This indicates corrupted data in Matrix
- The tool will break the circular reference
- Check the detailed report for affected comments

---

**Problem**: `Validation error: Thread depth exceeds maximum`

**Cause**: Comment thread is too deeply nested

**Solution**:
- This is rare but possible with very long threads
- The tool will flatten excessively deep threads
- Check the detailed report for affected comments

---

## Debugging Tips

### Enable Debug Logging

```bash
# Method 1: Command line flag
python migrate.py --verbose

# Method 2: Config file
# In config.yaml:
logging:
  level: "DEBUG"
```

### Check Log Files

```bash
# View entire log
cat migration.log

# View recent entries
tail -100 migration.log

# Follow log in real-time
tail -f migration.log

# Search for errors
grep "ERROR" migration.log

# Search for specific page
grep "page-id" migration.log
```

### Test Individual Components

```bash
# Test Matrix extraction
python -c "from extractors.matrix_extractor import MatrixCommentExtractor; print('OK')"

# Test configuration
python -c "from utils.config import Config; c = Config('config.yaml'); print('OK')"

# Test page ID mapper
python -c "from utils.page_id_mapper import PageIdMapper; print('OK')"
```

### Isolate Problems

```bash
# Test with one page
python migrate.py --pages /p/test-page --verbose

# Test with dry run
python migrate.py --dry-run --verbose

# Test specific component
python utils/demo_matrix_extractor.py
```

### Check API Connectivity

```bash
# Test Matrix API
curl https://matrix.cactus.chat:8448/_matrix/client/versions

# Test Cloudflare API
curl https://comments.your-site.com/api/health

# Test wrangler
wrangler d1 execute <database-id> --command "SELECT 1"
```

### Verify Data

```bash
# Check comment count in database
wrangler d1 execute <database-id> --command "SELECT COUNT(*) FROM comments"

# Check specific page comments
wrangler d1 execute <database-id> --command "SELECT * FROM comments WHERE page_id='/p/test-page'"

# Check for orphaned replies
wrangler d1 execute <database-id> --command "SELECT * FROM comments WHERE parent_id IS NOT NULL AND parent_id NOT IN (SELECT id FROM comments)"
```

---

## FAQ

### Q: Can I run the migration multiple times?

**A**: Yes, but be careful:
- Dry-run mode can be run unlimited times (no data changes)
- Actual migration uses content hashes to prevent duplicates
- Already-imported comments will be skipped
- Use `--pages` to re-import specific pages if needed

---

### Q: What happens if migration is interrupted?

**A**: The tool saves checkpoints:
- Progress is saved after each page
- Use `--resume` to continue from where you left off
- ID mappings are preserved for parent-child relationships
- No data corruption (each page is atomic)

---

### Q: How do I handle failed imports?

**A**: Several options:
1. Use `--resume` to retry failed pages
2. Check `failed_comments_*.json` for details
3. Fix issues and re-run for specific pages with `--pages`
4. Some failures may require manual intervention

---

### Q: Can I migrate only certain pages?

**A**: Yes:
```bash
python migrate.py --pages /p/page-1 /p/page-2 /p/page-3
```

---

### Q: How long does migration take?

**A**: Depends on:
- Number of pages and comments
- Network speed
- Batch size and delay settings
- API rate limits

Estimate: ~1-2 seconds per comment with default settings

---

### Q: Will timestamps be preserved?

**A**: Yes:
- Original Matrix timestamps are preserved
- Comments appear with their original creation time
- Uses direct D1 database access to set `created_at`

---

### Q: What about comment IDs?

**A**: 
- New comment IDs are generated by the database
- Matrix event IDs are mapped to new IDs
- Mapping is saved in reports and checkpoint
- Parent-child relationships are preserved

---

### Q: Can I rollback the migration?

**A**: Yes, if you have backups:
1. Restore D1 database from backup
2. Restore Hugo content files from backup
3. Clear any partial migrations

**Important**: Make backups before migrating!

---

### Q: What if I find bugs?

**A**: 
1. Run with `--verbose` to get detailed logs
2. Check `migration.log` for error details
3. Review the implementation summaries in the docs
4. Report issues with log excerpts and error messages

---

### Q: Can I customize the migration?

**A**: Yes:
- Edit `config.yaml` for settings
- Modify Python code for custom behavior
- Use environment variables for different environments
- See the design document for architecture details

---

### Q: How do I verify migration success?

**A**: Multiple ways:
1. Check summary report for success rate
2. Spot-check comments on your site
3. Compare comment counts:
   ```bash
   # Count in database
   wrangler d1 execute <db-id> --command "SELECT COUNT(*) FROM comments"
   ```
4. Verify parent-child relationships work
5. Check that timestamps are correct

---

### Q: What about duplicate comments?

**A**: 
- Tool uses content hashes to detect duplicates
- Duplicates are automatically skipped
- Logged as warnings in the report
- No duplicate comments will be imported

---

### Q: Can I migrate from multiple Matrix servers?

**A**: 
- Current implementation supports one Matrix server
- To migrate from multiple servers:
  1. Run migration for each server separately
  2. Use different config files
  3. Merge results manually if needed

---

### Q: What if my Hugo structure is different?

**A**: 
- Tool assumes standard Hugo structure
- If different, you may need to:
  1. Adjust `hugo.content_dir` in config
  2. Modify `PageIdMapper` for custom structure
  3. Ensure permalinks are correctly extracted

---

### Q: How do I test without affecting production?

**A**: 
1. Use `--dry-run` mode (no data changes)
2. Test with a staging environment
3. Use `--pages` to test with a few pages first
4. Make database backups before migrating

---

## Getting More Help

If this guide doesn't solve your problem:

1. **Check the documentation**:
   - [INSTALLATION.md](INSTALLATION.md) - Installation issues
   - [CONFIGURATION.md](CONFIGURATION.md) - Configuration issues
   - [USAGE.md](USAGE.md) - Usage questions
   - Design document - Architecture questions

2. **Review implementation summaries**:
   - Each component has a summary document
   - Check `*_IMPLEMENTATION.md` files
   - Review verification checklists

3. **Enable debug logging**:
   ```bash
   python migrate.py --verbose
   ```

4. **Check the logs**:
   ```bash
   tail -100 migration.log
   ```

5. **Test components individually**:
   ```bash
   python utils/demo_*.py
   ```

6. **Report issues**:
   - Include error messages
   - Include relevant log excerpts
   - Include configuration (redact sensitive values)
   - Include steps to reproduce
