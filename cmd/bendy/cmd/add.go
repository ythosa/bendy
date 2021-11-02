package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ythosa/bendy/internal/index"
	"github.com/ythosa/bendy/internal/storage"
)

type AddFileCommand struct {
	filesStorage storage.Files
	indexStorage storage.Index
	indexer      *index.Indexer
}

func NewAddFileCommand(
	filesStorage storage.Files,
	indexStorage storage.Index,
	indexer *index.Indexer,
) *AddFileCommand {
	return &AddFileCommand{
		filesStorage: filesStorage,
		indexStorage: indexStorage,
		indexer:      indexer,
	}
}

func (a *AddFileCommand) getCLI() *cobra.Command {
	return &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Adds file to index",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]

			fmt.Printf("Adding file %s...\n", filename)

			if err := a.filesStorage.Put(filename); err != nil {
				fmt.Printf("Error while adding file to index: %s", err)

				return
			}

			files, err := a.filesStorage.Get()
			if err != nil {
				fmt.Printf("Error while getting files to index: %s", err)

				return
			}

			i, err := a.indexer.IndexFiles(files)
			if err != nil {
				fmt.Printf("Error while indexing files: %s", err)
				_ = a.filesStorage.Delete(filename)

				return
			}

			if err := a.indexStorage.Set(i); err != nil {
				fmt.Printf("Error while updating indexes: %s", err)

				return
			}

			fmt.Printf("File successfully added")
		},
	}
}
