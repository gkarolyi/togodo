package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/spf13/cobra"
)

// NewConfigCmd creates a Cobra command for configuration management
func NewConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config [key] [value]",
		Short: "Manage configuration settings",
		Long: `Manage configuration settings for togodo.

# Show a configuration value
togodo config todo_txt_path

# Set a configuration value
togodo config todo_txt_path ~/todo.txt

# List all configuration
togodo config
`,
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) == 0 {
				// No args: list all config (Task 3)
				return fmt.Errorf("config list not yet implemented")
			} else if len(args) == 1 {
				// One arg: read config value (Task 1)
				key := args[0]
				result, err := cmd.ConfigRead(key)
				if err != nil {
					return err
				}

				fmt.Fprintln(command.OutOrStdout(), result.Value)
				return nil
			} else if len(args) == 2 {
				// Two args: write config value (Task 2)
				return fmt.Errorf("config write not yet implemented")
			} else {
				return fmt.Errorf("too many arguments")
			}
		},
	}

	return configCmd
}
