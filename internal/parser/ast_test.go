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
		{
			name:  "test",
			input: `\[`,
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "[",
						Leaf:     true,
						Type:     TypeText,
						Name:     []string{"test"},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			require := r.New(t)

			output, err := Parse(test.name, test.input)
			require.NoError(err)

			require.Equal(test.output, output)
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		charRange *Range
		input     rune
		output    bool
	}{
		{
			input: 'a',
			charRange: &Range{
				CharList: "abc",
			},
			output: true,
		},
		{
			input: 'a',
			charRange: &Range{
				Bounds: []*Bounds{
					{
						Low:  'a',
						High: 'z',
					},
				},
			},
			output: true,
		},
		{
			input: 'A',
			charRange: &Range{
				CharList: "abc",
			},
			output: false,
		},
		{
			input: 'A',
			charRange: &Range{
				Bounds: []*Bounds{
					{
						Low:  'a',
						High: 'z',
					},
				},
			},
			output: false,
		},
		{
			input: 'a',
			charRange: &Range{
				Inverse:  true,
				CharList: "abc",
			},
			output: false,
		},
		{
			input: 'a',
			charRange: &Range{
				Inverse: true,
				Bounds: []*Bounds{
					{
						Low:  'a',
						High: 'z',
					},
				},
			},
			output: false,
		},
		{
			input: 'A',
			charRange: &Range{
				Inverse:  true,
				CharList: "abc",
			},
			output: true,
		},
		{
			input: 'A',
			charRange: &Range{
				Inverse: true,
				Bounds: []*Bounds{
					{
						Low:  'a',
						High: 'z',
					},
				},
			},
			output: true,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			require := r.New(t)
			require.Equal(test.output, test.charRange.Matches(test.input))
		})
	}
}

