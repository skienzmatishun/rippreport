package processor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/rippreport/go-content-pipeline/internal/client"
	"github.com/rippreport/go-content-pipeline/internal/models"
)

const (
	sidebarHooksFilename = "sidebar-hooks.yaml"
	hookMaxPromptTokens  = 20000
	hookMaxOutputTokens  = 4096
)

var (
	thinkBlockRe     = regexp.MustCompile(`(?is)<think>.*?</think>|<thinking>.*?</thinking>|\[THINK\].*?\[/THINK\]|<thought>.*?</thought>`)
	channelPairRe    = regexp.MustCompile(`(?s)<\|channel\>[^<]*<channel\|>`)
	channelTokenRe   = regexp.MustCompile(`<\|/?channel\>[^\n]*|<channel\|>|<\|im_start\|>|<\|im_end\|>|<\|endoftext\|>|<\|turn\>|</?turn\|>`)
	shortcodeRe      = regexp.MustCompile(`\{\{<\s*[^>]*>\}\}|\{\{%\s*[^%]*%\}\}`)
	htmlTagRe        = regexp.MustCompile(`(?is)<[^>]+>`)
	urlRe            = regexp.MustCompile(`https?://\S+`)
	extraBlankLineRe = regexp.MustCompile(`\n{3,}`)
)

// GenerateHook produces a Bridge Brief connecting two articles, matching
// recent_placement/hook_generator (THE CONFLICT / DISCOVERY / PIVOT / KEY ENTITIES).
func GenerateHook(ctx context.Context, llamaClient client.LlamaClient, modelName string, cache *CompressionCache, req models.HookRequest) (*models.Hook, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid hook request: %w", err)
	}

	mainContent := stripArticleNoise(req.MainPost.Body)
	widgetContent := stripArticleNoise(req.WidgetPost.Body)

	// Prefer cached compressed articles when available (no LLM call needed).
	// Fall back to on-the-fly LLM compression only when --use-compressed is set
	// and no cache entry exists. Otherwise truncate.
	mainContent = resolveContent(ctx, llamaClient, cache, req.MainPost, mainContent, req.UseCompressed)
	widgetContent = resolveContent(ctx, llamaClient, cache, req.WidgetPost, widgetContent, req.UseCompressed)

	mainFM := req.MainPost.FrontMatter
	widgetFM := req.WidgetPost.FrontMatter

	brief := stubBridgeBrief(mainFM.Title, widgetFM.Title)

	if llamaClient != nil {
		systemPrompt, userPrompt := buildBridgeBriefPrompts(req.MainPost, req.WidgetPost, mainContent, widgetContent)
		resp, err := llamaClient.Complete(ctx, models.CompletionRequest{
			Messages: []models.ChatMessage{
				{Role: "system", Content: systemPrompt},
				{Role: "user", Content: userPrompt},
			},
			Model:       modelName,
			Temperature: 1.0,
			TopP:        0.95,
			TopK:        20,
			MaxTokens:   hookMaxOutputTokens,
		})
		if err != nil {
			return nil, fmt.Errorf("hook completion failed for %s -> %s: %w", req.MainPost.Slug, req.WidgetPost.Slug, err)
		}
		cleaned := cleanHookBrief(resp.Text)
		if isGarbageBrief(cleaned) {
			return nil, fmt.Errorf("hook model returned unusable brief for %s -> %s", req.MainPost.Slug, req.WidgetPost.Slug)
		}
		brief = cleaned
	}

	now := time.Now().UTC()
	hook := &models.Hook{
		WidgetSlug:  req.WidgetPost.Slug,
		WidgetTitle: widgetFM.Title,
		Brief:       brief,
		WidgetDate:  formatDate(widgetFM.Date),
		Generated:   now,
		ModelName:   modelName,
	}
	if err := hook.Validate(); err != nil {
		return nil, fmt.Errorf("generated hook validation failed: %w", err)
	}
	return hook, nil
}

