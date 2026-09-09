package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// ServerMode defines the operational mode of llama-server.
type ServerMode string

const (
	ModeEmbedding  ServerMode = "embedding"
	ModeReranker   ServerMode = "reranker"
	ModeCompletion ServerMode = "completion"
)

const (
	defaultBinaryPath     = "/Users/ryandunphy/.unsloth/llama.cpp/build/bin/llama-server"
	defaultModelsDir      = "/Volumes/1tb/models"
	defaultHost           = "127.0.0.1"
	defaultPort           = 8080
	defaultContextSize    = 4096
	defaultGPULayers      = 99
	defaultStartupTimeout = 0 // 0 means no timeout: wait indefinitely until server is ready or process exits
)

// ManagerConfig holds configuration for the ServerManager.
type ManagerConfig struct {
	Managed        bool          `yaml:"managed"`
	BinaryPath     string        `yaml:"binary_path"`
	ModelsDir      string        `yaml:"models_dir"`
	Host           string        `yaml:"host"`
	Port           int           `yaml:"port"`
	ContextSize    int           `yaml:"context_size"`
	GPULayers      int           `yaml:"gpu_layers"`
	StartupTimeout time.Duration `yaml:"startup_timeout"`
	// Speculative decoding options
	DraftModel    string `yaml:"draft_model,omitempty"`      // Path/name of draft model for speculative decoding
	SpecType      string `yaml:"spec_type,omitempty"`        // Type: draft-dflash, draft-eagle3, draft-dspark, ngram-mod, etc.
	SpecDraftNMax int    `yaml:"spec_draft_n_max,omitempty"` // Max number of tokens to draft (default: 3)
	// Qwen 3.8 specific options
	ReasoningEffort  string `yaml:"reasoning_effort,omitempty"`  // For Qwen 3.8: xhigh, high, medium, low, none
	PreserveThinking bool   `yaml:"preserve_thinking,omitempty"` // For Qwen 3.8: preserve thinking traces in conversation
}

// DefaultManagerConfig returns default settings for ServerManager.
func DefaultManagerConfig() ManagerConfig {
	return ManagerConfig{
		Managed:        true,
		BinaryPath:     defaultBinaryPath,
		ModelsDir:      defaultModelsDir,
		Host:           defaultHost,
		Port:           defaultPort,
		ContextSize:    defaultContextSize,
		GPULayers:      defaultGPULayers,
		StartupTimeout: defaultStartupTimeout,
	}
}

// ServerManager manages the lifecycle of llama-server subprocesses.
type ServerManager struct {
	cfg          ManagerConfig
	mu           sync.Mutex
	cmd          *exec.Cmd
	exitCh       chan error
	currentModel string
	currentMode  ServerMode
	currentPath  string
	stderrBuf    *bytes.Buffer
}

// NewServerManager creates a new ServerManager instance.
func NewServerManager(cfg ManagerConfig) *ServerManager {
	if cfg.BinaryPath == "" {
		cfg.BinaryPath = defaultBinaryPath
	}
	if cfg.ModelsDir == "" {
		cfg.ModelsDir = defaultModelsDir
	}
	if cfg.Host == "" {
		cfg.Host = defaultHost
	}
	if cfg.Port <= 0 {
		cfg.Port = defaultPort
	}
	if cfg.ContextSize <= 0 {
		cfg.ContextSize = defaultContextSize
	}
	if cfg.GPULayers == 0 {
		cfg.GPULayers = defaultGPULayers
	}
	if cfg.StartupTimeout < 0 {
		cfg.StartupTimeout = 0
	}

	return &ServerManager{
		cfg: cfg,
	}
}

// GetBaseURL returns the HTTP URL for the managed llama-server.
func (sm *ServerManager) GetBaseURL() string {
	return fmt.Sprintf("http://%s:%d", sm.cfg.Host, sm.cfg.Port)
}

