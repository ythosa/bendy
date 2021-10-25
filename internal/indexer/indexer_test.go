package indexer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ythosa/bendy/internal/config"
	"github.com/ythosa/bendy/internal/indexer"
	"github.com/ythosa/bendy/internal/normalizer"
)

func TestNewIndexer(t *testing.T) {
	t.Parallel()

	ix := indexer.NewIndexer(normalizer.NewEnglishNormalizer(), config.Get().Index)

	assert.NotNil(t, ix)
}
