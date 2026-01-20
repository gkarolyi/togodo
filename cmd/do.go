package cmd

import (
	"fmt"
	"strconv"

	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// parseLineNumbers converts CLI line numbers (1-based) to repository indices (0-based)
func parseLineNumbers(args []string) ([]int, error) {
	indices := make([]int, len(args))
	for i, arg := range args {
		lineNumber, err := strconv.Atoi(arg)
		if err != nil {
			return nil, fmt.Errorf("failed to convert arg to int: %w", err)
		}
		if lineNumber < 1 {
			return nil, fmt.Errorf("line number must be positive, got %d", lineNumber)
		}
		indices[i] = lineNumber - 1 // Convert to 0-based
	}
	return indices, nil
}

// NewDoCmd creates a new cobra command for toggling todos.
func NewDoCmd(repo todotxtlib.TodoRepository, presenter *cli.Presenter) *cobra.Command {
	return &cobra.Command{
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
			// Parse line numbers (convert from 1-based to 0-based)
			indices, err := parseLineNumbers(args)
			if err != nil {
				return err
			}

			// Toggle todos
			toggledTodos := make([]todotxtlib.Todo, 0, len(indices))
			for _, index := range indices {
				todo, err := repo.ToggleDone(index)
				if err != nil {
					return fmt.Errorf("failed to toggle todo at index %d: %w", index, err)
				}
				toggledTodos = append(toggledTodos, todo)
			}

			// Sort and save
			repo.Sort(nil)
			if err := repo.Save(); err != nil {
				return fmt.Errorf("failed to save todos: %w", err)
			}

			// Present
			for _, todo := range toggledTodos {
				presenter.Print(todo)
			}
			return nil
		},
	}
}
