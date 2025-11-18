# Data Validator Quick Start Guide

## Quick Usage

### 1. Basic Validation

```python
from migrate_comments.utils.data_validator import DataValidator

# Create validator
validator = DataValidator()

# Validate extracted comments
extraction_summary = validator.validate_extracted_comments(matrix_comments)
print(f"Valid: {extraction_summary.valid_items}/{extraction_summary.total_items}")

# Validate transformed comments
transformation_summary = validator.validate_transformed_comments(transformed_comments)
print(f"Valid: {transformation_summary.valid_items}/{transformation_summary.total_items}")

# Validate thread structure
thread_result = validator.validate_thread_structure(matrix_comments)
if not thread_result.is_valid:
    print(f"Thread errors: {thread_result.errors}")
```

### 2. Comprehensive Validation

```python
# Validate entire migration batch
results = validator.validate_migration_batch(
    extracted_comments=matrix_comments,
    transformed_comments=transformed_comments
)

if results['overall_valid']:
    print("✓ All validation passed - ready to import")
else:
    print("✗ Validation failed - review issues before importing")
```

### 3. Filter Invalid Comments

```python
# Get only valid comments
valid_extracted = validator.filter_valid_comments(
    matrix_comments,
    extraction_summary
)

valid_transformed = validator.filter_valid_comments(
    transformed_comments,
    transformation_summary
)

# Import only valid comments
for comment in valid_transformed:
    importer.import_comment(comment)
```

## Validation Checks

### Extraction Validation
- ✓ event_id present
- ✓ sender present
- ✓ display_name present
- ✓ content not empty
- ✓ timestamp valid
- ✓ room_id present

### Transformation Validation
- ✓ page_id present and starts with '/'
- ✓ author_name present
- ✓ content not empty
- ✓ created_at in ISO 8601 format
- ✓ content_hash in SHA-256 format
- ✓ matrix_event_id present
- ✓ ip_address is "migrated"
- ✓ parent_id valid if present

### Thread Validation
- ✓ Parent comments exist
- ✓ No circular references
- ✓ Thread depth reasonable
- ✓ Orphaned replies reported
- ✓ Chronological order checked

## Demo Script

Run the demo to see validation in action:

```bash
python migrate_comments/utils/demo_data_validator.py
```

## Test Suite

Run tests to verify implementation:

```bash
python -m pytest migrate_comments/utils/test_data_validator.py -v
```

All 26 tests should pass.

## Integration Example

```python
from migrate_comments.extractors.matrix_extractor import MatrixCommentExtractor
from migrate_comments.transformers.comment_transformer import CommentTransformer
from migrate_comments.utils.data_validator import DataValidator

# Initialize components
extractor = MatrixCommentExtractor(...)
transformer = CommentTransformer(...)
validator = DataValidator()

# Extract comments
matrix_comments = extractor.extract_comments_for_page(section_id)

# Validate extraction
extraction_summary = validator.validate_extracted_comments(matrix_comments)
if extraction_summary.invalid_items > 0:
    print(f"Warning: {extraction_summary.invalid_items} invalid extracted comments")
    matrix_comments = validator.filter_valid_comments(matrix_comments, extraction_summary)

# Transform comments
transformed_comments = transformer.transform_batch(matrix_comments, page_id)

# Validate transformation
transformation_summary = validator.validate_transformed_comments(transformed_comments)
if transformation_summary.invalid_items > 0:
    print(f"Warning: {transformation_summary.invalid_items} invalid transformed comments")
    transformed_comments = validator.filter_valid_comments(transformed_comments, transformation_summary)

# Validate thread structure
thread_result = validator.validate_thread_structure(matrix_comments)
if not thread_result.is_valid:
    print(f"Error: Thread validation failed - {thread_result.errors}")
    # Decide whether to proceed

# Import valid comments
for comment in transformed_comments:
    importer.import_comment(comment)
```

## Common Issues

### Issue: High number of orphaned replies
**Cause**: Parent comments were deleted or not extracted
**Solution**: Review extraction logic, consider importing in chronological order

### Issue: Circular references detected
**Cause**: Data corruption in Matrix system
**Solution**: Break circular references by clearing reply_to for one comment

### Issue: Excessive thread depth
**Cause**: Very long conversation threads
**Solution**: This is a warning only - comments can still be imported

### Issue: Invalid timestamps
**Cause**: Timestamp format mismatch
**Solution**: Verify timestamp conversion logic in transformer

## Best Practices

1. **Always validate after extraction** - Catch issues early
2. **Filter invalid comments** - Don't attempt to import invalid data
3. **Log validation results** - Keep audit trail of issues
4. **Use comprehensive validation in dry-run** - Test before actual migration
5. **Review warnings** - They may indicate data quality issues
6. **Handle errors gracefully** - Continue migration even if some comments fail

## Performance Tips

- Validation is fast - don't skip it
- Validate in batches for large migrations
- Use filter_valid_comments to remove invalid data early
- Log validation summaries for monitoring

## Support

For detailed documentation, see:
- `DATA_VALIDATOR_IMPLEMENTATION.md` - Full implementation details
- `test_data_validator.py` - Test examples
- `demo_data_validator.py` - Working examples
