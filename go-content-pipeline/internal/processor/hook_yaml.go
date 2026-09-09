package processor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
	"gopkg.in/yaml.v3"
)

type hookEntryYAML struct {
	Brief       string `yaml:"brief"`
	Text        string `yaml:"text,omitempty"`
	WidgetDate  string `yaml:"widget_date,omitempty"`
	GeneratedAt string `yaml:"generated_at,omitempty"`
	ModelName   string `yaml:"model_name,omitempty"`
}

type hookMetaYAML struct {
	GeneratedAt      string `yaml:"generated_at"`
	GeneratorVersion string `yaml:"generator_version"`
	PostSlug         string `yaml:"post_slug"`
	MainDate         string `yaml:"main_date,omitempty"`
	RefinedBy        string `yaml:"refined_by,omitempty"`
	RefinedAt        string `yaml:"refined_at,omitempty"`
	ModelName        string `yaml:"model_name,omitempty"`
}

func formatTimestamp(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("2006-01-02T15:04:05.000000Z")
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func parseTimestamp(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	if t, err := time.Parse("2006-01-02T15:04:05.000000Z", s); err == nil {
		return t
	}
	return time.Time{}
}

func metadataToYAML(meta models.HookMetadata) hookMetaYAML {
	out := hookMetaYAML{
		GeneratedAt:      formatTimestamp(meta.GeneratedAt),
		GeneratorVersion: meta.GeneratorVersion,
		PostSlug:         meta.PostSlug,
		MainDate:         meta.MainDate,
		RefinedBy:        meta.RefinedBy,
		ModelName:        meta.ModelName,
	}
	if !meta.RefinedAt.IsZero() {
		out.RefinedAt = formatTimestamp(meta.RefinedAt)
	}
	return out
}

func metadataFromYAML(in hookMetaYAML) models.HookMetadata {
	return models.HookMetadata{
		GeneratedAt:      parseTimestamp(in.GeneratedAt),
		GeneratorVersion: in.GeneratorVersion,
		PostSlug:         in.PostSlug,
		MainDate:         in.MainDate,
		RefinedBy:        in.RefinedBy,
		RefinedAt:        parseTimestamp(in.RefinedAt),
		ModelName:        in.ModelName,
	}
}

func encodeNode(v interface{}) (*yaml.Node, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	_ = enc.Close()

	var doc yaml.Node
	if err := yaml.Unmarshal(buf.Bytes(), &doc); err != nil {
		return nil, err
	}
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 {
		return doc.Content[0], nil
	}
	return &doc, nil
}

func scalarKey(v string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: v}
}

