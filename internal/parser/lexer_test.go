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
					kind:  text,
				},
			},
		},
		{
			input: "test*",
			output: []*token{
				{
					value: "test",
					kind:  text,
				},
				{
					value: "*",
					kind:  wildcard,
				},
			},
		},
		{
			input: "test1*test2",
			output: []*token{
				{
					value: "test1",
					kind:  text,
				},
				{
					value: "*",
					kind:  wildcard,
				},
				{
					value: "test2",
					kind:  text,
				},
			},
		},
		{
			input: "*",
			output: []*token{
				{
					value: "*",
					kind:  wildcard,
				},
			},
		},
		{
			input:  "",
			output: []*token{},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			require := r.New(t)

			output := make([]*token, 0)
			l := newLexer(test.input)

			for l.Next() {
				output = append(output, l.Scan())
			}

			require.Equal(test.output, output)
		})
	}
}
