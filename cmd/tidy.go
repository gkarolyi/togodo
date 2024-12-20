package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/internal/todolib"

	"github.com/spf13/cobra"
)

// tidyCmd represents the add command
var tidyCmd = &cobra.Command{
	Use:   "tidy",
	Short: "Tidy up your todo.txt.",
	Long: `Cleans up your todo.txt by removing done tasks, and prints the tasks that were removed.

# tidy up your todo.txt
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
