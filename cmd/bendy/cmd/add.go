package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds file to index",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Adding file %s...\n", args[0])
		if err := storage.Put(args[0]); err != nil {
			fmt.Printf("Error while adding file to index: %s", err)

			return
		}
	},
}
