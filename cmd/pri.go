package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// parsePriorityArgs parses CLI arguments for the pri command
// Returns line indices (0-based) and the priority string
func parsePriorityArgs(args []string) ([]int, string, error) {
	if len(args) < 2 {
		return nil, "", fmt.Errorf("pri requires at least a line number and priority")
	}

	priority := args[len(args)-1]
	lineNumberArgs := args[:len(args)-1]

	indices, err := parseLineNumbers(lineNumberArgs)
	if err != nil {
		return nil, "", err
	}

	return indices, priority, nil
}

// NewPriCmd creates a new cobra command for setting priority.
func NewPriCmd(repo todotxtlib.TodoRepository, presenter *cli.Presenter) *cobra.Command {
	return &cobra.Command{
		Use:   "pri [LINE NUMBER]...",
		Short: "Set the priority of a todo item",
		Long: `Set the priority of a todo item.

# set the priority of the todo on line 1 to A
togodo pri 1 A

# set the priority of the todos on lines 1, 2, and 3 to B
togodo pri 1 2 3 B
`,

		Args:    cobra.MinimumNArgs(2),
		Aliases: []string{"p"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse priority arguments
			indices, priority, err := parsePriorityArgs(args)
			if err != nil {
				return err
			}

			// Set priorities
			updatedTodos := make([]todotxtlib.Todo, 0, len(indices))
			for _, index := range indices {
				todo, err := repo.SetPriority(index, priority)
				if err != nil {
					return fmt.Errorf("failed to set priority for todo at index %d: %w", index, err)
				}
				updatedTodos = append(updatedTodos, todo)
			}

			// Note: Pri command doesn't sort - preserves user's order
			if err := repo.Save(); err != nil {
				return fmt.Errorf("failed to save todos: %w", err)
			}

			// Present
			for _, todo := range updatedTodos {
				presenter.Print(todo)
			}
			return nil
		},
	}
}
