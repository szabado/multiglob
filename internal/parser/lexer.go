package parser

import (
	"bytes"
	"strings"
	"text/scanner"
)

const (
	wildcardRune = '*'
)

type lexerTokenType int

const (
	eof lexerTokenType = iota
	wildcard
	text
)

type lexer struct {
	source   *scanner.Scanner
	finished bool
	current  *token
}

func newLexer(source string) *lexer {
	var sc scanner.Scanner

	l := &lexer{
		source:   sc.Init(strings.NewReader(source)),
		finished: false,
	}

	return l
}

func (l *lexer) Scan() *token {
	return &*l.current
}

func (l *lexer) Next() bool {
	switch r := l.source.Next(); getTokenType(r) {
	case wildcard:
		for getTokenType(l.source.Peek()) == wildcard {
			l.source.Next()
		}
		l.current = &token{
			value: string(r),
			kind:  wildcard,
		}

	case text:
		var value bytes.Buffer
		value.WriteRune(r)

		for getTokenType(l.source.Peek()) == text {
			value.WriteRune(l.source.Next())
		}

		l.current = &token{
			value: value.String(),
			kind:  text,
		}

	default:
		fallthrough
	case eof:
		l.finished = true
		l.current = &token{
			value: "",
			kind:  eof,
		}
	}

	return !l.finished
}

func getTokenType(r rune) lexerTokenType {
	switch r {
	case wildcardRune:
		return wildcard
	case scanner.EOF:
		return eof
	default:
		return text
	}
}

type token struct {
	value string
	kind  lexerTokenType
}
