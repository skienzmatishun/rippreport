# Google Fallback Handler Implementation

## Overview

The `GoogleFallbackHandler` class provides comprehensive fallback strategies for broken links using Google search, Playwright verification, and LM Studio analysis. This implementation completes Task 5 of the Archive Link Verifier specification.

## Implementation Summary

### Files Created

1. **google_fallback_handler.py** (520 lines)
   - Main implementation of GoogleFallbackHandler class
   - All 8 subtasks implemented

2. **test_google_fallback_handler.py** (350 lines)
   - Comprehensive test suite with 18 tests
   - All tests passing ✓

3. **demo_google_fallback.py** (200 lines)
   - Demo script showing usage of all features

## Features Implemented

### 1. Class Structure (Subtask 5.1) ✓
- Initialized with LM Studio client and Playwright timeout
- Context manager support for resource cleanup
- Lazy Playwright initialization
- Screenshot directory management
- **Requirements met:** 7, 8, 9, 10, 11, 12

### 2. Google Search Functionality (Subtask 5.2) ✓
- Generous rate limiting: 15s base + 10-20s random delay
- Browser-like headers to avoid detection
- URL extraction from Google results
- Site-restricted search support
- Graceful error handling
- **Requirements met:** 7, 11

### 3. PDF Relocation Finder (Subtask 5.3) ✓
- Google search for exact PDF URL
- Playwright screenshot capture
- LM Studio vision analysis for PDF verification
- Content-type header fallback
- Validates PDF vs error page
- **Requirements met:** 7, 14

### 4. Malformed CDN Link Fixer (Subtask 5.4) ✓
- Detects cdn.rippreport.com URLs without extensions
- Extracts and validates link text as URL
- Tests link text URL accessibility
- Archive search fallback suggestion
- **Requirements met:** 8

