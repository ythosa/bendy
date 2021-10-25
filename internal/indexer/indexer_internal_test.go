package indexer

import (
	"container/list"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ythosa/bendy/internal/normalizer"
)

func TestIndexer_IndexFiles(t *testing.T) {
	filePaths := []string{"./test_index_files1", "./test_index_files2", "./test_index_files3"}
	indxr := NewIndexer(normalizer.NewEnglishNormalizer())

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

func TestMergeIndexingResults(t *testing.T) {
	// there is no generated index files
	_, err := mergeIndexingResults(1)
	assert.NotNil(t, err)

	// ok
	const inputFilePath = "./merge_indexing_results"

	i := NewIndexer(normalizer.NewEnglishNormalizer())

	f, _ := os.Create(inputFilePath)
	_, _ = f.WriteString("input data :)")
	_ = i.indexFile(inputFilePath, 0)
	_ = os.Remove(inputFilePath)

	f, _ = os.Create(inputFilePath)
	_, _ = f.WriteString("input data kek :)")
	_ = i.indexFile(inputFilePath, 1)
	_ = os.Remove(inputFilePath)

	ii, err := mergeIndexingResults(2)
	assert.Nil(t, err)
	compareLists(t, sliceToList([]DocID{0, 1}), ii["input"])
	compareLists(t, sliceToList([]DocID{0, 1}), ii["data"])
	compareLists(t, sliceToList([]DocID{1}), ii["kek"])
}

func TestIndexer_IndexFile(t *testing.T) {
	t.Parallel()

	i := NewIndexer(normalizer.NewEnglishNormalizer())

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
	_ = os.Remove(getFilenameFromDocID(docID))
}

func TestInsertWithKeepSorting(t *testing.T) {
	t.Parallel()

	type args struct {
		l     *list.List
		docID DocID
	}

	testCases := []struct {
		input    args
		expected *list.List
	}{
		{
			input: args{
				l:     sliceToList([]DocID{1, 2, 5, 6}),
				docID: 3,
			},
			expected: sliceToList([]DocID{1, 2, 3, 5, 6}),
		},
		{
			input: args{
				l:     sliceToList([]DocID{2}),
				docID: 1,
			},
			expected: sliceToList([]DocID{1, 2}),
		},
		{
			input: args{
				l:     sliceToList([]DocID{2}),
				docID: 3,
			},
			expected: sliceToList([]DocID{2, 3}),
		},
	}

	for _, tc := range testCases {
		insertToListWithKeepSorting(tc.input.l, tc.input.docID)
		compareLists(t, tc.expected, tc.input.l)
	}
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

func TestCheckIsFilePathsValid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		filePaths     []string
		expectedError bool
	}{
		{
			filePaths:     []string{"./indexer.go"},
			expectedError: false,
		},
		{
			filePaths:     []string{"./invalid_file_path"},
			expectedError: true,
		},
		{
			filePaths:     []string{"./indexer", "./invalid_file_path"},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		if tc.expectedError {
			assert.NotNil(t, checkIsFilePathsValid(tc.filePaths))
		} else {
			assert.Nil(t, checkIsFilePathsValid(tc.filePaths))
		}
	}
}

func TestGetFileNameFromDocID(t *testing.T) {
	t.Parallel()

	filename := getFilenameFromDocID(1)
	assert.Equal(t, "./1", filename)
}

func TestDumpingDictionary(t *testing.T) {
	t.Parallel()

	toDump := []string{"it", "is", "dump", "value"}
	filePath := "./dump_test"

	assert.Nil(t, encodeDictionaryToFile(toDump, filePath))

	decoded, err := decodeDictionaryFromFile(filePath)
	assert.Nil(t, err)

	assert.Equal(t, reflect.DeepEqual(toDump, decoded), true)

	assert.Nil(t, os.Remove(filePath))
}

func TestEncodeDictionaryToFile(t *testing.T) {
	t.Parallel()

	// creating file error
	invalidFilePath := "./kek/encode_test"
	assert.NotNil(t, encodeDictionaryToFile(nil, invalidFilePath))
	_ = os.Remove(invalidFilePath)
}

func TestDecodeDictionaryFromFile(t *testing.T) {
	t.Parallel()

	// file is not exist error
	filePath := "./decode_test"
	decoded, err := decodeDictionaryFromFile(filePath)

	assert.Nil(t, decoded)
	assert.NotNil(t, err)

	// decoding error
	f, _ := os.Create(filePath)
	_, _ = f.WriteString("kek")
	decoded, err = decodeDictionaryFromFile(filePath)

	assert.Nil(t, decoded)
	assert.NotNil(t, err)

	assert.Nil(t, os.Remove(filePath))
}
