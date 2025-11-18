# Quick Start Guide

Get up and running with the Comment Migration Tool in 5 minutes.

## Prerequisites

- Python 3.8+ installed
- Wrangler CLI installed and authenticated
- Access to your Cloudflare D1 database

## Installation (2 minutes)

```bash
# 1. Navigate to the migration tool directory
cd migrate_comments

# 2. Create virtual environment
python3 -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# 3. Install dependencies
pip install -r requirements.txt

# 4. Verify installation
python test_setup.py
```

## Configuration (2 minutes)

```bash
# 1. Copy example config
cp config.yaml.example config.yaml

# 2. Edit config.yaml with your values
# Required fields:
#   - matrix.site_name: "your-site.com"
#   - cloudflare.api_base: "https://comments.your-site.com/api"
#   - cloudflare.database_id: "your-database-id"
#   - hugo.base_url: "https://your-site.com"
```

**Get your database ID:**
```bash
wrangler d1 list
```

## Test Run (1 minute)

```bash
# Preview migration without importing
python migrate.py --dry-run

# Check the preview report
cat dry_run_reports/dry_run_preview_*.txt
```

## Full Migration

```bash
# Run the migration
python migrate.py

# Check results
cat migration_reports/summary_report_*.txt
```

## What's Next?

1. **If migration succeeded**: 
   - Review the summary report
   - Spot-check comments on your site
   - Update Hugo shortcodes (see Task 11)

2. **If migration failed**:
   - Check [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
   - Review `migration.log`
   - Use `python migrate.py --resume` to retry

3. **For more details**:
   - [INSTALLATION.md](INSTALLATION.md) - Detailed installation
   - [CONFIGURATION.md](CONFIGURATION.md) - All config options
   - [USAGE.md](USAGE.md) - Complete usage guide
   - [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Problem solving

## Common Commands

```bash
# Preview without importing
python migrate.py --dry-run

# Full migration
python migrate.py

# Resume after failure
python migrate.py --resume

# Migrate specific pages
python migrate.py --pages /p/page-1 /p/page-2

# Debug mode
python migrate.py --verbose

# Get help
python migrate.py --help
```

## Need Help?

1. Check [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
2. Review logs: `tail -100 migration.log`
3. Run with debug: `python migrate.py --verbose`
