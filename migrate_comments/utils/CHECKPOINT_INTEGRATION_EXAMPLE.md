# CheckpointManager Integration Example

## Complete Migration Workflow with Checkpointing

This example demonstrates how the CheckpointManager integrates with other migration components to provide robust recovery capabilities.

## Full Integration Example

```python
#!/usr/bin/env python3
"""
Complete migration script with checkpoint support.

This script demonstrates the full integration of CheckpointManager
with other migration components.
"""

import sys
import time
from pathlib import Path

# Import migration components
from utils import (
    Config,
    setup_logger_from_config,
    CheckpointManager
)
from extractors.matrix_extractor import MatrixCommentExtractor
from transformers.comment_transformer import CommentTransformer
from importers.cloudflare_client import CloudflareCommentClient
from utils.page_id_mapper import PageIdMapper


def main():
    """Main migration function with checkpoint support."""
    
    # 1. Load configuration
    print("Loading configuration...")
    config = Config('config.yaml')
    
    # 2. Set up logging
    logger = setup_logger_from_config(config.config)
    logger.info("Starting comment migration with checkpoint support")
    
    # 3. Initialize checkpoint manager
    checkpoint_file = config.get('migration.checkpoint_file', '.migration_checkpoint.json')
    checkpoint = CheckpointManager(
        checkpoint_file=checkpoint_file,
        logger=logger
    )
    
    # Check if resuming from previous run
    if checkpoint.started_at:
        logger.info(
            f"Resuming migration from checkpoint: "
            f"{len(checkpoint.completed_pages)} pages already completed"
        )
    else:
        logger.info("Starting new migration")
    
    # 4. Initialize components
    logger.info("Initializing migration components...")
    
    # Matrix extractor
    matrix_config = config.get_matrix_config()
    extractor = MatrixCommentExtractor(
        homeserver_url=matrix_config['homeserver_url'],
        server_name=matrix_config['server_name'],
        site_name=matrix_config['site_name'],
        logger=logger
    )
    
    # Page ID mapper
    hugo_config = config.get_hugo_config()
    page_mapper = PageIdMapper(
        content_dir=hugo_config['content_dir'],
        base_url=hugo_config['base_url'],
        logger=logger
    )
    
    # Comment transformer
    transformer = CommentTransformer(
        page_id_mapper=page_mapper,
        logger=logger
    )
    
    # Cloudflare client
    cloudflare_config = config.get_cloudflare_config()
    importer = CloudflareCommentClient(
        api_base=cloudflare_config['api_base'],
        database_id=cloudflare_config['database_id'],
        database_name=cloudflare_config['database_name'],
        wrangler_path=cloudflare_config.get('wrangler_path', 'wrangler'),
        logger=logger
    )
    
    # 5. Scan Hugo pages for section IDs
    logger.info("Scanning Hugo pages...")
    section_mapping = page_mapper.scan_hugo_pages()
    pages = list(section_mapping.items())  # [(section_id, page_id), ...]
    
    logger.info(f"Found {len(pages)} pages with chat shortcodes")
    
    # 6. Start migration (or resume)
    if not checkpoint.started_at:
        checkpoint.start_migration(total_pages=len(pages))
    
    # 7. Process each page
    start_time = time.time()
    
    for section_id, page_id in pages:
        # Skip already-completed pages
        if checkpoint.should_skip_page(page_id):
            logger.info(f"Skipping completed page: {page_id}")
            continue
        
        # Check if page previously failed (for logging)
        if checkpoint.is_page_failed(page_id):
            logger.info(f"Retrying failed page: {page_id}")
        
        page_start_time = time.time()
        
        try:
            # Extract comments from Matrix
            logger.info(f"Extracting comments for {page_id} (section: {section_id})")
            matrix_comments = extractor.extract_comments_for_page(section_id)
            
            if not matrix_comments:
                logger.info(f"No comments found for {page_id}")
                checkpoint.mark_page_complete(page_id, section_id, 0)
                continue
            
            logger.info(f"Extracted {len(matrix_comments)} comments")
            checkpoint.update_statistics(comments_extracted=len(matrix_comments))
            
            # Transform comments
            logger.info(f"Transforming comments for {page_id}")
            transformed_comments = transformer.transform_batch(
                matrix_comments,
                page_id
            )
            
            # Import comments (maintaining parent-child relationships)
            logger.info(f"Importing {len(transformed_comments)} comments")
            imported_count = 0
            failed_count = 0
            
            for comment in transformed_comments:
                try:
                    # Resolve parent ID if this is a reply
                    if comment.get('reply_to'):
                        parent_matrix_id = comment['reply_to']
                        parent_new_id = checkpoint.get_new_comment_id(parent_matrix_id)
                        
                        if parent_new_id:
                            comment['parent_id'] = parent_new_id
                        else:
                            logger.warning(
                                f"Parent comment not found for reply: {parent_matrix_id}"
                            )
                            comment['parent_id'] = None
                    
                    # Import comment
                    result = importer.create_comment_with_timestamp(comment)
                    new_comment_id = result['id']
                    
                    # Store ID mapping for future parent resolution
                    matrix_event_id = comment['matrix_event_id']
                    checkpoint.add_id_mapping(matrix_event_id, new_comment_id)
                    
                    imported_count += 1
                    
                except Exception as e:
                    logger.error(f"Failed to import comment: {e}")
                    failed_count += 1
            
            # Update statistics
            checkpoint.update_statistics(
                comments_imported=imported_count,
                failed_imports=failed_count
            )
            
            # Mark page complete
            page_duration = time.time() - page_start_time
            checkpoint.mark_page_complete(page_id, section_id, imported_count)
            
            logger.info(
                f"Completed {page_id}: {imported_count}/{len(transformed_comments)} "
                f"comments imported in {page_duration:.2f}s"
            )
            
        except Exception as e:
            # Mark page as failed
            logger.error(f"Failed to migrate page {page_id}: {e}")
            checkpoint.mark_page_failed(
                page_id=page_id,
                section_id=section_id,
                error=str(e),
                comment_count=len(matrix_comments) if 'matrix_comments' in locals() else 0
            )
            
            # Continue with next page
            continue
    
    # 8. Generate final report
    total_duration = time.time() - start_time
    progress = checkpoint.get_progress_summary()
    
    logger.info("=" * 60)
    logger.info("Migration Complete")
    logger.info("=" * 60)
    logger.info(f"Total pages: {progress['total_pages']}")
    logger.info(f"Completed: {progress['completed_pages']}")
    logger.info(f"Failed: {progress['failed_pages']}")
    logger.info(f"Comments imported: {progress['total_comments_imported']}")
    logger.info(f"Failed imports: {progress['failed_imports']}")
    logger.info(f"Duration: {total_duration:.2f}s")
    logger.info("=" * 60)
    
    # 9. Clear checkpoint if all pages completed successfully
    if progress['failed_pages'] == 0 and progress['remaining_pages'] == 0:
        logger.info("All pages migrated successfully, clearing checkpoint")
        checkpoint.clear_checkpoint()
        return 0
    else:
        logger.warning(
            f"{progress['failed_pages']} pages failed, "
            f"{progress['remaining_pages']} pages remaining"
        )
        logger.info("Checkpoint preserved for recovery")
        return 1


if __name__ == '__main__':
    try:
        sys.exit(main())
    except KeyboardInterrupt:
        print("\nMigration interrupted by user")
        print("Progress has been saved. Run again to resume.")
        sys.exit(130)
    except Exception as e:
        print(f"Fatal error: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
```

