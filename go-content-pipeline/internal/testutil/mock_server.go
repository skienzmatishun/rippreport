package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// MockLlamaServer represents a mock llama.cpp HTTP server for testing.
type MockLlamaServer struct {
	Server *httptest.Server
	
	// Request tracking
	mu                sync.Mutex
	EmbeddingRequests []EmbeddingRequest
	CompletionRequests []CompletionRequest
	HealthCheckCount   int
	
	// Response configuration
	EmbeddingResponse  func(EmbeddingRequest) ([]float32, error)
	CompletionResponse func(CompletionRequest) (string, error)
	HealthCheckOK      bool
	
	// Error simulation
	SimulateTimeout    bool
	SimulateServerError bool
	SimulateRateLimitError bool
	ErrorAfterN        int
	requestCount       int
}

// EmbeddingRequest represents an embedding API request.
type EmbeddingRequest struct {
	Content string `json:"content"`
}

// CompletionRequest represents a completion API request.
type CompletionRequest struct {
	Prompt      string   `json:"prompt"`
	Temperature float32  `json:"temperature,omitempty"`
	TopP        float32  `json:"top_p,omitempty"`
	TopK        int      `json:"top_k,omitempty"`
	MaxTokens   int      `json:"n_predict,omitempty"`
	Stop        []string `json:"stop,omitempty"`
}

// NewMockLlamaServer creates a new mock llama.cpp server.
func NewMockLlamaServer(t *testing.T) *MockLlamaServer {
	t.Helper()
	
	mock := &MockLlamaServer{
		EmbeddingRequests:  []EmbeddingRequest{},
		CompletionRequests: []CompletionRequest{},
		HealthCheckOK:      true,
	}
	
	// Default embedding response: return simple vector
	mock.EmbeddingResponse = func(req EmbeddingRequest) ([]float32, error) {
		return TestEmbedding(384, 0.5), nil
	}
	
	// Default completion response: return score
	mock.CompletionResponse = func(req CompletionRequest) (string, error) {
		return "85", nil
	}
	
	mock.Server = httptest.NewServer(http.HandlerFunc(mock.handler))
	
	t.Cleanup(func() {
		mock.Server.Close()
	})
	
	return mock
}

