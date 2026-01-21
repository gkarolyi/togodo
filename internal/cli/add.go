package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewAddCmd creates a Cobra command for adding todos
func NewAddCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "add [TASK]",
		Short: "Add a new todo item to the list",
		Long: `Adds a new task to the list and prints the newly added task.

# add "Buy milk" to the list
togodo add "Buy milk"

# add multiple words
togodo add Buy milk and eggs
`,
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"a"},
		RunE: func(command *cobra.Command, args []string) error {
			// Check if auto-dating is enabled
			autoDate := viper.GetBool("auto_add_creation_date")

			// Call business logic
			result, err := cmd.Add(repo, args, autoDate)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumber, result.Todo.Text)
			fmt.Fprintf(command.OutOrStdout(), "TODO: %d added.\n", result.LineNumber)
			return nil
		},
	}
}
