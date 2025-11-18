# Matrix Authentication Fix

## Problem

The migration tool was finding 115 pages but 0 comments. Investigation revealed two issues:

1. **Missing Authentication**: Matrix rooms required authentication to access messages
2. **Wrong Room Alias Format**: Room aliases were missing the "comments_" prefix

## Root Cause

### Issue 1: Authentication Required

When attempting to access Matrix room messages, the API returned:
```json
{"errcode":"M_MISSING_TOKEN","error":"Missing access token"}
```

The Cactus Chat Matrix rooms require authentication to read messages, even though they're publicly accessible.

### Issue 2: Incorrect Room Alias

The extractor was constructing room aliases as:
```
#rippreport.com_mljp23:cactus.chat
```

But Cactus Chat actually uses:
```
#comments_rippreport.com_mljp23:cactus.chat
```

The "comments_" prefix was missing.

## Solution

### Fix 1: Guest Authentication

Added guest user registration to the Matrix extractor:

```python
def _ensure_authenticated(self):
    """Register as guest user and obtain access token."""
    if self._access_token:
        return
    
    url = f"{self.homeserver_url}/_matrix/client/r0/register?kind=guest"
    response = self.session.post(url, json={}, timeout=30)
    
    data = response.json()
    self._access_token = data.get('access_token')
    
    # Update session headers with auth token
    self.session.headers.update({
        'Authorization': f'Bearer {self._access_token}'
    })
```

This method is called before any API requests that require authentication (getting messages, display names, etc.).

### Fix 2: Correct Room Alias Format

Updated the room alias construction:

```python
# OLD (incorrect)
room_alias = f"#{self.site_name}_{page_section_id}:{self.server_name}"

# NEW (correct)
room_alias = f"#comments_{self.site_name}_{page_section_id}:{self.server_name}"
```

## Changes Made

### File: `migrate_comments/extractors/matrix_extractor.py`

1. **Added `_access_token` attribute** to store guest access token
2. **Added `_ensure_authenticated()` method** to register as guest
3. **Updated `get_room_messages()`** to call `_ensure_authenticated()`
4. **Updated `get_display_name()`** to call `_ensure_authenticated()`
5. **Fixed room alias format** to include "comments_" prefix

### File: `migrate_comments/config.yaml`

1. **Updated `hugo.content_dir`** from `"content"` to `"../content"` to fix path resolution

## Verification

### Before Fix
```
Total Pages:       115
Comments Found:    0
```

### After Fix
```
Total Pages:       115
Comments Found:    655
Successful Pages:  115
```

### Test Results

Testing with sample sections:
- `h452`: Found 32 comments
- `jfmfhbc`: Found 4 comments
- Total across all pages: 655 comments

## Technical Details

### Guest Authentication Flow

1. Tool makes POST request to `/_matrix/client/r0/register?kind=guest`
2. Matrix server returns access token
3. Tool adds `Authorization: Bearer <token>` header to all subsequent requests
4. Token is cached for the session

### Room Discovery

1. Construct room alias: `#comments_{site_name}_{section_id}:{server_name}`
2. Resolve alias to room ID via `/_matrix/client/r0/directory/room/{alias}`
3. Fetch messages from room using room ID

### Message Access

With guest authentication, the tool can:
- Read room messages
- Get user display names
- Access room state
- Paginate through message history

## Impact

This fix enables the migration tool to:
- ✅ Successfully extract comments from Cactus Chat
- ✅ Access all 115 pages with comments
- ✅ Retrieve all 655 comments
- ✅ Preserve comment metadata (author, timestamp, replies)
- ✅ Complete the migration process

## Future Considerations

### Guest Token Expiration

Guest tokens may expire after some time. If migrations take a very long time, the tool may need to:
- Detect 401 errors during migration
- Re-register as guest
- Retry failed requests

Currently not implemented as migrations are expected to complete quickly.

### Alternative Authentication

If guest access is disabled in the future, alternatives include:
- User authentication with username/password
- Application service tokens
- Direct database access (if available)

## Testing

To test the fix:

```bash
# Test authentication
cd migrate_comments
python test_extractor_auth.py

# Test full dry-run
python migrate.py --dry-run

# Check results
cat migration_reports/dry_run_preview_*.txt
```

## Date

Fixed: November 17, 2024

## Status

✅ **RESOLVED** - Migration tool now successfully extracts all comments from Matrix.
