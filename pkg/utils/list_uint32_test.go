package utils_test

import (
	"container/list"
	"testing"

	"github.com/ythosa/bendy/pkg/utils"
)

func TestInsertWithKeepSorting(t *testing.T) {
	t.Parallel()

	type args struct {
		l     *list.List
		docID utils.SliceValue
	}

	testCases := []struct {
		input    args
		expected *list.List
	}{
		{
			input: args{
				l:     utils.SliceToList([]utils.SliceValue{1, 2, 5, 6}),
				docID: 3,
			},
			expected: utils.SliceToList([]utils.SliceValue{1, 2, 3, 5, 6}),
		},
		{
			input: args{
				l:     utils.SliceToList([]utils.SliceValue{2}),
				docID: 1,
			},
			expected: utils.SliceToList([]utils.SliceValue{1, 2}),
		},
		{
			input: args{
				l:     utils.SliceToList([]utils.SliceValue{2}),
				docID: 3,
			},
			expected: utils.SliceToList([]utils.SliceValue{2, 3}),
		},
	}

	for _, tc := range testCases {
		utils.InsertToListWithKeepSorting(tc.input.l, tc.input.docID)
		utils.CompareLists(t, tc.expected, tc.input.l)
	}
}
