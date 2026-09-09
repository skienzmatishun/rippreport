package errors

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestErrorCollector(t *testing.T) {
	ec := NewErrorCollector()

	ec.Add(NewNetworkError("fetch", "conn reset", nil, true))
	ec.Add(NewFileIOError("save", "disk full", nil, false))
	ec.Add(errors.New("plain error"))

	if ec.Count() != 3 {
		t.Errorf("expected count 3, got %d", ec.Count())
	}

	grouped := ec.GroupByCategory()
	if len(grouped[ErrorCategoryNetwork]) != 1 {
		t.Errorf("expected 1 network error")
	}
	if len(grouped[ErrorCategoryFileIO]) != 1 {
		t.Errorf("expected 1 file_io error")
	}
	if len(grouped[ErrorCategoryProcessing]) != 1 {
		t.Errorf("expected 1 processing error for plain error fallback")
	}

	summary := ec.FormatSummary(PipelineSummary{
		TotalPosts:     10,
		ProcessedPosts: 8,
		FailedPosts:    2,
		TotalDuration:  2 * time.Second,
	})

	if !strings.Contains(summary, "Pipeline Summary") {
		t.Errorf("summary missing header")
	}
	if !strings.Contains(summary, "network") {
		t.Errorf("summary missing category breakdown")
	}
}
