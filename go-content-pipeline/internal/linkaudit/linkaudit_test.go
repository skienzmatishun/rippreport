package linkaudit

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/archive"
)

func put(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
func save(t *testing.T, path string, data any) {
	t.Helper()
	if err := archive.SaveJSON(path, data); err != nil {
		t.Fatal(err)
	}
}
func sample(t *testing.T) (string, RepairOptions) {
	t.Helper()
	dir := t.TempDir()
	a := archive.DefaultOptions()
	a.ContentDir = filepath.Join(dir, "content", "p")
	a.BackupDir = filepath.Join(dir, "backups")
	a.UsePlaywright = false
	a.UseLMStudio = false
	put(t, filepath.Join(a.ContentDir, "sample", "index.md"), "---\ntitle: Sample\ndate: 2020-01-02\n---\n[old](https://old.test/a \"title\")\n")
	return dir, RepairOptions{Input: filepath.Join(dir, "input.json"), Output: filepath.Join(dir, "output.json"), CSV: filepath.Join(dir, "output.csv"), Archive: a, MinRedirects: 1, Resume: true, CaseInsensitive: true}
}

func TestExtractLinksAndExclusions(t *testing.T) {
	source := `<a href="../relative?a=1&amp;b=2">Some <b>text</b></a><a href="/root">root</a><a href="//other.test/file">other</a><a href="#anchor">skip</a><a href="mailto:a@b.test">mail</a><a href="tel:123">tel</a><a href="javascript:alert(1)">js</a><a href="https://m.facebook.com/p">fb</a><a href="https://facebook.com.evil.test/">keep</a><a href="https://x.test/p/privacy-policy/">privacy</a>`
	o := DefaultCheckOptions()
	links, err := ExtractLinks(source, "https://site.test/p/post/", o.SkipDomains, o.SkipPaths)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 4 || links[0].URL != "https://site.test/p/relative?a=1&b=2" || links[0].OriginalHref != "../relative?a=1&b=2" || links[0].Text != "Some text" {
		t.Fatalf("%+v", links)
	}
	if !ShouldSkip("http://localhost:1313/p/other/", o.SkipDomains, o.SkipPaths) {
		t.Fatal("localhost not excluded")
	}
}

func TestCheckerRedirectsTimeoutResumeAndRecheck(t *testing.T) {
	dir, _ := sample(t)
	var hits atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/p/sample/":
			hits.Add(1)
			fmt.Fprint(w, `<a href="/ok">ok</a><a href="/missing.pdf">missing</a><a href="/redirect">redirect</a><a href="/loop">loop</a><a href="/slow">slow</a>`)
		case "/ok":
			fmt.Fprint(w, "ok")
		case "/missing.pdf":
			http.NotFound(w, r)
		case "/redirect":
			http.Redirect(w, r, "/ok", http.StatusMovedPermanently)
		case "/loop":
			http.Redirect(w, r, "/loop", http.StatusFound)
		case "/slow":
			select {
			case <-r.Context().Done():
			case <-time.After(time.Second):
			}
		}
	}))
	defer server.Close()
	o := DefaultCheckOptions()
	o.ContentDir = filepath.Join(dir, "content", "p")
	o.Report = filepath.Join(dir, "audit.json")
	o.BaseURL = server.URL
	o.SkipDomains = nil
	o.Delay = 0
	o.Timeout = 150 * time.Millisecond
	o.MaxRedirects = 2
	c := NewChecker(o, nil)
	report, err := c.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if report.Summary.TotalLinks != 5 || report.Summary.BrokenLinks != 3 {
		t.Fatalf("%+v", report.Summary)
	}
	links := report.Posts[0].Links
	if links[2].Redirects != 1 || links[2].FinalURL == nil || !strings.HasSuffix(*links[2].FinalURL, "/ok") {
		t.Fatalf("%+v", links[2])
	}
	if links[3].Status != "too_many_redirects" || links[4].Status != "timeout" {
		t.Fatalf("%+v", links)
	}
	if _, err := c.Run(context.Background()); err != nil || hits.Load() != 1 {
		t.Fatalf("resume hits %d, err %v", hits.Load(), err)
	}
	c.Options.Resume = false
	if _, err := c.Run(context.Background()); err != nil || hits.Load() != 2 {
		t.Fatalf("recheck hits %d, err %v", hits.Load(), err)
	}
	_, _, err = archive.LoadReport(o.Report)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckerConcurrencyAndCancellation(t *testing.T) {
	dir, _ := sample(t)
	root := filepath.Join(dir, "content", "p")
	put(t, filepath.Join(root, "second", "index.md"), "post")
	var active, maximum atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/p/") {
			fmt.Fprint(w, strings.Repeat(`<a href="/ok">ok</a>`, 8))
			return
		}
		n := active.Add(1)
		defer active.Add(-1)
		for {
			previous := maximum.Load()
			if n <= previous || maximum.CompareAndSwap(previous, n) {
				break
			}
		}
		time.Sleep(5 * time.Millisecond)
		fmt.Fprint(w, "ok")
	}))
	defer server.Close()
	o := DefaultCheckOptions()
	o.ContentDir = root
	o.Report = filepath.Join(dir, "audit.json")
	o.BaseURL = server.URL
	o.SkipDomains = nil
	o.Delay = 0
	o.LinkWorkers = 2
	report, err := NewChecker(o, nil).Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if report.Summary.TotalLinks != 16 || maximum.Load() > 2 {
		t.Fatalf("%+v max=%d", report.Summary, maximum.Load())
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := NewChecker(o, nil).Run(ctx); err == nil {
		t.Fatal("ignored cancellation")
	}
}

