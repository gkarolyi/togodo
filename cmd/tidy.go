package cmd

import (
	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewTidyCmd creates a new cobra command for tidying up the todo list.
func NewTidyCmd(service todotxtlib.TodoService, presenter *cli.Presenter) *cobra.Command {
	return &cobra.Command{
		Use:   "tidy",
		Short: "Tidy up your todo.txt by removing done tasks",
		Long: `Cleans up your todo.txt by removing done tasks, and prints the tasks that were removed.

# tidy up your todo.txt
togodo tidy`,
		Args:    cobra.NoArgs,
		Aliases: []string{"clean"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Business logic - delegated to service
			todos, err := service.RemoveDoneTodos()
			if err != nil {
				return err
			}

			// Presentation logic - handled by presenter
			return presenter.PrintList(todos)
		},
	}
}
