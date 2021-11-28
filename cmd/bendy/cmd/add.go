package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ythosa/bendy/internal/storage"
)

type AddFileCommand struct {
	filesStorage storage.Files
}

func NewAddFileCommand(filesStorage storage.Files) *AddFileCommand {
	return &AddFileCommand{
		filesStorage: filesStorage,
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

			fmt.Printf("File successfully added")
		},
	}
}
