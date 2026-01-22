package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewPriCmd creates a Cobra command for setting priority
func NewPriCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "pri [LINE_NUMBER] [PRIORITY]",
		Short: "Set priority of a todo item",
		Long: `Sets the priority of a todo item (A, B, C, etc.).

# set task 1 to priority A
togodo pri 1 A
`,
		Args:    cobra.ExactArgs(2),
		Aliases: []string{"p"},
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

			// Get priority (normalize to uppercase)
			priority := strings.ToUpper(args[1])

			// Call business logic
			result, err := cmd.SetPriority(repo, []int{index}, priority)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			for i, todo := range result.UpdatedTodos {
				fmt.Fprintf(command.OutOrStdout(), "%d %s\n", lineNum, todo.Text)

				// Show appropriate message based on whether it was a re-prioritization
				oldPriority := result.OldPriorities[i]
				if oldPriority != "" {
					// Re-prioritization: task already had a priority
					fmt.Fprintf(command.OutOrStdout(), "TODO: %d re-prioritized from (%s) to (%s).\n", lineNum, oldPriority, priority)
				} else {
					// First-time prioritization
					fmt.Fprintf(command.OutOrStdout(), "TODO: %d prioritized (%s).\n", lineNum, priority)
				}
			}
			return nil
		},
	}
}
