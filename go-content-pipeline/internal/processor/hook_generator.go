package processor

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/rippreport/go-content-pipeline/internal/client"
	"github.com/rippreport/go-content-pipeline/internal/config"
	"github.com/rippreport/go-content-pipeline/internal/models"
)

const (
	sidebarHooksFilename = "sidebar-hooks.yaml"
	hookMaxPromptTokens  = 20000
	hookMaxOutputTokens  = 30000

	// Retry configuration matching Python's hook_generator.py
	hookMaxRetries      = 3
	hookRetryBackoffSec = 1 // base seconds for exponential backoff

	// Batch prompt constants matching Python's generate_hook_batch
	hookBatchMaxOutputTokens = 30000

	// RepetitionPenalty matching Python's REPETITION_PENALTY = 1.25
	hookRepetitionPenalty float32 = 1.25
)

var (
	thinkBlockRe   = regexp.MustCompile(`(?is)<think>.*?</think>|<thinking>.*?</thinking>|\[THINK\].*?\[/THINK\]|<thought>.*?</thought>`)
	channelPairRe  = regexp.MustCompile(`(?s)<\|channel>[^<]*<channel\|>`)
	channelTokenRe = regexp.MustCompile(`<\|/?channel>[^\n]*|<channel\|>|<\|im_start\|>|<\|im_end\|>|<\|endoftext\|>|<\|turn>|</?turn\|>`)

	// Content stripping patterns matching Python's strip_shortcodes:
	// Self-closing shortcodes: {{< shortcode >}}
	shortcodeSelfClosingRe = regexp.MustCompile(`\{\{<\s*[^>]+>\}\}`)
	// Percent-delimited shortcodes: {{% shortcode %}}
	shortcodePctRe = regexp.MustCompile(`\{\{%\s*[^%]+%\}\}`)
	// Paired shortcodes opener: {{< figure ... >}} — name captured in group 1
	shortcodePairedOpenRe = regexp.MustCompile(`\{\{[<%]\s*(\w+)`)

	// HTML element stripping matching Python's strip_shortcodes
	iframePairedRe       = regexp.MustCompile(`(?is)<iframe[^>]*>.*?</iframe>`)
	iframeSelfClosingRe  = regexp.MustCompile(`(?i)<iframe[^>]*/?>`)
	imgTagRe             = regexp.MustCompile(`(?i)<img[^>]*/?>`)
	anchorOpenRe         = regexp.MustCompile(`(?i)<a[^>]*>`)
	anchorCloseRe        = regexp.MustCompile(`(?i)</a>`)
	divPairedRe          = regexp.MustCompile(`(?is)<div[^>]*>.*?</div>`)
	scriptPairedRe       = regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	stylePairedRe        = regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	spanPairedRe         = regexp.MustCompile(`(?is)<span[^>]*>.*?</span>`)
	htmlTagRe            = regexp.MustCompile(`(?is)<[^>]+>`)
	urlRe                = regexp.MustCompile(`https?://\S+`)
	extraBlankLineRe     = regexp.MustCompile(`\n{3,}`)
	trailingWhitespaceRe = regexp.MustCompile(`(?m)[ \t]+$`)
)

// Reasoning indicators matching Python's strip_thinking function.
var reasoningIndicators = []string{
	"we need", "need answer", "need produce", "need output",
	"must output", "need analyze", "need identify", "need extract",
	"need ensure", "need include", "need connect", "need maybe",
	"let's", "i think", "could use", "should not", "probably",
	"the instruction", "user says", "user asks", "the prompt",
	"need decide", "need final", "draft:", "potential structure",
}

