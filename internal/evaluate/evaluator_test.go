package evaluate_test

import (
	"testing"

	"github.com/ythosa/bendy/internal/evaluate"
	"github.com/ythosa/bendy/internal/evaluate/lexer"
	"github.com/ythosa/bendy/internal/evaluate/object"
	"github.com/ythosa/bendy/internal/evaluate/parser"
	"github.com/ythosa/bendy/internal/index"
)

func getTestEvaluatorObject() *evaluate.Evaluator {
	ii := make(index.InvertIndex)
	ii["kek"] = index.SliceToList([]index.DocID{1, 3, 5})
	ii["lol"] = index.SliceToList([]index.DocID{2, 3, 4})
	ii["puk"] = index.SliceToList([]index.DocID{6})

	e := evaluate.NewEvaluator(ii, index.SliceToList([]index.DocID{1, 2, 3, 4, 5, 6, 7}))

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
