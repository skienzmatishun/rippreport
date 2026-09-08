package testutil

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// TempDir creates a temporary directory for testing and returns its path.
// The directory will be automatically cleaned up when the test completes.
func TempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "go-content-pipeline-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	return dir
}

// TempFile creates a temporary file with the given content and returns its path.
// The file will be automatically cleaned up when the test completes.
func TempFile(t *testing.T, content string) string {
	t.Helper()
	file, err := os.CreateTemp("", "go-content-pipeline-test-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	path := file.Name()
	if _, err := file.WriteString(content); err != nil {
		file.Close()
		t.Fatalf("failed to write to temp file: %v", err)
	}
	file.Close()
	t.Cleanup(func() {
		os.Remove(path)
	})
	return path
}

// WriteFile writes content to a file within a test directory.
// The file will be created if it doesn't exist, along with any necessary parent directories.
func WriteFile(t *testing.T, path, content string) {
	t.Helper()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create directory %s: %v", dir, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file %s: %v", path, err)
	}
}

// ReadFile reads the content of a file and returns it as a string.
// Fails the test if the file cannot be read.
func ReadFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	return string(content)
}

// FileExists checks if a file exists at the given path.
func FileExists(t *testing.T, path string) bool {
	t.Helper()
	_, err := os.Stat(path)
	return err == nil
}

// CalculateHash calculates the MD5 hash of the given content.
func CalculateHash(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}

// CopyFile copies a file from src to dst.
func CopyFile(t *testing.T, src, dst string) {
	t.Helper()
	srcFile, err := os.Open(src)
	if err != nil {
		t.Fatalf("failed to open source file %s: %v", src, err)
	}
	defer srcFile.Close()

	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create directory %s: %v", dir, err)
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		t.Fatalf("failed to create destination file %s: %v", dst, err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		t.Fatalf("failed to copy file from %s to %s: %v", src, dst, err)
	}
}

// AssertNoError fails the test if err is not nil.
func AssertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

// AssertError fails the test if err is nil.
func AssertError(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error but got nil", msg)
	}
}

// AssertEqual fails the test if got != want.
func AssertEqual(t *testing.T, got, want interface{}, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %v, want %v", msg, got, want)
	}
}

// AssertNotEqual fails the test if got == want.
func AssertNotEqual(t *testing.T, got, want interface{}, msg string) {
	t.Helper()
	if got == want {
		t.Fatalf("%s: got %v, expected different value", msg, got)
	}
}

// AssertTrue fails the test if condition is false.
func AssertTrue(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Fatalf("%s: expected true but got false", msg)
	}
}

// AssertFalse fails the test if condition is true.
func AssertFalse(t *testing.T, condition bool, msg string) {
	t.Helper()
	if condition {
		t.Fatalf("%s: expected false but got true", msg)
	}
}

// AssertContains fails the test if substring is not in str.
func AssertContains(t *testing.T, str, substring, msg string) {
	t.Helper()
	if !contains(str, substring) {
		t.Fatalf("%s: %q does not contain %q", msg, str, substring)
	}
}

// AssertNotContains fails the test if substring is in str.
func AssertNotContains(t *testing.T, str, substring, msg string) {
	t.Helper()
	if contains(str, substring) {
		t.Fatalf("%s: %q contains %q", msg, str, substring)
	}
}

func contains(str, substring string) bool {
	return len(str) >= len(substring) && (str == substring || containsHelper(str, substring))
}

func containsHelper(str, substring string) bool {
	for i := 0; i <= len(str)-len(substring); i++ {
		if str[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}

// AssertFloatEqual fails the test if got and want differ by more than epsilon.
func AssertFloatEqual(t *testing.T, got, want, epsilon float64, msg string) {
	t.Helper()
	diff := got - want
	if diff < 0 {
		diff = -diff
	}
	if diff > epsilon {
		t.Fatalf("%s: got %.10f, want %.10f (diff %.10f > epsilon %.10f)", msg, got, want, diff, epsilon)
	}
}

// CreateHugoSiteStructure creates a minimal Hugo site structure for testing.
func CreateHugoSiteStructure(t *testing.T, baseDir string) {
	t.Helper()
	contentDir := filepath.Join(baseDir, "content", "p")
	if err := os.MkdirAll(contentDir, 0755); err != nil {
		t.Fatalf("failed to create Hugo content directory: %v", err)
	}
}

// MustParseInt converts a string to an int, failing the test on error.
func MustParseInt(t *testing.T, s string) int {
	t.Helper()
	var result int
	if _, err := fmt.Sscanf(s, "%d", &result); err != nil {
		t.Fatalf("failed to parse int from %q: %v", s, err)
	}
	return result
}

// MustParseFloat converts a string to a float64, failing the test on error.
func MustParseFloat(t *testing.T, s string) float64 {
	t.Helper()
	var result float64
	if _, err := fmt.Sscanf(s, "%f", &result); err != nil {
		t.Fatalf("failed to parse float from %q: %v", s, err)
	}
	return result
}
