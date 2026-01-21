package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewListCmd creates a Cobra command for listing todos
func NewListCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "list [FILTER]",
		Short: "List all todo items",
		Long: `Lists all todo items, optionally filtered by search term.

# list all tasks
togodo list

# list tasks containing "milk"
togodo list milk
`,
		Aliases: []string{"l", "ls"},
		RunE: func(command *cobra.Command, args []string) error {
			// Get search filter if provided
			searchQuery := ""
			if len(args) > 0 {
				searchQuery = args[0]
			}

			// Call business logic
			result, err := cmd.List(repo, searchQuery)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			for i, todo := range result.Todos {
				fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumbers[i], todo.Text)
			}
			fmt.Fprintln(command.OutOrStdout(), "--")
			fmt.Fprintf(command.OutOrStdout(), "TODO: %d of %d tasks shown\n", result.ShownCount, result.TotalCount)
			return nil
		},
	}
}
