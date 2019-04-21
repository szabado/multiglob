package parser

import (
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"

	"github.com/szabado/multiglob/internal/parser/lexer"
)

type NodeType int

//go:generate stringer -type=NodeType
const (
	TypeRoot NodeType = iota
	TypeAny
	TypeText
	TypeRange
)

func newBounds(low, high rune) (*Bounds, error) {
	if high < low {
		return nil, errors.Errorf("character range (%c, %c) is out of order",
			low,
			high)
	}

	return &Bounds{
		Low:  low,
		High: high,
	}, nil
}

type Bounds struct {
	Low, High rune
}

func (b *Bounds) Contains(r rune) bool {
	return b.Low <= r && r <= b.High
}

type Range struct {
	Repeated bool
	Inverse  bool
	Bounds   []*Bounds
	CharList string
}

// Matches returns true if the rune is matched by the Range.
func (r *Range) Matches(ru rune) bool {
	if strings.ContainsRune(r.CharList, ru) {
		return !r.Inverse
	}

	for _, bound := range r.Bounds {
		if bound.Contains(ru) {
			return !r.Inverse
		}
	}

	return r.Inverse
}

type Node struct {
	Type     NodeType
	Value    string
	Children []*Node
	Leaf     bool
	Name     []string // Only valid on leaf nodes. List of names of patterns terminate that on this leaf node
	Range    *Range
}

func (n *Node) canMerge(n2 *Node) bool {
	if n == nil || n2 == nil {
		return true
	}

	if n.Type != n2.Type {
		return false
	}

	switch n.Type {
	case TypeRoot:
		// Do nothing, Root nodes can always merge
	case TypeAny:
		// Do nothing, Any nodes can always merge
	case TypeText:
		return n.Value == n2.Value
	case TypeRange:
		// Merging ranges can be messy. Avoid it.
		return false
	}
	return true
}

func (n *Node) merge(n2 *Node) *Node {
	if n == nil {
		return n2
	} else if n2 == nil {
		return n
	}

	var children []*Node
	if len(n.Children)+len(n2.Children) != 0 {
		children = make([]*Node, len(n.Children), len(n.Children)+len(n2.Children))

		copy(children, n.Children)
		children = append(children, n2.Children...)
	}

	for i := 0; i < len(children); i++ {
		for j := i + 1; j < len(children); j++ {
			child1, child2 := children[i], children[j]

			if !child1.canMerge(child2) {
				continue
			}

			if j+1 >= len(children) {
				children = children[:j]
			} else {
				children = append(children[:j], children[j+1:]...)
			}

			children[i] = child1.merge(child2)
		}
	}

	return &Node{
		Children: children,
		Type:     n.Type,
		Value:    n.Value,
		Leaf:     n.Leaf || n2.Leaf,
		Name:     mergeNames(n, n2),
	}
}

func (n *Node) compress() {
	if len(n.Children) != 1 {
		for _, child := range n.Children {
			child.compress()
		}
		return
	}

	child := n.Children[0]
	child.compress()

	if n.Type != TypeText || child.Type != TypeText || n.Leaf {
		return
	}

	n.Value += child.Value
	n.Children = child.Children
	n.Leaf = child.Leaf
	n.Name = mergeNames(n, child)
}

// Index returns the first index of the Node's expression in the string.
func (n *Node) Index(s string) int {
	switch n.Type {
	case TypeAny:
		return 0
	case TypeText:
		return strings.Index(s, n.Value)
	case TypeRange:
		i := 0
		for r, l := utf8.DecodeRuneInString(s); !(r == utf8.RuneError && l < 2); r, l = utf8.DecodeRuneInString(s[i:]) {
			if n.Range.Matches(r) {
				return i
			}
			i += l
		}
	}
	return -1
}

// LastIndex returns the last index of the Node's expression in the string.
// For ranges, it returns the beginning of the last blob
func (n *Node) LastIndex(s string) int {
	switch n.Type {
	case TypeAny:
		return len(s) - 1
	case TypeText:
		return strings.LastIndex(s, n.Value)
	case TypeRange:
		i := len(s)
		inBlob := false
		for r, l := utf8.DecodeLastRuneInString(s); !(r == utf8.RuneError && l < 2); r, l = utf8.DecodeLastRuneInString(s[:i]) {
			contains := n.Range.Matches(r)
			if contains && !inBlob {
				inBlob = true
			} else if !contains && inBlob {
				return i
			}

			i -= l
		}
	}
	return -1
}

