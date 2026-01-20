package cmd

import (
	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewAddCmd creates a new cobra command for adding todos.
func NewAddCmd(service todotxtlib.TodoService, presenter *cli.Presenter) *cobra.Command {
	return &cobra.Command{
		Use:   "add [TASK]",
		Short: "Add a new todo item to the list",
		Long: `Adds a new task to the list and prints the newly added task.
If [TASK] contains multiple lines, each line is added as a separate task.

# add "Buy milk" to the list
togodo add "Buy milk"

# add multiple tasks to the list
togodo add "Buy milk
Buy eggs
Buy bread"
`,

		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"a"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Business logic - delegated to service
			todos, err := service.AddTodos(args)
			if err != nil {
				return err
			}

			// Presentation logic - handled by presenter
			for _, todo := range todos {
				presenter.Print(todo)
			}
			return nil
		},
	}
}
