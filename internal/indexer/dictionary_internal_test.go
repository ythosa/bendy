package indexer

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
