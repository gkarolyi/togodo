package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewTidyCmd creates a new cobra command for tidying up the todo list.
func NewTidyCmd(repo todotxtlib.TodoRepository, presenter *cli.Presenter) *cobra.Command {
	return &cobra.Command{
		Use:   "tidy",
		Short: "Tidy up your todo.txt by removing done tasks",
		Long: `Cleans up your todo.txt by removing done tasks, and prints the tasks that were removed.

# tidy up your todo.txt
togodo tidy`,
		Args:    cobra.NoArgs,
		Aliases: []string{"clean"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get done todos before removing
			doneFilter := todotxtlib.Filter{Done: "true"}
			doneTodos, err := repo.Filter(doneFilter)
			if err != nil {
				return fmt.Errorf("failed to list done todos: %w", err)
			}

			// Get all todos to iterate
			allTodos, err := repo.ListAll()
			if err != nil {
				return fmt.Errorf("failed to list all todos: %w", err)
			}

			// Remove backwards to avoid index shifting
			for i := len(allTodos) - 1; i >= 0; i-- {
				if allTodos[i].Done {
					if _, err := repo.Remove(i); err != nil {
						return fmt.Errorf("failed to remove todo at index %d: %w", i, err)
					}
				}
			}

			// Sort and save
			repo.Sort(nil)
			if err := repo.Save(); err != nil {
				return fmt.Errorf("failed to save todos: %w", err)
			}

			// Present
			return presenter.PrintList(doneTodos)
		},
	}
}
