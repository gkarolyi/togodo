package cmd

import (
	"fmt"
	"togodo/internal/todolib"

	"github.com/spf13/cobra"
)

// tidyCmd represents the add command
var tidyCmd = &cobra.Command{
	Use:   "tidy",
	Short: "Tidy up the todo list",
	Long: `Removes all done items from the list. For example:

	togodo tidy`,

	Args:    cobra.NoArgs,
	Aliases: []string{"clean"},
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := todolib.New(TodoTxtPath)
		if err != nil {
			fmt.Println(err)
		}
		todos, err := repo.Tidy()
		if err != nil {
			fmt.Println(err)
		}
		for _, todo := range todos {
			fmt.Println(todo.Text)
		}
	},
}

func init() {
	rootCmd.AddCommand(tidyCmd)
}
