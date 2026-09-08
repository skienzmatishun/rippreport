package testutil

import (
	"fmt"
	"testing"
)

// TableTest represents a single test case in a table-driven test.
type TableTest struct {
	Name string
	Run  func(t *testing.T)
}

// RunTableTests executes a slice of table-driven tests.
func RunTableTests(t *testing.T, tests []TableTest) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.Name, tt.Run)
	}
}

// StringTest represents a test case with string input and expected output.
type StringTest struct {
	Name     string
	Input    string
	Expected string
	ShouldError bool
}

// RunStringTests executes string-based table tests.
func RunStringTests(t *testing.T, tests []StringTest, fn func(string) (string, error)) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			result, err := fn(tt.Input)
			
			if tt.ShouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			if result != tt.Expected {
				t.Errorf("got %q, want %q", result, tt.Expected)
			}
		})
	}
}

// IntTest represents a test case with integer input and expected output.
type IntTest struct {
	Name     string
	Input    int
	Expected int
	ShouldError bool
}

// RunIntTests executes integer-based table tests.
func RunIntTests(t *testing.T, tests []IntTest, fn func(int) (int, error)) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			result, err := fn(tt.Input)
			
			if tt.ShouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			if result != tt.Expected {
				t.Errorf("got %d, want %d", result, tt.Expected)
			}
		})
	}
}

// FloatTest represents a test case with float input and expected output.
type FloatTest struct {
	Name     string
	Input    float64
	Expected float64
	Epsilon  float64
	ShouldError bool
}

// RunFloatTests executes float-based table tests with epsilon comparison.
func RunFloatTests(t *testing.T, tests []FloatTest, fn func(float64) (float64, error)) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			result, err := fn(tt.Input)
			
			if tt.ShouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			epsilon := tt.Epsilon
			if epsilon == 0 {
				epsilon = 0.0001
			}
			
			diff := result - tt.Expected
			if diff < 0 {
				diff = -diff
			}
			
			if diff > epsilon {
				t.Errorf("got %.10f, want %.10f (diff %.10f > epsilon %.10f)", 
					result, tt.Expected, diff, epsilon)
			}
		})
	}
}

// VectorTest represents a test case with vector input and expected output.
type VectorTest struct {
	Name     string
	InputA   []float32
	InputB   []float32
	Expected float32
	Epsilon  float32
	ShouldError bool
}

// RunVectorTests executes vector-based table tests (e.g., for similarity calculations).
func RunVectorTests(t *testing.T, tests []VectorTest, fn func([]float32, []float32) (float32, error)) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			result, err := fn(tt.InputA, tt.InputB)
			
			if tt.ShouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			epsilon := tt.Epsilon
			if epsilon == 0 {
				epsilon = 0.0001
			}
			
			diff := result - tt.Expected
			if diff < 0 {
				diff = -diff
			}
			
			if diff > epsilon {
				t.Errorf("got %.10f, want %.10f (diff %.10f > epsilon %.10f)", 
					result, tt.Expected, diff, epsilon)
			}
		})
	}
}

// ParseTest represents a test case for parsing operations.
type ParseTest struct {
	Name        string
	Input       string
	Expected    interface{}
	ShouldError bool
	ErrorMsg    string
}

// ValidationTest represents a test case for validation operations.
type ValidationTest struct {
	Name     string
	Input    interface{}
	IsValid  bool
	ErrorMsg string
}

// RunValidationTests executes validation-based table tests.
func RunValidationTests(t *testing.T, tests []ValidationTest, fn func(interface{}) error) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := fn(tt.Input)
			
			if tt.IsValid {
				if err != nil {
					t.Errorf("expected valid input but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected invalid input but got no error")
				} else if tt.ErrorMsg != "" && !contains(err.Error(), tt.ErrorMsg) {
					t.Errorf("expected error message to contain %q but got %q", tt.ErrorMsg, err.Error())
				}
			}
		})
	}
}

// TransformTest represents a test case for transformation operations.
type TransformTest struct {
	Name     string
	Input    interface{}
	Expected interface{}
	ShouldError bool
}

