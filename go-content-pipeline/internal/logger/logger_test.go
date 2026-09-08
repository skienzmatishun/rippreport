package logger

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestLevel_String tests the string representation of log levels
func TestLevel_String(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		expected string
	}{
		{"DEBUG level", DEBUG, "DEBUG"},
		{"INFO level", INFO, "INFO"},
		{"WARN level", WARN, "WARN"},
		{"ERROR level", ERROR, "ERROR"},
		{"Unknown level", Level(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("Level.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestParseLevel tests parsing log levels from strings
func TestParseLevel(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  Level
		expectErr bool
	}{
		{"DEBUG string", "DEBUG", DEBUG, false},
		{"INFO string", "INFO", INFO, false},
		{"WARN string", "WARN", WARN, false},
		{"ERROR string", "ERROR", ERROR, false},
		{"Invalid string", "INVALID", INFO, true},
		{"Empty string", "", INFO, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLevel(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("ParseLevel() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ParseLevel() unexpected error: %v", err)
				}
				if got != tt.expected {
					t.Errorf("ParseLevel() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

// TestNewLogger tests logger creation
func TestNewLogger(t *testing.T) {
	t.Run("Console only logger", func(t *testing.T) {
		logger, err := NewLogger(INFO, "", true)
		if err != nil {
			t.Fatalf("NewLogger() error = %v", err)
		}
		defer logger.Close()

		if logger == nil {
			t.Fatal("NewLogger() returned nil logger")
		}
	})

	t.Run("File logger", func(t *testing.T) {
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "test.log")

		logger, err := NewLogger(INFO, logFile, false)
		if err != nil {
			t.Fatalf("NewLogger() error = %v", err)
		}
		defer logger.Close()

		// Verify file was created
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Errorf("Log file was not created: %s", logFile)
		}
	})

	t.Run("File and console logger", func(t *testing.T) {
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "test.log")

		logger, err := NewLogger(INFO, logFile, true)
		if err != nil {
			t.Fatalf("NewLogger() error = %v", err)
		}
		defer logger.Close()

		// Verify file was created
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Errorf("Log file was not created: %s", logFile)
		}
	})

	t.Run("Invalid file path", func(t *testing.T) {
		_, err := NewLogger(INFO, "/invalid/path/to/log.txt", false)
		if err == nil {
			t.Error("NewLogger() expected error for invalid path but got none")
		}
	})
}

// TestLogLevel_Filtering tests that log level filtering works correctly
func TestLogLevel_Filtering(t *testing.T) {
	tests := []struct {
		name         string
		loggerLevel  Level
		messageLevel Level
		shouldLog    bool
	}{
		{"DEBUG logger logs DEBUG", DEBUG, DEBUG, true},
		{"DEBUG logger logs INFO", DEBUG, INFO, true},
		{"DEBUG logger logs WARN", DEBUG, WARN, true},
		{"DEBUG logger logs ERROR", DEBUG, ERROR, true},
		{"INFO logger skips DEBUG", INFO, DEBUG, false},
		{"INFO logger logs INFO", INFO, INFO, true},
		{"INFO logger logs WARN", INFO, WARN, true},
		{"INFO logger logs ERROR", INFO, ERROR, true},
		{"WARN logger skips DEBUG", WARN, DEBUG, false},
		{"WARN logger skips INFO", WARN, INFO, false},
		{"WARN logger logs WARN", WARN, WARN, true},
		{"WARN logger logs ERROR", WARN, ERROR, true},
		{"ERROR logger skips DEBUG", ERROR, DEBUG, false},
		{"ERROR logger skips INFO", ERROR, INFO, false},
		{"ERROR logger skips WARN", ERROR, WARN, false},
		{"ERROR logger logs ERROR", ERROR, ERROR, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			logFile := filepath.Join(tmpDir, "test.log")

			logger, err := NewLogger(tt.loggerLevel, logFile, false)
			if err != nil {
				t.Fatalf("NewLogger() error = %v", err)
			}
			defer logger.Close()

			// Log a message at the test level
			switch tt.messageLevel {
			case DEBUG:
				logger.Debug("test message")
			case INFO:
				logger.Info("test message")
			case WARN:
				logger.Warn("test message")
			case ERROR:
				logger.Error("test message")
			}

			// Close to flush
			logger.Close()

			// Read log file
			content, err := os.ReadFile(logFile)
			if err != nil {
				t.Fatalf("Failed to read log file: %v", err)
			}

			hasContent := len(content) > 0
			if hasContent != tt.shouldLog {
				t.Errorf("Expected shouldLog=%v but got hasContent=%v", tt.shouldLog, hasContent)
			}

			if tt.shouldLog && !strings.Contains(string(content), "test message") {
				t.Errorf("Log file should contain 'test message' but doesn't. Content: %s", string(content))
			}
		})
	}
}

// TestLogger_FileOutput tests writing to file
func TestLogger_FileOutput(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	logger, err := NewLogger(DEBUG, logFile, false)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	// Write various log levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Close to flush
	logger.Close()

	// Read and verify file content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	expectedMessages := []string{"debug message", "info message", "warn message", "error message"}
	for _, msg := range expectedMessages {
		if !strings.Contains(string(content), msg) {
			t.Errorf("Log file should contain '%s' but doesn't. Content: %s", msg, string(content))
		}
	}

	// Verify log levels are included
	expectedLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	for _, level := range expectedLevels {
		if !strings.Contains(string(content), level) {
			t.Errorf("Log file should contain level '%s' but doesn't", level)
		}
	}
}

// TestLogger_ConsoleOutput tests writing to console
func TestLogger_ConsoleOutput(t *testing.T) {
	// We can't easily capture os.Stdout/Stderr in a test, but we can test
	// that the logger doesn't panic when console is enabled
	logger, err := NewLogger(DEBUG, "", true)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}
	defer logger.Close()

	// These should not panic
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
}

// TestLogger_TimestampFormat tests ISO 8601 timestamp formatting
func TestLogger_TimestampFormat(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	logger, err := NewLogger(INFO, logFile, false)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	logger.Info("test message")
	logger.Close()

	// Read log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// Check for ISO 8601 timestamp format (RFC3339)
	// Example: 2025-01-20T12:34:56Z
	if !strings.Contains(string(content), "T") || !strings.Contains(string(content), "Z") {
		t.Errorf("Log should contain ISO 8601 timestamp but got: %s", string(content))
	}

	// Verify it starts with a timestamp
	line := strings.TrimSpace(string(content))
	if !strings.HasPrefix(line, "[2") {
		t.Errorf("Log line should start with [2XXX... timestamp but got: %s", line)
	}
}

// TestLogger_ContextFields tests adding context fields
func TestLogger_ContextFields(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	logger, err := NewLogger(INFO, logFile, false)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}
	defer logger.Close()

	// Log with fields
	logger.Info("test message",
		String("post_slug", "example-post"),
		String("operation", "generate_embedding"),
		Int("count", 42),
	)

	logger.Close()

	// Read and verify
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// Verify fields are present
	contentStr := string(content)
	if !strings.Contains(contentStr, "post_slug=example-post") {
		t.Errorf("Log should contain post_slug field but got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "operation=generate_embedding") {
		t.Errorf("Log should contain operation field but got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "count=42") {
		t.Errorf("Log should contain count field but got: %s", contentStr)
	}
}

// TestLogger_With tests creating loggers with persistent context fields
func TestLogger_With(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	logger, err := NewLogger(INFO, logFile, false)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}
	defer logger.Close()

	// Create a logger with context
	contextLogger := logger.With(
		String("post_slug", "example-post"),
		String("operation", "ranking"),
	)

	// Log with the context logger
	contextLogger.Info("first message")
	contextLogger.Info("second message", Int("score", 95))

	logger.Close()

	// Read and verify
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 2 {
		t.Fatalf("Expected 2 log lines but got %d", len(lines))
	}

	// Both lines should have the context fields
	for i, line := range lines {
		if !strings.Contains(line, "post_slug=example-post") {
			t.Errorf("Line %d should contain post_slug field: %s", i+1, line)
		}
		if !strings.Contains(line, "operation=ranking") {
			t.Errorf("Line %d should contain operation field: %s", i+1, line)
		}
	}

	// Second line should also have the score
	if !strings.Contains(lines[1], "score=95") {
		t.Errorf("Line 2 should contain score field: %s", lines[1])
	}
}

// TestLogger_SetLevel tests changing log level dynamically
func TestLogger_SetLevel(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	logger, err := NewLogger(INFO, logFile, false)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}
	defer logger.Close()

	// This should be logged (level is INFO)
	logger.Info("info message 1")

	// This should not be logged (DEBUG < INFO)
	logger.Debug("debug message 1")

	// Change level to DEBUG
	logger.SetLevel(DEBUG)

	// Now DEBUG should be logged
	logger.Debug("debug message 2")

	// Change level to ERROR
	logger.SetLevel(ERROR)

	// Now INFO should not be logged
	logger.Info("info message 2")

	logger.Close()

	// Read and verify
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)

	// Should contain
	if !strings.Contains(contentStr, "info message 1") {
		t.Error("Should contain 'info message 1'")
	}
	if !strings.Contains(contentStr, "debug message 2") {
		t.Error("Should contain 'debug message 2'")
	}

	// Should NOT contain
	if strings.Contains(contentStr, "debug message 1") {
		t.Error("Should NOT contain 'debug message 1'")
	}
	if strings.Contains(contentStr, "info message 2") {
		t.Error("Should NOT contain 'info message 2'")
	}
}

