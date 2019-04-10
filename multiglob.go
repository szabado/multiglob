package multiglob

import (
	"github.com/szabado/multiglob/internal/parser"
	"strings"
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
	stack := stack{
		sl: make([]frame, 0),
	}

	stack.pushChildren(mg.node, rawInput, false)

	for !stack.empty() {
		f := stack.pop()

		//fmt.Println("frame: ", f.input, f.any, f.n)
		switch f.n.Type {
		case parser.TypeAny:
			if f.n.Leaf {
				//fmt.Println("any leaf node")
				return true
			} else {
				//fmt.Println("any node. Adding children")
				stack.pushChildren(f.n, f.input, true)
			}
		case parser.TypeText:
			//fmt.Println("text type")
			if f.any {
				//fmt.Println("following an any")
				if f.n.Leaf && strings.HasSuffix(f.input, f.n.Value) {
					return true
				} else if i := strings.Index(f.input, f.n.Value); i >= 0 {
					trunc := f.input[i+len(f.n.Value):]
					stack.push(frame{
						input: trunc,
						any:   true,
						n:     f.n,
					})
					stack.pushChildren(f.n, trunc, false)
				}
			} else {
				//fmt.Println("direct matching")
				if f.n.Leaf && f.n.Value == f.input {
					return true
				} else if strings.HasPrefix(f.input, f.n.Value) {
					trunc := f.input[len(f.n.Value):]
					stack.push(frame{
						input: trunc,
						any:   true,
						n:     f.n,
					})
					stack.pushChildren(f.n, trunc, false)
				}
			}
		}
	}

	return false
}

type frame struct {
	n     *parser.Node
	input string
	any   bool
}

type stack struct {
	sl []frame
}

func (s *stack) pop() frame {
	i := len(s.sl) - 1
	n := s.sl[i]
	s.sl = s.sl[:i]
	return n
}

func (s *stack) empty() bool {
	return len(s.sl) == 0
}

func (s *stack) push(frames ...frame) {
	//fmt.Println(frames)
	for _, n := range frames {
		s.sl = append(s.sl, n)
	}
}

func (s *stack) pushChildren(n *parser.Node, input string, any bool) {
	for _, c := range n.Children {
		s.push(frame{
			input: input,
			n:     c,
			any:   any,
		})
	}
}
