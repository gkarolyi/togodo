package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/internal/todolib"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new todo item to the list",
	Long: `Add a new todo item to the list. For example:

	togodo add "Buy milk"`,

	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"a"},
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := todolib.New(TodoTxtPath)
		if err != nil {
			fmt.Println(err)
		}
		todo, err := repo.Add(args[0])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(todo.Text)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