## Key Integration Points

### 1. Initialization

```python
# Create checkpoint manager early in the process
checkpoint = CheckpointManager(
    checkpoint_file='.migration_checkpoint.json',
    logger=logger
)

# Check if resuming
if checkpoint.started_at:
    logger.info("Resuming from checkpoint")
else:
    checkpoint.start_migration(total_pages=len(pages))
```

### 2. Page Processing Loop

```python
for section_id, page_id in pages:
    # Skip completed pages
    if checkpoint.should_skip_page(page_id):
        continue
    
    # Check if retrying failed page
    if checkpoint.is_page_failed(page_id):
        logger.info("Retrying failed page")
    
    try:
        # ... extract, transform, import ...
        checkpoint.mark_page_complete(page_id, section_id, count)
    except Exception as e:
        checkpoint.mark_page_failed(page_id, section_id, str(e), count)
```

### 3. ID Mapping for Parent Resolution

```python
# During import, store mapping
result = importer.create_comment(comment)
checkpoint.add_id_mapping(comment['matrix_event_id'], result['id'])

# Later, resolve parent for reply
if comment.get('reply_to'):
    parent_new_id = checkpoint.get_new_comment_id(comment['reply_to'])
    if parent_new_id:
        comment['parent_id'] = parent_new_id
```

