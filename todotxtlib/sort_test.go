package todotxtlib

import (
	"testing"
)

func TestSort_Apply(t *testing.T) {
	tests := []struct {
		name     string
		sort     Sort
		input    []Todo
		expected []Todo
	}{
		{
			name: "basic text sorting with done items last (ascending)",
			sort: NewDefaultSort(),
			input: []Todo{
				{Text: "x done todo", Done: true},
				{Text: "todo 2", Done: false},
				{Text: "todo 4", Done: false},
				{Text: "todo 1", Done: false},
				{Text: "x another done todo", Done: true},
				{Text: "todo 3", Done: false},
			},
			expected: []Todo{
				{Text: "todo 1", Done: false},
				{Text: "todo 2", Done: false},
				{Text: "todo 3", Done: false},
				{Text: "todo 4", Done: false},
				{Text: "x another done todo", Done: true},
				{Text: "x done todo", Done: true},
			},
		},
		{
			name: "basic text sorting with done items first (descending)",
			sort: Sort{Order: Descending},
			input: []Todo{
				{Text: "x done todo", Done: true},
				{Text: "todo 2", Done: false},
				{Text: "todo 4", Done: false},
				{Text: "todo 1", Done: false},
				{Text: "x another done todo", Done: true},
				{Text: "todo 3", Done: false},
			},
			expected: []Todo{
				{Text: "x done todo", Done: true},
				{Text: "x another done todo", Done: true},
				{Text: "todo 4", Done: false},
				{Text: "todo 3", Done: false},
				{Text: "todo 2", Done: false},
				{Text: "todo 1", Done: false},
			},
		},
		{
			name: "text sorting with same text but different done status (ascending)",
			sort: NewDefaultSort(),
			input: []Todo{
				{Text: "same todo", Done: true},
				{Text: "same todo", Done: false},
				{Text: "same todo", Done: true},
				{Text: "same todo", Done: false},
			},
			expected: []Todo{
				{Text: "same todo", Done: false},
				{Text: "same todo", Done: false},
				{Text: "same todo", Done: true},
				{Text: "same todo", Done: true},
			},
		},
		{
			name: "text sorting with same text but different done status (descending)",
			sort: Sort{Order: Descending},
			input: []Todo{
				{Text: "same todo", Done: true},
				{Text: "same todo", Done: false},
				{Text: "same todo", Done: true},
				{Text: "same todo", Done: false},
			},
			expected: []Todo{
				{Text: "same todo", Done: true},
				{Text: "same todo", Done: true},
				{Text: "same todo", Done: false},
				{Text: "same todo", Done: false},
			},
		},
		{
			name: "text sorting with special characters (ascending)",
			sort: NewDefaultSort(),
			input: []Todo{
				{Text: "!important todo", Done: false},
				{Text: "regular todo", Done: false},
				{Text: "#tagged todo", Done: false},
				{Text: "x done !important", Done: true},
			},
			expected: []Todo{
				{Text: "!important todo", Done: false},
				{Text: "#tagged todo", Done: false},
				{Text: "regular todo", Done: false},
				{Text: "x done !important", Done: true},
			},
		},
		{
			name: "text sorting with special characters (descending)",
			sort: Sort{Order: Descending},
			input: []Todo{
				{Text: "!important todo", Done: false},
				{Text: "regular todo", Done: false},
				{Text: "#tagged todo", Done: false},
				{Text: "x done !important", Done: true},
			},
			expected: []Todo{
				{Text: "x done !important", Done: true},
				{Text: "regular todo", Done: false},
				{Text: "#tagged todo", Done: false},
				{Text: "!important todo", Done: false},
			},
		},
		{
			name: "text sorting with mixed case (ascending)",
			sort: NewDefaultSort(),
			input: []Todo{
				{Text: "Todo with Capital", Done: false},
				{Text: "todo with lowercase", Done: false},
				{Text: "x Done with Capital", Done: true},
				{Text: "x done with lowercase", Done: true},
			},
			expected: []Todo{
				{Text: "Todo with Capital", Done: false},
				{Text: "todo with lowercase", Done: false},
				{Text: "x Done with Capital", Done: true},
				{Text: "x done with lowercase", Done: true},
			},
		},
		{
			name: "text sorting with mixed case (descending)",
			sort: Sort{Order: Descending},
			input: []Todo{
				{Text: "Todo with Capital", Done: false},
				{Text: "todo with lowercase", Done: false},
				{Text: "x Done with Capital", Done: true},
				{Text: "x done with lowercase", Done: true},
			},
			expected: []Todo{
				{Text: "x done with lowercase", Done: true},
				{Text: "x Done with Capital", Done: true},
				{Text: "todo with lowercase", Done: false},
				{Text: "Todo with Capital", Done: false},
			},
		},
		{
			name:     "empty input",
			sort:     NewDefaultSort(),
			input:    []Todo{},
			expected: []Todo{},
		},
		{
			name: "single todo",
			sort: NewDefaultSort(),
			input: []Todo{
				{Text: "single todo", Done: false},
			},
			expected: []Todo{
				{Text: "single todo", Done: false},
			},
		},
		{
			name: "all done todos (ascending)",
			sort: NewDefaultSort(),
			input: []Todo{
				{Text: "x done 3", Done: true},
				{Text: "x done 1", Done: true},
				{Text: "x done 2", Done: true},
			},
			expected: []Todo{
				{Text: "x done 1", Done: true},
				{Text: "x done 2", Done: true},
				{Text: "x done 3", Done: true},
			},
		},
		{
			name: "all done todos (descending)",
			sort: Sort{Order: Descending},
			input: []Todo{
				{Text: "x done 3", Done: true},
				{Text: "x done 1", Done: true},
				{Text: "x done 2", Done: true},
			},
			expected: []Todo{
				{Text: "x done 3", Done: true},
				{Text: "x done 2", Done: true},
				{Text: "x done 1", Done: true},
			},
		},
		{
			name: "all not done todos (ascending)",
			sort: NewDefaultSort(),
			input: []Todo{
				{Text: "todo 3", Done: false},
				{Text: "todo 1", Done: false},
				{Text: "todo 2", Done: false},
			},
			expected: []Todo{
				{Text: "todo 1", Done: false},
				{Text: "todo 2", Done: false},
				{Text: "todo 3", Done: false},
			},
		},
		{
			name: "all not done todos (descending)",
			sort: Sort{Order: Descending},
			input: []Todo{
				{Text: "todo 3", Done: false},
				{Text: "todo 1", Done: false},
				{Text: "todo 2", Done: false},
			},
			expected: []Todo{
				{Text: "todo 3", Done: false},
				{Text: "todo 2", Done: false},
				{Text: "todo 1", Done: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of the input slice to avoid modifying the original
			todos := make([]Todo, len(tt.input))
			copy(todos, tt.input)

			tt.sort.Apply(todos)

			if len(todos) != len(tt.expected) {
				t.Errorf("Sort.Apply() returned %d todos, want %d", len(todos), len(tt.expected))
				return
			}

			for i := range todos {
				if todos[i].Text != tt.expected[i].Text {
					t.Errorf("Sort.Apply() todo[%d] = %q, want %q", i, todos[i].Text, tt.expected[i].Text)
				}
				if todos[i].Done != tt.expected[i].Done {
					t.Errorf("Sort.Apply() todo[%d] done = %v, want %v", i, todos[i].Done, tt.expected[i].Done)
				}
			}
		})
	}
}
