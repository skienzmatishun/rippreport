# Main Orchestration Script Implementation

## Overview

The `main.py` script is the main entry point for the Archive Link Verifier system. It orchestrates all components to verify and replace broken links in Hugo posts with verified archive URLs.

## Implementation Summary

### Task 9.1: Entry Point and Argument Parsing ✓

**Implemented:**
- Command-line argument parser with configuration overrides
- Support for all major configuration options:
  - File paths (--report, --progress, --content-dir)
  - Behavior flags (--dry-run, --no-playwright, --no-lmstudio, --no-google)
  - Logging configuration (--log-level, --log-file)
- Component initialization with configuration
- Proper logging setup

**Requirements Addressed:** Requirement 13 (Configuration and Dry Run Mode)

### Task 9.2: Link Processing Decision Tree ✓

**Implemented:**
- `determine_strategy()` function that analyzes each entry and determines the appropriate processing strategy
- Strategies implemented:
  - `verify_fixed`: Verify existing archive URLs
  - `pdf_fallback`: Find relocated PDFs via Google
  - `malformed_cdn`: Fix malformed CDN links
  - `tag_conversion`: Convert tag pages to Google searches
  - `missing_p_path`: Fix missing /p/ paths
  - `paywall_check`: Verify paywall content accessibility
  - `skip`: Skip entries with no applicable strategy

**Requirements Addressed:** Requirements 2, 3, 4, 5, 7, 8, 9, 10, 11, 12

### Task 9.3: Status="fixed" Verification Flow ✓

**Implemented:**
- `process_verify_fixed()` function
- Verifies archive URLs with Playwright + LM Studio
- Extracts timestamp from archive URL
- Tries previous snapshots if current archive is blank
- Updates report with verification results

**Requirements Addressed:** Requirements 2, 3, 5

### Task 9.4: PDF Fallback Flow ✓

**Implemented:**
- `process_pdf_fallback()` function
- Detects PDF links without archives
- Uses Google search to find relocated PDFs
- Captures screenshot and verifies with vision model
- Verifies PDF validity (visual + content-type)

**Requirements Addressed:** Requirements 7, 14

### Task 9.5: Malformed CDN Link Flow ✓

**Implemented:**
- `process_malformed_cdn()` function
- Detects malformed cdn.rippreport.com links
- Extracts and verifies link text URL
- Replaces or searches for archive

**Requirements Addressed:** Requirement 8

### Task 9.6: Tag Conversion Flow ✓