// GenerateHook produces a Bridge Brief connecting two articles, matching
// recent_placement/hook_generator (THE CONFLICT / DISCOVERY / PIVOT / KEY ENTITIES).
//
// Includes retry logic with exponential backoff (up to 3 attempts), matching
// Python's generate_hook.
func GenerateHook(ctx context.Context, llamaClient client.LlamaClient, modelName string, cache *CompressionCache, cfg *config.Config, req models.HookRequest) (*models.Hook, error) {
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

		// Retry logic with exponential backoff matching Python's generate_hook
		var lastErr error
		for attempt := 0; attempt < hookMaxRetries; attempt++ {
			resp, err := llamaClient.Complete(ctx, models.CompletionRequest{
				Messages: []models.ChatMessage{
					{Role: "system", Content: systemPrompt},
					{Role: "user", Content: userPrompt},
				},
				Model:             modelName,
				Temperature:       1.0,
				TopP:              0.95,
				TopK:              20,
				MaxTokens:         hookMaxOutputTokens,
				RepetitionPenalty: hookRepetitionPenalty,
			})
			if err != nil {
				lastErr = err
				// Don't retry connection errors — fail fast (matching Python)
				if hookCompletionCannotRetry(err) {
					return nil, fmt.Errorf("hook completion failed for %s -> %s: %w", req.MainPost.Slug, req.WidgetPost.Slug, err)
				}
				if attempt < hookMaxRetries-1 {
					backoff := time.Duration(hookRetryBackoffSec*(1<<uint(attempt))) * time.Second
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					case <-time.After(backoff):
					}
				}
				continue
			}

			cleaned := cleanHookBrief(resp.Text)
			if isGarbageBrief(cleaned) {
				lastErr = fmt.Errorf("hook model returned unusable brief for %s -> %s", req.MainPost.Slug, req.WidgetPost.Slug)
				if attempt < hookMaxRetries-1 {
					backoff := time.Duration(hookRetryBackoffSec*(1<<uint(attempt))) * time.Second
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					case <-time.After(backoff):
					}
				}
				continue
			}

			brief = cleaned
			lastErr = nil
			break
		}

		if lastErr != nil {
			return nil, fmt.Errorf("hook generation failed for %s -> %s after %d attempts: %w",
				req.MainPost.Slug, req.WidgetPost.Slug, hookMaxRetries, lastErr)
		}
	}

	now := time.Now().UTC()

	var llmText string
	if llamaClient != nil && brief != stubBridgeBrief(mainFM.Title, widgetFM.Title) {
		var err error
		llmText, err = generateHookText(ctx, llamaClient, modelName, cfg, brief)
		if err != nil {
			return nil, fmt.Errorf("hook text for %s -> %s: %w", req.MainPost.Slug, req.WidgetPost.Slug, err)
		}
	}

	hook := &models.Hook{
		WidgetSlug:  req.WidgetPost.Slug,
		WidgetTitle: widgetFM.Title,
		Brief:       brief,
		LLMText:     llmText,
		WidgetDate:  formatDate(widgetFM.Date),
		Generated:   now,
		ModelName:   modelName,
	}
	if err := hook.Validate(); err != nil {
		return nil, fmt.Errorf("generated hook validation failed: %w", err)
	}
	return hook, nil
}

// generateHookText turns a bridge brief into the initial sidebar copy.
func generateHookText(ctx context.Context, llamaClient client.LlamaClient, modelName string, cfg *config.Config, brief string) (string, error) {
	hookSystemPrompt := `You are a sidebar hook writer for The Ripp Report.

Your task: Transform the structured bridge brief into a compelling 2-sentence sidebar hook.

Style:
- Direct and conversational
- Use specific names and details from the brief
- Connect the widget article to why readers should care
- 220-320 characters total
- NO passive voice

Output ONLY the hook text, nothing else.`

	hookUserPrompt := fmt.Sprintf("Create a sidebar hook from this bridge brief:\n\n%s", brief)

	maxTokens := cfg.Processing.HookLLMTextMaxTokens
	if maxTokens == 0 {
		maxTokens = 30000 // Default if not configured
	}

	hookResp, err := llamaClient.Complete(ctx, models.CompletionRequest{
		Messages: []models.ChatMessage{
			{Role: "system", Content: hookSystemPrompt},
			{Role: "user", Content: hookUserPrompt},
		},
		Model:       modelName,
		Temperature: 0.9,
		TopP:        0.95,
		MaxTokens:   maxTokens,
	})
	if err != nil {
		return "", fmt.Errorf("sidebar text completion failed: %w", err)
	}
	text := strings.TrimSpace(hookResp.Text)
	if text == "" {
		return "", fmt.Errorf("sidebar text completion returned empty text")
	}
	return text, nil
}

