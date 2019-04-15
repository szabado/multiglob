package lexer

import (
	"strings"
	"text/scanner"
)

const (
	asteriskRune     = '*'
	openBracketRune  = '['
	closeBracketRune = ']'
	escapeRune       = '\\'
	dashRune         = '-'
	caretRune        = '^'
)

type LexerTokenType int

const (
	eof LexerTokenType = iota
	Asterisk
	Text
	Bracket
	Backslash
	Dash
	Caret
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
	case Asterisk:
		for getTokenType(l.source.Peek()) == Asterisk {
			l.source.Next()
		}
		l.current = &Token{
			Value: string(r),
			Type:  Asterisk,
		}

	case Text:
		l.current = &Token{
			Value: string(r),
			Type:  Text,
		}
	case Bracket, Backslash, Dash, Caret:
		l.current = &Token{
			Value: string(r),
			Type:  getTokenType(r),
		}

	default:
		fallthrough
	case eof:
		l.finished = true
		l.current = &Token{
			Value: "",
			Type:  eof,
		}
	}

	return !l.finished
}

// TODO: Test Peek
func (l *Lexer) Peek() (r rune, ok bool) {
	r = l.source.Peek()
	if getTokenType(r) == eof {
		return rune(0), false
	}
	return r, true
}

func getTokenType(r rune) LexerTokenType {
	switch r {
	case caretRune:
		return Caret
	case asteriskRune:
		return Asterisk
	case scanner.EOF:
		return eof
	case openBracketRune, closeBracketRune:
		return Bracket
	case escapeRune:
		return Backslash
	case dashRune:
		return Dash
	default:
		return Text
	}
}

type Token struct {
	Value string
	Type  LexerTokenType
}
