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

func TestApp_GenerateFrontpageHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := NewApp(&stdout, &stderr)

	code := app.Run([]string{"pipeline", "generate-frontpage", "--help"})
	if code != 0 {
		t.Errorf("expected exit code 0 for help, got %d", code)
	}
	out := stdout.String()
	if !strings.Contains(out, "Usage of generate-frontpage") {
		t.Errorf("expected usage output, got: %s", out)
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

func TestBuildFrontMatter(t *testing.T) {
	title := "Test Podcast Episode"
	slug := "test-podcast-episode"
	desc := "Episode description"
	youtubeID := "L6dyOWiaei4"
	categories := []string{"dixie-podcast"}
	tags := []string{"test-tag"}
	body := "This is a summary.\n\n{{< youtubecomments >}}\n{{< aicomments >}}\n"

	result := buildFrontMatter(title, slug, desc, youtubeID, categories, tags, true, body)

	if !strings.Contains(result, "title: Test Podcast Episode") {
		t.Errorf("expected title in frontmatter, got: %s", result)
	}
	if !strings.Contains(result, "youtube: L6dyOWiaei4") {
		t.Errorf("expected youtube ID in frontmatter, got: %s", result)
	}
	if !strings.Contains(result, "thumbnail: img/test-podcast-episode.png") {
		t.Errorf("expected thumbnail in frontmatter, got: %s", result)
	}
	if !strings.Contains(result, "- dixie-podcast") {
		t.Errorf("expected category in frontmatter, got: %s", result)
	}
	if !strings.Contains(result, "{{< youtubecomments >}}") {
		t.Errorf("expected youtube comments shortcode in body, got: %s", result)
	}
	if !strings.Contains(result, "{{< aicomments >}}") {
		t.Errorf("expected aicomments shortcode in body, got: %s", result)
	}
}

func TestApp_NewPostRequiresTitle(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := NewApp(&stdout, &stderr)

	code := app.Run([]string{"pipeline", "new-post", "--email"})
	if code != 1 {
		t.Errorf("expected exit code 1 when --title is missing, got %d", code)
	}
	if !strings.Contains(stderr.String(), "--title flag is required") {
		t.Errorf("expected required title error message, got: %s", stderr.String())
	}
}

func TestApp_EmailPostRequiresTitle(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := NewApp(&stdout, &stderr)

	code := app.Run([]string{"pipeline", "email-post"})
	if code != 1 {
		t.Errorf("expected exit code 1 when --title is missing for email-post, got %d", code)
	}
	if !strings.Contains(stderr.String(), "--title flag is required") {
		t.Errorf("expected required title error message, got: %s", stderr.String())
	}
}

func TestGeminiHookFlags(t *testing.T) {
	for _, command := range []string{"generate-hooks", "pipeline"} {
		t.Run(command, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			app := NewApp(&stdout, &stderr)
			if code := app.Run([]string{"pipeline", command, "--help"}); code != 0 {
				t.Fatalf("help failed: %s", stderr.String())
			}
			for _, flag := range []string{"-plugsky", "-plugsky-model", "-plugsky-max-tokens", "plugsky-lite", "-gemini", "-gemini-model", "-gemini-max-tokens", "gemini-3.5-flash-lite", "-openrouter", "-cloudflare", "-cloudflare-model", "-cloudflare-max-tokens", "-cloudflare-account-id", "@cf/google/gemma-4-26b-a4b-it"} {
				if !strings.Contains(stdout.String(), flag) {
					t.Errorf("missing %s", flag)
				}
			}
			if strings.Contains(stdout.String(), "huggingface") {
				t.Fatal("removed provider still advertised")
			}
		})
	}
}
