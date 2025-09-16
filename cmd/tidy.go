package cmd

import (
	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/internal/injector"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

func executeTidy(repo *todotxtlib.Repository, presenter *cli.Presenter) error {
	// Get all done todos before removing them
	doneTodos, err := repo.ListDone()
	if err != nil {
		return err
	}

	// Remove done todos from the repository
	// We need to iterate backwards to avoid index shifting issues
	allTodos, err := repo.ListAll()
	if err != nil {
		return err
	}

	for i := len(allTodos) - 1; i >= 0; i-- {
		if allTodos[i].Done {
			_, err := repo.Remove(i)
			if err != nil {
				return err
			}
		}
	}

	// Sort the remaining todos
	repo.SortDefault()

	// Save the repository
	err = repo.Save()
	if err != nil {
		return err
	}

	// Print the removed todos if any
	if len(doneTodos) > 0 {
		return presenter.PrintList(doneTodos)
	}

	return nil
}

var tidyCmd = &cobra.Command{
	Use:   "tidy",
	Short: "Tidy up your todo.txt by removing done tasks",
	Long: `Cleans up your todo.txt by removing done tasks, and prints the tasks that were removed.

# tidy up your todo.txt
togodo tidy`,
	Args:    cobra.NoArgs,
	Aliases: []string{"clean"},
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := injector.CreateRepository()
		if err != nil {
			return err
		}
		presenter := injector.CreateCLIPresenter()
		return executeTidy(repo, presenter)
	},
}

func init() {
	rootCmd.AddCommand(tidyCmd)
}
