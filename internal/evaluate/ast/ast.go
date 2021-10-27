package ast

import (
	"bytes"

	"github.com/ythosa/bendy/internal/evaluate/token"
)

// Node is interface for simple node (anyone) element in the AST tree.
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement is interface for statements elements in the AST tree.
type Statement interface {
	Node
	statementNode()
}

// Expression is interface for expressions elements in the AST tree.
type Expression interface {
	Node
	expressionNode()
}

// Request is type for program - higher element of AST tree.
type Request struct {
	Statements []Statement
}

// TokenLiteral returns token literal of the node.
func (p *Request) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

// String returns string representation of the node.
func (p *Request) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// ExpressionStatement is type for expression statements in the AST tree.
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral returns token literal of the node.
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// String returns string representation of the node.
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

// PrefixExpression is type for prefix expressions in the AST tree.
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral returns token literal of the node.
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

// String returns string representation of the node.
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression is type for infix expressions in the AST tree.
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral returns token literal of the node.
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String returns string representation of the node.
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// WordLiteral is type for word literals.
type WordLiteral struct {
	Token token.Token
	Value string
}

func (sl *WordLiteral) expressionNode() {}

// TokenLiteral returns token literal of the node.
func (sl *WordLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

// String returns string representation of the node.
func (sl *WordLiteral) String() string {
	return sl.Token.Literal
}