func (m *MockLlamaServer) handler(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	m.requestCount++
	reqCount := m.requestCount
	m.mu.Unlock()
	
	// Simulate errors after N requests
	if m.ErrorAfterN > 0 && reqCount > m.ErrorAfterN {
		http.Error(w, "simulated error", http.StatusInternalServerError)
		return
	}
	
	// Simulate timeout (client will timeout, not server response)
	if m.SimulateTimeout {
		// In real testing, the client timeout will trigger
		http.Error(w, "timeout simulation", http.StatusRequestTimeout)
		return
	}
	
	// Simulate server error
	if m.SimulateServerError {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	
	// Simulate rate limit error
	if m.SimulateRateLimitError {
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	
	switch r.URL.Path {
	case "/embedding":
		m.handleEmbedding(w, r)
	case "/completion":
		m.handleCompletion(w, r)
	case "/health":
		m.handleHealthCheck(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (m *MockLlamaServer) handleEmbedding(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req EmbeddingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	
	m.mu.Lock()
	m.EmbeddingRequests = append(m.EmbeddingRequests, req)
	m.mu.Unlock()
	
	embedding, err := m.EmbeddingResponse(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"embedding": embedding,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (m *MockLlamaServer) handleCompletion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req CompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	
	m.mu.Lock()
	m.CompletionRequests = append(m.CompletionRequests, req)
	m.mu.Unlock()
	
	content, err := m.CompletionResponse(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"content":           content,
		"tokens_evaluated":  len(strings.Split(req.Prompt, " ")),
		"tokens_predicted":  len(strings.Split(content, " ")),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (m *MockLlamaServer) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	m.HealthCheckCount++
	m.mu.Unlock()
	
	if !m.HealthCheckOK {
		http.Error(w, "unhealthy", http.StatusServiceUnavailable)
		return
	}
	
	response := map[string]interface{}{
		"status":       "ok",
		"model_loaded": true,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetEmbeddingRequests returns all embedding requests received.
func (m *MockLlamaServer) GetEmbeddingRequests() []EmbeddingRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	requests := make([]EmbeddingRequest, len(m.EmbeddingRequests))
	copy(requests, m.EmbeddingRequests)
	return requests
}

// GetCompletionRequests returns all completion requests received.
func (m *MockLlamaServer) GetCompletionRequests() []CompletionRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	requests := make([]CompletionRequest, len(m.CompletionRequests))
	copy(requests, m.CompletionRequests)
	return requests
}

// GetHealthCheckCount returns the number of health check requests.
func (m *MockLlamaServer) GetHealthCheckCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.HealthCheckCount
}

// Reset clears all tracked requests.
func (m *MockLlamaServer) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.EmbeddingRequests = []EmbeddingRequest{}
	m.CompletionRequests = []CompletionRequest{}
	m.HealthCheckCount = 0
	m.requestCount = 0
	m.SimulateTimeout = false
	m.SimulateServerError = false
	m.SimulateRateLimitError = false
	m.ErrorAfterN = 0
}

// SetEmbeddingHandler sets a custom embedding response handler.
func (m *MockLlamaServer) SetEmbeddingHandler(fn func(EmbeddingRequest) ([]float32, error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.EmbeddingResponse = fn
}

// SetCompletionHandler sets a custom completion response handler.
func (m *MockLlamaServer) SetCompletionHandler(fn func(CompletionRequest) (string, error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CompletionResponse = fn
}

// URL returns the server's base URL.
func (m *MockLlamaServer) URL() string {
	return m.Server.URL
}

// AssertEmbeddingRequestCount fails if the number of embedding requests doesn't match expected.
func (m *MockLlamaServer) AssertEmbeddingRequestCount(t *testing.T, expected int, msg string) {
	t.Helper()
	got := len(m.GetEmbeddingRequests())
	if got != expected {
		t.Errorf("%s: got %d embedding requests, want %d", msg, got, expected)
	}
}

// AssertCompletionRequestCount fails if the number of completion requests doesn't match expected.
func (m *MockLlamaServer) AssertCompletionRequestCount(t *testing.T, expected int, msg string) {
	t.Helper()
	got := len(m.GetCompletionRequests())
	if got != expected {
		t.Errorf("%s: got %d completion requests, want %d", msg, got, expected)
	}
}

// AssertHealthCheckCount fails if the number of health checks doesn't match expected.
func (m *MockLlamaServer) AssertHealthCheckCount(t *testing.T, expected int, msg string) {
	t.Helper()
	got := m.GetHealthCheckCount()
	if got != expected {
		t.Errorf("%s: got %d health checks, want %d", msg, got, expected)
	}
}

// WaitForRequests waits until at least N total requests have been received.
// Useful for testing concurrent operations.
func (m *MockLlamaServer) WaitForRequests(minRequests int) {
	for {
		m.mu.Lock()
		total := len(m.EmbeddingRequests) + len(m.CompletionRequests) + m.HealthCheckCount
		m.mu.Unlock()
		
		if total >= minRequests {
			return
		}
		
		// Small sleep to avoid busy waiting
		// In real tests, use proper synchronization
	}
}

// SimulateModelNotLoaded configures the server to return model not loaded errors.
func (m *MockLlamaServer) SimulateModelNotLoaded() {
	m.SetEmbeddingHandler(func(req EmbeddingRequest) ([]float32, error) {
		return nil, fmt.Errorf("model not loaded")
	})
	m.SetCompletionHandler(func(req CompletionRequest) (string, error) {
		return "", fmt.Errorf("model not loaded")
	})
	m.HealthCheckOK = false
}

// SimulateContextTooLarge configures the server to return context window exceeded errors.
func (m *MockLlamaServer) SimulateContextTooLarge() {
	m.SetEmbeddingHandler(func(req EmbeddingRequest) ([]float32, error) {
		if len(req.Content) > 1000 {
			return nil, fmt.Errorf("context window exceeded")
		}
		return TestEmbedding(384, 0.5), nil
	})
}

// SimulateBatchScoring configures the server to handle batch scoring requests.
func (m *MockLlamaServer) SimulateBatchScoring(scores []int) {
	m.SetCompletionHandler(func(req CompletionRequest) (string, error) {
		// Return scores as JSON array
		scoresJSON, _ := json.Marshal(map[string][]int{"scores": scores})
		return string(scoresJSON), nil
	})
}

// SimulateBatchHooks configures the server to handle batch hook generation requests.
func (m *MockLlamaServer) SimulateBatchHooks(hooks []string) {
	m.SetCompletionHandler(func(req CompletionRequest) (string, error) {
		// Return hooks as JSON array
		hooksJSON, _ := json.Marshal(map[string][]string{"hooks": hooks})
		return string(hooksJSON), nil
	})
}
