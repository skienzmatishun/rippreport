package archive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func write(t *testing.T, file, text string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(text), 0644); err != nil {
		t.Fatal(err)
	}
}

func fixture(t *testing.T) Options {
	t.Helper()
	root := t.TempDir()
	o := DefaultOptions()
	o.Report = filepath.Join(root, "report.json")
	o.Progress = filepath.Join(root, "progress.json")
	o.ContentDir = filepath.Join(root, "content", "p")
	o.BackupDir = filepath.Join(root, "backups")
	o.ScreenshotDir = filepath.Join(root, "screenshots")
	o.UsePlaywright = false
	o.UseLMStudio = false
	o.Google = false
	write(t, filepath.Join(o.ContentDir, "sample", "index.md"), "---\ntitle: Sample\n---\n[text](https://old.example/doc.pdf \"title\")\n")
	return o
}

func TestStrategies(t *testing.T) {
	o := DefaultOptions()
	for _, tc := range []struct{ url, status, archive, want string }{
		{"https://old.example/a", "fixed", "https://web.archive.org/web/20200101000000/https://old.example/a", "verify_fixed"},
		{"https://old.example/A.PDF?download=1", "broken", "", "pdf_fallback"},
		{"https://cdn.rippreport.com/no-extension", "broken", "", "malformed_cdn"},
		{"https://rippreport.com/tags/corruption/", "broken", "", "tag_conversion"},
		{"https://rippreport.com/sample/", "broken", "", "missing_p_path"},
		{"https://rippreport.com.evil.test/tags/corruption/", "broken", "", "paywall_check"},
		{"https://other.example/article", "broken", "", "paywall_check"},
	} {
		t.Run(tc.want+tc.url, func(t *testing.T) {
			got := Strategy(Entry{"original_href": tc.url, "status": tc.status, "archive_url": tc.archive}, o)
			if got != tc.want {
				t.Fatalf("got %s, want %s", got, tc.want)
			}
		})
	}
	o.Google = false
	if got := Strategy(Entry{"original_href": "https://old.example/a.pdf"}, o); got != "skip" {
		t.Fatal(got)
	}
}

func TestReplaceURLs(t *testing.T) {
	old := "https://example.com/doc.pdf"
	target := "https://web.archive.org/web/20200101000000/" + old
	input := "[text](" + old + " \"title\") <" + old + "> " + old + ".\n<a href=\"" + old + "\">PDF</a> [ref]: " + old + "\n" + old + "?other=1 " + old + "/child https://mirror.example/" + old
	got, n := ReplaceURLs(input, old, target)
	if n != 5 {
		t.Fatalf("replaced %d: %s", n, got)
	}
	if !strings.Contains(got, "[text]("+target+" \"title\")") || !strings.Contains(got, old+"?other=1") || strings.Contains(got, "/"+target) {
		t.Fatal(got)
	}
	if _, n := ReplaceURLs(got, old, target); n != 0 {
		t.Fatal("reprocessed embedded URL")
	}
}

func TestHTMLVerification(t *testing.T) {
	text := strings.Repeat("Meaningful news article text with details. ", 30)
	for _, tc := range []struct {
		html  string
		valid bool
	}{
		{"<article><h1>Story</h1><p>" + text + "</p></article>", true},
		{"<div id='wm-ipp-base'><p>" + text + "</p></div><p>blank</p>", false},
		{"<div class='wayback-banner'><div><p>" + text + "</p></div></div>", false},
		{"<script>" + text + "</script><p>blank</p>", false},
		{"<p>Access denied " + text + "</p>", false},
	} {
		if got := inspectHTML(tc.html).substantial(false); got != tc.valid {
			t.Errorf("valid = %v for %s", got, tc.html[:30])
		}
	}
	if inspectHTML("<article><p>Subscribe subscription sign in " + text + "</p></article>").substantial(true) {
		t.Fatal("accepted subscription gate")
	}
}