func formatHookYAML(meta models.HookMetadata, hooks []*models.Hook) (string, error) {
	if meta.GeneratedAt.IsZero() {
		meta.GeneratedAt = time.Now().UTC()
	}
	if meta.GeneratorVersion == "" {
		meta.GeneratorVersion = models.HookGeneratorVersion
	}

	if meta.ModelName == "" && len(hooks) > 0 && hooks[0].ModelName != "" {
		meta.ModelName = hooks[0].ModelName
	}

	metaNode, err := encodeNode(metadataToYAML(meta))
	if err != nil {
		return "", fmt.Errorf("failed to encode hook metadata: %w", err)
	}

	hooksNode := &yaml.Node{Kind: yaml.MappingNode}
	for _, h := range hooks {
		entry := hookEntryYAML{
			Brief:       h.Brief,
			Text:        h.Text,
			WidgetDate:  h.WidgetDate,
			GeneratedAt: formatTimestamp(h.Generated),
			ModelName:   h.ModelName,
		}
		val, err := encodeNode(entry)
		if err != nil {
			return "", fmt.Errorf("failed to encode hook %s: %w", h.WidgetSlug, err)
		}
		hooksNode.Content = append(hooksNode.Content, scalarKey(h.WidgetSlug), val)
	}

	root := &yaml.Node{Kind: yaml.MappingNode}
	root.Content = append(root.Content,
		scalarKey("metadata"), metaNode,
		scalarKey("hooks"), hooksNode,
	)

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(root); err != nil {
		return "", fmt.Errorf("failed to encode hook YAML: %w", err)
	}
	if err := enc.Close(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func parseHookYAML(data []byte) (models.HookMetadata, []*models.Hook, error) {
	var raw map[string]yaml.Node
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return models.HookMetadata{}, nil, err
	}

	var meta models.HookMetadata
	if node, ok := raw["metadata"]; ok {
		var my hookMetaYAML
		if err := node.Decode(&my); err != nil {
			return models.HookMetadata{}, nil, fmt.Errorf("failed to parse hook metadata: %w", err)
		}
		meta = metadataFromYAML(my)
	} else {
		// Legacy Go stub: post_slug / generated at top level
		if n, ok := raw["post_slug"]; ok {
			meta.PostSlug = strings.TrimSpace(n.Value)
		}
		if n, ok := raw["generated"]; ok {
			meta.GeneratedAt = parseTimestamp(n.Value)
		}
		meta.GeneratorVersion = models.HookGeneratorVersion
	}

	hooksNode, ok := raw["hooks"]
	if !ok {
		return models.HookMetadata{}, nil, fmt.Errorf("hooks section missing")
	}

	var hooks []*models.Hook
	switch hooksNode.Kind {
	case yaml.MappingNode:
		for i := 0; i+1 < len(hooksNode.Content); i += 2 {
			slug := hooksNode.Content[i].Value
			var entry hookEntryYAML
			if err := hooksNode.Content[i+1].Decode(&entry); err != nil {
				return models.HookMetadata{}, nil, fmt.Errorf("failed to parse hook %s: %w", slug, err)
			}
			modelName := entry.ModelName
			if modelName == "" {
				modelName = meta.ModelName
			}
			hooks = append(hooks, &models.Hook{
				WidgetSlug: slug,
				Brief:      entry.Brief,
				Text:       entry.Text,
				WidgetDate: entry.WidgetDate,
				Generated:  parseTimestamp(entry.GeneratedAt),
				ModelName:  modelName,
			})
		}
	case yaml.SequenceNode:
		for _, item := range hooksNode.Content {
			var legacy struct {
				WidgetSlug  string `yaml:"widget_slug"`
				WidgetTitle string `yaml:"widget_title"`
				Brief       string `yaml:"brief"`
				Generated   string `yaml:"generated"`
				ModelName   string `yaml:"model_name"`
			}
			if err := item.Decode(&legacy); err != nil {
				return models.HookMetadata{}, nil, fmt.Errorf("failed to parse legacy hook: %w", err)
			}
			hooks = append(hooks, &models.Hook{
				WidgetSlug:  legacy.WidgetSlug,
				WidgetTitle: legacy.WidgetTitle,
				Brief:       legacy.Brief,
				Generated:   parseTimestamp(legacy.Generated),
				ModelName:   legacy.ModelName,
			})
		}
	default:
		return models.HookMetadata{}, nil, fmt.Errorf("hooks must be a map keyed by widget slug")
	}

	if meta.GeneratedAt.IsZero() {
		meta.GeneratedAt = time.Now().UTC()
	}
	for _, h := range hooks {
		if h.Generated.IsZero() {
			h.Generated = meta.GeneratedAt
		}
		if h.ModelName == "" {
			h.ModelName = meta.ModelName
		}
	}

	return meta, hooks, nil
}

func writeHookFile(path string, meta models.HookMetadata, hooks []*models.Hook) error {
	if len(hooks) == 0 {
		return fmt.Errorf("hooks cannot be empty")
	}
	for i, h := range hooks {
		if err := h.Validate(); err != nil {
			return fmt.Errorf("hook %d: %w", i, err)
		}
	}

	yamlStr, err := formatHookYAML(meta, hooks)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	tempFile, err := os.CreateTemp(dir, ".hooks_*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp hooks file: %w", err)
	}
	tempPath := tempFile.Name()

	success := false
	defer func() {
		if !success {
			_ = os.Remove(tempPath)
		}
	}()

	if _, err := tempFile.WriteString(yamlStr); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to write hook YAML: %w", err)
	}
	if err := tempFile.Sync(); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to sync hook file: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp hook file: %w", err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("failed to atomically replace %s: %w", path, err)
	}
	success = true
	return nil
}