func buildBridgeBriefPrompts(main, widget *models.Post, mainContent, widgetContent string) (systemPrompt, userPrompt string) {
	mainDate := main.FrontMatter.Date
	widgetDate := widget.FrontMatter.Date

	discoveryInstruction := "2. THE DISCOVERY: What Article 2 revealed"
	temporalGuidance := ""
	if !mainDate.IsZero() && !widgetDate.IsZero() {
		if widgetDate.Before(mainDate) {
			discoveryInstruction = "2. THE DISCOVERY: What Article 2 showed or revealed at the time (background/context)"
			temporalGuidance = "Article 2 is OLDER than Article 1. Frame it as background/historical context."
		} else if widgetDate.After(mainDate) {
			discoveryInstruction = "2. THE DISCOVERY: What Article 2 proves or shows since then (update)"
			temporalGuidance = "Article 2 is NEWER than Article 1. Frame it as an update."
		} else {
			discoveryInstruction = "2. THE DISCOVERY: What Article 2 showed or revealed at the time (background/context)"
			temporalGuidance = "Article 2 is OLDER than Article 1. Frame it as background/historical context."
		}
	}

	systemPrompt = "You are a Private Investigator analyzing two articles for The Ripp Report. " +
		"Your job is to extract the raw facts that connect these articles - NOT to write polished prose.\n\n"
	if temporalGuidance != "" {
		systemPrompt += "TEMPORAL CONTEXT: " + temporalGuidance + "\n\n"
	}
	systemPrompt += "OUTPUT A BRIDGE BRIEF containing:\n" +
		"1. THE CONFLICT: What happened in the past? (Include specific names, dates, events)\n" +
		discoveryInstruction + "\n" +
		"3. THE PIVOT: Why is Article 2 the 'smoking gun' or key context for Article 1?\n" +
		"4. KEY ENTITIES: List specific names, organizations, or locations to mention\n\n" +
		"RULES:\n" +
		"- Output raw facts as a structured list. No flowery language needed.\n" +
		"- Use specific details: names, dates, dollar amounts, positions/titles.\n" +
		"- Start directly with 'THE CONFLICT:' - NO preamble, NO reasoning, NO explanation of your process"

	var meta strings.Builder
	meta.WriteString("\n[MAIN ARTICLE METADATA]\n")
	meta.WriteString("Title: " + main.FrontMatter.Title + "\n")
	if main.FrontMatter.AltTags != "" {
		meta.WriteString("Image Context: " + main.FrontMatter.AltTags + "\n")
	}
	if len(main.FrontMatter.Tags) > 0 {
		meta.WriteString("Key Topics: " + strings.Join(main.FrontMatter.Tags, ", ") + "\n")
	}

	meta.WriteString("\n[WIDGET ARTICLE METADATA]\n")
	meta.WriteString("Title: " + widget.FrontMatter.Title + "\n")
	if widget.FrontMatter.AltTags != "" {
		meta.WriteString("Image Context: " + widget.FrontMatter.AltTags + "\n")
	}
	if len(widget.FrontMatter.Tags) > 0 {
		meta.WriteString("Key Topics: " + strings.Join(widget.FrontMatter.Tags, ", ") + "\n")
	}

	userPrompt = meta.String() + `
[MAIN ARTICLE CONTENT]
` + mainContent + `

[WIDGET ARTICLE CONTENT]
` + widgetContent + `

[INSTRUCTION]
Analyze the Main Article and the Widget Article. Produce a Bridge Brief containing:
The Conflict: What happened in the past? (Names/Dates)
The Discovery: What does the second article prove today?
The Pivot: Why is Article 2 the 'smoking gun' for Article 1?
Key Entities: A list of the specific names (Lowery, Wilters, etc.) to mention.

Output this as a raw list. No flowery language needed.`

	return systemPrompt, userPrompt
}

// resolveContent picks the best available content for a post:
//  1. Cached compressed article (always preferred, no LLM call)
//  2. On-the-fly LLM compression (only when useCompressed is true)
//  3. Head-tail truncation of raw content (fallback)
func resolveContent(ctx context.Context, llamaClient client.LlamaClient, cache *CompressionCache, post *models.Post, stripped string, useCompressed bool) string {
	// 1. Check cache first — free, no LLM call needed
	if cache != nil {
		if cached, ok := cache.GetCached(post.Slug); ok && cached != "" {
			return cached
		}
	}

	// 2. On-the-fly LLM compression when explicitly requested
	if useCompressed && llamaClient != nil && cache != nil {
		compressed, err := Compress(ctx, llamaClient, cache, post, 800)
		if err == nil {
			return compressed
		}
	}

	// 3. Fallback: truncate raw content
	return Truncate(stripped, hookMaxPromptTokens)
}

