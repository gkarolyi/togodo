package cmd

import (
	"fmt"
	"strconv"

	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

func executeDo(repo todotxtlib.TodoRepository, args []string) ([]todotxtlib.Todo, error) {
	var toggledTodos []todotxtlib.Todo

	for _, arg := range args {
		lineNumber, err := strconv.Atoi(arg)
		if err != nil {
			return nil, fmt.Errorf("failed to convert arg to int: %w", err)
		}
		todo, err := repo.ToggleDone(lineNumber - 1)
		if err != nil {
			return nil, fmt.Errorf("failed to toggle todo at line %d: %w", lineNumber, err)
		}
		toggledTodos = append(toggledTodos, todo)
	}

	repo.SortDefault()
	err := repo.Save()
	if err != nil {
		return nil, err
	}

	return toggledTodos, nil
}

// NewDoCmd creates a new cobra command for toggling todos.
func NewDoCmd(repo todotxtlib.TodoRepository) *cobra.Command {
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
			_, err := executeDo(repo, args)
			return err
		},
	}
}
