package cmd

import (
	"fmt"
	"strconv"

	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

func executePri(repo *todotxtlib.Repository, args []string) error {
	priority := args[len(args)-1]
	for _, arg := range args[:len(args)-1] {
		lineNumber, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("failed to convert arg to int: %w", err)
		}
		_, err = repo.SetPriority(lineNumber-1, priority)
		if err != nil {
			return fmt.Errorf("failed to set priority for todo at line %d: %w", lineNumber, err)
		}
	}
	return repo.Save()
}

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
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := createRepository()
		if err != nil {
			return err
		}
		return executePri(repo, args)
	},
}

func init() {
	rootCmd.AddCommand(priCmd)
}
