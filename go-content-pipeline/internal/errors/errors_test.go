package errors

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestPipelineError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *PipelineError
		expected string
	}{
		{
			name: "error without post slug",
			err: &PipelineError{
				Category:  ErrorCategoryNetwork,
				Operation: "HTTP request",
				Message:   "connection refused",
			},
			expected: "[network] HTTP request: connection refused",
		},
		{
			name: "error with post slug",
			err: &PipelineError{
				Category:  ErrorCategoryFileIO,
				Operation: "read post",
				PostSlug:  "example-post",
				Message:   "file not found",
			},
			expected: "[file_io] read post for post 'example-post': file not found",
		},
		{
			name: "parsing error",
			err: &PipelineError{
				Category:  ErrorCategoryParsing,
				Operation: "parse YAML",
				PostSlug:  "test-post",
				Message:   "invalid YAML syntax",
			},
			expected: "[parsing] parse YAML for post 'test-post': invalid YAML syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("PipelineError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPipelineError_Unwrap(t *testing.T) {
	causeErr := errors.New("underlying error")
	pipelineErr := &PipelineError{
		Category:  ErrorCategoryNetwork,
		Operation: "test operation",
		Message:   "test message",
		Cause:     causeErr,
	}

	unwrapped := pipelineErr.Unwrap()
	if unwrapped != causeErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, causeErr)
	}

	// Test with errors.Is
	if !errors.Is(pipelineErr, causeErr) {
		t.Error("errors.Is() should return true for wrapped error")
	}
}

func TestPipelineError_IsRetryable(t *testing.T) {
	tests := []struct {
		name      string
		retryable bool
	}{
		{"retryable error", true},
		{"non-retryable error", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &PipelineError{
				Category:  ErrorCategoryNetwork,
				Operation: "test",
				Message:   "test",
				Retryable: tt.retryable,
			}

			if got := err.IsRetryable(); got != tt.retryable {
				t.Errorf("IsRetryable() = %v, want %v", got, tt.retryable)
			}
		})
	}
}

func TestPipelineError_WithPostSlug(t *testing.T) {
	original := &PipelineError{
		Category:  ErrorCategoryProcessing,
		Operation: "generate embedding",
		Message:   "failed to generate",
		Retryable: true,
	}

	withSlug := original.WithPostSlug("test-post")

	// Verify the new error has the slug
	if withSlug.PostSlug != "test-post" {
		t.Errorf("WithPostSlug() PostSlug = %v, want %v", withSlug.PostSlug, "test-post")
	}

	// Verify other fields are copied
	if withSlug.Category != original.Category {
		t.Error("WithPostSlug() should preserve Category")
	}
	if withSlug.Operation != original.Operation {
		t.Error("WithPostSlug() should preserve Operation")
	}
	if withSlug.Message != original.Message {
		t.Error("WithPostSlug() should preserve Message")
	}
	if withSlug.Retryable != original.Retryable {
		t.Error("WithPostSlug() should preserve Retryable")
	}

	// Verify original is unchanged
	if original.PostSlug != "" {
		t.Error("WithPostSlug() should not modify original error")
	}
}

func TestErrorCategory_Constants(t *testing.T) {
	categories := []ErrorCategory{
		ErrorCategoryNetwork,
		ErrorCategoryFileIO,
		ErrorCategoryParsing,
		ErrorCategoryValidation,
		ErrorCategoryProcessing,
		ErrorCategoryTimeout,
	}

	expectedValues := []string{
		"network",
		"file_io",
		"parsing",
		"validation",
		"processing",
		"timeout",
	}

	for i, cat := range categories {
		if string(cat) != expectedValues[i] {
			t.Errorf("ErrorCategory[%d] = %v, want %v", i, cat, expectedValues[i])
		}
	}
}

