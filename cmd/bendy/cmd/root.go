package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type RootCommand struct {
	commands []Command
}

func NewRootCommand(commands []Command) *RootCommand {
	return &RootCommand{commands: commands}
}

func (r *RootCommand) getCLI() *cobra.Command {
	root := &cobra.Command{
		Use:   "bendy",
		Short: "Bendy is bool search engine",
		Long: `Bendy is bool search engine. 
You can help with the development here: https://github.com/ythosa/bendy.`,
	}

	for _, command := range r.commands {
		root.AddCommand(command.getCLI())
	}

	return root
}

func (r *RootCommand) Execute() {
	if err := r.getCLI().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
