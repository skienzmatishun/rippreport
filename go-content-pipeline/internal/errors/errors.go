package errors

import (
	"fmt"
	"time"
)

// ErrorCategory represents the category of an error for classification and handling
type ErrorCategory string

const (
	// ErrorCategoryNetwork represents network-related errors (connection, timeout, etc.)
	ErrorCategoryNetwork ErrorCategory = "network"
	
	// ErrorCategoryFileIO represents file I/O errors (read, write, permissions, etc.)
	ErrorCategoryFileIO ErrorCategory = "file_io"
	
	// ErrorCategoryParsing represents parsing errors (YAML, JSON, front matter, etc.)
	ErrorCategoryParsing ErrorCategory = "parsing"
	
	// ErrorCategoryValidation represents validation errors (invalid data, missing fields, etc.)
	ErrorCategoryValidation ErrorCategory = "validation"
	
	// ErrorCategoryProcessing represents processing errors (embedding generation, scoring, etc.)
	ErrorCategoryProcessing ErrorCategory = "processing"
	
	// ErrorCategoryTimeout represents timeout errors (API calls, operations exceeding time limits)
	ErrorCategoryTimeout ErrorCategory = "timeout"
)

// PipelineError represents a structured error in the content pipeline with context
type PipelineError struct {
	// Category classifies the error for handling and reporting
	Category ErrorCategory
	
	// Operation describes what operation was being performed when the error occurred
	Operation string
	
	// PostSlug identifies the post being processed (empty if not post-specific)
	PostSlug string
	
	// Message provides a human-readable description of the error
	Message string
	
	// Cause is the underlying error that caused this error (can be nil)
	Cause error
	
	// Retryable indicates whether the operation can be retried
	Retryable bool
	
	// Timestamp records when the error occurred
	Timestamp time.Time
}

// Error implements the error interface
func (e *PipelineError) Error() string {
	if e.PostSlug != "" {
		return fmt.Sprintf("[%s] %s for post '%s': %s",
			e.Category, e.Operation, e.PostSlug, e.Message)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Category, e.Operation, e.Message)
}

// Unwrap returns the underlying cause error, enabling error unwrapping
func (e *PipelineError) Unwrap() error {
	return e.Cause
}

// IsRetryable returns whether this error indicates an operation that can be retried
func (e *PipelineError) IsRetryable() bool {
	return e.Retryable
}

// WithPostSlug returns a new PipelineError with the PostSlug field set
func (e *PipelineError) WithPostSlug(slug string) *PipelineError {
	newErr := *e
	newErr.PostSlug = slug
	return &newErr
}

// NewNetworkError creates a new network-related error
func NewNetworkError(operation string, message string, cause error, retryable bool) *PipelineError {
	return &PipelineError{
		Category:  ErrorCategoryNetwork,
		Operation: operation,
		Message:   message,
		Cause:     cause,
		Retryable: retryable,
		Timestamp: time.Now(),
	}
}

// NewFileIOError creates a new file I/O error
func NewFileIOError(operation string, message string, cause error, retryable bool) *PipelineError {
	return &PipelineError{
		Category:  ErrorCategoryFileIO,
		Operation: operation,
		Message:   message,
		Cause:     cause,
		Retryable: retryable,
		Timestamp: time.Now(),
	}
}

// NewParsingError creates a new parsing error
func NewParsingError(operation string, message string, cause error, retryable bool) *PipelineError {
	return &PipelineError{
		Category:  ErrorCategoryParsing,
		Operation: operation,
		Message:   message,
		Cause:     cause,
		Retryable: retryable,
		Timestamp: time.Now(),
	}
}

// NewValidationError creates a new validation error
func NewValidationError(operation string, message string, cause error, retryable bool) *PipelineError {
	return &PipelineError{
		Category:  ErrorCategoryValidation,
		Operation: operation,
		Message:   message,
		Cause:     cause,
		Retryable: retryable,
		Timestamp: time.Now(),
	}
}

// NewProcessingError creates a new processing error
func NewProcessingError(operation string, message string, cause error, retryable bool) *PipelineError {
	return &PipelineError{
		Category:  ErrorCategoryProcessing,
		Operation: operation,
		Message:   message,
		Cause:     cause,
		Retryable: retryable,
		Timestamp: time.Now(),
	}
}

// NewTimeoutError creates a new timeout error
func NewTimeoutError(operation string, message string, cause error, retryable bool) *PipelineError {
	return &PipelineError{
		Category:  ErrorCategoryTimeout,
		Operation: operation,
		Message:   message,
		Cause:     cause,
		Retryable: retryable,
		Timestamp: time.Now(),
	}
}

// Wrap wraps an existing error with pipeline context
func Wrap(category ErrorCategory, operation string, message string, cause error, retryable bool) *PipelineError {
	return &PipelineError{
		Category:  category,
		Operation: operation,
		Message:   message,
		Cause:     cause,
		Retryable: retryable,
		Timestamp: time.Now(),
	}
}

// IsPipelineError checks if an error is a PipelineError
func IsPipelineError(err error) bool {
	_, ok := err.(*PipelineError)
	return ok
}

// GetCategory extracts the error category from an error if it's a PipelineError
func GetCategory(err error) (ErrorCategory, bool) {
	if pErr, ok := err.(*PipelineError); ok {
		return pErr.Category, true
	}
	return "", false
}

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	if pErr, ok := err.(*PipelineError); ok {
		return pErr.Retryable
	}
	return false
}