func TestBundleSelectionAndCorruptReports(t *testing.T) {
	dir, o := sample(t)
	put(t, filepath.Join(o.Archive.ContentDir, "z", "index.md"), "post")
	names, err := bundles(o.Archive.ContentDir, []string{"sample"})
	if err != nil || len(names) != 1 || names[0] != "sample" {
		t.Fatalf("%v %v", names, err)
	}
	if _, err := bundles(o.Archive.ContentDir, []string{"missing"}); err == nil {
		t.Fatal("ignored missing selected bundle")
	}
	path := filepath.Join(dir, "bad.json")
	for _, bad := range []string{`{`, `{"entries":[]}`, `{"posts":null}`, `{"posts":[{}]}`} {
		put(t, path, bad)
		if _, err := LoadReport(path); err == nil {
			t.Fatalf("accepted %s", bad)
		}
	}
}

func TestAnalysisRecountsAndGroups(t *testing.T) {
	code := 404
	report := &Report{Posts: []Post{{URL: "/p/sample/", Accessible: true, Links: []Link{{URL: "https://x.test/missing.pdf", Status: "ok", StatusCode: &code, Broken: true}, {URL: "https://x.test/timeout", Status: "timeout", Broken: true}}}, {URL: "/p/offline/", Links: []Link{}}}}
	a := Analyze(report)
	if a.Summary.BrokenLinks != 2 || a.Summary.InaccessiblePosts != 1 || len(a.ByDomain["x.test"]) != 2 || len(a.ByStatus["unknown"]) != 1 || len(a.ByError["timeout"]) != 1 {
		t.Fatalf("%+v", a)
	}
	var out bytes.Buffer
	a.Print(&out)
	if !strings.Contains(out.String(), "Posts with most broken links") {
		t.Fatal(out.String())
	}
}