// GenerateHookBatchSinglePrompt generates multiple Bridge Briefs in a single
// batched LLM request, matching Python's generate_hook_batch. All pairs are
// assembled into one prompt with === PAIR N === markers and the response is
// parsed using === BRIEF N === markers.
//
// Falls back to individual GenerateHook calls if parsing fails or counts
// don't match.
func GenerateHookBatchSinglePrompt(ctx context.Context, llamaClient client.LlamaClient, modelName string, cache *CompressionCache, cfg *config.Config, reqs []models.HookRequest) ([]*models.Hook, error) {
	if len(reqs) == 0 {
		return nil, nil
	}

	// Single request — just use the regular path
	if len(reqs) == 1 || llamaClient == nil {
		return generateHooksIndividually(ctx, llamaClient, modelName, cache, cfg, reqs)
	}

	// Prepare each pair's content
	type preparedPair struct {
		mainContent   string
		widgetContent string
		mainSlug      string
		widgetSlug    string
		mainFM        *models.FrontMatter
		widgetFM      *models.FrontMatter
		temporalFrame string
	}

	pairs := make([]preparedPair, len(reqs))
	for i, req := range reqs {
		mainContent := stripArticleNoise(req.MainPost.Body)
		widgetContent := stripArticleNoise(req.WidgetPost.Body)
		mainContent = resolveContent(ctx, llamaClient, cache, req.MainPost, mainContent, req.UseCompressed)
		widgetContent = resolveContent(ctx, llamaClient, cache, req.WidgetPost, widgetContent, req.UseCompressed)

		temporal := determineTemporalFraming(req.MainPost.FrontMatter.Date, req.WidgetPost.FrontMatter.Date)

		pairs[i] = preparedPair{
			mainContent:   mainContent,
			widgetContent: widgetContent,
			mainSlug:      req.MainPost.Slug,
			widgetSlug:    req.WidgetPost.Slug,
			mainFM:        req.MainPost.FrontMatter,
			widgetFM:      req.WidgetPost.FrontMatter,
			temporalFrame: temporal,
		}
	}

	// Build the batch prompt (matching Python's generate_hook_batch)
	var batchItems strings.Builder
	for i, p := range pairs {
		metaCtx := buildMetadataContext(p.mainFM, p.widgetFM)
		fmt.Fprintf(&batchItems, "\n=== PAIR %d: %s → %s ===\n", i+1, p.mainSlug, p.widgetSlug)
		fmt.Fprintf(&batchItems, "[TEMPORAL CONTEXT: Article 2 is from %s]\n", p.temporalFrame)
		batchItems.WriteString(metaCtx)
		fmt.Fprintf(&batchItems, "\n[MAIN ARTICLE CONTENT]\n%s\n\n[WIDGET ARTICLE CONTENT]\n%s\n", p.mainContent, p.widgetContent)
	}

	systemPrompt := "You are a Private Investigator analyzing article pairs for The Ripp Report. " +
		"For each pair, extract the raw facts that connect them - NOT polished prose.\n\n" +
		"For EACH pair, output a Bridge Brief containing:\n" +
		"1. THE CONFLICT: What happened in the past? (Include specific names, dates, events)\n" +
		"2. THE DISCOVERY: What Article 2 revealed (adjust based on TEMPORAL CONTEXT)\n" +
		"3. THE PIVOT: Why is Article 2 the 'smoking gun' or key context for Article 1?\n" +
		"4. KEY ENTITIES: List specific names, organizations, or locations to mention\n\n" +
		"IMPORTANT: Each pair includes TEMPORAL CONTEXT showing when Article 2 is from:\n" +
		"- If 'at the time': Article 2 is older - frame as background/context\n" +
		"- If 'since then': Article 2 is newer - frame as update/development\n" +
		"- Adjust THE DISCOVERY section accordingly\n\n" +
		"CRITICAL FORMAT:\n" +
		"Output each brief with a header like:\n" +
		"=== BRIEF 1 ===\n" +
		"[brief content]\n\n" +
		"=== BRIEF 2 ===\n" +
		"[brief content]\n\n" +
		"etc.\n\n" +
		"RULES:\n" +
		"- Output raw facts as a structured list. No flowery language needed.\n" +
		"- Use specific details: names, dates, dollar amounts, positions/titles.\n" +
		"- Keep each brief focused and concise."

	userPrompt := batchItems.String() + "\n[INSTRUCTION]\nGenerate a Bridge Brief for each pair above. Use the format:\n=== BRIEF 1 ===\n[content]\n\n=== BRIEF 2 ===\n[content]\netc."

	// Retry logic with exponential backoff
	var lastErr error
	for attempt := 0; attempt < hookMaxRetries; attempt++ {
		resp, err := llamaClient.Complete(ctx, models.CompletionRequest{
			Messages: []models.ChatMessage{
				{Role: "system", Content: systemPrompt},
				{Role: "user", Content: userPrompt},
			},
			Model:             modelName,
			Temperature:       1.0,
			TopP:              0.95,
			TopK:              20,
			MaxTokens:         hookBatchMaxOutputTokens,
			RepetitionPenalty: hookRepetitionPenalty,
		})
		if err != nil {
			lastErr = err
			if hookCompletionCannotRetry(err) {
				return nil, fmt.Errorf("batch hook completion failed: %w", err)
			}
			if attempt < hookMaxRetries-1 {
				backoff := time.Duration(hookRetryBackoffSec*(1<<uint(attempt))) * time.Second
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(backoff):
				}
			}
			continue
		}

		if resp.Text == "" {
			lastErr = fmt.Errorf("no response returned for batch")
			if attempt < hookMaxRetries-1 {
				backoff := time.Duration(hookRetryBackoffSec*(1<<uint(attempt))) * time.Second
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(backoff):
				}
			}
			continue
		}

		// Parse the batched response: === BRIEF N === markers
		briefs := parseBatchBriefs(resp.Text, len(reqs))
		if len(briefs) != len(reqs) {
			lastErr = fmt.Errorf("expected %d briefs but got %d, falling back to individual", len(reqs), len(briefs))
			if attempt < hookMaxRetries-1 {
				backoff := time.Duration(hookRetryBackoffSec*(1<<uint(attempt))) * time.Second
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(backoff):
				}
			}
			continue
		}

		// Build Hook objects from parsed briefs
		now := time.Now().UTC()
		hooks := make([]*models.Hook, 0, len(reqs))
		for i, briefText := range briefs {
			cleaned := cleanHookBrief(briefText)
			if isGarbageBrief(cleaned) {
				// For individual garbage briefs, use the stub
				cleaned = stubBridgeBrief(pairs[i].mainFM.Title, pairs[i].widgetFM.Title)
			}
			var llmText string
			if cleaned != stubBridgeBrief(pairs[i].mainFM.Title, pairs[i].widgetFM.Title) {
				var err error
				llmText, err = generateHookText(ctx, llamaClient, modelName, cfg, cleaned)
				if err != nil {
					return nil, fmt.Errorf("hook text for %s -> %s: %w", pairs[i].mainSlug, pairs[i].widgetSlug, err)
				}
			}
			hooks = append(hooks, &models.Hook{
				WidgetSlug:  pairs[i].widgetSlug,
				WidgetTitle: pairs[i].widgetFM.Title,
				Brief:       cleaned,
				LLMText:     llmText,
				WidgetDate:  formatDate(pairs[i].widgetFM.Date),
				Generated:   now,
				ModelName:   modelName,
			})
		}

		return hooks, nil
	}

	// All batch attempts failed — fall back to individual generation
	fmt.Printf("    ⚠ batch generation failed (%v), falling back to individual hooks\n", lastErr)
	return generateHooksIndividually(ctx, llamaClient, modelName, cache, cfg, reqs)
}

