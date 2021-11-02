package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ythosa/bendy/internal/index"
	"github.com/ythosa/bendy/internal/storage"
)

type RemoveFileCommand struct {
	filesStorage storage.Files
	indexStorage storage.Index
	indexer      *index.Indexer
}

func NewRemoveFileCommand(
	filesStorage storage.Files,
	indexStorage storage.Index,
	indexer *index.Indexer,
) *RemoveFileCommand {
	return &RemoveFileCommand{
		filesStorage: filesStorage,
		indexStorage: indexStorage,
		indexer:      indexer,
	}
}

func (r *RemoveFileCommand) getCLI() *cobra.Command {
	return &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "Removes file from index",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]
			fmt.Printf("Removing file %s...\n", filename)

			if err := r.filesStorage.Delete(filename); err != nil {
				fmt.Printf("Error while deleting file from index: %s", err)

				return
			}

			files, err := r.filesStorage.Get()
			if err != nil {
				fmt.Printf("Error while getting indexing files: %s", err)

				return
			}

			i, err := r.indexer.IndexFiles(files)
			if err != nil {
				fmt.Printf("Error while indexing files: %s", err)

				return
			}

			if err := r.indexStorage.Set(i); err != nil {
				fmt.Printf("Error while updating indexes: %s", err)

				return
			}

			fmt.Printf("File successfully removed")
		},
	}
}
