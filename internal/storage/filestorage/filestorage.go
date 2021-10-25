package filestorage

import (
	"github.com/ythosa/bendy/internal/config"
	"github.com/ythosa/bendy/internal/storage"
)

func NewStorage(cfg config.Storage) *storage.Storage {
	return &storage.Storage{
		Index: nil,
		Files: NewFilesImpl(cfg.IndexingFilesFilenamesPath),
	}
}
