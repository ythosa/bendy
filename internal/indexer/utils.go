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
