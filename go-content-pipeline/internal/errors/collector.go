package errors

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// PipelineSummary contains aggregate statistics for pipeline execution.
//
// Requirements: 17.6, 19.7
type PipelineSummary struct {
	TotalPosts         int
	ProcessedPosts     int
	FailedPosts        int
	SkippedPosts       int
	EmbeddingsComputed int
	CacheHits          int
	TotalAPICalls      int
	TotalDuration      time.Duration
}

// ErrorCollector gathers errors during pipeline execution grouped by category.
//
// Requirements: 17.1, 17.5, 17.6, 21.1, 22.4
type ErrorCollector struct {
	mu     sync.RWMutex
	errors []error
}

// NewErrorCollector creates a new ErrorCollector.
func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{
		errors: make([]error, 0),
	}
}

// Add appends an error to the collector in a thread-safe manner.
func (ec *ErrorCollector) Add(err error) {
	if err == nil {
		return
	}
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.errors = append(ec.errors, err)
}

// Count returns the total number of collected errors.
func (ec *ErrorCollector) Count() int {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	return len(ec.errors)
}

// Errors returns a copy of all collected errors.
func (ec *ErrorCollector) Errors() []error {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	copied := make([]error, len(ec.errors))
	copy(copied, ec.errors)
	return copied
}

// GroupByCategory groups collected errors by ErrorCategory.
func (ec *ErrorCollector) GroupByCategory() map[ErrorCategory][]error {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	grouped := make(map[ErrorCategory][]error)
	for _, err := range ec.errors {
		cat := ErrorCategoryProcessing
		if pErr, ok := err.(*PipelineError); ok {
			cat = pErr.Category
		}
		grouped[cat] = append(grouped[cat], err)
	}

	return grouped
}

// FormatSummary generates a formatted text report of the pipeline summary and errors.
func (ec *ErrorCollector) FormatSummary(stats PipelineSummary) string {
	var sb strings.Builder

	sb.WriteString("\n=================== Pipeline Summary ===================\n")
	sb.WriteString(fmt.Sprintf("Duration:            %s\n", stats.TotalDuration.Round(time.Millisecond)))
	sb.WriteString(fmt.Sprintf("Total Posts:         %d\n", stats.TotalPosts))
	sb.WriteString(fmt.Sprintf("Successfully Done:   %d\n", stats.ProcessedPosts))
	sb.WriteString(fmt.Sprintf("Failed:              %d\n", stats.FailedPosts))
	sb.WriteString(fmt.Sprintf("Skipped:             %d\n", stats.SkippedPosts))
	sb.WriteString(fmt.Sprintf("Embeddings Computed: %d\n", stats.EmbeddingsComputed))
	sb.WriteString(fmt.Sprintf("Cache Hits:          %d\n", stats.CacheHits))
	sb.WriteString(fmt.Sprintf("API Calls:           %d\n", stats.TotalAPICalls))
	sb.WriteString(fmt.Sprintf("Errors Encountered:  %d\n", ec.Count()))

	grouped := ec.GroupByCategory()
	if len(grouped) > 0 {
		sb.WriteString("\n------------------- Error Breakdown --------------------\n")
		for cat, errList := range grouped {
			sb.WriteString(fmt.Sprintf("[%s]: %d errors\n", cat, len(errList)))
			for i, e := range errList {
				if i >= 5 {
					sb.WriteString(fmt.Sprintf("  ... and %d more %s errors\n", len(errList)-5, cat))
					break
				}
				sb.WriteString(fmt.Sprintf("  - %s\n", e.Error()))
			}
		}
	}
	sb.WriteString("========================================================\n")

	return sb.String()
}
