package cli

import (
	"fmt"
	"strconv"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewDelCmd creates a Cobra command for deleting todos
func NewDelCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "del ITEM#...",
		Short: "Delete todo item(s)",
		Long: `Deletes one or more todo items from the list.

# delete task 1
togodo del 1

# delete multiple tasks
togodo del 1 3 5
`,
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"rm"},
		RunE: func(command *cobra.Command, args []string) error {
			// Determine if this is a term delete or task delete
			// If we have 2+ args and the second arg is not a number, treat as term delete
			if len(args) >= 2 {
				_, err := strconv.Atoi(args[1])
				if err != nil {
					// Second arg is not a number - this is a term delete
					return handleDelTerm(command, repo, args)
				}
			}

			// Handle as task delete (single or multiple)
			return handleDelTasks(command, repo, args)
		},
	}
}

func handleDelTasks(command *cobra.Command, repo todotxtlib.TodoRepository, args []string) error {
	// Parse all line numbers (1-based)
	indices := make([]int, 0, len(args))
	for _, arg := range args {
		lineNum, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("invalid line number: %s", arg)
		}

		// Find the array index for this line number
		index := repo.FindIndexByLineNumber(lineNum)
		if index == -1 {
			return fmt.Errorf("TODO: No task %d.", lineNum)
		}
		indices = append(indices, index)
	}

	// Call business logic
	result, err := cmd.Del(repo, indices)
	if err != nil {
		return err
	}

	// Format output - show each deleted task
	for _, deleted := range result.DeletedTodos {
		fmt.Fprintf(command.OutOrStdout(), "%d %s\n", deleted.LineNumber, deleted.Text)
	}
	if len(result.DeletedTodos) == 1 {
		fmt.Fprintf(command.OutOrStdout(), "TODO: %d deleted.\n", result.DeletedTodos[0].LineNumber)
	} else {
		fmt.Fprintf(command.OutOrStdout(), "TODO: %d tasks deleted.\n", len(result.DeletedTodos))
	}
	return nil
}

func handleDelTerm(command *cobra.Command, repo todotxtlib.TodoRepository, args []string) error {
	// Parse line number (first arg)
	lineNum, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid line number: %s", args[0])
	}

	// Find the array index
	index := repo.FindIndexByLineNumber(lineNum)
	if index == -1 {
		return fmt.Errorf("TODO: No task %d.", lineNum)
	}

	// Join remaining args as the term to delete
	term := args[1]
	if len(args) > 2 {
		// If multiple words, join them
		term = args[1] // For now, just use the second arg as the term
	}

	// Call business logic
	result, err := cmd.DelTerm(repo, index, term)
	if err != nil {
		return err
	}

	// Format output
	fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.ModifiedTodo.LineNumber, result.ModifiedTodo.Text)
	fmt.Fprintf(command.OutOrStdout(), "TODO: Removed '%s' from task.\n", result.RemovedTerm)
	return nil
}