func TestRunRepairResumeBackupRestore(t *testing.T) {
	o := fixture(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "<h1>Story</h1><p>"+strings.Repeat("Article content. ", 40)+"</p>")
	}))
	defer server.Close()
	entry := Entry{"post_url": "https://rippreport.com/p/sample/", "original_href": "https://old.example/doc.pdf", "status": "fixed", "archive_url": server.URL, "custom": "preserve"}
	if err := saveJSON(o.Report, map[string]any{"metadata": map[string]any{"test": "keep"}, "entries": []Entry{entry}}); err != nil {
		t.Fatal(err)
	}
	post := filepath.Join(o.ContentDir, "sample", "index.md")
	original, _ := os.ReadFile(post)
	v := New(o, nil)
	result, err := v.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result.Replaced != 1 {
		t.Fatalf("%+v", result)
	}
	data, _ := os.ReadFile(post)
	if !bytes.Contains(data, []byte(server.URL)) {
		t.Fatal(string(data))
	}
	_, entries, err := LoadReport(o.Report)
	if err != nil {
		t.Fatal(err)
	}
	if entries[0].str("custom") != "preserve" || entries[0].str("status") != "replaced" {
		t.Fatal(entries)
	}
	result, err = v.Run(context.Background())
	if err != nil || result.Skipped != 1 {
		t.Fatalf("%+v %v", result, err)
	}
	manifests, _ := filepath.Glob(filepath.Join(o.BackupDir, "*", "manifest.json"))
	if len(manifests) != 1 {
		t.Fatal(manifests)
	}
	if err := RestoreBackup(manifests[0]); err != nil {
		t.Fatal(err)
	}
	restored, _ := os.ReadFile(post)
	if !bytes.Equal(restored, original) {
		t.Fatal("restore did not recover original")
	}
}

func TestDryRunLeavesStateUntouched(t *testing.T) {
	o := fixture(t)
	o.DryRun = true
	write(t, filepath.Join(o.ContentDir, "sample", "index.md"), "[tag](https://rippreport.com/tags/corruption/)")
	write(t, o.Report, `{"entries":[{"post_url":"/p/sample/","original_href":"https://rippreport.com/tags/corruption/"}]}`)
	before, _ := os.ReadFile(o.Report)
	result, err := New(o, nil).Run(context.Background())
	if err != nil || result.WouldReplace != 1 {
		t.Fatalf("%+v %v", result, err)
	}
	after, _ := os.ReadFile(o.Report)
	if !bytes.Equal(before, after) {
		t.Fatal("dry run rewrote report")
	}
	for _, p := range []string{o.Progress, o.BackupDir, o.Report + ".lock"} {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Errorf("dry run created %s", p)
		}
	}
}

func TestExistingArchiveCanBeReplacedAndSameURLVerified(t *testing.T) {
	o := fixture(t)
	post := filepath.Join(o.ContentDir, "sample", "index.md")
	write(t, post, "[link](https://web.archive.org/web/20200101000000/https://old.example/doc.pdf)")
	e := Entry{"post_url": "/p/sample/", "original_href": "https://old.example/doc.pdf", "archive_url": "https://web.archive.org/web/20200101000000/https://old.example/doc.pdf"}
	manager := &BackupManager{Base: o.BackupDir}
	if _, err := replace(o, e, "https://new.example/doc.pdf", manager); err != nil {
		t.Fatal(err)
	}
	// Idempotence after a report/progress save interruption.
	if _, err := replace(o, e, "https://new.example/doc.pdf", manager); err != nil {
		t.Fatal(err)
	}
	write(t, post, "[link](https://old.example/doc.pdf)")
	if _, err := replace(o, e, e.str("original_href"), manager); err != nil {
		t.Fatal(err)
	}
}

func TestPostPathRejectsTraversalAndEscapingSymlink(t *testing.T) {
	o := fixture(t)
	outside := t.TempDir()
	write(t, filepath.Join(outside, "index.md"), "outside")
	if err := os.Symlink(outside, filepath.Join(o.ContentDir, "escape")); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"/p/../outside/", "/p/%2e%2e/outside/", "/p/escape/", "/p/", "/about/"} {
		if got, err := PostPath(o.ContentDir, p); err == nil {
			t.Errorf("accepted %s as %s", p, got)
		}
	}
}

