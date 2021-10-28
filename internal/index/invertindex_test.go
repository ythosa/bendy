package index_test

import (
	"testing"

	"github.com/ythosa/bendy/internal/index"
)

func TestCap(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		l1       []index.DocID
		l2       []index.DocID
		expected []index.DocID
	}{
		{
			l1:       []index.DocID{1, 2, 3, 4, 5},
			l2:       []index.DocID{1, 2},
			expected: []index.DocID{1, 2},
		},
		{
			l1:       []index.DocID{1, 2},
			l2:       []index.DocID{1, 2, 3, 4, 5},
			expected: []index.DocID{1, 2},
		},
		{
			l1:       []index.DocID{},
			l2:       []index.DocID{1, 2, 3, 4, 5},
			expected: []index.DocID{},
		},
		{
			l1:       []index.DocID{1, 2, 3, 4, 5},
			l2:       []index.DocID{},
			expected: []index.DocID{},
		},
	}

	for _, tc := range testCases {
		r := index.Cap(index.SliceToList(tc.l1), index.SliceToList(tc.l2))
		index.CompareLists(t, index.SliceToList(tc.expected), r)
	}
}

func TestCup(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		l1       []index.DocID
		l2       []index.DocID
		expected []index.DocID
	}{
		{
			l1:       []index.DocID{1, 2, 3, 4, 5},
			l2:       []index.DocID{1, 2, 3, 4, 5, 6},
			expected: []index.DocID{1, 2, 3, 4, 5, 6},
		},
		{
			l1:       []index.DocID{1, 2, 3, 4, 5, 6},
			l2:       []index.DocID{1, 2, 3, 4, 5},
			expected: []index.DocID{1, 2, 3, 4, 5, 6},
		},
		{
			l1:       []index.DocID{1, 2, 3, 4, 5},
			l2:       []index.DocID{},
			expected: []index.DocID{1, 2, 3, 4, 5},
		},
		{
			l1:       []index.DocID{},
			l2:       []index.DocID{1, 2, 3, 4, 5},
			expected: []index.DocID{1, 2, 3, 4, 5},
		},
	}

	for _, tc := range testCases {
		r := index.Cup(index.SliceToList(tc.l1), index.SliceToList(tc.l2))
		index.CompareLists(t, index.SliceToList(tc.expected), r)
	}
}

func TestInvert(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		l        []index.DocID
		all      []index.DocID
		expected []index.DocID
	}{
		{
			l:        []index.DocID{1, 2, 3, 4, 5},
			all:      []index.DocID{1, 2, 3, 4, 5, 6},
			expected: []index.DocID{6},
		},
		{
			l:        []index.DocID{1},
			all:      []index.DocID{1, 2, 3, 4, 5},
			expected: []index.DocID{2, 3, 4, 5},
		},
		{
			l:        []index.DocID{1, 2, 3, 4, 5},
			all:      []index.DocID{},
			expected: []index.DocID{},
		},
		{
			l:        []index.DocID{},
			all:      []index.DocID{1, 2, 3, 4, 5},
			expected: []index.DocID{1, 2, 3, 4, 5},
		},
	}

	for _, tc := range testCases {
		r := index.Invert(index.SliceToList(tc.l), index.SliceToList(tc.all))
		index.CompareLists(t, index.SliceToList(tc.expected), r)
	}
}
