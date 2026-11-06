package archive

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func atomicWrite(name string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(name), ".archive-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if err := f.Chmod(mode); err != nil {
		f.Close()
		return err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), name)
}

func saveJSON(name string, data any) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return atomicWrite(name, append(b, '\n'), 0644)
}

// Lock a stable sidecar, not the atomically replaced data inode. The kernel
// releases the lock if the process exits, including after an interruption.
func acquire(ctx context.Context, name string, timeout time.Duration) (func(), error) {
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	deadline := time.Now().Add(timeout)
	for {
		err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			return func() { syscall.Flock(int(f.Fd()), syscall.LOCK_UN); f.Close() }, nil
		}
		if err != syscall.EWOULDBLOCK && err != syscall.EAGAIN {
			f.Close()
			return nil, err
		}
		if time.Now().After(deadline) {
			f.Close()
			return nil, fmt.Errorf("timed out waiting for lock %s", name)
		}
		if err := pause(ctx, 100*time.Millisecond); err != nil {
			f.Close()
			return nil, err
		}
	}
}

func readObject(name string) (map[string]any, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var data map[string]any
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s: %w", name, err)
	}
	if data == nil {
		return nil, fmt.Errorf("%s must contain a JSON object", name)
	}
	return data, nil
}

// LoadReport accepts archive-fix reports (including legacy 'processed') and
// link-audit reports. Existing metadata and fields survive updates.
func LoadReport(name string) (map[string]any, []Entry, error) {
	data, err := readObject(name)
	if err != nil {
		return nil, nil, err
	}
	var rows []any
	if raw, ok := data["entries"]; ok {
		rows, ok = raw.([]any)
		if !ok {
			return nil, nil, fmt.Errorf("report entries must be an array")
		}
	} else if raw, ok := data["processed"]; ok {
		rows, ok = raw.([]any)
		if !ok {
			return nil, nil, fmt.Errorf("report processed must be an array")
		}
		delete(data, "processed")
	} else if raw, ok := data["posts"]; ok {
		posts, ok := raw.([]any)
		if !ok {
			return nil, nil, fmt.Errorf("report posts must be an array")
		}
		for _, p := range posts {
			post, ok := p.(map[string]any)
			if !ok {
				return nil, nil, fmt.Errorf("invalid post in report")
			}
			links, _ := post["links"].([]any)
			for _, l := range links {
				link, ok := l.(map[string]any)
				if !ok {
					return nil, nil, fmt.Errorf("invalid link in report")
				}
				broken, _ := link["is_broken"].(bool)
				if !broken {
					continue
				}
				e := map[string]any{}
				for k, value := range link {
					e[k] = value
				}
				e["post_url"] = post["post_url"]
				if s, _ := e["original_href"].(string); s == "" {
					e["original_href"] = e["url"]
				}
				if _, ok := e["link_text"]; !ok {
					e["link_text"] = e["text"]
				}
				for _, k := range []string{"post_title", "post_date", "post_excerpt"} {
					if value, ok := post[k]; ok {
						e[k] = value
					}
				}
				rows = append(rows, e)
			}
		}
	} else {
		return nil, nil, fmt.Errorf("report must contain entries, processed, or posts")
	}
	entries := make([]Entry, 0, len(rows))
	for i, row := range rows {
		e, ok := row.(map[string]any)
		if !ok {
			return nil, nil, fmt.Errorf("entry %d must be an object", i)
		}
		entries = append(entries, Entry(e))
	}
	data["entries"] = entries
	return data, entries, nil
}

// PostPath rejects traversal and symlinks that would escape the content root.
func PostPath(root, postURL string) (string, error) {
	u, err := url.Parse(postURL)
	if err != nil {
		return "", err
	}
	p := strings.Trim(u.Path, "/")
	if !strings.HasPrefix(p, "p/") {
		return "", fmt.Errorf("post URL must start with /p/")
	}
	slug := strings.TrimPrefix(p, "p/")
	for _, part := range strings.Split(slug, "/") {
		if part == "" || part == "." || part == ".." || strings.Contains(part, "\\") {
			return "", fmt.Errorf("invalid post slug")
		}
	}
	base, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}
	base, err = filepath.Abs(base)
	if err != nil {
		return "", err
	}
	file, err := filepath.EvalSymlinks(filepath.Join(base, filepath.FromSlash(slug), "index.md"))
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(base, file)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("post escapes content directory")
	}
	info, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("post is not a regular file")
	}
	return file, nil
}

// BackupManager keeps one verified original per file per run and supports
// restoring individual files or an entire run via its manifest.
type BackupManager struct {
	Base, Dir string
	Files     map[string]string
}