func TestRedirectUpdatePreviewBackupAndFilters(t *testing.T) {
	_, o := sample(t)
	code := 200
	final := "https://new.test/a"
	o.CSV = ""
	input := Report{Posts: []Post{{URL: "/p/sample/", Accessible: true, Links: []Link{{URL: "https://old.test/a", OriginalHref: "https://old.test/a", FinalURL: &final, StatusCode: &code, Redirects: 2}, {URL: "https://bad.test/a", FinalURL: &final, StatusCode: &code, Redirects: 2, Broken: true}}}}}
	save(t, o.Input, input)
	o.Archive.DryRun = true
	stats, err := UpdateRedirects(context.Background(), o, nil)
	if err != nil || stats.Updated != 1 || stats.Filtered != 1 {
		t.Fatalf("%+v %v", stats, err)
	}
	if _, err := os.Stat(o.Output); !os.IsNotExist(err) {
		t.Fatal("dry run created report")
	}
	if _, err := os.Stat(o.Archive.BackupDir); !os.IsNotExist(err) {
		t.Fatal("dry run created backups")
	}
	o.Archive.DryRun = false
	stats, err = UpdateRedirects(context.Background(), o, nil)
	if err != nil || stats.Updated != 1 {
		t.Fatalf("%+v %v", stats, err)
	}
	post, _ := os.ReadFile(filepath.Join(o.Archive.ContentDir, "sample", "index.md"))
	if !strings.Contains(string(post), `[old](https://new.test/a "title")`) {
		t.Fatal(string(post))
	}
	manifests, _ := filepath.Glob(filepath.Join(o.Archive.BackupDir, "*", "manifest.json"))
	if len(manifests) != 1 {
		t.Fatal(manifests)
	}
	if err := archive.RestoreBackup(manifests[0]); err != nil {
		t.Fatal(err)
	}
	o.Domain = "elsewhere"
	stats, err = UpdateRedirects(context.Background(), o, nil)
	if err != nil || stats.Updated != 0 {
		t.Fatalf("%+v %v", stats, err)
	}
}

func TestRedirectRelativeAndCaseInsensitive(t *testing.T) {
	_, o := sample(t)
	o.CSV = ""
	post := filepath.Join(o.Archive.ContentDir, "sample", "index.md")
	put(t, post, `[relative](/Old/Path "title") https://old.test/ab`)
	code := 200
	target := "https://new.test/destination"
	save(t, o.Input, Report{Posts: []Post{{URL: "/p/sample/", Links: []Link{{URL: "https://site.test/old/path", OriginalHref: "/old/path", FinalURL: &target, StatusCode: &code, Redirects: 1}}}}})
	stats, err := UpdateRedirects(context.Background(), o, nil)
	if err != nil || stats.Updated != 1 {
		t.Fatalf("%+v %v", stats, err)
	}
	data, _ := os.ReadFile(post)
	if string(data) != `[relative](https://new.test/destination "title") https://old.test/ab` {
		t.Fatal(string(data))
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestArchiveDiscoveryCSVResumeAndVerifierCompatibility(t *testing.T) {
	_, o := sample(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cdx" {
			if r.URL.Query().Get("url") != "https://old.test/a" {
				t.Error(r.URL.RawQuery)
			}
			fmt.Fprint(w, `[["timestamp","original","statuscode"],["20200102000000","https://old.test/a","200"],["20190102000000","https://old.test/a","200"]]`)
			return
		}
		if strings.Contains(r.URL.Path, "20200102000000") {
			fmt.Fprint(w, "<p>blank</p>")
			return
		}
		fmt.Fprint(w, "<h1>Archived article</h1><p>"+strings.Repeat("Meaningful archived article content. ", 20)+"</p>")
	}))
	defer server.Close()
	o.Archive.CDXURL = server.URL + "/cdx"
	// Redirect Wayback requests into a local fixture while preserving real
	// archive URLs in the persisted report and Markdown replacement.
	previous := http.DefaultTransport
	defer func() { http.DefaultTransport = previous }()
	http.DefaultTransport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Hostname() == "web.archive.org" {
			clone := r.Clone(r.Context())
			u := *r.URL
			base, _ := http.NewRequest("GET", server.URL, nil)
			u.Scheme, u.Host = base.URL.Scheme, base.URL.Host
			clone.URL = &u
			return previous.RoundTrip(clone)
		}
		return previous.RoundTrip(r)
	})
	save(t, o.Input, Report{Posts: []Post{{URL: "/p/sample/", Accessible: true, Links: []Link{{URL: "https://old.test/a", OriginalHref: "https://old.test/a", Text: "article", Broken: true}}}}})
	stats, err := ArchiveBroken(context.Background(), o, nil)
	if err != nil || stats.Updated != 1 {
		t.Fatalf("%+v %v", stats, err)
	}
	_, entries, err := archive.LoadReport(o.Output)
	if err != nil || len(entries) != 1 || entries[0]["status"] != "fixed" || !strings.Contains(entries[0]["archive_url"].(string), "20190102000000") || entries[0]["post_date"] != "2020-01-02T00:00:00Z" {
		t.Fatalf("%+v %v", entries, err)
	}
	f, err := os.Open(o.CSV)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := csv.NewReader(f).ReadAll()
	f.Close()
	if err != nil || len(rows) != 2 {
		t.Fatalf("%v %v", rows, err)
	}
	stats, err = ArchiveBroken(context.Background(), o, nil)
	if err != nil || stats.Skipped != 1 {
		t.Fatalf("%+v %v", stats, err)
	}
	vOptions := o.Archive
	vOptions.Report = o.Output
	vOptions.Progress = filepath.Join(filepath.Dir(o.Output), "progress.json")
	verified, err := archive.New(vOptions, nil).Run(context.Background())
	if err != nil || verified.Replaced != 1 {
		t.Fatalf("%+v %v", verified, err)
	}
}

