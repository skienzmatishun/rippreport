# Migration Reporter Quick Start Guide

## Basic Usage

### 1. Initialize Reporter

```python
from migration_reporter import MigrationReporter
import logging

# Set up logger
logger = logging.getLogger('migration')

# Create reporter
reporter = MigrationReporter(
    output_dir="migration_reports",
    logger=logger
)

# Start migration tracking
reporter.start_migration(total_pages=100)
```

### 2. Log Operations

```python
# Log extraction
reporter.log_extraction(
    page_id="/p/my-article/",
    section_id="my-article",
    comment_count=5,
    duration=0.5,
    success=True
)

# Log transformation
reporter.log_transformation(
    page_id="/p/my-article/",
    input_count=5,
    output_count=5,
    duration=0.1,
    success=True
)

# Log successful import
reporter.log_import(
    page_id="/p/my-article/",
    comment_id=1001,
    duration=0.05,
    success=True,
    matrix_event_id="$matrix_event_123"
)

# Log failed import
reporter.log_import(
    page_id="/p/my-article/",
    comment_id=None,
    duration=0.05,
    success=False,
    error="API rate limit exceeded",
    matrix_event_id="$matrix_event_124"
)

# Mark page complete
reporter.log_page_complete("/p/my-article/", success=True)
```

### 3. Generate Reports

```python
# End migration
reporter.end_migration()

# Generate all reports at once
reports = reporter.generate_all_reports()

# Or generate individual reports
summary = reporter.generate_summary_report()
detailed = reporter.generate_detailed_report()
id_mapping = reporter.generate_id_mapping_report()
json_file = reporter.export_to_json()
failed_file = reporter.export_failed_comments()
```

## Common Patterns

### Pattern 1: Complete Page Migration

```python
import time

page_id = "/p/article/"
section_id = "article"

# Start timing
start_time = time.time()

# Extract
try:
    comments = extractor.extract_comments_for_page(section_id)
    extraction_time = time.time() - start_time
    
    reporter.log_extraction(
        page_id=page_id,
        section_id=section_id,
        comment_count=len(comments),
        duration=extraction_time,
        success=True
    )
except Exception as e:
    extraction_time = time.time() - start_time
    reporter.log_extraction(
        page_id=page_id,
        section_id=section_id,
        comment_count=0,
        duration=extraction_time,
        success=False,
        error=str(e)
    )
    reporter.log_page_complete(page_id, success=False)
    return

# Transform
start_time = time.time()
transformed = transformer.transform_batch(comments, page_id)
transformation_time = time.time() - start_time

reporter.log_transformation(
    page_id=page_id,
    input_count=len(comments),
    output_count=len(transformed),
    duration=transformation_time,
    success=True
)

# Import
for comment in transformed:
    start_time = time.time()
    try:
        result = client.create_comment(comment)
        import_time = time.time() - start_time
        
        reporter.log_import(
            page_id=page_id,
            comment_id=result['id'],
            duration=import_time,
            success=True,
            matrix_event_id=comment.get('matrix_event_id')
        )
    except Exception as e:
        import_time = time.time() - start_time
        reporter.log_import(
            page_id=page_id,
            comment_id=None,
            duration=import_time,
            success=False,
            error=str(e),
            matrix_event_id=comment.get('matrix_event_id')
        )

# Mark complete
reporter.log_page_complete(page_id, success=True)
```

### Pattern 2: Using with CheckpointManager

```python
from checkpoint_manager import CheckpointManager
from migration_reporter import MigrationReporter

# Initialize both
checkpoint = CheckpointManager()
reporter = MigrationReporter()

# Start migration
reporter.start_migration(total_pages=100)

# After migration, sync ID mappings
reporter.set_id_mapping(checkpoint.id_mapping)

# Generate reports
reporter.end_migration()
reports = reporter.generate_all_reports()
```

### Pattern 3: Filtering Operation Logs

```python
# Get all failed imports
failed_imports = reporter.get_operation_logs(
    operation_type='import',
    success=False
)

# Get all operations for a specific page
page_ops = reporter.get_operation_logs(
    page_id='/p/my-article/'
)

# Get all failed extractions
failed_extractions = reporter.get_operation_logs(
    operation_type='extraction',
    success=False
)
```

### Pattern 4: Custom Report Analysis

```python
# Access raw data
summary = reporter.summary
page_reports = reporter.page_reports
operation_logs = reporter.operation_logs
id_mapping = reporter.id_mapping

# Calculate custom metrics
avg_extraction_time = sum(
    log.duration for log in reporter.get_operation_logs(operation_type='extraction')
) / len(reporter.get_operation_logs(operation_type='extraction'))

print(f"Average extraction time: {avg_extraction_time:.2f}s")

# Find slowest pages
slowest_pages = sorted(
    reporter.page_reports.items(),
    key=lambda x: x[1].total_duration,
    reverse=True
)[:10]

for page_id, report in slowest_pages:
    print(f"{page_id}: {report.total_duration:.2f}s")
```

## Report Output Locations

All reports are saved to the `output_dir` (default: `migration_reports/`):

- `summary_report_YYYYMMDD_HHMMSS.txt` - Summary report
- `detailed_report_YYYYMMDD_HHMMSS.txt` - Detailed report
- `id_mapping_YYYYMMDD_HHMMSS.txt` - ID mapping table
- `migration_report_YYYYMMDD_HHMMSS.json` - Complete JSON export
- `failed_comments_YYYYMMDD_HHMMSS.json` - Failed comments (if any)

## Tips

1. **Always call `start_migration()` first** - This initializes timing and statistics

2. **Always call `end_migration()` before generating reports** - This finalizes timing

3. **Log operations in order** - Extraction → Transformation → Import(s) → Page Complete

4. **Use consistent page IDs** - Use the same page ID across all operations for a page

5. **Include Matrix event IDs** - This enables ID mapping for parent-child relationships

6. **Handle errors gracefully** - Always log failures with error messages

7. **Generate reports at the end** - Wait until migration is complete for accurate statistics

## Troubleshooting

### Reports are empty
- Make sure you called `start_migration()` and `end_migration()`
- Check that you logged operations between start and end

### ID mapping is missing
- Ensure you pass `matrix_event_id` to `log_import()`
- Check that imports were successful (ID mapping only stored on success)

### Statistics don't match
- Verify you're calling `log_page_complete()` for each page
- Check that you're using consistent page IDs

### Reports not saved
- Check that `output_dir` exists and is writable
- Verify no exceptions during report generation

## Examples

See the demo script for complete examples:
```bash
python migrate_comments/utils/demo_migration_reporter.py
```

See the test suite for usage patterns:
```bash
python migrate_comments/utils/test_migration_reporter.py
```