func (b *BackupManager) Backup(file, root string) (string, error) {
	if b.Files == nil {
		b.Files = map[string]string{}
	}
	if existing, ok := b.Files[file]; ok {
		return existing, nil
	}
	if b.Dir == "" {
		if err := os.MkdirAll(b.Base, 0755); err != nil {
			return "", err
		}
		dir, err := os.MkdirTemp(b.Base, "archive-links-backup-"+time.Now().Format("20060102-150405")+"-")
		if err != nil {
			return "", err
		}
		b.Dir = dir
		b.Dir, err = filepath.Abs(b.Dir)
		if err != nil {
			return "", err
		}
	}
	base, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}
	base, err = filepath.Abs(base)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(base, file)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("backup source escapes content root")
	}
	destination := filepath.Join(b.Dir, "content", "p", rel)
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	if err := atomicWrite(destination, data, info.Mode().Perm()); err != nil {
		return "", err
	}
	check, err := os.ReadFile(destination)
	if err != nil {
		return "", err
	}
	if sha256.Sum256(data) != sha256.Sum256(check) {
		return "", fmt.Errorf("backup checksum mismatch")
	}
	b.Files[file] = destination
	if err := saveJSON(filepath.Join(b.Dir, "manifest.json"), b.Files); err != nil {
		delete(b.Files, file)
		return "", err
	}
	return destination, nil
}

func RestoreBackup(manifest string) error {
	b, err := os.ReadFile(manifest)
	if err != nil {
		return err
	}
	var files map[string]string
	if err := json.Unmarshal(b, &files); err != nil {
		return err
	}
	for original, backup := range files {
		data, err := os.ReadFile(backup)
		if err != nil {
			return err
		}
		info, err := os.Stat(backup)
		if err != nil {
			return err
		}
		if err := atomicWrite(original, data, info.Mode().Perm()); err != nil {
			return err
		}
	}
	return nil
}

func replace(o Options, e Entry, target string, backups *BackupManager) (string, error) {
	file, err := PostPath(o.ContentDir, e.str("post_url"))
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	updated, count := ReplaceURLs(string(data), e.str("original_href"), target)
	// Audit fixes may already have installed an invalid archive. Replace that
	// exact destination when choosing a verified earlier snapshot.
	if count == 0 && e.str("archive_url") != "" {
		updated, count = ReplaceURLs(string(data), e.str("archive_url"), target)
	}
	if count == 0 {
		_, present := ReplaceURLs(string(data), target, "https://archive-verifier.invalid/present")
		if present > 0 {
			return "", nil
		}
		return "", fmt.Errorf("neither original nor replacement URL found in post")
	}
	if o.DryRun {
		return "", nil
	}
	backup, err := backups.Backup(file, o.ContentDir)
	if err != nil {
		return "", fmt.Errorf("backup failed: %w", err)
	}
	info, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	current, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	if !bytes.Equal(current, data) {
		return "", fmt.Errorf("post changed during verification; retry")
	}
	if err := atomicWrite(file, []byte(updated), info.Mode().Perm()); err != nil {
		return backup, err
	}
	return backup, nil
}

func progressStats(entries []Entry) map[string]int {
	stats := map[string]int{"total_processed": len(entries), "successfully_replaced": 0, "verified": 0, "verification_failed": 0, "replacement_failed": 0, "skipped": 0, "failed": 0}
	for _, e := range entries {
		key := e.str("status")
		if key == "replaced" {
			key = "successfully_replaced"
		}
		stats[key]++
	}
	return stats
}

