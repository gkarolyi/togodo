package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// AddResult contains the result of an Add operation
type AddResult struct {
	Todo       todotxtlib.Todo
	LineNumber int
}

// Add adds a new todo task to the repository
func Add(repo todotxtlib.TodoRepository, args []string, autoDate bool) (AddResult, error) {
	if len(args) == 0 {
		return AddResult{}, fmt.Errorf("task text required")
	}

	// Join args into single task
	text := strings.Join(args, " ")

	// Add creation date if enabled and not already present
	if autoDate {
		text = addCreationDate(text)
	}

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

	// Return the todo with its assigned line number
	return AddResult{Todo: todo, LineNumber: todo.LineNumber}, nil
}

var priorityRe = regexp.MustCompile(`^\(([A-Z])\) `)
var dateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)

// addCreationDate adds a creation date to the task text if not already present
func addCreationDate(text string) string {
	// Check if text starts with priority
	remaining := text
	var priorityPrefix string

	if priorityRe.MatchString(text) {
		match := priorityRe.FindString(text)
		priorityPrefix = match
		remaining = text[len(match):]
	}

	// Check if date already present after priority
	if dateRe.MatchString(remaining) {
		// Date already present, return as-is
		return text
	}

	// Add today's date
	today := time.Now().Format("2006-01-02")

	if priorityPrefix != "" {
		// Insert date after priority: (A) YYYY-MM-DD text
		return priorityPrefix + today + " " + remaining
	}
	// Insert date at beginning: YYYY-MM-DD text
	return today + " " + text
}
