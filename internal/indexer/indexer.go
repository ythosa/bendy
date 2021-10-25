package indexer

import (
	"container/list"
	"context"
	"encoding/gob"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"github.com/ythosa/bendy/internal/decoder"
	"github.com/ythosa/bendy/internal/normalizer"
)

const maxOpenFilesCount = 1000

type DocID uint32

type Indexer struct {
	normalizer normalizer.Normalizer
}

func NewIndexer(normalizer normalizer.Normalizer) *Indexer {
	return &Indexer{
		normalizer: normalizer,
	}
}

func (i *Indexer) IndexFiles(filePaths []string) (InvertIndex, error) {
	if err := checkIsFilePathsValid(filePaths); err != nil {
		return nil, err
	}

	ctx := context.TODO()
	sem := semaphore.NewWeighted(maxOpenFilesCount)
	errs, _ := errgroup.WithContext(context.Background())

	for docID, filePath := range filePaths {
		filePath := filePath
		docID := DocID(docID)

		if err := sem.Acquire(ctx, 1); err != nil {
			return nil, fmt.Errorf("failed to acquire semaphore: %w", err)
		}

		errs.Go(func() error {
			defer sem.Release(1)

			return i.indexFile(filePath, docID)
		})
	}

	if err := errs.Wait(); err != nil {
		return nil, fmt.Errorf("error while indexing file: %w", err)
	}

	return mergeIndexingResults(len(filePaths))
}

func (i *Indexer) indexFile(filePath string, docID DocID) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error while opening file: %w", err)
	}

	dec := decoder.NewTXTDecoder(file)
	tokens := make(map[string]bool)
	dictionary := make([]string, 0)

	for decoded, ok := dec.DecodeNext(); ok; decoded, ok = dec.DecodeNext() {
		normalized := i.normalizer.Normalize(decoded)
		if normalized == "" {
			continue
		}

		if _, ok := tokens[normalized]; !ok {
			tokens[normalized] = true

			dictionary = append(dictionary, normalized)
		}
	}

	return encodeDictionaryToFile(dictionary, getFilenameFromDocID(docID))
}

func mergeIndexingResults(resultsCount int) (InvertIndex, error) {
	invertIndex := make(InvertIndex)

	for docID := DocID(0); int(docID) < resultsCount; docID++ {
		filename := getFilenameFromDocID(docID)

		terms, err := decodeDictionaryFromFile(filename)
		if err != nil {
			return nil, err
		}

		if err := os.Remove(filename); err != nil {
			return nil, fmt.Errorf("error removing file: %w", err)
		}

		for _, t := range terms {
			if docIDs, ok := invertIndex[t]; ok {
				insertToListWithKeepSorting(docIDs, docID)
			} else {
				invertIndex[t] = list.New()
				invertIndex[t].PushBack(docID)
			}
		}
	}

	return invertIndex, nil
}

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

func getFilenameFromDocID(docID DocID) string {
	return fmt.Sprintf("./%d", docID)
}

func encodeDictionaryToFile(dictionary []string, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error while creating file: %w", err)
	}

	e := gob.NewEncoder(f)
	if err := e.Encode(dictionary); err != nil {
		return fmt.Errorf("error while encoding: %w", err)
	}

	_ = f.Close()

	return nil
}

func decodeDictionaryFromFile(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while opening file: %w", err)
	}

	var terms []string

	d := gob.NewDecoder(f)
	if err := d.Decode(&terms); err != nil {
		return nil, fmt.Errorf("error while decoding terms from file: %w", err)
	}

	_ = f.Close()

	return terms, nil
}
