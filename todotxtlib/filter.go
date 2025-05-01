package todotxtlib

import "strings"

// Filter holds criteria for filtering todos
type Filter struct {
	Done     string // "true", "false", or "" (empty string means don't filter by done status)
	Priority string
	Project  string
	Context  string
	Text     string
}

// Apply applies the filter criteria to a list of todos and returns the matching ones
func (f Filter) Apply(todos []Todo) []Todo {
	var filtered []Todo

	for _, todo := range todos {
		if f.matches(todo) {
			filtered = append(filtered, todo)
		}
	}

	return filtered
}

// matches checks if a todo matches all the filter criteria
func (f Filter) matches(todo Todo) bool {
	// Check done status
	if f.Done != "" {
		wantDone := f.Done == "true"
		if todo.Done != wantDone {
			return false
		}
	}

	// Check priority
	if f.Priority != "" && todo.Priority != f.Priority {
		return false
	}

	// Check project
	if f.Project != "" {
		found := false
		for _, project := range todo.Projects {
			if project == f.Project {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check context
	if f.Context != "" {
		found := false
		for _, context := range todo.Contexts {
			if context == f.Context {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check text content
	if f.Text != "" && !strings.Contains(todo.Text, f.Text) {
		return false
	}

	return true
}
