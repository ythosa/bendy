package indexer_test

import (
	"bytes"
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ythosa/bendy/internal/decoder"
	"github.com/ythosa/bendy/internal/indexer"
	"github.com/ythosa/bendy/internal/normalizer"
)

func TestNewIndexer(t *testing.T) {
	t.Parallel()

	dec := decoder.NewTXTDecoder(bytes.NewBufferString(""))
	norm := normalizer.NewEnglishNormalizer()
	ix := indexer.NewIndexer(dec, norm)

	assert.NotNil(t, ix)
}

func TestMergeInvertIndexes(t *testing.T) {
	t.Parallel()

	ii1 := make(indexer.InvertIndex)
	ii1["kek"] = sliceToList([]int{1, 2, 3})
	ii1["lol"] = sliceToList([]int{5, 6})

	ii2 := make(indexer.InvertIndex)
	ii2["kek"] = sliceToList([]int{1, 2, 4, 10})

	expected := make(indexer.InvertIndex)
	expected["kek"] = sliceToList([]int{1, 2, 3, 4, 10})
	expected["lol"] = sliceToList([]int{5, 6})

	result := indexer.MergeInvertIndexes(ii1, ii2)
	compareLists(t, expected["kek"], result["kek"])
	compareLists(t, expected["lol"], result["lol"])
}

func TestMergeSortedListsDistrict(t *testing.T) {
	t.Parallel()

	type args struct {
		l1 *list.List
		l2 *list.List
	}

	testCases := []struct {
		input    args
		expected *list.List
	}{
		{
			input: args{
				sliceToList([]int{1, 3, 4, 5}),
				sliceToList([]int{2, 4}),
			},
			expected: sliceToList([]int{1, 2, 3, 4, 5}),
		},
		{
			input: args{
				sliceToList([]int{2, 4}),
				sliceToList([]int{1, 3, 5}),
			},
			expected: sliceToList([]int{1, 2, 3, 4, 5}),
		},
		{
			input: args{
				sliceToList([]int{}),
				sliceToList([]int{1, 3, 5}),
			},
			expected: sliceToList([]int{1, 3, 5}),
		},
		{
			input: args{
				sliceToList([]int{}),
				sliceToList([]int{}),
			},
			expected: sliceToList([]int{}),
		},
	}

	for _, tc := range testCases {
		result := indexer.MergeSortedListsDistrict(tc.input.l1, tc.input.l2)
		compareLists(t, tc.expected, result)
	}
}

func sliceToList(slice []int) *list.List {
	l := list.New()
	for _, v := range slice {
		l.PushBack(v)
	}

	return l
}

func compareLists(t *testing.T, expected *list.List, actual *list.List) {
	t.Helper()

	assert.Equal(t, expected.Len(), actual.Len(), "lists are different sizes")

	expectedElement := expected.Front()
	actualElement := actual.Front()

	for expectedElement != nil {
		expectedValue, _ := expectedElement.Value.(int)
		actualValue, _ := actualElement.Value.(int)
		assert.Equal(t, expectedValue, actualValue)

		expectedElement = expectedElement.Next()
		actualElement = actualElement.Next()
	}
}
