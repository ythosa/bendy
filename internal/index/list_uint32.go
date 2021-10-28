package index

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ListToSlice(l *list.List) []DocID {
	result := make([]DocID, 0)
	for e := l.Front(); e != nil; e = e.Next() {
		result = append(result, e.Value.(DocID))
	}

	return result
}

func MapOnListsToMapOnSlices(m map[string]*list.List) map[string][]DocID {
	result := make(map[string][]DocID)
	for k, v := range m {
		result[k] = ListToSlice(v)
	}

	return result
}

func CompareLists(t *testing.T, expected *list.List, actual *list.List) {
	t.Helper()

	assert.Equal(t, expected.Len(), actual.Len(), "lists are different sizes")

	expectedElement := expected.Front()
	actualElement := actual.Front()

	for expectedElement != nil {
		expectedValue, _ := expectedElement.Value.(DocID)
		actualValue, _ := actualElement.Value.(DocID)
		assert.Equal(t, expectedValue, actualValue)

		expectedElement = expectedElement.Next()
		actualElement = actualElement.Next()
	}
}
