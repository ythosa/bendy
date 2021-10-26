package indexer_test

import (
	"testing"

	"github.com/ythosa/bendy/internal/indexer"
)

func TestCap(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		l1       []indexer.DocID
		l2       []indexer.DocID
		expected []indexer.DocID
	}{
		{
			l1:       []indexer.DocID{1, 2, 3, 4, 5},
			l2:       []indexer.DocID{1, 2},
			expected: []indexer.DocID{1, 2},
		},
		{
			l1:       []indexer.DocID{1, 2},
			l2:       []indexer.DocID{1, 2, 3, 4, 5},
			expected: []indexer.DocID{1, 2},
		},
		{
			l1:       []indexer.DocID{},
			l2:       []indexer.DocID{1, 2, 3, 4, 5},
			expected: []indexer.DocID{},
		},
		{
			l1:       []indexer.DocID{1, 2, 3, 4, 5},
			l2:       []indexer.DocID{},
			expected: []indexer.DocID{},
		},
	}

	for _, tc := range testCases {
		r := indexer.Cap(indexer.SliceToList(tc.l1), indexer.SliceToList(tc.l2))
		indexer.CompareLists(t, indexer.SliceToList(tc.expected), r)
	}
}

func TestCup(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		l1       []indexer.DocID
		l2       []indexer.DocID
		expected []indexer.DocID
	}{
		{
			l1:       []indexer.DocID{1, 2, 3, 4, 5},
			l2:       []indexer.DocID{1, 2, 3, 4, 5, 6},
			expected: []indexer.DocID{1, 2, 3, 4, 5, 6},
		},
		{
			l1:       []indexer.DocID{1, 2, 3, 4, 5, 6},
			l2:       []indexer.DocID{1, 2, 3, 4, 5},
			expected: []indexer.DocID{1, 2, 3, 4, 5, 6},
		},
		{
			l1:       []indexer.DocID{1, 2, 3, 4, 5},
			l2:       []indexer.DocID{},
			expected: []indexer.DocID{1, 2, 3, 4, 5},
		},
		{
			l1:       []indexer.DocID{},
			l2:       []indexer.DocID{1, 2, 3, 4, 5},
			expected: []indexer.DocID{1, 2, 3, 4, 5},
		},
	}

	for _, tc := range testCases {
		r := indexer.Cup(indexer.SliceToList(tc.l1), indexer.SliceToList(tc.l2))
		indexer.CompareLists(t, indexer.SliceToList(tc.expected), r)
	}
}

func TestInvert(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		l        []indexer.DocID
		all      []indexer.DocID
		expected []indexer.DocID
	}{
		{
			l:        []indexer.DocID{1, 2, 3, 4, 5},
			all:      []indexer.DocID{1, 2, 3, 4, 5, 6},
			expected: []indexer.DocID{6},
		},
		{
			l:        []indexer.DocID{1},
			all:      []indexer.DocID{1, 2, 3, 4, 5},
			expected: []indexer.DocID{2, 3, 4, 5},
		},
		{
			l:        []indexer.DocID{1, 2, 3, 4, 5},
			all:      []indexer.DocID{},
			expected: []indexer.DocID{},
		},
		{
			l:        []indexer.DocID{},
			all:      []indexer.DocID{1, 2, 3, 4, 5},
			expected: []indexer.DocID{1, 2, 3, 4, 5},
		},
	}

	for _, tc := range testCases {
		r := indexer.Invert(indexer.SliceToList(tc.l), indexer.SliceToList(tc.all))
		indexer.CompareLists(t, indexer.SliceToList(tc.expected), r)
	}
}
