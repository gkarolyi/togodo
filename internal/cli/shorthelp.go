package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewShorthelpCmd creates a Cobra command for showing short help
func NewShorthelpCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "shorthelp",
		Short: "Show condensed help",
		Long:  `Shows a condensed list of available commands.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Show short version of help - just command names and short descriptions
			fmt.Fprintln(cmd.OutOrStdout(), "Usage: togodo [command]")
			fmt.Fprintln(cmd.OutOrStdout(), "")
			fmt.Fprintln(cmd.OutOrStdout(), "Available commands:")

			for _, command := range rootCmd.Commands() {
				// Skip hidden commands and help itself
				if command.Hidden || command.Name() == "completion" {
					continue
				}

				// Show command name and short description
				fmt.Fprintf(cmd.OutOrStdout(), "  %-12s %s\n", command.Name(), command.Short)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "")
			fmt.Fprintln(cmd.OutOrStdout(), "Use \"togodo [command] --help\" for more information about a command.")
		},
	}
}
