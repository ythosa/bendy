package lexer

import "github.com/ythosa/bendy/internal/eval/token"

// Lexer is type for lexer which turns code into a sequence of tokens.
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// New returns new lexer.
func New(input string) *Lexer {
	l := Lexer{input: input}
	l.readChar()

	return &l
}

// NextToken returns next token of the code.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '!':
		tok = newToken(token.NOT, l.ch)
	case '&':
		tok = newToken(token.AND, l.ch)
	case '|':
		tok = newToken(token.OR, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '"':
		tok.Type = token.WORD
		tok.Literal = l.readString()
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()

	return tok
}

func newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	var str []byte

	for {
		l.readChar()

		if l.ch == 0 {
			break
		}

		if l.ch == '"' {
			break
		}

		str = append(str, l.ch)
	}

	return string(str)
}
