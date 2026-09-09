package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestServerManager_ResolveModelPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested model files
	subDir := filepath.Join(tmpDir, "org", "repo")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	modelPath := filepath.Join(subDir, "qwen3-7b-instruct.Q4_K_M.gguf")
	if err := os.WriteFile(modelPath, []byte("mock-gguf"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	sm := NewServerManager(ManagerConfig{
		ModelsDir: tmpDir,
	})

	// 1. Resolve by exact relative path
	p1, err := sm.ResolveModelPath("org/repo/qwen3-7b-instruct.Q4_K_M.gguf")
	if err != nil {
		t.Fatalf("Resolve by relative path failed: %v", err)
	}
	if p1 != modelPath {
		t.Errorf("expected %s, got %s", modelPath, p1)
	}

	// 2. Resolve by absolute path
	p2, err := sm.ResolveModelPath(modelPath)
	if err != nil {
		t.Fatalf("Resolve by absolute path failed: %v", err)
	}
	if p2 != modelPath {
		t.Errorf("expected %s, got %s", modelPath, p2)
	}

	// 3. Resolve by partial model name search
	p3, err := sm.ResolveModelPath("qwen3-7b")
	if err != nil {
		t.Fatalf("Resolve by partial name failed: %v", err)
	}
	if p3 != modelPath {
		t.Errorf("expected %s, got %s", modelPath, p3)
	}

	// 4. Missing model returns error
	_, err = sm.ResolveModelPath("nonexistent-model-xyz")
	if err == nil {
		t.Errorf("expected error for nonexistent model")
	}
}

func TestServerManager_UnmanagedMode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	sm := NewServerManager(ManagerConfig{
		Managed: false,
		Host:    "127.0.0.1",
		Port:    8080,
	})

	// Unmanaged check against closed port returns error
	err := sm.EnsureModel(context.Background(), "any-model", ModeCompletion)
	if err == nil {
		t.Errorf("expected error when external server is not reachable")
	}
}

func TestDefaultManagerConfig(t *testing.T) {
	cfg := DefaultManagerConfig()
	if !cfg.Managed {
		t.Errorf("expected Managed true by default")
	}
	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
	if cfg.ModelsDir != defaultModelsDir {
		t.Errorf("expected default models dir %s, got %s", defaultModelsDir, cfg.ModelsDir)
	}
	if cfg.BinaryPath != defaultBinaryPath {
		t.Errorf("expected default binary path %s, got %s", defaultBinaryPath, cfg.BinaryPath)
	}
	if cfg.StartupTimeout != 0 {
		t.Errorf("expected 0 startup timeout (no timeout), got %v", cfg.StartupTimeout)
	}
}
