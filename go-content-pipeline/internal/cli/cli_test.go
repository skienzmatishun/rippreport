package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestApp_Help(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := NewApp(&stdout, &stderr)

	code := app.Run([]string{"pipeline", "--help"})
	if code != 0 {
		t.Errorf("expected exit code 0 for help, got %d", code)
	}

	out := stdout.String()
	if !strings.Contains(out, "generate-related") {
		t.Errorf("help output missing generate-related command")
	}
	if !strings.Contains(out, "generate-hooks") {
		t.Errorf("help output missing generate-hooks command")
	}
}

func TestApp_UnknownCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := NewApp(&stdout, &stderr)

	code := app.Run([]string{"pipeline", "unknown-cmd"})
	if code != 1 {
		t.Errorf("expected exit code 1 for unknown command, got %d", code)
	}

	errOut := stderr.String()
	if !strings.Contains(errOut, "Unknown command") {
		t.Errorf("stderr missing unknown command error message: %s", errOut)
	}
}

func TestApp_GenerateRelatedDryRun(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := NewApp(&stdout, &stderr)

	code := app.Run([]string{"pipeline", "generate-related", "--dry-run", "--vector-only", "--posts", "some-slug"})
	// Should complete with exit code 0 (even if post doesn't exist, discovery returns 0 posts)
	if code != 0 {
		t.Errorf("expected exit code 0 for dry-run, got %d. stderr: %s", code, stderr.String())
	}
}

func TestApp_ParseSlugs(t *testing.T) {
	slugs := parseSlugs("post-1, post-2,post-3  ")
	if len(slugs) != 3 {
		t.Fatalf("expected 3 slugs, got %d", len(slugs))
	}
	if slugs[0] != "post-1" || slugs[1] != "post-2" || slugs[2] != "post-3" {
		t.Errorf("unexpected parsed slugs: %v", slugs)
	}

	empty := parseSlugs("")
	if empty != nil {
		t.Errorf("expected nil for empty string, got %v", empty)
	}
}
