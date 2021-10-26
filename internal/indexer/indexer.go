package indexer

import (
	"container/list"
	"context"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"github.com/ythosa/bendy/internal/config"
	"github.com/ythosa/bendy/internal/decoder"
	"github.com/ythosa/bendy/internal/normalizer"
	"github.com/ythosa/bendy/pkg/utils"
)

type DocID uint32

type Indexer struct {
	normalizer normalizer.Normalizer
	config     *config.Index
}

func NewIndexer(normalizer normalizer.Normalizer, config *config.Index) *Indexer {
	return &Indexer{
		normalizer: normalizer,
		config:     config,
	}
}

func (i *Indexer) IndexFiles(filePaths []string) (InvertIndex, error) {
	if err := utils.CheckIsFilePathsValid(filePaths); err != nil {
		return nil, fmt.Errorf("error while checking files: %w", err)
	}

	ctx := context.TODO()
	sem := semaphore.NewWeighted(i.config.MaxOpenFilesCount)
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

	return i.mergeIndexingResults(len(filePaths))
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

	return encodeDictionaryToFile(dictionary, i.getFilenameFromDocID(docID))
}

func (i *Indexer) mergeIndexingResults(resultsCount int) (InvertIndex, error) {
	invertIndex := make(InvertIndex)

	for docID := DocID(0); int(docID) < resultsCount; docID++ {
		filename := i.getFilenameFromDocID(docID)

		terms, err := decodeDictionaryFromFile(filename)
		if err != nil {
			return nil, err
		}

		if err := os.Remove(filename); err != nil {
			return nil, fmt.Errorf("error removing file: %w", err)
		}

		for _, t := range terms {
			if docIDs, ok := invertIndex[t]; ok {
				InsertToListWithKeepSorting(docIDs, docID)
			} else {
				invertIndex[t] = list.New()
				invertIndex[t].PushBack(docID)
			}
		}
	}

	return invertIndex, nil
}

func (i *Indexer) getFilenameFromDocID(docID DocID) string {
	return fmt.Sprintf("%s%d", i.config.TempFilesStoragePath, docID)
}
