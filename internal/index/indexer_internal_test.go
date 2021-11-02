package index

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ythosa/bendy/internal/config"
	"github.com/ythosa/bendy/internal/decoding"
	"github.com/ythosa/bendy/internal/normalizing"
)

func getTestIndexer() *Indexer {
	return NewIndexer(
		decoding.NewDecoderImpl(),
		normalizing.NewNormalizerImpl(),
		config.Get().Index,
	)
}

func TestIndexer_IndexFiles(t *testing.T) {
	filePaths := []string{"./test_index_files1.txt", "./test_index_files2.txt", "./test_index_files3.txt"}
	indxr := getTestIndexer()

	// indexing error: no files
	_, err := indxr.IndexFiles(filePaths)

	assert.NotNil(t, err)

	// ok
	f, _ := os.Create(filePaths[0])
	_, _ = f.WriteString("kek keek keeek")
	f, _ = os.Create(filePaths[1])
	_, _ = f.WriteString("kek keek")
	f, _ = os.Create(filePaths[2])
	_, _ = f.WriteString("kek")

	invertIndex, err := indxr.IndexFiles(filePaths)
	assert.Nil(t, err)

	_ = os.Remove(filePaths[0])
	_ = os.Remove(filePaths[1])
	_ = os.Remove(filePaths[2])

	expectedInvertIndex := make(InvertIndex)
	expectedInvertIndex["kek"] = NewIndex(SliceToList([]DocID{0, 1, 2}))
	expectedInvertIndex["keek"] = NewIndex(SliceToList([]DocID{0, 1}))
	expectedInvertIndex["keeek"] = NewIndex(SliceToList([]DocID{0}))

	for k, v := range invertIndex {
		expectedValue, ok := expectedInvertIndex[k]
		assert.Equal(t, ok, true)
		CompareLists(t, expectedValue.DocIDs, v.DocIDs)
	}
}

func TestIndexer_MergeIndexingResults(t *testing.T) {
	// there is no generated index files
	i := getTestIndexer()
	_, err := i.mergeIndexingResults(1)
	assert.NotNil(t, err)

	// ok
	const inputFilePath = "./merge_indexing_results.txt"

	f, _ := os.Create(inputFilePath)
	_, _ = f.WriteString("input data :)")
	_ = i.indexFile(inputFilePath, 0)
	_ = os.Remove(inputFilePath)

	f, _ = os.Create(inputFilePath)
	_, _ = f.WriteString("input data kek :)")
	_ = i.indexFile(inputFilePath, 1)
	_ = os.Remove(inputFilePath)

	ii, err := i.mergeIndexingResults(2)
	assert.Nil(t, err)
	CompareLists(t, SliceToList([]DocID{0, 1}), ii["input"].DocIDs)
	CompareLists(t, SliceToList([]DocID{0, 1}), ii["data"].DocIDs)
	CompareLists(t, SliceToList([]DocID{1}), ii["kek"].DocIDs)
}

func TestIndexer_IndexFile(t *testing.T) {
	t.Parallel()

	i := getTestIndexer()

	// error creating file
	assert.NotNil(t, i.indexFile("./kek/kek", 1))

	// ok
	const (
		validFilename = "test_index_file.txt"
		docID         = 1
	)

	f, _ := os.Create(validFilename)
	_, _ = f.WriteString("kek kek1 !")

	assert.Nil(t, i.indexFile(validFilename, docID))

	// remove generated files
	_ = os.Remove(validFilename)
	_ = os.Remove(i.getFilenameFromDocID(docID))
}

func TestIndexer_GetFileNameFromDocID(t *testing.T) {
	t.Parallel()

	indxr := getTestIndexer()
	filename := indxr.getFilenameFromDocID(1)
	assert.Equal(t, fmt.Sprintf("%s1", config.Get().Index.TempFilesStoragePath), filename)
}
