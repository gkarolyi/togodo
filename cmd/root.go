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
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "togodo",
	Short: "A CLI tool for managing your todo.txt",
	Long:  `togodo is a CLI tool for managing your todo.txt file.`,

	// This is where the TUI will be called from eventually
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var TodoTxtPath string

func init() {
	// Try to open a todo.txt file in the current directory first
	if TodoTxtPath == "" {
		if _, err := os.Stat("todo.txt"); os.IsNotExist(err) {
			os.Exit(1)
		} else {
			TodoTxtPath = "todo.txt"
		}
	}

	// If that fails, try to open a todo.txt file in the directory specified by the TODO_TXT_PATH environment variable
	if TodoTxtPath == "" {
		envTodoTxtPath := os.Getenv("TODO_TXT_PATH")
		if _, err := os.Stat(envTodoTxtPath); os.IsNotExist(err) {
			os.Exit(1)
		} else {
			TodoTxtPath = envTodoTxtPath
		}
	}

	// Finally, try to open a todo.txt file in the user's home directory
	if TodoTxtPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			os.Exit(1)
		}
		TodoTxtPath = homeDir + "/todo.txt"
		if _, err := os.Stat(TodoTxtPath); os.IsNotExist(err) {
			os.Exit(1)
		}
	}

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
