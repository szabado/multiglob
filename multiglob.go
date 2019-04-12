package multiglob

import (
	"errors"
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
func (m *Builder) Compile() (*MultiGlob, error) {
	var final *parser.Node
	for _, p := range m.patterns {
		if final == nil {
			final = p
		} else {
			final = parser.Merge(final, p)
		}
	}

	patterns := make(map[string]*parser.Node)
	for k, v := range m.patterns {
		patterns[k] = v
	}

	return &MultiGlob{
		node:     final,
		patterns: patterns,
	}, nil
}

// MustCompile wraps Compile, and panics if there is an error.
func (m *Builder) MustCompile() *MultiGlob {
	mg, err := m.Compile()
	if err != nil {
		panic(err)
	}
	return mg
}

// MultiGlob is a matcher that is built from a collection of patterns. See Builder.
type MultiGlob struct {
	node     *parser.Node
	patterns map[string]*parser.Node
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

// FindPattern returns one pattern out of the set of patterns that matches input.
// There is no guarantee as to which of the patterns will be returned. Returns true
// if a pattern was matched.
func (mg *MultiGlob) FindPattern(input string) (string, bool) {
	results, ok := match(mg.node, input, false, false)
	if !ok || len(results) < 1 {
		return "", false
	}
	return results[0], true
}

// FindAllGlobs returns a map of pattern names to globs extracted using each pattern.
// It uses all the patterns returned FindAllPatterns. See FindGlobs for an explanation
// of glob extraction.
func (mg *MultiGlob) FindAllGlobs(input string) map[string][]string {
	patternNames := mg.FindAllPatterns(input)

	globs := make(map[string][]string)
	for _, name := range patternNames {
		g, _ := extractGlobs(input, mg.patterns[name])
		globs[name] = g
	}

	return globs
}

// FindGlobs finds a matching pattern using FindPattern, and then extracts the globs
// from the input based on that pattern. It also returns the name of the pattern
// matched. This uses a greedy matching algorithm. For example:
//
//   Input:         "test"
//   Pattern Found: "t*t"
//   Globs:         ["es"]
//
//   Input:         "pen pineapple apple pen"
//   Pattern Found: "*apple*"
//   Globs:         ["pen pineapple ", " pen"]
func (mg *MultiGlob) FindGlobs(input string) (name string, globs []string, matched bool) {
	name, ok := mg.FindPattern(input)
	if !ok {
		return "", nil, false
	}

	globs, _ = extractGlobs(input, mg.patterns[name])
	return name, globs, true
}

// FindGlobsForPattern extracts the globs from input using the named pattern.
func (mg *MultiGlob) FindGlobsForPattern(input, name string) (globs []string, err error) {
	if ast, ok := mg.patterns[name]; !ok {
		return nil, errors.New("pattern not found")
	} else {
		globs, err := extractGlobs(input, ast)
		if err != nil {
			return nil, errors.New("pattern did not match input")
		}
		return globs, nil
	}
}

var errTextNotFound = errors.New("text not found")

// extractGlobs returns the globs based on the pattern. It either returns a nil error or
// errTextNotFound
func extractGlobs(input string, ast *parser.Node) ([]string, error) {
	// TODO: decrease allows by making this ceil(half the tree size)
	globs := make([]string, 0)

	for leafConsumed := false; !leafConsumed && ast != nil; {
		switch ast.Type {
		case parser.TypeText:
			if !strings.HasPrefix(input, ast.Value) {
				return nil, errTextNotFound
			}
			input = strings.TrimPrefix(input, ast.Value)
			if ast.Leaf {
				leafConsumed = true
			}
		case parser.TypeAny:
			// It's globbing time, baby!
			if ast.Leaf {
				globs = append(globs, input)
				leafConsumed = true
				break
			} else if input == "" {
				return nil, errTextNotFound
			}

			// Consume as much as possible, and then slowly consume less until we find a
			// match or can't consume any less
			for globbed := input; globbed != ""; {
				child := ast.Children[0]

				globEnds := strings.LastIndex(globbed, child.Value)
				if globEnds < 0 {
					return nil, errTextNotFound
				}

				globbed = globbed[:globEnds]

				subglobs, err := extractGlobs(input[globEnds:], child)
				if err != nil {
					continue
				}

				// we found our match!
				globs = append(globs, globbed)
				globs = append(globs, subglobs...)
				leafConsumed = true
				break
			}
		}

		if !ast.Leaf {
			ast = ast.Children[0]
		}
	}

	return globs, nil
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
