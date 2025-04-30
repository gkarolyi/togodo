package todotxtlib

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func createTestTodos() []Todo {
	return []Todo{
		NewTodo("(A) test todo 1 +project2 @context1"),
		NewTodo("(B) test todo 2 +project1 @context2"),
		NewTodo("x (C) test todo 3"),
	}
}

func setupTestRepository(tb testing.TB) *Repository {
	// Create a buffer with test todos
	var buf bytes.Buffer
	for _, todo := range createTestTodos() {
		buf.WriteString(todo.Text + "\n")
	}

	// Create reader and writer that work with the buffer
	reader := NewBufferReader(&buf)
	writer := NewBufferWriter(&buf)

	// Create repository with the buffer-based reader and writer
	repo, err := NewRepository(reader, writer)
	if err != nil {
		tb.Fatalf("Failed to create test repository: %v", err)
	}

	return repo
}

func TestNewRepository(t *testing.T) {
	t.Run("with a buffer reader and writer", func(t *testing.T) {
		t.Run("creates new repository with empty buffer", func(t *testing.T) {
			var buf bytes.Buffer
			reader := NewBufferReader(&buf)
			writer := NewBufferWriter(&buf)

			repo, err := NewRepository(reader, writer)
			if err != nil {
				t.Fatalf("NewRepository() error = %v, want nil", err)
			}
			if repo == nil {
				t.Fatal("NewRepository() returned nil repository")
			}

			todos, err := repo.ListTodos()
			if err != nil {
				t.Errorf("ListTodos() error = %v, want nil", err)
			}
			if len(todos) != 0 {
				t.Errorf("ListTodos() returned %d todos, want 0", len(todos))
			}
		})

		t.Run("creates new repository with existing content", func(t *testing.T) {
			var buf bytes.Buffer
			buf.WriteString("Test todo 1\nTest todo 2\n")
			reader := NewBufferReader(&buf)
			writer := NewBufferWriter(&buf)

			repo, err := NewRepository(reader, writer)
			if err != nil {
				t.Fatalf("NewRepository() error = %v, want nil", err)
			}

			todos, err := repo.ListTodos()
			if err != nil {
				t.Fatalf("ListTodos() error = %v, want nil", err)
			}
			if len(todos) != 2 {
				t.Errorf("ListTodos() returned %d todos, want 2", len(todos))
			}
		})
	})

	t.Run("with a file reader and writer", func(t *testing.T) {
		t.Run("creates new repository with empty file", func(t *testing.T) {
			tempDir := t.TempDir()
			tempFile := filepath.Join(tempDir, "test.todo.txt")

			reader := NewFileReader(tempFile)
			writer := NewFileWriter(tempFile)

			repo, err := NewRepository(reader, writer)
			if err != nil {
				t.Fatalf("NewRepository() error = %v, want nil", err)
			}

			todos, err := repo.ListTodos()
			if err != nil {
				t.Fatalf("ListTodos() error = %v, want nil", err)
			}
			if len(todos) != 0 {
				t.Errorf("ListTodos() returned %d todos, want 0", len(todos))
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

			repo, err := NewRepository(reader, writer)
			if err != nil {
				t.Fatalf("NewRepository() error = %v, want nil", err)
			}

			todos, err := repo.ListTodos()
			if err != nil {
				t.Fatalf("ListTodos() error = %v, want nil", err)
			}
			if len(todos) != 3 {
				t.Errorf("ListTodos() returned %d todos, want 3", len(todos))
			}
		})
	})
}

func TestRepository_Add(t *testing.T) {
	repo := setupTestRepository(t)

	t.Run("adds new todo", func(t *testing.T) {
		todo, err := repo.Add("Test todo 3")
		if err != nil {
			t.Errorf("Add() error = %v, want nil", err)
		}
		if todo.Text != "Test todo 3" {
			t.Errorf("Add() returned todo with text %s, want Test todo 3", todo.Text)
		}

		todos, err := repo.ListTodos()
		if err != nil {
			t.Errorf("ListTodos() error = %v, want nil", err)
		}
		if len(todos) != 4 {
			t.Errorf("expected 4 todos, got %d", len(todos))
		}
	})
}

func TestRepository_Remove(t *testing.T) {
	repo := setupTestRepository(t)

	t.Run("removes todo by index", func(t *testing.T) {
		todo, err := repo.Remove(0)
		if err != nil {
			t.Errorf("Remove() error = %v, want nil", err)
		}
		if todo.Text != "(A) test todo 1 +project2 @context1" {
			t.Errorf("Remove() returned todo with text %s, want (A) test todo 1 +project2 @context1", todo.Text)
		}

		todos, err := repo.ListTodos()
		if err != nil {
			t.Errorf("ListTodos() error = %v, want nil", err)
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

func TestRepository_ListTodos(t *testing.T) {
	repo := setupTestRepository(t)

	t.Run("returns all todos", func(t *testing.T) {
		todos, err := repo.ListTodos()
		if err != nil {
			t.Errorf("ListTodos() error = %v, want nil", err)
		}
		if len(todos) != 3 {
			t.Errorf("ListTodos() returned %d todos, want 3", len(todos))
		}
	})
}

func TestRepository_ListDone(t *testing.T) {
	repo := setupTestRepository(t)

	t.Run("returns only done todos", func(t *testing.T) {
		done, err := repo.ListDone()
		if err != nil {
			t.Errorf("ListDone() error = %v, want nil", err)
		}
		if len(done) != 1 {
			t.Errorf("ListDone() returned %d todos, want 1", len(done))
		}
		if !done[0].Done {
			t.Error("ListDone() returned non-done todo")
		}
	})
}

func TestRepository_ListProjects(t *testing.T) {
	repo := setupTestRepository(t)

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
	repo := setupTestRepository(t)

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

func TestRepository_Save(t *testing.T) {
	t.Run("saves todos to buffer", func(t *testing.T) {
		var buf bytes.Buffer
		reader := NewBufferReader(&buf)
		writer := NewBufferWriter(&buf)

		repo, err := NewRepository(reader, writer)
		if err != nil {
			t.Fatalf("NewRepository() error = %v, want nil", err)
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
