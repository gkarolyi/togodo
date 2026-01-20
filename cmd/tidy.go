package cmd

import (
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

func executeTidy(repo todotxtlib.TodoRepository) ([]todotxtlib.Todo, error) {
	// Get all done todos before removing them
	doneTodos, err := repo.ListDone()
	if err != nil {
		return nil, err
	}

	// Remove done todos from the repository
	// We need to iterate backwards to avoid index shifting issues
	allTodos, err := repo.ListAll()
	if err != nil {
		return nil, err
	}

	for i := len(allTodos) - 1; i >= 0; i-- {
		if allTodos[i].Done {
			_, err := repo.Remove(i)
			if err != nil {
				return nil, err
			}
		}
	}

	// Sort the remaining todos
	repo.SortDefault()

	// Save the repository
	err = repo.Save()
	if err != nil {
		return nil, err
	}

	return doneTodos, nil
}

// NewTidyCmd creates a new cobra command for tidying up the todo list.
func NewTidyCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "tidy",
		Short: "Tidy up your todo.txt by removing done tasks",
		Long: `Cleans up your todo.txt by removing done tasks, and prints the tasks that were removed.

# tidy up your todo.txt
togodo tidy`,
		Args:    cobra.NoArgs,
		Aliases: []string{"clean"},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := executeTidy(repo)
			return err
		},
	}
}
