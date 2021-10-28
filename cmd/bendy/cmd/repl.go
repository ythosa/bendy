package cmd

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/ythosa/bendy/internal/evaluate"
	"github.com/ythosa/bendy/internal/evaluate/lexer"
	"github.com/ythosa/bendy/internal/evaluate/parser"
	"github.com/ythosa/bendy/internal/index"
)

const prompt = ">> "

func init() {
	rootCmd.AddCommand(replCmd)
}

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Returns indexing files",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		i, err := storage.Index.Get()
		if err != nil {
			fmt.Printf("Error while getting indexes: %s", err)

			return
		}

		docIDs, err := getAllDocIDs()
		if err != nil {
			fmt.Printf("Error getting all doc ids: %s", err)

			return
		}

		evaluator := evaluate.NewEvaluator(i, docIDs)

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print(prompt)
			scanned := scanner.Scan()
			if !scanned {
				return
			}

			line := scanner.Text()
			l := lexer.New(line)
			p := parser.New(l)

			request := p.ParseRequest()
			if len(p.Errors()) != 0 {
				printParserErrors(os.Stdout, p.Errors())
				continue
			}

			evaluated := evaluator.Eval(request)
			if evaluated != nil {
				_, _ = io.WriteString(os.Stdout, evaluated.Inspect())
				_, _ = io.WriteString(os.Stdout, "\n")
			}
		}
	},
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}

func getAllDocIDs() (*list.List, error) {
	files, err := storage.Files.Get()
	if err != nil {
		return nil, err
	}

	docIDs := make([]index.DocID, len(files))
	for i, _ := range files {
		docIDs[i] = index.DocID(i)
	}

	return index.SliceToList(docIDs), nil
}
