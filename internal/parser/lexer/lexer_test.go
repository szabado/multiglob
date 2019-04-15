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
					Value: "t",
					Type:  Text,
				},
				{
					Value: "e",
					Type:  Text,
				},
				{
					Value: "s",
					Type:  Text,
				},
				{
					Value: "t",
					Type:  Text,
				},
			},
		},
		{
			input: "test*",
			output: []*Token{
				{
					Value: "t",
					Type:  Text,
				},
				{
					Value: "e",
					Type:  Text,
				},
				{
					Value: "s",
					Type:  Text,
				},
				{
					Value: "t",
					Type:  Text,
				},
				{
					Value: "*",
					Type:  Asterisk,
				},
			},
		},
		{
			input: "test1*test2",
			output: []*Token{
				{
					Value: "t",
					Type:  Text,
				},
				{
					Value: "e",
					Type:  Text,
				},
				{
					Value: "s",
					Type:  Text,
				},
				{
					Value: "t",
					Type:  Text,
				},
				{
					Value: "1",
					Type:  Text,
				},
				{
					Value: "*",
					Type:  Asterisk,
				},
				{
					Value: "t",
					Type:  Text,
				},
				{
					Value: "e",
					Type:  Text,
				},
				{
					Value: "s",
					Type:  Text,
				},
				{
					Value: "t",
					Type:  Text,
				},
				{
					Value: "2",
					Type:  Text,
				},
			},
		},
		{
			input: "*",
			output: []*Token{
				{
					Value: "*",
					Type:  Asterisk,
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
					Type:  Asterisk,
				},
			},
		},
		{
			input: `\`,
			output: []*Token{
				{
					Value: `\`,
					Type:  Backslash,
				},
			},
		},
		{
			input: `-`,
			output: []*Token{
				{
					Value: `-`,
					Type:  Dash,
				},
			},
		},
		{
			input: `][`,
			output: []*Token{
				{
					Value: `]`,
					Type:  Bracket,
				},
				{
					Value: `[`,
					Type:  Bracket,
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
