package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewReplaceCmd creates a Cobra command for replacing todos
func NewReplaceCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "replace ITEM# \"UPDATED ITEM\"",
		Short: "Replace a todo item with new text",
		Long: `Replaces the entire text of a todo item.

# replace task 1 with new text
togodo replace 1 "new task text"
`,
		Args: cobra.ExactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			// Parse line number (1-based)
			lineNum, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid line number: %s", args[0])
			}

			// Convert to 0-based index
			index := lineNum - 1

			// Get new text (join remaining args)
			newText := strings.Join(args[1:], " ")

			// Call business logic
			result, err := cmd.Replace(repo, index, newText)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumber, result.OldTodo.Text)
			fmt.Fprintf(command.OutOrStdout(), "TODO: Replaced task with:\n")
			fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumber, result.NewTodo.Text)
			return nil
		},
	}
}
