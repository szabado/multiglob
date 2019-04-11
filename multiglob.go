package multiglob

import (
	"strings"

	"github.com/szabado/multiglob/internal/parser"
)

// Builder builds a MultiGlob.
type Builder struct {
	patterns map[string]*parser.Node
}

// New returns a new Builder that can be used to create a MultiGlob.
func New() *Builder {
	return &Builder{
		patterns: make(map[string]*parser.Node),
	}
}

// AddPattern adds the provided pattern to the builder and parses it.
func (m *Builder) AddPattern(name, pattern string) error {
	m.patterns[name] = parser.Parse(name, pattern)
	return nil
}

// MustAddPattern wraps AddPattern, and panics if there is an error.
func (m *Builder) MustAddPattern(name, pattern string) {
	err := m.AddPattern(name, pattern)
	if err != nil {
		panic(err)
	}
}

// Compile merges all the compiled patterns into one MultiGlob and returns it.
func (m *Builder) Compile() (MultiGlob, error) {
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
	}, nil
}

// MustCompile wraps Compile, and panics if there is an error.
func (m *Builder) MustCompile() MultiGlob {
	mg, err := m.Compile()
	if err != nil {
		panic(err)
	}
	return mg
}

// MultiGlob is a matcher that is built from a collection of patterns. See Builder.
type MultiGlob struct {
	node *parser.Node
}

// Match determines if any pattern matches the provided string.
func (mg *MultiGlob) Match(input string) bool {
	_, matched := match(mg.node, input, false, false)
	return matched
}

// FindAllPatterns returns a list containing all patterns that matched this input.
func (mg *MultiGlob) FindAllPatterns(input string) []string {
	results, _ := match(mg.node, input, false, true)
	return results
}

func match(node *parser.Node, input string, any, exhaustive bool) ([]string, bool) {
	var (
		childInput string
		childAny   bool
		results    []string
	)

	switch node.Type {
	case parser.TypeAny:
		if node.Leaf {
			if !exhaustive {
				return node.Name, true
			}
			results = merge(results, node.Name)
		}

		childInput = input
		childAny = true

	case parser.TypeText:
		if any {
			if node.Leaf && strings.HasSuffix(input, node.Value) {
				if !exhaustive {
					return node.Name, true
				}
				results = merge(results, node.Name)
			} else if i := strings.Index(input, node.Value); i < 0 {
				return nil, false
			} else {
				trunc := input[i+len(node.Value):]
				if r, ok := match(node, trunc, true, exhaustive); ok {
					if !exhaustive {
						return r, true
					}
					results = merge(results, node.Name)
				}

				childInput = trunc
				childAny = false
			}
		} else {
			if node.Leaf && node.Value == input {
				if !exhaustive {
					return node.Name, true
				}
				results = merge(results, node.Name)
			} else if !strings.HasPrefix(input, node.Value) {
				return nil, false
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
		sl, ok := match(c, childInput, childAny, exhaustive)
		if ok {
			if exhaustive {
				results = merge(results, sl)
			} else {
				return sl, true
			}
		}
	}
	return results, len(results) != 0
}

func merge(sl1, sl2 []string) []string {
	if sl2 == nil {
		return sl1
	} else if sl1 == nil {
		return sl2
	} else {
		return append(sl1, sl2...)
	}
}
