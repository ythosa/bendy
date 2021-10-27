package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ythosa/bendy/internal/config"
	"github.com/ythosa/bendy/internal/storage/filestorage"
)

var storage = filestorage.NewStorage(config.Get().Storage)

var rootCmd = &cobra.Command{
	Use:   "bendy",
	Short: "Bendy is bool search engine",
	Long:  `Bendy is bool search engine :)`,
}

func Execute() error {
	return rootCmd.Execute()
}
