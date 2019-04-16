package multiglob

import (
	"strings"

	"github.com/pkg/errors"

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
	p, err := parser.Parse(name, pattern)
	if err != nil {
		return errors.Wrap(err, "failed to add pattern")
	}
	m.patterns[name] = p
	return err
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
	_, matched := match(mg.node, input, false)
	return matched
}

// FindAllPatterns returns a list containing all patterns that matched this input.
func (mg *MultiGlob) FindAllPatterns(input string) []string {
	results, _ := match(mg.node, input, true)
	duplicates := make(map[string]bool)

	cleaned := make([]string, 0, len(results))
	for _, result := range results {
		if duplicates[result] {
			continue
		}
		duplicates[result] = true
		cleaned = append(cleaned, result)
	}

	return cleaned
}

// FindPattern returns one pattern out of the set of patterns that matches input.
// There is no guarantee as to which of the patterns will be returned. Returns true
// if a pattern was matched.
func (mg *MultiGlob) FindPattern(input string) (string, bool) {
	results, ok := match(mg.node, input, false)
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
	globs := make([]string, 0)

	for leafConsumed := false; !leafConsumed && ast != nil; {
		switch ast.Type {
		case parser.TypeText:
			if !strings.HasPrefix(input, ast.Value) {
				return nil, errTextNotFound
			}
			input = trimString(input, len(ast.Value))
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

func match(node *parser.Node, input string, exhaustive bool) ([]string, bool) {
	var (
		results []string
	)

	switch node.Type {
	case parser.TypeAny:
		if node.Leaf {
			if !exhaustive {
				return node.Name, true
			}
			results = merge(results, node.Name)
		}

		for _, child := range node.Children {
			tempInput := input
			for i := child.Index(tempInput); i >= 0; i = child.Index(tempInput) {
				if i < 0 {
					continue
				}

				names, ok := match(child, trimString(tempInput, i), exhaustive)
				tempInput = trimString(tempInput, i+len(child.Value))

				if !ok {
					continue
				}

				if !exhaustive {
					return names, true
				}
				results = merge(results, names)
			}
		}
	case parser.TypeText:
		if node.Leaf && node.Value == input {
			if !exhaustive {
				return node.Name, true
			}
			results = merge(results, node.Name)
		} else if !strings.HasPrefix(input, node.Value) {
			return nil, false
		}

		input = trimString(input, len(node.Value))

		for _, c := range node.Children {
			names, ok := match(c, input, exhaustive)
			if !ok {
				continue
			}
			if !exhaustive {
				return names, true
			}
			results = merge(results, names)
		}

	case parser.TypeRange:
		short := input
		for _, r := range input {
			if !isConsumable(r, node.Range) {
				break
			}

			short = strings.TrimPrefix(short, string(r))

			for _, child := range node.Children {
				names, ok := match(child, short, exhaustive)
				if !ok {
					continue
				}

				if !exhaustive {
					return names, true
				}
				results = merge(results, names)
			}

			if !node.Range.Repeated {
				break
			}
		}

		if node.Leaf && short == "" && len(short) != len(input) {
			results = append(results, node.Name...)
		}
	case parser.TypeRoot:
		for _, c := range node.Children {
			names, ok := match(c, input, exhaustive)
			if !ok {
				continue
			}
			if !exhaustive {
				return names, true
			}
			results = merge(results, names)
		}
	}

	return results, len(results) != 0
}

func trimString(s string, prefixLen int) string {
	if len(s) <= prefixLen {
		return ""
	}
	return s[prefixLen:]
}

// isConsumable checks if the rune can be consumed, based on the rules in rnge. If rnge is an inverse
// parser.Range, it'll return true if the rune doesn't match any of the rnge's rules. If rnge is a
// normal range, it'll return true if the rune matches any of the rnge's rules.
func isConsumable(r rune, rnge *parser.Range) bool {
	contains := rnge.Contains(r)
	if rnge.Inverse {
		return !contains
	}

	return contains
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
