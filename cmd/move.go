package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// MoveResult contains the result of moving a todo
type MoveResult struct {
	Todo       todotxtlib.Todo
	LineNumber int
	SourceFile string
	DestFile   string
}

// Move moves a todo from one file to another
func Move(sourceRepo, destRepo todotxtlib.TodoRepository, lineNumber int) (MoveResult, error) {
	// Find index by line number
	index := sourceRepo.FindIndexByLineNumber(lineNumber)
	if index == -1 {
		return MoveResult{}, fmt.Errorf("TODO: No task %d.", lineNumber)
	}

	// Get the todo
	todo, err := sourceRepo.Get(index)
	if err != nil {
		return MoveResult{}, fmt.Errorf("failed to get task: %w", err)
	}

	// Add to destination
	if _, err := destRepo.Add(todo.Text); err != nil {
		return MoveResult{}, fmt.Errorf("failed to add to destination: %w", err)
	}

	// Remove from source
	if _, err := sourceRepo.Remove(index); err != nil {
		return MoveResult{}, fmt.Errorf("failed to remove from source: %w", err)
	}

	// Save both
	if err := sourceRepo.Save(); err != nil {
		return MoveResult{}, fmt.Errorf("failed to save source: %w", err)
	}

	if err := destRepo.Save(); err != nil {
		return MoveResult{}, fmt.Errorf("failed to save destination: %w", err)
	}

	return MoveResult{
		Todo:       todo,
		LineNumber: lineNumber,
	}, nil
}
