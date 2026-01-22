package cli

import (
	"fmt"
	"strings"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewListpriCmd creates a Cobra command for listing todos by priority
func NewListpriCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "listpri [PRIORITY]",
		Short: "List todos with specific priority",
		Long: `Lists all todos with the specified priority.

# list all A-priority tasks
togodo listpri A
`,
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"lsp"},
		RunE: func(command *cobra.Command, args []string) error {
			// Default to listing all priorities
			priority := ""
			if len(args) > 0 {
				priority = strings.ToUpper(args[0])
			}

			// Call business logic
			result, err := cmd.Listpri(repo, priority)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			for i, todo := range result.Todos {
				fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumbers[i], todo.Text)
			}
			fmt.Fprintln(command.OutOrStdout(), "--")
			fmt.Fprintf(command.OutOrStdout(), "TODO: %d of %d tasks shown\n", len(result.Todos), result.TotalCount)
			return nil
		},
	}
}
