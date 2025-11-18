# Shortcode Replacer Verification Checklist

## Pre-Replacement Verification

### Environment Check
- [ ] Python 3.8+ installed
- [ ] Required dependencies installed (`pip install -r requirements.txt`)
- [ ] Hugo content directory exists and is accessible
- [ ] Sufficient disk space for backups

### Content Check
- [ ] Content directory contains markdown files
- [ ] Files contain `{{< chat` shortcodes
- [ ] Files are not currently open in editors
- [ ] Git repository is clean (optional but recommended)

### Configuration Check
- [ ] Content directory path is correct
- [ ] Backup directory path is writable
- [ ] Logging is configured properly

## Dry Run Verification

### Run Dry Run
```bash
python migrate_comments/utils/demo_shortcode_replacer.py
```

### Check Output
- [ ] Correct number of files found
- [ ] Shortcode IDs are correct
- [ ] Line numbers are accurate
- [ ] No unexpected errors
- [ ] Replacement preview looks correct

### Sample Verification
- [ ] Manually check 3-5 sample files
- [ ] Verify shortcode format matches pattern
- [ ] Confirm IDs will be preserved
- [ ] Check for any special cases

## Replacement Verification

### Run Replacement
```python
from utils.shortcode_replacer import ShortcodeReplacer

replacer = ShortcodeReplacer("content")
replacer.scan_files()
stats = replacer.replace_all(dry_run=False)
```

### Check Statistics
- [ ] `stats['successful']` matches expected file count
- [ ] `stats['failed']` is 0
- [ ] `stats['total_replacements']` matches expected shortcode count
- [ ] No errors in `stats['errors']`
- [ ] No issues in `stats['verification_issues']`

### Backup Verification
- [ ] Backup directory created with timestamp
- [ ] `backup_metadata.json` exists
- [ ] All modified files backed up
- [ ] Backup files are readable
- [ ] Metadata contains correct file paths

## Post-Replacement Verification

### File Content Check
```bash
# Check for remaining chat shortcodes
grep -r "{{< chat" content/

# Should return no results
```

- [ ] No `{{< chat` patterns remain
- [ ] All replaced with `{{< aicomments`
- [ ] Shortcode IDs preserved correctly
- [ ] No duplicate replacements

### Sample File Inspection
Manually check 5-10 files:
- [ ] Shortcode format is correct
- [ ] IDs match original
- [ ] No extra spaces or formatting issues
- [ ] Markdown structure intact
- [ ] No broken links or references

### Hugo Build Check
```bash
hugo
```

- [ ] Hugo builds without errors
- [ ] No shortcode-related warnings
- [ ] Generated HTML looks correct
- [ ] No broken pages

### Site Functionality Check
```bash
hugo server
```

Visit several pages and verify:
- [ ] Pages load correctly
- [ ] Comment sections appear
- [ ] No JavaScript errors in console
- [ ] Shortcodes render properly
- [ ] No layout issues

## Automated Verification

### Run Verification Script
```python
from pathlib import Path
from utils.shortcode_replacer import ShortcodeReplacer

replacer = ShortcodeReplacer("content")
replacer.scan_files()

# Verify all modified files
issues_found = []
for file_info in replacer.files_to_modify:
    success, issues = replacer.verify_replacement(file_info['path'])
    if not success:
        issues_found.append({
            'file': file_info['relative_path'],
            'issues': issues
        })

if issues_found:
    print(f"✗ Found issues in {len(issues_found)} files")
    for item in issues_found:
        print(f"\n{item['file']}:")
        for issue in item['issues']:
            print(f"  - {issue}")
else:
    print("✓ All files verified successfully")
```

### Check Results
- [ ] All files pass verification
- [ ] No remaining chat shortcodes
- [ ] No duplicate replacements
- [ ] No malformed shortcodes

## Rollback Test (Optional)

### Test Rollback Capability
```python
from pathlib import Path
from utils.shortcode_replacer import ShortcodeReplacer

# Find backup metadata
backup_file = Path("backups/shortcode-replacement/[timestamp]/backup_metadata.json")

# Test rollback on one file
replacer = ShortcodeReplacer("content")
stats = replacer.rollback(backup_file)
```

- [ ] Rollback completes successfully
- [ ] Original content restored
- [ ] Files are readable
- [ ] No data loss

### Re-run Replacement
After testing rollback:
- [ ] Re-run replacement
- [ ] Verify results again
- [ ] Confirm everything works

## Integration Testing

### Test with Comment System
- [ ] Visit pages with comments
- [ ] Verify comments load
- [ ] Test posting new comments
- [ ] Test replying to comments
- [ ] Check comment threading

### Test Different Page Types
- [ ] Regular posts
- [ ] Pages with multiple shortcodes
- [ ] Pages with other shortcodes
- [ ] Pages with special characters
- [ ] Pages in subdirectories

## Performance Verification

### Check Processing Time
- [ ] Replacement completed in reasonable time
- [ ] No timeouts or hangs
- [ ] Memory usage acceptable
- [ ] No performance degradation

### Check File Sizes
- [ ] File sizes unchanged (except for shortcode text)
- [ ] No unexpected file growth
- [ ] Backup sizes reasonable

## Documentation Verification

### Check Documentation
- [ ] README updated if needed
- [ ] Implementation docs accurate
- [ ] Quick start guide tested
- [ ] API reference correct

### Check Examples
- [ ] Demo script works
- [ ] Code examples run
- [ ] Test cases pass

## Final Checklist

### Before Deployment
- [ ] All verifications passed
- [ ] No errors or warnings
- [ ] Backups are safe
- [ ] Git changes committed
- [ ] Deployment plan ready

### After Deployment
- [ ] Site deployed successfully
- [ ] Comments work on live site
- [ ] No user-reported issues
- [ ] Monitoring shows no errors
- [ ] Backups can be archived

## Issue Resolution

### If Issues Found

1. **Stop immediately**
2. **Document the issue**
3. **Check logs for details**
4. **Determine if rollback needed**
5. **Fix the issue**
6. **Re-run verification**

### Common Issues and Solutions

#### Issue: Some chat shortcodes remain
**Solution**: Check regex pattern, may need to handle edge cases

#### Issue: Verification fails
**Solution**: Inspect specific files, check for malformed shortcodes

#### Issue: Hugo build fails
**Solution**: Check Hugo logs, verify shortcode syntax

#### Issue: Comments don't load
**Solution**: Check JavaScript console, verify API endpoints

## Sign-Off

### Verification Complete
- [ ] All checks passed
- [ ] No issues found
- [ ] Ready for deployment
- [ ] Backups secured
- [ ] Documentation updated

**Verified by:** _______________  
**Date:** _______________  
**Notes:** _______________

## Rollback Plan

### If Rollback Needed

1. Locate backup metadata file
2. Run rollback script
3. Verify restoration
4. Test site functionality
5. Investigate root cause
6. Fix issues
7. Re-run replacement

### Rollback Command
```python
from pathlib import Path
from utils.shortcode_replacer import ShortcodeReplacer

replacer = ShortcodeReplacer("content")
backup_file = Path("backups/shortcode-replacement/[timestamp]/backup_metadata.json")
stats = replacer.rollback(backup_file)

print(f"Restored {stats['successful']}/{stats['total_files']} files")
```

## Support

If you encounter issues:
1. Check logs in `migration.log`
2. Review error messages
3. Consult documentation
4. Check test cases for examples
5. Contact support if needed
