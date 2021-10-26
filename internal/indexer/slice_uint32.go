package indexer

import (
	"container/list"
)

func InsertToListWithKeepSorting(l *list.List, v DocID) {
	var previousElement *list.Element

	for currentElement := l.Front(); currentElement != nil; currentElement = currentElement.Next() {
		if value, _ := currentElement.Value.(DocID); value > v {
			break
		}

		previousElement = currentElement
	}

	if previousElement != nil {
		l.InsertAfter(v, previousElement)
	} else {
		l.InsertBefore(v, l.Front())
	}
}

func SliceToList(slice []DocID) *list.List {
	l := list.New()
	for _, v := range slice {
		l.PushBack(v)
	}

	return l
}

func MapOnSlicesToMapOnLists(m map[string][]DocID) map[string]*list.List {
	result := make(map[string]*list.List)
	for k, v := range m {
		result[k] = SliceToList(v)
	}

	return result
}
