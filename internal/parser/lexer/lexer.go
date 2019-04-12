package lexer

import (
	"bytes"
	"strings"
	"text/scanner"
)

const (
	asteriskRune = '*'
	openBracketRune = '['
	closeBracketRune = ']'
	escapeRune = '\\'
	dashRune = '-'
)

type LexerTokenType int

const (
	EOF       LexerTokenType = iota // Tokens less than EOF should never exit this package
	Asterisk
	Text
	Bracket
	Backslash
	Dash
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
			Kind:  Asterisk,
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
	case Bracket, Backslash, Dash:
		l.current = &Token{
			Value: string(r),
			Kind: getTokenType(r),
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
	case asteriskRune:
		return Asterisk
	case scanner.EOF:
		return EOF
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
	Kind  LexerTokenType
}
