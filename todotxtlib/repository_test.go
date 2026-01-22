package todotxtlib

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileRepository(t *testing.T) {
	t.Run("with a buffer reader and writer", func(t *testing.T) {
		t.Run("creates new repository with empty buffer", func(t *testing.T) {
			var buf bytes.Buffer
			reader := NewBufferReader(&buf)
			writer := NewBufferWriter(&buf)

			repo, err := NewFileRepository(reader, writer)
			if err != nil {
				t.Fatalf("NewFileRepository() error = %v, want nil", err)
			}
			if repo == nil {
				t.Fatal("NewFileRepository() returned nil repository")
			}

			todos, err := repo.ListAll()
			if err != nil {
				t.Errorf("ListAll() error = %v, want nil", err)
			}
			if len(todos) != 0 {
				t.Errorf("ListAll() returned %d todos, want 0", len(todos))
			}
		})

		t.Run("creates new repository with existing content", func(t *testing.T) {
			var buf bytes.Buffer
			buf.WriteString("Test todo 1\nTest todo 2\n")
			reader := NewBufferReader(&buf)
			writer := NewBufferWriter(&buf)

			repo, err := NewFileRepository(reader, writer)
			if err != nil {
				t.Fatalf("NewFileRepository() error = %v, want nil", err)
			}

			todos, err := repo.ListAll()
			if err != nil {
				t.Fatalf("ListAll() error = %v, want nil", err)
			}
			if len(todos) != 2 {
				t.Errorf("ListAll() returned %d todos, want 2", len(todos))
			}
		})
	})

	t.Run("with a file reader and writer", func(t *testing.T) {
		t.Run("creates new repository with empty file", func(t *testing.T) {
			tempDir := t.TempDir()
			tempFile := filepath.Join(tempDir, "test.todo.txt")

			reader := NewFileReader(tempFile)
			writer := NewFileWriter(tempFile)

			repo, err := NewFileRepository(reader, writer)
			if err != nil {
				t.Fatalf("NewFileRepository() error = %v, want nil", err)
			}

			todos, err := repo.ListAll()
			if err != nil {
				t.Fatalf("ListAll() error = %v, want nil", err)
			}
			if len(todos) != 0 {
				t.Errorf("ListAll() returned %d todos, want 0", len(todos))
			}
		})

		t.Run("creates new repository with existing content", func(t *testing.T) {
			tempDir := t.TempDir()
			tempFile := filepath.Join(tempDir, "test.todo.txt")
			testData := createTestTodos()
			content := ""
			for _, todo := range testData {
				content += todo.Text + "\n"
			}
			if err := os.WriteFile(tempFile, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			reader := NewFileReader(tempFile)
			writer := NewFileWriter(tempFile)

			repo, err := NewFileRepository(reader, writer)
			if err != nil {
				t.Fatalf("NewFileRepository() error = %v, want nil", err)
			}

			todos, err := repo.ListAll()
			if err != nil {
				t.Fatalf("ListAll() error = %v, want nil", err)
			}
			if len(todos) != 3 {
				t.Errorf("ListAll() returned %d todos, want 3", len(todos))
			}
		})
	})
}

func TestRepository_Add(t *testing.T) {
	repo, _ := setupTestRepository(t)

	t.Run("adds new todo", func(t *testing.T) {
		todo, err := repo.Add("Test todo 3")
		if err != nil {
			t.Errorf("Add() error = %v, want nil", err)
		}
		if todo.Text != "Test todo 3" {
			t.Errorf("Add() returned todo with text %s, want Test todo 3", todo.Text)
		}

		todos, err := repo.ListAll()
		if err != nil {
			t.Errorf("ListAll() error = %v, want nil", err)
		}
		if len(todos) != 4 {
			t.Errorf("expected 4 todos, got %d", len(todos))
		}
	})
}

func TestRepository_Remove(t *testing.T) {
	repo, _ := setupTestRepository(t)

	t.Run("removes todo by index", func(t *testing.T) {
		todo, err := repo.Remove(0)
		if err != nil {
			t.Errorf("Remove() error = %v, want nil", err)
		}
		if todo.Text != "(A) test todo 1 +project2 @context1" {
			t.Errorf("Remove() returned todo with text %s, want (A) test todo 1 +project2 @context1", todo.Text)
		}

		todos, err := repo.ListAll()
		if err != nil {
			t.Errorf("ListAll() error = %v, want nil", err)
		}
		if len(todos) != 2 {
			t.Errorf("expected 2 todos, got %d", len(todos))
		}
	})

	t.Run("returns error for invalid index", func(t *testing.T) {
		_, err := repo.Remove(999)
		if err == nil {
			t.Error("Remove() expected error for invalid index, got nil")
		}
	})
}

func TestRepository_Update(t *testing.T) {
	repo, _ := setupTestRepository(t)

	t.Run("updates todo by index", func(t *testing.T) {
		todo, err := repo.Update(0, NewTodo("Test todo 1 updated"))
		if err != nil {
			t.Errorf("Update() error = %v, want nil", err)
		}
		if todo.Text != "Test todo 1 updated" {
			t.Errorf("Update() returned todo with text %s, want Test todo 1 updated", todo.Text)
		}
	})

	t.Run("returns error for invalid index", func(t *testing.T) {
		_, err := repo.Update(999, NewTodo("Test todo 1 updated"))
		if err == nil {
			t.Error("Update() expected error for invalid index, got nil")
		}
	})
}

func TestRepository_ToggleDone(t *testing.T) {
	repo, _ := setupTestRepository(t)

	t.Run("sets a not done todo to done", func(t *testing.T) {
		todo, err := repo.ToggleDone(0)
		if err != nil {
			t.Errorf("ToggleDone() error = %v, want nil", err)
		}
		if !todo.Done {
			t.Error("ToggleDone() returned todo with done status false, want true")
		}
	})

	t.Run("sets a done todo to not done", func(t *testing.T) {
		todo, err := repo.ToggleDone(2)
		if err != nil {
			t.Errorf("ToggleDone() error = %v, want nil", err)
		}
		if todo.Done {
			t.Error("ToggleDone() returned todo with done status true, want false")
		}
	})
}

func TestRepository_SetPriority(t *testing.T) {
	repo, _ := setupTestRepository(t)

	t.Run("sets priority for a todo", func(t *testing.T) {
		todo, err := repo.SetPriority(0, "A")
		if err != nil {
			t.Errorf("SetPriority() error = %v, want nil", err)
		}
		if todo.Priority != "A" {
			t.Errorf("SetPriority() returned todo with priority %s, want A", todo.Priority)
		}
	})

	t.Run("returns error for invalid index", func(t *testing.T) {
		_, err := repo.SetPriority(999, "A")
		if err == nil {
			t.Error("SetPriority() expected error for invalid index, got nil")
		}
	})
}

func TestRepository_Filter(t *testing.T) {
	repo, _ := setupTestRepository(t)

	t.Run("filters todos with a combined filter", func(t *testing.T) {
		filter := Filter{
			Done:     "true",
			Priority: "",
			Project:  "+project1",
			Context:  "@context1",
		}
		filtered, err := repo.Filter(filter)
		if err != nil {
			t.Errorf("Filter() error = %v, want nil", err)
		}
		if len(filtered) != 1 {
			t.Errorf("Filter() returned %d todos, want 1", len(filtered))
		}
		if !filtered[0].Equals(NewTodo("x (C) test todo 3 +project1 @context1")) {
			t.Errorf("Filter() returned todo %v, want x (C) test todo 3 +project1 @context1", filtered[0].Text)
		}
	})
}

func TestRepository_Search(t *testing.T) {
	repo, _ := setupTestRepository(t)

	t.Run("searches todos by text", func(t *testing.T) {
		filter := Filter{Text: "test todo 1"}
		filtered, err := repo.Filter(filter)
		if err != nil {
			t.Errorf("Filter() error = %v, want nil", err)
		}
		if len(filtered) != 1 {
			t.Errorf("Filter() returned %d todos, want 1", len(filtered))
		}
		if !filtered[0].Equals(NewTodo("(A) test todo 1 +project2 @context1")) {
			t.Errorf("Filter() returned todo %v, want (A) test todo 1 +project2 @context1", filtered[0].Text)
		}
	})

	t.Run("returns empty list for non-matching query", func(t *testing.T) {
		filter := Filter{Text: "non-matching"}
		filtered, err := repo.Filter(filter)
		if err != nil {
			t.Errorf("Filter() error = %v, want nil", err)
		}
		if len(filtered) != 0 {
			t.Errorf("Filter() returned %d todos, want 0", len(filtered))
		}
	})
}

func TestRepository_Sort(t *testing.T) {
	defaultSort := NewDefaultSort()
	descendingSort := Sort{Field: Text, Order: Descending}

	tests := []struct {
		name     string
		sort     *Sort
		expected []Todo
	}{
		{
			name: "sorts todos by text ascending with done items last",
			sort: &defaultSort,
			expected: []Todo{
				NewTodo("(A) test todo 1 +project2 @context1"),
				NewTodo("(B) test todo 2 +project1 @context2"),
				NewTodo("x (C) test todo 3 +project1 @context1"),
			},
		},
		{
			name: "sorts todos by text descending with done items first",
			sort: &descendingSort,
			expected: []Todo{
				NewTodo("x (C) test todo 3 +project1 @context1"),
				NewTodo("(B) test todo 2 +project1 @context2"),
				NewTodo("(A) test todo 1 +project2 @context1"),
			},
		},
		{
			name: "nil sort uses default sort (ascending with done items last)",
			sort: nil,
			expected: []Todo{
				NewTodo("(A) test todo 1 +project2 @context1"),
				NewTodo("(B) test todo 2 +project1 @context2"),
				NewTodo("x (C) test todo 3 +project1 @context1"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, _ := setupTestRepository(t)

			// Sort the todos
			repo.Sort(tt.sort)

			// Verify the sorted todos
			todos, err := repo.ListAll()
			if err != nil {
				t.Errorf("ListAll() error = %v, want nil", err)
				return
			}

			if len(todos) != len(tt.expected) {
				t.Errorf("Sort() returned %d todos, want %d", len(todos), len(tt.expected))
				return
			}

			for i := range todos {
				if !todos[i].Equals(tt.expected[i]) {
					t.Errorf("Sort() todo[%d] = %v, want %v", i, todos[i], tt.expected[i])
				}
			}
		})
	}
}

func TestRepository_ListProjects(t *testing.T) {
	repo, _ := setupTestRepository(t)

	t.Run("returns unique sorted projects", func(t *testing.T) {
		projects, err := repo.ListProjects()
		if err != nil {
			t.Errorf("ListProjects() error = %v, want nil", err)
		}
		if len(projects) != 2 {
			t.Errorf("ListProjects() returned %d projects, want 2", len(projects))
		}
		if projects[0] != "+project1" || projects[1] != "+project2" {
			t.Errorf("ListProjects() returned projects in wrong order, got %v, want [+project1 +project2]", projects)
		}
	})
}

func TestRepository_ListContexts(t *testing.T) {
	repo, _ := setupTestRepository(t)

	t.Run("returns unique sorted contexts", func(t *testing.T) {
		contexts, err := repo.ListContexts()
		if err != nil {
			t.Errorf("ListContexts() error = %v, want nil", err)
		}
		if len(contexts) != 2 {
			t.Errorf("ListContexts() returned %d contexts, want 2", len(contexts))
		}
		if contexts[0] != "@context1" || contexts[1] != "@context2" {
			t.Errorf("ListContexts() returned contexts in wrong order, got %v, want [@context1 @context2]", contexts)
		}
	})
}

func TestRepository_Get(t *testing.T) {
	repo, _ := setupTestRepository(t)

	t.Run("gets todo by index", func(t *testing.T) {
		todo, err := repo.Get(0)
		if err != nil {
			t.Errorf("Get() error = %v, want nil", err)
		}
		if todo.Text != "(A) test todo 1 +project2 @context1" {
			t.Errorf("Get() returned todo with text %s, want (A) test todo 1 +project2 @context1", todo.Text)
		}
	})

	t.Run("returns error for invalid index", func(t *testing.T) {
		_, err := repo.Get(999)
		if err == nil {
			t.Error("Get() expected error for invalid index, got nil")
		}
	})

	t.Run("returns error for negative index", func(t *testing.T) {
		_, err := repo.Get(-1)
		if err == nil {
			t.Error("Get() expected error for negative index, got nil")
		}
	})
}

func TestRepository_Save(t *testing.T) {
	t.Run("saves todos to buffer", func(t *testing.T) {
		var buf bytes.Buffer
		reader := NewBufferReader(&buf)
		writer := NewBufferWriter(&buf)

		repo, err := NewFileRepository(reader, writer)
		if err != nil {
			t.Fatalf("NewFileRepository() error = %v, want nil", err)
		}

		// Add some todos
		repo.Add("Test todo 1")
		repo.Add("Test todo 2")

		if err := repo.Save(); err != nil {
			t.Errorf("Save() error = %v, want nil", err)
		}

		expected := "Test todo 1\nTest todo 2\n"
		if buf.String() != expected {
			t.Errorf("Save() wrote %q, want %q", buf.String(), expected)
		}
	})
}
