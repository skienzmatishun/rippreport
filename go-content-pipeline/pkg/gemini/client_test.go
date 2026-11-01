package gemini

import (
	"testing"
)

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "Valid JSON array",
			input: `[
				{"slug": "test1", "text": "Test text 1"},
				{"slug": "test2", "text": "Test text 2"}
			]`,
			wantErr: false,
		},
		{
			name: "JSON with extra text before",
			input: `Here is the refined output:
			[
				{"slug": "test1", "text": "Test text 1"}
			]`,
			wantErr: false,
		},
		{
			name: "JSON with extra text after",
			input: `[
				{"slug": "test1", "text": "Test text 1"}
			]
			Hope this helps!`,
			wantErr: false,
		},
		{
			name:    "No JSON array",
			input:   "Just some text without JSON",
			wantErr: true,
		},
		{
			name:    "Malformed JSON",
			input:   `[{"slug": "test", "text": "unclosed`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(result) == 0 {
				t.Errorf("extractJSON() returned empty result for valid input")
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "Valid config",
			cfg: Config{
				APIKey: "test-key",
			},
			wantErr: false,
		},
		{
			name:    "Missing API key",
			cfg:     Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
