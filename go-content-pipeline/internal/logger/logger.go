package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Level represents the severity of a log message
type Level int

const (
	// DEBUG level for detailed debugging information
	DEBUG Level = iota
	// INFO level for general informational messages
	INFO
	// WARN level for warning messages
	WARN
	// ERROR level for error messages
	ERROR
)

// String returns the string representation of a log level
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel converts a string to a Level
func ParseLevel(s string) (Level, error) {
	switch s {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	default:
		return INFO, fmt.Errorf("invalid log level: %s", s)
	}
}

// Logger provides structured logging with multiple output destinations
type Logger interface {
	// Debug logs a DEBUG level message
	Debug(msg string, fields ...Field)

	// Info logs an INFO level message
	Info(msg string, fields ...Field)

	// Warn logs a WARN level message
	Warn(msg string, fields ...Field)

	// Error logs an ERROR level message
	Error(msg string, fields ...Field)

	// With returns a new logger with additional context fields
	With(fields ...Field) Logger

	// SetLevel changes the minimum log level
	SetLevel(level Level)

	// Close closes the log file (if any)
	Close() error
}

// Field represents a structured logging field with key-value pairs
type Field struct {
	Key   string
	Value interface{}
}

// NewField creates a new logging field
func NewField(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// Convenience functions for creating common fields
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Float32(key string, value float32) Field {
	return Field{Key: key, Value: value}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	if err == nil {
		return Field{Key: "error", Value: nil}
	}
	return Field{Key: "error", Value: err.Error()}
}

// standardLogger implements the Logger interface
type standardLogger struct {
	level      Level
	fileWriter io.WriteCloser
	consoleOut io.Writer
	consoleErr io.Writer
	mu         sync.Mutex
	fields     []Field
}

// NewLogger creates a new logger instance
// If filePath is empty, logs will only be written to console
// If consoleEnabled is false, logs will only be written to file
func NewLogger(level Level, filePath string, consoleEnabled bool) (Logger, error) {
	var fileWriter io.WriteCloser
	var err error

	// Open log file if path is provided
	if filePath != "" {
		fileWriter, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file %s: %w", filePath, err)
		}
	}

	consoleOut := io.Writer(nil)
	consoleErr := io.Writer(nil)
	if consoleEnabled {
		consoleOut = os.Stdout
		consoleErr = os.Stderr
	}

	return &standardLogger{
		level:      level,
		fileWriter: fileWriter,
		consoleOut: consoleOut,
		consoleErr: consoleErr,
		fields:     make([]Field, 0),
	}, nil
}

// Debug logs a DEBUG level message
func (l *standardLogger) Debug(msg string, fields ...Field) {
	l.log(DEBUG, msg, fields...)
}

// Info logs an INFO level message
func (l *standardLogger) Info(msg string, fields ...Field) {
	l.log(INFO, msg, fields...)
}

// Warn logs a WARN level message
func (l *standardLogger) Warn(msg string, fields ...Field) {
	l.log(WARN, msg, fields...)
}

// Error logs an ERROR level message
func (l *standardLogger) Error(msg string, fields ...Field) {
	l.log(ERROR, msg, fields...)
}

// With returns a new logger with additional context fields
func (l *standardLogger) With(fields ...Field) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Create a new logger with combined fields
	newFields := make([]Field, 0, len(l.fields)+len(fields))
	newFields = append(newFields, l.fields...)
	newFields = append(newFields, fields...)

	return &standardLogger{
		level:      l.level,
		fileWriter: l.fileWriter,
		consoleOut: l.consoleOut,
		consoleErr: l.consoleErr,
		fields:     newFields,
	}
}

// SetLevel changes the minimum log level
func (l *standardLogger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// Close closes the log file
func (l *standardLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.fileWriter != nil {
		return l.fileWriter.Close()
	}
	return nil
}

// log writes a log message to configured outputs
func (l *standardLogger) log(level Level, msg string, fields ...Field) {
	// Check if we should log this level
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Format timestamp in ISO 8601 format
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Combine context fields with message fields
	allFields := make([]Field, 0, len(l.fields)+len(fields))
	allFields = append(allFields, l.fields...)
	allFields = append(allFields, fields...)

	// Format the log message
	logLine := l.formatLogLine(timestamp, level, msg, allFields)

	// Write to file if configured
	if l.fileWriter != nil {
		l.fileWriter.Write([]byte(logLine + "\n"))
	}

	// Write to console if configured
	if l.consoleOut != nil || l.consoleErr != nil {
		// ERROR level goes to stderr, others to stdout
		if level == ERROR && l.consoleErr != nil {
			l.consoleErr.Write([]byte(logLine + "\n"))
		} else if l.consoleOut != nil {
			l.consoleOut.Write([]byte(logLine + "\n"))
		}
	}
}

// formatLogLine creates a formatted log line with timestamp, level, message, and fields
func (l *standardLogger) formatLogLine(timestamp string, level Level, msg string, fields []Field) string {
	// Start with timestamp, level, and message
	line := fmt.Sprintf("[%s] %s: %s", timestamp, level.String(), msg)

	// Add fields if present
	if len(fields) > 0 {
		line += " |"
		for _, field := range fields {
			line += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		}
	}

	return line
}

// Global default logger (can be replaced)
var defaultLogger Logger

// Initialize default logger to console-only with INFO level
func init() {
	defaultLogger, _ = NewLogger(INFO, "", true)
}

// SetDefault sets the default global logger
func SetDefault(logger Logger) {
	defaultLogger = logger
}

// GetDefault returns the default global logger
func GetDefault() Logger {
	return defaultLogger
}

// Convenience functions that use the default logger
func Debug(msg string, fields ...Field) {
	defaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	defaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	defaultLogger.Warn(msg, fields...)
}

func Errorf(msg string, fields ...Field) {
	defaultLogger.Error(msg, fields...)
}

func With(fields ...Field) Logger {
	return defaultLogger.With(fields...)
}
