package youtube

import (
	"testing"
)

func TestExtractVideoID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "bare video ID",
			input: "L6dyOWiaei4",
			want:  "L6dyOWiaei4",
		},
		{
			name:  "standard watch URL",
			input: "https://www.youtube.com/watch?v=L6dyOWiaei4",
			want:  "L6dyOWiaei4",
		},
		{
			name:  "watch URL with extra params",
			input: "https://www.youtube.com/watch?v=L6dyOWiaei4&t=120",
			want:  "L6dyOWiaei4",
		},
		{
			name:  "short URL",
			input: "https://youtu.be/L6dyOWiaei4",
			want:  "L6dyOWiaei4",
		},
		{
			name:  "embed URL",
			input: "https://www.youtube.com/embed/L6dyOWiaei4",
			want:  "L6dyOWiaei4",
		},
		{
			name:  "v URL",
			input: "https://www.youtube.com/v/L6dyOWiaei4",
			want:  "L6dyOWiaei4",
		},
		{
			name:  "URL with http",
			input: "http://www.youtube.com/watch?v=L6dyOWiaei4",
			want:  "L6dyOWiaei4",
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid URL",
			input:   "https://example.com/video",
			wantErr: true,
		},
		{
			name:    "youtube URL without video ID",
			input:   "https://www.youtube.com/channel/UCtest",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractVideoID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractVideoID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractVideoID() = %q, want %q", got, tt.want)
			}
		})
	}
}

