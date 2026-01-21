package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewPrependCmd creates a Cobra command for prepending to todos
func NewPrependCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "prepend ITEM# \"TEXT TO PREPEND\"",
		Short: "Prepend text to a todo item",
		Long: `Prepends text to the beginning of a todo item, preserving priority.

# prepend text to task 2
togodo prepend 2 "really"
`,
		Args:    cobra.MinimumNArgs(2),
		Aliases: []string{"prep"},
		RunE: func(command *cobra.Command, args []string) error {
			// Parse line number (1-based)
			lineNum, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid line number: %s", args[0])
			}

			// Convert to 0-based index
			index := lineNum - 1

			// Get text to prepend
			text := strings.Join(args[1:], " ")

			// Call business logic
			result, err := cmd.Prepend(repo, index, text)
			if err != nil {
				return err
			}

			// Format output
			fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumber, result.Todo.Text)
			return nil
		},
	}
}
