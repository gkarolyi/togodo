package todotxtlib

import (
	"os"
	"path/filepath"
	"testing"
)

var testTodos = []Todo{
	{Text: "Test todo 1"},
	{Text: "Test todo 2"},
}

func createTestRepository(tempDir string) Repository {
	reader := NewFileReader()
	writer := NewFileWriter()
	tempTodoTxt := filepath.Join(tempDir, "todo.txt")

	return &repository{
		todos:  testTodos,
		reader: reader,
		writer: writer,
		path:   tempTodoTxt,
	}
}

func TestNewRepository(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "todo.txt")

	// Test creating new repository
	repo, err := NewRepository(tempFile)
	if err != nil {
		t.Fatalf("NewRepository() error = %v, want nil", err)
	}
	if repo == nil {
		t.Fatal("NewRepository() returned nil repository")
	}
}

func TestRepository_Save(t *testing.T) {
	// Create repository with test todos
	tempDir := t.TempDir()
	repo := createTestRepository(tempDir)

	// Test that Save does not return an error
	if err := repo.Save(); err != nil {
		t.Errorf("Save() error = %v, want nil", err)
	}

	// Test that the file was created
	// Test that the file exists and is not empty
	todoFile := filepath.Join(tempDir, "todo.txt")
	fileInfo, err := os.Stat(todoFile)
	if err != nil {
		t.Errorf("todo.txt does not exist: %v", err)
	}
	if fileInfo.Size() == 0 {
		t.Error("todo.txt is empty")
	}
}

// func TestRepository_CRUDOperations(t *testing.T) {
// 	tempDir := t.TempDir()
// 	tempFile := filepath.Join(tempDir, "test.todo.txt")

// 	// Create repository with initial todos
// 	content := "Buy groceries\nCall mom\n"
// 	if err := os.WriteFile(tempFile, []byte(content), 0644); err != nil {
// 		t.Fatalf("Failed to create test file: %v", err)
// 	}

// 	repo, err := NewRepository(tempFile)
// 	if err != nil {
// 		t.Fatalf("NewRepository() error = %v", err)
// 	}

// 	// Test Add
// 	newTodo := Todo{Text: "Pay bills"}
// 	if err := repo.Add(newTodo); err != nil {
// 		t.Errorf("Add() error = %v, want nil", err)
// 	}

// 	// Test Find
// 	todos, err := repo.Find(Filter{Query: "bills"})
// 	if err != nil {
// 		t.Errorf("Find() error = %v, want nil", err)
// 	}
// 	if len(todos) != 1 || todos[0].Text != "Pay bills" {
// 		t.Errorf("Find() got = %v, want [Pay bills]", todos)
// 	}

// 	// Test Update
// 	updatedTodo := todos[0]
// 	updatedTodo.Text = "Pay bills tomorrow"
// 	if err := repo.Update(updatedTodo); err != nil {
// 		t.Errorf("Update() error = %v, want nil", err)
// 	}

// 	// Verify update
// 	todos, err = repo.Find(Filter{Query: "tomorrow"})
// 	if err != nil {
// 		t.Errorf("Find() after update error = %v, want nil", err)
// 	}
// 	if len(todos) != 1 || todos[0].Text != "Pay bills tomorrow" {
// 		t.Errorf("Find() after update got = %v, want [Pay bills tomorrow]", todos)
// 	}

// 	// Test Remove
// 	if err := repo.Remove(updatedTodo); err != nil {
// 		t.Errorf("Remove() error = %v, want nil", err)
// 	}

// 	// Verify remove
// 	todos, err = repo.Find(Filter{Query: "tomorrow"})
// 	if err != nil {
// 		t.Errorf("Find() after remove error = %v, want nil", err)
// 	}
// 	if len(todos) != 0 {
// 		t.Errorf("Find() after remove got = %v, want []", todos)
// 	}
// }

// func TestRepository_StatusOperations(t *testing.T) {
// 	tempDir := t.TempDir()
// 	tempFile := filepath.Join(tempDir, "test.todo.txt")

// 	// Create repository with initial todo
// 	content := "Buy groceries\n"
// 	if err := os.WriteFile(tempFile, []byte(content), 0644); err != nil {
// 		t.Fatalf("Failed to create test file: %v", err)
// 	}

// 	repo, err := NewRepository(tempFile)
// 	if err != nil {
// 		t.Fatalf("NewRepository() error = %v", err)
// 	}

// 	// Get the todo
// 	todos, err := repo.All()
// 	if err != nil {
// 		t.Fatalf("All() error = %v", err)
// 	}
// 	if len(todos) != 1 {
// 		t.Fatalf("All() got %d todos, want 1", len(todos))
// 	}
// 	todo := todos[0]

// 	// Test ToggleDone
// 	if err := repo.ToggleDone(todo); err != nil {
// 		t.Errorf("ToggleDone() error = %v, want nil", err)
// 	}

// 	// Verify toggle
// 	todos, err = repo.All()
// 	if err != nil {
// 		t.Errorf("All() after toggle error = %v", err)
// 	}
// 	if !todos[0].Done {
// 		t.Error("ToggleDone() did not mark todo as done")
// 	}

// 	// Test SetPriority
// 	if err := repo.SetPriority(todo, "A"); err != nil {
// 		t.Errorf("SetPriority() error = %v, want nil", err)
// 	}

// 	// Verify priority
// 	todos, err = repo.All()
// 	if err != nil {
// 		t.Errorf("All() after priority error = %v", err)
// 	}
// 	if todos[0].Priority != "A" {
// 		t.Errorf("SetPriority() got priority = %v, want A", todos[0].Priority)
// 	}
// }
