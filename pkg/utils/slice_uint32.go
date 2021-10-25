package utils

import "container/list"

type SliceValue uint32

func InsertToListWithKeepSorting(l *list.List, v SliceValue) {
	var previousElement *list.Element

	for currentElement := l.Front(); currentElement != nil; currentElement = currentElement.Next() {
		if value, _ := currentElement.Value.(SliceValue); value > v {
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

func SliceToList(slice []SliceValue) *list.List {
	l := list.New()
	for _, v := range slice {
		l.PushBack(v)
	}

	return l
}

func MapOnSlicesToMapOnLists(m map[string][]SliceValue) map[string]*list.List {
	result := make(map[string]*list.List)
	for k, v := range m {
		result[k] = SliceToList(v)
	}

	return result
}
