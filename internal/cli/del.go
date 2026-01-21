package cli

import (
	"fmt"
	"strconv"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewDelCmd creates a Cobra command for deleting todos
func NewDelCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "del ITEM#",
		Short: "Delete a todo item",
		Long: `Deletes a todo item from the list.

# delete task 1
togodo del 1
`,
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"rm"},
		RunE: func(command *cobra.Command, args []string) error {
			// Parse line number (1-based)
			lineNum, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid line number: %s", args[0])
			}

			// Find the array index for this line number
			index := repo.FindIndexByLineNumber(lineNum)
			if index == -1 {
				return fmt.Errorf("TODO: No task %d.", lineNum)
			}

			// Call business logic
			result, err := cmd.Del(repo, index)
			if err != nil {
				return err
			}

			// Format output
			fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumber, result.DeletedTodo.Text)
			fmt.Fprintf(command.OutOrStdout(), "TODO: %d deleted.\n", result.LineNumber)
			return nil
		},
	}
}
