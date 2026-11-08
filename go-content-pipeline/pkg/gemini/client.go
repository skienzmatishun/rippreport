package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultModel   = "gemini-3.7-flash"
	defaultBaseURL = "https://generativelanguage.googleapis.com/v1beta"
	defaultTimeout = 120 * time.Second
)

// Client represents a Gemini API client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Config holds configuration for the Gemini client
type Config struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
}

// HookItem represents a hook to be refined
type HookItem struct {
	Slug          string   `json:"slug"`
	Brief         string   `json:"brief"`
	Text          string   `json:"text,omitempty"` // Current hook text, supplied separately for collision rewrites.
	AvoidExamples []string `json:"avoid_examples,omitempty"`
}

// RefinedHookItem represents a refined hook
type RefinedHookItem struct {
	Slug string `json:"slug"`
	Text string `json:"text"`
}

// NewClient creates a new Gemini API client
func NewClient(cfg Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = defaultTimeout
	}

	return &Client{
		apiKey:  cfg.APIKey,
		baseURL: cfg.BaseURL,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}, nil
}

// RefineHooks refines a list of hooks using the Gemini API
func (c *Client) RefineHooks(ctx context.Context, items []HookItem, model string) ([]RefinedHookItem, error) {
	return c.refineHooks(ctx, items, model, "")
}

// RewriteCollision rewrites only the second hook, keeping the first as context.
func (c *Client) RewriteCollision(ctx context.Context, first, second HookItem, model string) ([]RefinedHookItem, error) {
	return c.refineHooks(ctx, []HookItem{first, second}, model, "These two hooks are too similar. Rewrite the second one by focusing on a different specific fact from its source brief. Write at a sixth-grade reading level. Use short sentences and common words. No metaphors, similes, analogies, idioms, or figurative comparisons. Each item contains its current hook in text and the complete source brief in brief. Rewrite the second item's text using its full brief as factual context. Preserve its facts and specific names; do not invent claims. Return ONLY the second hook, with its original slug. The first hook is context only. Any avoid_examples are rejected previous versions: do not repeat their wording, structure, or factual focus. Do not replace an old metaphor with a new one. Use a different fact supported by the full brief. Remove references to Article 1, Article 2, the first article, or the second article; describe the actual people, events, and facts instead.")
}

