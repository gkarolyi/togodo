package cmd

import (
	"fmt"
	"strings"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// DelResult contains the result of a Del operation
type DelResult struct {
	DeletedTodos []todotxtlib.Todo
}

// DelTermResult contains the result of removing a term from a task
type DelTermResult struct {
	ModifiedTodo todotxtlib.Todo
	RemovedTerm  string
}

// Del deletes todos from the repository by indices (0-based)
func Del(repo todotxtlib.TodoRepository, indices []int) (DelResult, error) {
	if len(indices) == 0 {
		return DelResult{}, fmt.Errorf("no indices provided")
	}

	// STEP 1: Validate all indices first (fail-fast)
	allTodos, err := repo.ListAll()
	if err != nil {
		return DelResult{}, fmt.Errorf("failed to list todos: %w", err)
	}

	for _, index := range indices {
		if index < 0 || index >= len(allTodos) {
			return DelResult{}, fmt.Errorf("invalid index: %d", index)
		}
	}

	// STEP 2: Sort indices in descending order to delete from end to beginning
	// This prevents index shifting issues
	sortedIndices := make([]int, len(indices))
	copy(sortedIndices, indices)
	for i := 0; i < len(sortedIndices); i++ {
		for j := i + 1; j < len(sortedIndices); j++ {
			if sortedIndices[i] < sortedIndices[j] {
				sortedIndices[i], sortedIndices[j] = sortedIndices[j], sortedIndices[i]
			}
		}
	}

	// STEP 3: Delete in descending order
	deletedTodos := make([]todotxtlib.Todo, 0, len(sortedIndices))
	for _, index := range sortedIndices {
		deleted, err := repo.Remove(index)
		if err != nil {
			return DelResult{}, fmt.Errorf("failed to remove todo at index %d: %w", index, err)
		}
		deletedTodos = append(deletedTodos, deleted)
	}

	// STEP 4: Save
	if err := repo.Save(); err != nil {
		return DelResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return DelResult{
		DeletedTodos: deletedTodos,
	}, nil
}

// DelTerm removes a term from a task at the given index
func DelTerm(repo todotxtlib.TodoRepository, index int, term string) (DelTermResult, error) {
	// Get the todo
	todo, err := repo.Get(index)
	if err != nil {
		return DelTermResult{}, fmt.Errorf("failed to get todo: %w", err)
	}

	// Remove the term from the text (case-sensitive, whole word match)
	oldText := todo.Text
	newText := removeTerm(todo.Text, term)

	if newText == oldText {
		// Term not found
		return DelTermResult{}, fmt.Errorf("term '%s' not found in task", term)
	}

	// Update the todo
	todo.Text = newText
	_, err = repo.Update(index, todo)
	if err != nil {
		return DelTermResult{}, fmt.Errorf("failed to update todo: %w", err)
	}

	// Save
	if err := repo.Save(); err != nil {
		return DelTermResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return DelTermResult{
		ModifiedTodo: todo,
		RemovedTerm:  term,
	}, nil
}

// removeTerm removes a term from text, handling word boundaries and extra spaces
func removeTerm(text, term string) string {
	// Simple approach: replace all occurrences of the term
	// Handle both " term " and " term" at end and "term " at start
	result := text

	// Try to remove with surrounding spaces first
	result = strings.ReplaceAll(result, " "+term+" ", " ")

	// Try to remove at the start
	if strings.HasPrefix(result, term+" ") {
		result = strings.TrimPrefix(result, term+" ")
	}

	// Try to remove at the end
	if strings.HasSuffix(result, " "+term) {
		result = strings.TrimSuffix(result, " "+term)
	}

	// If the whole text is just the term
	if result == term {
		result = ""
	}

	// Clean up multiple spaces
	result = strings.Join(strings.Fields(result), " ")

	return result
}