// Run saves the report before progress, making interruption retries idempotent.
// Dry runs do not modify reports, progress, backups, or content.
func (v *Verifier) Run(ctx context.Context) (Summary, error) {
	o := v.Options
	var summary Summary
	if err := o.Validate(); err != nil {
		return summary, err
	}
	if !o.DryRun {
		unlock, err := acquire(ctx, o.Report+".lock", o.LockTimeout)
		if err != nil {
			return summary, err
		}
		defer unlock()
		unlockProgress, err := acquire(ctx, o.Progress+".lock", o.LockTimeout)
		if err != nil {
			return summary, err
		}
		defer unlockProgress()
	}
	report, entries, err := LoadReport(o.Report)
	if err != nil {
		return summary, err
	}
	summary.Total = len(entries)
	progress := map[string]any{}
	processed := []Entry{}
	if p, err := readObject(o.Progress); err == nil {
		progress = p
		rows, ok := p["processed_links"].([]any)
		if !ok {
			return summary, fmt.Errorf("progress missing processed_links array")
		}
		for _, raw := range rows {
			e, ok := raw.(map[string]any)
			if !ok {
				return summary, fmt.Errorf("invalid progress entry")
			}
			processed = append(processed, Entry(e))
		}
	} else if !os.IsNotExist(err) {
		return summary, err
	}
	done := map[string]bool{}
	for _, e := range processed {
		done[e.key()] = true
	}
	backups := &BackupManager{Base: o.BackupDir}
	for i, e := range entries {
		if err := ctx.Err(); err != nil {
			return summary, err
		}
		if done[e.key()] || e.str("status") == "replaced" || e.str("status") == "verified" {
			summary.Skipped++
			continue
		}
		fmt.Fprintf(v.Log, "[%d/%d] %s: %s\n", i+1, len(entries), Strategy(e, o), e.str("original_href"))
		status := "skipped"
		reason := "No applicable strategy"
		var verification Verification
		var target, backup string
		if Strategy(e, o) == "skip" {
			summary.Skipped++
		} else if e.str("post_url") == "" || !webURL(e.str("original_href")) {
			status = "failed"
			reason = "Missing post_url or invalid original_href"
			summary.Failed++
		} else {
			target, verification = v.Process(ctx, e)
			if err := ctx.Err(); err != nil {
				return summary, err
			}
			reason = verification.Reason
			status = "verification_failed"
			if verification.Valid {
				backup, err = replace(o, e, target, backups)
				if err != nil {
					status = "replacement_failed"
					reason = err.Error()
					summary.Failed++
				} else if target == e.str("original_href") {
					status = "verified"
					summary.Verified++
				} else if o.DryRun {
					status = "would_replace"
					summary.WouldReplace++
				} else {
					status = "replaced"
					summary.Replaced++
				}
			} else {
				summary.Failed++
			}
		}
		fmt.Fprintf(v.Log, "  %s: %s\n", status, reason)
		if o.DryRun {
			continue
		}
		now := time.Now().UTC().Format(time.RFC3339)
		e["status"] = status
		e["reason"] = reason
		e["verification"] = verification
		if target != "" && (status == "replaced" || status == "verified") {
			e["archive_url"] = target
			e["replacement"] = map[string]any{"success": true, "backup_path": backup, "replaced_at": now}
		}
		processed = append(processed, Entry{"post_url": e.str("post_url"), "original_href": e.str("original_href"), "status": status, "timestamp": now, "new_url": target, "reason": reason, "backup_path": backup})
		progress["processed_links"] = processed
		progress["statistics"] = progressStats(processed)
		progress["last_updated"] = now
		metadata, ok := report["metadata"].(map[string]any)
		if !ok {
			metadata = map[string]any{}
			report["metadata"] = metadata
		}
		metadata["last_updated"] = now
		if err := saveJSON(o.Report, report); err != nil {
			return summary, err
		}
		if err := saveJSON(o.Progress, progress); err != nil {
			return summary, err
		}
		done[e.key()] = true
	}
	return summary, nil
}

// ExportMissingPDFs groups unresolved PDF URLs and their referring posts using
// the recorded audit results; it does not recheck URLs over the network.
func ExportMissingPDFs(reportPath, output string) (int, error) {
	_, entries, err := LoadReport(reportPath)
	if err != nil {
		return 0, err
	}
	byURL := map[string]map[string]any{}
	var pdfs []map[string]any
	occurrences := 0
	for _, e := range entries {
		raw := e.str("original_href")
		if !isPDF(raw) || e.str("status") == "replaced" || e.str("status") == "verified" {
			continue
		}
		if e.str("status") == "fixed" && e.str("archive_url") != "" {
			continue
		}
		item, ok := byURL[raw]
		if !ok {
			item = map[string]any{"url": raw, "references": []Entry{}}
			byURL[raw] = item
			pdfs = append(pdfs, item)
		}
		ref := Entry{}
		for _, k := range []string{"post_url", "link_text", "status", "is_broken", "status_code", "final_url", "error", "reason", "archive_url"} {
			if value, ok := e[k]; ok {
				ref[k] = value
			}
		}
		item["references"] = append(item["references"].([]Entry), ref)
		occurrences++
	}
	if pdfs == nil {
		pdfs = []map[string]any{}
	}
	err = saveJSON(output, map[string]any{"source_report": reportPath, "generated_at": time.Now().UTC().Format(time.RFC3339), "description": "Unresolved PDF links from recorded audit results; URLs have not been rechecked.", "total_missing_pdfs": len(pdfs), "total_references": occurrences, "pdfs": pdfs})
	return len(pdfs), err
}