// parseBatchBriefs extracts briefs from a batch response using === BRIEF N === markers.
// Matching Python's regex: r'=== BRIEF (\d+) ===\s*(.*?)(?=\n=== BRIEF \d+ ===|\Z)'
var batchBriefHeaderRe = regexp.MustCompile(`=== BRIEF \d+ ===`)

func parseBatchBriefs(response string, expectedCount int) []string {
	// Find all header positions
	locs := batchBriefHeaderRe.FindAllStringIndex(response, -1)
	if len(locs) != expectedCount {
		return nil
	}

	briefs := make([]string, len(locs))
	for i, loc := range locs {
		// Content starts after the header
		start := loc[1]
		// Content ends at next header or end of string
		var end int
		if i+1 < len(locs) {
			end = locs[i+1][0]
		} else {
			end = len(response)
		}
		briefs[i] = strings.TrimSpace(response[start:end])
	}
	return briefs
}

// generateHooksIndividually generates hooks one-by-one, used as fallback.
func generateHooksIndividually(ctx context.Context, llamaClient client.LlamaClient, modelName string, cache *CompressionCache, cfg *config.Config, reqs []models.HookRequest) ([]*models.Hook, error) {
	hooks := make([]*models.Hook, 0, len(reqs))
	var firstErr error
	for i, req := range reqs {
		fmt.Printf("      [%d/%d] → %s...", i+1, len(reqs), req.WidgetPost.Slug)
		hook, hErr := GenerateHook(ctx, llamaClient, modelName, cache, cfg, req)
		if hErr != nil {
			fmt.Printf(" ✗ %v\n", hErr)
			if firstErr == nil {
				firstErr = hErr
			}
			continue
		}
		fmt.Printf(" ✓\n")
		hooks = append(hooks, hook)
	}
	if len(hooks) == 0 {
		if firstErr != nil {
			return nil, firstErr
		}
		return nil, fmt.Errorf("no hooks generated")
	}
	return hooks, nil
}

