package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewTidyCmd creates a Cobra command for removing done tasks
func NewTidyCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "tidy",
		Short: "Remove completed tasks",
		Long: `Removes all completed tasks from the todo list.

# remove all done tasks
togodo tidy
`,
		Aliases: []string{"clean"},
		RunE: func(command *cobra.Command, args []string) error {
			// Call business logic
			result, err := cmd.Tidy(repo)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			if len(result.RemovedTodos) == 0 {
				fmt.Println("TODO: No completed tasks to remove.")
			} else {
				for _, todo := range result.RemovedTodos {
					fmt.Println(todo.Text)
				}
				fmt.Printf("TODO: %d completed task(s) removed.\n", len(result.RemovedTodos))
			}
			return nil
		},
	}
}
