# ProgressTracker Implementation

## Overview

The `ProgressTracker` class provides robust progress tracking for the archive link verification and replacement process. It enables resumability, prevents duplicate work, and provides comprehensive statistics.

## Features

✅ **Progress Persistence**: Saves progress to disk after each link is processed
✅ **File Locking**: Uses `fcntl` to prevent concurrent modifications
✅ **Duplicate Detection**: Checks both post_url and original_href to identify duplicates
✅ **Resumability**: Can resume from any point after interruption
✅ **Statistics**: Tracks success rates and various status types
✅ **Error Recovery**: Handles corrupted progress files gracefully
✅ **Summary Generation**: Provides formatted progress summaries

## Implementation Details

### Class Structure

```python
class ProgressTracker:
    def __init__(self, progress_file: str = "archive_verifier_progress.json")
    def load(self) -> Dict
    def save(self, lock_timeout: int = 30, retry_interval: int = 1)
    def is_processed(self, post_url: str, original_href: str) -> bool
    def get_processed_entry(self, post_url: str, original_href: str) -> Optional[Dict]
    def mark_processed(self, post_url: str, original_href: str, status: str, details: Dict)
    def get_summary(self) -> Dict
    def get_remaining_count(self, total_links: int) -> int
    def print_summary()
```

### Progress File Format

```json
{
  "last_updated": "2025-10-19T12:34:56",
  "processed_links": [
    {
      "post_url": "http://localhost:1313/p/test/",
      "original_href": "https://example.com/page",
      "status": "replaced",
      "timestamp": "2025-10-19T12:30:00",
      "archive_url": "https://web.archive.org/web/20200101/...",
      "backup_path": "backups/archive-links-backup-20251019-123000/test/index.md"
    }
  ],
  "statistics": {
    "total_processed": 150,
    "successfully_replaced": 120,
    "verification_failed": 20,
    "replacement_failed": 10,
    "skipped": 0
  }
}
```

### Status Types

- **`replaced`**: Successfully replaced with verified archive
- **`verification_failed`**: Archive verification failed
- **`replacement_failed`**: File replacement failed
- **`skipped`**: Skipped (e.g., already manually fixed)
- **`failed`**: General failure

## Usage Examples

### Basic Usage

```python
from archive_link_verifier.progress_tracker import ProgressTracker

# Initialize tracker
tracker = ProgressTracker("archive_verifier_progress.json")

# Load existing progress
tracker.load()

# Mark a link as processed
tracker.mark_processed(
    post_url="http://localhost:1313/p/test-post/",
    original_href="https://example.com/broken-link",
    status="replaced",
    details={
        "archive_url": "https://web.archive.org/web/20200101/example.com/page",
        "verification_method": "playwright+lmstudio",
        "backup_path": "backups/archive-links-backup-20251019-123000/test-post/index.md"
    }
)

# Save progress
tracker.save()

# Print summary
tracker.print_summary()
```

### Duplicate Detection

```python
# Check if link already processed
if tracker.is_processed(post_url, original_href):
    print("Already processed, skipping...")
else:
    # Process the link
    process_link(post_url, original_href)
    tracker.mark_processed(post_url, original_href, "replaced", details)
```

### Resumability

```python
# First run
tracker = ProgressTracker()
tracker.load()  # Load existing progress

for entry in report_entries:
    if tracker.is_processed(entry["post_url"], entry["original_href"]):
        continue  # Skip already processed
    
    # Process link
    result = process_link(entry)
    tracker.mark_processed(
        entry["post_url"],
        entry["original_href"],
        result["status"],
        result["details"]
    )
    tracker.save()  # Save after each link

# If interrupted, next run will resume from where it left off
```

### Statistics

```python
# Get summary statistics
summary = tracker.get_summary()
print(f"Total processed: {summary['total_processed']}")
print(f"Success rate: {summary['success_rate']:.1f}%")

# Calculate remaining work
total_links = len(report_entries)
remaining = tracker.get_remaining_count(total_links)
print(f"Remaining: {remaining} links")
```

## File Locking

The `save()` method uses `fcntl.flock()` to prevent concurrent modifications:

```python
# Acquire exclusive lock (non-blocking)
fcntl.flock(f.fileno(), fcntl.LOCK_EX | fcntl.LOCK_NB)

# Write data
json.dump(progress_data, f)

# Release lock (automatic on file close)
fcntl.flock(f.fileno(), fcntl.LOCK_UN)
```

