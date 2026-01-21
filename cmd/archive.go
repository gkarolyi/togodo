package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// ArchiveResult contains the result of an Archive operation
type ArchiveResult struct {
	ArchivedTodos []todotxtlib.Todo
}

// Archive moves completed todos from todo.txt to done.txt
// Returns the list of archived todos
func Archive(repo todotxtlib.TodoRepository, doneReader todotxtlib.Reader, doneWriter todotxtlib.Writer) (ArchiveResult, error) {
	// Get all todos
	allTodos, err := repo.ListAll()
	if err != nil {
		return ArchiveResult{}, fmt.Errorf("failed to list todos: %w", err)
	}

	// Separate completed and non-completed todos
	var completedTodos []todotxtlib.Todo
	var activeTodos []todotxtlib.Todo

	for _, todo := range allTodos {
		if todo.Done {
			completedTodos = append(completedTodos, todo)
		} else {
			activeTodos = append(activeTodos, todo)
		}
	}

	// If no completed todos, return early
	if len(completedTodos) == 0 {
		return ArchiveResult{ArchivedTodos: []todotxtlib.Todo{}}, nil
	}

	// Read existing done.txt content
	existingDoneTodos, err := doneReader.Read()
	if err != nil {
		return ArchiveResult{}, fmt.Errorf("failed to read done.txt: %w", err)
	}

	// Combine existing and new done todos
	allDoneTodos := append(existingDoneTodos, completedTodos...)

	// Write all done todos to done.txt
	if err := doneWriter.Write(allDoneTodos); err != nil {
		return ArchiveResult{}, fmt.Errorf("failed to write to done.txt: %w", err)
	}

	// Update repository with only active todos
	// We need to clear and re-add all active todos
	// First, remove all todos in reverse order to avoid index issues
	for i := len(allTodos) - 1; i >= 0; i-- {
		_, err := repo.Remove(i)
		if err != nil {
			return ArchiveResult{}, fmt.Errorf("failed to remove todo at index %d: %w", i, err)
		}
	}

	// Add back only active todos
	for _, todo := range activeTodos {
		_, err := repo.Add(todo.Text)
		if err != nil {
			return ArchiveResult{}, fmt.Errorf("failed to add active todo: %w", err)
		}
	}

	// Save the repository
	if err := repo.Save(); err != nil {
		return ArchiveResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return ArchiveResult{
		ArchivedTodos: completedTodos,
	}, nil
}
