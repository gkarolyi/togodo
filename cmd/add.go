package cmd

import (
	"github.com/spf13/cobra"
)

func executeAdd(base *BaseCommand, args []string) error {
	added := make([]string, 0)
	for _, todoText := range args {
		todo, err := base.Repository.Add(todoText)
		if err != nil {
			return err
		}
		added = append(added, todo.Text)
	}

	base.Repository.SortDefault()

	return base.Save()
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
		base := NewDefaultBaseCommand()
		return executeAdd(base, args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
