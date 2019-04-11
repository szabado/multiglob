package multiglob

import (
	"sort"
	"testing"

	"github.com/szabado/multiglob/internal/parser"

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
				"*test",
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
				"test*",
				"*hit",
			},
			output: true,
		},
		{
			input: "foo",
			patterns: []string{
				"test",
			},
			output: false,
		},
		{
			input: "foo",
			patterns: []string{
				"test*",
			},
			output: false,
		},
		{
			input: "foo",
			patterns: []string{
				"*test",
			},
			output: false,
		},
		{
			input: "foo",
			patterns: []string{
				"tes*t",
			},
			output: false,
		},
		{
			input: "test.hit",
			patterns: []string{
				"test*hit.hit",
			},
			output: false,
		},
		{
			input: "test.hit",
			patterns: []string{
				"*test",
				"hit*",
			},
			output: false,
		},
		{
			input: "test.hit",
			patterns: []string{
				"*test",
				"hit*",
			},
			output: false,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			require := r.New(t)

			b := New()
			for _, pattern := range test.patterns {
				b.MustAddPattern(pattern, pattern)
			}

			mg := b.MustCompile()

			require.Equal(test.output, mg.Match(test.input))
		})
	}
}

func TestAddPattern(t *testing.T) {
	require := r.New(t)

	b := New()
	b.MustAddPattern("test", "pattern")

	require.Equal(1, len(b.patterns))
	require.Equal(parser.Parse("test", "pattern"), b.patterns["test"])
}

func TestFindAllPatterns(t *testing.T) {
	tests := []struct {
		patterns map[string]string
		input    string
		output   []string
	}{
		{
			input: "test",
			patterns: map[string]string{
				"a": "test",
			},
			output: []string{
				"a",
			},
		},
		{
			input: "foobar",
			patterns: map[string]string{
				"a": "foo*",
				"b": "*bar",
			},
			output: []string{
				"a",
				"b",
			},
		},
		{
			input: "whoops",
			patterns: map[string]string{
				"a": "foobar",
			},
			output: nil,
		},
		{
			input: "foo",
			patterns: map[string]string{
				"a": "foo",
				"b": "foo*",
			},
			output: []string{
				"a",
				"b",
			},
		},
		{
			input: "bar",
			patterns: map[string]string{
				"a": "*",
				"b": "*bar",
			},
			output: []string{
				"a",
				"b",
			},
		},
		{
			input: "hit hit hit",
			patterns: map[string]string{
				"a": "*hit*",
			},
			output: []string{
				"a",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			require := r.New(t)

			b := New()
			for name, pattern := range test.patterns {
				b.MustAddPattern(name, pattern)
			}

			mg := b.MustCompile()

			sort.Strings(test.output)
			output := mg.FindAllPatterns(test.input)

			sort.Strings(output)

			require.Equal(test.output, output)
		})
	}
}