func TestArchiveDiscoveryFailureAndDryRun(t *testing.T) {
	_, o := sample(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, `[["timestamp","original","statuscode"]]`) }))
	defer server.Close()
	o.Archive.CDXURL = server.URL
	save(t, o.Input, Report{Posts: []Post{{URL: "/p/sample/", Accessible: true, Links: []Link{{URL: "https://old.test/a", OriginalHref: "https://old.test/a", Broken: true}}}}})
	o.Archive.DryRun = true
	stats, err := ArchiveBroken(context.Background(), o, nil)
	if err != nil || stats.Failed != 1 {
		t.Fatalf("%+v %v", stats, err)
	}
	for _, p := range []string{o.Output, o.CSV, o.Archive.BackupDir} {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Errorf("dry run created %s", p)
		}
	}
	o.Archive.DryRun = false
	stats, err = ArchiveBroken(context.Background(), o, nil)
	if err != nil || stats.Failed != 1 {
		t.Fatalf("%+v %v", stats, err)
	}
}

func TestWordPressIssues(t *testing.T) {
	content := "---\ntitle: '[gallery]'\n---\n[gallery ids=\"1,2\"]\n[caption]text[/caption]\n[gallery](https://x.test)\n[video][ref]\n[audio]: https://x.test\n`[embed]`\n```html\n[playlist]\n```\n~~~\n[gallery]\n~~~\n    [audio]\n\\[gallery]\n{{< example text=\"[audio]\" >}}\n<pre>[video]</pre>\n<iframe src='https://x.test'></iframe>\n"
	issues := WordPressIssues(content)
	if len(issues) != 3 {
		t.Fatalf("%v", issues)
	}
}

func TestContentAuditOfflineModelThinkingAndResume(t *testing.T) {
	dir, repair := sample(t)
	put(t, filepath.Join(repair.Archive.ContentDir, "sample", "index.md"), "[gallery ids=1]")
	o := ContentOptions{ContentDir: repair.Archive.ContentDir, Report: filepath.Join(dir, "content.json"), Resume: true, Model: repair.Archive}
	results, err := AuditContent(context.Background(), o, nil)
	if err != nil || len(results) != 1 || !results[0].HasIssue {
		t.Fatalf("%+v %v", results, err)
	}
	put(t, filepath.Join(repair.Archive.ContentDir, "sample", "index.md"), "clean")
	results, err = AuditContent(context.Background(), o, nil)
	if err != nil || !results[0].HasIssue {
		t.Fatalf("resume: %+v %v", results, err)
	}
	o.Resume = false
	o.UseModel = true
	o.Model.Retries = 1
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if !bytes.Contains(body, []byte("has_issue")) {
			t.Error(string(body))
		}
		response := map[string]any{"choices": []any{map[string]any{"message": map[string]string{"content": "<think>{ignored reasoning}</think>\n```json\n{\"has_issue\":false,\"issues\":[]}\n```"}}}}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	o.Model.ModelURL = server.URL
	results, err = AuditContent(context.Background(), o, nil)
	if err != nil || results[0].HasIssue || results[0].Error != "" {
		t.Fatalf("%+v %v", results, err)
	}
	o.Model.ModelURL = server.URL + "/invalid"
	o.Model.ModelTimeout = time.Nanosecond
	results, err = AuditContent(context.Background(), o, nil)
	if err != nil || results[0].Error == "" {
		t.Fatalf("%+v %v", results, err)
	}
}
