package processor

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// BackupManager handles creating verified backups of files before modification.
//
// Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 8.6, 8.7
type BackupManager struct {
	mu        sync.Mutex
	backupDir string
}

// NewBackupManager creates a new BackupManager with the specified backup directory.
func NewBackupManager(backupDir string) *BackupManager {
	return &BackupManager{
		backupDir: backupDir,
	}
}

// Backup creates a timestamped copy of filePath in the backup directory and verifies its checksum.
// Returns the path of the created backup file on success.
//
// Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 8.6
func (bm *BackupManager) Backup(filePath string) (string, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	// Verify original file exists and is readable
	origInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to access original file %s: %w", filePath, err)
	}
	if origInfo.IsDir() {
		return "", fmt.Errorf("cannot backup directory: %s", filePath)
	}

	timestamp := time.Now().Format("20060102-150405")
	fileName := filepath.Base(filePath)
	parentDirName := filepath.Base(filepath.Dir(filePath))

	// Preserve directory structure: {backupDir}/{timestamp}/{parentDir}/{filename}
	destDir := filepath.Join(bm.backupDir, timestamp, parentDirName)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory structure %s: %w", destDir, err)
	}

	backupPath := filepath.Join(destDir, fileName)

	// Copy file contents
	if err := copyFile(filePath, backupPath); err != nil {
		_ = os.Remove(backupPath)
		return "", fmt.Errorf("failed to copy file to backup: %w", err)
	}

	// Verify backup via checksum
	if err := bm.verifyChecksum(filePath, backupPath); err != nil {
		_ = os.Remove(backupPath)
		return "", fmt.Errorf("backup verification failed: %w", err)
	}

	return backupPath, nil
}

// VerifyBackup checks that original and backup files exist and have identical MD5 checksums.
//
// Requirements: 8.6
func (bm *BackupManager) VerifyBackup(originalPath, backupPath string) error {
	return bm.verifyChecksum(originalPath, backupPath)
}

func (bm *BackupManager) verifyChecksum(fileA, fileB string) error {
	hashA, err := calculateFileMD5(fileA)
	if err != nil {
		return fmt.Errorf("failed calculating checksum for %s: %w", fileA, err)
	}

	hashB, err := calculateFileMD5(fileB)
	if err != nil {
		return fmt.Errorf("failed calculating checksum for %s: %w", fileB, err)
	}

	if hashA != hashB {
		return fmt.Errorf("checksum mismatch between %s (%s) and %s (%s)", fileA, hashA, fileB, hashB)
	}

	return nil
}

func calculateFileMD5(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := md5.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}
