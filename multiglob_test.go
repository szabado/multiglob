package multiglob

import (
	"testing"

	r "github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		patterns []string
		input    string
		output   bool
	}{
		{
			input: "test",
			patterns: []string{
				"test",
			},
			output: true,
		},
		{
			input: "test",
			patterns: []string{
				"*",
			},
			output: true,
		},
		{
			input: "test",
			patterns: []string{
				"test*",
			},
			output: true,
		},
		{
			input: "test",
			patterns: []string{
				"tes*t",
			},
			output: true,
		},
		{
			input: "test.hit.hit.hit",
			patterns: []string{
				"test*hit",
			},
			output: true,
		},
		{
			input: "test.hit",
			patterns: []string{
				"test*hit.hit",
			},
			output: false,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			require := r.New(t)

			b := New()
			for _, pattern := range test.patterns {
				b.AddPattern(pattern, pattern)
			}

			mg := b.Build()

			require.Equal(test.output, mg.Match(test.input))
		})
	}
}
