package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewAppendCmd creates a Cobra command for appending to todos
func NewAppendCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "append ITEM# \"TEXT TO APPEND\"",
		Short: "Append text to a todo item",
		Long: `Appends text to the end of a todo item.

# append text to task 1
togodo append 1 "additional text"
`,
		Args:    cobra.MinimumNArgs(2),
		Aliases: []string{"app"},
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

			// Get text to append
			text := strings.Join(args[1:], " ")

			// Call business logic
			result, err := cmd.Append(repo, index, text)
			if err != nil {
				return err
			}

			// Format output
			fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumber, result.Todo.Text)
			return nil
		},
	}
}
