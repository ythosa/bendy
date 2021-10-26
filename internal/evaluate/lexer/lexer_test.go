package lexer_test

import (
	"testing"

	"github.com/ythosa/bendy/internal/evaluate/lexer"
	"github.com/ythosa/bendy/internal/evaluate/token"
)

func TestNextToken(t *testing.T) {
	input := `"LOL" & "KEK" | ("BOO" & (!"AAA"))`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.WORD, "LOL"},
		{token.AND, "&"},
		{token.WORD, "KEK"},
		{token.OR, "|"},
		{token.LPAREN, "("},
		{token.WORD, "BOO"},
		{token.AND, "&"},
		{token.LPAREN, "("},
		{token.NOT, "!"},
		{token.WORD, "AAA"},
		{token.RPAREN, ")"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	}

	l := lexer.New(input)
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
