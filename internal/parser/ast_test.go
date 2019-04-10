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
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "test",
						Type:     TypeText,
						Leaf:     true,
					},
				},
			},
		},
		{
			input: "test*",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "test",
						Type:  TypeText,
						Leaf:  false,
						Children: []*Node{
							{
								Children: nil,
								Value:    "*",
								Type:     TypeAny,
								Leaf:     true,
							},
						},
					},
				},
			},
		},
		{
			input: "test1*test2",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "test1",
						Type:  TypeText,
						Leaf:  false,
						Children: []*Node{
							{
								Value: "*",
								Type:  TypeAny,
								Leaf:  false,
								Children: []*Node{
									{
										Value:    "test2",
										Type:     TypeText,
										Leaf:     true,
										Children: nil,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: "*",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "*",
						Leaf:     true,
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
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "",
						Leaf:     true,
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

			//require.Equal(test.input, output.String())
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
						Leaf:     true,
					},
					{
						Children: nil,
						Value:    "test2",
						Type:     TypeText,
						Leaf:     true,
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
								Leaf:  true,
								Children: []*Node{
									{
										Children: nil,
										Value:    "2",
										Type:     TypeText,
										Leaf:     true,
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
						Leaf:  true,
						Children: []*Node{
							{
								Value: "*",
								Type:  TypeAny,
								Children: []*Node{
									{
										Children: nil,
										Value:    "b",
										Type:     TypeText,
										Leaf:     true,
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
								Leaf:     true,
							},
							{
								Value:    "b",
								Type:     TypeText,
								Children: nil,
								Leaf:     true,
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
			//require.Equal(test.outputString, final.String())
		})
	}
}
