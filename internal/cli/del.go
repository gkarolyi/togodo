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
		Use:   "del ITEM#...",
		Short: "Delete todo item(s)",
		Long: `Deletes one or more todo items from the list.

# delete task 1
togodo del 1

# delete multiple tasks
togodo del 1 3 5
`,
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"rm"},
		RunE: func(command *cobra.Command, args []string) error {
			// Parse all line numbers (1-based)
			indices := make([]int, 0, len(args))
			for _, arg := range args {
				lineNum, err := strconv.Atoi(arg)
				if err != nil {
					return fmt.Errorf("invalid line number: %s", arg)
				}

				// Find the array index for this line number
				index := repo.FindIndexByLineNumber(lineNum)
				if index == -1 {
					return fmt.Errorf("TODO: No task %d.", lineNum)
				}
				indices = append(indices, index)
			}

			// Call business logic
			result, err := cmd.Del(repo, indices)
			if err != nil {
				return err
			}

			// Format output - show each deleted task
			for _, deleted := range result.DeletedTodos {
				fmt.Fprintf(command.OutOrStdout(), "%d %s\n", deleted.LineNumber, deleted.Text)
			}
			if len(result.DeletedTodos) == 1 {
				fmt.Fprintf(command.OutOrStdout(), "TODO: %d deleted.\n", result.DeletedTodos[0].LineNumber)
			} else {
				fmt.Fprintf(command.OutOrStdout(), "TODO: %d tasks deleted.\n", len(result.DeletedTodos))
			}
			return nil
		},
	}
}
