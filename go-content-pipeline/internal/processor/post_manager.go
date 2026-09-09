package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
	"github.com/rippreport/go-content-pipeline/internal/parser"
)

// PostManager handles Hugo post discovery and management
type PostManager struct {
	postsDirectory string
}

// NewPostManager creates a new PostManager instance
func NewPostManager(postsDirectory string) *PostManager {
	return &PostManager{
		postsDirectory: postsDirectory,
	}
}

// DiscoverPosts walks the content/p/ directory and returns all posts matching the filter criteria.
// Posts are sorted by date descending (newest first).
func (pm *PostManager) DiscoverPosts(filter models.PostFilter) ([]models.PostPath, error) {
	// Ensure "holiday" is always in the exclude list
	if !containsIgnoreCase(filter.ExcludeCategories, "holiday") {
		filter.ExcludeCategories = append(filter.ExcludeCategories, "holiday")
	}

	var discoveredPaths []models.PostPath
	var postsWithDates []postWithDate

	// Walk the posts directory
	err := filepath.Walk(pm.postsDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process index.md files
		if info.Name() != "index.md" {
			return nil
		}

		postPath := models.PostPath(path)

		// Extract slug from path (parent directory name)
		slug := extractSlugFromPath(path)
		if slug == "" {
			return nil // Skip if we can't extract a slug
		}

		// If specific slugs are requested, check if this is one of them
		if len(filter.SpecificSlugs) > 0 {
			if !containsSlug(filter.SpecificSlugs, slug) {
				return nil // Skip posts not in the specific list
			}
		}

		// Read and parse the post to apply filters
		post, err := pm.readAndParsePost(postPath, slug)
		if err != nil {
			// Log the error but continue processing other posts
			// In production, this should use the logger
			return nil
		}

		// Apply filter criteria
		if !filter.Matches(post) {
			return nil
		}

		// Add to list with date for sorting
		postsWithDates = append(postsWithDates, postWithDate{
			path: postPath,
			date: post.FrontMatter.Date,
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk posts directory: %w", err)
	}

	// Sort by date descending (newest first)
	sort.Slice(postsWithDates, func(i, j int) bool {
		return postsWithDates[i].date.After(postsWithDates[j].date)
	})

	// Extract just the paths
	for _, pwd := range postsWithDates {
		discoveredPaths = append(discoveredPaths, pwd.path)
	}

	return discoveredPaths, nil
}

// postWithDate is a helper struct for sorting posts by date
type postWithDate struct {
	path models.PostPath
	date time.Time
}

// readAndParsePost reads a post file and parses its front matter for filtering purposes
func (pm *PostManager) readAndParsePost(path models.PostPath, slug string) (*models.Post, error) {
	content, err := os.ReadFile(string(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read post file: %w", err)
	}

	// Parse front matter
	// This is a minimal parse just for filtering - full parsing would use the parser package
	post := &models.Post{
		Path:        path,
		Slug:        slug,
		FrontMatter: &models.FrontMatter{},
	}

	// For now, we need to parse the front matter
	// We'll use the parser package once it's available for full parsing
	// For discovery, we need at least: categories, title, date
	err = pm.quickParseFrontMatter(string(content), post.FrontMatter)
	if err != nil {
		return nil, fmt.Errorf("failed to parse front matter: %w", err)
	}

	return post, nil
}

// quickParseFrontMatter does a minimal parse of front matter for filtering purposes
// This is temporary until the full parser is integrated
func (pm *PostManager) quickParseFrontMatter(content string, fm *models.FrontMatter) error {
	// This is a simplified parser for discovery purposes
	// It extracts the fields needed for filtering: title, date, categories

	lines := strings.Split(content, "\n")
	if len(lines) < 3 {
		return fmt.Errorf("content too short to contain front matter")
	}

	// Check for YAML delimiter
	if !strings.HasPrefix(strings.TrimSpace(lines[0]), "---") {
		return fmt.Errorf("missing YAML front matter delimiter")
	}

	// Find the closing delimiter
	endIdx := -1
	for i := 1; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "---") {
			endIdx = i
			break
		}
	}

	if endIdx == -1 {
		return fmt.Errorf("missing closing YAML delimiter")
	}

	// Parse key-value pairs
	for i := 1; i < endIdx; i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first colon
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "title":
			// Remove quotes if present
			fm.Title = strings.Trim(value, `"'`)
		case "date":
			// Parse date
			dateStr := strings.Trim(value, `"'`)
			parsedDate, err := time.Parse(time.RFC3339, dateStr)
			if err != nil {
				// Try date-only format
				parsedDate, err = time.Parse("2006-01-02", dateStr)
				if err != nil {
					return fmt.Errorf("failed to parse date: %w", err)
				}
			}
			fm.Date = parsedDate
		case "categories":
			// Handle categories array or string
			if strings.HasPrefix(value, "[") {
				// Array format: [cat1, cat2]
				value = strings.Trim(value, "[]")
				cats := strings.Split(value, ",")
				for _, cat := range cats {
					cat = strings.TrimSpace(cat)
					cat = strings.Trim(cat, `"'`)
					if cat != "" {
						fm.Categories = append(fm.Categories, cat)
					}
				}
			} else if value != "" {
				// String format or list item
				value = strings.Trim(value, `"'`)
				if value != "" {
					fm.Categories = append(fm.Categories, value)
				}
			}
		}
	}

	// Check for list-style categories (- item)
	inCategories := false
	for i := 1; i < endIdx; i++ {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "categories:") {
			inCategories = true
			continue
		}
		if inCategories {
			if strings.HasPrefix(line, "- ") {
				cat := strings.TrimSpace(strings.TrimPrefix(line, "- "))
				cat = strings.Trim(cat, `"'`)
				if cat != "" && !containsIgnoreCase(fm.Categories, cat) {
					fm.Categories = append(fm.Categories, cat)
				}
			} else if !strings.HasPrefix(line, " ") && line != "" {
				// End of categories section
				inCategories = false
			}
		}
	}

	return nil
}