func TestNewNetworkError(t *testing.T) {
	cause := errors.New("connection timeout")
	err := NewNetworkError("fetch data", "failed to connect", cause, true)

	if err.Category != ErrorCategoryNetwork {
		t.Errorf("Category = %v, want %v", err.Category, ErrorCategoryNetwork)
	}
	if err.Operation != "fetch data" {
		t.Errorf("Operation = %v, want %v", err.Operation, "fetch data")
	}
	if err.Message != "failed to connect" {
		t.Errorf("Message = %v, want %v", err.Message, "failed to connect")
	}
	if err.Cause != cause {
		t.Errorf("Cause = %v, want %v", err.Cause, cause)
	}
	if !err.Retryable {
		t.Error("Retryable should be true")
	}
	if err.Timestamp.IsZero() {
		t.Error("Timestamp should be set")
	}
	if time.Since(err.Timestamp) > time.Second {
		t.Error("Timestamp should be recent")
	}
}

func TestNewFileIOError(t *testing.T) {
	cause := errors.New("permission denied")
	err := NewFileIOError("read file", "cannot read file", cause, false)

	if err.Category != ErrorCategoryFileIO {
		t.Errorf("Category = %v, want %v", err.Category, ErrorCategoryFileIO)
	}
	if err.Retryable {
		t.Error("Retryable should be false")
	}
}

