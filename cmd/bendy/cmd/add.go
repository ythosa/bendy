package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addFileToIndexCmd)
}

var addFileToIndexCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds file to index",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		fmt.Printf("Adding file %s...\n", filename)

		if err := storage.Put(filename); err != nil {
			fmt.Printf("Error while adding file to index: %s", err)

			return
		}

		files, err := storage.Files.Get()
		if err != nil {
			fmt.Printf("Error while getting files to index: %s", err)

			return
		}

		i, err := indexer.IndexFiles(files)
		if err != nil {
			fmt.Printf("Error while indexing files: %s", err)
			_ = storage.Files.Delete(filename)

			return
		}

		if err := storage.Index.Set(i); err != nil {
			fmt.Printf("Error while updating indexes: %s", err)

			return
		}

		fmt.Printf("File successfully added")
	},
}
