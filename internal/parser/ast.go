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
	if n == nil && n2 == nil {
		panic("wut")
	}

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
	default:
		panic("wut")
	case TypeText:
		return n.Value == n2.Value
	}
	return true
}

func (n *Node) merge(n2 *Node) *Node {
	if n == nil && n2 == nil {
		panic("wut")
	}

	if n == nil {
		return n2
	} else if n2 == nil {
		return n
	}

	if n.Type != n2.Type {
		// TODO: remove this
		panic("wut")
	}

	children := make([]*Node, len(n.Children))

	fmt.Printf("%#v\n", n.Children)
	fmt.Printf("%#v\n", n2.Children)
	fmt.Println(n.Children)
	fmt.Println(n2.Children)
	copy(children, n.Children)
	children = append(children, n2.Children...)

	for i := 0; i < len(children); i++ {
		for j := i + 1; j < len(children); j++ {
			child1, child2 := children[i], children[j]

			if !child1.canMerge(child2) {
				continue
			}
			fmt.Println("changing array", children)

			if j+1 >= len(children) {
				children = children[:j]
			} else {
				children = append(children[:j], children[j+1:]...)
			}

			fmt.Println("changed array", children)
			fmt.Println("Merging ", child1, child2)
			children[i] = child1.merge(child2)
			fmt.Println("Post merge: ", children[i])
		}
	}
	fmt.Println("final: ", children)

	return &Node{
		Children: children,
		Type:     n.Type,
		Value:    n.Value,
	}
}

func parse(l *lexer) *Node {
	if !l.Next() {
		return nil
	}

	token := l.Scan()

	child := parse(l)

	children := []*Node{
		child,
	}

	if child == nil {
		children = nil
	}

	return &Node{
		Children: children,
		Type:     getNodeType(token.kind),
		Value:    token.value,
	}
}

func getNodeType(tokenType lexerTokenType) NodeType {
	switch tokenType {
	case wildcard:
		return TypeAny
	case text:
		return TypeText
	case eof:
		fallthrough
	default:
		panic("this should never happen")
	}
}

func Parse(input string) *Node {
	root := newRootNode(nil)

	if n := parse(newLexer(input)); n != nil {
		root.Children = []*Node{n}
	} else {
		root.Children = []*Node{
			{
				Children: nil,
				Value:    "",
				Type:     TypeText,
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
	}
}

func Merge(root1, root2 *Node) *Node {
	return root1.merge(root2)
}
