package indexer

import (
	"container/list"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ythosa/bendy/internal/config"
	"github.com/ythosa/bendy/internal/normalizer"
)

func TestIndexer_IndexFiles(t *testing.T) {
	filePaths := []string{"./test_index_files1", "./test_index_files2", "./test_index_files3"}
	indxr := NewIndexer(normalizer.NewEnglishNormalizer(), config.Get().Index)

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
	expectedInvertIndex["kek"] = sliceToList([]DocID{0, 1, 2})
	expectedInvertIndex["keek"] = sliceToList([]DocID{0, 1})
	expectedInvertIndex["keeek"] = sliceToList([]DocID{0})

	for k, v := range invertIndex {
		expectedValue, ok := expectedInvertIndex[k]
		assert.Equal(t, ok, true)
		compareLists(t, expectedValue, v)
	}
}

func TestIndexer_MergeIndexingResults(t *testing.T) {
	// there is no generated index files
	i := NewIndexer(normalizer.NewEnglishNormalizer(), config.Get().Index)
	_, err := i.mergeIndexingResults(1)
	assert.NotNil(t, err)

	// ok
	const inputFilePath = "./merge_indexing_results"

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
	compareLists(t, sliceToList([]DocID{0, 1}), ii["input"])
	compareLists(t, sliceToList([]DocID{0, 1}), ii["data"])
	compareLists(t, sliceToList([]DocID{1}), ii["kek"])
}

func TestIndexer_IndexFile(t *testing.T) {
	t.Parallel()

	i := NewIndexer(normalizer.NewEnglishNormalizer(), config.Get().Index)

	// error creating file
	assert.NotNil(t, i.indexFile("./kek/kek", 1))

	// ok
	const (
		validFilename = "test_index_file"
		docID         = 1
	)

	f, _ := os.Create(validFilename)
	_, _ = f.WriteString("kek kek1 !")

	assert.Nil(t, i.indexFile(validFilename, docID))

	// remove generated files
	_ = os.Remove(validFilename)
	_ = os.Remove(i.getFilenameFromDocID(docID))
}

func sliceToList(slice []DocID) *list.List {
	l := list.New()
	for _, v := range slice {
		l.PushBack(v)
	}

	return l
}

func compareLists(t *testing.T, expected *list.List, actual *list.List) {
	t.Helper()

	assert.Equal(t, expected.Len(), actual.Len(), "lists are different sizes")

	expectedElement := expected.Front()
	actualElement := actual.Front()

	for expectedElement != nil {
		expectedValue, _ := expectedElement.Value.(DocID)
		actualValue, _ := actualElement.Value.(DocID)
		assert.Equal(t, expectedValue, actualValue)

		expectedElement = expectedElement.Next()
		actualElement = actualElement.Next()
	}
}

func TestIndexer_GetFileNameFromDocID(t *testing.T) {
	t.Parallel()

	indxr := NewIndexer(normalizer.NewEnglishNormalizer(), config.Get().Index)
	filename := indxr.getFilenameFromDocID(1)
	assert.Equal(t, fmt.Sprintf("%s1", config.Get().Index.TempFilesStoragePath), filename)
}
