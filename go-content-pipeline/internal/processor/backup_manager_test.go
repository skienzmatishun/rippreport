package processor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBackupManager_BackupAndVerify(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "content", "my-post")
	backupDir := filepath.Join(tmpDir, "backups")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}

	sourceFile := filepath.Join(sourceDir, "index.md")
	content := "Original post markdown content for backup test.\nLine 2."
	if err := os.WriteFile(sourceFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	bm := NewBackupManager(backupDir)
	backupPath, err := bm.Backup(sourceFile)
	if err != nil {
		t.Fatalf("Backup failed: %v", err)
	}

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Fatalf("expected backup file at %s, but does not exist", backupPath)
	}

	// Verify checksum
	if err := bm.VerifyBackup(sourceFile, backupPath); err != nil {
		t.Errorf("VerifyBackup failed: %v", err)
	}

	// Check that backup content equals original
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("failed to read backup: %v", err)
	}
	if string(backupContent) != content {
		t.Errorf("backup content mismatch:\nGot: %s\nExpected: %s", backupContent, content)
	}
}

func TestBackupManager_ErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()
	bm := NewBackupManager(filepath.Join(tmpDir, "backups"))

	t.Run("missing source file", func(t *testing.T) {
		_, err := bm.Backup(filepath.Join(tmpDir, "nonexistent.md"))
		if err == nil {
			t.Error("expected error for missing source file, got nil")
		}
	})

	t.Run("directory instead of file", func(t *testing.T) {
		dir := filepath.Join(tmpDir, "some-dir")
		_ = os.MkdirAll(dir, 0755)
		_, err := bm.Backup(dir)
		if err == nil {
			t.Error("expected error when backing up directory, got nil")
		}
	})

	t.Run("corrupted checksum verification", func(t *testing.T) {
		fileA := filepath.Join(tmpDir, "a.txt")
		fileB := filepath.Join(tmpDir, "b.txt")
		_ = os.WriteFile(fileA, []byte("alpha"), 0644)
		_ = os.WriteFile(fileB, []byte("beta"), 0644)

		err := bm.VerifyBackup(fileA, fileB)
		if err == nil {
			t.Error("expected checksum mismatch error, got nil")
		}
	})
}
