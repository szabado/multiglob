package parser

import (
	"fmt"
	"testing"

	r "github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output *Node
	}{
		{
			name:  "test1",
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
						Name:     []string{"test1"},
					},
				},
			},
		},
		{
			name:  "test2",
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
								Name:     []string{"test2"},
							},
						},
					},
				},
			},
		},
		{
			name:  "testIII",
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
										Name:     []string{"testIII"},
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
			name:  "testiv",
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
						Name:     []string{"testiv"},
					},
				},
			},
		},
		{
			name:  "testfive",
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
						Name:     []string{"testfive"},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			require := r.New(t)

			output := Parse(test.name, test.input)

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
						Leaf:     true,
						Name:     []string{"0"},
					},
					{
						Children: nil,
						Value:    "test2",
						Type:     TypeText,
						Leaf:     true,
						Name:     []string{"1"},
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
								Name:  []string{"0"},
								Children: []*Node{
									{
										Children: nil,
										Value:    "2",
										Type:     TypeText,
										Leaf:     true,
										Name:     []string{"1"},
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
				"a*b",
				"a",
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
						Name:  []string{"1"},
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
										Name:     []string{"0"},
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
								Name:     []string{"0"},
							},
							{
								Value:    "b",
								Type:     TypeText,
								Children: nil,
								Leaf:     true,
								Name:     []string{"1"},
							},
						},
					},
				},
			},
		},
		{
			inputs: []string{
				"test",
				"test",
			},
			outputString: "test",
			output: &Node{
				Value: "",
				Type:  TypeRoot,
				Children: []*Node{
					{
						Children: nil,
						Value:    "test",
						Type:     TypeText,
						Leaf:     true,
						Name:     []string{"0", "1"},
					},
				},
			},
		},
		{
			inputs: []string{
				"*",
				"*",
			},
			outputString: "*",
			output: &Node{
				Value: "",
				Type:  TypeRoot,
				Children: []*Node{
					{
						Children: nil,
						Value:    "*",
						Type:     TypeAny,
						Leaf:     true,
						Name:     []string{"0", "1"},
					},
				},
			},
		},
		{
			inputs: []string{
				"a",
				"*",
			},
			outputString: "(a|*)",
			output: &Node{
				Value: "",
				Type:  TypeRoot,
				Children: []*Node{
					{
						Children: nil,
						Value:    "a",
						Type:     TypeText,
						Leaf:     true,
						Name:     []string{"0"},
					},
					{
						Children: nil,
						Value:    "*",
						Type:     TypeAny,
						Leaf:     true,
						Name:     []string{"1"},
					},
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			require := r.New(t)

			outputs := make([]*Node, 0)
			for inputNum, input := range test.inputs {
				outputs = append(outputs, Parse(fmt.Sprint(inputNum), input))
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

func TestMultipleMerges(t *testing.T) {
	require := r.New(t)

	astA := Parse("A", "a")
	//astA := Parse("b", "b")
	astC := Parse("C", "c")

	ast := Merge(astA, astC)
	ast = Merge(astA, ast)

	require.Equal(&Node{
		Value: "",
		Type:  TypeRoot,
		Children: []*Node{
			{
				Children: nil,
				Value:    "a",
				Type:     TypeText,
				Leaf:     true,
				Name:     []string{"A", "A"},
			},
			{
				Children: nil,
				Value:    "c",
				Type:     TypeText,
				Leaf:     true,
				Name:     []string{"C"},
			},
		},
	}, ast)
}
