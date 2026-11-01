package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
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
	Slug  string `json:"slug"`
	Brief string `json:"brief"`
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
	if model == "" {
		model = "gemini-3.7-flash"
	}

	systemInstruction := `You are the Editor-in-Chief of The Ripp Report. 
Your task: Transform raw "Bridge Briefs" into aggressive, 2-sentence "Contextual Bridges."

A bridge is NOT a summary. It is a reason to click. Use the PIVOT section of the brief as your primary inspiration.

Sentence 1: The setup. Connect to the theme of the Main Article.
Sentence 2: The payoff. Use the PIVOT and DISCOVERY to show why the Widget Article is the "smoking gun."

Style: Cynical, hard-hitting, investigative. Use specific names (Wilters, Lowery, etc.).
Constraints: 220-320 characters. NO passive voice. 6th-grade vocabulary.
Variety: Do not start every hook with the same phrase.

INPUT: A JSON list of objects with 'slug' and 'brief'.
OUTPUT: A JSON list of objects with 'slug' and 'text'.`

	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal items: %w", err)
	}

	userPrompt := fmt.Sprintf("Refine these hooks:\n%s", string(itemsJSON))

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

	return refined, nil
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
