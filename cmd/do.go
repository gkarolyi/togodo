/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"
	"togodo/internal/todolib"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := todolib.New(TodoTxtPath)
		if err != nil {
			fmt.Println(err)
		}
		var lineNumbers []int
		for _, arg := range args {
			lineNumber, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Invalid argument:", arg)
				continue
			}
			lineNumbers = append(lineNumbers, lineNumber)
		}
		todos, err := repo.Toggle(lineNumbers)
		if err != nil {
			fmt.Println(err)
		}
		for _, todo := range todos {
			fmt.Println(todo.Text)
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
