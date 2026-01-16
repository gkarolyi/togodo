package todotxtlib

import (
	"slices"
	"strings"
	"testing"
)

func TestFilter_Apply(t *testing.T) {
	repo, _ := setupTestRepository(t)
	todos, err := repo.ListAll()
	if err != nil {
		t.Fatalf("ListAll() error = %v", err)
	}

	t.Run("filter by done status", func(t *testing.T) {
		filter := Filter{Done: "true"}
		filtered := filter.Apply(todos)

		if len(filtered) != 1 {
			t.Errorf("Filter.Apply() returned %d todos, want 1", len(filtered))
		}
		for _, todo := range filtered {
			if !todo.Done {
				t.Error("Filter.Apply() returned non-done todo")
			}
		}
	})

	t.Run("filter by priority", func(t *testing.T) {
		filter := Filter{Priority: "A"}
		filtered := filter.Apply(todos)

		if len(filtered) != 1 {
			t.Errorf("Filter.Apply() returned %d todos, want 1", len(filtered))
		}
		if filtered[0].Priority != "A" {
			t.Errorf("Filter.Apply() returned todo with priority %s, want A", filtered[0].Priority)
		}
	})

	t.Run("filter by project", func(t *testing.T) {
		filter := Filter{Project: "+project1"}
		filtered := filter.Apply(todos)

		if len(filtered) != 2 {
			t.Errorf("Filter.Apply() returned %d todos, want 2", len(filtered))
		}
		for _, todo := range filtered {
			if !slices.Contains(todo.Projects, "+project1") {
				t.Error("Filter.Apply() returned todo without +project1")
			}
		}
	})

	t.Run("filter by context", func(t *testing.T) {
		filter := Filter{Context: "@context1"}
		filtered := filter.Apply(todos)

		if len(filtered) != 2 {
			t.Errorf("Filter.Apply() returned %d todos, want 2", len(filtered))
		}
		for _, todo := range filtered {
			if !slices.Contains(todo.Contexts, "@context1") {
				t.Error("Filter.Apply() returned todo without @context1")
			}
		}
	})

	t.Run("filter by text content", func(t *testing.T) {
		filter := Filter{Text: "todo 1"}
		filtered := filter.Apply(todos)

		if len(filtered) != 1 {
			t.Errorf("Filter.Apply() returned %d todos, want 1", len(filtered))
		}
		if !strings.Contains(filtered[0].Text, "todo 1") {
			t.Error("Filter.Apply() returned todo without 'todo 1' in text")
		}
	})

	t.Run("filter by multiple criteria", func(t *testing.T) {
		filter := Filter{
			Done:     "false",
			Priority: "A",
			Project:  "+project2",
		}
		filtered := filter.Apply(todos)

		if len(filtered) != 1 {
			t.Errorf("Filter.Apply() returned %d todos, want 1", len(filtered))
		}
		todo := filtered[0]
		if todo.Done {
			t.Error("Filter.Apply() returned done todo")
		}
		if todo.Priority != "A" {
			t.Errorf("Filter.Apply() returned todo with priority %s, want A", todo.Priority)
		}
		if !slices.Contains(todo.Projects, "+project2") {
			t.Error("Filter.Apply() returned todo without +project2")
		}
	})

	t.Run("empty filter returns all todos", func(t *testing.T) {
		filter := Filter{}
		filtered := filter.Apply(todos)

		if len(filtered) != len(todos) {
			t.Errorf("Filter.Apply() returned %d todos, want %d", len(filtered), len(todos))
		}
	})
}

func TestFilter_matches(t *testing.T) {
	todo := NewTodo("(A) test todo +project1 @context1")

	t.Run("matches done status", func(t *testing.T) {
		filter := Filter{Done: "false"}
		if !filter.matches(todo) {
			t.Error("Filter.matches() should match todo with matching done status")
		}
	})

	t.Run("matches priority", func(t *testing.T) {
		filter := Filter{Priority: "A"}
		if !filter.matches(todo) {
			t.Error("Filter.matches() should match todo with matching priority")
		}
	})

	t.Run("matches project", func(t *testing.T) {
		filter := Filter{Project: "+project1"}
		if !filter.matches(todo) {
			t.Error("Filter.matches() should match todo with matching project")
		}
	})

	t.Run("matches context", func(t *testing.T) {
		filter := Filter{Context: "@context1"}
		if !filter.matches(todo) {
			t.Error("Filter.matches() should match todo with matching context")
		}
	})

	t.Run("matches text content", func(t *testing.T) {
		filter := Filter{Text: "test"}
		if !filter.matches(todo) {
			t.Error("Filter.matches() should match todo with matching text content")
		}
	})

	t.Run("does not match non-matching done status", func(t *testing.T) {
		filter := Filter{Done: "true"}
		if filter.matches(todo) {
			t.Error("Filter.matches() should not match todo with non-matching done status")
		}
	})

	t.Run("does not match non-matching priority", func(t *testing.T) {
		filter := Filter{Priority: "B"}
		if filter.matches(todo) {
			t.Error("Filter.matches() should not match todo with non-matching priority")
		}
	})

	t.Run("does not match non-matching project", func(t *testing.T) {
		filter := Filter{Project: "+project2"}
		if filter.matches(todo) {
			t.Error("Filter.matches() should not match todo with non-matching project")
		}
	})

	t.Run("does not match non-matching context", func(t *testing.T) {
		filter := Filter{Context: "@context2"}
		if filter.matches(todo) {
			t.Error("Filter.matches() should not match todo with non-matching context")
		}
	})

	t.Run("does not match non-matching text content", func(t *testing.T) {
		filter := Filter{Text: "nonexistent"}
		if filter.matches(todo) {
			t.Error("Filter.matches() should not match todo with non-matching text content")
		}
	})
}
