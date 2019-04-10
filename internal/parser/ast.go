package parser

import "fmt"

type NodeType int

//go:generate stringer -type=NodeType
const (
	TypeNothing NodeType = iota
	TypeAny
	TypeText
)

type Node struct {
	Type     NodeType
	Value    string
	Children []*Node
}

func newNode(t NodeType) *Node {
	switch t {
	case TypeAny:
		fallthrough
	case TypeText:
		return &Node{
			Type: t,
		}

	case TypeNothing:
		fallthrough
	default:
		return nil
	}
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
		return TypeNothing
	}
}

func Parse(input string) *Node {
	n := parse(newLexer(input))
	if n == nil {
		n = &Node{
			Value:    "",
			Type:     TypeText,
			Children: nil,
		}
	}
	return n
}
