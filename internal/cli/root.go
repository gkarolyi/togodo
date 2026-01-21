/*
Copyright © 2024 Gergely Karolyi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cli

import (
	"fmt"
	"os"

	"github.com/gkarolyi/togodo/internal/config"
	"github.com/gkarolyi/togodo/internal/tui"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command and its subcommands, injecting dependencies.
func NewRootCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "togodo",
		Short: "A CLI tool for managing your todo.txt",
		Long:  `togodo is a CLI tool for managing your todo.txt file.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := tui.Run(repo)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("file", "f", "", "Specify the todo.txt file to use")

	// Set up persistent pre-run to handle --file flag globally
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if file, _ := cmd.Flags().GetString("file"); file != "" {
			config.SetTodoTxtPath(file)
		}
	}

	// Add subcommands from internal/cli
	rootCmd.AddCommand(NewAddCmd(repo))
	rootCmd.AddCommand(NewListCmd(repo))
	rootCmd.AddCommand(NewListpriCmd(repo))
	rootCmd.AddCommand(NewListconCmd(repo))
	rootCmd.AddCommand(NewDoCmd(repo))
	rootCmd.AddCommand(NewPriCmd(repo))
	rootCmd.AddCommand(NewTidyCmd(repo))
	rootCmd.AddCommand(NewReplaceCmd(repo))
	// TODO: Config command needs to be migrated
	// rootCmd.AddCommand(NewConfigCmd(presenter))

	return rootCmd
}

// initConfig reads in config file and ENV variables.
func initConfig() {
	if err := config.InitConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
		os.Exit(1)
	}
}
