package parser

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
	"gopkg.in/yaml.v3"
)

// ParseFrontMatter parses Hugo front matter and body from markdown content.
// It detects the delimiter type, extracts front matter and body as separate strings,
// parses the YAML front matter, handles string or array categories, and preserves extra fields.
//
// Requirements: 1.1, 1.2, 1.6, 20.1, 20.2, 20.3, 20.4, 20.5, 20.6, 20.7, 20.8
func ParseFrontMatter(content string) (*models.FrontMatter, string, models.Delimiter, error) {
	delimiter, err := DetectDelimiter(content)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to detect delimiter: %w", err)
	}

	if delimiter != models.DelimiterYAML {
		return nil, "", delimiter, fmt.Errorf("only YAML front matter (---) is currently supported, got TOML (+++)")
	}

	// Split content into lines preserving newlines
	lines := strings.Split(content, "\n")
	if len(lines) < 3 {
		return nil, "", 0, fmt.Errorf("content too short to contain front matter")
	}

	// Find opening delimiter
	startIdx := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == yamlDelimiter {
			startIdx = i + 1
			break
		}
	}

	if startIdx == -1 {
		return nil, "", 0, fmt.Errorf("missing opening front matter delimiter")
	}

	// Find closing delimiter
	endIdx := -1
	for i := startIdx; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == yamlDelimiter {
			endIdx = i
			break
		}
	}

	if endIdx == -1 {
		return nil, "", 0, fmt.Errorf("missing closing front matter delimiter")
	}

	frontMatterLines := lines[startIdx:endIdx]
	frontMatterStr := strings.Join(frontMatterLines, "\n")

	// Body is everything after closing delimiter line
	var body string
	if endIdx+1 < len(lines) {
		body = strings.Join(lines[endIdx+1:], "\n")
	}

	// Parse YAML front matter using yaml.Node to handle polymorphic types (e.g. categories as string vs slice)
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(frontMatterStr), &root); err != nil {
		return nil, "", 0, fmt.Errorf("failed to unmarshal YAML front matter: %w", err)
	}

	// Preprocess node tree if needed (normalize categories to sequence, parse custom dates)
	if len(root.Content) > 0 && root.Content[0].Kind == yaml.MappingNode {
		normalizeMappingNode(root.Content[0])
	}

	var fm models.FrontMatter
	if err := root.Decode(&fm); err != nil {
		// Fallback: try direct unmarshal or return error
		return nil, "", 0, fmt.Errorf("failed to decode front matter into struct: %w", err)
	}

	// Fallback check for categories in Extra if not captured
	if len(fm.Categories) == 0 && fm.Extra != nil {
		if catVal, ok := fm.Extra["categories"]; ok {
			switch c := catVal.(type) {
			case string:
				if strings.TrimSpace(c) != "" {
					fm.Categories = []string{c}
				}
				delete(fm.Extra, "categories")
			case []interface{}:
				for _, item := range c {
					if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
						fm.Categories = append(fm.Categories, s)
					}
				}
				delete(fm.Extra, "categories")
			}
		}
	}

	return &fm, body, delimiter, nil
}

// normalizeMappingNode adjusts scalar categories into a sequence node so yaml.Decode succeeds
func normalizeMappingNode(mapNode *yaml.Node) {
	for i := 0; i < len(mapNode.Content)-1; i += 2 {
		keyNode := mapNode.Content[i]
		valNode := mapNode.Content[i+1]

		if keyNode.Value == "categories" {
			if valNode.Kind == yaml.ScalarNode && valNode.Value != "" {
				// Convert single scalar to sequence
				mapNode.Content[i+1] = &yaml.Node{
					Kind:  yaml.SequenceNode,
					Tag:   "!!seq",
					Style: 0,
					Content: []*yaml.Node{
						{
							Kind:  yaml.ScalarNode,
							Tag:   "!!str",
							Value: valNode.Value,
						},
					},
				}
			}
		} else if keyNode.Value == "date" {
			// If date is scalar with date-only format YYYY-MM-DD, try to parse it
			if valNode.Kind == yaml.ScalarNode {
				if t, err := time.Parse("2006-01-02", valNode.Value); err == nil {
					valNode.Tag = "!!timestamp"
					valNode.Value = t.Format(time.RFC3339)
				}
			}
		}
	}
}

// FormatYAML marshals an arbitrary data structure to a YAML string with 2-space indentation.
//
// Requirements: 21.1, 21.2, 21.3, 21.4, 21.5, 21.6, 21.8
func FormatYAML(data interface{}) (string, error) {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(data); err != nil {
		return "", fmt.Errorf("failed to format YAML: %w", err)
	}

	if err := encoder.Close(); err != nil {
		return "", fmt.Errorf("failed to close YAML encoder: %w", err)
	}

	return buf.String(), nil
}

// SerializeFrontMatter serializes FrontMatter into formatted markdown front matter with delimiters.
//
// Requirements: 1.2, 1.4, 20.7, 21.1, 21.8
func SerializeFrontMatter(fm *models.FrontMatter, delimiter models.Delimiter) (string, error) {
	if fm == nil {
		return "", fmt.Errorf("front matter cannot be nil")
	}

	delimStr := yamlDelimiter
	if delimiter == models.DelimiterTOML {
		delimStr = tomlDelimiter
	}

	yamlStr, err := FormatYAML(fm)
	if err != nil {
		return "", fmt.Errorf("failed to serialize front matter: %w", err)
	}

	// Ensure yamlStr is trimmed of trailing newline for consistent wrapping
	yamlStr = strings.TrimSpace(yamlStr)

	var sb strings.Builder
	sb.WriteString(delimStr)
	sb.WriteString("\n")
	sb.WriteString(yamlStr)
	sb.WriteString("\n")
	sb.WriteString(delimStr)
	sb.WriteString("\n")

	return sb.String(), nil
}
