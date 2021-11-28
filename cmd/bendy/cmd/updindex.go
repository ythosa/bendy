package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ythosa/bendy/internal/index"
	"github.com/ythosa/bendy/internal/storage"
)

type UpdateIndexCommand struct {
	filesStorage storage.Files
	indexStorage storage.Index
	indexer      *index.Indexer
}

func NewUpdateIndexCommand(
	filesStorage storage.Files,
	indexStorage storage.Index,
	indexer *index.Indexer,
) *UpdateIndexCommand {
	return &UpdateIndexCommand{
		filesStorage: filesStorage,
		indexStorage: indexStorage,
		indexer:      indexer,
	}
}

func (u *UpdateIndexCommand) getCLI() *cobra.Command {
	return &cobra.Command{
		Use:     "update",
		Aliases: []string{"u"},
		Short:   "Updates invert index",
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]

			files, err := u.filesStorage.Get()
			if err != nil {
				fmt.Printf("Error while getting files to index: %s", err)

				return
			}

			i, err := u.indexer.IndexFiles(files)
			if err != nil {
				fmt.Printf("Error while indexing files: %s", err)
				_ = u.filesStorage.Delete(filename)

				return
			}

			if err := u.indexStorage.Set(i); err != nil {
				fmt.Printf("Error while updating indexes: %s", err)

				return
			}

			fmt.Printf("Invert index successfully updated")
		},
	}
}
