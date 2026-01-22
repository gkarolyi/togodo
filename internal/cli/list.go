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

# list without colors
togodo -p list
togodo list --plain
`,
		Aliases: []string{"l", "ls"},
		RunE: func(command *cobra.Command, args []string) error {
			// Check for plain mode (either local or global flag)
			plainMode, _ := command.Flags().GetBool("plain")
			if !plainMode {
				// Check parent (root) flag
				if parent := command.Parent(); parent != nil {
					plainMode, _ = parent.PersistentFlags().GetBool("plain")
				}
			}

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

			// Choose formatter based on plain flag
			var formatter TodoFormatter
			if plainMode {
				formatter = NewPlainFormatter()
			} else {
				formatter = NewLipglossFormatter()
			}

			// Format output with optional highlighting
			for i, todo := range result.Todos {
				formatted := formatter.Format(todo)
				fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumbers[i], formatted)
			}
			fmt.Fprintln(command.OutOrStdout(), "--")
			fmt.Fprintf(command.OutOrStdout(), "TODO: %d of %d tasks shown\n", result.ShownCount, result.TotalCount)
			return nil
		},
	}
}
