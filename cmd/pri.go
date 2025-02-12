package cmd

import (
	"fmt"
	"strconv"

	"github.com/gkarolyi/togodo/internal/todolib"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var priCmd = &cobra.Command{
	Use:   "pri [LINE NUMBER]...",
	Short: "Set the priority of a todo item",
	Long: `Set the priority of a todo item.

# set the priority of the todo on line 1 to A
togodo pri 1 A

# set the priority of the todos on lines 1, 2, and 3 to B
togodo pri 1 2 3 B
`,

	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"p"},
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := todolib.New(TodoTxtPath)
		if err != nil {
			fmt.Println(err)
		}

		var lineNumbers []int
		var priority string
		for i, arg := range args {
			lineNumber, err := strconv.Atoi(arg)
			if err != nil {
				if i == len(args)-1 {
					priority = arg
				} else {
					fmt.Println("Invalid argument:", arg)
				}
				continue
			}
			lineNumbers = append(lineNumbers, lineNumber)
		}

		todos, err := repo.SetPriority(lineNumbers, priority)
		if err != nil {
			fmt.Println(err)
		}
		for _, todo := range todos {
			todolib.Render(todo)
		}
	},
}

func init() {
	rootCmd.AddCommand(priCmd)
}
