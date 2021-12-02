package cmd

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/ythosa/bendy/internal/eval"
	"github.com/ythosa/bendy/internal/eval/lexer"
	"github.com/ythosa/bendy/internal/eval/object"
	"github.com/ythosa/bendy/internal/eval/parser"
	"github.com/ythosa/bendy/internal/index"
	"github.com/ythosa/bendy/internal/storage"
)

type REPLCommand struct {
	filesStorage storage.Files
	indexStorage storage.Index
}

func NewReplCommand(filesStorage storage.Files, indexStorage storage.Index) *REPLCommand {
	return &REPLCommand{filesStorage: filesStorage, indexStorage: indexStorage}
}

func (r *REPLCommand) getCLI() *cobra.Command {
	return &cobra.Command{
		Use:   "repl",
		Short: "Returns indexing files",
		Long: `Using:
1) write some request and press ENTER;
2) use CTRL + D to exit.

Search format:
1) & - for AND operation;
2) | - for OR operation;
3) ! - for OR operation;
4) "..." - for the word you want to find;
5) () - brackets are used the way you are used to.

Example: "word1" | ("word2" & "word3") | !"word4" & "word5"`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			i, err := r.indexStorage.Get()
			if err != nil {
				fmt.Printf("Error while getting indexes: %s", err)

				return
			}

			files, err := r.filesStorage.Get()
			if err != nil {
				fmt.Printf("Error while getting all files: %s", err)

				return
			}

			evaluator := eval.NewEvaluator(i, getAllDocIDs(files))

			scanner := bufio.NewScanner(os.Stdin)
			for {
				fmt.Print(">> ")
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
				switch v := evaluated.(type) {
				case *object.DocIDs:
					printDocuments(files, v.Value)
				case nil:
					_, _ = fmt.Fprintf(os.Stdout, "Empty result\n")
				default:
					_, _ = fmt.Fprintf(os.Stdout, "%s\n", evaluated.Inspect())
				}
			}
		},
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}

func printDocuments(files []string, docIDs *list.List) {
	for e := docIDs.Front(); e != nil; e = e.Next() {
		docID, _ := e.Value.(index.DocID)
		_, _ = fmt.Fprintf(os.Stdout, "%d) %s\n", docID, files[docID])
	}
}

func getAllDocIDs(files []string) *list.List {
	docIDs := make([]index.DocID, len(files))
	for i := range files {
		docIDs[i] = index.DocID(i)
	}

	return index.SliceToList(docIDs)
}
