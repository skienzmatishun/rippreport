package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

// TestDetectDelimiter_RealHugoPosts tests delimiter detection on actual Hugo posts from the content directory
func TestDetectDelimiter_RealHugoPosts(t *testing.T) {
	contentDir := "/Volumes/1tb/rippreport/content/p"

	// Check if the content directory exists
	if _, err := os.Stat(contentDir); os.IsNotExist(err) {
		t.Skipf("Content directory not found: %s", contentDir)
		return
	}

	// Test a few known posts
	testPosts := []string{
		"justice-for/index.md",
		"bonum-est-faciet-2022/index.md",
		"state-of-the-city-2019/index.md",
	}

	for _, postPath := range testPosts {
		t.Run(postPath, func(t *testing.T) {
			fullPath := filepath.Join(contentDir, postPath)
			
			// Read the file
			content, err := os.ReadFile(fullPath)
			if err != nil {
				if os.IsNotExist(err) {
					t.Skipf("Post not found: %s", fullPath)
					return
				}
				t.Fatalf("Failed to read post: %v", err)
			}

			// Detect delimiter
			delimiter, err := DetectDelimiter(string(content))
			if err != nil {
				t.Errorf("Failed to detect delimiter in %s: %v", postPath, err)
				return
			}

			// All known posts should use YAML (---)
			if delimiter != models.DelimiterYAML {
				t.Errorf("Expected YAML delimiter for %s, got %v", postPath, delimiter)
			}

			t.Logf("✓ Successfully detected YAML delimiter in %s", postPath)
		})
	}
}

// TestDetectDelimiter_ScanAllPosts scans all posts in the content directory
// This test is skipped by default but can be run with -run flag
func TestDetectDelimiter_ScanAllPosts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping comprehensive scan in short mode")
	}

	contentDir := "/Volumes/1tb/rippreport/content/p"

	// Check if the content directory exists
	if _, err := os.Stat(contentDir); os.IsNotExist(err) {
		t.Skipf("Content directory not found: %s", contentDir)
		return
	}

	var totalPosts, successfulPosts, failedPosts int
	var yamlCount, tomlCount int

	// Walk through all post directories
	entries, err := os.ReadDir(contentDir)
	if err != nil {
		t.Fatalf("Failed to read content directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		indexPath := filepath.Join(contentDir, entry.Name(), "index.md")
		
		// Check if index.md exists
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			continue
		}

		totalPosts++

		// Read the file
		content, err := os.ReadFile(indexPath)
		if err != nil {
			t.Logf("Failed to read %s: %v", indexPath, err)
			failedPosts++
			continue
		}

		// Detect delimiter
		delimiter, err := DetectDelimiter(string(content))
		if err != nil {
			t.Logf("Failed to detect delimiter in %s: %v", indexPath, err)
			failedPosts++
			continue
		}

		successfulPosts++
		if delimiter == models.DelimiterYAML {
			yamlCount++
		} else {
			tomlCount++
		}
	}

	// Report statistics
	t.Logf("Scanned %d posts:", totalPosts)
	t.Logf("  ✓ Successful: %d", successfulPosts)
	t.Logf("  ✗ Failed: %d", failedPosts)
	t.Logf("  YAML (---): %d", yamlCount)
	t.Logf("  TOML (+++): %d", tomlCount)

	if failedPosts > 0 {
		t.Errorf("%d posts failed delimiter detection", failedPosts)
	}
}
