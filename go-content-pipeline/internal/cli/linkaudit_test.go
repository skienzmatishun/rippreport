package cli

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLinkAuditHelp(t *testing.T) {
	for _, command := range []string{"check-links", "link-checker", "analyze-broken-links", "archive-links", "update-redirected-links", "audit-content", "content-audit"} {
		t.Run(command, func(t *testing.T) {
			var out, stderr bytes.Buffer
			code := NewApp(&out, &stderr).Run([]string{"pipeline", command, "--help"})
			if code != 0 || !strings.Contains(out.String(), "report") {
				t.Fatalf("%d: %s %s", code, out.String(), stderr.String())
			}
		})
	}
}

func TestLinkAuditCommandsEndToEnd(t *testing.T) {
	dir := t.TempDir()
	content := filepath.Join(dir, "content", "p")
	bundle := filepath.Join(content, "sample")
	if err := os.MkdirAll(bundle, 0755); err != nil {
		t.Fatal(err)
	}
	post := filepath.Join(bundle, "index.md")
	if err := os.WriteFile(post, []byte("[gallery ids=1]\n[old](/redirect)\n"), 0644); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/p/sample/":
			fmt.Fprint(w, `<a href="/redirect">old</a><a href="/missing.pdf">PDF</a>`)
		case "/redirect":
			http.Redirect(w, r, "/ok", 301)
		case "/ok":
			fmt.Fprint(w, "ok")
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	report := filepath.Join(dir, "audit.json")
	var out, stderr bytes.Buffer
	app := NewApp(&out, &stderr)
	code := app.Run([]string{"pipeline", "check-links", "--content-dir", content, "--base-url", server.URL, "--report", report, "--skip-domains=", "--request-delay=0s"})
	if code != 0 {
		t.Fatalf("checker %d: %s", code, stderr.String())
	}
	analysis := filepath.Join(dir, "analysis.json")
	code = app.Run([]string{"pipeline", "analyze-broken-links", "--report", report, "--output", analysis})
	if code != 0 {
		t.Fatalf("analysis %d: %s", code, stderr.String())
	}
	redirects := filepath.Join(dir, "redirects.json")
	code = app.Run([]string{"pipeline", "update-redirected-links", "--report", report, "--content-dir", content, "--output", redirects, "--backup-dir", filepath.Join(dir, "backups")})
	if code != 0 || !strings.Contains(out.String(), "would_update") {
		t.Fatalf("preview %d: %s %s", code, out.String(), stderr.String())
	}
	if _, err := os.Stat(redirects); !os.IsNotExist(err) {
		t.Fatal("default preview wrote report")
	}
	code = app.Run([]string{"pipeline", "update-redirected-links", "--report", report, "--content-dir", content, "--output", redirects, "--backup-dir", filepath.Join(dir, "backups"), "--update"})
	if code != 0 {
		t.Fatalf("update %d: %s", code, stderr.String())
	}
	data, _ := os.ReadFile(post)
	if !strings.Contains(string(data), server.URL+"/ok") {
		t.Fatal(string(data))
	}
	code = app.Run([]string{"pipeline", "audit-content", "--content-dir", content, "--report", filepath.Join(dir, "content.json"), "--no-lmstudio", "--request-delay=0s"})
	if code != 0 || !strings.Contains(out.String(), "1/1 posts") {
		t.Fatalf("content %d: %s %s", code, out.String(), stderr.String())
	}
}

func TestLinkAuditInvalidOptions(t *testing.T) {
	for _, args := range [][]string{
		{"check-links", "--post-workers=0"},
		{"archive-links", "--report=/tmp/shared.json", "--output=/tmp/shared.json"},
		{"update-redirected-links", "--min-redirects=0"},
		{"audit-content", "--request-delay=-1s"},
		{"analyze-broken-links", "--report=/tmp/shared.json", "--output=/tmp/shared.json"},
	} {
		var out, stderr bytes.Buffer
		if code := NewApp(&out, &stderr).Run(append([]string{"pipeline"}, args...)); code == 0 {
			t.Fatalf("accepted %v", args)
		}
	}
	var out, stderr bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if code := NewApp(&out, &stderr).auditError(ctx, ctx.Err()); code != 130 {
		t.Fatal(code)
	}
}