**Implemented:**
- `process_tag_conversion()` function
- Detects rippreport.com/tags/* URLs
- Converts to Google site search

**Requirements Addressed:** Requirement 10

### Task 9.7: Missing /p/ Path Flow ✓

**Implemented:**
- `process_missing_p_path()` function
- Detects rippreport.com URLs without /p/
- Searches Google for correct URL
- Verifies relevance with LM Studio (when enabled)

**Requirements Addressed:** Requirements 11, 12

### Task 9.8: Paywall Verification Flow ✓

**Implemented:**
- `process_paywall_check()` function
- For links without archives, tries Playwright verification
- Captures screenshot and analyzes with vision model
- Checks if content is accessible and shows correct article
- Keeps original link if accessible and correct

**Requirements Addressed:** Requirements 9, 13

### Task 9.9: Link Replacement ✓

**Implemented:**
- `replace_link_in_post()` function
- Creates timestamped backup via BackupManager
- Replaces link in post file via LinkReplacer
- Updates report with replacement status
- Saves progress after each link

**Requirements Addressed:** Requirements 1, 4, 6

### Task 9.10: Error Handling and Resumability ✓

**Implemented:**
- Graceful handling of interruptions (Ctrl+C)
- Skips already-processed links on restart
- Displays progress on resume
- Try-except blocks around all processing
- Continues processing on individual entry failures
- Saves progress after each entry

**Requirements Addressed:** Requirement 6

### Task 9.11: Final Summary Generation ✓

**Implemented:**
- Generates comprehensive summary report via ProgressTracker
- Displays statistics via ReportManager
- Saves updated report files
- Returns appropriate exit codes

**Requirements Addressed:** Requirement 6

## Usage

### Basic Usage

```bash
# Run with default settings
python3 -m archive_link_verifier.main

# Dry run mode (no file modifications)
python3 -m archive_link_verifier.main --dry-run

# Custom report and progress files
python3 -m archive_link_verifier.main --report my_report.json --progress my_progress.json

# Disable specific features
python3 -m archive_link_verifier.main --no-playwright --no-lmstudio
python3 -m archive_link_verifier.main --no-google

# Verbose logging
python3 -m archive_link_verifier.main --log-level DEBUG
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `--report` | Path to archive fix report | `archive_fix_report.json` |
| `--progress` | Path to progress file | `archive_verifier_progress.json` |
| `--content-dir` | Content directory | `content/p` |
| `--dry-run` | Don't modify files | `False` |
| `--no-playwright` | Disable Playwright | `False` |
| `--no-lmstudio` | Disable LM Studio | `False` |
| `--no-google` | Disable Google fallback | `False` |
| `--log-level` | Logging level | `INFO` |
| `--log-file` | Log file path | `archive_verifier.log` |

## Architecture

### Component Flow

```
main.py
├── parse_arguments() - Parse CLI args
├── initialize_components() - Initialize all components
│   ├── BackupManager
│   ├── LMStudioClient
│   ├── ProgressTracker
│   ├── ReportManager
│   ├── LinkReplacer
│   ├── ArchiveVerifier (config)
│   └── GoogleFallbackHandler (config)
├── Load progress and report
├── For each entry:
│   ├── determine_strategy() - Decide processing approach
│   ├── process_link() - Execute strategy
│   │   ├── process_verify_fixed()
│   │   ├── process_pdf_fallback()
│   │   ├── process_malformed_cdn()
│   │   ├── process_tag_conversion()
│   │   ├── process_missing_p_path()
│   │   └── process_paywall_check()
│   ├── replace_link_in_post() - Replace link with backup
│   ├── Update report
│   └── Save progress
└── Generate final summary
```

### Decision Tree

```
Entry → Already processed? → Yes → Skip
                          → No → Determine strategy
                                 ├── status="fixed" → Verify archive
                                 ├── PDF + no archive → Google search
                                 ├── Malformed CDN → Fix link text
                                 ├── Tag page → Convert to search
                                 ├── Missing /p/ → Google search
                                 ├── No archive → Paywall check
                                 └── None applicable → Skip
```

## Error Handling

The script handles errors at multiple levels:

1. **Component Initialization**: Fails fast if components can't be initialized
2. **Progress/Report Loading**: Handles corrupted files gracefully
3. **Entry Processing**: Continues on individual entry failures
4. **Interruption**: Saves progress on Ctrl+C
5. **File Operations**: Validates paths and handles missing files

## Exit Codes

- `0`: Success
- `1`: Error (initialization, fatal error)
- `130`: Interrupted by user (Ctrl+C)

## Resumability

The script is fully resumable:

1. Progress is saved after each entry
2. On restart, already-processed entries are skipped
3. Progress summary is displayed on resume
4. Ctrl+C saves progress before exiting

## Testing

To test the script:

```bash
# Test help output
python3 -m archive_link_verifier.main --help

# Test dry run mode
python3 -m archive_link_verifier.main --dry-run

# Test with minimal features
python3 -m archive_link_verifier.main --no-playwright --no-lmstudio --no-google --dry-run
```

## Implementation Notes

- Uses context managers for Playwright resources (ArchiveVerifier, GoogleFallbackHandler)
- Saves progress after each entry for maximum resumability
- Supports graceful interruption with Ctrl+C
- All strategies are modular and can be disabled via configuration
- Comprehensive logging at INFO level, DEBUG for troubleshooting
- Dry-run mode for safe testing

## Requirements Coverage

All requirements from the design document are addressed:

- ✓ Requirement 1: Timestamped backups (via BackupManager)
- ✓ Requirement 2: Archive verification (via ArchiveVerifier)
- ✓ Requirement 3: LM Studio analysis (via LMStudioClient)
- ✓ Requirement 4: Link replacement (via LinkReplacer)
- ✓ Requirement 5: Re-verification (process_verify_fixed)
- ✓ Requirement 6: Progress tracking (via ProgressTracker)
- ✓ Requirement 7: Google PDF search (process_pdf_fallback)
- ✓ Requirement 8: Malformed CDN links (process_malformed_cdn)
- ✓ Requirement 9: Paywall verification (process_paywall_check)
- ✓ Requirement 10: Tag conversion (process_tag_conversion)
- ✓ Requirement 11: Missing /p/ path (process_missing_p_path)
- ✓ Requirement 12: Relevance verification (via LMStudioClient)
- ✓ Requirement 13: Visual verification (via LMStudioClient vision)
- ✓ Requirement 14: PDF visual verification (via LMStudioClient vision)
- ✓ Requirement 15: Configuration and dry-run (argument parsing)
