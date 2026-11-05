package processor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/config"
	"github.com/rippreport/go-content-pipeline/internal/models"
)

func TestGenerateHook_Basic(t *testing.T) {
	main := &models.Post{
		Path: "/p/main/index.md",
		Slug: "main",
		FrontMatter: &models.FrontMatter{
			Title: "Grand Hotel Renovation",
			Date:  time.Now(),
		},
		Body: "Long article describing the history and preservation of the historic Grand Hotel.",
	}
	widget := &models.Post{
		Path: "/p/widget/index.md",
		Slug: "widget",
		FrontMatter: &models.FrontMatter{
			Title: "Point Clear Coastal Heritage",
			Date:  time.Now(),
		},
		Body: "Historic overview of Point Clear, Alabama.",
	}

	req := models.HookRequest{
		MainPost:   main,
		WidgetPost: widget,
	}

	hook, err := GenerateHook(context.Background(), nil, "test-model", nil, config.DefaultConfig(), req)
	if err != nil {
		t.Fatalf("GenerateHook failed: %v", err)
	}

	if hook.WidgetSlug != "widget" {
		t.Errorf("expected widget slug 'widget', got %s", hook.WidgetSlug)
	}
	if hook.WidgetTitle != "Point Clear Coastal Heritage" {
		t.Errorf("expected widget title, got %s", hook.WidgetTitle)
	}
	if hook.Brief == "" {
		t.Errorf("brief should not be empty")
	}
}

func TestSaveAndLoadHooks(t *testing.T) {
	tmpDir := t.TempDir()
	postDir := filepath.Join(tmpDir, "sample-post")
	if err := os.MkdirAll(postDir, 0755); err != nil {
		t.Fatalf("failed to create post dir: %v", err)
	}

	hooks := []*models.Hook{
		{
			WidgetSlug:  "related-one",
			WidgetTitle: "Related One",
			Brief:       "Discover how this connects to related topic one.",
			Generated:   time.Now(),
			ModelName:   "test-model",
		},
		{
			WidgetSlug:  "related-two",
			WidgetTitle: "Related Two",
			Brief:       "Deep dive into the financial records.",
			Generated:   time.Now(),
			ModelName:   "test-model",
		},
	}

	meta := models.HookMetadata{
		GeneratedAt:      time.Now().UTC(),
		GeneratorVersion: models.HookGeneratorVersion,
		PostSlug:         "sample-post",
	}
	if err := SaveHooks(postDir, meta, hooks); err != nil {
		t.Fatalf("SaveHooks failed: %v", err)
	}

	loaded, err := LoadHooks(postDir)
	if err != nil {
		t.Fatalf("LoadHooks failed: %v", err)
	}

	if len(loaded) != 2 {
		t.Fatalf("expected 2 loaded hooks, got %d", len(loaded))
	}
	if loaded[0].WidgetSlug != "related-one" || loaded[1].WidgetSlug != "related-two" {
		t.Errorf("unexpected loaded hooks: %+v", loaded)
	}
}

// --- New tests for features ported from Python hook_generator.py ---

