package tui

import (
	"reflect"
	"testing"
)

func TestSplitArgs(t *testing.T) {
	for _, tt := range []struct {
		input string
		want  []string
		bad   bool
	}{
		{``, nil, false},
		{`--posts foo,bar --config "path with spaces.yaml"`, []string{"--posts", "foo,bar", "--config", "path with spaces.yaml"}, false},
		{`--title 'A new post' --description ""`, []string{"--title", "A new post", "--description", ""}, false},
		{`a\ b '$HOME' $(date)`, []string{"a b", "$HOME", "$(date)"}, false},
		{`"unfinished`, nil, true},
		{`trailing\`, nil, true},
	} {
		t.Run(tt.input, func(t *testing.T) {
			got, err := splitArgs(tt.input)
			if (err != nil) != tt.bad {
				t.Fatalf("error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
