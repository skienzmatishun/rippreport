# Link and migration audits

The workflows in `audit/link_audit` are available directly in the Go pipeline:

| Original script | Pipeline command |
| --- | --- |
| `link_checker.py` | `check-links` (alias `link-checker`) |
| `analyze_broken_links.py` | `analyze-broken-links` |
| `archive_links.py` | `archive-links` |
| `update_redirected_links.py` | `update-redirected-links` |
| `content_audit.py` | `audit-content` (alias `content-audit`) |

These commands do not need the original Python scripts. Optional Playwright
rendering for archive verification uses the embedded adapter described in the
[archive verifier guide](ARCHIVE_VERIFIER.md).

Build and run from `go-content-pipeline/`:

```sh
go build -o bin/pipeline ./cmd/pipeline
./bin/pipeline check-links --help
```

Default content/report paths resolve from either the site root or the Go module
directory. Explicit paths are relative to the current working directory. Audit
commands are independent of the recommendation pipeline's model configuration.

## Check rendered post links

Start Hugo separately, then run:

```sh
./bin/pipeline check-links --base-url http://localhost:1313
./bin/pipeline check-links --posts sample-post,another-post
./bin/pipeline check-links --resume=false --posts sample-post
```

The checker scans bundle directories with `index.md`, fetches their rendered
`/p/{slug}/` pages, extracts anchor destinations and labels, resolves relative
URLs against the fetched post URL, and checks HTTP(S) destinations with GET
requests. Fragment-only links and non-HTTP schemes such as `mailto:`, `tel:`,
and `javascript:` are excluded. Link labels are limited to 100 characters.

Defaults match the Python checker:

- Site: `http://localhost:1313`; content: `content/p`.
- Report: `audit/broken_links_report.json`.
- HTTP timeout: 10 seconds; maximum redirects: 10.
- Five post workers, with at most 20 outgoing link checks globally.
- A 500 ms delay after each link request per worker.
- Excluded domains: Google, Facebook/fb.com, Twitter/X, localhost, sibforms.com,
  and their subdomains. Privacy-policy path suffixes are also excluded.

Configure these with `--base-url`, `--content-dir`, `--report`,
`--request-timeout`, `--max-redirects`, `--post-workers`, `--link-workers`,
`--request-delay`, `--skip-domains`, and `--skip-paths`. Domains and paths are
comma-separated; an empty value disables that exclusion list. Domain matching
uses host boundaries, so excluding `facebook.com` does not exclude
`facebook.com.example.org`. The default localhost exclusion includes relative
links resolved to the local Hugo server; use `--skip-domains` without localhost
to audit those links too.

The report retains the original `posts`, `links`, and `summary` schema, including
`original_href`, final URL, redirect count, HTTP status, transport errors, and
`is_broken`. A link can have `status: "ok"` and HTTP 404; `is_broken` records the
failure, matching the original script. Inaccessible posts are recorded separately
and counted in `summary.inaccessible_posts`.

Progress is saved atomically after each complete post. Reruns skip recorded
posts, including inaccessible ones; `--resume=false` rechecks the selected
bundles and replaces their results. Cancelled partial posts are not recorded as
completed. Results are sorted by post URL on completion. Rechecking produces
the link-audit schema and removes derived `entries` added by a previous
`verify-archives` run; use a separate archive repair report for that workflow.

## Analyze failures

```sh
./bin/pipeline analyze-broken-links
./bin/pipeline analyze-broken-links --output /tmp/broken-link-analysis.json
```

The command recalculates counts from the report, prints success rate and
inaccessible-post totals, groups broken links by domain, HTTP status, and error
type, and lists the ten posts with the most failures. Console groups show
bounded examples; optional JSON retains all references. `--report` selects a
different input. This command performs no network requests.

## Discover and install Wayback archives

```sh
./bin/pipeline archive-links --dry-run
./bin/pipeline archive-links
./bin/pipeline archive-links --resume=false --max-snapshots 30
```

This is the discovery stage before `verify-archives`: it reads broken links from
accessible posts, finds successful Wayback CDX captures, tries them newest first,
and installs the first verified archive. It records the post date/context when
YAML front matter can be parsed, every attempted snapshot, the selected URL,
verification result, and failure reason.

Default outputs are `audit/archive_fix_report.json` and
`audit/archive_fix_report.csv`. Change these with `--output` and `--csv`; pass
`--csv=` to disable CSV. JSON uses `entries`, accepted by `verify-archives`;
legacy `processed` reports can be resumed. Existing attempts are skipped unless
`--resume=false` is set. A successful repair is marked `fixed` for subsequent
archive re-verification.

HTTP HTML verification is enabled by default. `--playwright` enables rendered
verification; `--python` selects its interpreter. `--lmstudio`, `--lmstudio-url`,
and `--model` enable optional text/vision analysis using the existing archive
verifier. `--request-timeout`, `--cdx-url`, `--max-snapshots`, `--lock-timeout`,
`--content-dir`, `--report`, and `--backup-dir` control inputs and limits.