func TestStripArticleNoise_HTMLElements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string // expected substring in output
		absent   string // substring that should NOT be in output
	}{
		{
			name:     "strips script tags with content",
			input:    "Hello <script>var x = 1;</script> World",
			contains: "Hello  World",
			absent:   "var x",
		},
		{
			name:     "strips style tags with content",
			input:    "Hello <style>.foo{color:red}</style> World",
			contains: "Hello  World",
			absent:   "color",
		},
		{
			name:     "strips img tags",
			input:    `Before <img src="photo.jpg" /> After`,
			contains: "Before  After",
			absent:   "photo.jpg",
		},
		{
			name:     "preserves anchor text",
			input:    `Click <a href="http://example.com">here for details</a> now`,
			contains: "here for details",
			absent:   "http://example.com",
		},
		{
			name:     "strips iframes",
			input:    `Before <iframe src="video.mp4"></iframe> After`,
			contains: "Before  After",
			absent:   "video.mp4",
		},
		{
			name:     "strips div with content",
			input:    `Text <div class="ad">ad content</div> more`,
			contains: "Text  more",
			absent:   "ad content",
		},
		{
			name:     "strips span with content",
			input:    `Before <span style="color:red">styled</span> after`,
			contains: "Before  after",
			absent:   "styled",
		},
		{
			name:     "strips shortcodes",
			input:    `Before {{< socialpreview url="foo" >}} after`,
			contains: "Before  after",
			absent:   "socialpreview",
		},
		{
			name:     "strips percent shortcodes",
			input:    `Before {{% highlight go %}} code {{% /highlight %}} after`,
			contains: "Before",
			absent:   "highlight",
		},
		{
			name:     "strips YAML front matter",
			input:    "---\ntitle: Test\n---\nActual content here",
			contains: "Actual content here",
			absent:   "title: Test",
		},
		{
			name:     "strips standalone URLs",
			input:    "Check out https://example.com/page for more info",
			contains: "Check out  for more info",
			absent:   "https://",
		},
		{
			name:     "collapses blank lines",
			input:    "Line 1\n\n\n\n\nLine 2",
			contains: "Line 1\n\nLine 2",
			absent:   "\n\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripArticleNoise(tt.input)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("expected output to contain %q, got %q", tt.contains, result)
			}
			if tt.absent != "" && strings.Contains(result, tt.absent) {
				t.Errorf("expected output to NOT contain %q, got %q", tt.absent, result)
			}
		})
	}
}

func TestCleanHookBrief_ReasoningStrip(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "strips think blocks",
			input:  "<think>internal reasoning</think>THE CONFLICT: Important facts",
			expect: "THE CONFLICT: Important facts",
		},
		{
			name:   "strips thinking blocks",
			input:  "<thinking>my analysis</thinking>THE CONFLICT: Data here",
			expect: "THE CONFLICT: Data here",
		},
		{
			name:   "strips THINK brackets",
			input:  "[THINK]reasoning[/THINK]THE CONFLICT: Result",
			expect: "THE CONFLICT: Result",
		},
		{
			name:   "strips reasoning preamble with 'we need'",
			input:  "We need to analyze this carefully. THE CONFLICT: The mayor...",
			expect: "THE CONFLICT: The mayor...",
		},
		{
			name:   "strips reasoning preamble with 'let's'",
			input:  "Let's look at the facts. THE CONFLICT: In 2020...",
			expect: "THE CONFLICT: In 2020...",
		},
		{
			name:   "strips reasoning with 'the instruction'",
			input:  "The instruction says to analyze. THE CONFLICT: Data",
			expect: "THE CONFLICT: Data",
		},
		{
			name:   "strips reasoning with 'draft:'",
			input:  "Draft: here is my attempt\nTHE CONFLICT: Facts",
			expect: "THE CONFLICT: Facts",
		},
		{
			name:   "keeps clean output untouched",
			input:  "THE CONFLICT: Mayor Jones in 2021...",
			expect: "THE CONFLICT: Mayor Jones in 2021...",
		},
		{
			name:   "strips chat template tokens",
			input:  "<|im_start|>THE CONFLICT: Data<|im_end|>",
			expect: "THE CONFLICT: Data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanHookBrief(tt.input)
			if result != tt.expect {
				t.Errorf("expected %q, got %q", tt.expect, result)
			}
		})
	}
}

func TestIsGarbageBrief(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		garbage bool
	}{
		{"empty", "", true},
		{"channel tokens", "THE CONFLICT <|channel> stuff", true},
		{"too few letters", "abc", true},
		{"no conflict marker", "This is a normal brief about important things that happened recently", true},
		{"digit spam", "THE CONFLICT 1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890 data", true},
		{"valid brief", "THE CONFLICT: Mayor Jones approved the zoning variance in January 2024, bypassing the standard review process.", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isGarbageBrief(tt.input); got != tt.garbage {
				t.Errorf("isGarbageBrief(%q) = %v, want %v", tt.input, got, tt.garbage)
			}
		})
	}
}

