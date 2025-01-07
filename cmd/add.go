package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/internal/todolib"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
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
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := todolib.New(TodoTxtPath)
		if err != nil {
			fmt.Println(err)
		}

		todos, err := repo.Add(args[0])
		if err != nil {
			fmt.Println(err)
		}
		for _, todo := range todos {
			todolib.Render(todo)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