// ResolveModelPath finds the actual GGUF file path from a model path, relative path, or model name.
func (sm *ServerManager) ResolveModelPath(modelPathOrName string) (string, error) {
	if modelPathOrName == "" {
		return "", fmt.Errorf("model path or name cannot be empty")
	}

	// 1. Check if absolute path that exists
	if filepath.IsAbs(modelPathOrName) {
		if _, err := os.Stat(modelPathOrName); err == nil {
			return modelPathOrName, nil
		}
	}

	// 2. Check directly under ModelsDir
	joined := filepath.Join(sm.cfg.ModelsDir, modelPathOrName)
	if _, err := os.Stat(joined); err == nil {
		return joined, nil
	}

	// 3. Search recursively under ModelsDir for matching .gguf file
	if _, err := os.Stat(sm.cfg.ModelsDir); err == nil {
		cleanTarget := strings.ToLower(modelPathOrName)
		var bestMatch string
		_ = filepath.Walk(sm.cfg.ModelsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if strings.HasSuffix(strings.ToLower(info.Name()), ".gguf") {
				lowerPath := strings.ToLower(path)
				lowerName := strings.ToLower(info.Name())
				// Check exact filename match (with or without .gguf)
				if lowerName == cleanTarget || lowerName == cleanTarget+".gguf" {
					bestMatch = path
					return io.EOF // Found exact match
				}
				// Check substring match
				if strings.Contains(lowerPath, cleanTarget) || strings.Contains(lowerName, cleanTarget) {
					if bestMatch == "" {
						bestMatch = path
					}
				}
			}
			return nil
		})
		if bestMatch != "" {
			return bestMatch, nil
		}
	}

	// If not found in ModelsDir, and it was a relative path, check current working dir
	if _, err := os.Stat(modelPathOrName); err == nil {
		abs, aErr := filepath.Abs(modelPathOrName)
		if aErr == nil {
			return abs, nil
		}
		return modelPathOrName, nil
	}

	return "", fmt.Errorf("model %q not found (checked %s and %s)", modelPathOrName, sm.cfg.ModelsDir, modelPathOrName)
}