### 5. Tag Page Converter (Subtask 5.5) ✓
- Extracts tag name from rippreport.com/tags/* URLs
- Generates Google site search URL
- Proper URL encoding
- **Requirements met:** 10

### 6. Missing /p/ Path Fixer (Subtask 5.6) ✓
- Detects rippreport.com URLs missing /p/ segment
- Google search for correct URL
- Verifies result is from rippreport.com
- Confidence scoring
- **Requirements met:** 11

### 7. Playwright Paywall Verification (Subtask 5.7) ✓
- Extended timeout (60s) for slow pages
- Screenshot capture
- LM Studio vision analysis (article vs homepage/paywall)
- Text-based content analysis fallback
- Paywall indicator detection
- **Requirements met:** 9, 13

### 8. LM Studio Relevance Verification (Subtask 5.8) ✓
- Fetches replacement page content
- Sends to LM Studio with post context
- Relevance analysis with confidence scoring
- Handles LM Studio unavailable scenario
- Content truncation for API limits
- **Requirements met:** 12

## Key Design Decisions

### Rate Limiting
- **Base delay:** 15 seconds between searches
- **Random delay:** 10-20 seconds additional
- **Total:** 25-35 seconds per search
- **Rationale:** Generous delays to avoid Google rate limiting/blocking

### Playwright Integration
- **Lazy initialization:** Browser only created when needed
- **Context manager:** Ensures proper cleanup
- **Extended timeouts:** 60 seconds for slow/paywall pages
- **Screenshot support:** For vision model analysis

### LM Studio Integration
- **Vision analysis:** For PDFs and paywall pages
- **Text analysis:** For relevance verification
- **Fallback:** Graceful degradation when unavailable
- **Content limits:** Respects API token limits

### Error Handling
- **Network errors:** Graceful handling with empty results
- **Playwright errors:** Fallback to text analysis
- **LM Studio errors:** Fallback to heuristic methods
- **Timeout errors:** Proper error reporting

## Test Coverage

### Test Classes
1. **TestGoogleSearch** (3 tests)
   - Basic search
   - Site-restricted search
   - Error handling

2. **TestPDFRelocation** (2 tests)
   - Content-type verification
   - No results handling

3. **TestMalformedCDNLinks** (3 tests)
   - Valid link text
   - Non-CDN URLs
   - URLs with extensions

4. **TestTagConversion** (2 tests)
   - Valid tag URLs
   - Invalid URLs

5. **TestMissingPPath** (3 tests)
   - Successful fix
   - Not applicable cases
   - No results

6. **TestPlaywrightVerification** (2 tests)
   - Text analysis
   - Paywall detection

7. **TestRelevanceVerification** (2 tests)
   - Successful verification
   - Fetch errors

8. **TestContextManager** (1 test)
   - Resource cleanup

### Test Results
```
18 passed in 10.20s
```

All tests passing with proper mocking of external dependencies.

## Usage Example

```python
from google_fallback_handler import GoogleFallbackHandler
from lm_studio_client import LMStudioClient

# Create handler with context manager
with GoogleFallbackHandler() as handler:
    # Search Google
    results = handler.search_google("test query", site="rippreport.com")
    
    # Find relocated PDF
    pdf_result = handler.find_relocated_pdf("https://old.com/doc.pdf")
    
    # Fix malformed CDN link
    fix_result = handler.fix_malformed_cdn_link(
        href="https://cdn.rippreport.com/abc",
        link_text="https://example.com/article"
    )
    
    # Convert tag to search
    search_url = handler.convert_tag_to_search(
        "https://rippreport.com/tags/corruption/"
    )
    
    # Fix missing /p/ path
    path_result = handler.fix_missing_p_path(
        "https://rippreport.com/article/"
    )
    
    # Verify with Playwright
    verify_result = handler.verify_with_playwright(
        "https://example.com/article",
        post_context={"title": "Test", "date": "2024-01-01", "excerpt": "..."}
    )
    
    # Verify relevance
    relevance = handler.verify_relevance_with_lm(
        post_title="Test",
        post_date="2024-01-01",
        post_excerpt="...",
        original_url="https://old.com/page",
        replacement_url="https://new.com/page"
    )
```

## Integration with Archive Verifier

The GoogleFallbackHandler is designed to work seamlessly with:
- **ArchiveVerifier:** For archive URL verification
- **LMStudioClient:** For content and relevance analysis
- **BackupManager:** For creating backups before modifications
- **LinkReplacer:** For replacing links in posts

## Configuration

All behavior is controlled through `config.py`:

```python
# Google Fallback Configuration
ENABLE_GOOGLE_FALLBACK = True
GOOGLE_SEARCH_DELAY = 15  # seconds
GOOGLE_RANDOM_DELAY_MIN = 10  # seconds
GOOGLE_RANDOM_DELAY_MAX = 20  # seconds
GOOGLE_MAX_RESULTS = 5
VERIFY_GOOGLE_RESULTS = True

# Link Type Handling
FIX_MALFORMED_CDN_LINKS = True
CONVERT_TAG_PAGES_TO_SEARCH = True
FIX_MISSING_P_PATH = True
VERIFY_PAYWALL_CONTENT = True

# Playwright Configuration
PLAYWRIGHT_TIMEOUT = 60  # seconds
PLAYWRIGHT_WAIT_TIME = 5  # seconds

# Screenshot Settings
SCREENSHOT_DIR = "temp_screenshots"
SCREENSHOT_CLEANUP = True
```

## Performance Considerations

### Rate Limiting
- Google searches are intentionally slow (25-35s each)
- This prevents rate limiting and blocking
- Consider caching results for repeated searches

### Resource Usage
- Playwright browser uses ~100-200MB RAM
- Screenshots are temporary and cleaned up
- LM Studio calls may take 5-10 seconds each

### Optimization Tips
1. Use context manager to ensure cleanup
2. Batch operations when possible
3. Cache Google search results
4. Disable vision analysis if not needed
5. Adjust timeouts based on network speed

## Requirements Verification

### Requirement 7: Google Search Fallback for PDFs ✓
- ✓ Searches Google for exact PDF URL
- ✓ Extracts top result URL
- ✓ Verifies PDF validity (status 200, content-type)
- ✓ Visual verification with screenshots
- ✓ Replaces broken link with new URL

### Requirement 8: Fix Malformed cdn.rippreport.com Links ✓
- ✓ Detects malformed CDN links without extensions
- ✓ Extracts link text content
- ✓ Verifies link text URL accessibility
- ✓ Replaces href with link text URL
- ✓ Falls back to archive search

### Requirement 9: Playwright Verification for Paywall Content ✓
- ✓ Uses Playwright with extended timeouts (60s)
- ✓ Checks for actual content vs error pages
- ✓ Keeps original link if accessible
- ✓ Visual verification with screenshots
- ✓ Marks as "verified_accessible"

### Requirement 10: Convert Missing Tag Pages to Google Search ✓
- ✓ Detects rippreport.com/tags/* URLs
- ✓ Extracts tag name from URL
- ✓ Creates Google site search URL
- ✓ Marks as "converted_to_search"

### Requirement 11: Fix Missing /p/ in rippreport.com URLs ✓
- ✓ Detects URLs missing /p/ path segment
- ✓ Searches Google with slug
- ✓ Extracts first result URL
- ✓ Verifies result is from rippreport.com
- ✓ Marks as "failed" if not from rippreport.com

### Requirement 12: Verify Google-Found Replacements with LM Studio ✓
- ✓ Fetches replacement page content
- ✓ Sends to LM Studio with post context
- ✓ Asks if content is relevant
- ✓ Proceeds with replacement if relevant
- ✓ Rejects if irrelevant
- ✓ Logs warning if LM Studio unavailable

### Requirement 13: Visual Verification with Vision-Capable LLM ✓
- ✓ Captures screenshots with Playwright
- ✓ Sends to vision-capable LM Studio model
- ✓ Includes post context in prompt
- ✓ Determines if page is article or homepage/paywall
- ✓ Marks as inaccessible if homepage/paywall
- ✓ Marks as verified if correct article
- ✓ Falls back to text analysis if unavailable

### Requirement 14: Visual PDF Verification ✓
- ✓ Captures PDF screenshots
- ✓ Sends to vision-capable LM Studio model
- ✓ Asks if image shows PDF content or error
- ✓ Marks as invalid if error page
- ✓ Proceeds if valid PDF
- ✓ Falls back to content-type verification

## Next Steps

This implementation completes Task 5. The next tasks in the specification are:

- **Task 6:** Implement LinkReplacer class
- **Task 7:** Implement ProgressTracker class
- **Task 8:** Implement ReportManager class
- **Task 9:** Implement main orchestration script
- **Task 10:** Create comprehensive test suite
- **Task 11:** Create documentation and usage guide
- **Task 12:** Perform manual testing and validation

## Conclusion

The GoogleFallbackHandler implementation is complete, tested, and ready for integration with the rest of the Archive Link Verifier system. All requirements have been met, and the code follows best practices for error handling, resource management, and testing.
