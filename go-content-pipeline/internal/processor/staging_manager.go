package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
	"github.com/rippreport/go-content-pipeline/internal/parser"
	"gopkg.in/yaml.v3"
)

const stagedFilename = "related-articles.staged.yaml"

// StagingManager manages staged rankings workflow for posts.
//
// Requirements: 10.1, 10.2, 10.3, 10.4, 10.5, 10.6, 10.7, 10.8
type StagingManager struct {
	mu            sync.Mutex
	contentDir    string
	postManager   *PostManager
	backupManager *BackupManager
}

// NewStagingManager creates a new StagingManager instance.
func NewStagingManager(contentDir string, postManager *PostManager, backupManager *BackupManager) *StagingManager {
	return &StagingManager{
		contentDir:    contentDir,
		postManager:   postManager,
		backupManager: backupManager,
	}
}

func (sm *StagingManager) stagingPath(postSlug string) string {
	return filepath.Join(sm.contentDir, postSlug, stagedFilename)
}

func (sm *StagingManager) postPath(postSlug string) string {
	return filepath.Join(sm.contentDir, postSlug, "index.md")
}

// SaveStaged writes staged related articles to {slug}/related-articles.staged.yaml.
//
// Requirements: 10.1, 10.2, 10.6
func (sm *StagingManager) SaveStaged(postSlug string, articles []models.RelatedArticle) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if postSlug == "" {
		return fmt.Errorf("post slug cannot be empty")
	}
	if len(articles) == 0 {
		return fmt.Errorf("articles cannot be empty")
	}

	for i, a := range articles {
		if err := a.Validate(); err != nil {
			return fmt.Errorf("invalid article %d: %w", i, err)
		}
	}

	staged := models.StagedRankings{
		PostSlug:  postSlug,
		Generated: time.Now(),
		Articles:  articles,
	}

	if err := staged.Validate(); err != nil {
		return fmt.Errorf("staged rankings validation failed: %w", err)
	}

	yamlStr, err := parser.FormatYAML(staged)
	if err != nil {
		return fmt.Errorf("failed to format staged rankings YAML: %w", err)
	}

	outPath := sm.stagingPath(postSlug)
	outDir := filepath.Dir(outPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for staging file %s: %w", outDir, err)
	}

	// Atomic write
	tempFile, err := os.CreateTemp(outDir, ".staged_*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp staging file: %w", err)
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
		return fmt.Errorf("failed to write staged rankings to temp file: %w", err)
	}

	if err := tempFile.Sync(); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to sync temp staging file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp staging file: %w", err)
	}

	if err := os.Rename(tempPath, outPath); err != nil {
		return fmt.Errorf("failed to atomically replace staging file %s: %w", outPath, err)
	}

	success = true
	return nil
}

// ReadStaged reads and validates staged rankings for a given post.
//
// Requirements: 10.6, 10.7
func (sm *StagingManager) ReadStaged(postSlug string) ([]models.RelatedArticle, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	filePath := sm.stagingPath(postSlug)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read staging file %s: %w", filePath, err)
	}

	var staged models.StagedRankings
	if err := yaml.Unmarshal(data, &staged); err != nil {
		return nil, fmt.Errorf("failed to parse staging file YAML: %w", err)
	}

	if err := staged.Validate(); err != nil {
		return nil, fmt.Errorf("invalid staging file format: %w", err)
	}

	return staged.Articles, nil
}

// ApplyStaged applies the staged rankings to the post's front matter,
// creating a backup before modifying the post and deleting the staging file upon success.
//
// Requirements: 10.3, 10.8
func (sm *StagingManager) ApplyStaged(postSlug string) error {
	articles, err := sm.ReadStaged(postSlug)
	if err != nil {
		return fmt.Errorf("cannot apply staged rankings for %s: %w", postSlug, err)
	}

	postFile := sm.postPath(postSlug)
	if _, err := os.Stat(postFile); os.IsNotExist(err) {
		return fmt.Errorf("target post file does not exist: %s", postFile)
	}

	// Create backup before applying
	if sm.backupManager != nil {
		if _, err := sm.backupManager.Backup(postFile); err != nil {
			return fmt.Errorf("backup failed prior to applying staged rankings: %w", err)
		}
	}

	// Update post front matter
	pm := sm.postManager
	if pm == nil {
		pm = NewPostManager(sm.contentDir)
	}

	if err := pm.UpdateRelatedArticles(models.PostPath(postFile), articles); err != nil {
		return fmt.Errorf("failed to update post related articles: %w", err)
	}

	// Delete staging file upon success
	stagingFile := sm.stagingPath(postSlug)
	if err := os.Remove(stagingFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove applied staging file %s: %w", stagingFile, err)
	}

	return nil
}

// ClearStaged deletes the staging files for the specified post slugs.
// If postSlugs is empty, it clears all staging files found in the content directory.
//
// Requirements: 10.4
func (sm *StagingManager) ClearStaged(postSlugs []string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if len(postSlugs) == 0 {
		var err error
		postSlugs, err = sm.listStagedLocked()
		if err != nil {
			return err
		}
	}

	var errs []string
	for _, slug := range postSlugs {
		path := sm.stagingPath(slug)
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			errs = append(errs, fmt.Sprintf("%s: %v", slug, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed clearing some staged files: %v", errs)
	}

	return nil
}

// ListStaged scans content directory and returns all slugs that have a staging file.
//
// Requirements: 10.5
func (sm *StagingManager) ListStaged() ([]string, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	return sm.listStagedLocked()
}

func (sm *StagingManager) listStagedLocked() ([]string, error) {
	var stagedSlugs []string

	err := filepath.Walk(sm.contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == stagedFilename {
			slug := filepath.Base(filepath.Dir(path))
			stagedSlugs = append(stagedSlugs, slug)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed scanning for staging files: %w", err)
	}

	return stagedSlugs, nil
}
