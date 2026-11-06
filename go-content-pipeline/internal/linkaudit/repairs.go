package linkaudit

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/archive"
	"github.com/rippreport/go-content-pipeline/internal/parser"
)

type RepairOptions struct {
	Input, Output, CSV      string
	Archive                 archive.Options
	MinRedirects            int
	Domain                  string
	Resume, CaseInsensitive bool
}

type RepairSummary struct {
	Checked       int `json:"total_links_checked"`
	Candidates    int `json:"candidates_found"`
	Filtered      int `json:"filtered_out"`
	SameURL       int `json:"skipped_same_url"`
	Updated       int `json:"successfully_updated"`
	Failed        int `json:"update_failed"`
	Skipped       int `json:"already_processed"`
	PostsModified int `json:"posts_modified"`
}

func key(e archive.Entry) string {
	return stringField(e, "post_url") + "\x00" + stringField(e, "original_href")
}
func stringField(e archive.Entry, k string) string { s, _ := e[k].(string); return s }

func writeCSV(path string, entries []archive.Entry) error {
	if path == "" {
		return nil
	}
	fields := []string{"post_url", "post_date", "original_href", "archive_url", "status", "reason", "attempted_snapshots", "verification", "backup_path"}
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.Write(fields); err != nil {
		return err
	}
	for _, e := range entries {
		values := make([]string, len(fields))
		for i, k := range fields {
			if e[k] == nil {
				continue
			}
			if s, ok := e[k].(string); ok {
				values[i] = s
			} else {
				b, _ := json.Marshal(e[k])
				values[i] = string(b)
			}
		}
		if err := w.Write(values); err != nil {
			return err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	return archive.WriteFile(path, buf.Bytes(), 0644)
}

func postContext(root, postURL string) (archive.Entry, error) {
	file, err := archive.PostPath(root, postURL)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	context := archive.Entry{}
	if fm, body, _, err := parser.ParseFrontMatter(string(data)); err == nil {
		context["post_title"] = fm.Title
		if !fm.Date.IsZero() {
			context["post_date"] = fm.Date.Format(time.RFC3339)
		}
		excerpt := []rune(body)
		if len(excerpt) > 1000 {
			excerpt = excerpt[:1000]
		}
		context["post_excerpt"] = string(excerpt)
	}
	return context, nil
}

func validateRepair(o RepairOptions) error {
	if o.Input == "" || o.Output == "" {
		return fmt.Errorf("input and output reports are required")
	}
	if err := o.Archive.Validate(); err != nil {
		return err
	}
	return distinct(o.Input, o.Output, o.CSV)
}

// ArchiveBroken discovers archives for broken links without an existing archive
// URL, filling the gap before the verify-archives workflow.
func ArchiveBroken(ctx context.Context, o RepairOptions, log io.Writer) (RepairSummary, error) {
	var stats RepairSummary
	if log == nil {
		log = io.Discard
	}
	if err := validateRepair(o); err != nil {
		return stats, err
	}
	unlock, err := archive.Lock(ctx, o.Output+".lock", o.Archive.LockTimeout)
	if err != nil {
		return stats, err
	}
	defer unlock()
	input, err := LoadReport(o.Input)
	if err != nil {
		return stats, err
	}
	report := map[string]any{"entries": []archive.Entry{}}
	entries := []archive.Entry{}
	if _, err := os.Stat(o.Output); err == nil {
		loaded, existing, err := archive.LoadReport(o.Output)
		if err != nil {
			return stats, err
		}
		report, entries = loaded, existing
	} else if !os.IsNotExist(err) {
		return stats, err
	}
	indices := map[string]int{}
	for i, e := range entries {
		indices[key(e)] = i
	}
	v := archive.New(o.Archive, log)
	backups := &archive.BackupManager{Base: o.Archive.BackupDir}
	modified := map[string]bool{}
	save := func() error {
		if o.Archive.DryRun {
			return nil
		}
		report["entries"] = entries
		report["statistics"] = stats
		report["generated_at"] = time.Now().UTC().Format(time.RFC3339)
		if err := archive.SaveJSON(o.Output, report); err != nil {
			return err
		}
		return writeCSV(o.CSV, entries)
	}
	for _, post := range input.Posts {
		if !post.Accessible {
			continue
		}
		for _, link := range post.Links {
			if err := ctx.Err(); err != nil {
				return stats, err
			}
			stats.Checked++
			if !link.Broken {
				stats.Filtered++
				continue
			}
			e := archive.Entry{"post_url": post.URL, "original_href": link.OriginalHref, "link_text": link.Text, "archive_url": nil, "status": "failed", "reason": "", "attempted_snapshots": []any{}}
			if link.OriginalHref == "" {
				e["original_href"] = link.URL
			}
			if _, done := indices[key(e)]; done && o.Resume {
				stats.Skipped++
				continue
			}
			stats.Candidates++
			context, err := postContext(o.Archive.ContentDir, post.URL)
			if err != nil {
				e["reason"] = err.Error()
			} else {
				for k, value := range context {
					e[k] = value
				}
				lookup := link.URL
				if !webURL(lookup) {
					lookup = stringField(e, "original_href")
				}
				snapshots, err := v.FindSnapshots(ctx, lookup, "")
				if err != nil {
					e["reason"] = "Wayback lookup: " + err.Error()
				} else if len(snapshots) == 0 {
					e["reason"] = "No snapshots found in Wayback Machine"
				} else {
					attempts := []any{}
					for _, candidate := range snapshots {
						if err := ctx.Err(); err != nil {
							return stats, err
						}
						u, _ := url.Parse(candidate)
						parts := strings.Split(strings.TrimPrefix(u.Path, "/web/"), "/")
						timestamp := parts[0]
						attempt := map[string]any{"timestamp": timestamp, "url": candidate}
						if date, err := time.Parse("20060102150405", timestamp); err == nil {
							attempt["datetime"] = date.Format(time.RFC3339)
						}
						attempts = append(attempts, attempt)
						verification := v.VerifyURL(ctx, candidate, e)
						e["verification"] = verification
						if !verification.Valid {
							continue
						}
						if err := ctx.Err(); err != nil {
							return stats, err
						}
						e["archive_url"] = candidate
						backup, err := archive.ReplaceLink(o.Archive, e, candidate, backups)
						if err != nil {
							e["reason"] = err.Error()
						} else {
							e["status"] = "fixed"
							if o.Archive.DryRun {
								e["status"] = "would_fix"
							}
							e["reason"] = "Verified archive replacement"
							e["backup_path"] = backup
							stats.Updated++
							modified[post.URL] = true
						}
						break
					}
					e["attempted_snapshots"] = attempts
					if stringField(e, "reason") == "" {
						e["reason"] = "All snapshots invalid or blank"
					}
				}
			}
			if err := ctx.Err(); err != nil {
				return stats, err
			}
			if stringField(e, "status") == "failed" {
				stats.Failed++
			}
			stats.PostsModified = len(modified)
			if i, ok := indices[key(e)]; ok {
				entries[i] = e
			} else {
				indices[key(e)] = len(entries)
				entries = append(entries, e)
			}
			fmt.Fprintf(log, "%s: %s (%s)\n", stringField(e, "status"), stringField(e, "original_href"), stringField(e, "reason"))
			if err := save(); err != nil {
				return stats, err
			}
		}
	}
	return stats, save()
}

func redirectEligible(link Link, o RepairOptions) (bool, string) {
	if link.FinalURL == nil || !webURL(*link.FinalURL) {
		return false, "No final URL"
	}
	if link.Redirects < o.MinRedirects {
		return false, "Too few redirects"
	}
	if strings.TrimRight(link.URL, "/") == strings.TrimRight(*link.FinalURL, "/") {
		return false, "URLs are identical"
	}
	if o.Domain != "" {
		u, err := url.Parse(link.URL)
		if err != nil || !strings.Contains(strings.ToLower(u.Hostname()), strings.ToLower(o.Domain)) {
			return false, "Domain filter"
		}
	}
	if link.Broken || link.StatusCode == nil || *link.StatusCode < 200 || *link.StatusCode >= 400 {
		return false, "Link is broken"
	}
	return true, ""
}

func replaceRedirect(o RepairOptions, e archive.Entry, target string, backups *archive.BackupManager) (string, error) {
	backup, err := archive.ReplaceLink(o.Archive, e, target, backups)
	if err == nil || !o.CaseInsensitive {
		return backup, err
	}
	file, pathErr := archive.PostPath(o.Archive.ContentDir, stringField(e, "post_url"))
	if pathErr != nil {
		return "", pathErr
	}
	data, readErr := os.ReadFile(file)
	if readErr != nil {
		return "", readErr
	}
	pattern := regexp.MustCompile("(?i)" + regexp.QuoteMeta(stringField(e, "original_href")))
	seen := map[string]bool{}
	replaced := false
	for _, match := range pattern.FindAllString(string(data), -1) {
		if seen[match] {
			continue
		}
		seen[match] = true
		if _, count := archive.ReplaceURLs(string(data), match, target); count == 0 {
			continue
		}
		variant := archive.Entry{"post_url": e["post_url"], "original_href": match}
		backup, err = archive.ReplaceLink(o.Archive, variant, target, backups)
		if err != nil {
			return backup, err
		}
		replaced = true
	}
	if !replaced {
		return "", err
	}
	return backup, nil
}

func UpdateRedirects(ctx context.Context, o RepairOptions, log io.Writer) (RepairSummary, error) {
	var stats RepairSummary
	if log == nil {
		log = io.Discard
	}
	if o.MinRedirects < 1 {
		return stats, fmt.Errorf("min-redirects must be positive")
	}
	if err := validateRepair(o); err != nil {
		return stats, err
	}
	unlock, err := archive.Lock(ctx, o.Output+".lock", o.Archive.LockTimeout)
	if err != nil {
		return stats, err
	}
	defer unlock()
	input, err := LoadReport(o.Input)
	if err != nil {
		return stats, err
	}
	updates := []archive.Entry{}
	modified := map[string]bool{}
	backups := &archive.BackupManager{Base: o.Archive.BackupDir}
	save := func() error {
		if o.Archive.DryRun {
			return nil
		}
		posts := []string{}
		for p := range modified {
			posts = append(posts, p)
		}
		sort.Strings(posts)
		return archive.SaveJSON(o.Output, map[string]any{"generated_at": time.Now().UTC().Format(time.RFC3339), "mode": "live_update", "filters": map[string]any{"min_redirects": o.MinRedirects, "domain_filter": o.Domain}, "statistics": stats, "posts_modified": posts, "updates": updates, "backup_path": backups.Dir})
	}
	for _, post := range input.Posts {
		for _, link := range post.Links {
			if err := ctx.Err(); err != nil {
				return stats, err
			}
			stats.Checked++
			if ok, reason := redirectEligible(link, o); !ok {
				if reason == "URLs are identical" {
					stats.SameURL++
				} else {
					stats.Filtered++
				}
				continue
			}
			stats.Candidates++
			old := link.OriginalHref
			if old == "" {
				old = link.URL
			}
			e := archive.Entry{"post_url": post.URL, "post": post.URL, "original_href": old, "old_url": link.URL, "new_url": *link.FinalURL, "redirects": link.Redirects, "link_text": link.Text, "timestamp": time.Now().UTC().Format(time.RFC3339)}
			backup, err := replaceRedirect(o, e, *link.FinalURL, backups)
			if err != nil {
				stats.Failed++
				e["status"] = "failed"
				e["error"] = err.Error()
			} else {
				stats.Updated++
				modified[post.URL] = true
				e["status"] = "updated"
				if o.Archive.DryRun {
					e["status"] = "would_update"
				}
				e["backup_path"] = backup
			}
			stats.PostsModified = len(modified)
			updates = append(updates, e)
			fmt.Fprintf(log, "%s: %s -> %s\n", e["status"], old, *link.FinalURL)
			if err := save(); err != nil {
				return stats, err
			}
		}
	}
	return stats, save()
}
