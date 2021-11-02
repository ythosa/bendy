package eval_test

import (
	"testing"

	"github.com/ythosa/bendy/internal/eval"
	"github.com/ythosa/bendy/internal/eval/lexer"
	"github.com/ythosa/bendy/internal/eval/object"
	"github.com/ythosa/bendy/internal/eval/parser"
	"github.com/ythosa/bendy/internal/index"
)

func getTestEvaluatorObject() *eval.Evaluator {
	ii := make(index.InvertIndex)
	ii["kek"] = index.NewIndex(index.SliceToList([]index.DocID{1, 3, 5}))
	ii["lol"] = index.NewIndex(index.SliceToList([]index.DocID{2, 3, 4}))
	ii["puk"] = index.NewIndex(index.SliceToList([]index.DocID{6}))

	e := eval.NewEvaluator(ii, index.SliceToList([]index.DocID{1, 2, 3, 4, 5, 6, 7}))

	return e
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseRequest()

	return getTestEvaluatorObject().Eval(program)
}

func TestEvaluator_Eval(t *testing.T) {
	testCases := []struct {
		input    string
		expected []index.DocID
	}{
		{
			input:    `"kek" & "lol"`,
			expected: []index.DocID{3},
		},
		{
			input:    `"kek" | "lol"`,
			expected: []index.DocID{1, 2, 3, 4, 5},
		},
		{
			input:    `"kek" & !"lol"`,
			expected: []index.DocID{1, 5},
		},
		{
			input:    `"kek" | !"puk"`,
			expected: []index.DocID{1, 2, 3, 4, 5, 7},
		},
		{
			input:    `!"kek" & ("lol" | "puk")`,
			expected: []index.DocID{2, 4, 6},
		},
		{
			input:    `"kek" & "ya kto"`,
			expected: []index.DocID{},
		},
	}

	for _, tc := range testCases {
		evaluated := testEval(tc.input)
		index.CompareLists(t, index.SliceToList(tc.expected), evaluated.(*object.DocIDs).Value)
	}
}