func TestParseRanges(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output *Node
		err    bool
	}{
		{
			name:  "test",
			input: "[abc]",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "",
						Type:     TypeRange,
						Leaf:     true,
						Name:     []string{"test"},
						Range: &Range{
							CharList: "abc",
						},
					},
				},
			},
		},
		{
			name:  "test",
			input: "[abc]+",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "",
						Type:     TypeRange,
						Leaf:     true,
						Name:     []string{"test"},
						Range: &Range{
							Repeated: true,
							CharList: "abc",
						},
					},
				},
			},
		},
		{
			name:  "test",
			input: "[^abc]",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "",
						Type:     TypeRange,
						Leaf:     true,
						Name:     []string{"test"},
						Range: &Range{
							CharList: "abc",
							Inverse:  true,
						},
					},
				},
			},
		},
		{
			name:  "test",
			input: "[a-c]",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "",
						Type:     TypeRange,
						Leaf:     true,
						Name:     []string{"test"},
						Range: &Range{
							Bounds: []*Bounds{
								{
									Low:  rune('a'),
									High: rune('c'),
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "test",
			input: "[^a-c]",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "",
						Type:     TypeRange,
						Leaf:     true,
						Name:     []string{"test"},
						Range: &Range{
							Inverse: true,
							Bounds: []*Bounds{
								{
									Low:  rune('a'),
									High: rune('c'),
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "test",
			input: "[ab-c]",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "",
						Type:     TypeRange,
						Leaf:     true,
						Name:     []string{"test"},
						Range: &Range{
							CharList: "a",
							Bounds: []*Bounds{
								{
									Low:  rune('b'),
									High: rune('c'),
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "test",
			input: "[^ab-c]",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "",
						Type:     TypeRange,
						Leaf:     true,
						Name:     []string{"test"},
						Range: &Range{
							Inverse:  true,
							CharList: "a",
							Bounds: []*Bounds{
								{
									Low:  rune('b'),
									High: rune('c'),
								},
							},
						},
					},
				},
			},
		},
		{
			name:   "test",
			input:  "[a",
			err:    true,
			output: nil,
		},
		{
			name:   "test",
			input:  "[^a",
			err:    true,
			output: nil,
		},
		{
			name:   "test",
			input:  "[b-a]",
			err:    true,
			output: nil,
		},
		{
			name:   "test",
			input:  "[^b-a]",
			err:    true,
			output: nil,
		},
		{
			name:  "test",
			input: "]a",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "]a",
						Type:     TypeText,
						Leaf:     true,
						Name:     []string{"test"},
					},
				},
			},
		},
		{
			name:  "test",
			input: "[-]",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "",
						Type:     TypeRange,
						Leaf:     true,
						Name:     []string{"test"},
						Range: &Range{
							CharList: "-",
						},
					},
				},
			},
		},
		{
			name:  "test",
			input: "[^-]",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "",
						Type:     TypeRange,
						Leaf:     true,
						Name:     []string{"test"},
						Range: &Range{
							Inverse:  true,
							CharList: "-",
						},
					},
				},
			},
		},
		{
			name:   "test",
			input:  "[a-]",
			output: nil,
			err:    true,
		},
		{
			name:   "test",
			input:  "[^a-]",
			output: nil,
			err:    true,
		},
		{
			name:  "test",
			input: "[]]",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "",
						Type:     TypeRange,
						Leaf:     true,
						Name:     []string{"test"},
						Range: &Range{
							CharList: "]",
						},
					},
				},
			},
		},
		{
			name:  "test",
			input: "[^]]",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "",
						Type:     TypeRange,
						Leaf:     true,
						Name:     []string{"test"},
						Range: &Range{
							Inverse:  true,
							CharList: "]",
						},
					},
				},
			},
		},
		{
			name:  "test",
			input: "[ ]]",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Range: &Range{
							CharList: " ",
						},
						Children: []*Node{
							{
								Leaf:  true,
								Name:  []string{"test"},
								Value: "]",
								Type:  TypeText,
							},
						},
					},
				},
			},
		},
		{
			name:  "test",
			input: "[^ ]]",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Range: &Range{
							Inverse:  true,
							CharList: " ",
						},
						Children: []*Node{
							{
								Leaf:  true,
								Name:  []string{"test"},
								Value: "]",
								Type:  TypeText,
							},
						},
					},
				},
			},
		},
		{
			name:  "test",
			input: "[[]",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "",
						Type:     TypeRange,
						Leaf:     true,
						Name:     []string{"test"},
						Range: &Range{
							CharList: "[",
						},
					},
				},
			},
		},
		{
			name:  "test",
			input: "[^[]",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Children: nil,
						Value:    "",
						Type:     TypeRange,
						Leaf:     true,
						Name:     []string{"test"},
						Range: &Range{
							Inverse:  true,
							CharList: "[",
						},
					},
				},
			},
		},
		{
			name:  "test",
			input: "[ []",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Leaf:  true,
						Name:  []string{"test"},
						Range: &Range{
							CharList: " [",
						},
						Children: nil,
					},
				},
			},
		},
		{
			name:  "test",
			input: "[^ []",
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Leaf:  true,
						Name:  []string{"test"},
						Range: &Range{
							Inverse:  true,
							CharList: " [",
						},
						Children: nil,
					},
				},
			},
		},
		{
			name:  "test",
			input: `[\^]`,
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Leaf:  true,
						Name:  []string{"test"},
						Range: &Range{
							CharList: "^",
						},
						Children: nil,
					},
				},
			},
		},
		{
			name:  "test",
			input: `[^\^]`,
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Leaf:  true,
						Name:  []string{"test"},
						Range: &Range{
							Inverse:  true,
							CharList: "^",
						},
						Children: nil,
					},
				},
			},
		},
		{
			name:  "test",
			input: `[ ^]`,
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Leaf:  true,
						Name:  []string{"test"},
						Range: &Range{
							CharList: " ^",
						},
						Children: nil,
					},
				},
			},
		},
		{
			name:  "test",
			input: `[^ ^]`,
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Leaf:  true,
						Name:  []string{"test"},
						Range: &Range{
							Inverse:  true,
							CharList: " ^",
						},
						Children: nil,
					},
				},
			},
		},
		{
			name:  "test",
			input: `[ \]]`,
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Leaf:  true,
						Name:  []string{"test"},
						Range: &Range{
							CharList: " ]",
						},
						Children: nil,
					},
				},
			},
		},
		{
			name:  "test",
			input: `[^ \]]`,
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Leaf:  true,
						Name:  []string{"test"},
						Range: &Range{
							Inverse:  true,
							CharList: " ]",
						},
						Children: nil,
					},
				},
			},
		},
		{
			name:  "test",
			input: `[ \-]`,
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Leaf:  true,
						Name:  []string{"test"},
						Range: &Range{
							CharList: " -",
						},
						Children: nil,
					},
				},
			},
		},
		{
			name:  "test",
			input: `[^ \-]`,
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Leaf:  true,
						Name:  []string{"test"},
						Range: &Range{
							Inverse:  true,
							CharList: " -",
						},
						Children: nil,
					},
				},
			},
		},
		{
			name:  "test",
			input: `[^^]`,
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Leaf:  true,
						Name:  []string{"test"},
						Range: &Range{
							Inverse:  true,
							CharList: "^",
						},
						Children: nil,
					},
				},
			},
		},
		{
			name:  "test",
			input: `[\--\-]`,
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Leaf:  true,
						Name:  []string{"test"},
						Range: &Range{
							Bounds: []*Bounds{
								{
									Low:  '-',
									High: '-',
								},
							},
						},
						Children: nil,
					},
				},
			},
		},
		{
			name:  "test",
			input: `[^\^-\^]`,
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Leaf:  true,
						Name:  []string{"test"},
						Range: &Range{
							Inverse: true,
							Bounds: []*Bounds{
								{
									Low:  '^',
									High: '^',
								},
							},
						},
						Children: nil,
					},
				},
			},
		},
		{
			name:  "test",
			input: `[^^-^]`,
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Leaf:  true,
						Name:  []string{"test"},
						Range: &Range{
							Inverse: true,
							Bounds: []*Bounds{
								{
									Low:  '^',
									High: '^',
								},
							},
						},
						Children: nil,
					},
				},
			},
		},
		{
			name:  "test",
			input: `[\-\--\-\-]`,
			output: &Node{
				Type:  TypeRoot,
				Value: "",
				Leaf:  false,
				Children: []*Node{
					{
						Value: "",
						Type:  TypeRange,
						Leaf:  true,
						Name:  []string{"test"},
						Range: &Range{
							CharList: "--",
							Bounds: []*Bounds{
								{
									Low:  '-',
									High: '-',
								},
							},
						},
						Children: nil,
					},
				},
			},
		},
		{
			name:  "test",
			input: `[\]`,
			err:   true,
		},
		{
			name:  "test",
			input: `[\a]`,
			err:   true,
		},
		{
			name:  "test",
			input: `\a`,
			err:   true,
		},
		{
			name:  `test`,
			input: `\`,
			err:   true,
		},
		{
			name:  `test`,
			input: `a[`,
			err:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			require := r.New(t)

			output, err := Parse(test.name, test.input)
			if test.err {
				require.Error(err)
			} else {
				require.NoError(err)
			}

			require.Equal(test.output, output)
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
		{
			inputs: []string{
				"[a]",
				"[ab]",
			},
			output: &Node{
				Value: "",
				Type:  TypeRoot,
				Children: []*Node{
					{
						Children: nil,
						Type:     TypeRange,
						Leaf:     true,
						Range: &Range{
							CharList: "a",
						},
						Name: []string{"0"},
					},
					{
						Children: nil,
						Type:     TypeRange,
						Leaf:     true,
						Range: &Range{
							CharList: "ab",
						},
						Name: []string{"1"},
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
				ast, err := Parse(fmt.Sprint(inputNum), input)
				require.NoError(err)
				outputs = append(outputs, ast)
			}

			final := outputs[0]
			for _, output := range outputs[1:] {
				final = Merge(final, output)
			}

			require.Equal(test.output, final)
		})
	}
}

func TestCanMerge(t *testing.T) {
	require := r.New(t)

	var n1, n2 *Node
	require.True(n1.canMerge(n2))
}

func TestIndex(t *testing.T) {
	tests := []struct {
		input        string
		inputPattern string
		output       int
	}{
		{
			input:        "12341234",
			inputPattern: "[234]",
			output:       1,
		},
		{
			input:        "12341234",
			inputPattern: "potato",
			output:       -1,
		},
		{
			input:        "12341234",
			inputPattern: "*",
			output:       0,
		},
		{
			input:        "12341234",
			inputPattern: "4",
			output:       3,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			require := r.New(t)

			node, err := Parse(fmt.Sprint(i), test.inputPattern)
			require.NoError(err)

			require.Equal(test.output, node.Children[0].Index(test.input))
		})
	}
}

func TestLastIndex(t *testing.T) {
	tests := []struct {
		input        string
		inputPattern string
		output       int
	}{
		{
			input:        "12341234",
			inputPattern: "[234]",
			output:       5,
		},
		{
			input:        "12341234",
			inputPattern: "potato",
			output:       -1,
		},
		{
			input:        "12341234",
			inputPattern: "*",
			output:       7,
		},
		{
			input:        "12341234",
			inputPattern: "4",
			output:       7,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			require := r.New(t)

			node, err := Parse(fmt.Sprint(i), test.inputPattern)
			require.NoError(err)

			require.Equal(test.output, node.Children[0].LastIndex(test.input))
		})
	}
}

func TestMultipleMerges(t *testing.T) {
	require := r.New(t)

	astA, err := Parse("A", "a")
	require.NoError(err)
	//astA := Parse("b", "b")
	astC, err := Parse("C", "c")
	require.NoError(err)

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

func TestCompress(t *testing.T) {
	tests := []struct {
		input  *Node
		output *Node
	}{
		{
			input: &Node{
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
								Value: "a",
								Type:  TypeText,
								Leaf:  false,
								Children: []*Node{
									{
										Value: "*",
										Type:  TypeText,
										Leaf:  true,
										Name:  []string{"2"},
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
					},
				},
			},
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
								Value: "a*",
								Type:  TypeText,
								Leaf:  true,
								Name:  []string{"2"},
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
			},
		},
		{
			input: &Node{
				Value: "",
				Type:  TypeRoot,
				Children: []*Node{
					{
						Value: "a",
						Type:  TypeText,
						Children: []*Node{
							{
								Value: "a",
								Type:  TypeText,
								Children: []*Node{
									{
										Value:    "b",
										Type:     TypeText,
										Leaf:     true,
										Name:     []string{"2"},
										Children: nil,
									},
								},
							},
						},
					},
				},
			},
			output: &Node{
				Value: "",
				Type:  TypeRoot,
				Children: []*Node{
					{
						Value:    "aab",
						Type:     TypeText,
						Leaf:     true,
						Name:     []string{"2"},
						Children: nil,
					},
				},
			},
		},
		{
			input: &Node{
				Value: "",
				Type:  TypeRoot,
				Children: []*Node{
					{
						Value: "a",
						Type:  TypeText,
						Children: []*Node{
							{
								Value: "a",
								Type:  TypeText,
								Children: []*Node{
									{
										Value:    "b",
										Type:     TypeText,
										Leaf:     true,
										Name:     []string{"2"},
										Children: nil,
									},
								},
							},
							{
								Value: "a",
								Type:  TypeText,
								Children: []*Node{
									{
										Value:    "b",
										Type:     TypeText,
										Leaf:     true,
										Name:     []string{"2"},
										Children: nil,
									},
								},
							},
						},
					},
				},
			},
			output: &Node{
				Value: "",
				Type:  TypeRoot,
				Children: []*Node{
					{
						Value: "a",
						Type:  TypeText,
						Children: []*Node{
							{
								Value:    "ab",
								Type:     TypeText,
								Leaf:     true,
								Name:     []string{"2"},
								Children: nil,
							},
							{
								Value:    "ab",
								Type:     TypeText,
								Leaf:     true,
								Name:     []string{"2"},
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

			test.input.compress()
			require.Equal(test.output, test.input)
		})
	}
}