// TestField_Constructors tests the field constructor functions
func TestField_Constructors(t *testing.T) {
	tests := []struct {
		name     string
		field    Field
		wantKey  string
		wantType string
	}{
		{"NewField", NewField("custom", "value"), "custom", "string"},
		{"String field", String("name", "test"), "name", "string"},
		{"Int field", Int("count", 42), "count", "int"},
		{"Int64 field", Int64("big", int64(1234567890)), "big", "int64"},
		{"Float32 field", Float32("score", float32(3.14)), "score", "float32"},
		{"Float64 field", Float64("precise", 3.14159265), "precise", "float64"},
		{"Bool field", Bool("active", true), "active", "bool"},
		{"Duration field", Duration("elapsed", time.Second), "elapsed", "time.Duration"},
		{"Error field", Error(errors.New("test error")), "error", "string"},
		{"Nil error field", Error(nil), "error", "nil"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field.Key != tt.wantKey {
				t.Errorf("Field.Key = %v, want %v", tt.field.Key, tt.wantKey)
			}

			// Check value is not nil (except for nil error)
			if tt.wantType != "nil" && tt.field.Value == nil {
				t.Errorf("Field.Value should not be nil")
			}
		})
	}
}

// TestLogger_AppendMode tests that logger appends to existing file
func TestLogger_AppendMode(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	// First logger writes initial content
	logger1, err := NewLogger(INFO, logFile, false)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}
	logger1.Info("first message")
	logger1.Close()

	// Second logger should append
	logger2, err := NewLogger(INFO, logFile, false)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}
	logger2.Info("second message")
	logger2.Close()

	// Read and verify both messages are present
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "first message") {
		t.Error("Should contain 'first message'")
	}
	if !strings.Contains(contentStr, "second message") {
		t.Error("Should contain 'second message'")
	}
}

