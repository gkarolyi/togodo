package cmd

import (
	"time"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// ReportResult contains the statistics for the report
type ReportResult struct {
	Date      time.Time
	Total     int // Total tasks in todo.txt (done + pending)
	Done      int // Completed tasks in done.txt
	Pending   int // Pending/open tasks in todo.txt
}

// Report generates task statistics
// It counts tasks from todo.txt (both pending and completed) and done.txt
func Report(repo todotxtlib.TodoRepository, doneReader todotxtlib.Reader) (ReportResult, error) {
	// Get all todos from todo.txt
	allTodos, err := repo.ListAll()
	if err != nil {
		return ReportResult{}, err
	}

	// Count pending tasks in todo.txt (non-completed)
	pendingCount := 0
	for _, todo := range allTodos {
		if !todo.Done {
			pendingCount++
		}
	}

	// Read done.txt to count archived tasks
	doneTodos, err := doneReader.Read()
	if err != nil {
		// If done.txt doesn't exist or is empty, that's okay
		doneTodos = []todotxtlib.Todo{}
	}

	// Total is the count of all tasks in todo.txt
	totalCount := len(allTodos)

	// Done count is the number of tasks in done.txt
	doneCount := len(doneTodos)

	return ReportResult{
		Date:    time.Now(),
		Total:   totalCount,
		Done:    doneCount,
		Pending: pendingCount,
	}, nil
}
