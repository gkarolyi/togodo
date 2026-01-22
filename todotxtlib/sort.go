package todotxtlib

import (
	"sort"
	"strings"
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
			// First: Sort by done status (not done first)
			if todos[i].Done != todos[j].Done {
				if s.Order == Descending {
					return todos[i].Done
				} else {
					return !todos[i].Done
				}
			}

			// Second: Sort by priority (A-Z, then no priority)
			iPri := todos[i].Priority
			jPri := todos[j].Priority

			// If one has priority and the other doesn't, prioritized comes first
			if (iPri != "") != (jPri != "") {
				if s.Order == Descending {
					return iPri == ""
				} else {
					return iPri != ""
				}
			}

			// If both have priority, sort by priority (A before B), then by text
			if iPri != "" && jPri != "" && iPri != jPri {
				if s.Order == Descending {
					return iPri > jPri
				} else {
					return iPri < jPri
				}
			}

			// Third: Sort by text (case-insensitive alphabetically, letters before symbols)
			iText := strings.ToLower(todos[i].Text)
			jText := strings.ToLower(todos[j].Text)

			// Custom comparison: letters (a-z) should come before special characters
			iIsLetter := len(iText) > 0 && ((iText[0] >= 'a' && iText[0] <= 'z') || (iText[0] >= 'A' && iText[0] <= 'Z'))
			jIsLetter := len(jText) > 0 && ((jText[0] >= 'a' && jText[0] <= 'z') || (jText[0] >= 'A' && jText[0] <= 'Z'))

			// If one starts with a letter and the other doesn't, letter comes first
			if iIsLetter != jIsLetter {
				if s.Order == Descending {
					return !iIsLetter
				} else {
					return iIsLetter
				}
			}

			// Otherwise, use standard string comparison
			if s.Order == Descending {
				return iText > jText
			} else {
				return iText < jText
			}
		})
	}
}
