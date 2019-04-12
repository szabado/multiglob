package multiglob

import (
	"fmt"
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

func requireOneOf(require *r.Assertions, options []string, result string) {
	for _, option := range options {
		if option == result {
			return
		}
	}
	require.FailNowf("FAIL", "No options matched the result. Options: %#v, Result: %#v", options, result)
}

func TestFindPattern(t *testing.T) {
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
			output, ok := mg.FindPattern(test.input)
			if len(test.output) == 0 {
				require.False(ok)
			} else {
				requireOneOf(require, test.output, output)
			}
		})
	}
}

func TestExtractGlobs(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		output  []string
		err     bool
	}{
		{
			input:   "test",
			pattern: "test",
			output:  []string{},
		},
		{
			input:   "foo",
			pattern: "f*",
			output: []string{
				"oo",
			},
		},
		{
			input:   "foobar",
			pattern: "*f*b*",
			output: []string{
				"",
				"oo",
				"ar",
			},
		},
		{
			input:   "pen pineapple apple pen",
			pattern: "*apple*",
			output: []string{
				"pen pineapple ",
				" pen",
			},
		},
		{
			input:   "pen",
			pattern: "foo",
			output:  nil,
			err:     true,
		},
		{
			input:   "pineapple",
			pattern: "*foo",
			output:  nil,
			err:     true,
		},
		{
			input:   "apple",
			pattern: "foo*",
			output:  nil,
			err:     true,
		},
		{
			input:   "aba",
			pattern: "a*a*a",
			output:  nil,
			err:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			require := r.New(t)

			output, err := extractGlobs(test.input, parser.Parse(test.pattern, test.pattern))
			if test.err {
				require.Error(err)
			} else {
				require.NoError(err)
			}

			require.Equal(test.output, output)
		})
	}
}

func TestFindGlobs(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		output  []string
		matched bool
	}{
		{
			input:   "test",
			pattern: "test",
			output:  []string{},
			matched: true,
		},
		{
			input:   "foo",
			pattern: "f*",
			output: []string{
				"oo",
			},
			matched: true,
		},
		{
			input:   "foobar",
			pattern: "*f*b*",
			output: []string{
				"",
				"oo",
				"ar",
			},
			matched: true,
		},
		{
			input:   "pen pineapple apple pen",
			pattern: "*apple*",
			output: []string{
				"pen pineapple ",
				" pen",
			},
			matched: true,
		},
		{
			input:   "pen",
			pattern: "foo",
			output:  nil,
			matched: false,
		},
		{
			input:   "pineapple",
			pattern: "*foo",
			output:  nil,
			matched: false,
		},
		{
			input:   "apple",
			pattern: "foo*",
			output:  nil,
			matched: false,
		},
		{
			input:   "aba",
			pattern: "a*a*a",
			output:  nil,
			matched: false,
		},
	}

	for i, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			require := r.New(t)

			b := New()
			b.MustAddPattern(fmt.Sprint(i), test.pattern)

			mg := b.MustCompile()

			name, matched, globs := mg.FindGlobs(test.input)
			if test.matched {
				require.True(matched)
				require.Equal(fmt.Sprint(i), name)
			} else {
				require.False(matched)
				require.Equal("", name)
			}
			require.Equal(test.output, globs)

			globs, err := mg.FindGlobsForPattern(test.input, name)
			if test.matched {
				require.NoError(err)
			} else {
				require.Error(err)
			}

			require.Equal(test.output, globs)
		})
	}
}

func TestFindAllGlobs(t *testing.T) {
	tests := []struct {
		patterns map[string]string
		input    string
		output   map[string][]string
	}{
		{
			input: "test",
			patterns: map[string]string{
				"a": "test",
			},
			output: map[string][]string{
				"a": {},
			},
		},
		{
			input: "pen pineapple apple pen",
			patterns: map[string]string{
				"a": "*apple*",
				"b": "*pen*",
			},
			output: map[string][]string{
				"a": {
					"pen pineapple ",
					" pen",
				},
				"b": {
					"pen pineapple apple ",
					"",
				},
			},
		},
		{
			input: "foobar",
			patterns: map[string]string{
				"a": "foo*",
				"b": "*apple*",
			},
			output: map[string][]string{
				"a": {
					"bar",
				},
			},
		},
		{
			input: "foo",
			patterns: map[string]string{
				"a": "*o",
				"b": "foo",
			},
			output: map[string][]string{
				"a": {
					"fo",
				},
				"b": {},
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

			output := mg.FindAllGlobs(test.input)
			require.Equal(test.output, output)
		})
	}
}
