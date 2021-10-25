package utils

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ListValue uint32

func ListToSlice(l *list.List) []ListValue {
	result := make([]ListValue, 0)
	for e := l.Front(); e != nil; e = e.Next() {
		result = append(result, e.Value.(ListValue))
	}

	return result
}

func MapOnListsToMapOnSlices(m map[string]*list.List) map[string][]ListValue {
	result := make(map[string][]ListValue)
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
		expectedValue, _ := expectedElement.Value.(ListValue)
		actualValue, _ := actualElement.Value.(ListValue)
		assert.Equal(t, expectedValue, actualValue)

		expectedElement = expectedElement.Next()
		actualElement = actualElement.Next()
	}
}
