package cmd

import (
	"github.com/gkarolyi/togodo/internal/injector"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

func executeAdd(repo *todotxtlib.Repository, args []string) error {
	for _, todoText := range args {
		_, err := repo.Add(todoText)
		if err != nil {
			return err
		}
	}

	repo.SortDefault()
	return repo.Save()
}

var addCmd = &cobra.Command{
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
		repo, err := injector.CreateRepository()
		if err != nil {
			return err
		}
		return executeAdd(repo, args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
