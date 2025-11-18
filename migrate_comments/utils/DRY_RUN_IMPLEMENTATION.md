# Dry Run Mode Implementation

## Overview

The dry-run mode allows testing the migration process without actually importing data to the database. It extracts and transforms comments, validates all data, checks for duplicates, verifies parent-child relationships, and generates preview reports with go/no-go recommendations.

## Implementation Summary

### Task 9.1: Dry-Run Extraction ✓

**Implemented in:** `dry_run.py` - `DryRunExecutor.run_dry_run_for_page()`

- **Extract comments normally**: Uses `MatrixCommentExtractor.extract_comments_for_page()` to extract comments from Matrix
- **Transform comments normally**: Uses `CommentTransformer.transform_batch()` to transform comments to new system format
- **Skip actual import operations**: No calls to `CloudflareCommentClient` - comments are only logged
- **Log what would be imported**: Logs each comment that would be imported with DEBUG level logging

### Task 9.2: Dry-Run Validation ✓

**Implemented in:** `dry_run.py` - `DryRunExecutor.run_dry_run_for_page()`

- **Run all validation checks**: Uses `DataValidator.validate_migration_batch()` to validate both extracted and transformed comments
- **Report potential issues**: Collects validation errors and warnings in `DryRunPageResult.issues` and `DryRunPageResult.warnings`
- **Check for duplicate content**: Tracks content hashes across all pages to detect duplicates (implemented in `_check_for_duplicates()`)
- **Verify parent-child relationships**: Validates thread structure and parent references (implemented in `_verify_parent_child_relationships()`)

### Task 9.3: Generate Dry-Run Report ✓

**Implemented in:** `dry_run.py` - `DryRunExecutor.generate_preview_report()`

- **Show statistics of what would be migrated**: Displays total pages, comments extracted, comments that would be imported
- **List potential issues and warnings**: Shows all issues and warnings found during validation
- **Show sample transformed comments**: Includes sample comments from each page (configurable via `max_sample_comments`)
- **Provide go/no-go recommendation**: Generates recommendation based on success rate, error rate, and warning rate (implemented in `_generate_recommendation()`)

## Components

### DryRunPageResult

Dataclass that stores results for a single page:

```python
@dataclass
class DryRunPageResult:
    page_id: str
    section_id: str
    comments_extracted: int = 0
    comments_transformed: int = 0
    would_import: int = 0
    validation_passed: bool = True
    extraction_duration: float = 0.0
    transformation_duration: float = 0.0
    validation_duration: float = 0.0
    issues: List[str] = field(default_factory=list)
    warnings: List[str] = field(default_factory=list)
    sample_comments: List[Dict[str, Any]] = field(default_factory=list)
```

### DryRunSummary

Dataclass that stores overall dry-run results:

```python
@dataclass
class DryRunSummary:
    total_pages: int = 0
    successful_pages: int = 0
    pages_with_issues: int = 0
    total_comments_extracted: int = 0
    total_comments_would_import: int = 0
    total_issues: int = 0
    total_warnings: int = 0
    total_duration: float = 0.0
    recommendation: str = ""
    started_at: Optional[str] = None
    completed_at: Optional[str] = None
```

### DryRunExecutor

Main class that orchestrates the dry-run process:

**Key Methods:**

- `run_dry_run_for_page(section_id, page_id)`: Run dry-run for a single page
- `run_dry_run(section_ids)`: Run dry-run for multiple pages
- `generate_preview_report()`: Generate text preview report
- `export_to_json()`: Export results to JSON format
- `_check_for_duplicates()`: Check for duplicate content across pages
- `_verify_parent_child_relationships()`: Verify thread structure
- `_generate_recommendation()`: Generate go/no-go recommendation

## Usage

### Basic Usage

```python
from migrate_comments.extractors.matrix_extractor import MatrixCommentExtractor
from migrate_comments.transformers.comment_transformer import CommentTransformer
from migrate_comments.utils.data_validator import DataValidator
from migrate_comments.utils.page_id_mapper import PageIdMapper
from migrate_comments.utils.dry_run import DryRunExecutor

# Initialize components
extractor = MatrixCommentExtractor(
    homeserver_url="https://matrix.cactus.chat:8448",
    server_name="cactus.chat",
    site_name="rippreport.com"
)

page_mapper = PageIdMapper(hugo_content_dir="content")
page_mapper.scan_hugo_pages()

transformer = CommentTransformer(page_id_mapper=page_mapper)
validator = DataValidator()

# Create dry-run executor
dry_run = DryRunExecutor(
    extractor=extractor,
    transformer=transformer,
    validator=validator,
    page_id_mapper=page_mapper,
    output_dir="dry_run_reports"
)

# Run dry-run for all pages
summary = dry_run.run_dry_run()

# Generate reports
preview_report = dry_run.generate_preview_report()
json_file = dry_run.export_to_json()

# Check recommendation
print(summary.recommendation)
```

### Run Dry-Run for Specific Pages

```python
# Run dry-run for specific section IDs
section_ids = ["section1", "section2", "section3"]
summary = dry_run.run_dry_run(section_ids=section_ids)
```

### Demo Script

Run the demo script to see dry-run in action:

```bash
python migrate_comments/utils/demo_dry_run.py
```

This will:
1. Initialize all components
2. Run dry-run for the first 3 pages
3. Generate preview report and JSON export
4. Display summary and recommendation

