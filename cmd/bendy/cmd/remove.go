package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ythosa/bendy/internal/storage"
)

type RemoveFileCommand struct {
	filesStorage storage.Files
}

func NewRemoveFileCommand(
	filesStorage storage.Files,
) *RemoveFileCommand {
	return &RemoveFileCommand{
		filesStorage: filesStorage,
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

			fmt.Printf("File successfully removed")
		},
	}
}
