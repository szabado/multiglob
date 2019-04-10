package parser

import "fmt"

type NodeType int

//go:generate stringer -type=NodeType
const (
	TypeRoot NodeType = iota
	TypeAny
	TypeText
	TypeLeaf
)

type Node struct {
	Type     NodeType
	Value    string
	Children []*Node
}

func (n *Node) String() string {
	// TODO: Fix this
	return fmt.Sprintf("%s%s", n.Value, n.Children)
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
		return TypeLeaf
	}
}

func Parse(input string) *Node {
	root := &Node{
		Value:    "",
		Type:     TypeRoot,
		Children: nil,
	}

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
