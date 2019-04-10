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

func (m *Builder) AddPattern(name, pattern string) error {
	m.MustAddPattern(name, pattern)
	return nil
}

func (m *Builder) MustAddPattern(name, pattern string) {
	m.patterns[name] = parser.Parse(pattern)
}

func (m *Builder) Compile() (MultiGlob, error){
	return m.MustCompile(), nil
}

func (m *Builder) MustCompile() MultiGlob {
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
	var (
		childInput string
		childAny   bool
	)

	switch node.Type {
	case parser.TypeAny:
		if node.Leaf {
			return true
		}

		childInput = input
		childAny = true

	case parser.TypeText:
		if any {
			if node.Leaf && strings.HasSuffix(input, node.Value) {
				return true
			} else if i := strings.Index(input, node.Value); i < 0 {
				return false
			} else {
				trunc := input[i+len(node.Value):]
				if match(node, trunc, true) {
					return true
				}

				childInput = trunc
				childAny = false
			}
		} else {
			if node.Leaf && node.Value == input {
				return true
			} else if !strings.HasPrefix(input, node.Value) {
				return false
			}

			trunc := input[len(node.Value):]

			childAny = false
			childInput = trunc
		}
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
