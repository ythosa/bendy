package index

import (
	"container/list"
	"context"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"github.com/ythosa/bendy/internal/config"
	"github.com/ythosa/bendy/internal/decoding"
	"github.com/ythosa/bendy/internal/normalizing"
	"github.com/ythosa/bendy/pkg/fcheck"
)

type Indexer struct {
	decoder    decoding.Decoder
	normalizer normalizing.Normalizer
	config     *config.Index
}

func NewIndexer(decoder decoding.Decoder, normalizer normalizing.Normalizer, config *config.Index) *Indexer {
	return &Indexer{
		decoder:    decoder,
		normalizer: normalizer,
		config:     config,
	}
}

func (i *Indexer) IndexFiles(filePaths []string) (InvertIndex, error) {
	if err := fcheck.CheckIsFilePathsValid(filePaths); err != nil {
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

	defer file.Close()

	decoder, err := i.decoder.GetDecoder(file)
	if err != nil {
		return fmt.Errorf("error while getting decoder for file: %w", err)
	}

	tokens := make(map[string]bool)
	dictionary := make([]string, 0)

	for decoded, ok := decoder.DecodeNext(); ok; decoded, ok = decoder.DecodeNext() {
		normalizer, err := i.normalizer.GetNormalizer(decoded)
		if err != nil {
			return fmt.Errorf("error while getting normalizer for %s: %w", decoded, err)
		}

		normalized := normalizer.Normalize(decoded)
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
			if index, ok := invertIndex[t]; ok {
				index.Insert(docID)
			} else {
				invertIndex[t] = NewIndex(list.New())
				invertIndex[t].DocIDs.PushBack(docID)
			}
		}
	}

	return invertIndex, nil
}

func (i *Indexer) getFilenameFromDocID(docID DocID) string {
	return fmt.Sprintf("%s%d", i.config.TempFilesStoragePath, docID)
}
