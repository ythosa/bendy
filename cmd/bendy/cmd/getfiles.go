package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getIndexingFilesCmd)
}

var getIndexingFilesCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "Returns indexing files",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		files, err := storage.Files.Get()
		if err != nil {
			fmt.Printf("Error while getting indexing files: %s", err)

			return
		}

		for i, f := range files {
			fmt.Printf("%d) %s\n", i+1, f)
		}
	},
}
