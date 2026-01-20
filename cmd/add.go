package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewAddCmd creates a new cobra command for adding todos.
func NewAddCmd(repo todotxtlib.TodoRepository, presenter *cli.Presenter) *cobra.Command {
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
			// Add todos
			addedTodos := make([]todotxtlib.Todo, 0, len(args))
			for _, text := range args {
				todo, err := repo.Add(text)
				if err != nil {
					return fmt.Errorf("failed to add todo: %w", err)
				}
				addedTodos = append(addedTodos, todo)
			}

			// Sort and save
			repo.Sort(nil)
			if err := repo.Save(); err != nil {
				return fmt.Errorf("failed to save todos: %w", err)
			}

			// Present
			for _, todo := range addedTodos {
				presenter.Print(todo)
			}
			return nil
		},
	}
}
