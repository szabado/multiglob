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
	plusRune         = '+'
)

// TokenType enumerates the possible token types returned by the Lexer. Any unexported types
// never exit the package.
type TokenType int

const (
	eof TokenType = iota
	Asterisk
	Text
	Bracket
	Backslash
	Dash
	Caret
	Plus
)

// Lexer is a tokenizer that returns individual runes along with their associated types.
type Lexer struct {
	source   *scanner.Scanner
	finished bool
	current  *Token
}

// New returns a new Lexer that wraps the given source string.
func New(source string) *Lexer {
	var sc scanner.Scanner

	l := &Lexer{
		source:   sc.Init(strings.NewReader(source)),
		finished: false,
	}

	return l
}

// Scan returns the current token.
func (l *Lexer) Scan() *Token {
	return &*l.current
}

// Next advances the lexer to the next token, and discards the current one. It must be called before
// any calls to Scan.
func (l *Lexer) Next() bool {
	r := l.source.Next()
	switch t := getTokenType(r); t {
	case Asterisk:
		for getTokenType(l.source.Peek()) == Asterisk {
			l.source.Next()
		}
		l.current = &Token{
			Value: string(r),
			Type:  Asterisk,
		}

	case Bracket, Backslash, Caret, Dash, Plus, Text:
		l.current = &Token{
			Value: string(r),
			Type:  t,
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

// Peek returns the next token without consuming the current one. If the current token
// is the last token, Peek returns nil. It can be called before the first call to Next.
func (l *Lexer) Peek() (token *Token) {
	r := l.source.Peek()
	t := getTokenType(r)
	if t == eof {
		return nil
	}

	return &Token{
		Value: string(r),
		Type:  t,
	}
}

func getTokenType(r rune) TokenType {
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
	case plusRune:
		return Plus
	default:
		return Text
	}
}

type Token struct {
	Value string
	Type  TokenType
}
