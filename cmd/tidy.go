package cmd

import (
	"github.com/spf13/cobra"
)

func executeTidy(base *BaseCommand) error {
	// Get all done todos before removing them
	doneTodos, err := base.Repository.ListDone()
	if err != nil {
		return err
	}

	// Remove done todos from the repository
	// We need to iterate backwards to avoid index shifting issues
	allTodos, err := base.Repository.ListAll()
	if err != nil {
		return err
	}

	removedCount := 0
	for i := len(allTodos) - 1; i >= 0; i-- {
		if allTodos[i].Done {
			_, err := base.Repository.Remove(i)
			if err != nil {
				return err
			}
			removedCount++
		}
	}

	// Sort the remaining todos
	base.Repository.SortDefault()

	// Save the repository
	err = base.Save()
	if err != nil {
		return err
	}

	// Print the removed todos if any
	if len(doneTodos) > 0 {
		return base.PrintList(doneTodos)
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
		base := NewDefaultBaseCommand()
		return executeTidy(base)
	},
}

func init() {
	rootCmd.AddCommand(tidyCmd)
}