func mergeNames(n1, n2 *Node) []string {
	if n1.Leaf && n2.Leaf {
		return append(n1.Name, n2.Name...)
	} else if n1.Leaf {
		return n1.Name
	} else {
		return n2.Name
	}
}

func parse(name string, l *lexer.Lexer) (*Node, error) {
	if !l.Next() {
		return nil, nil
	}

	node := &Node{}

	token := l.Scan()
	switch token.Type {
	case lexer.Asterisk:
		node.Type = TypeAny
		node.Value = string(token.Value)
	case lexer.Bracket:
		if token.Value == ']' {
			node.Value = string(token.Value)
			node.Type = TypeText
			break
		}

		var (
			rnge          = &Range{}
			charCount     = 0
			previous      rune
			validChars    strings.Builder
			previousValid = false
			parsingBounds = false
			normalChar    = false
			escaped       = false
		)

		for finished := false; !finished; charCount++ {
			if !l.Next() {
				return nil, errors.New("unclosed range missing ]")
			}

			token = l.Scan()
			switch token.Type {
			case lexer.Caret:
				if charCount == 0 {
					rnge.Inverse = true
					normalChar = false
				} else {
					normalChar = true
				}
				escaped = false
			case lexer.Dash:
				if charCount == 0 || charCount == 1 && rnge.Inverse || escaped {
					normalChar = true
				} else {
					parsingBounds = true
					previousValid = false
					normalChar = false
				}
				escaped = false
			case lexer.Bracket:
				if charCount == 0 || charCount == 1 && rnge.Inverse || escaped {
					normalChar = true
				} else if token.Value == ']' {
					// Close this, handle error cases
					if parsingBounds {
						return nil, errors.Errorf("invalid range syntax %c-", previous)
					}

					if previousValid {
						validChars.WriteRune(previous)
					}
					normalChar = false
					finished = true
				} else {
					normalChar = true
				}
				escaped = false
			case lexer.Backslash:
				escaped = true
				normalChar = false
			default:
				// Treat anything unhandled as text
				fallthrough
			case lexer.Text:
				normalChar = true
				if escaped {
					return nil, errors.Errorf(`unknown escaping: \%c`, token.Value)
				}
			}

			if !normalChar {
				continue
			}

			if parsingBounds {
				b, err := newBounds(previous, token.Value)
				if err != nil {
					return nil, err
				}
				rnge.Bounds = append(rnge.Bounds, b)
				parsingBounds = false
			} else {
				if previousValid {
					validChars.WriteRune(previous)
				}
				previous = token.Value
				previousValid = true
			}
		}

		node.Type = TypeRange
		node.Range = rnge
		rnge.CharList = validChars.String()

		if nextToken, ok := l.Peek(); ok {
			if nextToken.Type == lexer.Plus {
				l.Next() // consume the plus
				node.Range.Repeated = true
			}
		}

	case lexer.Backslash:
		if !l.Next() {
			return nil, errors.New("escape found at end of pattern")
		}

		nextToken := l.Scan()
		switch nextToken.Type {
		case lexer.Bracket, lexer.Asterisk, lexer.Backslash:
			node.Value = string(nextToken.Value)
			node.Type = TypeText
		default:
			return nil, errors.Errorf(`unknown character escaping: \%c`, token.Value)
		}

		// anything other than asterisk, bracket, backslash is an error
	case lexer.Caret, lexer.Dash, lexer.Text:
		node.Value = string(token.Value)
		node.Type = TypeText
	}

	child, err := parse(name, l)
	if err != nil {
		return nil, err
	}

	if child != nil {
		node.Children = []*Node{
			child,
		}
	}

	node.Leaf = node.Children == nil
	if node.Leaf {
		node.Name = []string{name}
	}

	return node, nil
}

func Parse(name, input string) (*Node, error) {
	root := &Node{
		Value: "",
		Type:  TypeRoot,
		Leaf:  false,
	}

	if n, err := parse(name, lexer.New(input)); err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s", input)
	} else if n != nil {
		root.Children = []*Node{n}
	} else {
		root.Children = []*Node{
			{
				Children: nil,
				Value:    "",
				Type:     TypeText,
				Leaf:     true,
				Name:     []string{name},
			},
		}
	}

	root.compress()
	return root, nil
}

func Merge(root1, root2 *Node) *Node {
	return root1.merge(root2)
}