// TestDefaultLogger tests the default global logger functions
func TestDefaultLogger(t *testing.T) {
	// Create a custom logger and set as default
	var buf bytes.Buffer
	customLogger := &testLogger{buf: &buf, level: DEBUG}
	SetDefault(customLogger)

	// Test that global functions use the default logger
	Debug("debug test")
	Info("info test")
	Warn("warn test")
	Errorf("error test")

	// Test GetDefault
	retrieved := GetDefault()
	if retrieved != customLogger {
		t.Error("GetDefault should return the same logger that was set")
	}

	// Test global With function
	contextLogger := With(String("context", "test"))
	contextLogger.Info("with context")

	output := buf.String()
	if !strings.Contains(output, "debug test") {
		t.Error("Should contain 'debug test'")
	}
	if !strings.Contains(output, "info test") {
		t.Error("Should contain 'info test'")
	}
	if !strings.Contains(output, "warn test") {
		t.Error("Should contain 'warn test'")
	}
	if !strings.Contains(output, "error test") {
		t.Error("Should contain 'error test'")
	}
	if !strings.Contains(output, "with context") {
		t.Error("Should contain 'with context'")
	}
}

// TestLogger_ConcurrentWrites tests thread safety
func TestLogger_ConcurrentWrites(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	logger, err := NewLogger(DEBUG, logFile, false)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}
	defer logger.Close()

	// Write concurrently from multiple goroutines
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				logger.Info("concurrent message", Int("goroutine", id), Int("iteration", j))
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < 10; i++ {
		<-done
	}

	logger.Close()

	// Verify file was written and contains expected number of lines
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	expectedLines := 1000 // 10 goroutines * 100 iterations
	if len(lines) != expectedLines {
		t.Errorf("Expected %d log lines but got %d", expectedLines, len(lines))
	}
}

// testLogger is a simple test implementation of Logger for testing global functions
type testLogger struct {
	buf    *bytes.Buffer
	level  Level
	fields []Field
}

func (l *testLogger) Debug(msg string, fields ...Field) {
	if DEBUG >= l.level {
		l.buf.WriteString("DEBUG: " + msg + "\n")
	}
}

func (l *testLogger) Info(msg string, fields ...Field) {
	if INFO >= l.level {
		l.buf.WriteString("INFO: " + msg + "\n")
	}
}

func (l *testLogger) Warn(msg string, fields ...Field) {
	if WARN >= l.level {
		l.buf.WriteString("WARN: " + msg + "\n")
	}
}

func (l *testLogger) Error(msg string, fields ...Field) {
	if ERROR >= l.level {
		l.buf.WriteString("ERROR: " + msg + "\n")
	}
}

func (l *testLogger) With(fields ...Field) Logger {
	newFields := make([]Field, len(l.fields)+len(fields))
	copy(newFields, l.fields)
	copy(newFields[len(l.fields):], fields)
	return &testLogger{buf: l.buf, level: l.level, fields: newFields}
}

func (l *testLogger) SetLevel(level Level) {
	l.level = level
}

func (l *testLogger) Close() error {
	return nil
}