func (c *Client) refineHooks(ctx context.Context, items []HookItem, model, instruction string) ([]RefinedHookItem, error) {
	if model == "" {
		model = DefaultModel
	}

	systemInstruction := `You are the Senior Editorial Writer for The Ripp Report, an independent watchdog publication.

Your task: Transform raw LLM-generated sidebar hooks into polished, engaging contextual bridges that match The Ripp Report's distinctive editorial voice.

## EDITORIAL VOICE

Your writing must be:
- **Direct and conversational**: Write as though speaking to an informed reader
- **Skeptical and confrontational**: Challenge official narratives and question convenient coincidences
- **Factual**: State what happened in plain, literal language. Preserve the source's distinction between allegations and established facts.
- **Investigative**: Connect dots and identify patterns
- **Populist**: Emphasize public interest over institutional protection

## STYLE REQUIREMENTS

Favor:
- Punchy, direct sentences
- Natural transitions
- Strong opening phrases
- Specific names and details (use them!)
- Rhetorical questions that highlight contradictions
- Clear statements of what people did and why it matters

Avoid:
- Academic vocabulary
- Corporate jargon
- Generic AI-sounding language
- Repetitive sentence structures
- Constant outrage (vary your tone)
- Starting every hook the same way

## TECHNICAL CONSTRAINTS

- Length: 220-320 characters total
- Two sentences preferred (can be one if punchy enough)
- Sixth-grade reading level: use common words, short sentences, and simple sentence structure
- NO metaphors, similes, analogies, idioms, or figurative comparisons
- Do not compare people or events to treasure hunters, pirates, games, machines, animals, or other unrelated subjects
- Remove figurative language already present in the input; state the underlying facts directly
- Do not add people, historical references, or claims absent from the source
- Do not refer to Article 1, Article 2, the first article, or the second article; name the actual subject instead
- NO passive voice
- Must use specific entity names from the source text
- OUTPUT ONLY THE HOOK TEXT - NO bullet points, NO entity lists, NO extra formatting

## VARIETY IS CRITICAL

Do NOT start every hook with the same phrase or structure. Mix it up:
- Some start with questions
- Some start with stark observations
- Some start with names
- Some start with locations or institutions

## INPUT/OUTPUT FORMAT

INPUT: A JSON list of objects with 'slug' and 'brief' (the raw LLM-generated text)
OUTPUT: A JSON list of objects with 'slug' and 'text' (your refined version - JUST the hook text, nothing else)

CRITICAL: The "text" field must contain ONLY the hook sentences. Do NOT add bullet lists, entity lists, or any other formatting after the hook.

REMEMBER: You're creating a reason to click, not a summary. Make it compelling.`

	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal items: %w", err)
	}

	userPrompt := fmt.Sprintf("Refine these sidebar hooks. Each 'brief' field contains raw LLM-generated text. Transform them into compelling, varied editorial hooks that match The Ripp Report's voice. Write at a sixth-grade reading level using short sentences and common words. No metaphors, similes, analogies, idioms, or figurative comparisons. State only source-supported facts. Never repeat the same opening structure:\n\n%s", string(itemsJSON))

	if instruction != "" {
		userPrompt = instruction + "\n\n" + string(itemsJSON)
	}

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": userPrompt},
				},
			},
		},
		"systemInstruction": map[string]interface{}{
			"parts": []map[string]string{
				{"text": systemInstruction},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":        1.0,
			"responseMimeType":   "application/json",
			"responseModalities": []string{"TEXT"},
		},
	}

	reqBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", c.baseURL, model, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(apiResp.Candidates) == 0 || len(apiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no text returned in response")
	}

	responseText := apiResp.Candidates[0].Content.Parts[0].Text

	// Parse the JSON response
	var refined []RefinedHookItem
	if err := json.Unmarshal([]byte(responseText), &refined); err != nil {
		// Try to extract JSON from the response if it has extra text
		refined, err = extractJSON(responseText)
		if err != nil {
			return nil, fmt.Errorf("failed to parse refined hooks: %w", err)
		}
	}

	// Post-process: remove any bullet lists that might have been added
	for i := range refined {
		refined[i].Text = stripBulletLists(refined[i].Text)
	}

	return refined, nil
}

// stripBulletLists removes bullet point lists from hook text
func stripBulletLists(text string) string {
	lines := strings.Split(text, "\n")
	var cleaned []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip lines that are bullet points (start with -, •, *, or numbered lists)
		if strings.HasPrefix(trimmed, "- ") ||
			strings.HasPrefix(trimmed, "* ") ||
			strings.HasPrefix(trimmed, "• ") ||
			(len(trimmed) > 2 && trimmed[0] >= '0' && trimmed[0] <= '9' && trimmed[1] == '.') {
			continue
		}
		// Keep non-bullet lines
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}

	return strings.Join(cleaned, " ")
}

// extractJSON attempts to extract a JSON array from text that may contain extra content
func extractJSON(text string) ([]RefinedHookItem, error) {
	// Find the first [ and match to its closing ]
	start := -1
	for i, ch := range text {
		if ch == '[' {
			start = i
			break
		}
	}

	if start == -1 {
		return nil, fmt.Errorf("no JSON array found in response")
	}

	depth := 0
	for i := start; i < len(text); i++ {
		switch text[i] {
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				jsonStr := text[start : i+1]
				var refined []RefinedHookItem
				if err := json.Unmarshal([]byte(jsonStr), &refined); err != nil {
					return nil, err
				}
				return refined, nil
			}
		}
	}

	return nil, fmt.Errorf("unmatched brackets in JSON")
}
