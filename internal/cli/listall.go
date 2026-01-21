package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewListallCmd creates a Cobra command for listing all todos
func NewListallCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "listall",
		Short: "List all todos including completed",
		Long: `Lists all todos including completed ones.

# list all tasks
togodo listall
`,
		Aliases: []string{"lsa"},
		RunE: func(command *cobra.Command, args []string) error {
			// Call business logic
			result, err := cmd.Listall(repo)
			if err != nil {
				return err
			}

			// Format output
			for _, todo := range result.Todos {
				fmt.Fprintf(command.OutOrStdout(), "%d %s\n", todo.LineNumber, todo.Text)
			}
			fmt.Fprintln(command.OutOrStdout(), "--")
			fmt.Fprintf(command.OutOrStdout(), "TODO: %d of %d tasks shown\n", result.TotalCount, result.TotalCount)
			return nil
		},
	}
}
