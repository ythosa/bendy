package indexer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ythosa/bendy/internal/indexer"
	"github.com/ythosa/bendy/internal/normalizer"
)

func TestNewIndexer(t *testing.T) {
	t.Parallel()

	ix := indexer.NewIndexer(normalizer.NewEnglishNormalizer())

	assert.NotNil(t, ix)
}
