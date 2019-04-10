package parser

import (
	"fmt"
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
				Type:  TypeRoot,
				Value: "",
				Children: []*Node{
					{
						Children: nil,
						Value:    "test",
						Type:     TypeText,
					},
				},
			},
		},
		{
			input: "test*",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Children: []*Node{
					{
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
			},
		},
		{
			input: "test1*test2",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Children: []*Node{
					{
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
			},
		},
		{
			input: "*",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Children: []*Node{
					{
						Children: nil,
						Value:    "*",
						Type:     TypeAny,
					},
				},
			},
		},
		{
			input: "",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Children: []*Node{
					{
						Children: nil,
						Value:    "",
						Type:     TypeText,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			require := r.New(t)

			output := Parse(test.input)

			require.Equal(test.output, output)

			require.Equal(test.input, output.String())
		})
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		inputs       []string
		output       *Node
		outputString string
	}{
		{
			inputs: []string{
				"test",
				"test2",
			},
			outputString: "(test|test2)",
			output: &Node{
				Value: "",
				Type:  TypeRoot,
				Children: []*Node{
					{
						Children: nil,
						Value:    "test",
						Type:     TypeText,
					},
					{
						Children: nil,
						Value:    "test2",
						Type:     TypeText,
					},
				},
			},
		},
		{
			inputs: []string{
				"test*",
				"test*2",
			},
			outputString: "test*2",
			output: &Node{
				Value: "",
				Type:  TypeRoot,
				Children: []*Node{
					{
						Value: "test",
						Type:  TypeText,
						Children: []*Node{
							{
								Value: "*",
								Type:  TypeAny,
								Children: []*Node{
									{
										Children: nil,
										Value:    "2",
										Type:     TypeText,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			inputs: []string{
				"a",
				"a*b",
			},
			outputString: "a*b",
			output: &Node{
				Value: "",
				Type:  TypeRoot,
				Children: []*Node{
					{
						Value: "a",
						Type:  TypeText,
						Children: []*Node{
							{
								Value: "*",
								Type:  TypeAny,
								Children: []*Node{
									{
										Children: nil,
										Value:    "b",
										Type:     TypeText,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			inputs: []string{
				"*a",
				"*b",
			},
			outputString: "*(a|b)",
			output: &Node{
				Value: "",
				Type:  TypeRoot,
				Children: []*Node{
					{
						Value: "*",
						Type:  TypeAny,
						Children: []*Node{
							{
								Value:    "a",
								Type:     TypeText,
								Children: nil,
							},
							{
								Value:    "b",
								Type:     TypeText,
								Children: nil,
							},
						},
					},
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			require := r.New(t)

			outputs := make([]*Node, 0)
			for _, input := range test.inputs {
				outputs = append(outputs, Parse(input))
			}

			fmt.Println(outputs)
			final := outputs[0]
			for _, output := range outputs[1:] {
				final = Merge(final, output)
			}

			require.Equal(test.output, final)
			require.Equal(test.outputString, final.String())
		})
	}
}
