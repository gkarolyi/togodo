package cmd

import (
	"log"
	"strings"
	"togodo/internal/todolib"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all todo items in your todo.txt",
	Long: `For example:

	togodo list`,
	Aliases: []string{"ls", "l"},
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := todolib.New(TodoTxtPath)
		var todos []todolib.Todo

		if err != nil {
			log.Fatal(err)
		}

		if len(args) == 0 {
			todos = repo.All()
		} else {
			query := strings.Join(args, " ")
			todos = repo.Filter(query)
		}

		for number, todo := range todos {
			todolib.Render(number, todo)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
