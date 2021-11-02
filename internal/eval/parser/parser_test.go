package parser_test

import (
	"testing"

	"github.com/ythosa/bendy/internal/eval/ast"
	"github.com/ythosa/bendy/internal/eval/lexer"
	"github.com/ythosa/bendy/internal/eval/parser"
)

func TestStringLiteralExpression(t *testing.T) {
	t.Parallel()

	input := `"hello";`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseRequest()
	checkParserErrors(t, p)

	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)

	literal, ok := stmt.Expression.(*ast.WordLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.StringLiteral. got=%T",
			stmt.Expression)
	}

	if literal.Value != "hello" {
		t.Errorf("literal.Value not %q. got=%q",
			"hello tanyushka", literal.Value)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	t.Parallel()

	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{`!"hello";`, "!", "hello"},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseRequest()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	t.Parallel()

	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{`"lol" & "kek"`, "lol", "&", "kek"},
		{`"lol" | "kek"`, "lol", "|", "kek"},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseRequest()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue,
			tt.operator, tt.rightValue) {
			return
		}
	}
}

func testInfixExpression(t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {
	t.Helper()

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)

		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)

		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	t.Helper()

	switch v := expected.(type) {
	case string:
		return testWordLiteral(t, exp, v)
	default:
		t.Errorf("type of exp not handled. got=%T", exp)

		return false
	}
}

func testWordLiteral(t *testing.T, exp ast.Expression, value string) bool {
	t.Helper()

	ident, ok := exp.(*ast.WordLiteral)

	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)

		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)

		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())

		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	t.Helper()

	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{
			`!"a" & "b"`,
			`((!a) & b)`,
		},
		{
			`!!"a"`,
			"(!(!a))",
		},
		{
			`"a" & "b" & "c"`,
			"((a & b) & c)",
		},
		{
			`"a" | "b" & "c"`,
			"(a | (b & c))",
		},
		{
			`"a" | "b" & !"c"`,
			"(a | (b & (!c)))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseRequest()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
