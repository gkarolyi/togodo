package todotxtlib

import (
	"sort"
)

// sortOrder represents the order in which to sort
type sortOrder int

const (
	Ascending sortOrder = iota
	Descending
)

// sortField represents the field to sort by
type sortField int

const (
	Text sortField = iota
	// TODO: Add other fields
)

// Sort represents sorting criteria for todos
type Sort struct {
	Field sortField // Field to sort by, e.g. Text
	Order sortOrder // Order to sort by, e.g. Ascending
}

// NewDefaultSort returns the default todo.txt sorting (alphabetical with done items last)
func NewDefaultSort() Sort {
	return Sort{
		Field: Text,
		Order: Ascending,
	}
}

// Apply sorts the todos according to the specified criteria
func (s Sort) Apply(todos []Todo) {
	if s.Field == Text {
		sort.SliceStable(todos, func(i, j int) bool {
			if todos[i].Done != todos[j].Done {
				if s.Order == Descending {
					return todos[i].Done
				} else {
					return !todos[i].Done
				}
			}

			if s.Order == Descending {
				return todos[i].Text > todos[j].Text
			} else {
				return todos[i].Text < todos[j].Text
			}
		})
	}
}
