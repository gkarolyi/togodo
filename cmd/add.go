package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// AddResult contains the result of an Add operation
type AddResult struct {
	Todo       todotxtlib.Todo
	LineNumber int
}

// AddMultilineResult contains the result of adding multiple tasks
type AddMultilineResult struct {
	Todos       []todotxtlib.Todo
	LineNumbers []int
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

// AddMultiple adds multiple tasks from text with embedded newlines
func AddMultiple(repo todotxtlib.TodoRepository, text string, autoDate bool) (AddMultilineResult, error) {
	// Split on newlines to create separate tasks
	lines := strings.Split(text, "\n")

	var todos []todotxtlib.Todo
	var lineNumbers []int

	for _, line := range lines {
		// Skip empty lines
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Add creation date if enabled and not already present
		taskText := line
		if autoDate {
			taskText = addCreationDate(taskText)
		}

		// Add to repository
		todo, err := repo.Add(taskText)
		if err != nil {
			return AddMultilineResult{}, fmt.Errorf("failed to add todo: %w", err)
		}

		todos = append(todos, todo)
		lineNumbers = append(lineNumbers, todo.LineNumber)
	}

	// Sort and save after adding all tasks
	repo.Sort(nil)
	if err := repo.Save(); err != nil {
		return AddMultilineResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return AddMultilineResult{
		Todos:       todos,
		LineNumbers: lineNumbers,
	}, nil
}

var priorityRe = regexp.MustCompile(`^\(([A-Z])\) `)
var dateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)

// getTodayDate returns today's date, or mocked date from TODO_TEST_TIME env var
func getTodayDate() string {
	// Check for test date mock (for test compatibility with todo.txt-cli)
	if testTime := os.Getenv("TODO_TEST_TIME"); testTime != "" {
		// TODO_TEST_TIME is a Unix timestamp
		if timestamp, err := strconv.ParseInt(testTime, 10, 64); err == nil {
			return time.Unix(timestamp, 0).Format("2006-01-02")
		}
	}
	// Use actual current time
	return time.Now().Format("2006-01-02")
}

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

	// Add today's date (or mocked date for testing)
	today := getTodayDate()

	if priorityPrefix != "" {
		// Insert date after priority: (A) YYYY-MM-DD text
		return priorityPrefix + today + " " + remaining
	}
	// Insert date at beginning: YYYY-MM-DD text
	return today + " " + text
}