If the lock cannot be acquired, the method retries with a configurable timeout:

```python
tracker.save(lock_timeout=30, retry_interval=1)
```

## Error Handling

### Corrupted Progress File

If the progress file is corrupted (invalid JSON), the tracker:

1. Backs up the corrupted file with `.corrupted` suffix
2. Starts with fresh progress data
3. Logs a warning

```python
# Corrupted file: progress.json
# Backed up to: progress.json.corrupted
# If backup exists: progress.json.corrupted.1, progress.json.corrupted.2, etc.
```

### Missing Progress File

If no progress file exists, the tracker starts with empty progress:

```python
{
  "last_updated": None,
  "processed_links": [],
  "statistics": {
    "total_processed": 0,
    "successfully_replaced": 0,
    "verification_failed": 0,
    "replacement_failed": 0,
    "skipped": 0
  }
}
```

### File Lock Timeout

If the file lock cannot be acquired within the timeout, a `TimeoutError` is raised:

```python
try:
    tracker.save(lock_timeout=30)
except TimeoutError as e:
    print(f"Could not save progress: {e}")
    # Handle error (e.g., retry later, alert user)
```

## Testing

Comprehensive test suite with 19 tests covering:

- ✅ Initialization
- ✅ Loading (existing, missing, corrupted files)
- ✅ Saving with file locking
- ✅ Duplicate detection
- ✅ Entry retrieval
- ✅ Status tracking
- ✅ Statistics calculation
- ✅ Summary generation
- ✅ Remaining count calculation

Run tests:

```bash
python -m pytest archive_link_verifier/test_progress_tracker.py -v
```

## Demo Script

Run the demo to see all features in action:

```bash
python archive_link_verifier/demo_progress_tracker.py
```

The demo shows:
1. Basic progress tracking
2. Duplicate detection
3. Resumability after interruption
4. Statistics and summary generation
5. Corrupted file recovery

## Integration with Main Script

The ProgressTracker will be integrated into the main orchestration script:

```python
# main.py
from archive_link_verifier.progress_tracker import ProgressTracker

def main():
    # Initialize tracker
    tracker = ProgressTracker("archive_verifier_progress.json")
    tracker.load()
    
    # Load report
    report = load_report("archive_fix_report.json")
    
    # Show resume info
    total = len(report["entries"])
    remaining = tracker.get_remaining_count(total)
    if remaining < total:
        print(f"Resuming: {total - remaining} already processed, {remaining} remaining")
    
    # Process each entry
    for entry in report["entries"]:
        # Skip if already processed
        if tracker.is_processed(entry["post_url"], entry["original_href"]):
            continue
        
        # Process link
        result = process_link(entry)
        
        # Track progress
        tracker.mark_processed(
            entry["post_url"],
            entry["original_href"],
            result["status"],
            result["details"]
        )
        
        # Save after each link (resumability)
        tracker.save()
    
    # Final summary
    tracker.print_summary()
```

## Requirements Satisfied

This implementation satisfies **Requirement 6** from the requirements document:

✅ **6.1**: Saves progress to disk after each link is processed
✅ **6.2**: Saves current state on interruption (via save after each link)
✅ **6.3**: Reads existing report and skips processed links on restart
✅ **6.4**: Checks both post_url and original_href for duplicates
✅ **6.5**: Can resume from replacement step (via status tracking)
✅ **6.6**: Logs errors, marks in report, continues with remaining links
✅ **6.7**: Generates summary report with statistics
✅ **6.8**: Displays resume message with processed/remaining counts
✅ **6.9**: Creates backup of corrupted file and starts fresh
✅ **6.10**: Uses file locking to prevent concurrent modifications

## Performance Considerations

- **Save Frequency**: Saves after each link (configurable if needed)
- **File Locking**: Non-blocking with retry mechanism
- **Memory Usage**: Loads entire progress file into memory (acceptable for typical use)
- **Disk I/O**: Minimal overhead with JSON serialization

## Future Enhancements

Potential improvements for future versions:

1. **Batch Saves**: Option to save every N links instead of every link
2. **Compression**: Compress progress file for large datasets
3. **Database Backend**: Use SQLite for very large progress tracking
4. **Progress Callbacks**: Notify external systems of progress updates
5. **Detailed Timing**: Track time spent on each link for performance analysis

## Conclusion

The ProgressTracker provides a robust, production-ready solution for tracking progress during archive link verification and replacement. It handles edge cases gracefully and ensures no work is lost due to interruptions.
