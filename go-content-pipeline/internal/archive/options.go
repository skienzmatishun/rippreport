// Package archive verifies audit reports and repairs links in Hugo bundles.
package archive

import (
	"fmt"
	"net/url"
	"time"
)

// Options keeps archive repair independent of recommendation model settings.
type Options struct {
	Report, Progress, ContentDir, BackupDir, ScreenshotDir                                     string
	Python, ModelURL, Model, CDXURL, SearchURL                                                 string
	DryRun, UsePlaywright, UseLMStudio, Vision, Google, VerifyGoogle                           bool
	FixCDN, ConvertTags, FixMissingPath, VerifyPaywall, CleanupScreenshots                     bool
	RequestTimeout, BrowserTimeout, ModelTimeout, Wait, SearchDelay, SearchJitter, LockTimeout time.Duration
	MaxSnapshots, MaxResults, MaxContent, Retries                                              int
}

func DefaultOptions() Options {
	return Options{
		Report: "audit/archive_fix_report.json", Progress: "audit/archive_verifier_progress.json",
		ContentDir: "content/p", BackupDir: "backups", ScreenshotDir: "temp_screenshots", Python: "python3",
		ModelURL: "http://localhost:1234/v1", Model: "google/gemma-3-12b",
		CDXURL: "https://web.archive.org/cdx/search/cdx", SearchURL: "https://www.google.com/search",
		UsePlaywright: true, UseLMStudio: true, Vision: true, Google: true, VerifyGoogle: true,
		FixCDN: true, ConvertTags: true, FixMissingPath: true, VerifyPaywall: true, CleanupScreenshots: true,
		RequestTimeout: 10 * time.Second, BrowserTimeout: 60 * time.Second, ModelTimeout: 30 * time.Second,
		Wait: 5 * time.Second, SearchDelay: 25 * time.Second, SearchJitter: 10 * time.Second, LockTimeout: 30 * time.Second,
		MaxSnapshots: 15, MaxResults: 5, MaxContent: 50000, Retries: 3,
	}
}

func (o Options) Validate() error {
	if o.Report == "" || o.Progress == "" || o.ContentDir == "" || o.BackupDir == "" || o.ScreenshotDir == "" {
		return fmt.Errorf("report, progress, content, backup and screenshot paths must not be empty")
	}
	if o.RequestTimeout <= 0 || o.BrowserTimeout <= 0 || o.ModelTimeout <= 0 || o.LockTimeout <= 0 || o.Wait < 0 || o.SearchDelay < 0 || o.SearchJitter < 0 {
		return fmt.Errorf("timeouts must be positive; waits and search delays must be nonnegative")
	}
	if o.MaxSnapshots < 1 || o.MaxResults < 1 || o.MaxContent < 1 || o.Retries < 1 {
		return fmt.Errorf("limits and retry count must be positive")
	}
	for _, raw := range []string{o.ModelURL, o.CDXURL, o.SearchURL} {
		if !webURL(raw) {
			return fmt.Errorf("invalid HTTP(S) endpoint: %s", raw)
		}
	}
	return nil
}

func webURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Host != "" && (u.Scheme == "https" || u.Scheme == "http")
}

// Entry preserves unknown audit fields when reports are updated.
type Entry map[string]any

func (e Entry) str(k string) string { s, _ := e[k].(string); return s }
func (e Entry) key() string         { return e.str("post_url") + "\x00" + e.str("original_href") }

type Verification struct {
	Valid      bool           `json:"is_valid"`
	Method     string         `json:"method"`
	Reason     string         `json:"reason"`
	Confidence float64        `json:"confidence,omitempty"`
	VerifiedAt string         `json:"verified_at"`
	Details    map[string]any `json:"details,omitempty"`
}

type Summary struct{ Total, Replaced, Verified, Failed, Skipped, WouldReplace int }

func Strategy(e Entry, o Options) string {
	original := e.str("original_href")
	u, _ := url.Parse(original)
	if e.str("status") == "fixed" && e.str("archive_url") != "" {
		return "verify_fixed"
	}
	if u == nil {
		return "skip"
	}
	if isPDF(original) && e.str("archive_url") == "" {
		if o.Google {
			return "pdf_fallback"
		}
		return "skip"
	}
	if o.FixCDN && u.Hostname() == "cdn.rippreport.com" && !hasExtension(u.Path) {
		return "malformed_cdn"
	}
	if siteHost(u.Hostname()) {
		if o.ConvertTags && tagName(u.Path) != "" {
			return "tag_conversion"
		}
		if o.FixMissingPath && missingSlug(u.Path) != "" {
			return "missing_p_path"
		}
	}
	if o.VerifyPaywall && e.str("archive_url") == "" {
		return "paywall_check"
	}
	return "skip"
}