// extractSlugFromPath extracts the slug (parent directory name) from a post path
func extractSlugFromPath(path string) string {
	// Get the directory containing index.md
	dir := filepath.Dir(path)
	// Get the base name of that directory
	return filepath.Base(dir)
}

// containsIgnoreCase checks if a slice contains a string (case-insensitive)
func containsIgnoreCase(slice []string, item string) bool {
	lowerItem := strings.ToLower(item)
	for _, s := range slice {
		if strings.ToLower(s) == lowerItem {
			return true
		}
	}
	return false
}

// containsSlug checks if a slice contains a specific slug
func containsSlug(slice []string, slug string) bool {
	for _, s := range slice {
		if s == slug {
			return true
		}
	}
	return false
}

// GetEmbeddingContent returns the appropriate content to use for embedding generation.
// If transcript.md exists in the post directory, it uses that content.
// Otherwise, it falls back to the post body.
func (pm *PostManager) GetEmbeddingContent(post *models.Post) (string, error) {
	// Try to read transcript.md from the post directory
	postDir := filepath.Dir(string(post.Path))
	transcriptPath := filepath.Join(postDir, "transcript.md")

	// Check if transcript.md exists
	if _, err := os.Stat(transcriptPath); err == nil {
		// Transcript exists, read and return it
		transcriptContent, err := os.ReadFile(transcriptPath)
		if err != nil {
			return "", fmt.Errorf("error reading transcript file: %w", err)
		}
		return string(transcriptContent), nil
	}

	// Transcript doesn't exist, return the post body
	return post.Body, nil
}

// ReadPost reads a Hugo post from the given path, parses its front matter,
// extracts the markdown body, calculates the content hash, and validates required fields.
//
// Requirements: 1.1, 1.6, 1.7
func (pm *PostManager) ReadPost(path models.PostPath) (*models.Post, error) {
	// Read file content
	content, err := os.ReadFile(string(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read post file %s: %w", path, err)
	}

	// Parse front matter using the parser package
	frontMatter, body, delimiter, err := parser.ParseFrontMatter(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse front matter: %w", err)
	}

	// Extract slug from path
	slug := extractSlugFromPath(string(path))
	if slug == "" {
		return nil, fmt.Errorf("failed to extract slug from path %s", path)
	}

	// Create post object
	post := &models.Post{
		Path:        path,
		Slug:        slug,
		FrontMatter: frontMatter,
		Body:        body,
	}

	// Calculate content hash
	post.UpdateContentHash()

	// Validate required fields (title, date)
	if err := post.FrontMatter.Validate(); err != nil {
		return nil, fmt.Errorf("post validation failed: %w", err)
	}

	_ = delimiter

	return post, nil
}

// UpdateRelatedArticles updates the related_articles field in a post's front matter,
// preserving all other front matter fields, formatting, and markdown body,
// using an atomic write pattern.
//
// Requirements: 1.2, 1.4
func (pm *PostManager) UpdateRelatedArticles(path models.PostPath, articles []models.RelatedArticle) error {
	filePath := string(path)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read post file for update %s: %w", path, err)
	}

	// Parse front matter
	frontMatter, body, delimiter, err := parser.ParseFrontMatter(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse post front matter for update: %w", err)
	}

	// Validate all articles
	for i, article := range articles {
		if err := article.Validate(); err != nil {
			return fmt.Errorf("invalid related article %d: %w", i, err)
		}
	}

	// Update related articles
	frontMatter.RelatedArticles = articles

	// Serialize front matter back with delimiters
	serializedFm, err := parser.SerializeFrontMatter(frontMatter, delimiter)
	if err != nil {
		return fmt.Errorf("failed to serialize updated front matter: %w", err)
	}

	// Construct updated document preserving body byte-for-byte
	updatedContent := serializedFm + body

	// Atomic write using a temporary file in the same directory
	dir := filepath.Dir(filePath)
	tempFile, err := os.CreateTemp(dir, ".index.md.tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file for atomic write: %w", err)
	}
	tempPath := tempFile.Name()

	// Ensure cleanup if rename does not succeed
	success := false
	defer func() {
		if !success {
			_ = os.Remove(tempPath)
		}
	}()

	if _, err := tempFile.WriteString(updatedContent); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	if err := tempFile.Sync(); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if err := os.Rename(tempPath, filePath); err != nil {
		return fmt.Errorf("failed to atomically replace post file %s: %w", path, err)
	}

	success = true
	return nil
}

