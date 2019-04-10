package parser

import (
	"testing"

	r "github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input  string
		output *Node
	}{
		{
			input: "test",
			output: &Node{
				Children: nil,
				Value:    "test",
				Type:     TypeText,
			},
		},
		{
			input: "test*",
			output: &Node{
				Children: []*Node{
					{
						Children: nil,
						Value:    "*",
						Type:     TypeAny,
					},
				},
				Value: "test",
				Type:  TypeText,
			},
		},
		{
			input: "test1*test2",
			output: &Node{
				Children: []*Node{
					{
						Children: []*Node{
							{
								Children: nil,
								Value:    "test2",
								Type:     TypeText,
							},
						},
						Value: "*",
						Type:  TypeAny,
					},
				},
				Value: "test1",
				Type:  TypeText,
			},
		},
		{
			input: "*",
			output: &Node{
				Children: nil,
				Value:    "*",
				Type:     TypeAny,
			},
		},
		{
			input: "",
			output: &Node{
				Children: nil,
				Value:    "",
				Type:     TypeText,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			require := r.New(t)

			output := Parse(test.input)

			require.Equal(test.output, output)
		})
	}
}
