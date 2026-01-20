package business

import (
	"fmt"
	"strings"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// AddResult contains the result of an Add operation
type AddResult struct {
	Todo       todotxtlib.Todo
	LineNumber int
}

// Add adds a new todo task to the repository
func Add(repo todotxtlib.TodoRepository, args []string) (AddResult, error) {
	if len(args) == 0 {
		return AddResult{}, fmt.Errorf("task text required")
	}

	// Join args into single task
	text := strings.Join(args, " ")

	// Add to repository
	todo, err := repo.Add(text)
	if err != nil {
		return AddResult{}, fmt.Errorf("failed to add todo: %w", err)
	}

	// Sort and save
	repo.Sort(nil)
	if err := repo.Save(); err != nil {
		return AddResult{}, fmt.Errorf("failed to save: %w", err)
	}

	// Find line number after sort
	allTodos, err := repo.ListAll()
	if err != nil {
		return AddResult{}, fmt.Errorf("failed to list todos: %w", err)
	}

	lineNumber := 1
	for i, t := range allTodos {
		if t.Text == todo.Text && t.Priority == todo.Priority {
			lineNumber = i + 1
			break
		}
	}

	return AddResult{Todo: todo, LineNumber: lineNumber}, nil
}