## Report Formats

### Preview Report (Text)

The preview report includes:

1. **Summary Section**
   - Timeline (started, completed, duration)
   - Page statistics (total, successful, with issues)
   - Comment statistics (extracted, would import, issues, warnings)

2. **Recommendation Section**
   - GO/NO-GO/CAUTION recommendation with reasoning

3. **Pages with Issues Section**
   - Lists all pages that have issues or warnings
   - Shows specific issues and warnings for each page

4. **Sample Comments Section**
   - Shows sample transformed comments from first 5 pages
   - Includes author, timestamp, content preview, and reply status

5. **Per-Page Statistics Section**
   - Lists all pages with extraction and import counts
   - Shows validation status (✓ or ✗) for each page

### JSON Export

The JSON export includes:

```json
{
  "summary": {
    "total_pages": 10,
    "successful_pages": 9,
    "pages_with_issues": 1,
    "total_comments_extracted": 150,
    "total_comments_would_import": 148,
    "total_issues": 2,
    "total_warnings": 5,
    "total_duration": 45.3,
    "recommendation": "GO: Migration looks good...",
    "started_at": "2024-01-01T00:00:00Z",
    "completed_at": "2024-01-01T00:00:45Z"
  },
  "page_results": {
    "/page1": {
      "page_id": "/page1",
      "section_id": "section1",
      "comments_extracted": 15,
      "comments_transformed": 15,
      "would_import": 15,
      "validation_passed": true,
      "issues": [],
      "warnings": [],
      "sample_comments": [...]
    },
    ...
  }
}
```

## Recommendation Logic

The dry-run executor generates recommendations based on:

### GO Recommendation

- Success rate ≥ 95%
- Error rate < 5%
- Warning rate < 20%

Example: "GO: Migration looks good (100.0% success rate, 150 comments ready to import)."

### CAUTION Recommendation

- Success rate 80-95%
- Error rate 5-10%
- Warning rate 20-50%

Example: "CAUTION: Moderate error rate (7.5% of comments have issues). Review issues carefully before proceeding."

### NO-GO Recommendation

- Success rate < 80%
- Error rate > 10%

Example: "NO-GO: High error rate (12.3% of comments have issues). Review and fix issues before proceeding."

## Validation Checks

The dry-run performs comprehensive validation:

### Extraction Validation

- All required Matrix fields present
- Valid timestamp format
- Non-empty content
- Display name exists

### Transformation Validation

- All required new system fields present
- Valid ISO 8601 timestamp format
- Valid SHA-256 content hash format
- Valid page_id format

### Thread Validation

- Parent comments exist before replies
- No circular references
- Thread depth is reasonable (≤ 10 levels)
- Reports orphaned replies

### Duplicate Detection

- Tracks content hashes across all pages
- Reports potential duplicates
- Shows which pages have duplicate content

## Testing

Run the test suite:

```bash
python -m pytest migrate_comments/utils/test_dry_run.py -v
```

Test coverage includes:

- DryRunPageResult creation and serialization
- DryRunSummary creation and serialization
- DryRunExecutor initialization
- Successful dry-run for page
- Dry-run with no comments
- Dry-run with extraction errors
- Dry-run with validation failures
- Sample comment collection
- Recommendation generation (GO, CAUTION, NO-GO)
- Duplicate detection

## Integration with Main Migration Script

The dry-run mode will be integrated into the main migration script via command-line flag:

```bash
# Run dry-run
python migrate_comments/main.py --dry-run

# Run dry-run for specific pages
python migrate_comments/main.py --dry-run --pages section1,section2,section3

# Run actual migration after dry-run passes
python migrate_comments/main.py
```

## Files

- `migrate_comments/utils/dry_run.py` - Main implementation
- `migrate_comments/utils/demo_dry_run.py` - Demo script
- `migrate_comments/utils/test_dry_run.py` - Unit tests
- `migrate_comments/utils/DRY_RUN_IMPLEMENTATION.md` - This documentation

## Requirements Mapping

### Requirement 6.1: Dry-Run Extraction ✓

- WHERE dry run mode is enabled, THE Migration Tool SHALL extract and transform comments without importing them

### Requirement 6.2: Dry-Run Preview ✓

- WHERE dry run mode is enabled, THE Migration Tool SHALL generate a preview report showing what would be migrated

### Requirement 6.3: Dry-Run Validation ✓

- WHERE dry run mode is enabled, THE Migration Tool SHALL validate all data transformations

### Requirement 6.4: Dry-Run Issue Detection ✓

- WHERE dry run mode is enabled, THE Migration Tool SHALL identify potential issues without modifying the database

### Requirement 6.5: Dry-Run Statistics ✓

- WHEN dry run mode completes, THE Migration Tool SHALL provide statistics on comments that would be migrated

## Next Steps

To complete the dry-run implementation:

1. ✓ Implement dry-run extraction (Task 9.1)
2. ✓ Implement dry-run validation (Task 9.2)
3. ✓ Generate dry-run report (Task 9.3)
4. Integrate with main migration orchestrator (Task 10)
5. Add --dry-run command-line flag (Task 10.1)

## Notes

- Dry-run mode does NOT make any database changes
- All extraction and transformation operations are performed normally
- Only the import step is skipped
- Reports are saved to the `dry_run_reports/` directory (configurable)
- Dry-run can be run multiple times without side effects
- Use dry-run to test migration before committing to actual import