func TestNewParsingError(t *testing.T) {
	cause := errors.New("invalid YAML")
	err := NewParsingError("parse front matter", "YAML syntax error", cause, false)

	if err.Category != ErrorCategoryParsing {
		t.Errorf("Category = %v, want %v", err.Category, ErrorCategoryParsing)
	}
	if err.Retryable {
		t.Error("Retryable should be false")
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("validate config", "missing required field", nil, false)

	if err.Category != ErrorCategoryValidation {
		t.Errorf("Category = %v, want %v", err.Category, ErrorCategoryValidation)
	}
	if err.Cause != nil {
		t.Error("Cause should be nil")
	}
}

func TestNewProcessingError(t *testing.T) {
	cause := errors.New("model not loaded")
	err := NewProcessingError("generate embedding", "embedding failed", cause, true)

	if err.Category != ErrorCategoryProcessing {
		t.Errorf("Category = %v, want %v", err.Category, ErrorCategoryProcessing)
	}
	if !err.Retryable {
		t.Error("Retryable should be true")
	}
}

func TestNewTimeoutError(t *testing.T) {
	cause := errors.New("context deadline exceeded")
	err := NewTimeoutError("API call", "request timed out", cause, true)

	if err.Category != ErrorCategoryTimeout {
		t.Errorf("Category = %v, want %v", err.Category, ErrorCategoryTimeout)
	}
	if !err.Retryable {
		t.Error("Retryable should be true")
	}
}

func TestWrap(t *testing.T) {
	cause := fmt.Errorf("original error")
	wrapped := Wrap(
		ErrorCategoryNetwork,
		"test operation",
		"wrapped message",
		cause,
		true,
	)

	if wrapped.Category != ErrorCategoryNetwork {
		t.Errorf("Category = %v, want %v", wrapped.Category, ErrorCategoryNetwork)
	}
	if wrapped.Operation != "test operation" {
		t.Errorf("Operation = %v, want %v", wrapped.Operation, "test operation")
	}
	if wrapped.Message != "wrapped message" {
		t.Errorf("Message = %v, want %v", wrapped.Message, "wrapped message")
	}
	if wrapped.Cause != cause {
		t.Errorf("Cause = %v, want %v", wrapped.Cause, cause)
	}
	if !wrapped.Retryable {
		t.Error("Retryable should be true")
	}
}

func TestIsPipelineError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "is pipeline error",
			err:      NewNetworkError("test", "test", nil, false),
			expected: true,
		},
		{
			name:     "is not pipeline error",
			err:      errors.New("standard error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPipelineError(tt.err); got != tt.expected {
				t.Errorf("IsPipelineError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetCategory(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedCat    ErrorCategory
		expectedExists bool
	}{
		{
			name:           "network error",
			err:            NewNetworkError("test", "test", nil, false),
			expectedCat:    ErrorCategoryNetwork,
			expectedExists: true,
		},
		{
			name:           "file io error",
			err:            NewFileIOError("test", "test", nil, false),
			expectedCat:    ErrorCategoryFileIO,
			expectedExists: true,
		},
		{
			name:           "standard error",
			err:            errors.New("test"),
			expectedCat:    "",
			expectedExists: false,
		},
		{
			name:           "nil error",
			err:            nil,
			expectedCat:    "",
			expectedExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cat, exists := GetCategory(tt.err)
			if cat != tt.expectedCat {
				t.Errorf("GetCategory() category = %v, want %v", cat, tt.expectedCat)
			}
			if exists != tt.expectedExists {
				t.Errorf("GetCategory() exists = %v, want %v", exists, tt.expectedExists)
			}
		})
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "retryable pipeline error",
			err:      NewNetworkError("test", "test", nil, true),
			expected: true,
		},
		{
			name:     "non-retryable pipeline error",
			err:      NewParsingError("test", "test", nil, false),
			expected: false,
		},
		{
			name:     "standard error",
			err:      errors.New("test"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRetryableError(tt.err); got != tt.expected {
				t.Errorf("IsRetryableError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrorChaining(t *testing.T) {
	// Create a chain of errors
	originalErr := errors.New("original error")
	wrappedErr := fmt.Errorf("wrapped: %w", originalErr)
	pipelineErr := NewNetworkError("test", "pipeline error", wrappedErr, true)

	// Test errors.Is works through the chain
	if !errors.Is(pipelineErr, originalErr) {
		t.Error("errors.Is should find original error in chain")
	}

	// Test we can unwrap to get the wrapped error
	if pipelineErr.Unwrap() != wrappedErr {
		t.Error("Unwrap should return wrapped error")
	}

	// Test error message contains context
	errMsg := pipelineErr.Error()
	if !strings.Contains(errMsg, "pipeline error") {
		t.Error("Error message should contain pipeline context")
	}
	if !strings.Contains(errMsg, "[network]") {
		t.Error("Error message should contain category")
	}
}

func TestErrorFormatting(t *testing.T) {
	tests := []struct {
		name           string
		err            *PipelineError
		shouldContain  []string
		shouldNotContain []string
	}{
		{
			name: "basic error",
			err: NewNetworkError(
				"fetch data",
				"connection failed",
				nil,
				true,
			),
			shouldContain: []string{"[network]", "fetch data", "connection failed"},
			shouldNotContain: []string{"for post"},
		},
		{
			name: "error with post slug",
			err: NewProcessingError(
				"generate embedding",
				"model unavailable",
				nil,
				false,
			).WithPostSlug("test-post"),
			shouldContain: []string{"[processing]", "generate embedding", "test-post", "model unavailable"},
			shouldNotContain: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMsg := tt.err.Error()
			
			for _, substr := range tt.shouldContain {
				if !strings.Contains(errMsg, substr) {
					t.Errorf("Error message should contain %q, got: %s", substr, errMsg)
				}
			}
			
			for _, substr := range tt.shouldNotContain {
				if strings.Contains(errMsg, substr) {
					t.Errorf("Error message should not contain %q, got: %s", substr, errMsg)
				}
			}
		})
	}
}

func TestTimestampIsRecent(t *testing.T) {
	err := NewNetworkError("test", "test", nil, false)
	
	if err.Timestamp.IsZero() {
		t.Fatal("Timestamp should not be zero")
	}
	
	elapsed := time.Since(err.Timestamp)
	if elapsed > time.Second {
		t.Errorf("Timestamp should be recent, got elapsed time: %v", elapsed)
	}
	if elapsed < 0 {
		t.Error("Timestamp should not be in the future")
	}
}

func TestNilCauseHandling(t *testing.T) {
	err := NewValidationError("test", "test", nil, false)
	
	if err.Cause != nil {
		t.Error("Cause should be nil")
	}
	
	// Unwrap should return nil for nil cause
	if err.Unwrap() != nil {
		t.Error("Unwrap should return nil when cause is nil")
	}
	
	// Error message should still work
	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}
}

func TestAllCategoriesAreUnique(t *testing.T) {
	categories := []ErrorCategory{
		ErrorCategoryNetwork,
		ErrorCategoryFileIO,
		ErrorCategoryParsing,
		ErrorCategoryValidation,
		ErrorCategoryProcessing,
		ErrorCategoryTimeout,
	}

	seen := make(map[ErrorCategory]bool)
	for _, cat := range categories {
		if seen[cat] {
			t.Errorf("Duplicate category found: %s", cat)
		}
		seen[cat] = true
	}

	if len(seen) != len(categories) {
		t.Errorf("Expected %d unique categories, got %d", len(categories), len(seen))
	}
}