// CompareTest represents a test case for comparison operations.
type CompareTest struct {
	Name     string
	A        interface{}
	B        interface{}
	Expected int // -1, 0, or 1
}

// RoundTripTest represents a test case for round-trip conversions.
type RoundTripTest struct {
	Name  string
	Input interface{}
}

// RunRoundTripTests tests that serialize -> deserialize -> serialize produces same output.
// This validates requirement 23.8 (round-trip property for YAML and hooks).
func RunRoundTripTests(t *testing.T, tests []RoundTripTest, 
	serialize func(interface{}) (string, error),
	deserialize func(string) (interface{}, error)) {
	
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// First serialization
			serialized1, err := serialize(tt.Input)
			if err != nil {
				t.Fatalf("first serialization failed: %v", err)
			}
			
			// Deserialization
			deserialized, err := deserialize(serialized1)
			if err != nil {
				t.Fatalf("deserialization failed: %v", err)
			}
			
			// Second serialization
			serialized2, err := serialize(deserialized)
			if err != nil {
				t.Fatalf("second serialization failed: %v", err)
			}
			
			// Compare serialized outputs
			if serialized1 != serialized2 {
				t.Errorf("round-trip failed:\nfirst:  %s\nsecond: %s", serialized1, serialized2)
			}
		})
	}
}

// PropertyTest represents a property-based test case.
type PropertyTest struct {
	Name     string
	Property func(t *testing.T) bool
	Samples  int
}

// RunPropertyTests executes property-based tests.
// Each property is tested with multiple samples.
func RunPropertyTests(t *testing.T, tests []PropertyTest) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			samples := tt.Samples
			if samples == 0 {
				samples = 100
			}
			
			failures := 0
			for i := 0; i < samples; i++ {
				if !tt.Property(t) {
					failures++
				}
			}
			
			if failures > 0 {
				t.Errorf("property failed %d out of %d times", failures, samples)
			}
		})
	}
}

// BenchmarkTest represents a benchmark test case.
type BenchmarkTest struct {
	Name string
	Run  func(b *testing.B)
}

// RunBenchmarks executes a slice of benchmark tests.
func RunBenchmarks(b *testing.B, tests []BenchmarkTest) {
	for _, bt := range tests {
		b.Run(bt.Name, bt.Run)
	}
}

// ErrorTest represents a test case focused on error handling.
type ErrorTest struct {
	Name          string
	Setup         func() error
	ExpectedError string
	ShouldRetry   bool
}

// RunErrorTests executes error handling tests.
func RunErrorTests(t *testing.T, tests []ErrorTest) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.Setup()
			
			if err == nil {
				t.Errorf("expected error but got none")
				return
			}
			
			if tt.ExpectedError != "" && !contains(err.Error(), tt.ExpectedError) {
				t.Errorf("expected error to contain %q but got %q", tt.ExpectedError, err.Error())
			}
		})
	}
}

// ConcurrencyTest represents a test case for concurrent operations.
type ConcurrencyTest struct {
	Name       string
	Goroutines int
	Operations int
	Run        func(goroutineID, operationID int) error
}

// RunConcurrencyTests executes concurrent operation tests.
func RunConcurrencyTests(t *testing.T, tests []ConcurrencyTest) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			errors := make(chan error, tt.Goroutines*tt.Operations)
			done := make(chan bool)
			
			for g := 0; g < tt.Goroutines; g++ {
				go func(goroutineID int) {
					for op := 0; op < tt.Operations; op++ {
						if err := tt.Run(goroutineID, op); err != nil {
							errors <- fmt.Errorf("goroutine %d operation %d: %v", goroutineID, op, err)
						}
					}
					done <- true
				}(g)
			}
			
			// Wait for all goroutines
			for g := 0; g < tt.Goroutines; g++ {
				<-done
			}
			close(errors)
			
			// Check for errors
			errorCount := 0
			for err := range errors {
				t.Error(err)
				errorCount++
			}
			
			if errorCount > 0 {
				t.Errorf("%d concurrent operations failed", errorCount)
			}
		})
	}
}
