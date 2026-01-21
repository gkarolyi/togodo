package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewDoCmd creates a Cobra command for toggling todo completion
func NewDoCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "do ITEM#[,ITEM#,...]",
		Short: "Mark todo item(s) as done (or undone)",
		Long: `Toggles the completion status of one or more todo items.

# mark task 1 as done
togodo do 1

# mark multiple tasks as done
togodo do 1 3 5

# mark multiple tasks using comma separator
togodo do 1,3,5
`,
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"x"},
		RunE: func(command *cobra.Command, args []string) error {
			// Parse line numbers - handle both space-separated and comma-separated
			var lineNumbers []int
			for _, arg := range args {
				// Check if this arg contains commas
				if strings.Contains(arg, ",") {
					// Split by comma and parse each
					parts := strings.Split(arg, ",")
					for _, part := range parts {
						lineNum, err := strconv.Atoi(strings.TrimSpace(part))
						if err != nil {
							return fmt.Errorf("invalid line number: %s", part)
						}
						lineNumbers = append(lineNumbers, lineNum)
					}
				} else {
					// Single line number
					lineNum, err := strconv.Atoi(arg)
					if err != nil {
						return fmt.Errorf("invalid line number: %s", arg)
					}
					lineNumbers = append(lineNumbers, lineNum)
				}
			}

			// Convert line numbers to array indices
			indices := make([]int, 0, len(lineNumbers))
			for _, lineNum := range lineNumbers {
				index := repo.FindIndexByLineNumber(lineNum)
				if index == -1 {
					return fmt.Errorf("TODO: No task %d.", lineNum)
				}
				indices = append(indices, index)
			}

			// Call business logic
			result, err := cmd.Do(repo, indices)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			for _, todo := range result.ToggledTodos {
				fmt.Fprintf(command.OutOrStdout(), "%d %s\n", todo.LineNumber, todo.Text)
				if todo.Done {
					fmt.Fprintf(command.OutOrStdout(), "TODO: %d marked as done.\n", todo.LineNumber)
				} else {
					fmt.Fprintf(command.OutOrStdout(), "TODO: %d marked as TODO.\n", todo.LineNumber)
				}
			}
			return nil
		},
	}
}
