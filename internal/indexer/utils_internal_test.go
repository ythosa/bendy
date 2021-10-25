package indexer

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertWithKeepSorting(t *testing.T) {
	t.Parallel()

	type args struct {
		l     *list.List
		docID DocID
	}

	testCases := []struct {
		input    args
		expected *list.List
	}{
		{
			input: args{
				l:     sliceToList([]DocID{1, 2, 5, 6}),
				docID: 3,
			},
			expected: sliceToList([]DocID{1, 2, 3, 5, 6}),
		},
		{
			input: args{
				l:     sliceToList([]DocID{2}),
				docID: 1,
			},
			expected: sliceToList([]DocID{1, 2}),
		},
		{
			input: args{
				l:     sliceToList([]DocID{2}),
				docID: 3,
			},
			expected: sliceToList([]DocID{2, 3}),
		},
	}

	for _, tc := range testCases {
		insertToListWithKeepSorting(tc.input.l, tc.input.docID)
		compareLists(t, tc.expected, tc.input.l)
	}
}

func TestCheckIsFilePathsValid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		filePaths     []string
		expectedError bool
	}{
		{
			filePaths:     []string{"./indexer.go"},
			expectedError: false,
		},
		{
			filePaths:     []string{"./invalid_file_path"},
			expectedError: true,
		},
		{
			filePaths:     []string{"./indexer", "./invalid_file_path"},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		if tc.expectedError {
			assert.NotNil(t, checkIsFilePathsValid(tc.filePaths))
		} else {
			assert.Nil(t, checkIsFilePathsValid(tc.filePaths))
		}
	}
}
