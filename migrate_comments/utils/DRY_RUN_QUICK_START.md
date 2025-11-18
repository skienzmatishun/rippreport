# Dry Run Mode - Quick Start Guide

## What is Dry Run Mode?

Dry run mode lets you test the migration process without actually importing any data to the database. It's like a "preview" that shows you:

- How many comments will be migrated
- What issues or warnings exist
- Sample of what the transformed comments look like
- A GO/NO-GO recommendation

## Quick Start

### 1. Run the Demo

The easiest way to see dry-run in action:

```bash
python migrate_comments/utils/demo_dry_run.py
```

This will:
- Process the first 3 pages from your Hugo site
- Show you what would be migrated
- Generate reports in `dry_run_demo_reports/`
- Give you a recommendation

### 2. Check the Reports

After running, check these files:

```
dry_run_demo_reports/
├── dry_run_preview_YYYYMMDD_HHMMSS.txt  # Human-readable report
└── dry_run_results_YYYYMMDD_HHMMSS.json # Machine-readable data
```

### 3. Review the Recommendation

The dry-run will give you one of three recommendations:

**✓ GO**: Everything looks good, safe to proceed with migration
```
GO: Migration looks good (100.0% success rate, 150 comments ready to import).
```

**⚠ CAUTION**: Some issues found, review before proceeding
```
CAUTION: Moderate error rate (7.5% of comments have issues). 
Review issues carefully before proceeding.
```

**✗ NO-GO**: Too many issues, fix problems before migrating
```
NO-GO: High error rate (12.3% of comments have issues). 
Review and fix issues before proceeding.
```

## Understanding the Preview Report

### Summary Section

```
SUMMARY
--------------------------------------------------------------------------------
Started:  2024-01-01T00:00:00Z
Completed: 2024-01-01T00:00:45Z
Duration: 45.30s

Total Pages:       10
Successful Pages:  9
Pages with Issues: 1

Comments Extracted:    150
Comments Would Import: 148
Total Issues:          2
Total Warnings:        5
```

**What to look for:**
- High success rate (≥95% is good)
- Low number of issues
- Reasonable number of warnings

### Recommendation Section

```
RECOMMENDATION
--------------------------------------------------------------------------------
GO: Migration looks good (90.0% success rate, 148 comments ready to import).
```

**What it means:**
- **GO**: Safe to proceed with actual migration
- **CAUTION**: Review issues, might be okay to proceed
- **NO-GO**: Fix issues before attempting migration

### Pages with Issues

```
PAGES WITH ISSUES
--------------------------------------------------------------------------------

/p/some-article:
  Issues:
    - 1 invalid extracted comments
  Warnings:
    - Thread validation: Found 2 orphaned replies
```

**What to do:**
- Review each issue
- Check if the page has problems in Matrix
- Decide if issues are acceptable or need fixing

### Sample Comments

```
SAMPLE COMMENTS
--------------------------------------------------------------------------------

/p/article-1:
  1. John Doe at 2021-01-01T00:00:00Z
     This is a sample comment that would be imported...
  2. Jane Smith at 2021-01-01T00:01:00Z
     This is a reply to the first comment
     (Reply to parent comment)
```

**What to check:**
- Author names look correct
- Timestamps are reasonable
- Content is properly formatted
- Reply relationships are preserved

## Common Issues and Solutions

### Issue: "No room found for section"

**Cause**: Page doesn't have any comments in Matrix

**Solution**: This is normal for pages without comments, not an error

### Issue: "Orphaned reply"

**Cause**: A reply references a parent comment that doesn't exist

**Solution**: 
- Check if parent was deleted in Matrix
- Reply will be imported as root comment (not a reply)

### Issue: "Invalid timestamp"

**Cause**: Comment has malformed timestamp

**Solution**: Check the comment in Matrix, may need manual fix

### Issue: "Duplicate content hash"

**Cause**: Same comment appears on multiple pages

**Solution**: 
- This is just a warning
- Both comments will be imported
- Check if it's spam or legitimate cross-posting

## Next Steps

### If Recommendation is GO

1. Review the preview report one more time
2. Run the actual migration:
   ```bash
   python migrate_comments/main.py
   ```

### If Recommendation is CAUTION

1. Review all issues and warnings
2. Decide if issues are acceptable
3. If acceptable, proceed with migration
4. If not, fix issues and run dry-run again

### If Recommendation is NO-GO

1. Review all issues carefully
2. Fix problems in source data or code
3. Run dry-run again to verify fixes
4. Repeat until you get GO or CAUTION

## Advanced Usage

### Run Dry-Run for Specific Pages

```python
from migrate_comments.utils.dry_run import DryRunExecutor

# ... initialize components ...

# Run for specific sections
section_ids = ["section1", "section2", "section3"]
summary = dry_run.run_dry_run(section_ids=section_ids)
```

### Customize Sample Size

```python
dry_run = DryRunExecutor(
    extractor=extractor,
    transformer=transformer,
    validator=validator,
    page_id_mapper=page_mapper,
    max_sample_comments=5  # Show 5 samples per page instead of 3
)
```

### Change Output Directory

```python
dry_run = DryRunExecutor(
    extractor=extractor,
    transformer=transformer,
    validator=validator,
    page_id_mapper=page_mapper,
    output_dir="my_custom_reports"
)
```

## Tips

1. **Always run dry-run first**: Never run actual migration without dry-run
2. **Review the full report**: Don't just look at the recommendation
3. **Check sample comments**: Make sure they look correct
4. **Run multiple times**: Dry-run is safe to run repeatedly
5. **Start small**: Test with a few pages first before running all pages

## Troubleshooting

### Dry-run is slow

**Solution**: Run for fewer pages first to test
```python
section_ids = list(page_mapper.get_all_section_ids())[:10]  # First 10 pages
summary = dry_run.run_dry_run(section_ids=section_ids)
```

### Can't find reports

**Solution**: Check the output directory
```python
print(f"Reports saved to: {dry_run.output_dir}")
```

### Getting connection errors

**Solution**: Check Matrix homeserver is accessible
```bash
curl https://matrix.cactus.chat:8448/_matrix/client/versions
```

## Questions?

- Check the full documentation: `DRY_RUN_IMPLEMENTATION.md`
- Review the demo script: `demo_dry_run.py`
- Run the tests: `python -m pytest test_dry_run.py -v`
