package parser

import (
	"bytes"
	"fmt"
)

type NodeType int

//go:generate stringer -type=NodeType
const (
	TypeRoot NodeType = iota
	TypeAny
	TypeText
)

type Node struct {
	Type     NodeType
	Value    string
	Children []*Node
	Leaf     bool
	Name     []string // Only valid on leaf nodes. List of names of patterns terminate that on this leaf node
}

func (n *Node) String() string {
	if l := len(n.Children); l == 0 {
		return fmt.Sprintf("%s", n.Value)
	} else if l == 1 {
		return fmt.Sprintf("%s%s", n.Value, n.Children[0])
	}

	b := bytes.Buffer{}

	b.WriteString(n.Value)
	b.WriteString("(")
	for i, child := range n.Children {
		b.WriteString(child.String())
		if i+1 < len(n.Children) {
			b.WriteString("|")
		}
	}

	b.WriteString(")")

	return b.String()
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
	}
	return true
}

func (n *Node) merge(n2 *Node) *Node {
	if n == nil {
		return n2
	} else if n2 == nil {
		return n
	}

	children := make([]*Node, len(n.Children))

	copy(children, n.Children)
	children = append(children, n2.Children...)

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

func mergeNames(n1, n2 *Node) []string {
	if n1.Leaf && n2.Leaf {
		return append(n1.Name, n2.Name...)
	} else if n1.Leaf {
		return n1.Name
	} else {
		return n2.Name
	}
}

func parse(name string, l *lexer) *Node {
	if !l.Next() {
		return nil
	}

	token := l.Scan()

	child := parse(name, l)

	children := []*Node{
		child,
	}

	if child == nil {
		children = nil
	}

	leaf := children == nil
	var nameSl []string
	if leaf {
		nameSl = []string{name}
	}

	return &Node{
		Children: children,
		Type:     getNodeType(token.kind),
		Value:    token.value,
		Leaf:     leaf,
		Name:     nameSl,
	}
}

func getNodeType(tokenType lexerTokenType) NodeType {
	switch tokenType {
	case LexerWildcard:
		return TypeAny
	case LexerText:
		return TypeText
	default:
		return NodeType(-1)
	}
}

func Parse(name, input string) *Node {
	root := newRootNode(nil)

	if n := parse(name, NewLexer(input)); n != nil {
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
	return root
}

func newRootNode(children []*Node) *Node {
	return &Node{
		Value:    "",
		Type:     TypeRoot,
		Children: children,
		Leaf:     false,
	}
}

func Merge(root1, root2 *Node) *Node {
	return root1.merge(root2)
}
