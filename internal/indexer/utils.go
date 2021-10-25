package indexer

import (
	"container/list"
	"fmt"
	"os"
)

func insertToListWithKeepSorting(l *list.List, docID DocID) {
	var previousElement *list.Element

	for currentElement := l.Front(); currentElement != nil; currentElement = currentElement.Next() {
		if value, _ := currentElement.Value.(DocID); value > docID {
			break
		}

		previousElement = currentElement
	}

	if previousElement != nil {
		l.InsertAfter(docID, previousElement)
	} else {
		l.InsertBefore(docID, l.Front())
	}
}

func checkIsFilePathsValid(filePaths []string) error {
	for _, fp := range filePaths {
		if _, err := os.Stat(fp); err != nil {
			return fmt.Errorf("error while getting file stat: %w", err)
		}
	}

	return nil
}

func sliceToList(slice []DocID) *list.List {
	l := list.New()
	for _, v := range slice {
		l.PushBack(v)
	}

	return l
}

func MapIndexOnSlicesToIndexOnLists(m map[string][]DocID) map[string]*list.List {
	result := make(map[string]*list.List)
	for k, v := range m {
		result[k] = sliceToList(v)
	}

	return result
}

func listToSlice(l *list.List) []DocID {
	result := make([]DocID, 0)
	for e := l.Front(); e != nil; e = e.Next() {
		result = append(result, e.Value.(DocID))
	}

	return result
}

func MapIndexOnListsToIndexOnSlices(m map[string]*list.List) map[string][]DocID {
	result := make(map[string][]DocID)
	for k, v := range m {
		result[k] = listToSlice(v)
	}

	return result
}
