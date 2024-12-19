package cmd

import (
	"fmt"
	"strconv"

	"github.com/gkarolyi/togodo/internal/todolib"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Toggle the done status of a todo item",
	Long: `Toggle the done status of one or more todo items using their line numbers. For example:

	togodo do 1 2 3`,

	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"x"},
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := todolib.New(TodoTxtPath)
		if err != nil {
			fmt.Println(err)
		}
		var lineNumbers []int
		for _, arg := range args {
			lineNumber, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Invalid argument:", arg)
				continue
			}
			lineNumbers = append(lineNumbers, lineNumber)
		}
		todos, err := repo.Toggle(lineNumbers)
		if err != nil {
			fmt.Println(err)
		}
		for _, todo := range todos {
			fmt.Println(todo.Text)
		}
	},
}
var xCmd = doCmd

func init() {
	rootCmd.AddCommand(doCmd)
	rootCmd.AddCommand(xCmd)
}