// EnsureModel ensures that llama-server is running with the specified model and mode.
// If the server is already running with this model and mode, it returns immediately.
// Otherwise, it terminates the existing server and launches a new one with the target model.
func (sm *ServerManager) EnsureModel(ctx context.Context, modelPathOrName string, mode ServerMode) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if !sm.cfg.Managed {
		// Unmanaged mode: verify external server is healthy
		return sm.checkHealth(ctx)
	}

	// Check if already running with the target model and mode
	if sm.cmd != nil && sm.cmd.Process != nil && sm.currentModel == modelPathOrName && sm.currentMode == mode {
		if err := sm.checkHealth(ctx); err == nil {
			return nil
		}
	}

	// Unload / stop current server if running
	if sm.cmd != nil {
		_ = sm.stopLocked()
	}

	// Resolve the model path
	resolvedPath, err := sm.ResolveModelPath(modelPathOrName)
	if err != nil {
		return fmt.Errorf("failed to resolve model: %w", err)
	}

	// Build arguments
	args := []string{
		"-m", resolvedPath,
		"--host", sm.cfg.Host,
		"--port", fmt.Sprintf("%d", sm.cfg.Port),
		"-c", fmt.Sprintf("%d", sm.cfg.ContextSize),
		"-ngl", fmt.Sprintf("%d", sm.cfg.GPULayers),
	}

	switch mode {
	case ModeEmbedding:
		args = append(args, "--embedding")
	case ModeReranker:
		args = append(args, "--reranking")
	case ModeCompletion:
		args = append(args, "--jinja")
	}

	// Add speculative decoding options if configured
	if sm.cfg.DraftModel != "" && sm.cfg.SpecType != "" {
		draftPath, err := sm.ResolveModelPath(sm.cfg.DraftModel)
		if err != nil {
			return fmt.Errorf("failed to resolve draft model: %w", err)
		}
		args = append(args, "-md", draftPath)
		args = append(args, "--spec-type", sm.cfg.SpecType)

		if sm.cfg.SpecDraftNMax > 0 {
			args = append(args, "--spec-draft-n-max", fmt.Sprintf("%d", sm.cfg.SpecDraftNMax))
		}
	}

	// Add Qwen 3.8 specific options via chat template kwargs
	if sm.cfg.ReasoningEffort != "" || sm.cfg.PreserveThinking {
		var kwargs []string
		if sm.cfg.ReasoningEffort != "" {
			kwargs = append(kwargs, fmt.Sprintf(`"reasoning_effort":"%s"`, sm.cfg.ReasoningEffort))
		}
		if sm.cfg.PreserveThinking {
			kwargs = append(kwargs, `"preserve_thinking":true`)
		}
		if len(kwargs) > 0 {
			kwargsStr := "{" + strings.Join(kwargs, ",") + "}"
			args = append(args, "--chat-template-kwargs", kwargsStr)
		}
	}

	sm.stderrBuf = new(bytes.Buffer)
	cmd := exec.Command(sm.cfg.BinaryPath, args...)
	cmd.Stderr = sm.stderrBuf
	cmd.Stdout = io.Discard

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start llama-server (%s): %w", sm.cfg.BinaryPath, err)
	}

	exitCh := make(chan error, 1)
	go func(c *exec.Cmd, ch chan error) {
		ch <- c.Wait()
		close(ch)
	}(cmd, exitCh)

	sm.cmd = cmd
	sm.exitCh = exitCh
	sm.currentModel = modelPathOrName
	sm.currentMode = mode
	sm.currentPath = resolvedPath

	// Poll health until server is ready (no timeout if StartupTimeout <= 0)
	var startupCtx context.Context
	var cancel context.CancelFunc
	if sm.cfg.StartupTimeout > 0 {
		startupCtx, cancel = context.WithTimeout(ctx, sm.cfg.StartupTimeout)
	} else {
		startupCtx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	startTime := time.Now()
	lastProgress := time.Now()

	for {
		select {
		case <-startupCtx.Done():
			_ = sm.stopLocked()
			stderrMsg := sm.stderrBuf.String()
			if len(stderrMsg) > 500 {
				stderrMsg = stderrMsg[len(stderrMsg)-500:]
			}
			if errors.Is(startupCtx.Err(), context.Canceled) {
				return fmt.Errorf("interrupted waiting for llama-server to load model %s: %w",
					filepath.Base(resolvedPath), startupCtx.Err())
			}
			return fmt.Errorf("timeout waiting for llama-server to load model %s (mode: %s): %w\nOutput: %s",
				filepath.Base(resolvedPath), mode, startupCtx.Err(), stderrMsg)
		case err := <-exitCh:
			_ = sm.stopLocked()
			stderrMsg := sm.stderrBuf.String()
			if len(stderrMsg) > 1000 {
				stderrMsg = stderrMsg[len(stderrMsg)-1000:]
			}
			return fmt.Errorf("llama-server exited unexpectedly (%v):\n%s", err, stderrMsg)
		case <-ticker.C:
			if err := sm.checkHealth(ctx); err == nil {
				return nil
			}
			if time.Since(lastProgress) >= 5*time.Second {
				fmt.Printf("   ... loading model %s into memory (%v elapsed)\n", filepath.Base(resolvedPath), time.Since(startTime).Round(time.Second))
				lastProgress = time.Now()
			}
		}
	}
}

// checkHealth sends a GET request to the /health endpoint.
func (sm *ServerManager) checkHealth(ctx context.Context) error {
	url := fmt.Sprintf("http://%s:%d/health", sm.cfg.Host, sm.cfg.Port)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}
	return nil
}

// Stop terminates the running llama-server process.
func (sm *ServerManager) Stop() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.stopLocked()
}

// stopLocked stops the process assuming the caller holds the lock.
func (sm *ServerManager) stopLocked() error {
	if sm.cmd == nil || sm.cmd.Process == nil {
		return nil
	}

	proc := sm.cmd.Process
	exitCh := sm.exitCh
	// Attempt graceful termination
	_ = proc.Signal(syscall.SIGTERM)

	if exitCh != nil {
		select {
		case <-exitCh:
			// Clean exit
		case <-time.After(5 * time.Second):
			// Force kill if not exited within 5 seconds
			_ = proc.Kill()
			<-exitCh
		}
	}

	sm.cmd = nil
	sm.exitCh = nil
	sm.currentModel = ""
	sm.currentMode = ""
	sm.currentPath = ""

	// Brief pause to allow OS port release
	time.Sleep(200 * time.Millisecond)
	return nil
}
