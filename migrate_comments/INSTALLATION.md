# Installation Guide

This guide will walk you through installing and setting up the Comment Migration Tool.

## Prerequisites

### Python Version Requirements

The Comment Migration Tool requires **Python 3.8 or higher**. We recommend Python 3.9 or later for best compatibility.

To check your Python version:

```bash
python --version
# or
python3 --version
```

If you need to install or upgrade Python:
- **macOS**: Use Homebrew: `brew install python@3.9`
- **Linux**: Use your package manager: `apt install python3.9` or `yum install python39`
- **Windows**: Download from [python.org](https://www.python.org/downloads/)

### Wrangler CLI Requirements

The migration tool uses Cloudflare's Wrangler CLI to interact with the D1 database. You must have Wrangler installed and authenticated.

#### Installing Wrangler

**Option 1: Using npm (recommended)**
```bash
npm install -g wrangler
```

**Option 2: Using yarn**
```bash
yarn global add wrangler
```

**Option 3: Using pnpm**
```bash
pnpm add -g wrangler
```

Verify installation:
```bash
wrangler --version
```

#### Authenticating Wrangler

Before running the migration, authenticate with Cloudflare:

```bash
wrangler login
```

This will open a browser window for you to authorize Wrangler with your Cloudflare account.

Verify authentication:
```bash
wrangler whoami
```

### System Requirements

- **Operating System**: macOS, Linux, or Windows
- **Disk Space**: At least 100MB free for dependencies and migration artifacts
- **Memory**: Minimum 512MB RAM (1GB+ recommended for large migrations)
- **Network**: Stable internet connection for API access

## Installation Steps

### Step 1: Clone or Navigate to the Project

Navigate to your Hugo site's root directory:

```bash
cd /path/to/your/hugo/site
```

The `migrate_comments/` directory should already exist in your project.

### Step 2: Create a Virtual Environment

Creating a virtual environment isolates the migration tool's dependencies from your system Python packages.

**On macOS/Linux:**
```bash
cd migrate_comments
python3 -m venv venv
source venv/bin/activate
```

**On Windows:**
```bash
cd migrate_comments
python -m venv venv
venv\Scripts\activate
```

You should see `(venv)` in your terminal prompt, indicating the virtual environment is active.

### Step 3: Install Python Dependencies

With the virtual environment activated, install the required packages:

```bash
pip install -r requirements.txt
```

This will install:
- **PyYAML** (≥6.0.1): YAML configuration file parsing
- **requests** (≥2.31.0): HTTP requests to Matrix and Cloudflare APIs
- **python-dateutil** (≥2.8.2): Date and time parsing/formatting
- **python-frontmatter** (≥1.0.0): Hugo markdown frontmatter parsing
- **colorlog** (≥6.7.0): Color-coded console logging
- **pytest** (≥7.4.0): Testing framework

### Step 4: Verify Installation

Run the setup verification script to ensure everything is installed correctly:

```bash
python test_setup.py
```

Expected output:
```
✓ Python version: 3.9.x
✓ All required packages installed
✓ Configuration system working
✓ Logging system working
✓ Wrangler CLI available
✓ Setup verification complete!
```

If you see any errors, refer to the [Troubleshooting](#troubleshooting) section below.

### Step 5: Configure the Tool

Copy the example configuration file:

```bash
cp config.yaml.example config.yaml
```

Edit `config.yaml` with your specific settings. See [CONFIGURATION.md](CONFIGURATION.md) for detailed configuration options.

**Minimum required configuration:**

```yaml
matrix:
  site_name: "your-site.com"

cloudflare:
  api_base: "https://comments.your-site.com/api"
  database_id: "your-database-id"

hugo:
  base_url: "https://your-site.com"
```

### Step 6: Test Configuration

Verify your configuration is valid:

```bash
python -c "from utils.config import Config; c = Config('config.yaml'); print('Configuration valid!')"
```

## Troubleshooting

### Python Version Issues

**Problem**: `python: command not found` or wrong version

**Solution**: 
- Try `python3` instead of `python`
- Install Python 3.8+ using your system's package manager
- On macOS: `brew install python@3.9`

### Virtual Environment Issues

**Problem**: `venv` module not found

**Solution**:
```bash
# On Ubuntu/Debian
sudo apt install python3-venv

# On macOS (if using system Python)
pip3 install virtualenv
python3 -m virtualenv venv
```

### Dependency Installation Failures

**Problem**: `pip install` fails with permission errors

**Solution**:
- Ensure your virtual environment is activated (you should see `(venv)` in prompt)
- If still failing, try: `pip install --user -r requirements.txt`

**Problem**: `pip install` fails with SSL errors

**Solution**:
```bash
pip install --trusted-host pypi.org --trusted-host files.pythonhosted.org -r requirements.txt
```

### Wrangler Issues

**Problem**: `wrangler: command not found`

**Solution**:
- Ensure Node.js is installed: `node --version`
- Reinstall Wrangler: `npm install -g wrangler`
- Add npm global bin to PATH: `export PATH="$(npm config get prefix)/bin:$PATH"`

**Problem**: Wrangler authentication fails

**Solution**:
- Clear existing auth: `wrangler logout`
- Re-authenticate: `wrangler login`
- Ensure you have access to the Cloudflare account with the D1 database

### Configuration Issues

**Problem**: `config.yaml` not found

**Solution**:
- Ensure you're in the `migrate_comments/` directory
- Copy the example: `cp config.yaml.example config.yaml`

**Problem**: YAML parsing errors

**Solution**:
- Check YAML syntax (indentation must be spaces, not tabs)
- Validate YAML online: https://www.yamllint.com/
- Ensure strings with special characters are quoted

## Uninstallation

To remove the migration tool:

1. Deactivate the virtual environment:
   ```bash
   deactivate
   ```

2. Remove the virtual environment:
   ```bash
   rm -rf venv
   ```

3. Optionally remove the entire `migrate_comments/` directory:
   ```bash
   cd ..
   rm -rf migrate_comments
   ```

## Next Steps

After successful installation:

1. Read [CONFIGURATION.md](CONFIGURATION.md) to configure the tool
2. Read [USAGE.md](USAGE.md) to learn how to run migrations
3. Run a dry-run migration to preview: `python migrate.py --dry-run`

## Getting Help

If you encounter issues not covered here:

1. Check [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for common problems
2. Review the logs in `migration.log`
3. Run with debug logging: Edit `config.yaml` and set `logging.level: DEBUG`
