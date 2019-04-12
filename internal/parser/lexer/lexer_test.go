package lexer

import (
	"testing"

	r "github.com/stretchr/testify/require"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input  string
		output []*Token
	}{
		{
			input: "test",
			output: []*Token{
				{
					Value: "test",
					Kind:  Text,
				},
			},
		},
		{
			input: "test*",
			output: []*Token{
				{
					Value: "test",
					Kind:  Text,
				},
				{
					Value: "*",
					Kind:  Wildcard,
				},
			},
		},
		{
			input: "test1*test2",
			output: []*Token{
				{
					Value: "test1",
					Kind:  Text,
				},
				{
					Value: "*",
					Kind:  Wildcard,
				},
				{
					Value: "test2",
					Kind:  Text,
				},
			},
		},
		{
			input: "*",
			output: []*Token{
				{
					Value: "*",
					Kind:  Wildcard,
				},
			},
		},
		{
			input:  "",
			output: []*Token{},
		},
		{
			input: "*****",
			output: []*Token{
				{
					Value: "*",
					Kind:  Wildcard,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			require := r.New(t)

			output := make([]*Token, 0)
			l := New(test.input)

			for l.Next() {
				output = append(output, l.Scan())
			}

			require.Equal(test.output, output)
		})
	}
}