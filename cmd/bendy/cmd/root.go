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
Use case: 
1) Provide path to bendy config folder with env variable BENDY_CONFIGS_FOLDER_PATH;
2) Provide path to bendy config name with env variable BENDY_CONFIG_NAME 
if you are not satisfied with the default config;
3) Add files with "bendy add";
4) Update index with "bendy update";
5) Run REPL and write your search queries!
You can use command "bendy <command> --help" to get help for <command>.

Btw you can help with the development here: https://github.com/ythosa/bendy.`,
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
