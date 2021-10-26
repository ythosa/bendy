package indexer_test

import (
	"container/list"
	"testing"

	"github.com/ythosa/bendy/internal/indexer"
)

func TestInsertWithKeepSorting(t *testing.T) {
	t.Parallel()

	type args struct {
		l     *list.List
		docID indexer.DocID
	}

	testCases := []struct {
		input    args
		expected *list.List
	}{
		{
			input: args{
				l:     indexer.SliceToList([]indexer.DocID{1, 2, 5, 6}),
				docID: 3,
			},
			expected: indexer.SliceToList([]indexer.DocID{1, 2, 3, 5, 6}),
		},
		{
			input: args{
				l:     indexer.SliceToList([]indexer.DocID{2}),
				docID: 1,
			},
			expected: indexer.SliceToList([]indexer.DocID{1, 2}),
		},
		{
			input: args{
				l:     indexer.SliceToList([]indexer.DocID{2}),
				docID: 3,
			},
			expected: indexer.SliceToList([]indexer.DocID{2, 3}),
		},
	}

	for _, tc := range testCases {
		indexer.InsertToListWithKeepSorting(tc.input.l, tc.input.docID)
		indexer.CompareLists(t, tc.expected, tc.input.l)
	}
}
