package cmd

import (
	"fmt"
	"strconv"

	"github.com/gkarolyi/togodo/internal/injector"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

func executeDo(repo *todotxtlib.Repository, args []string) error {
	for _, arg := range args {
		lineNumber, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("failed to convert arg to int: %w", err)
		}
		_, err = repo.ToggleDone(lineNumber - 1)
		if err != nil {
			return fmt.Errorf("failed to toggle todo at line %d: %w", lineNumber, err)
		}
	}

	repo.SortDefault()
	return repo.Save()
}

var doCmd = &cobra.Command{
	Use:   "do [LINE NUMBER]...",
	Short: "Toggle the done status of a todo item",
	Long: `Marks a task as done or not done depending on its current status, and prints the toggled task.
If [LINE_NUMBER] contains multiple line numbers, each todo will be toggled.

# toggle the done status of the task on line 1
togodo do 1

# toggle the done status of the tasks on lines 1, 2, and 3
togodo do 1 2 3
`,

	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"x"},
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := injector.CreateRepository()
		if err != nil {
			return err
		}
		return executeDo(repo, args)
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
