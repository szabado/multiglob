package multiglob

import (
	"fmt"
	"sort"
	"testing"

	r "github.com/stretchr/testify/require"

	"github.com/szabado/multiglob/internal/parser"
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
		{
			input: "ba",
			patterns: []string{
				"[ab]+",
			},
			output: true,
		},
		{
			input: "bb",
			patterns: []string{
				"[a]+bb",
			},
			output: false,
		},
		{
			input: "bb",
			patterns: []string{
				"[a]bb",
			},
			output: false,
		},
		{
			input: "aa",
			patterns: []string{
				"[a]+a",
			},
			output: true,
		},
		{
			input: "this.is.not.a.dance",
			patterns: []string{
				"this[^.]dance",
			},
			output: false,
		},
		{
			input: "this.is.not.a.dance",
			patterns: []string{
				"this[isnota.]dance",
			},
			output: false,
		},
		{
			input: "this.is.not.a.dance",
			patterns: []string{
				"this[isnota.]+dance",
			},
			output: true,
		},
		{
			input: "this.is.not.a.dance",
			patterns: []string{
				"this[*]+dance",
			},
			output: false,
		},
		{
			input: "",
			patterns: []string{
				"[abc]",
			},
			output: false,
		},
		{
			input: "abcdef",
			patterns: []string{
				"*[def]+",
			},
			output: true,
		},
		{
			input: "abc",
			patterns: []string{
				"*[^b]*",
			},
			output: true,
		},
		{
			input: "abc",
			patterns: []string{
				"[^b]+",
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
	ast, err := parser.Parse("test", "pattern")
	require.NoError(err)
	require.Equal(ast, b.patterns["test"])
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
			output: []string{},
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
		{
			input: "pen pineapple apple pen",
			patterns: map[string]string{
				"a": "*apple*",
				"b": "*pen*",
			},
			output: []string{
				"a",
				"b",
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
			output:  nil,
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

			ast, err := parser.Parse(test.pattern, test.pattern)
			require.NoError(err)

			output, err := extractGlobs(test.input, ast)
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
			output:  nil,
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
		{
			input:   "abcba",
			pattern: "[abc]+",
			output: []string{
				"abcba",
			},
			matched: true,
		},
		{
			input:   "abcdef",
			pattern: "*[d-f]+",
			output: []string{
				"abc",
				"def",
			},
			matched: true,
		},
		{
			input:   "abc",
			pattern: "[abc][abc][abc]",
			output: []string{
				"a",
				"b",
				"c",
			},
			matched: true,
		},
		{
			input:   "abc",
			pattern: "*[^b]*",
			output: []string{
				"ab",
				"c",
				"",
			},
			matched: true,
		},
		{
			input:   "abc",
			pattern: "[^b]+",
			output:  nil,
			matched: false,
		},
		{
			input:   "abc",
			pattern: "a[^b]c",
			output:  nil,
			matched: false,
		},
		{
			input:   "this.is.a.test",
			pattern: "this.[^.]+*test",
			output: []string{
				"is",
				".a.",
			},
			matched: true,
		},
		{
			input:   "this.is.a.test",
			pattern: "this.[^.]*test",
			output: []string{
				"i",
				"s.a.",
			},
			matched: true,
		},
	}

	for i, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			require := r.New(t)

			b := New()
			b.MustAddPattern(fmt.Sprint(i), test.pattern)

			mg := b.MustCompile()

			name, globs, matched := mg.FindGlobs(test.input)
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
				"a": nil,
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
				"b": nil,
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
