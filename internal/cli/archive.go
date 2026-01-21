package cli

import (
	"fmt"
	"path/filepath"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/internal/config"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewArchiveCmd creates a Cobra command for archiving completed tasks
func NewArchiveCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "archive",
		Short: "Move completed tasks from todo.txt to done.txt",
		Long: `Archive moves all completed tasks from todo.txt to done.txt.
The completed tasks are appended to done.txt and removed from todo.txt.

# archive completed tasks
togodo archive
`,
		Args: cobra.NoArgs,
		RunE: func(command *cobra.Command, args []string) error {
			// Get done.txt path (same directory as todo.txt, different name)
			todoPath := config.GetTodoTxtPath()
			doneFilePath := getDoneFilePath(todoPath)

			// Create reader and writer for done.txt
			doneReader := todotxtlib.NewFileReader(doneFilePath)
			doneWriter := todotxtlib.NewFileWriter(doneFilePath)

			// Call business logic
			result, err := cmd.Archive(repo, doneReader, doneWriter)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			if len(result.ArchivedTodos) == 0 {
				fmt.Fprintf(command.OutOrStdout(), "TODO: todo.txt does not contain any done tasks.\n")
			} else {
				// Print each archived todo
				for _, todo := range result.ArchivedTodos {
					fmt.Fprintf(command.OutOrStdout(), "%s\n", todo.Text)
				}
				fmt.Fprintf(command.OutOrStdout(), "TODO: todo.txt archived.\n")
			}

			return nil
		},
	}
}

// getDoneFilePath returns the path to done.txt based on the todo.txt path
func getDoneFilePath(todoPath string) string {
	dir := filepath.Dir(todoPath)
	return filepath.Join(dir, "done.txt")
}
