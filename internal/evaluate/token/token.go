package token

// Type is type of token
type Type string

// Token is type for token which contains type and literal
type Token struct {
	Type    Type
	Literal string
}

// Names of tokens
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// identifies
	WORD = "WORD"

	// operators
	AND = "&"
	OR  = "|"
	NOT = "!"

	// Delimiters
	LPAREN    = "("
	RPAREN    = ")"
	SEMICOLON = ";"
)
