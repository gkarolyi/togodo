package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewListprojCmd creates a Cobra command for listing projects
func NewListprojCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "listproj",
		Short: "List all projects",
		Long: `Lists all projects (+project) found in todos.

# list all projects
togodo listproj
`,
		Aliases: []string{"lsprj"},
		RunE: func(command *cobra.Command, args []string) error {
			// Call business logic
			result, err := cmd.Listproj(repo)
			if err != nil {
				return err
			}

			// Format output
			for _, project := range result.Projects {
				fmt.Fprintln(command.OutOrStdout(), project)
			}
			return nil
		},
	}
}