func stubBridgeBrief(mainTitle, widgetTitle string) string {
	return "THE CONFLICT:\n- " + mainTitle + "\n\nTHE DISCOVERY:\n- " + widgetTitle +
		"\n\nTHE PIVOT:\n- The related article supplies key context for the main article.\n\nKEY ENTITIES:\n- See both articles"
}

func stripArticleNoise(text string) string {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "---") {
		parts := strings.SplitN(text, "---", 3)
		if len(parts) >= 3 {
			text = parts[2]
		}
	}
	text = shortcodeRe.ReplaceAllString(text, "")
	text = htmlTagRe.ReplaceAllString(text, "")
	text = urlRe.ReplaceAllString(text, "")
	text = extraBlankLineRe.ReplaceAllString(text, "\n\n")
	return strings.TrimSpace(text)
}

func cleanHookBrief(text string) string {
	text = thinkBlockRe.ReplaceAllString(text, "")
	if idx := strings.LastIndex(text, "<channel|>"); idx != -1 {
		text = text[idx+len("<channel|>"):]
	}
	text = channelPairRe.ReplaceAllString(text, "")
	text = channelTokenRe.ReplaceAllString(text, "")

	upper := strings.ToUpper(text)
	if i := strings.Index(upper, "THE CONFLICT"); i > 0 {
		before := strings.ToLower(strings.TrimSpace(text[:i]))
		for _, indicator := range []string{
			"we need", "need to", "let's", "i think", "the instruction",
			"user says", "the prompt", "draft:", "analyze",
		} {
			if strings.Contains(before, indicator) {
				text = text[i:]
				break
			}
		}
	}

	return strings.TrimSpace(text)
}

func isGarbageBrief(text string) bool {
	if text == "" {
		return true
	}
	if strings.Contains(text, "<|channel>") || strings.Contains(text, "<channel|>") {
		return true
	}
	letters, digits := 0, 0
	for _, r := range text {
		if unicode.IsLetter(r) {
			letters++
		}
		if unicode.IsDigit(r) {
			digits++
		}
	}
	if letters < 40 {
		return true
	}
	if digits > letters*3 {
		return true
	}
	if !strings.Contains(strings.ToUpper(text), "THE CONFLICT") {
		return true
	}
	return false
}

// SaveHooks writes production sidebar-hooks.yaml (metadata + slug-keyed briefs).
func SaveHooks(postDir string, meta models.HookMetadata, hooks []*models.Hook) error {
	if postDir == "" {
		return fmt.Errorf("post directory cannot be empty")
	}
	if meta.PostSlug == "" {
		return fmt.Errorf("post slug cannot be empty")
	}
	if meta.GeneratorVersion == "" {
		meta.GeneratorVersion = models.HookGeneratorVersion
	}
	if meta.GeneratedAt.IsZero() {
		meta.GeneratedAt = time.Now().UTC()
	}
	hookFile := models.HookFile{Metadata: meta, Hooks: hooks}
	if err := hookFile.Validate(); err != nil {
		return fmt.Errorf("hook file validation failed: %w", err)
	}
	return writeHookFile(filepath.Join(postDir, sidebarHooksFilename), meta, hooks)
}

// LoadHooks reads sidebar-hooks.yaml (production map format or legacy list).
func LoadHooks(postDir string) ([]*models.Hook, error) {
	file, err := LoadHookFile(filepath.Join(postDir, sidebarHooksFilename))
	if err != nil {
		return nil, err
	}
	return file.Hooks, nil
}

// LoadHookFile reads a sidebar-hooks.yaml (or staged) document.
func LoadHookFile(path string) (*models.HookFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read sidebar hooks file %s: %w", path, err)
	}
	meta, hooks, err := parseHookYAML(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sidebar hooks YAML: %w", err)
	}
	if meta.GeneratorVersion == "" {
		meta.GeneratorVersion = models.HookGeneratorVersion
	}
	if meta.GeneratedAt.IsZero() {
		meta.GeneratedAt = time.Now().UTC()
	}
	file := &models.HookFile{Metadata: meta, Hooks: hooks}
	if err := file.Validate(); err != nil {
		return nil, fmt.Errorf("sidebar hooks validation failed: %w", err)
	}
	return file, nil
}
