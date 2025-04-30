package todotxtlib

import (
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

func setupTestRepository(tb testing.TB) (func(tb testing.TB), Repository) {
	tempDir := tb.TempDir()
	reader := NewFileReader()
	writer := NewFileWriter()

	repo := &repository{
		todos:  createTestTodos(),
		reader: reader,
		writer: writer,
		path:   filepath.Join(tempDir, "todo.txt"),
	}

	teardown := func(tb testing.TB) {
		os.RemoveAll(tempDir)
	}

	return teardown, repo
}

func TestNewRepository(t *testing.T) {
	t.Run("creates new repository with empty file", func(t *testing.T) {
		tempDir := t.TempDir()
		tempFile := filepath.Join(tempDir, "todo.txt")

		repo, err := NewRepository(tempFile)
		if err != nil {
			t.Fatalf("NewRepository() error = %v, want nil", err)
		}
		if repo == nil {
			t.Fatal("NewRepository() returned nil repository")
		}
	})

	t.Run("creates new repository with existing file", func(t *testing.T) {
		tempDir := t.TempDir()
		tempFile := filepath.Join(tempDir, "todo.txt")

		// Create an initial todo file
		if err := os.WriteFile(tempFile, []byte("Test todo\n"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		repo, err := NewRepository(tempFile)
		if err != nil {
			t.Fatalf("NewRepository() error = %v, want nil", err)
		}
		if repo == nil {
			t.Fatal("NewRepository() returned nil repository")
		}

		todos, err := repo.ListTodos()
		if err != nil {
			t.Fatalf("ListTodos() error = %v, want nil", err)
		}
		if len(todos) != 1 {
			t.Errorf("ListTodos() returned %d todos, want 1", len(todos))
		}
	})
}

func TestRepository_Add(t *testing.T) {
	teardown, repo := setupTestRepository(t)
	defer teardown(t)

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
	teardown, repo := setupTestRepository(t)
	defer teardown(t)

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
	teardown, repo := setupTestRepository(t)
	defer teardown(t)

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
	teardown, repo := setupTestRepository(t)
	defer teardown(t)

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
	teardown, repo := setupTestRepository(t)
	defer teardown(t)

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
	teardown, repo := setupTestRepository(t)
	defer teardown(t)

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
	t.Run("saves todos to file", func(t *testing.T) {
		tempDir := t.TempDir()
		tempFile := filepath.Join(tempDir, "todo.txt")

		repo, err := NewRepository(tempFile)
		if err != nil {
			t.Fatalf("NewRepository() error = %v, want nil", err)
		}

		// Add some todos
		repo.Add("Test todo 1")
		repo.Add("Test todo 2")

		if err := repo.Save(); err != nil {
			t.Errorf("Save() error = %v, want nil", err)
		}

		// Verify file contents
		content, err := os.ReadFile(tempFile)
		if err != nil {
			t.Fatalf("Failed to read saved file: %v", err)
		}

		expected := "Test todo 1\nTest todo 2\n"
		if string(content) != expected {
			t.Errorf("Save() wrote %q, want %q", string(content), expected)
		}
	})
}