func TestParseBatchBriefs(t *testing.T) {
	tests := []struct {
		name     string
		response string
		expected int
		wantNil  bool
	}{
		{
			name: "two briefs",
			response: `=== BRIEF 1 ===
THE CONFLICT: Mayor Jones...

THE DISCOVERY: Article 2 reveals...

=== BRIEF 2 ===
THE CONFLICT: Senator Smith...

THE DISCOVERY: The report shows...`,
			expected: 2,
		},
		{
			name: "three briefs",
			response: `=== BRIEF 1 ===
First brief content

=== BRIEF 2 ===
Second brief content

=== BRIEF 3 ===
Third brief content`,
			expected: 3,
		},
		{
			name:     "count mismatch returns nil",
			response: "=== BRIEF 1 ===\nOnly one brief",
			expected: 3,
			wantNil:  true,
		},
		{
			name:     "no briefs returns nil",
			response: "just some random text",
			expected: 2,
			wantNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseBatchBriefs(tt.response, tt.expected)
			if tt.wantNil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
				return
			}
			if len(result) != tt.expected {
				t.Fatalf("expected %d briefs, got %d", tt.expected, len(result))
			}
			for i, b := range result {
				if b == "" {
					t.Errorf("brief %d is empty", i)
				}
			}
		})
	}
}

func TestDetermineTemporalFraming(t *testing.T) {
	now := time.Now()
	past := now.Add(-30 * 24 * time.Hour)
	future := now.Add(30 * 24 * time.Hour)

	tests := []struct {
		name       string
		mainDate   time.Time
		widgetDate time.Time
		expected   string
	}{
		{"widget older", now, past, "at the time"},
		{"widget newer", now, future, "since then"},
		{"same date", now, now, "at the time"},
		{"main zero", time.Time{}, now, "today"},
		{"widget zero", now, time.Time{}, "today"},
		{"both zero", time.Time{}, time.Time{}, "today"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineTemporalFraming(tt.mainDate, tt.widgetDate)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestBuildMetadataContext(t *testing.T) {
	mainFM := &models.FrontMatter{
		Title:   "Main Article Title",
		AltTags: "Photo of courthouse",
		Tags:    []string{"corruption", "local-government"},
		Date:    time.Now(),
	}
	widgetFM := &models.FrontMatter{
		Title: "Widget Article Title",
		Tags:  []string{"investigation"},
		Date:  time.Now(),
	}

	result := buildMetadataContext(mainFM, widgetFM)

	if !strings.Contains(result, "[MAIN ARTICLE METADATA]") {
		t.Error("expected [MAIN ARTICLE METADATA] section")
	}
	if !strings.Contains(result, "Title: Main Article Title") {
		t.Error("expected main title")
	}
	if !strings.Contains(result, "Image Context: Photo of courthouse") {
		t.Error("expected main alttags")
	}
	if !strings.Contains(result, "Key Topics: corruption, local-government") {
		t.Error("expected main tags")
	}
	if !strings.Contains(result, "[WIDGET ARTICLE METADATA]") {
		t.Error("expected [WIDGET ARTICLE METADATA] section")
	}
	if !strings.Contains(result, "Title: Widget Article Title") {
		t.Error("expected widget title")
	}
}

func TestIsConnectionError(t *testing.T) {
	tests := []struct {
		name   string
		errMsg string
		expect bool
	}{
		{"connection refused", "dial tcp 127.0.0.1:8080: connection refused", true},
		{"no such host", "dial tcp: lookup invalid.host: no such host", true},
		{"network unreachable", "dial tcp: network is unreachable", true},
		{"dial tcp", "dial tcp 10.0.0.1:8080: i/o timeout", true},
		{"timeout", "request timed out after 5 minutes", false},
		{"generic error", "model returned garbage", false},
		{"nil error", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.errMsg != "" {
				err = &testErr{msg: tt.errMsg}
			}
			if got := isConnectionError(err); got != tt.expect {
				t.Errorf("isConnectionError(%q) = %v, want %v", tt.errMsg, got, tt.expect)
			}
		})
	}
}

type testErr struct{ msg string }

func (e *testErr) Error() string { return e.msg }

func TestRepetitionPenaltyInCompletionRequest(t *testing.T) {
	// Verify the RepetitionPenalty field exists and can be set on CompletionRequest
	req := models.CompletionRequest{
		Prompt:            "test",
		Temperature:       1.0,
		RepetitionPenalty: 1.25,
	}
	if req.RepetitionPenalty != 1.25 {
		t.Errorf("expected RepetitionPenalty 1.25, got %f", req.RepetitionPenalty)
	}
}
