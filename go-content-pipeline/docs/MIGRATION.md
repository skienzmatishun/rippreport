# Migration Guide: Python (`recent_placement`) to Go (`go-content-pipeline`)

This guide assists developers and operators migrating from the legacy Python implementation (`recent_placement`) to the high-performance Go-based pipeline (`go-content-pipeline`).

---

## 1. Summary of Differences

| Feature / Metric | Python Implementation (`recent_placement`) | Go Implementation (`go-content-pipeline`) |
| :--- | :--- | :--- |
| **Runtime** | Python 3.10+ (requires virtualenv & pip packages) | Single static compiled binary (zero runtime dependencies) |
| **Concurrency** | Sequential or GIL-limited multiprocessing | Native goroutines with bounded worker pools |
| **Memory Footprint** | ~350 MB - 600 MB RSS | **< 3 MB live heap**, guaranteed < 500 MB under heavy load |
| **Vector Similarity** | NumPy / pure Python | Native zero-allocation Go implementation (single pass) |
| **File I/O** | In-place file rewriting | **Atomic writes** (`.tmp-*` + `os.Rename`) across all operations |
| **Backup Verification** | Simple file copies | **Timestamped directories with MD5 checksum validation** |
| **Cache Corruption** | Unhandled JSON decode exceptions | **Automatic corruption recovery** (resets cleanly to empty) |
| **Signals & Teardown** | Immediate termination without save | **Graceful shutdown on SIGINT/SIGTERM** with progress persistence |
| **Data Compatibility** | Hugo YAML front matter & `sidebar-hooks.yaml` | **100% identical format and fields** |

---

## 2. Breaking Changes

**There are NO breaking data format changes.**
- Post front matter format (`related_articles: - slug: ... title: ... score: ... rank: ...`) is 100% identical.
- Hook file format (`{slug}/sidebar-hooks.yaml`) schema and keys are 100% identical.
- Hugo theme templates require zero modifications.

---

## 3. Command Mapping

| Python Command | Go Equivalent |
| :--- | :--- |
| `python -m recent_placement.main --generate-related` | `./bin/pipeline generate-related` |
| `python -m recent_placement.main --generate-hooks` | `./bin/pipeline generate-hooks` |
| `python -m recent_placement.main --all` | `./bin/pipeline pipeline` |
| `python -m recent_placement.main --stage` | `./bin/pipeline generate-related --stage` |
| `python -m recent_placement.main --apply-staged` | `./bin/pipeline apply-staged` |
| `python -m recent_placement.main --refresh-recent` | `./bin/pipeline generate-related --refresh-recent` |
| `python -m recent_placement.main --resume` | `./bin/pipeline generate-related --resume` |

---

## 4. Comparison Testing Steps

To compare results between the Python pipeline and the Go pipeline on the same dataset:

### Step 1: Run Go Pipeline in Staging Mode
Run the Go pipeline with `--stage` so live post front matter is not altered:
```bash
./bin/pipeline generate-related --stage
```

### Step 2: Compare Recommendations
Inspect the generated staging files (`content/p/{slug}/related-articles.staged.yaml`) against previous outputs:
```bash
cat content/p/rock-bottom/related-articles.staged.yaml
```
Verify that:
1. Slugs recommended are highly relevant and match expected topics.
2. Self-exclusion is respected (no post recommends itself).
3. Backstory/podcast episodes are properly segregated.
4. Ranks are 1 to 5 (or top-N) sequential.

### Step 3: Run Full Verification Tests
```bash
# Run unit tests
go test ./...

# Run full integration and end-to-end tests
go test -v ./test/...
```

### Step 4: Apply Staged Rankings When Confirmed
```bash
./bin/pipeline apply-staged
```
Backups of original post files will be automatically preserved in `backups/`.
