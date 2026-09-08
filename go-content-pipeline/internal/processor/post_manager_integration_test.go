package processor

import (
	"os"
	"testing"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

// TestDiscoverRealPosts tests post discovery on the actual content directory
// This is an integration test that only runs if the content directory exists
func TestDiscoverRealPosts(t *testing.T) {
	// This test only runs if we're in the rippreport directory structure
	contentDir := "/Volumes/1tb/rippreport/content/p"

	// Check if directory exists
	if _, err := os.Stat(contentDir); os.IsNotExist(err) {
		t.Skip("Skipping integration test: content directory not found")
		return
	}

	pm := NewPostManager(contentDir)

	t.Run("discover all non-holiday posts", func(t *testing.T) {
		filter := models.PostFilter{}
		paths, err := pm.DiscoverPosts(filter)
		if err != nil {
			t.Fatalf("DiscoverPosts failed: %v", err)
		}

		if len(paths) == 0 {
			t.Error("Expected to find some posts")
		}

		t.Logf("Found %d posts", len(paths))

		// Verify all returned paths are valid
		for i, path := range paths {
			if i < 5 { // Log first 5 for inspection
				t.Logf("Post %d: %s", i+1, path)
			}

			// Verify file exists
			if _, err := os.Stat(string(path)); os.IsNotExist(err) {
				t.Errorf("Post file does not exist: %s", path)
			}
		}
	})

	t.Run("exclude backstory category", func(t *testing.T) {
		filter := models.PostFilter{
			ExcludeCategories: []string{"backstory"},
		}
		paths, err := pm.DiscoverPosts(filter)
		if err != nil {
			t.Fatalf("DiscoverPosts failed: %v", err)
		}

		t.Logf("Found %d non-backstory posts", len(paths))

		// Verify no backstory posts in results
		for _, path := range paths {
			post, err := pm.readAndParsePost(path, extractSlugFromPath(string(path)))
			if err != nil {
				continue // Skip posts that can't be parsed
			}

			if post.HasCategory("backstory") {
				t.Errorf("Found backstory post that should be excluded: %s", path)
			}
		}
	})

	t.Run("exclude title patterns", func(t *testing.T) {
		filter := models.PostFilter{
			ExcludeTitles: []string{"wonderful wednesday", "backstory podcast"},
		}
		paths, err := pm.DiscoverPosts(filter)
		if err != nil {
			t.Fatalf("DiscoverPosts failed: %v", err)
		}

		t.Logf("Found %d posts after title exclusions", len(paths))
	})

	t.Run("specific slugs", func(t *testing.T) {
		// Test with a few known slugs (these should exist in most setups)
		filter := models.PostFilter{
			SpecificSlugs: []string{"9-11-24", "9-11-25"},
		}
		paths, err := pm.DiscoverPosts(filter)
		if err != nil {
			t.Fatalf("DiscoverPosts failed: %v", err)
		}

		// Should find at least some of the requested slugs
		t.Logf("Found %d of %d requested posts", len(paths), len(filter.SpecificSlugs))

		// Verify returned posts match requested slugs
		for _, path := range paths {
			slug := extractSlugFromPath(string(path))
			found := false
			for _, requestedSlug := range filter.SpecificSlugs {
				if slug == requestedSlug {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Found unexpected slug: %s", slug)
			}
		}
	})
}