func TestReportsAndMissingPDFExport(t *testing.T) {
	dir := t.TempDir()
	report := filepath.Join(dir, "audit.json")
	output := filepath.Join(dir, "pdfs.json")
	write(t, report, `{"posts":[{"post_url":"/p/a/","links":[{"original_href":"https://x.test/MISSING.PDF?download=1","is_broken":true,"status_code":404},{"original_href":"https://x.test/ok.pdf","is_broken":false}]},{"post_url":"/p/b/","links":[{"original_href":"https://x.test/MISSING.PDF?download=1","is_broken":true}]}]}`)
	count, err := ExportMissingPDFs(report, output)
	if err != nil || count != 1 {
		t.Fatalf("%d %v", count, err)
	}
	d, _ := readObject(output)
	if d["total_references"] != float64(2) {
		t.Fatal(d)
	}
	write(t, report, `{"processed":[{"original_href":"https://x.test/a.pdf","status":"fixed","archive_url":"https://archive.test/a"},{"original_href":"https://x.test/b.pdf","status":"verification_failed"}],"extra":42}`)
	d, entries, err := LoadReport(report)
	if err != nil || len(entries) != 2 || d["extra"] != float64(42) {
		t.Fatalf("%v %v", d, err)
	}
	count, err = ExportMissingPDFs(report, output)
	if err != nil || count != 1 {
		t.Fatalf("%d %v", count, err)
	}
	write(t, report, `{"entries":42}`)
	if _, _, err := LoadReport(report); err == nil {
		t.Fatal("accepted invalid entries")
	}
}

func TestModelAndHTTPFallback(t *testing.T) {
	for _, vision := range []bool{false, true} {
		t.Run(fmt.Sprint(vision), func(t *testing.T) {
			o := fixture(t)
			o.UseLMStudio = true
			o.Retries = 1
			o.Vision = vision
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/v1/chat/completions" {
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Error(err)
					}
					fmt.Fprint(w, `{"choices":[{"message":{"content":"`+"```json\\n"+`{\"is_valid\":false,\"reason\":\"Wrong article\",\"confidence\":0.9}`+"\\n```"+`"}}]}`)
					return
				}
				fmt.Fprint(w, "<article><p>"+strings.Repeat("Substantial page content. ", 50)+"</p></article>")
			}))
			defer server.Close()
			o.ModelURL = server.URL
			r := New(o, nil).verify(context.Background(), server.URL, Entry{"original_href": server.URL}, false)
			if r.Valid || !strings.Contains(r.Method, "lmstudio") || r.Reason != "Wrong article" {
				t.Fatalf("%+v", r)
			}
			o.ModelURL = server.URL + "/unavailable"
			r = New(o, nil).verify(context.Background(), server.URL, Entry{}, false)
			if !r.Valid || r.Method != "http+html" {
				t.Fatalf("%+v", r)
			}
		})
	}
}

func TestPDFSearchAndSnapshotOrdering(t *testing.T) {
	o := fixture(t)
	o.Google = true
	o.SearchDelay = 0
	o.SearchJitter = 0
	var base string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/search":
			fmt.Fprintf(w, `<a href="/url?q=%s/file.pdf">PDF</a><a href="%s/file.pdf">duplicate</a>`, base, base)
		case "/file.pdf":
			w.Header().Set("Content-Type", "application/pdf")
		case "/cdx":
			if r.URL.Query().Get("limit") != "-15" {
				t.Error(r.URL.RawQuery)
			}
			fmt.Fprint(w, `[["timestamp","original","statuscode"],["20190101000000","https://x.test/a","200"],["20210101000000","https://x.test/a","200"],["20200101000000","https://x.test/a","200"]]`)
		}
	}))
	defer server.Close()
	base = server.URL
	o.SearchURL = base + "/search"
	o.CDXURL = base + "/cdx"
	v := New(o, nil)
	target, r := v.Process(context.Background(), Entry{"original_href": "https://x.test/old.pdf"})
	if !r.Valid || target != base+"/file.pdf" {
		t.Fatalf("%s %+v", target, r)
	}
	snapshots, err := v.snapshots(context.Background(), "https://x.test/a", "20210101000000")
	if err != nil || len(snapshots) != 2 || !strings.Contains(snapshots[0], "20200101000000") {
		t.Fatalf("%v %v", snapshots, err)
	}
}

func TestLockCancellationAndCorruptProgress(t *testing.T) {
	o := fixture(t)
	unlock, err := acquire(context.Background(), o.Report+".lock", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := acquire(ctx, o.Report+".lock", time.Second); err == nil {
		t.Fatal("ignored cancellation")
	}
	write(t, o.Report, `{"entries":[]}`)
	write(t, o.Progress, "corrupted")
	o.DryRun = true
	if _, err := New(o, nil).Run(context.Background()); err == nil {
		t.Fatal("ignored corrupt progress")
	}
	data, _ := os.ReadFile(o.Progress)
	if string(data) != "corrupted" {
		t.Fatal("destroyed corrupt state")
	}
}
