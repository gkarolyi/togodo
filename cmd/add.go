package cmd

import (
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

func executeAdd(repo todotxtlib.TodoRepository, args []string) ([]todotxtlib.Todo, error) {
	var addedTodos []todotxtlib.Todo

	for _, todoText := range args {
		todo, err := repo.Add(todoText)
		if err != nil {
			return nil, err
		}
		addedTodos = append(addedTodos, todo)
	}

	repo.SortDefault()
	err := repo.Save()
	if err != nil {
		return nil, err
	}

	return addedTodos, nil
}

// NewAddCmd creates a new cobra command for adding todos.
func NewAddCmd(repo todotxtlib.TodoRepository) *cobra.Command {
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
			_, err := executeAdd(repo, args)
			return err
		},
	}
}
