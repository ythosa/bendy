package evaluate

import (
	"container/list"
	"fmt"

	"github.com/ythosa/bendy/internal/evaluate/ast"
	"github.com/ythosa/bendy/internal/evaluate/object"
	"github.com/ythosa/bendy/internal/indexer"
)

type Evaluator struct {
	InvertIndex indexer.InvertIndex
	AllDocIDs   *list.List
}

// Eval evaluates node of AST tree
func (e *Evaluator) Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Request:
		return e.evalRequest(node)

	// Expressions
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression)

	case *ast.WordLiteral:
		return &object.DocIDs{
			Value: e.InvertIndex[node.Value],
		}

	case *ast.PrefixExpression:
		right := e.Eval(node.Right)
		if isError(right) {
			return right
		}

		return e.evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := e.Eval(node.Left)
		if isError(left) {
			return left
		}

		right := e.Eval(node.Right)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	}

	return nil
}

func (e *Evaluator) evalRequest(program *ast.Request) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = e.Eval(statement)
	}

	return result
}

func (e *Evaluator) evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return e.evalNotOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func (e *Evaluator) evalNotOperatorExpression(right object.Object) object.Object {
	switch right.Type() {
	case object.DocIDsObj:
		return &object.DocIDs{Value: indexer.Invert(right.(*object.DocIDs).Value, e.AllDocIDs)}
	default:
		return newError("unknown operator: ! %s", right.Type())
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.DocIDsObj && right.Type() == object.DocIDsObj:
		return evalInfixDocIDsExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalInfixDocIDsExpression(
	operator string,
	left object.Object,
	right object.Object,
) object.Object {
	leftVal := left.(*object.DocIDs).Value
	rightVal := right.(*object.DocIDs).Value

	switch operator {
	case "&":
		return &object.DocIDs{Value: indexer.Cap(leftVal, rightVal)}
	case "|":
		return &object.DocIDs{Value: indexer.Cup(leftVal, rightVal)}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorObj
	}

	return false
}
