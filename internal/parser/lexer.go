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
	LexerEOF lexerTokenType = iota
	LexerWildcard
	LexerText
)

type lexer struct {
	source   *scanner.Scanner
	finished bool
	current  *token
}

func NewLexer(source string) *lexer {
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
	case LexerWildcard:
		for getTokenType(l.source.Peek()) == LexerWildcard {
			l.source.Next()
		}
		l.current = &token{
			value: string(r),
			kind:  LexerWildcard,
		}

	case LexerText:
		var value bytes.Buffer
		value.WriteRune(r)

		for getTokenType(l.source.Peek()) == LexerText {
			value.WriteRune(l.source.Next())
		}

		l.current = &token{
			value: value.String(),
			kind:  LexerText,
		}

	default:
		fallthrough
	case LexerEOF:
		l.finished = true
		l.current = &token{
			value: "",
			kind:  LexerEOF,
		}
	}

	return !l.finished
}

func getTokenType(r rune) lexerTokenType {
	switch r {
	case wildcardRune:
		return LexerWildcard
	case scanner.EOF:
		return LexerEOF
	default:
		return LexerText
	}
}

type token struct {
	value string
	kind  lexerTokenType
}