### 4. Statistics Tracking

```python
# Update statistics as you go
checkpoint.update_statistics(
    comments_extracted=len(matrix_comments),
    comments_imported=imported_count,
    failed_imports=failed_count
)
```

### 5. Final Cleanup

```python
# Get final progress
progress = checkpoint.get_progress_summary()

# Clear checkpoint if successful
if progress['failed_pages'] == 0:
    checkpoint.clear_checkpoint()
```

## Recovery Scenarios

### Scenario 1: Network Failure

```
Initial run:
- Pages 1-50: ✓ Completed
- Page 51: ✗ Network timeout
- Pages 52-100: Not attempted

Checkpoint saved with:
- 50 completed pages
- 1 failed page
- ID mappings for 500 comments

On restart:
- Loads checkpoint
- Skips pages 1-50
- Retries page 51
- Continues with pages 52-100
```

### Scenario 2: API Rate Limiting

```
Initial run:
- Pages 1-30: ✓ Completed
- Page 31: ✗ Rate limit exceeded
- Process terminated

Checkpoint saved with:
- 30 completed pages
- 1 failed page
- ID mappings for 300 comments

On restart (after rate limit reset):
- Loads checkpoint
- Skips pages 1-30
- Retries page 31
- Continues normally
```

### Scenario 3: User Interruption (Ctrl+C)

```
Initial run:
- Pages 1-75: ✓ Completed
- User presses Ctrl+C

Checkpoint saved with:
- 75 completed pages
- ID mappings for 750 comments

On restart:
- Loads checkpoint
- Skips pages 1-75
- Continues with pages 76-100
```

## Error Handling Best Practices

### 1. Graceful Degradation

```python
try:
    # Attempt migration
    migrate_page(page)
    checkpoint.mark_page_complete(page.id, page.section_id, count)
except Exception as e:
    # Log error but continue
    logger.error(f"Failed: {e}")
    checkpoint.mark_page_failed(page.id, page.section_id, str(e), count)
    # Don't raise - continue with next page
```

### 2. Atomic Operations

```python
# Checkpoint saves are atomic
# If save fails, previous checkpoint is preserved
checkpoint.mark_page_complete(...)  # Saves automatically
```

### 3. Signal Handling

```python
import signal

def signal_handler(signum, frame):
    logger.info("Received interrupt signal, saving checkpoint...")
    # Checkpoint is already saved after each page
    logger.info("Checkpoint saved. Safe to exit.")
    sys.exit(130)

signal.signal(signal.SIGINT, signal_handler)
signal.signal(signal.SIGTERM, signal_handler)
```

## Monitoring Progress

### Real-time Progress

```python
# Get progress at any time
progress = checkpoint.get_progress_summary()

print(f"Progress: {progress['completed_pages']}/{progress['total_pages']}")
print(f"Success rate: {progress['progress_percentage']:.1f}%")
print(f"Comments imported: {progress['total_comments_imported']}")
```

### Progress Bar Integration

```python
from tqdm import tqdm

pages = list(section_mapping.items())
progress_bar = tqdm(total=len(pages), desc="Migrating pages")

# Update based on checkpoint
progress_bar.update(len(checkpoint.completed_pages))

for section_id, page_id in pages:
    if checkpoint.should_skip_page(page_id):
        continue
    
    # ... migrate page ...
    
    progress_bar.update(1)

progress_bar.close()
```

## Conclusion

The CheckpointManager provides robust recovery capabilities that integrate seamlessly with the migration workflow. Key benefits:

1. **Automatic recovery** - No manual intervention needed
2. **Minimal overhead** - Checkpoint saves are fast and atomic
3. **Complete state preservation** - All progress and mappings saved
4. **Flexible retry logic** - Failed pages can be retried
5. **Progress tracking** - Real-time visibility into migration status

The integration is straightforward and requires minimal changes to the migration logic, making it easy to add checkpoint support to any migration workflow.
