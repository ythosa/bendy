package index_test

import (
	"container/list"
	"testing"

	"github.com/ythosa/bendy/internal/index"
)

func TestInsertWithKeepSorting(t *testing.T) {
	t.Parallel()

	type args struct {
		l     *list.List
		docID index.DocID
	}

	testCases := []struct {
		input    args
		expected *list.List
	}{
		{
			input: args{
				l:     index.SliceToList([]index.DocID{1, 2, 5, 6}),
				docID: 3,
			},
			expected: index.SliceToList([]index.DocID{1, 2, 3, 5, 6}),
		},
		{
			input: args{
				l:     index.SliceToList([]index.DocID{2}),
				docID: 1,
			},
			expected: index.SliceToList([]index.DocID{1, 2}),
		},
		{
			input: args{
				l:     index.SliceToList([]index.DocID{2}),
				docID: 3,
			},
			expected: index.SliceToList([]index.DocID{2, 3}),
		},
	}

	for _, tc := range testCases {
		index.InsertToListWithKeepSorting(tc.input.l, tc.input.docID)
		index.CompareLists(t, tc.expected, tc.input.l)
	}
}
