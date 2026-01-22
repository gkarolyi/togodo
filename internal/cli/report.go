package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/internal/config"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewReportCmd creates a Cobra command for generating task statistics
func NewReportCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "report",
		Short: "Generate task statistics",
		Long: `Report generates statistics about tasks in todo.txt and done.txt.
Shows the current date, total tasks, completed tasks, and pending tasks.

# generate report
togodo report
`,
		Args: cobra.NoArgs,
		RunE: func(command *cobra.Command, args []string) error {
			// Get done.txt path (same directory as todo.txt, different name)
			todoPath := config.GetTodoTxtPath()
			doneFilePath := getDoneFilePath(todoPath)

			// Create reader for done.txt
			doneReader := todotxtlib.NewFileReader(doneFilePath)

			// Call business logic
			result, err := cmd.Report(repo, doneReader)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli format
			// Format: YYYY-MM-DDTHH:MM:SS TOTAL DONE PENDING
			timestamp := result.Date.Format("2006-01-02T15:04:05")
			fmt.Fprintf(command.OutOrStdout(), "%s %d %d %d\n",
				timestamp,
				result.Total,
				result.Done,
				result.Pending)

			return nil
		},
	}
}
