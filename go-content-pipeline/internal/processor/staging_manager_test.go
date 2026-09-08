package processor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

func TestStagingManager_SaveAndRead(t *testing.T) {
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	backupDir := filepath.Join(tmpDir, "backups")

	pm := NewPostManager(contentDir)
	bm := NewBackupManager(backupDir)
	sm := NewStagingManager(contentDir, pm, bm)

	articles := []models.RelatedArticle{
		{
			Slug:  "recommended-post-1",
			Title: "Recommended Post 1",
			Score: 92,
			Rank:  1,
		},
		{
			Slug:  "recommended-post-2",
			Title: "Recommended Post 2",
			Score: 84,
			Rank:  2,
		},
	}

	err := sm.SaveStaged("post-alpha", articles)
	if err != nil {
		t.Fatalf("SaveStaged failed: %v", err)
	}

	readArticles, err := sm.ReadStaged("post-alpha")
	if err != nil {
		t.Fatalf("ReadStaged failed: %v", err)
	}

	if len(readArticles) != 2 {
		t.Fatalf("expected 2 articles, got %d", len(readArticles))
	}
	if readArticles[0].Slug != "recommended-post-1" || readArticles[0].Score != 92 {
		t.Errorf("unexpected article 0: %+v", readArticles[0])
	}
}

func TestStagingManager_ApplyAndClear(t *testing.T) {
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	backupDir := filepath.Join(tmpDir, "backups")

	// Create a mock post
	postDir := filepath.Join(contentDir, "post-beta")
	if err := os.MkdirAll(postDir, 0755); err != nil {
		t.Fatalf("failed to create post dir: %v", err)
	}
	initialPost := `---
title: "Beta Post"
date: 2025-01-15
categories:
  - news
---

Initial Body`
	if err := os.WriteFile(filepath.Join(postDir, "index.md"), []byte(initialPost), 0644); err != nil {
		t.Fatalf("failed to write post: %v", err)
	}

	pm := NewPostManager(contentDir)
	bm := NewBackupManager(backupDir)
	sm := NewStagingManager(contentDir, pm, bm)

	stagedArticles := []models.RelatedArticle{
		{
			Slug:  "staged-winner",
			Title: "Staged Winner",
			Score: 99,
			Rank:  1,
		},
	}

	if err := sm.SaveStaged("post-beta", stagedArticles); err != nil {
		t.Fatalf("SaveStaged failed: %v", err)
	}

	// List staged
	stagedList, err := sm.ListStaged()
	if err != nil {
		t.Fatalf("ListStaged failed: %v", err)
	}
	if len(stagedList) != 1 || stagedList[0] != "post-beta" {
		t.Errorf("expected [post-beta], got %v", stagedList)
	}

	// Apply staged
	if err := sm.ApplyStaged("post-beta"); err != nil {
		t.Fatalf("ApplyStaged failed: %v", err)
	}

	// Staging file should be deleted after apply
	stagedPath := filepath.Join(postDir, "related-articles.staged.yaml")
	if _, err := os.Stat(stagedPath); !os.IsNotExist(err) {
		t.Errorf("expected staging file to be deleted after apply")
	}

	// Post should have updated related articles
	post, err := pm.ReadPost(models.PostPath(filepath.Join(postDir, "index.md")))
	if err != nil {
		t.Fatalf("ReadPost failed: %v", err)
	}
	if len(post.FrontMatter.RelatedArticles) != 1 || post.FrontMatter.RelatedArticles[0].Slug != "staged-winner" {
		t.Errorf("expected post to have staged article, got %+v", post.FrontMatter.RelatedArticles)
	}

	// ClearStaged on non-existent shouldn't error
	if err := sm.ClearStaged([]string{"post-beta"}); err != nil {
		t.Errorf("ClearStaged failed: %v", err)
	}
}

func TestStagingManager_ValidationErrors(t *testing.T) {
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	sm := NewStagingManager(contentDir, nil, nil)

	t.Run("empty post slug", func(t *testing.T) {
		err := sm.SaveStaged("", []models.RelatedArticle{
			{Slug: "a", Title: "b", Score: 10, Rank: 1},
		})
		if err == nil {
			t.Error("expected error for empty post slug")
		}
	})

	t.Run("empty articles", func(t *testing.T) {
		err := sm.SaveStaged("slug", nil)
		if err == nil {
			t.Error("expected error for empty articles")
		}
	})

	t.Run("invalid article in slice", func(t *testing.T) {
		err := sm.SaveStaged("slug", []models.RelatedArticle{
			{Slug: "", Title: "no slug", Score: 10, Rank: 1},
		})
		if err == nil {
			t.Error("expected error for invalid article")
		}
	})
}
