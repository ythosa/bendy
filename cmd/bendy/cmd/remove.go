package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeFileFromIndexCmd)
}

var removeFileFromIndexCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Removes file from index",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		fmt.Printf("Removing file %s...\n", filename)

		if err := storage.Files.Delete(filename); err != nil {
			fmt.Printf("Error while deleting file from index: %s", err)

			return
		}

		files, err := storage.Files.Get()
		if err != nil {
			fmt.Printf("Error while getting indexing files: %s", err)

			return
		}

		i, err := indexer.IndexFiles(files)
		if err != nil {
			fmt.Printf("Error while indexing files: %s", err)

			return
		}

		if err := storage.Index.Set(i); err != nil {
			fmt.Printf("Error while updating indexes: %s", err)

			return
		}

		fmt.Printf("File successfully removed")
	},
}
