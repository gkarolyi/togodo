package todotxtlib

import (
	"sort"
)

// SortOrder represents the order in which to sort
type SortOrder int

const (
	Ascending SortOrder = iota
	Descending
)

// SortField represents the field to sort by
type SortField int

const (
	Text SortField = iota
	// TODO: Add other fields
)

// Sort represents sorting criteria for todos
type Sort struct {
	Field SortField
	Order SortOrder
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
