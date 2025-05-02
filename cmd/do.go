package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func executeDo(base *BaseCommand, args []string) error {
	for _, arg := range args {
		lineNumber, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("failed to convert arg to int: %w", err)
		}
		base.Repository.ToggleDone(lineNumber - 1)
	}

	base.Repository.SortDefault()

	return base.Write()
}

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
	RunE: func(cmd *cobra.Command, args []string) error {
		base := NewDefaultBaseCommand()
		return executeDo(base, args)
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
