package main

import (
	"github.com/ythosa/bendy/cmd/bendy/cmd"
	"github.com/ythosa/bendy/internal/config"
	"github.com/ythosa/bendy/internal/decoding"
	"github.com/ythosa/bendy/internal/index"
	"github.com/ythosa/bendy/internal/normalizing"
	"github.com/ythosa/bendy/internal/storage"
	"github.com/ythosa/bendy/internal/storage/filestorage"
)

func getSubcommands(storage *storage.Storage, indexer *index.Indexer) []cmd.Command {
	return []cmd.Command{
		cmd.NewAddFileCommand(storage.Files, storage.Index, indexer),
		cmd.NewGetFilesCommand(storage.Files),
		cmd.NewRemoveFileCommand(storage.Files, storage.Index, indexer),
		cmd.NewReplCommand(storage.Files, storage.Index),
	}
}

func main() {
	indexer := index.NewIndexer(decoding.NewDecoderImpl(), normalizing.NewNormalizerImpl(), config.Get().Index)
	s := filestorage.NewStorage(config.Get().Storage)

	rootCmd := cmd.NewRootCommand(getSubcommands(s, indexer))

	rootCmd.Execute()
}
