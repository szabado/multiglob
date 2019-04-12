package lexer

import (
	"bytes"
	"strings"
	"text/scanner"
)

const (
	wildcardRune = '*'
)

type LexerTokenType int

const (
	EOF LexerTokenType = iota
	Wildcard
	Text
)

type Lexer struct {
	source   *scanner.Scanner
	finished bool
	current  *Token
}

func New(source string) *Lexer {
	var sc scanner.Scanner

	l := &Lexer{
		source:   sc.Init(strings.NewReader(source)),
		finished: false,
	}

	return l
}

func (l *Lexer) Scan() *Token {
	return &*l.current
}

func (l *Lexer) Next() bool {
	switch r := l.source.Next(); getTokenType(r) {
	case Wildcard:
		for getTokenType(l.source.Peek()) == Wildcard {
			l.source.Next()
		}
		l.current = &Token{
			Value: string(r),
			Kind:  Wildcard,
		}

	case Text:
		var value bytes.Buffer
		value.WriteRune(r)

		for getTokenType(l.source.Peek()) == Text {
			value.WriteRune(l.source.Next())
		}

		l.current = &Token{
			Value: value.String(),
			Kind:  Text,
		}

	default:
		fallthrough
	case EOF:
		l.finished = true
		l.current = &Token{
			Value: "",
			Kind:  EOF,
		}
	}

	return !l.finished
}

func getTokenType(r rune) LexerTokenType {
	switch r {
	case wildcardRune:
		return Wildcard
	case scanner.EOF:
		return EOF
	default:
		return Text
	}
}

type Token struct {
	Value string
	Kind  LexerTokenType
}
