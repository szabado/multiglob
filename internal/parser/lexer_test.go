package parser

import (
	"testing"

	r "github.com/stretchr/testify/require"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input  string
		output []*token
	}{
		{
			input: "test",
			output: []*token{
				{
					value: "test",
					kind:  LexerText,
				},
			},
		},
		{
			input: "test*",
			output: []*token{
				{
					value: "test",
					kind:  LexerText,
				},
				{
					value: "*",
					kind:  LexerWildcard,
				},
			},
		},
		{
			input: "test1*test2",
			output: []*token{
				{
					value: "test1",
					kind:  LexerText,
				},
				{
					value: "*",
					kind:  LexerWildcard,
				},
				{
					value: "test2",
					kind:  LexerText,
				},
			},
		},
		{
			input: "*",
			output: []*token{
				{
					value: "*",
					kind:  LexerWildcard,
				},
			},
		},
		{
			input:  "",
			output: []*token{},
		},
		{
			input:  "*****",
			output: []*token{
				{
					value: "*",
					kind: LexerWildcard,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			require := r.New(t)

			output := make([]*token, 0)
			l := NewLexer(test.input)

			for l.Next() {
				output = append(output, l.Scan())
			}

			require.Equal(test.output, output)
		})
	}
}