// determineTemporalFraming returns the temporal framing string matching
// Python's temporal_framing logic.
func determineTemporalFraming(mainDate, widgetDate time.Time) string {
	if mainDate.IsZero() || widgetDate.IsZero() {
		return "today"
	}
	if widgetDate.Before(mainDate) {
		return "at the time"
	}
	if widgetDate.After(mainDate) {
		return "since then"
	}
	return "at the time"
}

// buildMetadataContext builds the [MAIN ARTICLE METADATA] / [WIDGET ARTICLE METADATA]
// block matching Python's generate_hook_batch and generate_hook.
func buildMetadataContext(mainFM, widgetFM *models.FrontMatter) string {
	var meta strings.Builder

	if mainFM.Title != "" || mainFM.AltTags != "" || len(mainFM.Tags) > 0 {
		meta.WriteString("\n[MAIN ARTICLE METADATA]")
		if mainFM.Title != "" {
			meta.WriteString("\nTitle: " + mainFM.Title)
		}
		if mainFM.AltTags != "" {
			meta.WriteString("\nImage Context: " + mainFM.AltTags)
		}
		if len(mainFM.Tags) > 0 {
			meta.WriteString("\nKey Topics: " + strings.Join(mainFM.Tags, ", "))
		}
	}

	if widgetFM.Title != "" || widgetFM.AltTags != "" || len(widgetFM.Tags) > 0 {
		meta.WriteString("\n\n[WIDGET ARTICLE METADATA]")
		if widgetFM.Title != "" {
			meta.WriteString("\nTitle: " + widgetFM.Title)
		}
		if widgetFM.AltTags != "" {
			meta.WriteString("\nImage Context: " + widgetFM.AltTags)
		}
		if len(widgetFM.Tags) > 0 {
			meta.WriteString("\nKey Topics: " + strings.Join(widgetFM.Tags, ", "))
		}
	}

	return meta.String()
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

	metaCtx := buildMetadataContext(main.FrontMatter, widget.FrontMatter)

	userPrompt = metaCtx + `
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

// stripArticleNoise removes Hugo shortcodes, HTML elements, and non-content
// markup from article text. Matches Python's strip_shortcodes exactly:
//   - YAML front matter removal
//   - Hugo shortcodes ({{< >}} and {{% %}})
//   - HTML iframes, img, div, script, style, span (content-aware removal)
//   - Anchor tags (preserves link text)
//   - All remaining HTML tags
//   - URLs
//   - Excessive blank lines
//   - Trailing whitespace per line
func stripArticleNoise(text string) string {
	text = strings.TrimSpace(text)

	// Remove YAML front matter if present
	if strings.HasPrefix(text, "---") {
		parts := strings.SplitN(text, "---", 3)
		if len(parts) >= 3 {
			text = parts[2]
		}
	}

	// Remove Hugo shortcodes (matching Python exactly)
	text = shortcodeSelfClosingRe.ReplaceAllString(text, "")
	text = stripPairedShortcodes(text)
	text = shortcodePctRe.ReplaceAllString(text, "")

	// Remove HTML elements with content (order matters — remove content-bearing first)
	text = iframePairedRe.ReplaceAllString(text, "")
	text = iframeSelfClosingRe.ReplaceAllString(text, "")
	text = imgTagRe.ReplaceAllString(text, "")

	// Remove anchor tags but keep the link text (matching Python)
	text = anchorOpenRe.ReplaceAllString(text, "")
	text = anchorCloseRe.ReplaceAllString(text, "")

	// Remove paired elements with their content
	text = divPairedRe.ReplaceAllString(text, "")
	text = scriptPairedRe.ReplaceAllString(text, "")
	text = stylePairedRe.ReplaceAllString(text, "")
	text = spanPairedRe.ReplaceAllString(text, "")

	// Strip any remaining HTML tags (catch-all for p, strong, em, br, hr, etc.)
	text = htmlTagRe.ReplaceAllString(text, "")

	// Remove standalone URLs
	text = urlRe.ReplaceAllString(text, "")

	// Normalize whitespace (matching Python)
	text = extraBlankLineRe.ReplaceAllString(text, "\n\n")
	text = trailingWhitespaceRe.ReplaceAllString(text, "")

	// One more pass to catch remaining multiple blank lines
	text = extraBlankLineRe.ReplaceAllString(text, "\n\n")

	return strings.TrimSpace(text)
}

// stripPairedShortcodes removes paired Hugo shortcodes like:
//
//	{{< figure ... >}}...{{< /figure >}}
//	{{% highlight go %}}...{{% /highlight %}}
//
// Go's regexp doesn't support backreferences, so we scan for openers,
// extract the shortcode name, and remove everything up to the matching closer.
func stripPairedShortcodes(text string) string {
	for {
		loc := shortcodePairedOpenRe.FindStringSubmatchIndex(text)
		if loc == nil {
			break
		}
		name := text[loc[2]:loc[3]]
		// Build the closer pattern for this specific shortcode name
		closerAngle := "{{< /" + name + " >}}"
		closerPct := "{{% /" + name + " %}}"

		// Find the end of the opener
		openerEnd := strings.Index(text[loc[0]:], ">}}")
		if openerEnd == -1 {
			openerEnd = strings.Index(text[loc[0]:], "%}}")
		}
		if openerEnd == -1 {
			break
		}
		openerEnd += loc[0] + 3 // past the >}} or %}}

		// Find the closer
		rest := text[openerEnd:]
		closerIdx := strings.Index(rest, closerAngle)
		closerLen := len(closerAngle)
		if closerIdx == -1 {
			closerIdx = strings.Index(rest, closerPct)
			closerLen = len(closerPct)
		}
		if closerIdx == -1 {
			// No closer found — just strip the opener
			text = text[:loc[0]] + text[openerEnd:]
			continue
		}

		// Remove opener + content + closer
		endPos := openerEnd + closerIdx + closerLen
		text = text[:loc[0]] + text[endPos:]
	}
	return text
}

// cleanHookBrief strips thinking blocks, channel tokens, and reasoning preamble.
// Matches Python's strip_thinking function exactly.
func cleanHookBrief(text string) string {
	// Remove <think>/<thinking>/<thought>/[THINK] blocks
	text = thinkBlockRe.ReplaceAllString(text, "")

	// Remove channel tokens
	if idx := strings.LastIndex(text, "<channel|>"); idx != -1 {
		text = text[idx+len("<channel|>"):]
	}
	text = channelPairRe.ReplaceAllString(text, "")
	text = channelTokenRe.ReplaceAllString(text, "")

	// Strip verbose reasoning text before THE CONFLICT (matching Python's reasoning_indicators)
	upper := strings.ToUpper(text)
	if i := strings.Index(upper, "THE CONFLICT"); i > 0 {
		before := strings.ToLower(strings.TrimSpace(text[:i]))
		for _, indicator := range reasoningIndicators {
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

// isConnectionError checks if an error is a connection-level error that
// shouldn't be retried (connection refused, unreachable, etc.).
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "no such host") ||
		strings.Contains(msg, "network is unreachable") ||
		strings.Contains(msg, "dial tcp")
}

// EstimateTokensBatch returns approximate total token count for a batch of text strings.
func EstimateTokensBatch(texts []string) int {
	total := 0
	for _, t := range texts {
		total += int(math.Ceil(float64(len(t)) / 4.0))
	}
	return total
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

func hookCompletionCannotRetry(err error) bool {
	var unavailable *client.OpenRouterUnavailableError
	return errors.As(err, &unavailable) || isConnectionError(err)
}
