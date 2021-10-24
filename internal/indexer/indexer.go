package indexer

import (
	"container/list"
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"sync/atomic"

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

	var docID DocID

	for _, fp := range filePaths {
		fp := fp

		if err := sem.Acquire(ctx, 1); err != nil {
			return nil, fmt.Errorf("failed to acquire semaphore: %w", err)
		}

		errs.Go(func() error {
			defer sem.Release(1)
			atomic.AddUint32((*uint32)(&docID), 1)

			return i.indexFile(fp, docID)
		})
	}

	if err := errs.Wait(); err != nil {
		return nil, fmt.Errorf("error while indexing files: %w", err)
	}

	return mergeIndexingResults(len(filePaths))
}

func mergeIndexingResults(resultsCount int) (InvertIndex, error) {
	invertIndex := make(InvertIndex)

	for i := 1; i <= resultsCount; i++ {
		terms, err := decodeDictionaryFromFile(getFilenameFromDocID(DocID(i)))
		if err != nil {
			return nil, err
		}

		for _, t := range terms {
			if docIDs, ok := invertIndex[t]; ok {
				insertWithKeepSorting(docIDs, DocID(i))
			} else {
				invertIndex[t] = list.New()
				invertIndex[t].PushBack(DocID(i))
			}
		}
	}

	return invertIndex, nil
}

func (i *Indexer) indexFile(filePath string, docID DocID) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error while opening file: %w", err)
	}

	d := decoder.NewTXTDecoder(f)
	tokens := make(map[string]bool)
	dictionary := make([]string, 0)

	for decoded, ok := d.DecodeNext(); ok; decoded, ok = d.DecodeNext() {
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

func insertWithKeepSorting(l *list.List, docID DocID) {
	var prev *list.Element

	for cur := l.Front(); cur != nil; cur = cur.Next() {
		if value, _ := cur.Value.(DocID); value > docID {
			break
		}

		prev = cur
	}

	if prev != nil {
		l.InsertAfter(docID, prev)
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
