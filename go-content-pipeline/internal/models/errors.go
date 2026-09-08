package models

import (
	"fmt"
)

// ErrorCategory represents the category of an error
type ErrorCategory string

const (
	// ErrorCategoryNetwork indicates a network-related error
	ErrorCategoryNetwork ErrorCategory = "network"
	// ErrorCategoryFileIO indicates a file I/O error
	ErrorCategoryFileIO ErrorCategory = "file_io"
	// ErrorCategoryParsing indicates a parsing error
	ErrorCategoryParsing ErrorCategory = "parsing"
	// ErrorCategoryValidation indicates a validation error
	ErrorCategoryValidation ErrorCategory = "validation"
	// ErrorCategoryProcessing indicates a processing error
	ErrorCategoryProcessing ErrorCategory = "processing"
	// ErrorCategoryTimeout indicates a timeout error
	ErrorCategoryTimeout ErrorCategory = "timeout"
)

// PipelineError represents an error that occurred during pipeline execution
type PipelineError struct {
	Category  ErrorCategory
	Operation string
	Message   string
	Retryable bool
	Cause     error
}

// Error implements the error interface
func (e *PipelineError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %s (caused by: %v)", e.Category, e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Category, e.Operation, e.Message)
}

// Unwrap returns the underlying cause error
func (e *PipelineError) Unwrap() error {
	return e.Cause
}

// IsRetryable checks if the error is retryable
func (e *PipelineError) IsRetryable() bool {
	return e.Retryable
}

// NewNetworkError creates a network-related error
func NewNetworkError(operation, message string, cause error) *PipelineError {
	return &PipelineError{
		Category:  ErrorCategoryNetwork,
		Operation: operation,
		Message:   message,
		Retryable: true,
		Cause:     cause,
	}
}

// NewFileIOError creates a file I/O error
func NewFileIOError(operation, message string, cause error) *PipelineError {
	return &PipelineError{
		Category:  ErrorCategoryFileIO,
		Operation: operation,
		Message:   message,
		Retryable: false,
		Cause:     cause,
	}
}

// NewParsingError creates a parsing error
func NewParsingError(operation, message string, cause error) *PipelineError {
	return &PipelineError{
		Category:  ErrorCategoryParsing,
		Operation: operation,
		Message:   message,
		Retryable: false,
		Cause:     cause,
	}
}

// NewValidationError creates a validation error
func NewValidationError(operation, message string, cause error) *PipelineError {
	return &PipelineError{
		Category:  ErrorCategoryValidation,
		Operation: operation,
		Message:   message,
		Retryable: false,
		Cause:     cause,
	}
}

// NewProcessingError creates a processing error
func NewProcessingError(operation, message string, cause error, retryable bool) *PipelineError {
	return &PipelineError{
		Category:  ErrorCategoryProcessing,
		Operation: operation,
		Message:   message,
		Retryable: retryable,
		Cause:     cause,
	}
}

// NewTimeoutError creates a timeout error
func NewTimeoutError(operation, message string, cause error) *PipelineError {
	return &PipelineError{
		Category:  ErrorCategoryTimeout,
		Operation: operation,
		Message:   message,
		Retryable: true,
		Cause:     cause,
	}
}

// ProcessingError represents an error during post processing
type ProcessingError struct {
	PostSlug string
	Phase    string
	Error    error
}

// LlamaError represents an error from the llama.cpp server
type LlamaError struct {
	Code       ErrorCode
	Message    string
	Retryable  bool
	StatusCode int
	Cause      error
}

// ErrorCode represents specific error codes from llama.cpp
type ErrorCode int

const (
	// ErrConnection indicates a connection error
	ErrConnection ErrorCode = iota
	// ErrTimeout indicates a timeout error
	ErrTimeout
	// ErrServerError indicates a 5xx server error
	ErrServerError
	// ErrClientError indicates a 4xx client error
	ErrClientError
	// ErrModelNotLoaded indicates the model is not loaded
	ErrModelNotLoaded
	// ErrContextTooLarge indicates the context is too large
	ErrContextTooLarge
	// ErrMalformedResponse indicates a malformed response
	ErrMalformedResponse
)

// Error implements the error interface
func (e *LlamaError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("llama error [%d]: %s (status: %d, caused by: %v)", 
			e.Code, e.Message, e.StatusCode, e.Cause)
	}
	return fmt.Sprintf("llama error [%d]: %s (status: %d)", 
		e.Code, e.Message, e.StatusCode)
}

// Unwrap returns the underlying cause error
func (e *LlamaError) Unwrap() error {
	return e.Cause
}

// IsRetryable checks if the error is retryable
func (e *LlamaError) IsRetryable() bool {
	return e.Retryable
}
