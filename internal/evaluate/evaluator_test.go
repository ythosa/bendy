package evaluate_test

import (
	"testing"

	"github.com/ythosa/bendy/internal/evaluate"
	"github.com/ythosa/bendy/internal/evaluate/lexer"
	"github.com/ythosa/bendy/internal/evaluate/object"
	"github.com/ythosa/bendy/internal/evaluate/parser"
	"github.com/ythosa/bendy/internal/indexer"
)

func getTestEvaluatorObject() *evaluate.Evaluator {
	ii := make(indexer.InvertIndex)
	ii["kek"] = indexer.SliceToList([]indexer.DocID{1, 3, 5})
	ii["lol"] = indexer.SliceToList([]indexer.DocID{2, 3, 4})
	ii["puk"] = indexer.SliceToList([]indexer.DocID{6})

	e := evaluate.NewEvaluator(ii, indexer.SliceToList([]indexer.DocID{1, 2, 3, 4, 5, 6, 7}))

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
		expected []indexer.DocID
	}{
		{
			input:    `"kek" & "lol"`,
			expected: []indexer.DocID{3},
		},
		{
			input:    `"kek" | "lol"`,
			expected: []indexer.DocID{1, 2, 3, 4, 5},
		},
		{
			input:    `"kek" & !"lol"`,
			expected: []indexer.DocID{1, 5},
		},
		{
			input:    `"kek" | !"puk"`,
			expected: []indexer.DocID{1, 2, 3, 4, 5, 7},
		},
		{
			input:    `!"kek" & ("lol" | "puk")`,
			expected: []indexer.DocID{2, 4, 6},
		},
		{
			input:    `"kek" & "ya kto"`,
			expected: []indexer.DocID{},
		},
	}

	for _, tc := range testCases {
		evaluated := testEval(tc.input)
		indexer.CompareLists(t, indexer.SliceToList(tc.expected), evaluated.(*object.DocIDs).Value)
	}
}
