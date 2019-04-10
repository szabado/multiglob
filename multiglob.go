package multiglob

import (
	"strings"

	"github.com/szabado/multiglob/internal/parser"
)

type Builder struct {
	patterns map[string]*parser.Node
}

func New() *Builder {
	return &Builder{
		patterns: make(map[string]*parser.Node),
	}
}

func (m *Builder) AddPattern(name, pattern string) {
	m.patterns[name] = parser.Parse(pattern)
}

func (m *Builder) Build() MultiGlob {
	var final *parser.Node
	for _, p := range m.patterns {
		if final == nil {
			final = p
		} else {
			final = parser.Merge(final, p)
		}
	}

	return MultiGlob{
		node: final,
	}
}

type MultiGlob struct {
	node *parser.Node
}

func (mg *MultiGlob) Match(rawInput string) bool {
	return match(mg.node, rawInput, false)
}

func match(node *parser.Node, input string, any bool) bool {
	var childInput string
	var childAny bool
	switch node.Type {
	case parser.TypeAny:
		if node.Leaf {
			//fmt.Println("any leaf node")
			return true
		} else {
			//fmt.Println("any node. Adding children")

			childInput = input
			childAny = true
		}
	case parser.TypeText:
		//fmt.Println("text type")
		if any {
			//fmt.Println("following an any")
			if node.Leaf && strings.HasSuffix(input, node.Value) {
				return true
			} else if i := strings.Index(input, node.Value); i >= 0 {
				trunc := input[i+len(node.Value):]
				if match(node, trunc, true) {
					return true
				}

				childInput = trunc
				childAny = false
			}
		} else {
			//fmt.Println("direct matching")
			//fmt.Println(input, " ", node.Value)
			if node.Leaf && node.Value == input {
				//fmt.Println("leaf matches")
				return true
			} else if strings.HasPrefix(input, node.Value) {
				trunc := input[len(node.Value):]

				if match(node, trunc, true) {
					return true
				}

				childAny = false
				childInput = trunc
			}
		}
	case parser.TypeRoot:
		fallthrough
	default:
		childInput = input
		childAny = any
	}

	for _, c := range node.Children {
		if match(c, childInput, childAny) {
			return true
		}
	}
	return false
}
