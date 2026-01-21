package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewListallCmd creates a Cobra command for listing all todos
func NewListallCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	listallCmd := &cobra.Command{
		Use:   "listall",
		Short: "List all todos including completed",
		Long: `Lists all todos including completed ones with color highlighting.

# list all tasks
togodo listall

# list all tasks without colors
togodo listall --plain
`,
		Aliases: []string{"lsa"},
		RunE: func(command *cobra.Command, args []string) error {
			// Get plain flag
			plainMode, _ := command.Flags().GetBool("plain")

			// Call business logic
			result, err := cmd.Listall(repo)
			if err != nil {
				return err
			}

			// Choose formatter based on plain flag
			var formatter TodoFormatter
			if plainMode {
				formatter = NewPlainFormatter()
			} else {
				formatter = NewLipglossFormatter()
			}

			// Format output with highlighting
			for _, todo := range result.Todos {
				formatted := formatter.Format(todo)
				fmt.Fprintf(command.OutOrStdout(), "%d %s\n", todo.LineNumber, formatted)
			}
			fmt.Fprintln(command.OutOrStdout(), "--")
			fmt.Fprintf(command.OutOrStdout(), "TODO: %d of %d tasks shown\n", result.TotalCount, result.TotalCount)
			return nil
		},
	}

	// Add flags
	listallCmd.Flags().BoolP("plain", "p", false, "Plain output without colors")

	return listallCmd
}
