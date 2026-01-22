package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewDeduplicateCmd creates a Cobra command for removing duplicate tasks
func NewDeduplicateCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "deduplicate",
		Short: "Remove duplicate tasks",
		Long: `Removes duplicate tasks from the todo list, keeping only the first occurrence of each unique task.

Duplicate detection is case-sensitive and based on exact text match.

# remove all duplicate tasks
togodo deduplicate
`,
		Aliases: []string{"dedup", "dedupe"},
		RunE: func(command *cobra.Command, args []string) error {
			// Call business logic
			result, err := cmd.Deduplicate(repo)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			if result.RemovedCount == 0 {
				fmt.Fprintln(command.OutOrStdout(), "TODO: No duplicates found.")
			} else {
				fmt.Fprintf(command.OutOrStdout(), "TODO: %d duplicate task(s) removed.\n", result.RemovedCount)
			}
			return nil
		},
	}
}
