# Archive link verification and missing PDFs

`pipeline verify-archives` integrates the functionality of
`audit/archive_link_verifier` into the Go pipeline. Verification, strategy
selection, HTTP requests, model calls, repairs, backups, reports, and progress
tracking run in Go. An embedded Python Playwright adapter renders pages and
captures screenshots. The binary does not need the original Python package.

## Usage

From `go-content-pipeline/`:

```sh
go build -o bin/pipeline ./cmd/pipeline
./bin/pipeline missing-pdfs --output ../audit/missing_pdfs.json
./bin/pipeline verify-archives --dry-run --no-playwright --no-lmstudio --no-google
./bin/pipeline verify-archives --missing-pdfs-output ../audit/missing_pdfs.json
./bin/pipeline verify-archives --help
```

Defaults resolve from either the site root or the Go module directory. The
command prefers `audit/archive_fix_report.json` and uses
`audit/broken_links_report.json` if the former is absent. Explicit paths are
relative to the working directory. `archive-link-verifier` is an alias.

`missing-pdfs` exports unresolved PDF destinations with their referring posts,
link labels, status codes, errors, and available repair metadata. Destinations
are deduplicated while all references are retained. `.PDF` and PDF URLs with
query strings are recognized. The export uses recorded audit results and does
not claim that a URL remains missing today. Broken audit entries may include
403s, timeouts, or other access failures as well as 404s. Entries with completed
repairs, verified originals, or existing successful archive fixes are excluded.

## Verification and repair strategies

| Strategy | Behavior |
| --- | --- |
| Existing archive | Verify `fixed` archive URLs, then try up to 15 earlier successful CDX captures, newest first. Ignore Wayback chrome when examining text. |
| Missing PDF | Search for the original PDF URL; verify successful PDF content types and, when available, screenshots. |
| Malformed CDN | Try a URL supplied as link text; optionally use model context to suggest and verify a destination. |
| Tag link | Convert a Ripp Report tag link to a Google site search. |
| Missing `/p/` | Use an existing local post bundle, or search for and verify a canonical Ripp Report URL. |
| Paywall | Check rendered or fetched content for an accessible article, subscription gates, and relevance to available post context. Keep accessible original URLs. |

Strategies follow the original verifier's priority order. Set `--no-google` to
disable network searches. Local bundle repairs and tag conversion remain
available. Other switches include `--fix-malformed-cdn=false`,
`--convert-tags=false`, `--fix-missing-p-path=false`, and
`--verify-paywall=false`.

## Browser and model setup

HTTP-only mode (`--no-playwright`) needs no Python packages. For rendered pages:

```sh
python3 -m venv /tmp/archive-verifier-venv
/tmp/archive-verifier-venv/bin/python -m pip install playwright
/tmp/archive-verifier-venv/bin/python -m playwright install chromium
./bin/pipeline verify-archives --python /tmp/archive-verifier-venv/bin/python --dry-run
```

Browser failures fall back to HTTP with a diagnostic. Browser navigation times
out after 60 seconds and waits another 5 seconds for dynamic content. Configure
these using `--browser-timeout` and `--browser-wait` (Go durations, such as `30s`).

The default model endpoint is `http://localhost:1234/v1`, with model
`google/gemma-3-12b`. Override `--lmstudio-url` and `--model` for LM Studio,
Unsloth, or another OpenAI-compatible server. Requests support text and PNG
screenshots, JSON extraction from fenced responses, bounded context, timeout,
and exponential retry backoff. `--no-lmstudio` uses HTML/content-type checks;
`--vision=false` uses text-only model analysis. Model failures fall back to
HTML analysis. Results record which method actually succeeded.

Screenshots are deleted after analysis by default. Use
`--cleanup-screenshots=false` to keep them and `--screenshot-dir` to choose a
location. Dry runs may render pages and create temporary screenshots, but never
write content, backups, reports, or progress.

Searches are separated by 25 seconds plus up to 10 seconds of jitter, matching
the previous verifier's conservative pacing. Configure `--search-delay`,
`--search-jitter`, and `--max-results` as needed. `--request-timeout`,
`--max-snapshots`, `--lmstudio-timeout`, `--lmstudio-retries`, and
`--max-content-length` control the remaining limits.

## Reports, backups, and resume

Accepted inputs are `{ "entries": [...] }`, legacy `{ "processed": [...] }`,
and link-audit `{ "posts": [{ "post_url": "...", "links": [...] }] }` reports.
Only links marked `is_broken` are imported from the last format. Importing a
link-audit report adds an `entries` array; original `posts` data is retained as
the audit snapshot. Subsequent runs use the repair entries. Unknown fields are
preserved. Missing or malformed input fails without overwriting it.

Each entry records verification metadata and one of `replaced`, `verified`,
`verification_failed`, `replacement_failed`, `skipped`, or `failed`. Accessible
originals receive `verified`, avoiding unnecessary writes. Existing invalid
archive destinations can be replaced with verified earlier captures. Exact URL
replacement preserves Markdown titles, HTML links, reference links, and bare
URLs while avoiding changes to URLs with a shared prefix.

Every changed content file receives a checksum-verified backup before writing.
One original per file is kept per run, including when multiple links change.
Backups include a manifest mapping original files to their saved copies:

```sh
./bin/pipeline verify-archives --restore-backup ../backups/archive-links-backup-TIMESTAMP-SUFFIX/manifest.json
```

Manifest paths should be kept in their original location for restoration.
Content paths are confined to the selected content directory, including checks
for symlink escapes. Content and JSON writes use temporary files and atomic
renames. Report/progress sidecar locks prevent competing runs using the same
files. Interrupted runs retain completed work, with report saved before
progress so retries can recognize installed replacements.

Rerunning the command resumes automatically using `(post_url, original_href)`.
Like the original, progress tracks unsuccessful attempts as processed too.
To retry selected failures, remove their records from `processed_links` or use
a new `--progress` path. To retry completed repairs, also reset their report
statuses. Failed backups prevent content changes.

`--log-file` adds a persistent execution log. `--log-level` accepts DEBUG, INFO,
WARNING, or ERROR; WARNING/ERROR suppress routine console progress. The final
summary is always printed. A fatal error or any failed entry returns exit code
1, help and successful runs return 0, and cancellation returns 130.

Tests use temporary bundles and local HTTP/model/search servers:

```sh
go test ./internal/archive
go test ./internal/cli -run 'TestArchive|TestMissingPDFs'
```
