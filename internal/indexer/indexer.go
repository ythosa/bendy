package indexer

import (
	"container/list"

	"github.com/ythosa/bendy/internal/decoder"
	"github.com/ythosa/bendy/internal/normalizer"
)

type Indexer struct {
	decoder    decoder.Decoder
	normalizer normalizer.Normalizer
}

func NewIndexer(decoder decoder.Decoder, normalizer normalizer.Normalizer) *Indexer {
	return &Indexer{
		decoder:    decoder,
		normalizer: normalizer,
	}
}

func MergeInvertIndexes(invertIndexes ...InvertIndex) InvertIndex {
	result := make(InvertIndex)

	for _, index := range invertIndexes {
		for term, docs := range index {
			if _, ok := result[term]; ok {
				result[term] = MergeSortedListsDistrict(result[term], docs)
			} else {
				result[term] = docs
			}
		}
	}

	return result
}

func MergeSortedListsDistrict(l1 *list.List, l2 *list.List) *list.List {
	result := list.New()
	used := make(map[int]bool)
	e1 := l1.Front()
	e2 := l2.Front()

	for e1 != nil && e2 != nil {
		value1, _ := e1.Value.(int)
		value2, _ := e2.Value.(int)

		var valueToPush int

		if value1 < value2 {
			valueToPush = value1
			e1 = e1.Next()
		} else {
			valueToPush = value2
			e2 = e2.Next()
		}

		if _, ok := used[valueToPush]; !ok {
			result.PushBack(valueToPush)

			used[valueToPush] = true
		}
	}

	for e1 != nil {
		value, _ := e1.Value.(int)
		if _, ok := used[value]; !ok {
			result.PushBack(value)

			used[value] = true
		}

		e1 = e1.Next()
	}

	for e2 != nil {
		value, _ := e2.Value.(int)
		if _, ok := used[value]; !ok {
			result.PushBack(value)

			used[value] = true
		}

		e2 = e2.Next()
	}

	return result
}
