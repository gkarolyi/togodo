package cli

import (
	"fmt"
	"strconv"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewDoCmd creates a Cobra command for toggling todo completion
func NewDoCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "do [LINE_NUMBER]",
		Short: "Mark a todo item as done (or undone)",
		Long: `Toggles the completion status of a todo item.

# mark task 1 as done
togodo do 1
`,
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"x"},
		RunE: func(command *cobra.Command, args []string) error {
			// Parse line number (1-based)
			lineNum, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid line number: %s", args[0])
			}

			// Convert to 0-based index
			index := lineNum - 1

			// Call business logic
			result, err := cmd.Do(repo, []int{index})
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			for _, todo := range result.ToggledTodos {
				fmt.Fprintf(command.OutOrStdout(), "%d %s\n", lineNum, todo.Text)
				if todo.Done {
					fmt.Fprintf(command.OutOrStdout(), "TODO: %d marked as done.\n", lineNum)
				} else {
					fmt.Fprintf(command.OutOrStdout(), "TODO: %d marked as TODO.\n", lineNum)
				}
			}
			return nil
		},
	}
}
