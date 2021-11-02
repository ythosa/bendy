package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ythosa/bendy/internal/storage"
)

type GetFilesCommand struct {
	storage storage.Files
}

func NewGetFilesCommand(storage storage.Files) *GetFilesCommand {
	return &GetFilesCommand{storage: storage}
}

func (g *GetFilesCommand) getCLI() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "Returns indexing files",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			files, err := g.storage.Get()
			if err != nil {
				fmt.Printf("Error while getting indexing files: %s", err)

				return
			}

			for i, f := range files {
				fmt.Printf("%d) %s\n", i+1, f)
			}
		},
	}
}
