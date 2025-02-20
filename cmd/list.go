package cmd

import (
	"github.com/gkarolyi/togodo/internal/todolib"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [FILTER]",
	Short: "List and filter items in your todo.txt",
	Long: `Lists tasks sorted in order of priority, with done items at the bottom of the list. Tasks can optionally be filtered
by passing an optional [FILTER] argument. If no filter is passed, list shows all items in your todo.txt file. Tasks are shown
with a line number to allow you to easily refer to them. For example:

# list all items in your todo.txt file
togodo list

# list all items in your todo.txt file that contain the string '@work'
togodo list '@work'
`,
	Aliases: []string{"ls", "l"},
	Args:    cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		todolib.List(TodoTxtPath, args)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
