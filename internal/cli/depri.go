package cli

import (
	"fmt"
	"strconv"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewDepriCmd creates a Cobra command for removing priority
func NewDepriCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "depri ITEM# [ITEM#...]",
		Short: "Remove priority from todo items",
		Long: `Removes the priority from one or more todo items.

# remove priority from task 1
togodo depri 1

# remove priority from multiple tasks
togodo depri 1 2 3
`,
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"dp"},
		RunE: func(command *cobra.Command, args []string) error {
			var lastErr error
			// Process each line number
			for _, arg := range args {
				// Parse line number (1-based)
				lineNum, err := strconv.Atoi(arg)
				if err != nil {
					return fmt.Errorf("invalid line number: %s", arg)
				}

				// Find the array index for this line number
				index := repo.FindIndexByLineNumber(lineNum)
				if index == -1 {
					lastErr = fmt.Errorf("TODO: No task %d.", lineNum)
					fmt.Fprintf(command.OutOrStderr(), "%v\n", lastErr)
					continue
				}

				// Call business logic
				result, err := cmd.Depri(repo, index)
				if err != nil {
					lastErr = err
					fmt.Fprintf(command.OutOrStderr(), "%v\n", err)
					continue
				}

				// Format output
				fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumber, result.Todo.Text)
				fmt.Fprintf(command.OutOrStdout(), "TODO: %d deprioritized.\n", result.LineNumber)
			}
			// Return error if any command failed
			return lastErr
		},
	}
}
