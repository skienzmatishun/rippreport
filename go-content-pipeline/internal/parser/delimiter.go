package parser

import (
	"fmt"
	"strings"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

const (
	yamlDelimiter = "---"
	tomlDelimiter = "+++"
)

// DetectDelimiter detects the front matter delimiter type in Hugo markdown content.
// It checks the first non-empty line of the content for YAML (---) or TOML (+++) delimiters.
// Returns an error if no valid delimiter is found or if the content is empty.
//
// Requirements: 20.1, 20.8
func DetectDelimiter(content string) (models.Delimiter, error) {
	if content == "" {
		return 0, fmt.Errorf("content is empty")
	}

	// Trim leading whitespace and get the first line
	trimmed := strings.TrimLeft(content, " \t\r\n")
	if trimmed == "" {
		return 0, fmt.Errorf("content contains only whitespace")
	}

	// Find the first line
	firstLineEnd := strings.IndexAny(trimmed, "\r\n")
	var firstLine string
	if firstLineEnd == -1 {
		// Content is a single line
		firstLine = trimmed
	} else {
		firstLine = trimmed[:firstLineEnd]
	}

	// Trim any trailing whitespace from the first line
	firstLine = strings.TrimRight(firstLine, " \t")

	// Check for YAML delimiter
	if firstLine == yamlDelimiter {
		return models.DelimiterYAML, nil
	}

	// Check for TOML delimiter
	if firstLine == tomlDelimiter {
		return models.DelimiterTOML, nil
	}

	// No valid delimiter found
	return 0, fmt.Errorf("invalid or missing front matter delimiter: expected '---' or '+++', got '%s'", firstLine)
}
