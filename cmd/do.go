package cmd

import (
	"fmt"
	"strconv"

	"github.com/gkarolyi/togodo/internal/todolib"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
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
			todolib.Render(todo)
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
