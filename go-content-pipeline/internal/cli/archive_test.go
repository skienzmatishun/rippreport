package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArchiveCommandHelp(t *testing.T) {
	for _, command := range []string{"verify-archives", "archive-link-verifier", "missing-pdfs"} {
		t.Run(command, func(t *testing.T) {
			var out, stderr bytes.Buffer
			if code := NewApp(&out, &stderr).Run([]string{"pipeline", command, "--help"}); code != 0 {
				t.Fatalf("code %d: %s", code, stderr.String())
			}
			if !strings.Contains(out.String(), "report") {
				t.Fatal(out.String())
			}
		})
	}
}

func TestMissingPDFsCommand(t *testing.T) {
	dir := t.TempDir()
	report := filepath.Join(dir, "audit.json")
	output := filepath.Join(dir, "missing.json")
	if err := os.WriteFile(report, []byte(`{"entries":[{"post_url":"/p/test/","original_href":"https://example.com/missing.pdf","status":"broken"}]}`), 0644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := NewApp(&stdout, &stderr).Run([]string{"pipeline", "missing-pdfs", "--report", report, "--output", output})
	if code != 0 {
		t.Fatalf("%d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "1 missing PDFs") {
		t.Fatal(stdout.String())
	}
	code = NewApp(&stdout, &stderr).Run([]string{"pipeline", "missing-pdfs", "--report", report, "--output", report})
	if code == 0 {
		t.Fatal("allowed overwriting report")
	}
}

func TestArchiveCommandDryRun(t *testing.T) {
	dir := t.TempDir()
	report := filepath.Join(dir, "audit.json")
	posts := filepath.Join(dir, "content", "p")
	bundle := filepath.Join(posts, "sample")
	if err := os.MkdirAll(bundle, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bundle, "index.md"), []byte("[tag](https://rippreport.com/tags/news/)"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(report, []byte(`{"entries":[{"post_url":"/p/sample/","original_href":"https://rippreport.com/tags/news/"}]}`), 0644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := NewApp(&stdout, &stderr).Run([]string{"pipeline", "verify-archives", "--report", report, "--content-dir", posts, "--progress", filepath.Join(dir, "progress.json"), "--dry-run", "--no-playwright", "--no-lmstudio", "--no-google"})
	if code != 0 {
		t.Fatalf("%d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "1 would replace") {
		t.Fatal(stdout.String())
	}
}
