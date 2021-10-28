package index_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ythosa/bendy/internal/config"
	"github.com/ythosa/bendy/internal/index"
	"github.com/ythosa/bendy/internal/normalizer"
)

func TestNewIndexer(t *testing.T) {
	t.Parallel()

	ix := index.NewIndexer(normalizer.NewEnglishNormalizer(), config.Get().Index)

	assert.NotNil(t, ix)
}