Every changed file receives a checksum-verified original backup, retaining
bundle structure and a restoration manifest. Reports are saved after each
attempt. Dry runs write no content, backups, or output reports. Stable sidecar
lock files may be created while a command holds a report lock.

## Update redirected URLs

```sh
./bin/pipeline update-redirected-links
./bin/pipeline update-redirected-links --domain lagniappe --min-redirects 2
./bin/pipeline update-redirected-links --update
```

The default previews changes. `--update` applies them; `--dry-run` always takes
precedence. Only successful, unbroken links with a final HTTP(S) URL and at least
`--min-redirects` hops are eligible. Trailing-slash-only differences are skipped.
`--domain` performs the original script's case-insensitive substring filter on
the original hostname.

Repairs use `original_href`, supporting relative links in Markdown. Exact URL
boundaries preserve link titles and avoid changing destinations with a shared
prefix. The original script's case-insensitive fallback is available by default;
disable it with `--case-insensitive=false`. Existing installed replacements are
recognized so interrupted work can be retried.

Live changes are backed up before writing. The output defaults to
`audit/redirected_links_update_report.json` and records filters, counts, modified
posts, individual updates/errors, and backup paths. `--output`, `--report`,
`--content-dir`, `--backup-dir`, and `--lock-timeout` customize these. Preview mode
prints its results and does not write a repair report.

## Audit WordPress migration shortcodes

```sh
./bin/pipeline audit-content
./bin/pipeline audit-content --no-model --request-delay=0s
./bin/pipeline audit-content --posts sample-post --resume=false
```

The default uses the pipeline's `config.yaml` settings for `llama_server.base_url`,
`llama_server.llm_model`, and managed process settings. In managed mode the audit
starts llama-server with the completion model and stops it on completion or
cancellation, just like the other pipeline commands. Use `--config` to select a
configuration, `--model` to override the GGUF path/name, or `--managed=false`
to connect to an existing server. `--llama-server-url` selects an existing server
and disables managed startup unless `--managed` is explicitly supplied.
The audit sends up to 8,000 characters per post, with a 600-second model timeout
and a two-second delay between posts. The conservative prompt only
flags WordPress `[gallery]`, `[embed]`, `[caption]`, `[video]`, `[audio]`, and
`[playlist]` shortcodes, including attributes/closing tags. It excludes ordinary
Markdown links, HTML, Hugo shortcodes, front matter, and code examples. Thinking
blocks and JSON fences in model responses are handled.

`--no-model` performs local detection without loading configuration or starting a model (`--no-lmstudio` remains a compatibility alias). It excludes front matter, fenced and
indented code, inline code, HTML code/pre blocks, Hugo shortcode bodies, escaped
brackets, and ordinary Markdown links/reference definitions.

The resumable report keeps the original JSON array format:

```json
[{"post":"sample-post","has_issue":true,"issues":["WordPress shortcode [gallery] at line 5"]}]
```

Its default path is `hugo_post_audit_report.json` in the site root. Flags include
`--report`, `--content-dir`, `--posts`, `--resume`, `--config`, `--managed`, `--llama-server-url`, `--model`,
`--model-timeout`, `--model-retries`, `--max-content-length`, and
`--request-delay`. Each completed post is saved atomically. Model/parse errors
are recorded explicitly and return a failing exit status. Failed records are
also skipped on resume; use `--resume=false` to retry them. Potential shortcode
issues alone do not make command execution fail.

## Combined workflow

```sh
./bin/pipeline check-links
./bin/pipeline analyze-broken-links
./bin/pipeline update-redirected-links --update
./bin/pipeline archive-links
./bin/pipeline verify-archives --report ../audit/archive_fix_report.json --dry-run
./bin/pipeline missing-pdfs --report ../audit/archive_fix_report.json --output ../audit/missing_pdfs.json
```

The PDF export is based on recorded results, not a fresh network check. A full
list of unresolved audited PDFs can also be exported directly from
`broken_links_report.json`.

Input JSON is validated before updates; malformed reports are left intact.
Output paths must differ from source reports. Source content paths use the
archive package's traversal/symlink confinement checks. Ctrl+C preserves
completed work and returns 130; fatal errors or failed repairs return 1.
Successfully completed audits and help return 0, even when an audit discovers
broken links or potential migration issues.

## Tests

```sh
go test -race ./internal/linkaudit ./internal/archive
go test ./internal/cli -run 'TestLinkAudit|TestArchive|TestMissingPDFs'
```

Tests use temporary bundles and local HTTP/model/CDX fixtures, covering link
extraction, exclusions, redirect limits, timeouts, concurrency, resume/recheck,
analysis, repairs, backups/restoration, CSV, model JSON/thinking parsing, and
offline shortcode detection. They do not depend on a running Hugo or model
server, and do not crawl production URLs.
