package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewListconCmd creates a Cobra command for listing contexts
func NewListconCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "listcon",
		Short: "List all contexts",
		Long: `Lists all contexts (@context) found in todos.

# list all contexts
togodo listcon
`,
		Aliases: []string{"lsc"},
		RunE: func(command *cobra.Command, args []string) error {
			// Call business logic
			result, err := cmd.Listcon(repo)
			if err != nil {
				return err
			}

			// Format output
			for _, context := range result.Contexts {
				fmt.Fprintln(command.OutOrStdout(), context)
			}
			return nil
		},
	}
}
